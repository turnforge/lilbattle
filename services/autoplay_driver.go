package services

import (
	"context"
	"fmt"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// AutoplayContextProvider returns the context to use for the next read/write
// against the GamesService. It receives the current player whose turn it is
// so callers can switch user identity per turn — needed for tests (where one
// caller drives all players) and for future dev-mode "swap as current player"
// behaviour. Pass nil to use the caller's ctx unchanged for every call.
type AutoplayContextProvider func(ctx context.Context, currentPlayer int32) context.Context

// RunAutoplayRequest captures every parameter RunAutoplay needs. Follows the
// project's (ctx, *Request) → (*Response, error) convention so future
// additions (policy name, RL-format flag, observer hooks) extend the request
// rather than the function signature.
type RunAutoplayRequest struct {
	// Svc is the GamesService used both by the presenter (read-only via
	// GetGame / GetOptionsAt) and by the driver (write via ProcessMoves) so
	// a single consistent backend handles the whole loop. Pass a singleton
	// or FSGamesService for in-process tests; pass the CLI's wired service
	// in production.
	Svc GamesService

	// GameID identifies the target game in the service.
	GameID string

	// Seed pins the picker's RNG for determinism. Same (initial state, seed,
	// service) tuple must produce identical action sequences.
	Seed int64

	// MaxTurns is a runaway safety cap. When the loop observes this many
	// turn cycles without a Finished flip, it returns an error so a stuck
	// game can't burn CI time or RL collection budget. Pass <= 0 to use the
	// default of 200.
	MaxTurns int

	// MaxMoves bounds the number of ProcessMoves calls before the loop
	// returns normally (no error). Pass 0 or negative for "run to
	// completion." Useful for short deterministic test runs, replays, and
	// step-debugging.
	MaxMoves int

	// CtxFor lets callers swap the request context per turn (e.g. to inject
	// the current player's identity). Pass nil to use the base ctx
	// unchanged for every call.
	CtxFor AutoplayContextProvider
}

// RunAutoplayResponse is the terminal state captured at the end of a
// RunAutoplay invocation — either a finished game (Finished=true,
// WinningPlayer set) or the state at the bound if MaxMoves stopped the loop
// before the game ended.
type RunAutoplayResponse struct {
	FinalState     *v1.GameState
	TurnsObserved  int
	HitSafetyCap   bool
	ActionsApplied int
}

// RunAutoplay drives the game to a terminal state via the presenter's
// NextMove policy. On each iteration it asks NextMove for an action; nil
// is interpreted as "end turn." Loops until GameState.Finished flips true,
// MaxTurns is exceeded, or MaxMoves is reached (whichever first).
//
// Errors propagate without retry per CONSTRAINTS.md ("no defensive error
// handling") — autoplay halts loudly on the first failure rather than
// silently skipping moves and continuing on stale state.
func RunAutoplay(ctx context.Context, req *RunAutoplayRequest) (*RunAutoplayResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("autoplay: nil request")
	}
	maxTurns := req.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 200
	}

	presenter := NewGameViewPresenter()
	presenter.GamesService = req.Svc
	presenter.SetSeed(req.Seed)

	resp := &RunAutoplayResponse{}
	var lastTurnCounter int32 = -1

	for {
		// Read uses the caller's base ctx; per-player swap only applies to
		// the ProcessMoves write below.
		getResp, err := req.Svc.GetGame(ctx, &v1.GetGameRequest{Id: req.GameID})
		if err != nil {
			return resp, fmt.Errorf("autoplay: reload state: %w", err)
		}
		resp.FinalState = getResp.State
		if getResp.State.Finished {
			return resp, nil
		}
		if getResp.State.TurnCounter != lastTurnCounter {
			lastTurnCounter = getResp.State.TurnCounter
			resp.TurnsObserved++
			if resp.TurnsObserved > maxTurns {
				resp.HitSafetyCap = true
				return resp, fmt.Errorf("autoplay: safety cap hit (max-turns=%d); aborting at turn %d to avoid runaway",
					maxTurns, getResp.State.TurnCounter)
			}
		}

		opt, err := presenter.NextMove(ctx, req.GameID)
		if err != nil {
			return resp, fmt.Errorf("autoplay: NextMove: %w", err)
		}

		var move *v1.GameMove
		if opt == nil {
			move = &v1.GameMove{MoveType: &v1.GameMove_EndTurn{EndTurn: &v1.EndTurnAction{}}}
		} else {
			move = OptionToMove(opt)
			if move == nil {
				return resp, fmt.Errorf("autoplay: unsupported GameOption variant: %T", opt.OptionType)
			}
		}

		moveCtx := ctx
		if req.CtxFor != nil {
			moveCtx = req.CtxFor(ctx, getResp.State.CurrentPlayer)
		}
		_, err = req.Svc.ProcessMoves(moveCtx, &v1.ProcessMovesRequest{
			GameId: req.GameID,
			Moves:  []*v1.GameMove{move},
		})
		if err != nil {
			return resp, fmt.Errorf("autoplay: ProcessMoves: %w", err)
		}
		resp.ActionsApplied++

		if req.MaxMoves > 0 && resp.ActionsApplied >= req.MaxMoves {
			// Re-read the post-move state so callers can inspect the
			// terminal world (which player's turn it is, whether the
			// final action ended the game, etc.).
			finalResp, err := req.Svc.GetGame(ctx, &v1.GetGameRequest{Id: req.GameID})
			if err != nil {
				return resp, fmt.Errorf("autoplay: reload final state: %w", err)
			}
			resp.FinalState = finalResp.State
			return resp, nil
		}
	}
}

// OptionToMove wraps a GameOption variant in the parallel GameMove variant.
// The underlying action protos (MoveUnitAction, AttackUnitAction, etc.) are
// shared between the two parents — only the wrapper changes.
func OptionToMove(opt *v1.GameOption) *v1.GameMove {
	switch o := opt.OptionType.(type) {
	case *v1.GameOption_Move:
		return &v1.GameMove{MoveType: &v1.GameMove_MoveUnit{MoveUnit: o.Move}}
	case *v1.GameOption_Attack:
		return &v1.GameMove{MoveType: &v1.GameMove_AttackUnit{AttackUnit: o.Attack}}
	case *v1.GameOption_Build:
		return &v1.GameMove{MoveType: &v1.GameMove_BuildUnit{BuildUnit: o.Build}}
	case *v1.GameOption_Capture:
		return &v1.GameMove{MoveType: &v1.GameMove_CaptureBuilding{CaptureBuilding: o.Capture}}
	case *v1.GameOption_Heal:
		return &v1.GameMove{MoveType: &v1.GameMove_HealUnit{HealUnit: o.Heal}}
	case *v1.GameOption_EndTurn:
		return &v1.GameMove{MoveType: &v1.GameMove_EndTurn{EndTurn: o.EndTurn}}
	}
	return nil
}
