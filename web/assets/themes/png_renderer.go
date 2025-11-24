package themes

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
	"github.com/turnforge/weewar/lib"
)

// PNGWorldRenderer renders worlds using PNG assets (for default theme)
type PNGWorldRenderer struct {
	theme     *DefaultTheme
	assetRoot string // Root directory for PNG assets on filesystem

	// Cache for loaded images
	tileCache  map[string]image.Image
	unitCache  map[string]image.Image
	cacheMutex sync.RWMutex
}

// NewPNGWorldRenderer creates a new PNG-based world renderer
func NewPNGWorldRenderer(theme Theme) (*PNGWorldRenderer, error) {
	defaultTheme, ok := theme.(*DefaultTheme)
	if !ok {
		return nil, fmt.Errorf("PNGWorldRenderer requires a DefaultTheme")
	}

	// Determine the asset root path - need to convert web path to filesystem path
	// The theme's basePath is "/static/assets/v1" which corresponds to "web/static/assets/v1"
	assetRoot := "web/static/assets/v1"

	return &PNGWorldRenderer{
		theme:     defaultTheme,
		assetRoot: assetRoot,
		tileCache: make(map[string]image.Image),
		unitCache: make(map[string]image.Image),
	}, nil
}

// SetAssetRoot allows overriding the asset root directory
func (r *PNGWorldRenderer) SetAssetRoot(root string) {
	r.assetRoot = root
}

