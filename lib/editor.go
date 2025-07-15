package weewar

import (
	"fmt"
	"math"
	"path/filepath"
)

// =============================================================================
// World Editor Core
// =============================================================================

// WorldEditor provides tools for creating and editing game worlds (maps, units, etc.)
type WorldEditor struct {
	// Current world being edited
	currentWorld *World

	// Editor state
	filename     string
	modified     bool
	brushTerrain int // Current terrain type for painting
	brushSize    int // Brush radius (0 = single hex, 1 = 7 hexes, etc.)

	// Rendering and viewport
	drawable        Drawable         // Platform-agnostic drawing interface
	layeredRenderer *LayeredRenderer // Fast layered rendering system
	canvasWidth     int
	canvasHeight    int
	scrollX         float64 // Horizontal scroll offset for viewport
	scrollY         float64 // Vertical scroll offset for viewport
}

// NewWorldEditor creates a new world editor instance
func NewWorldEditor() *WorldEditor {
	return &WorldEditor{
		currentWorld: nil,
		filename:     "",
		modified:     false,
		brushTerrain: 1, // Default to grass
		brushSize:    0, // Single hex brush
		scrollX:      0, // No initial scroll offset
		scrollY:      0, // No initial scroll offset
	}
}

// =============================================================================
// Map Management
// =============================================================================

// NewWorld creates a new 1x1 world for editing (use Add/Remove methods to expand)
func (e *WorldEditor) NewWorld() error {
	// Create new map with single tile at origin (Q=0, R=0)
	gameMap := NewMapWithBounds(0, 0, 0, 0)
	e.filename = ""
	e.modified = false

	// Create single tile at origin with default terrain (grass)
	coord := CubeCoord{Q: 0, R: 0}
	tile := NewTile(coord, 1) // Grass terrain
	gameMap.AddTile(tile)

	// Create world with map and initialize units by player
	e.currentWorld = &World{
		Map:           gameMap,
		UnitsByPlayer: make([][]*Unit, 2), // Start with 2 players
		PlayerCount:   2,
	}

	// Update layered renderer with new world
	if e.layeredRenderer != nil {
		e.layeredRenderer.SetWorld(e.currentWorld)
	}

	return nil
}

// LoadMap loads an existing map for editing
func (e *WorldEditor) LoadMap(filename string) error {
	// TODO: Implement map loading from file
	// For now, create a placeholder implementation
	return fmt.Errorf("map loading not yet implemented")
}

// SaveMap saves the current map to file
func (e *WorldEditor) SaveMap(filename string) error {
	if e.currentWorld == nil {
		return fmt.Errorf("no map to save")
	}

	// TODO: Implement map saving to file
	// For now, just update the filename and mark as unmodified
	e.filename = filename
	e.modified = false

	return nil
}

// GetCurrentMap returns the map being edited (read-only access)
func (e *WorldEditor) GetCurrentMap() *Map {
	return e.currentWorld.Map
}

// IsModified returns whether the map has unsaved changes
func (e *WorldEditor) IsModified() bool {
	return e.modified
}

// CalculateCanvasSize returns the optimal canvas size for the current map
func (e *WorldEditor) CalculateCanvasSize(rows, cols int) (width, height int) {
	if e.currentWorld == nil {
		return 400, 300 // Default size for empty editor
	}

	// Use the world renderer to calculate proper tile dimensions
	renderer := &BaseRenderer{}
	tempWorld := &World{Map: e.currentWorld.Map}

	// Use a reasonable base canvas size for calculation
	options := renderer.CalculateRenderOptions(800, 600, tempWorld)

	mapWidth, mapHeight := e.currentWorld.Map.CanvasSize(options.TileWidth, options.TileHeight, options.YIncrement)

	// Calculate canvas dimensions with minimal padding
	width = int(mapWidth + math.Max(e.scrollX, 0))
	height = int(mapHeight + math.Max(e.scrollY, 0))
	fmt.Println("Ok here, rows, cols, w, h: ", rows, cols, width, height, options)

	// Ensure minimum size for usability
	if width < 200 {
		width = 200
	}
	if height < 200 {
		height = 200
	}

	return width, height
}

