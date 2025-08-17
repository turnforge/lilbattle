package weewar

import v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"

// createTestUnit creates a test unit with given parameters
func CreateTestUnit(q, r int, player, unitType int) *v1.Unit {
	return &v1.Unit{
		Q:               int32(q),
		R:               int32(r),
		Player:          int32(player),
		UnitType:        int32(unitType),
		AvailableHealth: 100,
		DistanceLeft:    3,
		TurnCounter:     1,
	}
}
