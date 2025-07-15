package weewar

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"
)

// =============================================================================
// Layer Interface and Base Types
// =============================================================================

// Layer represents a single rendering layer (terrain, units, UI, etc.)
type Layer interface {
	// Core rendering
	Render(world *World, drawable Drawable, options LayerRenderOptions)
	
	// Dirty tracking for efficient updates
	MarkDirty(coord CubeCoord)
	MarkAllDirty()
	ClearDirty()
	HasDirty() bool
	
	// Lifecycle management
	Resize(width, height int)
	SetAssetProvider(provider AssetProvider)
	
	// Layer identification
	GetName() string
}

// LayerRenderOptions contains rendering parameters for layers
type LayerRenderOptions struct {
	// Hex grid parameters
	TileWidth  float64
	TileHeight float64
	YIncrement float64
	
	// Viewport parameters
	ScrollX float64
	ScrollY float64
	
	// Visual options
	ShowGrid bool
	ShowCoordinates bool
}

// LayerScheduler interface for layers to request renders
type LayerScheduler interface {
	ScheduleRender()
}

// BaseLayer provides common functionality for all layers
type BaseLayer struct {
	name string
	width, height int
	buffer *Buffer
	
	// Dirty tracking
	dirtyCoords map[CubeCoord]bool
	allDirty bool
	
	// Asset provider
	assetProvider AssetProvider
	
	// Renderer reference for scheduling
	scheduler LayerScheduler
}

// NewBaseLayer creates a new base layer
func NewBaseLayer(name string, width, height int, scheduler LayerScheduler) *BaseLayer {
	return &BaseLayer{
		name: name,
		width: width,
		height: height,
		buffer: NewBuffer(width, height),
		dirtyCoords: make(map[CubeCoord]bool),
		allDirty: true, // Start with everything dirty
		scheduler: scheduler,
	}
}

// Common BaseLayer methods
func (bl *BaseLayer) GetName() string {
	return bl.name
}

func (bl *BaseLayer) MarkDirty(coord CubeCoord) {
	bl.dirtyCoords[coord] = true
	if bl.scheduler != nil {
		bl.scheduler.ScheduleRender()
	}
}

func (bl *BaseLayer) MarkAllDirty() {
	bl.allDirty = true
	bl.dirtyCoords = make(map[CubeCoord]bool)
	if bl.scheduler != nil {
		bl.scheduler.ScheduleRender()
	}
}

func (bl *BaseLayer) ClearDirty() {
	bl.dirtyCoords = make(map[CubeCoord]bool)
	bl.allDirty = false
}

func (bl *BaseLayer) HasDirty() bool {
	return bl.allDirty || len(bl.dirtyCoords) > 0
}

func (bl *BaseLayer) Resize(width, height int) {
	bl.width = width
	bl.height = height
	bl.buffer = NewBuffer(width, height)
	bl.MarkAllDirty()
}

func (bl *BaseLayer) SetAssetProvider(provider AssetProvider) {
	bl.assetProvider = provider
	bl.MarkAllDirty()
}

func (bl *BaseLayer) GetBuffer() *Buffer {
	return bl.buffer
}

// =============================================================================
// TileLayer - Terrain Rendering
// =============================================================================

// TileLayer handles rendering of terrain tiles
type TileLayer struct {
	*BaseLayer
	terrainSprites map[int]image.Image // Cached terrain sprites
}

// NewTileLayer creates a new tile layer
func NewTileLayer(width, height int, scheduler LayerScheduler) *TileLayer {
	return &TileLayer{
		BaseLayer: NewBaseLayer("terrain", width, height, scheduler),
		terrainSprites: make(map[int]image.Image),
	}
}

// Render renders terrain tiles to the layer buffer
func (tl *TileLayer) Render(world *World, drawable Drawable, options LayerRenderOptions) {
	if world == nil || world.Map == nil {
		return
	}
	
	// Clear buffer if full rebuild needed
	if tl.allDirty {
		tl.buffer.Clear()
		
		// Render all tiles
		for coord, tile := range world.Map.Tiles {
			if tile != nil {
				tl.renderTile(world, coord, tile, options)
			}
		}
		
		tl.allDirty = false
	} else {
		// Render only dirty tiles
		for coord := range tl.dirtyCoords {
			tile := world.Map.TileAt(coord)
			tl.renderTile(world, coord, tile, options)
		}
	}
	
	// Clear dirty tracking
	tl.ClearDirty()
}

// renderTile renders a single terrain tile
func (tl *TileLayer) renderTile(world *World, coord CubeCoord, tile *Tile, options LayerRenderOptions) {
	if tile == nil {
		return
	}
	
	// Get pixel position using Map's coordinate system
	x, y := world.Map.CenterXYForTile(coord, options.TileWidth, options.TileHeight, options.YIncrement)
	
	// Apply viewport offset
	x += options.ScrollX
	y += options.ScrollY
	
	// Try to use real terrain sprite if available
	if tl.assetProvider != nil && tl.assetProvider.HasTileAsset(tile.TileType) {
		tl.renderTerrainSprite(tile.TileType, x, y, options)
	} else {
		// Fallback to colored hexagon
		color := tl.getTerrainColor(tile.TileType)
		tl.drawSimpleHexToBuffer(x, y, color, options)
	}
}

