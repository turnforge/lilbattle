package tests

import (
	"context"
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	"github.com/turnforge/lilbattle/lib"
	"github.com/turnforge/lilbattle/services"
	"github.com/turnforge/lilbattle/services/singleton"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// autoplayCtxFor swaps the context's userID per current player so
// authz.CanSubmitMoves accepts every turn. The smoke test drives both
// players' turns through one process, so it must impersonate whichever
// player is acting next.
func autoplayCtxFor(_ context.Context, currentPlayer int32) context.Context {
	if currentPlayer == 1 {
		return ContextWithUserID(TestUserID)
	}
	return ContextWithUserID("player-2")
}


// newAutoplayTestGame builds a small 2-player skirmish suited to driving
// autoplay to completion quickly. Two soldiers on opposite sides of a
// 5x5 grass map; one will eventually annihilate the other.
func newAutoplayTestGame(t *testing.T) *singleton.SingletonGamesService {
	t.Helper()

	protoWorld := &v1.WorldData{}
	world := lib.NewWorld("test", protoWorld)
	for q := -2; q <= 2; q++ {
		for r := -2; r <= 2; r++ {
			coord := lib.AxialCoord{Q: q, R: r}
			world.AddTile(lib.NewTile(coord, lib.TileTypeGrass))
		}
	}

	world.AddUnit(&v1.Unit{
		Q: -2, R: 0, Player: 1, UnitType: UnitTypeSoldier,
		AvailableHealth: 10, DistanceLeft: 3,
	})
	world.AddUnit(&v1.Unit{
		Q: 2, R: 0, Player: 2, UnitType: UnitTypeSoldier,
		AvailableHealth: 10, DistanceLeft: 3,
	})

	rulesEngine, err := lib.LoadRulesEngineFromFile(RULES_DATA_FILE, DAMAGE_DATA_FILE)
	if err != nil {
		t.Fatalf("Failed to load rules engine: %v", err)
	}

	game := &v1.Game{
		Id:   "autoplay-test",
		Name: "Autoplay Test",
		Config: &v1.GameConfiguration{
			Players: []*v1.GamePlayer{
				{PlayerId: 1, UserId: TestUserID},
				{PlayerId: 2, UserId: "player-2"},
			},
		},
	}
	gameState := &v1.GameState{
		GameId:        game.Id,
		CurrentPlayer: 1,
		TurnCounter:   1,
		PlayerStates: map[int32]*v1.PlayerState{
			1: {Coins: 1000, IsActive: true},
			2: {Coins: 1000, IsActive: true},
		},
	}
	rtGame := lib.NewGame(game, gameState, world, rulesEngine, 12345)

	svc := singleton.NewSingletonGamesService()
	svc.SingletonGame = game
	svc.SingletonGameState = gameState
	svc.SingletonGameState.WorldData = convertRuntimeWorldToProto(world)
	svc.SingletonGameState.UpdatedAt = timestamppb.Now()
	svc.SingletonGameMoveHistory = &v1.GameMoveHistory{Groups: []*v1.GameMoveGroup{}}
	svc.RuntimeGame = rtGame
	svc.Self = svc

	return svc
}

// TestAutoplay_Smoke_RunsToCompletion pins the cardinal property: a fresh
// game driven by autoplay reaches a terminal Finished=true state with a
// declared winner, within the safety cap. If this fails, the picker is
// broken, NextMove is wrong, or the loop has stalled — all three are
// catastrophic and should block CI.
func TestAutoplay_Smoke_RunsToCompletion(t *testing.T) {
	svc := newAutoplayTestGame(t)
	ctx := AuthenticatedContext()

	result, err := services.RunAutoplay(ctx, &services.RunAutoplayRequest{
		Svc: svc, GameID: svc.SingletonGame.Id, Seed: 42,
		MaxTurns: 200, CtxFor: autoplayCtxFor,
	})
	if err != nil {
		t.Fatalf("RunAutoplay failed: %v (turns=%d, actions=%d, hitCap=%v)",
			err, result.TurnsObserved, result.ActionsApplied, result.HitSafetyCap)
	}

	if !result.FinalState.Finished {
		t.Errorf("expected Finished=true; got false (turns=%d, actions=%d)",
			result.TurnsObserved, result.ActionsApplied)
	}
	if result.FinalState.WinningPlayer == 0 {
		t.Errorf("expected WinningPlayer != 0; got 0")
	}
	if result.HitSafetyCap {
		t.Errorf("expected game to finish within safety cap; hit cap at turn %d",
			result.FinalState.TurnCounter)
	}
}

// TestAutoplay_Determinism_SameSeedSameOutcome pins seed determinism over
// the first MaxMoves steps. A bounded MaxMoves keeps the test cheap and
// isolates the failure mode to the picker / option-ordering layer —
// combat-RNG drift requires more actions to compound, and a small bound
// rarely reaches combat. Two runs from identical starting states with the
// same seed must produce identical action counts and identical post-bound
// state.
func TestAutoplay_Determinism_SameSeedSameOutcome(t *testing.T) {
	svc1 := newAutoplayTestGame(t)
	svc2 := newAutoplayTestGame(t)
	ctx := AuthenticatedContext()

	const moves = 10
	r1, err := services.RunAutoplay(ctx, &services.RunAutoplayRequest{
		Svc: svc1, GameID: svc1.SingletonGame.Id, Seed: 7,
		MaxTurns: 200, MaxMoves: moves, CtxFor: autoplayCtxFor,
	})
	if err != nil {
		t.Fatalf("run 1 failed: %v", err)
	}
	r2, err := services.RunAutoplay(ctx, &services.RunAutoplayRequest{
		Svc: svc2, GameID: svc2.SingletonGame.Id, Seed: 7,
		MaxTurns: 200, MaxMoves: moves, CtxFor: autoplayCtxFor,
	})
	if err != nil {
		t.Fatalf("run 2 failed: %v", err)
	}

	if r1.ActionsApplied != r2.ActionsApplied {
		t.Errorf("same-seed bounded runs disagree on action count: %d vs %d",
			r1.ActionsApplied, r2.ActionsApplied)
	}
	if r1.FinalState.CurrentPlayer != r2.FinalState.CurrentPlayer {
		t.Errorf("same-seed bounded runs disagree on CurrentPlayer after %d moves: %d vs %d",
			moves, r1.FinalState.CurrentPlayer, r2.FinalState.CurrentPlayer)
	}
	if r1.FinalState.TurnCounter != r2.FinalState.TurnCounter {
		t.Errorf("same-seed bounded runs disagree on turn after %d moves: %d vs %d",
			moves, r1.FinalState.TurnCounter, r2.FinalState.TurnCounter)
	}
}

// TestAutoplay_NextMove_OnlyValidOptions pins that NextMove never returns
// an option that would fail ProcessMoves validation. We exercise it by
// running a full autoplay loop and asserting ProcessMoves never errored
// for an action NextMove selected — RunAutoplay would have returned an
// error in that case.
//
// The test is implicit in the smoke test, but called out explicitly here
// so any future regression that silently drops invalid options (e.g. by
// catching the ProcessMoves error) breaks this test even if the smoke
// test still completes.
func TestAutoplay_NextMove_OnlyValidOptions(t *testing.T) {
	svc := newAutoplayTestGame(t)
	ctx := AuthenticatedContext()

	result, err := services.RunAutoplay(ctx, &services.RunAutoplayRequest{
		Svc: svc, GameID: svc.SingletonGame.Id, Seed: 13,
		MaxTurns: 200, CtxFor: autoplayCtxFor,
	})
	if err != nil {
		t.Fatalf("RunAutoplay failed (suggests an invalid option was picked or stale state): %v", err)
	}
	if result.ActionsApplied == 0 {
		t.Errorf("expected ActionsApplied > 0; got 0 (did NextMove always return nil?)")
	}
}

// TestNextMove_EmptyOptions_ReturnsNil pins the empty-pool contract on
// the presenter side. A game in Finished state has no actionable options
// (GetAllOptions returns empty), and NextMove must convert that to nil.
// Drivers depend on this nil signal to call EndTurn.
func TestNextMove_EmptyOptions_ReturnsNil(t *testing.T) {
	svc := newAutoplayTestGame(t)
	svc.SingletonGameState.Finished = true // simulate end-of-game
	ctx := AuthenticatedContext()

	presenter := services.NewGameViewPresenter()
	presenter.GamesService = svc

	opt, err := presenter.NextMove(ctx, svc.SingletonGame.Id)
	if err != nil {
		t.Fatalf("NextMove on finished game returned error: %v", err)
	}
	if opt != nil {
		t.Errorf("expected nil option on finished game; got %v", opt.OptionType)
	}
}
