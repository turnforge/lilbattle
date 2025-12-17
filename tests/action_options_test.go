package tests

import (
	"context"
	"testing"

	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
	"github.com/turnforge/weewar/services/singleton"
)

// TestGetOptionsAt_AttackAfterMove tests that a unit with action_order ["move", "attack"]
// can still attack after moving, even if it hasn't exhausted all movement points.
// This is a common scenario where the player moves a unit and then wants to attack.
func TestGetOptionsAt_AttackAfterMove(t *testing.T) {
	// Create a scenario:
	// - Anti-aircraft (type 6) at position 1,0 with action_order ["move", "attack"]
	// - Enemy unit at position 2,0 (within attack range)
	// - Unit has already moved this turn (progression_step=0, but no more move tiles available)
	// - Unit should be able to attack the enemy

	game := &v1.Game{
		Id:   "test-game",
		Name: "Test Game",
		Config: &v1.GameConfiguration{
			Settings: &v1.GameSettings{},
		},
	}

	gameState := &v1.GameState{
		CurrentPlayer: 1,
		TurnCounter:   1,
		PlayerStates: map[int32]*v1.PlayerState{
			1: {Coins: 1000},
			2: {Coins: 1000},
		},
		WorldData: &v1.WorldData{
			TilesMap: map[string]*v1.Tile{
				"0,0": {Q: 0, R: 0, TileType: 5}, // Grass
				"1,0": {Q: 1, R: 0, TileType: 5}, // Grass - our unit here
				"2,0": {Q: 2, R: 0, TileType: 5}, // Grass - enemy here
			},
			UnitsMap: map[string]*v1.Unit{
				"1,0": {
					Q:                1,
					R:                0,
					Player:           1,
					UnitType:         6, // Anti-aircraft - action_order: ["move", "attack"]
					AvailableHealth:  10,
					DistanceLeft:     0.5,            // Some movement left, but not enough to move anywhere useful
					ProgressionStep:  0,              // Still at move step
					LastToppedupTurn: 1,
					Shortcut:         "A1",
				},
				"2,0": {
					Q:                2,
					R:                0,
					Player:           2,  // Enemy
					UnitType:         1,  // Trooper
					AvailableHealth:  10,
					DistanceLeft:     3,
					LastToppedupTurn: 1,
					Shortcut:         "B1",
				},
			},
		},
	}

	gamesService := singleton.NewSingletonGamesService()
	gamesService.SingletonGame = game
	gamesService.SingletonGameState = gameState
	gamesService.SingletonGameMoveHistory = &v1.GameMoveHistory{}
	gamesService.Self = gamesService

	ctx := context.Background()

	// Get options for our unit at 1,0
	resp, err := gamesService.GetOptionsAt(ctx, &v1.GetOptionsAtRequest{
		GameId: "test-game",
		Q:      1,
		R:      0,
	})

	if err != nil {
		t.Fatalf("GetOptionsAt failed: %v", err)
	}

	// Count attack options - there should be at least one (attack enemy at 2,0)
	attackOptionCount := 0
	for _, opt := range resp.Options {
		if opt.GetAttack() != nil {
			attackOptionCount++
			t.Logf("Found attack option: attack at (%d,%d)", opt.GetAttack().DefenderQ, opt.GetAttack().DefenderR)
		}
		if opt.GetMove() != nil {
			t.Logf("Found move option: move to (%d,%d)", opt.GetMove().ToQ, opt.GetMove().ToR)
		}
	}

	if attackOptionCount == 0 {
		t.Errorf("Expected attack options for unit that can still attack after move, got 0 attack options")
		t.Logf("Total options: %d", len(resp.Options))
	}
}

// TestGetOptionsAt_AfterFullMove tests that after a unit exhausts all movement,
// it should show attack options if action_order allows attack after move.
func TestGetOptionsAt_AfterFullMove(t *testing.T) {
	game := &v1.Game{
		Id:   "test-game",
		Name: "Test Game",
		Config: &v1.GameConfiguration{
			Settings: &v1.GameSettings{},
		},
	}

	gameState := &v1.GameState{
		CurrentPlayer: 1,
		TurnCounter:   1,
		PlayerStates: map[int32]*v1.PlayerState{
			1: {Coins: 1000},
			2: {Coins: 1000},
		},
		WorldData: &v1.WorldData{
			TilesMap: map[string]*v1.Tile{
				"0,0": {Q: 0, R: 0, TileType: 5},
				"1,0": {Q: 1, R: 0, TileType: 5},
				"2,0": {Q: 2, R: 0, TileType: 5},
			},
			UnitsMap: map[string]*v1.Unit{
				"1,0": {
					Q:                1,
					R:                0,
					Player:           1,
					UnitType:         6, // Anti-aircraft - action_order: ["move", "attack"]
					AvailableHealth:  10,
					DistanceLeft:     0,               // No movement left
					ProgressionStep:  1,               // Advanced to attack step
					LastToppedupTurn: 1,
					Shortcut:         "A1",
				},
				"2,0": {
					Q:                2,
					R:                0,
					Player:           2,
					UnitType:         1,
					AvailableHealth:  10,
					DistanceLeft:     3,
					LastToppedupTurn: 1,
					Shortcut:         "B1",
				},
			},
		},
	}

	gamesService := singleton.NewSingletonGamesService()
	gamesService.SingletonGame = game
	gamesService.SingletonGameState = gameState
	gamesService.SingletonGameMoveHistory = &v1.GameMoveHistory{}
	gamesService.Self = gamesService

	ctx := context.Background()

	resp, err := gamesService.GetOptionsAt(ctx, &v1.GetOptionsAtRequest{
		GameId: "test-game",
		Q:      1,
		R:      0,
	})

	if err != nil {
		t.Fatalf("GetOptionsAt failed: %v", err)
	}

	attackOptionCount := 0
	for _, opt := range resp.Options {
		if opt.GetAttack() != nil {
			attackOptionCount++
		}
	}

	if attackOptionCount == 0 {
		t.Errorf("Expected attack options after full move (progression_step=1), got 0")
	}

	t.Logf("Got %d total options, %d attack options", len(resp.Options), attackOptionCount)
}