// renderTerrainSprite renders a terrain sprite
func (tl *TileLayer) renderTerrainSprite(tileType int, x, y float64, options LayerRenderOptions) {
	// Check cache first
	cachedSprite, exists := tl.terrainSprites[tileType]
	if !exists {
		// Load and cache sprite
		img, err := tl.assetProvider.GetTileImage(tileType)
		if err != nil {
			// Fallback to colored hex
			color := tl.getTerrainColor(tileType)
			tl.drawSimpleHexToBuffer(x, y, color, options)
			return
		}
		tl.terrainSprites[tileType] = img
		cachedSprite = img
	}
	
	// Draw sprite to buffer
	tl.drawImageToBuffer(cachedSprite, x, y, options.TileWidth, options.TileHeight)
}

// getTerrainColor returns color for terrain type
func (tl *TileLayer) getTerrainColor(terrainType int) string {
	switch terrainType {
	case 1: // Grass
		return "#228B22"
	case 2: // Desert
		return "#EECBAD"
	case 3: // Water
		return "#4169E1"
	case 4: // Mountain
		return "#8B8989"
	case 5: // Rock
		return "#696969"
	default:
		return "#C8C8C8"
	}
}

// drawSimpleHexToBuffer draws a colored hexagon
func (tl *TileLayer) drawSimpleHexToBuffer(x, y float64, colorStr string, options LayerRenderOptions) {
	hexColor := tl.parseHexColor(colorStr)
	bufferImg := tl.buffer.GetImageData()
	
	// Draw ellipse approximation
	radiusX := int(options.TileWidth / 2)
	radiusY := int(options.TileHeight / 2)
	centerX, centerY := int(x), int(y)
	
	for dy := -radiusY; dy <= radiusY; dy++ {
		for dx := -radiusX; dx <= radiusX; dx++ {
			if float64(dx*dx)/float64(radiusX*radiusX)+float64(dy*dy)/float64(radiusY*radiusY) <= 1.0 {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && py >= 0 && px < tl.width && py < tl.height {
					rgba := color.RGBA{R: hexColor.R, G: hexColor.G, B: hexColor.B, A: hexColor.A}
					bufferImg.Set(px, py, rgba)
				}
			}
		}
	}
}

// drawImageToBuffer draws an image to the buffer
func (tl *TileLayer) drawImageToBuffer(img image.Image, x, y, width, height float64) {
	bufferImg := tl.buffer.GetImageData()
	
	destRect := image.Rect(
		int(x-width/2), int(y-height/2),
		int(x+width/2), int(y+height/2),
	)
	
	draw.DrawMask(bufferImg, destRect, img, image.Point{}, nil, image.Point{}, draw.Over)
}

// parseHexColor converts hex color string to Color
func (tl *TileLayer) parseHexColor(hexColor string) Color {
	if len(hexColor) > 0 && hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}
	
	if len(hexColor) != 6 {
		return Color{R: 34, G: 139, B: 34, A: 255}
	}
	
	var red, green, blue uint8
	fmt.Sscanf(hexColor[0:2], "%02x", &red)
	fmt.Sscanf(hexColor[2:4], "%02x", &green)
	fmt.Sscanf(hexColor[4:6], "%02x", &blue)
	
	return Color{R: red, G: green, B: blue, A: 255}
}

// =============================================================================
// UnitLayer - Unit Rendering
// =============================================================================

// UnitLayer handles rendering of units
type UnitLayer struct {
	*BaseLayer
	unitSprites map[string]image.Image // Cached unit sprites by "type_player"
}

// NewUnitLayer creates a new unit layer
func NewUnitLayer(width, height int, scheduler LayerScheduler) *UnitLayer {
	return &UnitLayer{
		BaseLayer: NewBaseLayer("units", width, height, scheduler),
		unitSprites: make(map[string]image.Image),
	}
}

// Render renders units to the layer buffer
func (ul *UnitLayer) Render(world *World, drawable Drawable, options LayerRenderOptions) {
	if world == nil {
		return
	}
	
	// Clear buffer if full rebuild needed
	if ul.allDirty {
		ul.buffer.Clear()
		
		// Render all units from all players
		for _, playerUnits := range world.UnitsByPlayer {
			for _, unit := range playerUnits {
				if unit != nil {
					ul.renderUnit(world, unit, options)
				}
			}
		}
		
		ul.allDirty = false
	} else {
		// Clear and render only dirty unit positions
		for coord := range ul.dirtyCoords {
			ul.clearHexArea(coord, options)
			
			// Find unit at this position
			unit := ul.findUnitAt(world, coord)
			if unit != nil {
				ul.renderUnit(world, unit, options)
			}
		}
	}
	
	// Clear dirty tracking
	ul.ClearDirty()
}

