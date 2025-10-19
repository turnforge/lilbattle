package themes_test

import (
	"fmt"
	"testing"

	"github.com/panyam/turnengine/games/weewar/web/assets/themes"
)

// Example showing how to create and use themes
func ExampleCreateTheme() {
	// Create themes using the registry
	defaultTheme, _ := themes.CreateTheme("default")
	fantasyTheme, _ := themes.CreateTheme("fantasy")
	modernTheme, _ := themes.CreateTheme("modern")

	// Use default theme (PNG-based)
	fmt.Println(defaultTheme.GetUnitName(1))  // "Infantry"
	fmt.Println(defaultTheme.GetThemeInfo().Name) // "Default (PNG)"

	// Use fantasy theme (SVG-based, loaded from embedded mapping.json)
	fmt.Println(fantasyTheme.GetUnitName(1))  // "Peasant"
	fmt.Println(fantasyTheme.GetThemeInfo().Name) // "Medieval Fantasy"

	// Use modern theme (SVG-based, loaded from embedded mapping.json)
	fmt.Println(modernTheme.GetUnitName(1))   // "Infantry"
	fmt.Println(modernTheme.GetThemeInfo().Name) // "Modern Military"

	// Output:
	// Infantry
	// Default (PNG)
	// Peasant
	// Medieval Fantasy
	// Infantry
	// Modern Military
}

// Example showing how to get asset paths
func Example_assetPaths() {
	defaultTheme := themes.NewDefaultTheme()
	fantasyTheme, _ := themes.NewFantasyTheme()

	// Default theme returns full PNG paths
	unitPath := defaultTheme.GetAssetPathForTemplate("unit", 1, 2)
	fmt.Println(unitPath) // /static/assets/v1/Units/1/2.png

	// Fantasy theme returns SVG template paths
	fantasyUnitPath := fantasyTheme.GetUnitAssetPath(1)
	fmt.Println(fantasyUnitPath) // /static/assets/themes/fantasy/Units/Peasant.svg

	// Output:
	// /static/assets/v1/Units/1/2.png
	// /static/assets/themes/fantasy/Units/Peasant.svg
}

// Test that all themes can be created
func TestAllThemes(t *testing.T) {
	themeNames := []string{"default", "fantasy", "modern"}

	for _, name := range themeNames {
		theme, err := themes.CreateTheme(name)
		if err != nil {
			t.Errorf("Failed to create theme %s: %v", name, err)
			continue
		}

		info := theme.GetThemeInfo()
		if info == nil {
			t.Errorf("Theme %s returned nil info", name)
			continue
		}

		t.Logf("Created theme: %s (version %s, type: %s)",
			info.Name, info.Version, info.AssetType)

		// Test that we can get unit names
		unitName := theme.GetUnitName(1)
		if unitName == "" {
			t.Errorf("Theme %s returned empty name for unit 1", name)
		}

		// Test that we can get terrain names
		terrainName := theme.GetTerrainName(1)
		if terrainName == "" {
			t.Errorf("Theme %s returned empty name for terrain 1", name)
		}
	}
}

// Test terrain classification
func TestTerrainClassification(t *testing.T) {
	theme := themes.NewDefaultTheme()

	// Test city tiles
	if !theme.IsCityTile(1) { // Land Base
		t.Error("Land Base should be classified as city tile")
	}

	// Test nature tiles
	if !theme.IsNatureTile(5) { // Grass
		t.Error("Grass should be classified as nature tile")
	}

	// Test bridge tiles
	if !theme.IsBridgeTile(17) { // Bridge (Regular)
		t.Error("Bridge should be classified as bridge tile")
	}
}