// GetFilename returns the current filename (empty if new map)
func (e *WorldEditor) GetFilename() string {
	return e.filename
}

// GetCanvasSize returns the current canvas dimensions
func (e *WorldEditor) GetCanvasSize() (width, height int) {
	return e.canvasWidth, e.canvasHeight
}

// GetLayeredRenderer returns the layered renderer for direct access
func (e *WorldEditor) GetLayeredRenderer() *LayeredRenderer {
	return e.layeredRenderer
}

// SetAssetProvider updates the asset provider for terrain/unit sprites
func (e *WorldEditor) SetAssetProvider(provider AssetProvider) {
	if e.layeredRenderer != nil {
		e.layeredRenderer.SetAssetProvider(provider)
	}
}

// =============================================================================
// Terrain Editing
// =============================================================================

// SetBrushTerrain sets the terrain type for painting
func (e *WorldEditor) SetBrushTerrain(terrainType int) error {
	if terrainType < 0 || terrainType >= len(terrainData) {
		return fmt.Errorf("invalid terrain type: %d", terrainType)
	}
	e.brushTerrain = terrainType
	return nil
}

// SetBrushSize sets the brush radius (0 = single hex, 1 = 7 hexes, etc.)
func (e *WorldEditor) SetBrushSize(size int) error {
	if size < 0 || size > 5 {
		return fmt.Errorf("invalid brush size: %d (must be 0-5)", size)
	}
	e.brushSize = size
	return nil
}

// PaintTerrain paints terrain at the specified display position
func (e *WorldEditor) PaintTerrain(row, col int) error {
	if e.currentWorld == nil {
		return fmt.Errorf("no map loaded")
	}

	// Convert to cube coordinate
	centerCoord := e.currentWorld.RowColToHex(row, col)

	// Get all positions to paint based on brush size
	positions := e.getBrushPositions(centerCoord)

	// Paint each position
	for _, coord := range positions {
		// Check if position is within map bounds
		displayRow, displayCol := e.currentWorld.HexToRowCol(coord)
		if displayRow < 0 || displayRow >= e.currentWorld.NumRows() ||
			displayCol < 0 || displayCol >= e.currentWorld.NumCols() {
			continue // Skip out-of-bounds positions
		}

		// Get existing tile or create new one
		tile := e.currentWorld.TileAt(coord)
		if tile == nil {
			tile = NewTile(displayRow, displayCol, e.brushTerrain)
			e.currentWorld.AddTile(tile)
		} else {
			tile.TileType = e.brushTerrain
		}
	}

	// Take snapshot after making changes
	e.takeSnapshot()
	e.modified = true

	// Mark affected tiles as dirty for efficient rendering
	if e.layeredRenderer != nil {
		for _, coord := range positions {
			e.layeredRenderer.MarkTerrainDirty(coord)
		}
	}

	return nil
}

// RemoveTerrain removes terrain at the specified position
func (e *WorldEditor) RemoveTerrain(row, col int) error {
	if e.currentWorld == nil {
		return fmt.Errorf("no map loaded")
	}

	coord := e.currentWorld.RowColToHex(row, col)
	e.currentWorld.DeleteTileCube(coord)

	e.takeSnapshot()
	e.modified = true

	// Mark affected tile as dirty for efficient rendering
	if e.layeredRenderer != nil {
		e.layeredRenderer.MarkTerrainDirty(coord)
	}

	return nil
}

