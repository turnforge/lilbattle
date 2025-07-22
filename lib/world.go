package weewar

import "fmt"

// =============================================================================
// World - Pure Game State Container
// =============================================================================

// World represents the pure game state without any rendering or UI concerns.
// This is the single source of truth for all game data.
type World struct {
	Map           *Map                 `json:"map"`   // The game map with terrain and tiles
	UnitsByPlayer [][]*Unit            `json:"units"` // All units in the game world by player ID
	UnitsByCoord  map[AxialCoord]*Unit `json:"-"`     // All units in the game world by player ID

	// Observer pattern for state changes
	WorldSubject `json:"-"`

	// Pure state metadata
	PlayerCount int `json:"playerCount"` // Number of players in the game
}

// ViewState represents UI-specific state that doesn't affect game logic.
// This includes visual concerns like selections, highlights, and camera position.
type ViewState struct {
	// Selection and highlighting
	SelectedUnit    *Unit      `json:"selectedUnit"`    // Currently selected unit
	HoveredTile     *Tile      `json:"hoveredTile"`     // Tile under cursor
	MovableTiles    []Position `json:"movableTiles"`    // Highlighted movement tiles
	AttackableTiles []Position `json:"attackableTiles"` // Highlighted attack tiles

	// Visual settings
	ShowGrid        bool `json:"showGrid"`        // Whether to show hex grid lines
	ShowCoordinates bool `json:"showCoordinates"` // Whether to show coordinate labels
	ShowPaths       bool `json:"showPaths"`       // Whether to show movement paths

	// Camera and viewport
	CameraX   float64 `json:"cameraX"`   // Camera X position
	CameraY   float64 `json:"cameraY"`   // Camera Y position
	ZoomLevel float64 `json:"zoomLevel"` // Zoom level (1.0 = normal)

	// Editor-specific state
	BrushTerrain int `json:"brushTerrain"` // Current terrain type for painting
	BrushSize    int `json:"brushSize"`    // Brush radius (0 = single hex)
}

// Note: Position type is already defined in game_interface.go

// =============================================================================
// World Creation and Management
// =============================================================================

// NewWorld creates a new game world with the specified parameters
func NewWorld(playerCount int, gameMap *Map) (*World, error) {
	if playerCount < 2 || playerCount > MaxUnits {
		return nil, fmt.Errorf("invalid player count: %d (must be 2-%d)", playerCount, MaxUnits)
	}

	if gameMap == nil {
		return nil, fmt.Errorf("map cannot be nil")
	}

	w := &World{
		Map:           gameMap,
		UnitsByPlayer: make([][]*Unit, 0),
		UnitsByCoord:  make(map[AxialCoord]*Unit),
		PlayerCount:   playerCount,
	}

	w.UnitsByPlayer = make([][]*Unit, playerCount)
	for i := range playerCount {
		w.UnitsByPlayer[i] = nil
	}

	return w, nil
}

// NewViewState creates a new view state with default settings
func NewViewState() *ViewState {
	return &ViewState{
		SelectedUnit:    nil,
		HoveredTile:     nil,
		MovableTiles:    make([]Position, 0),
		AttackableTiles: make([]Position, 0),
		ShowGrid:        true,
		ShowCoordinates: false,
		ShowPaths:       true,
		CameraX:         0.0,
		CameraY:         0.0,
		ZoomLevel:       1.0,
		BrushTerrain:    1, // Default to grass
		BrushSize:       0, // Single hex brush
	}
}

// =============================================================================
// World State Access Methods
// =============================================================================

// GetMapSize returns the dimensions of the world map
func (w *World) GetMapSizeRect() (rows, cols int) {
	if w.Map == nil {
		return 0, 0
	}
	return w.Map.NumRows(), w.Map.NumCols()
}

func (w *World) UnitAt(coord AxialCoord) *Unit {
	return w.UnitsByCoord[coord]
}

// GetTileAt returns the tile at the specified cube coordinates
func (w *World) GetTileAt(coord AxialCoord) *Tile {
	if w.Map == nil {
		return nil
	}
	return w.Map.TileAt(coord)
}

// GetUnitsAt returns all units at the specified display coordinates
func (w *World) GetUnitsAt(coord AxialCoord) []*Unit {
	// Use efficient O(1) lookup with UnitsByCoord
	if unit := w.UnitsByCoord[coord]; unit != nil {
		return []*Unit{unit}
	}
	return []*Unit{}
}

// GetPlayerUnits returns all units belonging to the specified player
func (w *World) GetPlayerUnits(playerID int) []*Unit {
	return w.UnitsByPlayer[playerID]
}

// =============================================================================
// World State Mutation Methods
// =============================================================================

