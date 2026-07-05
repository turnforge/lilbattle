package lib

import (
	"testing"
)

// Issue 156 changes the victory rule from "last player with units" to
// "last non-eliminated player," where elimination means BOTH zero units
// AND zero owned bases (build-capable tiles). Bases produce units given
// coins; a lone unit can capture bases. A player with either resource
// still has a recovery path and stays in the game.

// TestVictory_TwoLivePlayers_NoWinner pins the baseline: two players
// with units and bases each, no winner should be declared.
func TestVictory_TwoLivePlayers_NoWinner(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		tile(3, 0, TileTypeLandBase, 2).
		unit(0, 1, 1, UnitTypeSoldier).
		unit(3, 1, 2, UnitTypeSoldier).
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if hasWinner {
		t.Errorf("expected no winner with both players alive; got winner=%d", winner)
	}
}

// TestVictory_TotalElimination_LastPlayerWins is the unchanged happy
// path — one player has neither units nor bases, the other has both.
// Works under both the old and new rule.
func TestVictory_TotalElimination_LastPlayerWins(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		unit(0, 1, 1, UnitTypeSoldier).
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if !hasWinner {
		t.Fatal("expected a winner when only player 1 has any resources")
	}
	if winner != 1 {
		t.Errorf("winner = %d, want 1", winner)
	}
}

// TestVictory_UnitsButNoBases_StillAlive is the FIRST behavioral shift.
// Under today's rule, a player with 0 units loses. Under the new rule,
// the player is still alive as long as they have units OR bases — a
// unit-only player can capture a base and recover, so keep them in the
// game. Red-before-green: this test fails under the "last player with
// units" impl (declares player 1 winner despite player 2 still having
// units). Under the new impl, no winner because player 2 still has a
// unit (recovery: could capture a neutral base).
func TestVictory_UnitsButNoBases_StillAlive(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		unit(0, 1, 1, UnitTypeSoldier).
		unit(3, 3, 2, UnitTypeSoldier). // nomad, no base
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if hasWinner {
		t.Errorf("expected no winner while player 2 has a unit (can capture a base to recover); got winner=%d", winner)
	}
}

// TestVictory_BasesButNoUnits_StillAlive is the SECOND behavioral shift.
// Under today's rule, a player with 0 units loses immediately, even
// with 3 bases and 500 coins queued up to build new units. Under the
// new rule, they're still alive — bases + coins = a build-recovery
// path. Red-before-green: this test fails under the old impl.
func TestVictory_BasesButNoUnits_StillAlive(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		unit(0, 1, 1, UnitTypeSoldier).
		tile(3, 3, TileTypeLandBase, 2). // player 2 has a base
		coins(2, 500).                   // and coins to build with
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if hasWinner {
		t.Errorf("expected no winner while player 2 has a base (can build a unit to recover); got winner=%d", winner)
	}
}

// TestVictory_FullyEliminated_OpponentWins pins the correctness of the
// new elimination rule: neither units nor bases → out of the game.
func TestVictory_FullyEliminated_OpponentWins(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		unit(0, 1, 1, UnitTypeSoldier).
		// player 2 has nothing — no units, no bases
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if !hasWinner {
		t.Fatal("expected a winner when player 2 has zero units and zero bases")
	}
	if winner != 1 {
		t.Errorf("winner = %d, want 1", winner)
	}
}

// TestVictory_NonBaseTilesDoNotCountAsResources pins the terrain
// distinction: only tiles that CAN build units count as bases for
// recovery. A player owning a grass tile with no unit is still
// eliminated — grass doesn't produce anything.
func TestVictory_NonBaseTilesDoNotCountAsResources(t *testing.T) {
	g := newTestGameBuilder().
		tile(0, 0, TileTypeLandBase, 1).
		unit(0, 1, 1, UnitTypeSoldier).
		tile(3, 3, TileTypeGrass, 2). // owned but non-buildable
		build()

	winner, hasWinner := g.checkVictoryConditions()
	if !hasWinner {
		t.Fatal("expected player 1 to win; grass tile is not a base and provides no recovery path")
	}
	if winner != 1 {
		t.Errorf("winner = %d, want 1", winner)
	}
}
