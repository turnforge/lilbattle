package services

import (
	"context"
	"fmt"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

type GamesServiceImpl interface {
	v1.GamesServiceServer
	SaveGame(game *v1.Game, state *v1.GameState, history *v1.GameMoveHistory) error
	GetRuntimeGame(gameId string) (*weewar.Game, error)
}

type BaseGamesServiceImpl struct {
	v1.UnimplementedGamesServiceServer
	WorldsService v1.WorldsServiceServer
	Self          GamesServiceImpl // The actual implementation
}

type BaseWorldsServiceImpl struct {
	v1.UnimplementedWorldsServiceServer
}

// ProcessMoves processes moves for an existing game on the wasm side.
// Unlike the service side games service - it wont persist any changes - it only will return the diffs.
func (s *BaseGamesServiceImpl) ProcessMoves(ctx context.Context, req *v1.ProcessMovesRequest) (resp *v1.ProcessMovesResponse, err error) {
	if len(req.Moves) == 0 {
		return nil, fmt.Errorf("at least one move is required")
	}

	gameresp, err := s.GetGame(ctx, &v1.GetGameRequest{Id: req.GameId})
	if err != nil || gameresp.Game == nil {
		return nil, err
	}
	if gameresp.State == nil {
		panic("Game state cannot be nil")
	}
	if gameresp.History == nil {
		panic("Game history cannot cannot be nil")
	}

	// Get the runtime game corresponding to this game Id, we can create it on the fly
	// or we can cache it somewhere, or in the case of wasm just have a singleton
	rtGame, err := s.Self.GetRuntimeGame(gameresp.Game.Id)
	if err != nil {
		return nil, err
	}

	// Get the moves validted by the move processor, it is upto the move processor
	// to decide how "transactional" it wants to be - ie fail after  N moves,
	// success only if all moves succeeds etc.  Note that at this point the game
	// state has not changed and neither has the Runtime Game object.  Both the
	// GameState and the Runtime Game are checkpointed at before the moves started
	var dmp weewar.DefaultMoveProcessor
	results, err := dmp.ProcessMoves(rtGame, req.Moves)
	if err != nil {
		return nil, err
	}
	resp = &v1.ProcessMovesResponse{
		MoveResults: results,
	}

	// Now that we have the results, we want to update our gamestate by applying the
	// results - this would also set the next "checkoint" to after the reuslts.
	// It is upto the storage to see how the runtime game is also updated.  For example
	// a storage that persists the gameState may just not do anythign and let it be
	// reconstructed on the next load
	s.ApplyChangeResults(rtGame, gameresp.Game, gameresp.State, gameresp.History)

	s.Self.SaveGame(gameresp.Game, gameresp.State, gameresp.History)

	return resp, err
}

func (b *BaseGamesServiceImpl) ApplyChangeResults(rtGame *weewar.Game, game *v1.Game, state *v1.GameState, history *v1.GameMoveHistory) error {
	return nil
}
