package services

import (
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

// ProtoToRuntimeGame converts protobuf game/state to runtime game
// This is WeeWar-specific and doesn't belong in TurnEngine
func ProtoToRuntimeGame(game *v1.Game, gameState *v1.GameState) *weewar.Game {
	// Create the runtime game from the protobuf data
	world := weewar.NewWorld(game.Name, gameState.WorldData)

	// Create the runtime game with loaded default rules engine
	rulesEngine := weewar.DefaultRulesEngine() // Use loaded default rules engine

	// Use NewGameFromState instead of NewGame to preserve unit stats
	return weewar.NewGame(game, gameState, world, rulesEngine, 12345) // Default seed
}