// renderUnit renders a single unit
func (ul *UnitLayer) renderUnit(world *World, unit *Unit, options LayerRenderOptions) {
	// Get pixel position using Map's coordinate system
	x, y := world.Map.CenterXYForTile(unit.Coord, options.TileWidth, options.TileHeight, options.YIncrement)
	
	// Apply viewport offset
	x += options.ScrollX
	y += options.ScrollY
	
	// Try to use real unit sprite if available
	if ul.assetProvider != nil && ul.assetProvider.HasUnitAsset(unit.UnitType, unit.PlayerID) {
		ul.renderUnitSprite(unit.UnitType, unit.PlayerID, x, y, options)
	} else {
		// Fallback to colored circle
		ul.drawSimpleUnitToBuffer(x, y, unit.PlayerID, options)
	}
}

// renderUnitSprite renders a unit sprite
func (ul *UnitLayer) renderUnitSprite(unitType, playerID int, x, y float64, options LayerRenderOptions) {
	// Check cache first
	spriteKey := fmt.Sprintf("%d_%d", unitType, playerID)
	cachedSprite, exists := ul.unitSprites[spriteKey]
	if !exists {
		// Load and cache sprite
		img, err := ul.assetProvider.GetUnitImage(unitType, playerID)
		if err != nil {
			// Fallback to colored circle
			ul.drawSimpleUnitToBuffer(x, y, playerID, options)
			return
		}
		ul.unitSprites[spriteKey] = img
		cachedSprite = img
	}
	
	// Draw sprite to buffer
	ul.drawImageToBuffer(cachedSprite, x, y, options.TileWidth, options.TileHeight)
}

// drawSimpleUnitToBuffer draws a colored circle for a unit
func (ul *UnitLayer) drawSimpleUnitToBuffer(x, y float64, playerID int, options LayerRenderOptions) {
	// Get player color
	var unitColor Color
	switch playerID {
	case 0:
		unitColor = Color{R: 255, G: 0, B: 0, A: 255} // Red
	case 1:
		unitColor = Color{R: 0, G: 0, B: 255, A: 255} // Blue
	case 2:
		unitColor = Color{R: 0, G: 255, B: 0, A: 255} // Green
	case 3:
		unitColor = Color{R: 255, G: 255, B: 0, A: 255} // Yellow
	default:
		unitColor = Color{R: 128, G: 128, B: 128, A: 255} // Gray
	}
	
	bufferImg := ul.buffer.GetImageData()
	
	// Draw smaller ellipse for units (60% of tile size)
	radiusX := int(options.TileWidth * 0.3)
	radiusY := int(options.TileHeight * 0.3)
	centerX, centerY := int(x), int(y)
	
	for dy := -radiusY; dy <= radiusY; dy++ {
		for dx := -radiusX; dx <= radiusX; dx++ {
			if float64(dx*dx)/float64(radiusX*radiusX)+float64(dy*dy)/float64(radiusY*radiusY) <= 1.0 {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && py >= 0 && px < ul.width && py < ul.height {
					rgba := color.RGBA{R: unitColor.R, G: unitColor.G, B: unitColor.B, A: unitColor.A}
					bufferImg.Set(px, py, rgba)
				}
			}
		}
	}
}

// drawImageToBuffer draws an image to the buffer
func (ul *UnitLayer) drawImageToBuffer(img image.Image, x, y, width, height float64) {
	bufferImg := ul.buffer.GetImageData()
	
	destRect := image.Rect(
		int(x-width/2), int(y-height/2),
		int(x+width/2), int(y+height/2),
	)
	
	draw.DrawMask(bufferImg, destRect, img, image.Point{}, nil, image.Point{}, draw.Over)
}

// clearHexArea clears a hexagonal area at the given coordinate
func (ul *UnitLayer) clearHexArea(coord CubeCoord, options LayerRenderOptions) {
	// For now, just clear the entire buffer - can optimize later
	ul.buffer.Clear()
}

// findUnitAt finds a unit at the given coordinate
func (ul *UnitLayer) findUnitAt(world *World, coord CubeCoord) *Unit {
	for _, playerUnits := range world.UnitsByPlayer {
		for _, unit := range playerUnits {
			if unit != nil && unit.Coord == coord {
				return unit
			}
		}
	}
	return nil
}

// =============================================================================
// LayeredRenderer - Coordinates Multiple Layers
// =============================================================================

// LayeredRenderer coordinates multiple rendering layers
type LayeredRenderer struct {
	// Canvas target
	canvasID string
	width    int
	height   int

	// Rendering layers (in order)
	layers []Layer
	
	// Output buffer for compositing
	outputBuffer *Buffer

	// Batching system
	batchTimer    *time.Timer
	batchInterval time.Duration
	renderPending bool

	// Rendering parameters
	renderOptions LayerRenderOptions

	// Current world reference
	currentWorld *World
}

