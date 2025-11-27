package services

import (
	"github.com/turnforge/weewar/lib"
	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
)

// Tile type constants for migration
const (
	TileTypeRoad          = 22
	TileTypeBridgeShallow = 18
	TileTypeBridgeRegular = 17
	TileTypeBridgeDeep    = 19
	TileTypePlains        = 5
	TileTypeWaterShallow  = 14
	TileTypeWaterRegular  = 10
	TileTypeWaterDeep     = 15
)

// MigrateWorldData converts old list-based WorldData to map-based storage.
// It also extracts crossings (roads, bridges) from tile types.
// This function is idempotent - calling it multiple times is safe.
func MigrateWorldData(wd *v1.WorldData) {
	if wd == nil {
		return
	}

	// Initialize maps if nil
	if wd.TilesMap == nil {
		wd.TilesMap = make(map[string]*v1.Tile)
	}
	if wd.UnitsMap == nil {
		wd.UnitsMap = make(map[string]*v1.Unit)
	}
	if wd.Crossings == nil {
		wd.Crossings = make(map[string]v1.CrossingType)
	}

	// Migrate tiles from deprecated list to map (if not already done)
	if len(wd.Tiles) > 0 && len(wd.TilesMap) == 0 {
		for _, tile := range wd.Tiles {
			key := lib.CoordKey(tile.Q, tile.R)
			wd.TilesMap[key] = tile
		}
	}

	// Migrate units from deprecated list to map (if not already done)
	if len(wd.Units) > 0 && len(wd.UnitsMap) == 0 {
		for _, unit := range wd.Units {
			key := lib.CoordKey(unit.Q, unit.R)
			wd.UnitsMap[key] = unit
		}
	}

	// Extract crossings from tile types
	extractCrossings(wd)

	// Clear deprecated lists after migration
	// Use empty slices instead of nil to ensure they're serialized as empty arrays
	// This ensures the old data is actually removed from storage on save
	wd.Tiles = []*v1.Tile{}
	wd.Units = []*v1.Unit{}
}

// extractCrossings extracts roads and bridges from tile types into the crossings map
// and updates the tile types to their underlying terrain.
func extractCrossings(wd *v1.WorldData) {
	for key, tile := range wd.TilesMap {
		switch tile.TileType {
		case TileTypeRoad:
			// Road -> Plains with road crossing
			wd.Crossings[key] = v1.CrossingType_CROSSING_TYPE_ROAD
			tile.TileType = TileTypePlains

		case TileTypeBridgeShallow:
			// Bridge over shallow water
			wd.Crossings[key] = v1.CrossingType_CROSSING_TYPE_BRIDGE
			tile.TileType = TileTypeWaterShallow

		case TileTypeBridgeRegular:
			// Bridge over regular water
			wd.Crossings[key] = v1.CrossingType_CROSSING_TYPE_BRIDGE
			tile.TileType = TileTypeWaterRegular

		case TileTypeBridgeDeep:
			// Bridge over deep water
			wd.Crossings[key] = v1.CrossingType_CROSSING_TYPE_BRIDGE
			tile.TileType = TileTypeWaterDeep
		}
	}
}

// GetCrossingType returns the crossing type at the given coordinates
func GetCrossingType(wd *v1.WorldData, q, r int32) v1.CrossingType {
	if wd == nil || wd.Crossings == nil {
		return v1.CrossingType_CROSSING_TYPE_UNSPECIFIED
	}
	key := lib.CoordKey(q, r)
	return wd.Crossings[key]
}

// HasCrossing checks if there's any crossing at the given coordinates
func HasCrossing(wd *v1.WorldData, q, r int32) bool {
	return GetCrossingType(wd, q, r) != v1.CrossingType_CROSSING_TYPE_UNSPECIFIED
}

// HasRoad checks if there's a road at the given coordinates
func HasRoad(wd *v1.WorldData, q, r int32) bool {
	return GetCrossingType(wd, q, r) == v1.CrossingType_CROSSING_TYPE_ROAD
}

// HasBridge checks if there's a bridge at the given coordinates
func HasBridge(wd *v1.WorldData, q, r int32) bool {
	return GetCrossingType(wd, q, r) == v1.CrossingType_CROSSING_TYPE_BRIDGE
}

// GetTileFromMap retrieves a tile from the map-based storage
func GetTileFromMap(wd *v1.WorldData, q, r int32) *v1.Tile {
	if wd == nil || wd.TilesMap == nil {
		return nil
	}
	key := lib.CoordKey(q, r)
	return wd.TilesMap[key]
}

// GetUnitFromMap retrieves a unit from the map-based storage
func GetUnitFromMap(wd *v1.WorldData, q, r int32) *v1.Unit {
	if wd == nil || wd.UnitsMap == nil {
		return nil
	}
	key := lib.CoordKey(q, r)
	return wd.UnitsMap[key]
}

// SetTileInMap adds or updates a tile in the map-based storage
func SetTileInMap(wd *v1.WorldData, tile *v1.Tile) {
	if wd == nil || tile == nil {
		return
	}
	if wd.TilesMap == nil {
		wd.TilesMap = make(map[string]*v1.Tile)
	}
	key := lib.CoordKey(tile.Q, tile.R)
	wd.TilesMap[key] = tile
}

// SetUnitInMap adds or updates a unit in the map-based storage
func SetUnitInMap(wd *v1.WorldData, unit *v1.Unit) {
	if wd == nil || unit == nil {
		return
	}
	if wd.UnitsMap == nil {
		wd.UnitsMap = make(map[string]*v1.Unit)
	}
	key := lib.CoordKey(unit.Q, unit.R)
	wd.UnitsMap[key] = unit
}

// RemoveUnitFromMap removes a unit from the map-based storage
func RemoveUnitFromMap(wd *v1.WorldData, q, r int32) {
	if wd == nil || wd.UnitsMap == nil {
		return
	}
	key := lib.CoordKey(q, r)
	delete(wd.UnitsMap, key)
}

// MoveUnitInMap moves a unit from one position to another in the map
func MoveUnitInMap(wd *v1.WorldData, unit *v1.Unit, toQ, toR int32) {
	if wd == nil || unit == nil {
		return
	}
	// Remove from old position
	RemoveUnitFromMap(wd, unit.Q, unit.R)
	// Update unit coordinates
	unit.Q = toQ
	unit.R = toR
	// Add to new position
	SetUnitInMap(wd, unit)
}

// SetCrossing sets or removes a crossing at the given coordinates
func SetCrossing(wd *v1.WorldData, q, r int32, crossingType v1.CrossingType) {
	if wd == nil {
		return
	}
	if wd.Crossings == nil {
		wd.Crossings = make(map[string]v1.CrossingType)
	}
	key := lib.CoordKey(q, r)
	if crossingType == v1.CrossingType_CROSSING_TYPE_UNSPECIFIED {
		delete(wd.Crossings, key)
	} else {
		wd.Crossings[key] = crossingType
	}
}
