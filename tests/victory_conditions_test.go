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

// TestVictoryConditions_TwoPlayer_LastStanding pins the baseline win path:
// P2 has no units when P1 ends turn, so checkVictoryConditions detects
// last-player-standing and the game ends with P1 as winner.
func TestVictoryConditions_TwoPlayer_LastStanding(t *testing.T) {
	setup := newVictoryTestGame(t, 2, []bool{true, false})

	_, err := setup.Game.EndTurn()
	if err != nil {
		t.Fatalf("EndTurn failed: %v", err)
	}

	if !setup.Game.GameState.Finished {
		t.Errorf("expected Finished=true after last-standing detected; got false")
	}
	if got := setup.Game.GameState.WinningPlayer; got != 1 {
		t.Errorf("expected WinningPlayer=1; got %d", got)
	}
	if got := setup.Game.GameState.Status; got != v1.GameStatus_GAME_STATUS_ENDED {
		t.Errorf("expected Status=ENDED; got %v", got)
	}
}

// TestVictoryConditions_ThreePlayer_MidGameElim_NoWinner pins two behaviors
// at once:
//
//  1. Multi-player elimination doesn't prematurely end the game — with P2
//     out but P1 and P3 still active, EndTurn must not declare a winner.
//  2. Eliminated players keep their turn slot in the current cycling impl —
//     EndTurn from P1 advances CurrentPlayer to 2 (not 3), even though P2
//     has nothing to act with. If this changes (e.g. skipping eliminated
//     players), the test forces the behavior change to be deliberate.
func TestVictoryConditions_ThreePlayer_MidGameElim_NoWinner(t *testing.T) {
	setup := newVictoryTestGame(t, 3, []bool{true, false, true})

	_, err := setup.Game.EndTurn()
	if err != nil {
		t.Fatalf("EndTurn failed: %v", err)
	}

	if setup.Game.GameState.Finished {
		t.Errorf("expected Finished=false with 2 of 3 players still alive; got true")
	}
	if got := setup.Game.CurrentPlayer; got != 2 {
		t.Errorf("expected CurrentPlayer=2 (eliminated players are not skipped in current impl); got %d", got)
	}
}

// TestVictoryConditions_ThreePlayer_LastSurvivor_Wins pins multi-player
// convergence: when only one of three players retains units, EndTurn
// declares that player the winner.
func TestVictoryConditions_ThreePlayer_LastSurvivor_Wins(t *testing.T) {
	setup := newVictoryTestGame(t, 3, []bool{true, false, false})

	_, err := setup.Game.EndTurn()
	if err != nil {
		t.Fatalf("EndTurn failed: %v", err)
	}

	if !setup.Game.GameState.Finished {
		t.Errorf("expected Finished=true with only P1 alive; got false")
	}
	if got := setup.Game.GameState.WinningPlayer; got != 1 {
		t.Errorf("expected WinningPlayer=1; got %d", got)
	}
}

// TestVictoryConditions_MidTurnKill_NotEndedUntilEndTurn pins the
// checked-only-on-EndTurn behavior. An attacker that eliminates the last
// enemy unit mid-turn does NOT see Finished flip until they end the turn.
// This is the deliberate single-check-point pattern in the current impl
// rather than a bug — pinning it ensures any future move to mid-turn
// victory detection is deliberate.
func TestVictoryConditions_MidTurnKill_NotEndedUntilEndTurn(t *testing.T) {
	setup := newVictoryTestGame(t, 2, []bool{true, true})

	// Replace the default defender at (2,0) with a 1-HP soldier adjacent
	// to a freshly-placed attacker so a single melee attack kills cleanly.
	defender := setup.World.UnitAt(lib.AxialCoord{Q: 2, R: 0})
	if defender == nil {
		t.Fatal("defender not found at (2,0)")
	}
	setup.World.RemoveUnit(defender)
	attacker := setup.World.UnitAt(lib.AxialCoord{Q: -2, R: 0})
	if attacker == nil {
		t.Fatal("attacker not found at (-2,0)")
	}
	setup.World.RemoveUnit(attacker)
	setup.World.AddUnit(&v1.Unit{
		Q: 0, R: 0, Player: 1, UnitType: UnitTypeSoldier,
		AvailableHealth: 10, DistanceLeft: 3,
	})
	setup.World.AddUnit(&v1.Unit{
		Q: 1, R: 0, Player: 2, UnitType: UnitTypeSoldier,
		AvailableHealth: 1, DistanceLeft: 3,
	})

	move := &v1.GameMove{
		MoveType: &v1.GameMove_AttackUnit{
			AttackUnit: &v1.AttackUnitAction{
				Attacker: &v1.Position{Q: 0, R: 0},
				Defender: &v1.Position{Q: 1, R: 0},
			},
		},
	}
	if err := setup.Game.ProcessMove(move); err != nil {
		t.Fatalf("ProcessMove (attack) failed: %v", err)
	}

	if setup.Game.GameState.Finished {
		t.Errorf("expected Finished=false mid-turn (victory checked only on EndTurn); got true")
	}
	if got := setup.Game.GameState.WinningPlayer; got != 0 {
		t.Errorf("expected WinningPlayer=0 mid-turn; got %d", got)
	}

	// Confirm the deferred check still fires correctly on EndTurn.
	if _, err := setup.Game.EndTurn(); err != nil {
		t.Fatalf("EndTurn after mid-turn kill failed: %v", err)
	}
	if !setup.Game.GameState.Finished {
		t.Errorf("expected Finished=true after EndTurn following mid-turn kill; got false")
	}
}

// TestVictoryConditions_MutualDestruction_NoWinner_BUG documents the
// current behavior when every player simultaneously runs out of units
// (e.g. mutual splash kill): checkVictoryConditions returns
// hasWinner=false because the "last player with units" check requires
// exactly one player to have units. The game then has no winner and no
// end state — a real bug, tracked at issue 157. Test stays skipped until
// that issue picks a tiebreaker rule.
func TestVictoryConditions_MutualDestruction_NoWinner_BUG(t *testing.T) {
	t.Skip("filed as issue 157 — mutual-destruction case needs tiebreaker decision before this can assert the right outcome")

	setup := newVictoryTestGame(t, 2, []bool{false, false})

	_, err := setup.Game.EndTurn()
	if err != nil {
		t.Fatalf("EndTurn failed: %v", err)
	}

	// Whatever rule issue 157 picks (draw / last-attacker / etc.), the
	// game must reach a terminal state. The current impl does not, so this
	// assertion will fail once unskipped — by design, pinning the fix.
	if !setup.Game.GameState.Finished {
		t.Errorf("expected Finished=true with mutual destruction; got false (game stuck)")
	}
}