// NewLayeredRenderer creates a new layered renderer with default tile dimensions
func NewLayeredRenderer(canvasID string, width, height int) (*LayeredRenderer, error) {
	return NewLayeredRendererWithTileSize(canvasID, width, height, DefaultTileWidth, DefaultTileHeight, DefaultYIncrement)
}

// NewLayeredRendererWithTileSize creates a new layered renderer with specified tile dimensions
func NewLayeredRendererWithTileSize(canvasID string, width, height int, tileWidth, tileHeight, yIncrement float64) (*LayeredRenderer, error) {
	// Create output buffer for compositing
	outputBuffer := NewBuffer(width, height)
	outputBuffer.Clear()

	renderer := &LayeredRenderer{
		canvasID:      canvasID,
		width:         width,
		height:        height,
		outputBuffer:  outputBuffer,
		batchInterval: 30 * time.Millisecond, // 33 FPS for prototyping
		renderPending: false,
		renderOptions: LayerRenderOptions{
			TileWidth:  tileWidth,
			TileHeight: tileHeight,
			YIncrement: yIncrement,
			ScrollX:    0,
			ScrollY:    0,
			ShowGrid:   false,
			ShowCoordinates: false,
		},
	}

	// Create layers in rendering order (pass renderer as scheduler)
	renderer.layers = []Layer{
		NewTileLayer(width, height, renderer),  // Terrain tiles (bottom layer)
		NewUnitLayer(width, height, renderer),  // Units (middle layer)
		// TODO: Add HighlightLayer, UILayer, etc.
	}

	return renderer, nil
}

// SetWorld updates the current world reference
func (r *LayeredRenderer) SetWorld(w *World) {
	r.currentWorld = w
	// Mark all layers as dirty when world changes
	for _, layer := range r.layers {
		layer.MarkAllDirty()
	}
}

// SetAssetProvider updates the asset provider for all layers
func (r *LayeredRenderer) SetAssetProvider(provider AssetProvider) {
	for _, layer := range r.layers {
		layer.SetAssetProvider(provider)
	}
}

// SetTileDimensions updates the tile rendering dimensions
func (r *LayeredRenderer) SetTileDimensions(tileWidth, tileHeight, yIncrement float64) {
	r.renderOptions.TileWidth = tileWidth
	r.renderOptions.TileHeight = tileHeight
	r.renderOptions.YIncrement = yIncrement
	
	// Mark all layers as dirty since dimensions changed
	for _, layer := range r.layers {
		layer.MarkAllDirty()
	}
}

// SetScroll updates the viewport scroll offset
func (r *LayeredRenderer) SetScroll(scrollX, scrollY float64) {
	r.renderOptions.ScrollX = scrollX
	r.renderOptions.ScrollY = scrollY
	
	// Mark all layers as dirty since viewport changed
	for _, layer := range r.layers {
		layer.MarkAllDirty()
	}
}

// ScheduleRender allows layers to request a render update
func (r *LayeredRenderer) ScheduleRender() {
	r.scheduleRender()
}

// GetLayerByName returns a layer by its name (for external access)
func (r *LayeredRenderer) GetLayerByName(name string) Layer {
	for _, layer := range r.layers {
		if layer.GetName() == name {
			return layer
		}
	}
	return nil
}

// SetTileDimensions updates the tile rendering dimensions
func (r *LayeredRenderer) SetTileDimensions(tileWidth, tileHeight, yIncrement float64) {
	r.tileWidth = tileWidth
	r.tileHeight = tileHeight
	r.yIncrement = yIncrement
	// Mark all terrain as dirty since dimensions changed
	r.MarkAllTerrainDirty()
	r.MarkAllUnitsDirty()
}

// MarkTerrainDirty marks a specific tile as needing terrain update
func (r *LayeredRenderer) MarkTerrainDirty(coord CubeCoord) {
	r.dirtyTerrain[coord] = true
	r.scheduleRender()
}

// MarkUnitDirty marks a specific position as needing unit update
func (r *LayeredRenderer) MarkUnitDirty(coord CubeCoord) {
	r.dirtyUnits[coord] = true
	r.scheduleRender()
}

// MarkUIDirty marks the UI layer as needing update
func (r *LayeredRenderer) MarkUIDirty() {
	r.dirtyUI = true
	r.scheduleRender()
}

// MarkAllTerrainDirty marks entire terrain layer for rebuild
func (r *LayeredRenderer) MarkAllTerrainDirty() {
	// Clear and mark for full rebuild
	for coord := range r.dirtyTerrain {
		delete(r.dirtyTerrain, coord)
	}
	r.dirtyTerrain[CubeCoord{Q: -999999, R: -999999}] = true // Special marker for "rebuild all"
	r.scheduleRender()
}