// FloodFill fills a connected region with the current brush terrain
func (e *WorldEditor) FloodFill(row, col int) error {
	if e.currentWorld == nil {
		return fmt.Errorf("no map loaded")
	}

	startCoord := e.currentWorld.RowColToHex(row, col)
	startTile := e.currentWorld.TileAtCube(startCoord)
	if startTile == nil {
		return fmt.Errorf("no tile at position (%d, %d)", row, col)
	}

	originalTerrain := startTile.TileType
	if originalTerrain == e.brushTerrain {
		return nil // Already the target terrain
	}

	// Use breadth-first search for flood fill
	visited := make(map[CubeCoord]bool)
	queue := []CubeCoord{startCoord}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		// Check if this position has the original terrain
		tile := e.currentWorld.TileAtCube(current)
		if tile == nil || tile.TileType != originalTerrain {
			continue
		}

		// Change terrain
		tile.TileType = e.brushTerrain

		// Add neighbors to queue
		neighbors := current.Neighbors()
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				// Check if neighbor is within bounds
				nRow, nCol := e.currentWorld.HexToRowCol(neighbor)
				if nRow >= 0 && nRow < e.currentWorld.NumRows() &&
					nCol >= 0 && nCol < e.currentWorld.NumCols() {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	e.takeSnapshot()
	e.modified = true

	// Mark entire terrain as dirty since flood fill can affect large areas
	if e.layeredRenderer != nil {
		e.layeredRenderer.MarkAllTerrainDirty()
	}

	return nil
}

// =============================================================================
// History Management (Undo/Redo)
// =============================================================================

// Undo reverts the last operation
func (e *WorldEditor) Undo() error {
	if e.historyPos <= 0 {
		return fmt.Errorf("nothing to undo")
	}

	e.historyPos--
	e.currentWorld = e.copyMap(e.history[e.historyPos])
	e.modified = true

	// Mark entire terrain as dirty since undo can affect entire state
	if e.layeredRenderer != nil {
		e.layeredRenderer.MarkAllTerrainDirty()
	}

	return nil
}

// Redo reapplies the next operation
func (e *WorldEditor) Redo() error {
	if e.historyPos >= len(e.history)-1 {
		return fmt.Errorf("nothing to redo")
	}

	e.historyPos++
	e.currentWorld = e.copyMap(e.history[e.historyPos])
	e.modified = true

	// Mark entire terrain as dirty since redo can affect entire state
	if e.layeredRenderer != nil {
		e.layeredRenderer.MarkAllTerrainDirty()
	}

	return nil
}

// CanUndo returns whether undo is available
func (e *WorldEditor) CanUndo() bool {
	return e.historyPos > 0
}

// CanRedo returns whether redo is available
func (e *WorldEditor) CanRedo() bool {
	return e.historyPos < len(e.history)-1
}

// =============================================================================
// Utility Methods
// =============================================================================

// getBrushPositions returns all cube coordinates affected by the current brush
func (e *WorldEditor) getBrushPositions(center CubeCoord) []CubeCoord {
	if e.brushSize == 0 {
		return []CubeCoord{center}
	}

	// Use the Range method from cube coordinates
	return center.Range(e.brushSize)
}

// takeSnapshot saves the current map state for undo
func (e *WorldEditor) takeSnapshot() {
	if e.currentWorld == nil {
		return
	}

	// Remove any redo history when taking a new snapshot
	if e.historyPos < len(e.history)-1 {
		e.history = e.history[:e.historyPos+1]
	}

	// Add current map to history
	mapCopy := e.copyMap(e.currentWorld)
	e.history = append(e.history, mapCopy)
	e.historyPos = len(e.history) - 1

	// Limit history size
	if len(e.history) > e.maxHistory {
		e.history = e.history[1:]
		e.historyPos--
	}
}

// copyMap creates a deep copy of a map
func (e *WorldEditor) copyMap(original *Map) *Map {
	if original == nil {
		return nil
	}

	copy := NewMap(original.NumRows(), original.NumCols(), false)

	// Copy all tiles
	for coord, tile := range original.Tiles {
		if tile != nil {
			newTile := &Tile{
				Row:      tile.Row,
				Col:      tile.Col,
				TileType: tile.TileType,
				Unit:     nil, // Don't copy units in editor
			}
			copy.AddTileCube(coord, newTile)
		}
	}

	return copy
}

// clearHistory clears the undo/redo history
func (e *WorldEditor) clearHistory() {
	e.history = make([]*Map, 0)
	e.historyPos = -1
}

// =============================================================================
// Map Information
// =============================================================================

// GetMapInfo returns information about the current map
func (e *WorldEditor) GetMapInfo() *MapInfo {
	if e.currentWorld == nil {
		return nil
	}

	// Count terrain types
	terrainCounts := make(map[int]int)
	totalTiles := 0

	for _, tile := range e.currentWorld.Tiles {
		if tile != nil {
			terrainCounts[tile.TileType]++
			totalTiles++
		}
	}

	return &MapInfo{
		Filename:      e.filename,
		Width:         e.currentWorld.NumCols(),
		Height:        e.currentWorld.NumRows(),
		TotalTiles:    totalTiles,
		TerrainCounts: terrainCounts,
		Modified:      e.modified,
	}
}

// MapInfo contains information about a map
type MapInfo struct {
	Filename      string
	Width         int
	Height        int
	TotalTiles    int
	TerrainCounts map[int]int
	Modified      bool
}

// =============================================================================
// Map Validation
// =============================================================================

// ValidateMap checks the map for common issues
func (e *WorldEditor) ValidateMap() []string {
	if e.currentWorld == nil {
		return []string{"No map loaded"}
	}

	var issues []string

	// Check for missing tiles (holes in the map)
	expectedTiles := e.currentWorld.NumRows() * e.currentWorld.NumCols()
	actualTiles := len(e.currentWorld.Tiles)

	if actualTiles < expectedTiles {
		issues = append(issues, fmt.Sprintf("Map has holes: %d tiles missing", expectedTiles-actualTiles))
	}

	// Check for invalid terrain types
	for _, tile := range e.currentWorld.Tiles {
		if tile != nil {
			if tile.TileType < 0 || tile.TileType >= len(terrainData) {
				issues = append(issues, fmt.Sprintf("Invalid terrain type %d at (%d, %d)",
					tile.TileType, tile.Row, tile.Col))
			}
		}
	}

	// Check map dimensions
	if e.currentWorld.NumRows() < 3 || e.currentWorld.NumCols() < 3 {
		issues = append(issues, "Map is very small (recommended minimum 3x3)")
	}

	if e.currentWorld.NumRows() > 50 || e.currentWorld.NumCols() > 50 {
		issues = append(issues, "Map is very large (may cause performance issues)")
	}

	return issues
}

// =============================================================================
// Export Functions
// =============================================================================

// ExportToGame converts the edited map to a Game instance for testing
func (e *WorldEditor) ExportToGame(playerCount int) (*Game, error) {
	if e.currentWorld == nil {
		return nil, fmt.Errorf("no map to export")
	}

	// Validate player count
	if playerCount < 2 || playerCount > 6 {
		return nil, fmt.Errorf("invalid player count: %d (must be 2-6)", playerCount)
	}

	// Create a copy of the map for the game
	gamemap := e.copyMap(e.currentWorld)

	// Create the game
	game, err := NewGame(playerCount, gamemap, 12345) // Use fixed seed for testing
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

// RenderToFile saves the current map as a PNG image
func (e *WorldEditor) RenderToFile(filename string, width, height int) error {
	if e.currentWorld == nil {
		return fmt.Errorf("no map to render")
	}

	// Create a temporary game for rendering
	game, err := e.ExportToGame(2)
	if err != nil {
		return fmt.Errorf("failed to create game for rendering: %w", err)
	}

	// Create buffer and render
	buffer := NewBuffer(width, height)

	// Calculate tile size based on map dimensions and buffer size
	tileWidth := float64(width) / float64(e.currentWorld.NumCols())
	tileHeight := float64(height) / float64(e.currentWorld.NumRows())
	yIncrement := tileHeight * 0.75 // Hex grid spacing

	err = game.RenderToBuffer(buffer, tileWidth, tileHeight, yIncrement)
	if err != nil {
		return fmt.Errorf("failed to render map: %w", err)
	}

	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != ".png" {
		filename += ".png"
	}

	return buffer.Save(filename)
}

// =============================================================================
// Canvas Management
// =============================================================================

// SetCanvas initializes the canvas for real-time rendering
func (e *WorldEditor) SetCanvas(canvasID string, width, height int) error {
	// Create new layered renderer for fast prototyping
	var err error
	e.layeredRenderer, err = NewLayeredRenderer(canvasID, width, height)
	if err != nil {
		return fmt.Errorf("failed to create layered renderer for '%s': %v", canvasID, err)
	}

	e.canvasWidth = width
	e.canvasHeight = height

	// If we have a current map, mark all terrain as dirty for initial render
	if e.currentWorld != nil {
		e.layeredRenderer.MarkAllTerrainDirty()
	}

	return nil
}

// SetCanvasSize resizes the canvas
func (e *WorldEditor) SetCanvasSize(width, height int) error {
	if e.layeredRenderer == nil {
		return fmt.Errorf("no layered renderer initialized")
	}

	e.canvasWidth = width
	e.canvasHeight = height

	// Resize the layered renderer (this will mark everything as dirty)
	err := e.layeredRenderer.Resize(width, height)
	if err != nil {
		return fmt.Errorf("failed to resize layered renderer: %v", err)
	}

	return nil
}

// renderFullMap renders the entire current map to the canvas
func (e *WorldEditor) renderFullMap() error {
	if e.drawable == nil || e.currentWorld == nil {
		return nil // No canvas or map to render
	}

	// Simplified rendering directly using FillPath for each tile
	tileWidth := float64(e.canvasWidth) / float64(e.currentWorld.NumCols())
	tileHeight := float64(e.canvasHeight) / float64(e.currentWorld.NumRows())

	// Render each tile as a hexagon
	for coord, tile := range e.currentWorld.Tiles {
		if tile == nil {
			continue
		}

		// Convert hex coordinates to display coordinates
		displayRow, displayCol := e.currentWorld.HexToRowCol(coord)

		// Calculate tile position
		x := float64(displayCol) * tileWidth
		y := float64(displayRow) * (tileHeight * 0.75) // Hex grid spacing

		// Offset even rows for hex grid
		if displayRow%2 == 0 {
			x += tileWidth * 0.5
		}

		// Create hexagon points
		hexPoints := createHexPoints(x+tileWidth/2, y+tileHeight/2, tileWidth*0.4)

		// Get terrain color
		var fillColor Color
		switch tile.TileType {
		case 1: // Grass
			fillColor = Color{R: 34, G: 139, B: 34, A: 255}
		case 2: // Desert
			fillColor = Color{R: 238, G: 203, B: 173, A: 255}
		case 3: // Water
			fillColor = Color{R: 65, G: 105, B: 225, A: 255}
		case 4: // Mountain
			fillColor = Color{R: 139, G: 137, B: 137, A: 255}
		case 5: // Rock
			fillColor = Color{R: 105, G: 105, B: 105, A: 255}
		default:
			fillColor = Color{R: 200, G: 200, B: 200, A: 255}
		}

		// Fill the hexagon
		e.drawable.FillPath(hexPoints, fillColor)

		// Draw border
		borderColor := Color{R: 0, G: 0, B: 0, A: 100}
		strokeProps := StrokeProperties{Width: 1.0, LineCap: "round", LineJoin: "round"}
		e.drawable.StrokePath(hexPoints, borderColor, strokeProps)
	}

	return nil
}

// renderTiles renders specific tiles to the canvas (for partial updates)
func (e *WorldEditor) renderTiles(coords []CubeCoord) error {
	if e.drawable == nil || e.currentWorld == nil || len(coords) == 0 {
		return nil // No canvas, map, or tiles to render
	}

	// For now, do a full render - we'll optimize for partial rendering later
	return e.renderFullMap()
}
