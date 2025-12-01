package themes

import (
	"fmt"
)

// ThemeRegistry provides a factory for creating themes by name
type ThemeRegistry struct {
	themes map[string]func(cityTerrains map[int32]bool) (Theme, error)
}

// NewThemeRegistry creates a new theme registry with all available themes
func NewThemeRegistry() *ThemeRegistry {
	registry := &ThemeRegistry{
		themes: make(map[string]func(cityTerrains map[int32]bool) (Theme, error)),
	}

	// Register default theme
	registry.Register("default", func(cityTerrains map[int32]bool) (Theme, error) {
		return NewDefaultTheme(cityTerrains), nil
	})

	// Register fantasy theme
	registry.Register("fantasy", func(cityTerrains map[int32]bool) (Theme, error) {
		return NewFantasyTheme(cityTerrains)
	})

	// Register modern theme
	registry.Register("modern", func(cityTerrains map[int32]bool) (Theme, error) {
		return NewModernTheme(cityTerrains)
	})

	return registry
}

// Register adds a theme factory to the registry
func (r *ThemeRegistry) Register(name string, factory func(cityTerrains map[int32]bool) (Theme, error)) {
	r.themes[name] = factory
}

// Create creates a theme instance by name
func (r *ThemeRegistry) Create(name string, cityTerrains map[int32]bool) (Theme, error) {
	factory, ok := r.themes[name]
	if !ok {
		return nil, fmt.Errorf("unknown theme: %s", name)
	}
	return factory(cityTerrains)
}

// GetAvailableThemes returns a list of all registered theme names
func (r *ThemeRegistry) GetAvailableThemes() []string {
	names := make([]string, 0, len(r.themes))
	for name := range r.themes {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the global theme registry
var DefaultRegistry = NewThemeRegistry()

// CreateTheme is a convenience function that uses the default registry
func CreateTheme(name string, cityTerrains map[int32]bool) (Theme, error) {
	return DefaultRegistry.Create(name, cityTerrains)
}