// MarkAllUnitsDirty marks all units in the current world as dirty
func (r *LayeredRenderer) MarkAllUnitsDirty() {
	if r.currentWorld == nil {
		return
	}

	// Clear existing dirty units
	for coord := range r.dirtyUnits {
		delete(r.dirtyUnits, coord)
	}

	// Mark all units from all players as dirty
	for _, playerUnits := range r.currentWorld.UnitsByPlayer {
		for _, unit := range playerUnits {
			if unit != nil {
				r.dirtyUnits[unit.Coord] = true
				fmt.Printf("MarkAllUnitsDirty: Marked unit at coord %v as dirty\n", unit.Coord)
			}
		}
	}

	fmt.Printf("MarkAllUnitsDirty: Total %d units marked as dirty\n", len(r.dirtyUnits))
	r.scheduleRender()
}

// scheduleRender schedules a batched render update
func (r *LayeredRenderer) scheduleRender() {
	if r.renderPending {
		return // Already scheduled
	}

	r.renderPending = true

	// Cancel existing timer
	if r.batchTimer != nil {
		r.batchTimer.Stop()
	}

	// Schedule new render
	r.batchTimer = time.AfterFunc(r.batchInterval, func() {
		r.performRender()
		r.renderPending = false
	})
}

// ForceRender immediately renders all dirty layers (for synchronous updates)
func (r *LayeredRenderer) ForceRender() {
	fmt.Printf("LayeredRenderer.ForceRender called - terrain dirty: %d, units dirty: %d, UI dirty: %v\n",
		len(r.dirtyTerrain), len(r.dirtyUnits), r.dirtyUI)

	// Debug: List the dirty units
	if len(r.dirtyUnits) > 0 {
		fmt.Printf("Dirty units: ")
		for coord := range r.dirtyUnits {
			fmt.Printf("%v ", coord)
		}
		fmt.Printf("\n")
	}

	if r.batchTimer != nil {
		r.batchTimer.Stop()
	}
	r.performRender()
	r.renderPending = false
	fmt.Printf("DEBUG: ForceRender() completed successfully\n")
}

// performRender executes the actual rendering of dirty layers
func (r *LayeredRenderer) performRender() {
	fmt.Printf("LayeredRenderer.performRender called\n")

	// Update terrain layer if dirty
	if len(r.dirtyTerrain) > 0 {
		fmt.Printf("Updating terrain layer with %d dirty tiles\n", len(r.dirtyTerrain))
		r.updateTerrainLayer()
	}

	// Update unit layer if dirty
	if len(r.dirtyUnits) > 0 {
		fmt.Printf("Updating unit layer with %d dirty positions\n", len(r.dirtyUnits))
		r.updateUnitLayer()
	}

	// Update UI layer if dirty
	if r.dirtyUI {
		fmt.Printf("Updating UI layer\n")
		r.updateUILayer()
		r.dirtyUI = false
	}

	// Composite all layers to main canvas
	fmt.Printf("Compositing layers to main canvas\n")
	r.composite()
	fmt.Printf("DEBUG: performRender() completed successfully\n")
}

// updateTerrainLayer renders dirty terrain tiles
func (r *LayeredRenderer) updateTerrainLayer() {
	// Check if full rebuild is needed
	_, fullRebuild := r.dirtyTerrain[CubeCoord{Q: -999999, R: -999999}]
	_, renderAllVisible := r.dirtyTerrain[CubeCoord{Q: -999998, R: -999998}]

	if fullRebuild || renderAllVisible {
		// Clear entire terrain buffer
		r.terrainBuffer.Clear()

		// Render all tiles in the current map
		if r.currentMap != nil {
			fmt.Printf("Rendering all %d tiles in map\n", len(r.currentMap.Tiles))
			for coord := range r.currentMap.Tiles {
				r.renderTerrainTile(coord)
			}
		}

		// Clear both markers
		delete(r.dirtyTerrain, CubeCoord{Q: -999999, R: -999999})
		delete(r.dirtyTerrain, CubeCoord{Q: -999998, R: -999998})
	} else {
		// Render individual dirty tiles
		for coord := range r.dirtyTerrain {
			r.renderTerrainTile(coord)
			delete(r.dirtyTerrain, coord)
		}
	}
}

