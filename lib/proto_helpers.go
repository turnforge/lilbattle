package weewar

import (
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// Proto Unit helper methods
func UnitGetCoord(u *v1.Unit) AxialCoord {
	return AxialCoord{Q: int(u.Q), R: int(u.R)}
}

func UnitSetCoord(u *v1.Unit, coord AxialCoord) {
	u.Q = int32(coord.Q)
	u.R = int32(coord.R)
}

// Proto Tile helper methods
func TileGetCoord(t *v1.Tile) AxialCoord {
	return AxialCoord{Q: int(t.Q), R: int(t.R)}
}

func TileSetCoord(t *v1.Tile, coord AxialCoord) {
	t.Q = int32(coord.Q)
	t.R = int32(coord.R)
}

// Proto factory functions
func NewUnit(unitType, player int, coord AxialCoord) *v1.Unit {
	return &v1.Unit{
		Q:        int32(coord.Q),
		R:        int32(coord.R),
		Player:   int32(player),
		UnitType: int32(unitType),
		// Initialize runtime state
		AvailableHealth: 100, // Will be set by rules engine
		DistanceLeft:    0,   // Will be set by rules engine
		TurnCounter:     0,
	}
}

func NewTile(coord AxialCoord, tileType int) *v1.Tile {
	return &v1.Tile{
		Q:        int32(coord.Q),
		R:        int32(coord.R),
		TileType: int32(tileType),
		Player:   0, // Default to neutral
	}
}

// Helper functions to convert between int and int32 for proto fields
func ProtoInt32(val int) int32 {
	return int32(val)
}

func ProtoInt(val int32) int {
	return int(val)
}
