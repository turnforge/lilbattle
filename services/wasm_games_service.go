package services

import (
	"context"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

type WasmGamesServiceImpl struct {
	BaseWorldsServiceImpl
	SingletonGame  *v1.Game
	SingletonWorld *v1.World
	RuntimeGame    *weewar.Game
}

func NewWasmGamesServiceImpl(game *v1.Game, gameStateData []byte) *WASMStorage {
	out := &WASMStorage{
		SingletonGame: game,
	}
	// TODO - create the runtime game data based on the game state
	return out
}

// GetGame returns a specific game with complete data including moves
func (s *WasmGamesServiceImpl) GetGame(ctx context.Context, req *v1.GetGameRequest) (resp *v1.GetGameResponse, err error) {
	return &v1.GetGameResponse{Game: s.SingletonGame}, nil
}
