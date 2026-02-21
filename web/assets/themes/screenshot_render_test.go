package themes_test

import (
	"testing"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	"github.com/turnforge/lilbattle/web/assets/themes"
)

// smallTestWorld creates a minimal tile/unit map for rendering tests
func smallTestWorld() (map[string]*v1.Tile, map[string]*v1.Unit) {
	tiles := map[string]*v1.Tile{
		"0,0": {Q: 0, R: 0, TileType: 1, Player: 1},  // Land Base
		"1,0": {Q: 1, R: 0, TileType: 5, Player: 0},  // Grass
		"0,1": {Q: 0, R: 1, TileType: 10, Player: 0}, // Water
	}
	units := map[string]*v1.Unit{
		"0,0": {Q: 0, R: 0, Player: 1, UnitType: 1, AvailableHealth: 10},
	}
	return tiles, units
}

func TestScreenshotRendering(t *testing.T) {
	cityTerrains := testCityTerrains()
	tiles, units := smallTestWorld()

	themeNames := []string{"default", "modern", "fantasy"}

	for _, themeName := range themeNames {
		t.Run(themeName, func(t *testing.T) {
			theme, err := themes.CreateTheme(themeName, cityTerrains)
			if err != nil {
				t.Fatalf("CreateTheme(%s): %v", themeName, err)
			}

			renderer, err := themes.CreateWorldRenderer(theme)
			if err != nil {
				t.Fatalf("CreateWorldRenderer(%s): %v", themeName, err)
			}

			imageBytes, contentType, err := renderer.Render(tiles, units, nil)
			if err != nil {
				t.Fatalf("Render(%s): %v", themeName, err)
			}

			if len(imageBytes) == 0 {
				t.Errorf("Render(%s) returned empty image", themeName)
			}

			info := theme.GetThemeInfo()
			switch info.AssetType {
			case "png":
				if contentType != "image/png" {
					t.Errorf("PNG theme %s returned content type %s", themeName, contentType)
				}
				// Check PNG magic bytes
				if len(imageBytes) < 4 || imageBytes[0] != 0x89 || imageBytes[1] != 'P' || imageBytes[2] != 'N' || imageBytes[3] != 'G' {
					t.Errorf("PNG theme %s output doesn't start with PNG magic bytes", themeName)
				}
			case "svg":
				if contentType != "image/svg+xml" {
					t.Errorf("SVG theme %s returned content type %s", themeName, contentType)
				}
			}

			t.Logf("%s: rendered %d bytes (%s)", themeName, len(imageBytes), contentType)
		})
	}
}

func TestScreenshotRenderingEmptyWorld(t *testing.T) {
	cityTerrains := testCityTerrains()
	emptyTiles := map[string]*v1.Tile{}
	emptyUnits := map[string]*v1.Unit{}

	theme, err := themes.CreateTheme("default", cityTerrains)
	if err != nil {
		t.Fatalf("CreateTheme: %v", err)
	}

	renderer, err := themes.CreateWorldRenderer(theme)
	if err != nil {
		t.Fatalf("CreateWorldRenderer: %v", err)
	}

	// Rendering an empty world should not panic
	_, _, err = renderer.Render(emptyTiles, emptyUnits, nil)
	if err != nil {
		t.Logf("Empty world render returned error (acceptable): %v", err)
	}
}
