package tests

import (
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	"github.com/turnforge/lilbattle/lib"
)

// newVictoryTestGame builds a minimal N-player game on a 5x5 grass map. Each
// player gets a starting soldier when seedUnits[i] is true. Callers can
// mutate the world before driving the game.
//
// The shared NewTestGameSetup helper in controller_test.go is hardcoded to
// 2 players; this sibling covers the 3+ player cases the multi-player
// elimination tests will need as the suite grows.
func newVictoryTestGame(t *testing.T, numPlayers int, seedUnits []bool) *TestGameSetup {
	t.Helper()
	if len(seedUnits) != numPlayers {
		t.Fatalf("seedUnits length (%d) must match numPlayers (%d)", len(seedUnits), numPlayers)
	}

	rulesEngine, err := LoadRulesEngineFromFile(RULES_DATA_FILE, DAMAGE_DATA_FILE)
	if err != nil {
		t.Fatalf("Failed to load rules engine: %v", err)
	}

	world := lib.NewWorld("test", nil)

	players := make([]*v1.GamePlayer, 0, numPlayers)
	playerStates := make(map[int32]*v1.PlayerState, numPlayers)
	for i := 1; i <= numPlayers; i++ {
		players = append(players, &v1.GamePlayer{PlayerId: int32(i), StartingCoins: 1000})
		playerStates[int32(i)] = &v1.PlayerState{Coins: 1000, IsActive: true}
	}

	game := &v1.Game{
		Id:   "test-game",
		Name: "Test Game",
		Config: &v1.GameConfiguration{
			Players: players,
		},
	}
	gameState := &v1.GameState{
		CurrentPlayer: 1,
		TurnCounter:   1,
		PlayerStates:  playerStates,
	}
	rtGame := lib.NewGame(game, gameState, world, rulesEngine, 12345)

	setup := &TestGameSetup{World: world, Game: rtGame, RulesEngine: rulesEngine}
	setup.AddGrassTiles(-2, 2, -2, 2)

	// One soldier per opted-in player, spread out so they don't share a hex.
	startCoords := [][2]int{{-2, 0}, {2, 0}, {0, -2}, {0, 2}, {-2, 2}, {2, -2}}
	for i, seed := range seedUnits {
		if !seed {
			continue
		}
		c := startCoords[i]
		setup.AddUnit(c[0], c[1], int32(i+1), UnitTypeSoldier)
	}
	return setup
}

// TestVictoryConditions_NoWinner_BothActive pins the no-premature-win path:
// while every player still has a unit, EndTurn must not declare a winner.
func TestVictoryConditions_NoWinner_BothActive(t *testing.T) {
	setup := newVictoryTestGame(t, 2, []bool{true, true})

	_, err := setup.Game.EndTurn()
	if err != nil {
		t.Fatalf("EndTurn failed: %v", err)
	}

	if setup.Game.GameState.Finished {
		t.Errorf("expected Finished=false while both players have units; got true")
	}
	if got := setup.Game.GameState.WinningPlayer; got != 0 {
		t.Errorf("expected WinningPlayer=0 while game in progress; got %d", got)
	}
	if got := setup.Game.GameState.Status; got == v1.GameStatus_GAME_STATUS_ENDED {
		t.Errorf("expected Status != ENDED while game in progress; got %v", got)
	}
}