// Render produces a composite PNG image of the world
func (r *PNGWorldRenderer) Render(tiles []*v1.Tile, units []*v1.Unit, options *lib.RenderOptions) ([]byte, string, error) {
	if options == nil {
		options = lib.DefaultRenderOptions()
	}

	if len(tiles) == 0 {
		return nil, "", fmt.Errorf("no tiles to render")
	}

	// Compute bounds
	minX, minY, width, height := computeBounds(tiles, units, options)

	// Create the output image
	outputImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Render tiles first (background layer)
	for _, tile := range tiles {
		if err := r.renderTile(outputImg, tile, minX, minY, options); err != nil {
			// Log but continue - don't fail entire render for one missing tile
			fmt.Printf("Warning: failed to render tile at (%d,%d): %v\n", tile.Q, tile.R, err)
		}
	}

	// Render units on top
	for _, unit := range units {
		if err := r.renderUnit(outputImg, unit, minX, minY, options); err != nil {
			fmt.Printf("Warning: failed to render unit at (%d,%d): %v\n", unit.Q, unit.R, err)
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, outputImg); err != nil {
		return nil, "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), "image/png", nil
}

// renderTile draws a single tile onto the output image
func (r *PNGWorldRenderer) renderTile(output *image.RGBA, tile *v1.Tile, offsetX, offsetY int, options *lib.RenderOptions) error {
	// Get tile image
	tileImg, err := r.getTileImage(tile.TileType, tile.Player)
	if err != nil {
		return err
	}

	// Calculate position (adjusted for offset)
	x, y := lib.HexToPixelInt32(tile.Q, tile.R, options)
	x -= offsetX
	y -= offsetY

	// Draw tile at position (x,y is top-left)
	r.drawImageAt(output, tileImg, x, y, options.TileWidth, options.TileHeight)
	return nil
}

// renderUnit draws a single unit onto the output image
func (r *PNGWorldRenderer) renderUnit(output *image.RGBA, unit *v1.Unit, offsetX, offsetY int, options *lib.RenderOptions) error {
	// Get unit image
	unitImg, err := r.getUnitImage(unit.UnitType, unit.Player)
	if err != nil {
		return err
	}

	// Calculate position (adjusted for offset)
	x, y := lib.HexToPixelInt32(unit.Q, unit.R, options)
	x -= offsetX
	y -= offsetY

	// Draw unit slightly smaller than tile (90% size) and centered within the tile
	unitWidth := int(float64(options.TileWidth) * 0.9)
	unitHeight := int(float64(options.TileHeight) * 0.9)
	unitX := x + (options.TileWidth-unitWidth)/2
	unitY := y + (options.TileHeight-unitHeight)/2
	r.drawImageAt(output, unitImg, unitX, unitY, unitWidth, unitHeight)
	return nil
}

// drawImageAt draws an image at the given top-left position with scaling and alpha blending
func (r *PNGWorldRenderer) drawImageAt(output *image.RGBA, src image.Image, x, y, width, height int) {
	srcBounds := src.Bounds()

	// Calculate destination rectangle (x,y is top-left corner)
	dstRect := image.Rect(x, y, x+width, y+height)

	// Use draw.Draw with draw.Over for proper alpha blending
	if srcBounds.Dx() == width && srcBounds.Dy() == height {
		draw.Draw(output, dstRect, src, srcBounds.Min, draw.Over)
	} else {
		// Scale the image first, then draw with alpha blending
		// Create a temporary scaled image
		scaled := image.NewRGBA(image.Rect(0, 0, width, height))

		scaleX := float64(srcBounds.Dx()) / float64(width)
		scaleY := float64(srcBounds.Dy()) / float64(height)

		for dy := 0; dy < height; dy++ {
			for dx := 0; dx < width; dx++ {
				srcX := int(float64(dx) * scaleX)
				srcY := int(float64(dy) * scaleY)
				if srcX < srcBounds.Dx() && srcY < srcBounds.Dy() {
					c := src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY)
					scaled.Set(dx, dy, c)
				}
			}
		}

		// Now draw the scaled image with proper alpha blending
		draw.Draw(output, dstRect, scaled, image.Point{}, draw.Over)
	}
}

// getTileImage loads and caches a tile image
func (r *PNGWorldRenderer) getTileImage(tileType, playerId int32) (image.Image, error) {
	// Nature terrains always use player 0
	effectivePlayer := playerId
	if r.theme.IsNatureTile(tileType) {
		effectivePlayer = 0
	}

	cacheKey := fmt.Sprintf("tile_%d_%d", tileType, effectivePlayer)

	// Check cache first
	r.cacheMutex.RLock()
	if img, ok := r.tileCache[cacheKey]; ok {
		r.cacheMutex.RUnlock()
		return img, nil
	}
	r.cacheMutex.RUnlock()

	// Load from filesystem
	path := filepath.Join(r.assetRoot, "Tiles", fmt.Sprintf("%d", tileType), fmt.Sprintf("%d.png", effectivePlayer))
	img, err := r.loadPNG(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load tile %d for player %d: %w", tileType, effectivePlayer, err)
	}

	// Cache it
	r.cacheMutex.Lock()
	r.tileCache[cacheKey] = img
	r.cacheMutex.Unlock()

	return img, nil
}

// getUnitImage loads and caches a unit image
func (r *PNGWorldRenderer) getUnitImage(unitType, playerId int32) (image.Image, error) {
	cacheKey := fmt.Sprintf("unit_%d_%d", unitType, playerId)

	// Check cache first
	r.cacheMutex.RLock()
	if img, ok := r.unitCache[cacheKey]; ok {
		r.cacheMutex.RUnlock()
		return img, nil
	}
	r.cacheMutex.RUnlock()

	// Load from filesystem
	path := filepath.Join(r.assetRoot, "Units", fmt.Sprintf("%d", unitType), fmt.Sprintf("%d.png", playerId))
	img, err := r.loadPNG(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load unit %d for player %d: %w", unitType, playerId, err)
	}

	// Cache it
	r.cacheMutex.Lock()
	r.unitCache[cacheKey] = img
	r.cacheMutex.Unlock()

	return img, nil
}

// loadPNG loads a PNG file from the filesystem
func (r *PNGWorldRenderer) loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	return img, nil
}

// ClearCache clears the image cache
func (r *PNGWorldRenderer) ClearCache() {
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()
	r.tileCache = make(map[string]image.Image)
	r.unitCache = make(map[string]image.Image)
}