// renderTerrainTile renders a single terrain tile using cached sprites
func (r *LayeredRenderer) renderTerrainTile(coord CubeCoord) {
	if r.currentMap == nil {
		return
	}

	// Get tile from current map
	tile := r.currentMap.TileAtCube(coord)

	// Calculate pixel position from hex coordinate
	x, y := r.hexToPixel(coord)

	if tile != nil {
		// Try to use real terrain sprite if asset provider is available
		if r.assetProvider != nil && r.assetProvider.HasTileAsset(tile.TileType) {
			// fmt.Printf("Rendering terrain sprite for tile type %d at (%f, %f)\n", tile.TileType, x, y)
			r.renderTerrainSprite(coord, tile.TileType, x, y)
		} else {
			// Fallback to colored hexagon using buffer operations
			if r.assetProvider == nil {
				fmt.Printf("AssetProvider is nil, using colored hex for tile type %d at (%f, %f)\n", tile.TileType, x, y)
			} else {
				fmt.Printf("AssetProvider.HasTileAsset(%d) returned false, using colored hex at (%f, %f)\n", tile.TileType, x, y)
			}
			color := r.getTerrainColor(tile.TileType)
			r.drawSimpleHexToBuffer(r.terrainBuffer, x, y, color)
		}
	}
}

// renderTerrainSprite renders a terrain sprite at the given position
func (r *LayeredRenderer) renderTerrainSprite(coord CubeCoord, tileType int, x, y float64) {
	// Check if we have a cached sprite for this tile type
	cachedSprite, exists := r.terrainSprites[tileType]
	if !exists {
		// Load and cache the terrain sprite
		img, err := r.assetProvider.GetTileImage(tileType)
		if err != nil {
			fmt.Printf("Failed to load terrain sprite for type %d: %v\n", tileType, err)
			// Fallback to colored hex
			color := r.getTerrainColor(tileType)
			r.drawSimpleHexToBuffer(r.terrainBuffer, x, y, color)
			return
		}

		// Debug: Check image properties
		bounds := img.Bounds()
		fmt.Printf("Loaded terrain sprite for type %d: size %dx%d, bounds %v\n", tileType, bounds.Dx(), bounds.Dy(), bounds)

		// Cache the image directly
		r.terrainSprites[tileType] = img
		cachedSprite = img
	}

	// Draw the sprite to the terrain buffer with proper alpha blending
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC in drawImageToBuffer: %v\n", r)
		}
	}()
	r.drawImageToBuffer(r.terrainBuffer, cachedSprite, x, y, r.tileWidth, r.tileHeight)
}

// getTerrainColor returns the color for a terrain type
func (r *LayeredRenderer) getTerrainColor(terrainType int) string {
	switch terrainType {
	case 1: // Grass
		return "#228B22"
	case 2: // Desert
		return "#EECBAD"
	case 3: // Water
		return "#4169E1"
	case 4: // Mountain
		return "#8B8989"
	case 5: // Rock
		return "#696969"
	default:
		return "#C8C8C8"
	}
}

// updateUnitLayer renders dirty unit positions
func (r *LayeredRenderer) updateUnitLayer() {
	fmt.Printf("updateUnitLayer called with %d dirty units\n", len(r.dirtyUnits))
	// Clear dirty areas and redraw units using buffer operations
	for coord := range r.dirtyUnits {
		fmt.Printf("Processing dirty unit at coord %v\n", coord)
		// Clear the specific hex area in unitBuffer first
		r.clearHexArea(r.unitBuffer, coord)

		// Get unit at this position from current map
		if r.currentMap != nil {
			tile := r.currentMap.TileAtCube(coord)
			if tile != nil && tile.Unit != nil {
				// Render unit sprite to unitBuffer
				r.renderUnitSprite(coord, tile.Unit)
			}
		}

		delete(r.dirtyUnits, coord)
	}
}

// updateUILayer renders UI elements (selection, hover, etc.)
func (r *LayeredRenderer) updateUILayer() {
	// Clear entire UI buffer
	r.uiBuffer.Clear()

	// TODO: Render current selection highlight to uiBuffer
	// TODO: Render hover highlight to uiBuffer
	// TODO: Render range indicators to uiBuffer, etc.
}

// composite just marks that layers need to be blitted
func (r *LayeredRenderer) composite() {
	// No complex compositing - just signal that buffers are ready for blitting
}

// GetTerrainBuffer returns the terrain buffer for external blitting
func (r *LayeredRenderer) GetTerrainBuffer() *Buffer {
	return r.terrainBuffer
}

// GetUnitBuffer returns the unit buffer for external blitting
func (r *LayeredRenderer) GetUnitBuffer() *Buffer {
	return r.unitBuffer
}

// GetUIBuffer returns the UI buffer for external blitting
func (r *LayeredRenderer) GetUIBuffer() *Buffer {
	return r.uiBuffer
}

// blendBuffers blends src buffer onto dst buffer with alpha blending
func (r *LayeredRenderer) blendBuffers(dst, src *Buffer) {
	dstImg := dst.GetImageData()
	srcImg := src.GetImageData()

	// Use Go's image/draw for proper alpha blending
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Over)
}

