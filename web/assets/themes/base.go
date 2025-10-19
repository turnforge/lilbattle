package themes

import (
	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// BaseTheme provides common functionality for all themes
// Mirrors BaseTheme.ts but focused on metadata, not asset loading
type BaseTheme struct {
	manifest *v1.ThemeManifest
}

// NewBaseTheme creates a new BaseTheme from a pre-loaded manifest
func NewBaseTheme(manifest *v1.ThemeManifest) *BaseTheme {
	return &BaseTheme{
		manifest: manifest,
	}
}

func (b *BaseTheme) GetUnitName(unitId int) string {
	if mapping, ok := b.manifest.Units[int32(unitId)]; ok {
		return mapping.Name
	}
	return ""
}

func (b *BaseTheme) GetTerrainName(terrainId int) string {
	if mapping, ok := b.manifest.Terrains[int32(terrainId)]; ok {
		return mapping.Name
	}
	return ""
}

func (b *BaseTheme) GetUnitDescription(unitId int) string {
	if mapping, ok := b.manifest.Units[int32(unitId)]; ok {
		return mapping.Description
	}
	return ""
}

func (b *BaseTheme) GetTerrainDescription(terrainId int) string {
	if mapping, ok := b.manifest.Terrains[int32(terrainId)]; ok {
		return mapping.Description
	}
	return ""
}

func (b *BaseTheme) GetUnitPath(unitId int) string {
	if mapping, ok := b.manifest.Units[int32(unitId)]; ok {
		// Return relative path from theme base
		return mapping.Image
	}
	return ""
}

func (b *BaseTheme) GetTilePath(terrainId int) string {
	if mapping, ok := b.manifest.Terrains[int32(terrainId)]; ok {
		// Return relative path from theme base
		return mapping.Image
	}
	return ""
}

func (b *BaseTheme) IsCityTile(terrainId int) bool {
	return IsCityTerrain(terrainId)
}

func (b *BaseTheme) IsNatureTile(terrainId int) bool {
	return IsNatureTerrain(terrainId)
}

func (b *BaseTheme) IsBridgeTile(terrainId int) bool {
	return IsBridgeTerrain(terrainId)
}

func (b *BaseTheme) GetThemeInfo() *v1.ThemeInfo {
	return b.manifest.ThemeInfo
}

func (b *BaseTheme) GetAvailableUnits() []int {
	units := make([]int, 0, len(b.manifest.Units))
	for id := range b.manifest.Units {
		units = append(units, int(id))
	}
	return units
}

func (b *BaseTheme) GetAvailableTerrains() []int {
	terrains := make([]int, 0, len(b.manifest.Terrains))
	for id := range b.manifest.Terrains {
		terrains = append(terrains, int(id))
	}
	return terrains
}

func (b *BaseTheme) HasUnit(unitId int) bool {
	_, ok := b.manifest.Units[int32(unitId)]
	return ok
}

func (b *BaseTheme) HasTerrain(terrainId int) bool {
	_, ok := b.manifest.Terrains[int32(terrainId)]
	return ok
}