// SetTileTypeCube changes the terrain type at the specified cube coordinates
func (w *World) SetTileType(coord AxialCoord, terrainType int) bool {
	if w.Map == nil {
		return false
	}

	// Get or create tile at position
	tile := w.Map.TileAt(coord)
	if tile == nil {
		// Create new tile
		tile = NewTile(coord, terrainType)
		w.Map.AddTile(tile)
	} else {
		// Update existing tile
		tile.TileType = terrainType
	}

	return true
}

// AddUnit adds a new unit to the world at the specified position
func (w *World) AddUnit(unit *Unit) (oldunit *Unit, err error) {
	if unit == nil {
		return nil, fmt.Errorf("unit is nil")
	}

	playerID := unit.PlayerID
	if playerID < 0 || playerID >= len(w.UnitsByPlayer) {
		return nil, fmt.Errorf("invalid player ID: %d", playerID)
	}

	oldunit = w.UnitAt(unit.Coord)

	// make sure to replace a unit here
	w.UnitsByPlayer[unit.PlayerID] = append(w.UnitsByPlayer[unit.PlayerID], unit)
	w.UnitsByCoord[unit.Coord] = unit
	return
}

// RemoveUnit removes a unit from the world
func (w *World) RemoveUnit(unit *Unit) error {
	if unit == nil {
		return fmt.Errorf("unit is nil")
	}

	tile := w.Map.TileAt(unit.Coord)
	if tile == nil {
		return fmt.Errorf("invalid tile")
	}
	p := unit.PlayerID
	delete(w.UnitsByCoord, unit.Coord)
	for i, u := range w.UnitsByPlayer[p] {
		if u == unit {
			// Remove unit from slice
			w.UnitsByPlayer[p] = append(w.UnitsByPlayer[p][:i], w.UnitsByPlayer[p][i+1:]...)
			break
		}
	}
	return nil
}

// MoveUnit moves a unit to a new position
func (w *World) MoveUnit(unit *Unit, newCoord AxialCoord) error {
	if unit == nil {
		return fmt.Errorf("unit is nil")
	}

	// Remove from old position
	delete(w.UnitsByCoord, unit.Coord)

	// Update unit position
	unit.Coord = newCoord

	// Add to new position
	w.UnitsByCoord[newCoord] = unit

	return nil
}

// =============================================================================
// World Validation and Utilities
// =============================================================================

// GetWorldBounds returns the bounding box of the world in display coordinates
func (w *World) GetWorldBoundsRect() (minRow, minCol, maxRow, maxCol int) {
	if w.Map == nil {
		return 0, 0, 0, 0
	}
	return 0, 0, w.Map.NumRows() - 1, w.Map.NumCols() - 1
}

// Clone creates a deep copy of the world state (useful for undo/redo systems)
func (w *World) Clone() *World {
	if w == nil {
		return nil
	}

	// Clone map
	var clonedMap *Map
	if w.Map != nil {
		clonedMap = NewMapRect(w.Map.NumRows(), w.Map.NumCols())
		for _, tile := range w.Map.Tiles {
			if tile != nil {
				newTile := tile.Clone()
				clonedMap.AddTile(newTile)
			}
		}
	}

	// Clone units
	clonedUnits := make([][]*Unit, len(w.UnitsByPlayer))
	for playerId, units := range w.UnitsByPlayer {
		for i, unit := range units {
			if unit != nil {
				clonedUnits[playerId][i] = unit.Clone()
			}
		}
	}

	return &World{
		Map:           clonedMap,
		UnitsByPlayer: clonedUnits,
		PlayerCount:   w.PlayerCount,
	}
}

// =============================================================================
// ViewState Management
// =============================================================================

// ClearSelection clears the current unit selection and highlights
func (vs *ViewState) ClearSelection() {
	vs.SelectedUnit = nil
	vs.MovableTiles = make([]Position, 0)
	vs.AttackableTiles = make([]Position, 0)
}

// SetSelection sets the selected unit and updates related highlights
func (vs *ViewState) SetSelection(unit *Unit, movableTiles, attackableTiles []Position) {
	vs.SelectedUnit = unit
	vs.MovableTiles = movableTiles
	vs.AttackableTiles = attackableTiles
}

// SetCamera updates the camera position and zoom
func (vs *ViewState) SetCamera(x, y, zoom float64) {
	vs.CameraX = x
	vs.CameraY = y
	vs.ZoomLevel = zoom
}

// SetBrush updates the brush settings for terrain editing
func (vs *ViewState) SetBrush(terrainType, brushSize int) {
	vs.BrushTerrain = terrainType
	vs.BrushSize = brushSize
}
