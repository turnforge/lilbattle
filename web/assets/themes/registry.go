package themes

import (
	"fmt"
)

// ThemeRegistry provides a factory for creating themes by name
type ThemeRegistry struct {
	themes map[string]func() (Theme, error)
}

// NewThemeRegistry creates a new theme registry with all available themes
func NewThemeRegistry() *ThemeRegistry {
	registry := &ThemeRegistry{
		themes: make(map[string]func() (Theme, error)),
	}

	// Register default theme
	registry.Register("default", func() (Theme, error) {
		return NewDefaultTheme(), nil
	})

	// Register fantasy theme
	registry.Register("fantasy", func() (Theme, error) {
		return NewFantasyTheme()
	})

	// Register modern theme
	registry.Register("modern", func() (Theme, error) {
		return NewModernTheme()
	})

	return registry
}

// Register adds a theme factory to the registry
func (r *ThemeRegistry) Register(name string, factory func() (Theme, error)) {
	r.themes[name] = factory
}

// Create creates a theme instance by name
func (r *ThemeRegistry) Create(name string) (Theme, error) {
	factory, ok := r.themes[name]
	if !ok {
		return nil, fmt.Errorf("unknown theme: %s", name)
	}
	return factory()
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
func CreateTheme(name string) (Theme, error) {
	return DefaultRegistry.Create(name)
}