// drawImageToBuffer draws an image to a buffer with proper alpha blending
func (r *LayeredRenderer) drawImageToBuffer(buffer *Buffer, img image.Image, x, y, width, height float64) {
	// Get the buffer's underlying image
	bufferImg := buffer.GetImageData()

	// Calculate destination rectangle (centered on x,y)
	destRect := image.Rect(
		int(x-width/2),
		int(y-height/2),
		int(x+width/2),
		int(y+height/2),
	)

	// fmt.Printf("DrawImageToBuffer: Drawing image at (%f,%f) size %fx%f, destRect %v, img bounds %v\n", x, y, width, height, destRect, img.Bounds())

	// Resize source image to match destination size if needed
	srcBounds := img.Bounds()
	if srcBounds.Dx() != int(width) || srcBounds.Dy() != int(height) {
		// fmt.Printf("Resizing image from %dx%d to %dx%d\n", srcBounds.Dx(), srcBounds.Dy(), int(width), int(height))
		// Create a resized version
		resized := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		// Simple nearest-neighbor scaling for now
		if int(width) > 0 && int(height) > 0 && srcBounds.Dx() > 0 && srcBounds.Dy() > 0 {
			for dy := 0; dy < int(height); dy++ {
				for dx := 0; dx < int(width); dx++ {
					srcX := dx * srcBounds.Dx() / int(width)
					srcY := dy * srcBounds.Dy() / int(height)
					resized.Set(dx, dy, img.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
				}
			}
		}
		img = resized
		// fmt.Printf("Image resizing completed\n")
	}

	// Draw the image using Go's image/draw with alpha blending
	draw.DrawMask(bufferImg, destRect, img, image.Point{}, nil, image.Point{}, draw.Over)
	// fmt.Printf("Image drawing completed successfully\n")
}

// drawSimpleHexToBuffer draws a simple colored hexagon to a buffer
func (r *LayeredRenderer) drawSimpleHexToBuffer(buffer *Buffer, x, y float64, colorStr string) {
	// Convert hex color string to Color struct
	hexColor := r.parseHexColor(colorStr)

	// Get the buffer's underlying image
	bufferImg := buffer.GetImageData()

	// Draw a simple filled ellipse as a placeholder for hexagon
	// Use tile dimensions to determine the shape
	radiusX := int(r.tileWidth / 2)
	radiusY := int(r.tileHeight / 2)
	centerX, centerY := int(x), int(y)

	fmt.Printf("Drawing simple hex at (%d, %d) with radii %dx%d, color %s (%d,%d,%d,%d)\n",
		centerX, centerY, radiusX, radiusY, colorStr, hexColor.R, hexColor.G, hexColor.B, hexColor.A)

	for dy := -radiusY; dy <= radiusY; dy++ {
		for dx := -radiusX; dx <= radiusX; dx++ {
			// Ellipse equation: (x/a)² + (y/b)² <= 1
			if float64(dx*dx)/float64(radiusX*radiusX)+float64(dy*dy)/float64(radiusY*radiusY) <= 1.0 {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && py >= 0 && px < r.width && py < r.height {
					// Convert our Color to color.RGBA
					rgba := color.RGBA{R: hexColor.R, G: hexColor.G, B: hexColor.B, A: hexColor.A}
					bufferImg.Set(px, py, rgba)
				}
			}
		}
	}
}

// parseHexColor converts a hex color string like "#228B22" to Color
func (r *LayeredRenderer) parseHexColor(hexColor string) Color {
	// Remove # if present
	if len(hexColor) > 0 && hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	// Default to green if parsing fails
	if len(hexColor) != 6 {
		return Color{R: 34, G: 139, B: 34, A: 255}
	}

	// Parse RGB components
	var red, green, blue uint8
	fmt.Sscanf(hexColor[0:2], "%02x", &red)
	fmt.Sscanf(hexColor[2:4], "%02x", &green)
	fmt.Sscanf(hexColor[4:6], "%02x", &blue)

	return Color{R: red, G: green, B: blue, A: 255}
}

// hexToPixel converts hex coordinates to pixel coordinates using the same logic as game.go
func (r *LayeredRenderer) hexToPixel(coord CubeCoord) (float64, float64) {
	// Use the same conversion as game.go XYForTile - convert to display coordinates first
	row := coord.R
	col := coord.Q + (coord.R+(coord.R&1))/2

	// Use the exact same calculation as game.go XYForTile
	x := float64(col)*r.tileWidth + r.tileWidth/2

	// Apply offset for alternating rows (hex grid staggering)
	isEvenRow := (row % 2) == 0
	// Assuming odd rows are offset (EvenRowsOffset() returns false)
	if !isEvenRow {
		x += r.tileWidth / 2
	}

	y := float64(row)*r.yIncrement + r.tileHeight/2

	return x, y
}

// clearHexArea clears a hexagonal area in the buffer at the given coordinate
func (r *LayeredRenderer) clearHexArea(buffer *Buffer, coord CubeCoord) {
	// Calculate pixel position
	x, y := r.hexToPixel(coord)

	// Get the buffer's underlying image
	bufferImg := buffer.GetImageData()

	// Clear an elliptical area (approximate hex area using tile dimensions)
	radiusX := int(r.tileWidth / 2)
	radiusY := int(r.tileHeight / 2)
	centerX, centerY := int(x), int(y)

	transparentColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	for dy := -radiusY; dy <= radiusY; dy++ {
		for dx := -radiusX; dx <= radiusX; dx++ {
			// Ellipse equation: (x/a)² + (y/b)² <= 1
			if float64(dx*dx)/float64(radiusX*radiusX)+float64(dy*dy)/float64(radiusY*radiusY) <= 1.0 {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && py >= 0 && px < r.width && py < r.height {
					bufferImg.Set(px, py, transparentColor)
				}
			}
		}
	}
}

// renderUnitSprite renders a unit sprite at the given coordinate
func (r *LayeredRenderer) renderUnitSprite(coord CubeCoord, unit *Unit) {
	// Calculate pixel position
	x, y := r.hexToPixel(coord)

	// Try to use real unit sprite if asset provider is available
	if r.assetProvider != nil && r.assetProvider.HasUnitAsset(unit.UnitType, unit.PlayerID) {
		// Check if we have a cached sprite for this unit type and player
		spriteKey := fmt.Sprintf("%d_%d", unit.UnitType, unit.PlayerID)
		cachedSprite, exists := r.unitSprites[spriteKey]
		if !exists {
			// Load and cache the unit sprite
			img, err := r.assetProvider.GetUnitImage(unit.UnitType, unit.PlayerID)
			if err != nil {
				fmt.Printf("Failed to load unit sprite for type %d, player %d: %v\n", unit.UnitType, unit.PlayerID, err)
				// Fallback to simple colored circle
				r.drawSimpleUnitToBuffer(r.unitBuffer, x, y, unit.PlayerID)
				return
			}

			// Cache the image
			r.unitSprites[spriteKey] = img
			cachedSprite = img
		}

		// Draw the sprite to the unit buffer
		fmt.Printf("Drawing unit sprite at position (%f, %f) with tileDimensions %fx%f, sprite bounds: %v\n",
			x, y, r.tileWidth, r.tileHeight, cachedSprite.Bounds())
		r.drawImageToBuffer(r.unitBuffer, cachedSprite, x, y, r.tileWidth, r.tileHeight)
	} else {
		// Fallback to simple colored circle
		fmt.Printf("Asset provider doesn't have unit asset, falling back to simple circle\n")
		r.drawSimpleUnitToBuffer(r.unitBuffer, x, y, unit.PlayerID)
	}
}

// drawSimpleUnitToBuffer draws a simple colored ellipse to represent a unit
func (r *LayeredRenderer) drawSimpleUnitToBuffer(buffer *Buffer, x, y float64, playerID int) {
	// Get player color
	var unitColor Color
	switch playerID {
	case 0:
		unitColor = Color{R: 255, G: 0, B: 0, A: 255} // Red
	case 1:
		unitColor = Color{R: 0, G: 0, B: 255, A: 255} // Blue
	default:
		unitColor = Color{R: 128, G: 128, B: 128, A: 255} // Gray
	}

	// Get the buffer's underlying image
	bufferImg := buffer.GetImageData()

	// Draw a smaller ellipse for units (60% of tile dimensions)
	radiusX := int(r.tileWidth * 0.3)  // 60% of half-width = 30% of full width
	radiusY := int(r.tileHeight * 0.3) // 60% of half-height = 30% of full height
	centerX, centerY := int(x), int(y)

	fmt.Printf("Drawing simple unit ellipse at (%d, %d) with radii %dx%d, player %d color (%d,%d,%d)\n",
		centerX, centerY, radiusX, radiusY, playerID, unitColor.R, unitColor.G, unitColor.B)

	for dy := -radiusY; dy <= radiusY; dy++ {
		for dx := -radiusX; dx <= radiusX; dx++ {
			// Ellipse equation: (x/a)² + (y/b)² <= 1
			if float64(dx*dx)/float64(radiusX*radiusX)+float64(dy*dy)/float64(radiusY*radiusY) <= 1.0 {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && py >= 0 && px < r.width && py < r.height {
					rgba := color.RGBA{R: unitColor.R, G: unitColor.G, B: unitColor.B, A: unitColor.A}
					bufferImg.Set(px, py, rgba)
				}
			}
		}
	}
}

// Resize updates the layer buffer sizes
func (r *LayeredRenderer) Resize(width, height int) error {
	r.width = width
	r.height = height

	// Recreate all layer buffers with new size
	r.terrainBuffer = NewBuffer(width, height)
	r.unitBuffer = NewBuffer(width, height)
	r.uiBuffer = NewBuffer(width, height)

	// Clear all buffers to transparent
	r.terrainBuffer.Clear()
	r.unitBuffer.Clear()
	r.uiBuffer.Clear()

	// Mark everything as dirty for redraw
	r.MarkAllTerrainDirty()
	r.MarkUIDirty()

	return nil
}