// TestGetOptionsAt_NoMoveOptionsAutoAdvance tests that when a unit has no valid move options
// (even with some DistanceLeft), it should auto-advance to the next action in action_order.
func TestGetOptionsAt_NoMoveOptionsAutoAdvance(t *testing.T) {
	// Scenario: Unit is surrounded by occupied tiles or impassable terrain
	// Even though DistanceLeft > 0 and ProgressionStep = 0, there are no move options
	// The system should auto-advance to allow attack

	game := &v1.Game{
		Id:   "test-game",
		Name: "Test Game",
		Config: &v1.GameConfiguration{
			Settings: &v1.GameSettings{},
		},
	}

	gameState := &v1.GameState{
		CurrentPlayer: 1,
		TurnCounter:   1,
		PlayerStates: map[int32]*v1.PlayerState{
			1: {Coins: 1000},
			2: {Coins: 1000},
		},
		WorldData: &v1.WorldData{
			TilesMap: map[string]*v1.Tile{
				"0,0": {Q: 0, R: 0, TileType: 5},   // Grass - friendly unit
				"1,0": {Q: 1, R: 0, TileType: 5},   // Grass - our test unit
				"-1,0": {Q: -1, R: 0, TileType: 5}, // Grass - friendly unit blocking
				"0,-1": {Q: 0, R: -1, TileType: 5}, // Grass - friendly unit blocking
				"1,-1": {Q: 1, R: -1, TileType: 5}, // Grass - friendly unit blocking
				"0,1":  {Q: 0, R: 1, TileType: 5},  // Grass - enemy unit here (can attack)
				"-1,1": {Q: -1, R: 1, TileType: 5}, // Grass - friendly unit blocking
			},
			UnitsMap: map[string]*v1.Unit{
				// Our unit in the center
				"1,0": {
					Q:                1,
					R:                0,
					Player:           1,
					UnitType:         6, // Anti-aircraft
					AvailableHealth:  10,
					DistanceLeft:     3, // Has full movement points
					ProgressionStep:  0, // At move step
					LastToppedupTurn: 1,
					Shortcut:         "A1",
				},
				// Surrounding friendly units (blocking movement)
				"0,0": {
					Q: 0, R: 0, Player: 1, UnitType: 1, AvailableHealth: 10, DistanceLeft: 3, LastToppedupTurn: 1,
				},
				"-1,0": {
					Q: -1, R: 0, Player: 1, UnitType: 1, AvailableHealth: 10, DistanceLeft: 3, LastToppedupTurn: 1,
				},
				"0,-1": {
					Q: 0, R: -1, Player: 1, UnitType: 1, AvailableHealth: 10, DistanceLeft: 3, LastToppedupTurn: 1,
				},
				"1,-1": {
					Q: 1, R: -1, Player: 1, UnitType: 1, AvailableHealth: 10, DistanceLeft: 3, LastToppedupTurn: 1,
				},
				"-1,1": {
					Q: -1, R: 1, Player: 1, UnitType: 1, AvailableHealth: 10, DistanceLeft: 3, LastToppedupTurn: 1,
				},
				// Enemy unit within attack range
				"0,1": {
					Q:                0,
					R:                1,
					Player:           2,
					UnitType:         1,
					AvailableHealth:  10,
					DistanceLeft:     3,
					LastToppedupTurn: 1,
					Shortcut:         "B1",
				},
			},
		},
	}

	gamesService := singleton.NewSingletonGamesService()
	gamesService.SingletonGame = game
	gamesService.SingletonGameState = gameState
	gamesService.SingletonGameMoveHistory = &v1.GameMoveHistory{}
	gamesService.Self = gamesService

	ctx := context.Background()

	resp, err := gamesService.GetOptionsAt(ctx, &v1.GetOptionsAtRequest{
		GameId: "test-game",
		Q:      1,
		R:      0,
	})

	if err != nil {
		t.Fatalf("GetOptionsAt failed: %v", err)
	}

	moveOptionCount := 0
	attackOptionCount := 0
	for _, opt := range resp.Options {
		if opt.GetMove() != nil {
			moveOptionCount++
		}
		if opt.GetAttack() != nil {
			attackOptionCount++
		}
	}

	// Should have no move options (surrounded)
	// But SHOULD have attack options (enemy in range)
	t.Logf("Move options: %d, Attack options: %d, Total: %d", moveOptionCount, attackOptionCount, len(resp.Options))

	if attackOptionCount == 0 {
		t.Errorf("Expected attack options when surrounded (no move options available), got 0")
	}
}
