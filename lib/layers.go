package weewar

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Layer represents a single rendering layer (terrain, units, UI, etc.)
type Layer interface {
	// Core rendering
	Render(world *World, options LayerRenderOptions)

	// Dirty tracking for efficient updates
	MarkDirty(coord CubeCoord)
	MarkAllDirty()
	ClearDirty()
	IsDirty() bool

	// Lifecycle management
	SetViewPort(x, y, width, height int)
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
	ShowGrid        bool
	ShowCoordinates bool
}

// LayerScheduler interface for layers to request renders
type LayerScheduler interface {
	ScheduleRender()
}

// BaseLayer provides common functionality for all layers
type BaseLayer struct {
	name                string
	x, y, width, height int
	buffer              *Buffer

	// Dirty tracking
	dirtyCoords map[CubeCoord]bool
	allDirty    bool

	// Asset provider
	assetProvider AssetProvider

	// Renderer reference for scheduling
	scheduler LayerScheduler
}

// NewBaseLayer creates a new base layer
func NewBaseLayer(name string, width, height int, scheduler LayerScheduler) *BaseLayer {
	return &BaseLayer{
		name:        name,
		width:       width,
		height:      height,
		buffer:      NewBuffer(width, height),
		dirtyCoords: make(map[CubeCoord]bool),
		allDirty:    true, // Start with everything dirty
		scheduler:   scheduler,
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

func (bl *BaseLayer) IsDirty() bool {
	return bl.allDirty || len(bl.dirtyCoords) > 0
}

func (bl *BaseLayer) SetViewPort(x, y, width, height int) {
	bl.x = x
	bl.y = y
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

// drawImageToBuffer draws an image to the buffer
func (ul *BaseLayer) drawImageToBuffer(img image.Image, x, y, width, height float64) {
	bufferImg := ul.buffer.GetImageData()

	destRect := image.Rect(
		int(x-width/2), int(y-height/2),
		int(x+width/2), int(y+height/2),
	)

	draw.DrawMask(bufferImg, destRect, img, image.Point{}, nil, image.Point{}, draw.Over)
}

// drawSimpleHexToBuffer draws a colored hexagon
func (tl *BaseLayer) drawSimpleHexToBuffer(x, y float64, hexColor Color, options LayerRenderOptions) {
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
		BaseLayer:      NewBaseLayer("terrain", width, height, scheduler),
		terrainSprites: make(map[int]image.Image),
	}
}

// Render renders terrain tiles to the layer buffer
func (tl *TileLayer) Render(world *World, options LayerRenderOptions) {
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
func (tl *TileLayer) getTerrainColor(terrainType int) Color {
	switch terrainType {
	case 1: // Grass
		return Color{R: 0x22, G: 0x8B, B: 0x22}
	case 2: // Desert
		return Color{R: 0xEE, G: 0xCB, B: 0xAD}
	case 3: // Water
		return Color{R: 0x41, G: 0x69, B: 0xE1}
	case 4: // Mountain
		return Color{R: 0x8B, G: 0x89, B: 0x89}
	case 5: // Rock
		return Color{R: 0x69, G: 0x69, B: 0x69}
	default:
		return Color{R: 0xC8, G: 0xC8, B: 0xC8}
	}
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
		BaseLayer:   NewBaseLayer("units", width, height, scheduler),
		unitSprites: make(map[string]image.Image),
	}
}

// Render renders units to the layer buffer
func (ul *UnitLayer) Render(world *World, options LayerRenderOptions) {
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
