//go:build js && wasm
// +build js,wasm

package weewar

// =============================================================================
// CanvasRenderer - HTML Canvas Implementation for WASM
// =============================================================================

// CanvasRenderer implements WorldRenderer for direct HTML Canvas rendering in WASM builds.
// It provides high-performance rendering by drawing directly to the canvas without PNG encoding.
type CanvasRenderer struct {
	BaseRenderer
}

// NewCanvasRenderer creates a new Canvas-based renderer for WASM
func NewCanvasRenderer() *CanvasRenderer {
	return &CanvasRenderer{}
}

// RenderWorld renders the complete world state to a CanvasBuffer
func (cr *CanvasRenderer) RenderWorld(world *World, viewState *ViewState, drawable Drawable, options WorldRenderOptions) {
	// Clear the canvas
	drawable.Clear()
	
	// Render all layers in order
	cr.RenderTerrain(world, viewState, drawable, options)
	cr.RenderHighlights(world, viewState, drawable, options)
	cr.RenderUnits(world, viewState, drawable, options)
	cr.RenderUI(world, viewState, drawable, options)
}

// RenderTerrain renders the terrain layer to a CanvasBuffer using proper hex grid calculations
func (cr *CanvasRenderer) RenderTerrain(world *World, viewState *ViewState, drawable Drawable, options WorldRenderOptions) {
	if world == nil || world.Map == nil {
		return
	}
	
	// Render each tile as a proper hexagon using exact hex grid mathematics
	for coord, tile := range world.Map.Tiles {
		if tile == nil {
			continue
		}
		
		// Convert cube coordinates to display coordinates
		displayRow, displayCol := world.Map.HexToDisplay(coord)
		
		// Calculate tile center position using proper hex grid spacing
		x, y := cr.CalculateTilePosition(displayRow, displayCol, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		// Create proper hexagon path with exact coordinates
		hexPath := cr.CreateHexagonPath(centerX, centerY, options)
		
		// Get terrain color
		fillColor := cr.GetTerrainColor(tile.TileType)
		
		// Fill the hexagon
		drawable.FillPath(hexPath, fillColor)
		
		// Draw border if grid is enabled
		if options.ShowGrid {
			borderColor := Color{R: 0, G: 0, B: 0, A: 128}
			strokeProps := StrokeProperties{Width: 1.0, LineCap: "round", LineJoin: "round"}
			drawable.StrokePath(hexPath, borderColor, strokeProps)
		}
	}
}

// RenderUnits renders the units layer to a CanvasBuffer
func (cr *CanvasRenderer) RenderUnits(world *World, viewState *ViewState, drawable Drawable, options WorldRenderOptions) {
	if world == nil {
		return
	}
	
	for _, unit := range world.Units {
		if unit == nil {
			continue
		}
		
		// Calculate unit position
		x, y := cr.CalculateTilePosition(unit.Row, unit.Col, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		// Get player color
		unitColor := cr.GetPlayerColor(unit.PlayerID)
		
		// Draw unit as a circle
		radius := (options.TileWidth + options.TileHeight) / 8 // Smaller than hex
		circlePoints := cr.createCirclePoints(centerX, centerY, radius, 12)
		
		// Fill unit circle
		drawable.FillPath(circlePoints, unitColor)
		
		// Draw unit border
		borderColor := Color{R: 0, G: 0, B: 0, A: 255}
		strokeProps := StrokeProperties{Width: 2.0, LineCap: "round", LineJoin: "round"}
		drawable.StrokePath(circlePoints, borderColor, strokeProps)
		
		// Draw unit type text
		if unit.UnitType > 0 {
			textColor := Color{R: 255, G: 255, B: 255, A: 255}
			fontSize := options.TileWidth / 6
			if fontSize < 10 {
				fontSize = 10
			}
			
			// Draw unit type number
			unitText := getUnitTypeText(unit.UnitType)
			drawable.DrawText(centerX-5, centerY+3, unitText, fontSize, textColor)
		}
	}
}

// RenderHighlights renders selection highlights and movement indicators
func (cr *CanvasRenderer) RenderHighlights(world *World, viewState *ViewState, drawable Drawable, options WorldRenderOptions) {
	if viewState == nil {
		return
	}
	
	// Highlight movable tiles (green overlay)
	for _, pos := range viewState.MovableTiles {
		x, y := cr.CalculateTilePosition(pos.Row, pos.Col, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		hexPath := cr.CreateHexagonPath(centerX, centerY, options)
		highlightColor := Color{R: 0, G: 255, B: 0, A: 64} // Transparent green
		drawable.FillPath(hexPath, highlightColor)
	}
	
	// Highlight attackable tiles (red overlay)
	for _, pos := range viewState.AttackableTiles {
		x, y := cr.CalculateTilePosition(pos.Row, pos.Col, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		hexPath := cr.CreateHexagonPath(centerX, centerY, options)
		highlightColor := Color{R: 255, G: 0, B: 0, A: 64} // Transparent red
		drawable.FillPath(hexPath, highlightColor)
	}
	
	// Highlight selected unit (yellow border)
	if viewState.SelectedUnit != nil {
		unit := viewState.SelectedUnit
		x, y := cr.CalculateTilePosition(unit.Row, unit.Col, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		hexPath := cr.CreateHexagonPath(centerX, centerY, options)
		selectionColor := Color{R: 255, G: 255, B: 0, A: 192} // Bright yellow
		strokeProps := StrokeProperties{Width: 3.0, LineCap: "round", LineJoin: "round"}
		drawable.StrokePath(hexPath, selectionColor, strokeProps)
	}
}

// RenderUI renders text overlays and UI elements
func (cr *CanvasRenderer) RenderUI(world *World, viewState *ViewState, drawable Drawable, options WorldRenderOptions) {
	if world == nil {
		return
	}
	
	// Render coordinate labels if enabled
	if options.ShowCoordinates && world.Map != nil {
		textColor := Color{R: 255, G: 255, B: 255, A: 255}
		backgroundColor := Color{R: 0, G: 0, B: 0, A: 128}
		
		for coord, tile := range world.Map.Tiles {
			if tile == nil {
				continue
			}
			
			displayRow, displayCol := world.Map.HexToDisplay(coord)
			x, y := cr.CalculateTilePosition(displayRow, displayCol, options)
			centerX := x + options.TileWidth/2
			centerY := y + options.TileHeight/2
			
			// Draw coordinate text (simplified for canvas)
			fontSize := options.TileWidth / 8
			if fontSize < 8 {
				fontSize = 8
			}
			
			// Draw row,col coordinates
			coordText := formatDisplayCoordinate(displayRow, displayCol)
			drawable.DrawTextWithStyle(centerX-10, centerY-5, coordText, fontSize, textColor, false, backgroundColor)
		}
	}
	
	// Render brush preview in editor mode
	if viewState.HoveredTile != nil && viewState.BrushSize >= 0 {
		// Show brush preview at hovered tile
		hoveredRow := viewState.HoveredTile.Row
		hoveredCol := viewState.HoveredTile.Col
		
		x, y := cr.CalculateTilePosition(hoveredRow, hoveredCol, options)
		centerX := x + options.TileWidth/2
		centerY := y + options.TileHeight/2
		
		// Draw brush preview as a dotted outline
		hexPath := cr.CreateHexagonPath(centerX, centerY, options)
		brushColor := Color{R: 255, G: 255, B: 255, A: 128}
		strokeProps := StrokeProperties{
			Width:       2.0,
			LineCap:     "round",
			LineJoin:    "round",
			DashPattern: []float64{5.0, 5.0}, // Dotted line
		}
		drawable.StrokePath(hexPath, brushColor, strokeProps)
	}
}

// =============================================================================
// Canvas-Specific Utility Functions
// =============================================================================

// createCirclePoints creates points for a circle approximation optimized for canvas rendering
func (cr *CanvasRenderer) createCirclePoints(centerX, centerY, radius float64, segments int) []Point {
	points := make([]Point, segments)
	for i := 0; i < segments; i++ {
		angle := float64(i) * 360.0 / float64(segments)
		angleRad := angle * 3.14159 / 180.0
		
		x := centerX + radius*cosApprox(angleRad)
		y := centerY + radius*sinApprox(angleRad)
		
		points[i] = Point{X: x, Y: y}
	}
	return points
}

// formatDisplayCoordinate formats display coordinates for canvas text rendering
func formatDisplayCoordinate(row, col int) string {
	// Simplified coordinate display for canvas
	// Can be enhanced later with full coordinate formatting
	return ""
	// return fmt.Sprintf("%d,%d", row, col)
}

// getUnitTypeText returns a text representation of the unit type
func getUnitTypeText(unitType int) string {
	// Simple numeric representation for now
	switch unitType {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	default:
		return "?"
	}
}