package services

import (
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
	pj "google.golang.org/protobuf/encoding/protojson"
)

type WasmGamesServiceImpl struct {
	BaseGamesServiceImpl
	SingletonGame            *v1.Game
	SingletonGameState       *v1.GameState
	SingletonGameMoveHistory *v1.GameMoveHistory
	SingletonWorld           *v1.World
	SingletonWorldData       *v1.WorldData

	RuntimeGame *weewar.Game
}

// NOTE - ONly API really needed here are "getters" and "move processors" so no Creations, Deletions, Listing or even
// GetGame needed - GetGame data is set when we create this
func NewWasmGamesServiceImpl(
	gameBytes []byte,
	gameStateBytes []byte,
	gameMoveHistoryBytes []byte,
	worldBytes []byte,
	worldDataBytes []byte,
) *WasmGamesServiceImpl {
	out := &WasmGamesServiceImpl{
		BaseGamesServiceImpl: BaseGamesServiceImpl{
			// WorldsService: SingletonWorldsService
		},
		SingletonGame:            &v1.Game{},
		SingletonGameState:       &v1.GameState{},
		SingletonGameMoveHistory: &v1.GameMoveHistory{},
		SingletonWorld:           &v1.World{},
		SingletonWorldData:       &v1.WorldData{},
	}
	out.Self = out

	// Now load data from the bytes
	if err := pj.Unmarshal(gameBytes, out.SingletonGame); err != nil {
		panic(err)
	}
	if err := pj.Unmarshal(gameStateBytes, out.SingletonGameState); err != nil {
		panic(err)
	}
	if err := pj.Unmarshal(gameMoveHistoryBytes, out.SingletonGameMoveHistory); err != nil {
		panic(err)
	}
	if err := pj.Unmarshal(worldBytes, out.SingletonWorld); err != nil {
		panic(err)
	}
	if err := pj.Unmarshal(worldDataBytes, out.SingletonWorldData); err != nil {
		panic(err)
	}
	return out
}

func (w *WasmGamesServiceImpl) GetRuntimeGame(gameId string) (*weewar.Game, error) {
	return w.RuntimeGame, nil
}

func (w *WasmGamesServiceImpl) SaveGame(game *v1.Game, state *v1.GameState, history *v1.GameMoveHistory) error {
	return nil
}
