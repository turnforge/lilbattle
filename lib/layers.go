package weewar

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
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
func (tl *BaseLayer) parseHexColor(hexColor string) Color {
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
// GridLayer - Grid Line and Coordinate Rendering
// =============================================================================

// NewGridLayer creates a new grid layer
func NewGridLayer(width, height int, scheduler LayerScheduler) *GridLayer {
	return &GridLayer{
		BaseLayer: NewBaseLayer("grid", width, height, scheduler),
	}
}

// Render renders hex grid lines and coordinates
func (gl *GridLayer) Render(world *World, options LayerRenderOptions) {
	if world == nil || world.Map == nil {
		return
	}

	// Only render if grid or coordinates are enabled
	if !options.ShowGrid && !options.ShowCoordinates {
		if !gl.allDirty {
			return // Nothing to render
		}
		// Clear buffer if switching from visible to hidden
		gl.buffer.Clear()
		gl.allDirty = false
		gl.ClearDirty()
		return
	}

	// Clear buffer for full redraw (grid/coordinates are view-dependent)
	gl.buffer.Clear()

	// Get all tiles in the map
	tiles := world.Map.Tiles

	for coord, tile := range tiles {
		if tile == nil {
			continue
		}

		// Get pixel position for this tile
		x, y := world.Map.CenterXYForTile(coord, options.TileWidth, options.TileHeight, options.YIncrement)

		// Apply viewport offset
		x += options.ScrollX
		y += options.ScrollY

		// Check if tile is within visible area
		if x < -options.TileWidth || x > float64(gl.width)+options.TileWidth ||
			y < -options.TileHeight || y > float64(gl.height)+options.TileHeight {
			continue
		}

		// Draw grid lines if enabled
		if options.ShowGrid {
			gl.drawHexGrid(x, y, options)
		}

		// Draw coordinates if enabled
		if options.ShowCoordinates {
			gl.drawCoordinates(coord, x, y, options)
		}
	}

	// Mark as clean
	gl.allDirty = false
	gl.ClearDirty()
}

// drawHexGrid draws hexagonal grid lines around a tile
func (gl *GridLayer) drawHexGrid(centerX, centerY float64, options LayerRenderOptions) {
	// Get hexagon vertices
	vertices := gl.getHexVertices(centerX, centerY, options.TileWidth, options.TileHeight)

	// Draw lines between vertices
	gridColor := color.RGBA{R: 64, G: 64, B: 64, A: 255} // Dark gray
	bufferImg := gl.buffer.GetImageData()

	for i := 0; i < len(vertices); i++ {
		x1, y1 := vertices[i][0], vertices[i][1]
		x2, y2 := vertices[(i+1)%len(vertices)][0], vertices[(i+1)%len(vertices)][1]

		gl.drawLine(bufferImg, int(x1), int(y1), int(x2), int(y2), gridColor)
	}
}

// drawCoordinates draws Q,R coordinates in the center of a hex
func (gl *GridLayer) drawCoordinates(coord CubeCoord, centerX, centerY float64, options LayerRenderOptions) {
	// Simple text rendering - draw coordinate text
	text := fmt.Sprintf("%d,%d", coord.Q, coord.R)

	// For now, draw a simple representation (can be enhanced with proper text rendering)
	gl.drawSimpleText(text, centerX, centerY)
}

// getHexVertices returns the vertices of a hexagon centered at (centerX, centerY)
func (gl *GridLayer) getHexVertices(centerX, centerY, tileWidth, tileHeight float64) [][2]float64 {
	// Hexagon vertices (flat-top orientation)
	vertices := make([][2]float64, 6)

	// Use actual tile dimensions for proper hexagon shape
	radiusX := tileWidth / 2
	radiusY := tileHeight / 2

	// Hexagon angles (flat-top)
	for i := 0; i < 6; i++ {
		angle := float64(i) * 60.0 * 3.14159 / 180.0 // Convert to radians
		vertices[i][0] = centerX + radiusX*math.Cos(angle)
		vertices[i][1] = centerY + radiusY*math.Sin(angle)
	}

	return vertices
}

// drawLine draws a line between two points using Bresenham's algorithm
func (gl *GridLayer) drawLine(img draw.Image, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	x, y := x1, y1

	var xInc, yInc int
	if x1 < x2 {
		xInc = 1
	} else {
		xInc = -1
	}
	if y1 < y2 {
		yInc = 1
	} else {
		yInc = -1
	}

	var err int
	if dx > dy {
		err = dx / 2
		for x != x2 {
			if x >= 0 && y >= 0 && x < gl.width && y < gl.height {
				img.Set(x, y, c)
			}
			err -= dy
			if err < 0 {
				y += yInc
				err += dx
			}
			x += xInc
		}
	} else {
		err = dy / 2
		for y != y2 {
			if x >= 0 && y >= 0 && x < gl.width && y < gl.height {
				img.Set(x, y, c)
			}
			err -= dx
			if err < 0 {
				x += xInc
				err += dy
			}
			y += yInc
		}
	}
}

// drawSimpleText draws simple text at the given position
func (gl *GridLayer) drawSimpleText(text string, centerX, centerY float64) {
	// For now, draw simple dots to represent coordinates
	// This can be enhanced with proper text rendering later
	bufferImg := gl.buffer.GetImageData()
	textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // White

	x, y := int(centerX), int(centerY)

	// Draw a small cross or dot to indicate coordinates
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			px, py := x+i, y+j
			if px >= 0 && py >= 0 && px < gl.width && py < gl.height {
				if (i == 0 && abs(j) <= 2) || (j == 0 && abs(i) <= 2) {
					bufferImg.Set(px, py, textColor)
				}
			}
		}
	}
}
