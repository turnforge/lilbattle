package themes

import (
	"fmt"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// DefaultTheme implements the Theme interface for PNG-based default assets
// Mirrors default.ts but focused on metadata
type DefaultTheme struct {
	basePath string

	// Simple mappings for the default theme (no mapping.json for v1 assets)
	unitNames      map[int]string
	terrainNames   map[int]string
	natureTerrains map[int]bool
}

// NewDefaultTheme creates a new default theme instance
func NewDefaultTheme() *DefaultTheme {
	// Initialize unit names (matches default.ts)
	unitNames := map[int]string{
		1: "Infantry", 2: "Mech", 3: "Recon", 4: "Tank", 5: "Medium Tank",
		6: "Neo Tank", 7: "APC", 8: "Artillery", 9: "Rocket", 10: "Anti-Air",
		12: "Fighter", 13: "Bomber", 14: "B-Copter", 15: "T-Copter",
		16: "Battleship", 17: "Cruiser", 18: "Lander", 19: "Sub",
		20: "Mech", 21: "Missile (Std)", 22: "Missile (Nuke)",
		24: "Sailboat", 25: "Artillery (Mega)", 26: "Artillery (Quick)",
		27: "Medic", 28: "Stratotanker", 29: "Engineer", 30: "Goliath RC",
		31: "Tugboat", 32: "Sea Mine", 33: "Drone", 37: "Cruiser",
		38: "Missile (Anti Air)", 39: "Aircraft Carrier", 40: "Miner",
		41: "Paratrooper", 44: "Anti Aircraft (Advanced)",
	}

	// Initialize terrain names (matches default.ts)
	terrainNames := map[int]string{
		0: "Clear", 1: "Land Base", 2: "Naval Base", 3: "Airport Base",
		4: "Desert", 5: "Grass", 6: "Hospital", 7: "Mountains", 8: "Swamp",
		9: "Forest", 10: "Water (Regular)", 12: "Lava", 14: "Water (Shallow)",
		15: "Water (Deep)", 16: "Missile Silo", 17: "Bridge (Regular)",
		18: "Bridge (Shallow)", 19: "Bridge (Deep)", 20: "Mines", 21: "City",
		22: "Road", 23: "Water (Rocky)", 25: "Guard Tower", 26: "Snow",
	}

	// Nature terrains (matches default.ts)
	natureTerrains := make(map[int]bool)
	for _, id := range []int{4, 5, 7, 8, 9, 10, 12, 14, 15, 17, 18, 19, 22, 23, 26} {
		natureTerrains[id] = true
	}

	return &DefaultTheme{
		basePath:       "/static/assets/v1",
		unitNames:      unitNames,
		terrainNames:   terrainNames,
		natureTerrains: natureTerrains,
	}
}

func (d *DefaultTheme) GetUnitName(unitId int) string {
	if name, ok := d.unitNames[unitId]; ok {
		return name
	}
	return ""
}

func (d *DefaultTheme) GetTerrainName(terrainId int) string {
	if name, ok := d.terrainNames[terrainId]; ok {
		return name
	}
	return ""
}

func (d *DefaultTheme) GetUnitDescription(unitId int) string {
	// Default theme has no custom descriptions
	return ""
}

func (d *DefaultTheme) GetTerrainDescription(terrainId int) string {
	// Default theme has no custom descriptions
	return ""
}

// GetUnitPath returns the directory path for a unit's assets
// For PNG themes, this is the directory containing player-colored variants
func (d *DefaultTheme) GetUnitPath(unitId int) string {
	if _, ok := d.unitNames[unitId]; ok {
		return fmt.Sprintf("%s/Units/%d", d.basePath, unitId)
	}
	return ""
}

// GetTilePath returns the directory path for a terrain's assets
func (d *DefaultTheme) GetTilePath(terrainId int) string {
	if _, ok := d.terrainNames[terrainId]; ok {
		return fmt.Sprintf("%s/Tiles/%d", d.basePath, terrainId)
	}
	return ""
}

// GetUnitAssetPath returns the full path to a specific unit+player PNG file
// This is a helper method specific to PNG themes
func (d *DefaultTheme) GetUnitAssetPath(unitId, playerId int) string {
	if _, ok := d.unitNames[unitId]; ok {
		return fmt.Sprintf("%s/Units/%d/%d.png", d.basePath, unitId, playerId)
	}
	return ""
}

// GetTileAssetPath returns the full path to a specific terrain+player PNG file
func (d *DefaultTheme) GetTileAssetPath(terrainId, playerId int) string {
	if _, ok := d.terrainNames[terrainId]; ok {
		// Nature terrains always use player 0 (neutral)
		effectivePlayer := playerId
		if d.natureTerrains[terrainId] {
			effectivePlayer = 0
		}
		return fmt.Sprintf("%s/Tiles/%d/%d.png", d.basePath, terrainId, effectivePlayer)
	}
	return ""
}

func (d *DefaultTheme) IsCityTile(terrainId int) bool {
	return IsCityTerrain(terrainId)
}

func (d *DefaultTheme) IsNatureTile(terrainId int) bool {
	return IsNatureTerrain(terrainId)
}

func (d *DefaultTheme) IsBridgeTile(terrainId int) bool {
	return IsBridgeTerrain(terrainId)
}

func (d *DefaultTheme) GetThemeInfo() *v1.ThemeInfo {
	return &v1.ThemeInfo{
		Name:                "Default (PNG)",
		Version:             "1.0.0",
		BasePath:            d.basePath,
		AssetType:           "png",
		NeedsPostProcessing: false,
	}
}

func (d *DefaultTheme) GetAvailableUnits() []int {
	units := make([]int, 0, len(d.unitNames))
	for id := range d.unitNames {
		units = append(units, id)
	}
	return units
}

func (d *DefaultTheme) GetAvailableTerrains() []int {
	terrains := make([]int, 0, len(d.terrainNames))
	for id := range d.terrainNames {
		terrains = append(terrains, id)
	}
	return terrains
}

func (d *DefaultTheme) HasUnit(unitId int) bool {
	_, ok := d.unitNames[unitId]
	return ok
}

func (d *DefaultTheme) HasTerrain(terrainId int) bool {
	_, ok := d.terrainNames[terrainId]
	return ok
}

// GetAssetPathForTemplate is a helper for templates to get either unit or tile paths
// Returns the full path to the PNG file
func (d *DefaultTheme) GetAssetPathForTemplate(assetType string, assetId, playerId int) string {
	switch assetType {
	case "unit":
		return d.GetUnitAssetPath(assetId, playerId)
	case "tile", "terrain":
		return d.GetTileAssetPath(assetId, playerId)
	default:
		return ""
	}
}
