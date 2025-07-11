package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/panyam/turnengine/games/weewar"
)

// Global editor instance for WASM
var globalEditor *weewar.MapEditor

// EditorResponse represents a JavaScript-friendly response
type EditorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Keep the program running
	c := make(chan struct{})
	
	// Register JavaScript functions
	js.Global().Set("editorCreate", js.FuncOf(createEditor))
	js.Global().Set("editorNewMap", js.FuncOf(newMap))
	js.Global().Set("editorLoadMap", js.FuncOf(loadMap))
	js.Global().Set("editorSaveMap", js.FuncOf(saveMap))
	js.Global().Set("editorPaintTerrain", js.FuncOf(paintTerrain))
	js.Global().Set("editorRemoveTerrain", js.FuncOf(removeTerrain))
	js.Global().Set("editorFloodFill", js.FuncOf(floodFill))
	js.Global().Set("editorSetBrushTerrain", js.FuncOf(setBrushTerrain))
	js.Global().Set("editorSetBrushSize", js.FuncOf(setBrushSize))
	js.Global().Set("editorUndo", js.FuncOf(undo))
	js.Global().Set("editorRedo", js.FuncOf(redo))
	js.Global().Set("editorCanUndo", js.FuncOf(canUndo))
	js.Global().Set("editorCanRedo", js.FuncOf(canRedo))
	js.Global().Set("editorGetMapInfo", js.FuncOf(getMapInfo))
	js.Global().Set("editorValidateMap", js.FuncOf(validateMap))
	js.Global().Set("editorRenderMap", js.FuncOf(renderMap))
	js.Global().Set("editorExportToGame", js.FuncOf(exportToGame))
	js.Global().Set("editorGetTerrainTypes", js.FuncOf(getTerrainTypes))
	
	fmt.Println("WeeWar Map Editor WASM loaded")
	<-c
}

// createEditor creates a new map editor instance
func createEditor(this js.Value, args []js.Value) interface{} {
	globalEditor = weewar.NewMapEditor()
	
	return createEditorResponse(true, "Map editor created", "", map[string]interface{}{
		"version": "1.0.0",
		"ready":   true,
	})
}

// newMap creates a new map with specified dimensions
func newMap(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 2 {
		return createEditorResponse(false, "", "Missing width/height arguments", nil)
	}
	
	rows := args[0].Int()
	cols := args[1].Int()
	
	err := globalEditor.NewMap(rows, cols)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to create map: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("New map created (%dx%d)", rows, cols), "", map[string]interface{}{
		"width":  cols,
		"height": rows,
	})
}

// loadMap loads a map from JSON data
func loadMap(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 1 {
		return createEditorResponse(false, "", "Missing map data argument", nil)
	}
	
	// For now, return not implemented since LoadMap is a placeholder
	return createEditorResponse(false, "", "Map loading not yet implemented", nil)
}

// saveMap saves the current map to JSON
func saveMap(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	filename := "map.json"
	if len(args) >= 1 {
		filename = args[0].String()
	}
	
	// For now, return not implemented since SaveMap is a placeholder
	err := globalEditor.SaveMap(filename)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to save map: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Map saved as %s", filename), "", nil)
}

// paintTerrain paints terrain at specified coordinates
func paintTerrain(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 2 {
		return createEditorResponse(false, "", "Missing row/col arguments", nil)
	}
	
	row := args[0].Int()
	col := args[1].Int()
	
	err := globalEditor.PaintTerrain(row, col)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to paint terrain: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Terrain painted at (%d, %d)", row, col), "", nil)
}

// removeTerrain removes terrain at specified coordinates
func removeTerrain(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 2 {
		return createEditorResponse(false, "", "Missing row/col arguments", nil)
	}
	
	row := args[0].Int()
	col := args[1].Int()
	
	err := globalEditor.RemoveTerrain(row, col)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to remove terrain: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Terrain removed at (%d, %d)", row, col), "", nil)
}

// floodFill performs flood fill at specified coordinates
func floodFill(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 2 {
		return createEditorResponse(false, "", "Missing row/col arguments", nil)
	}
	
	row := args[0].Int()
	col := args[1].Int()
	
	err := globalEditor.FloodFill(row, col)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to flood fill: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Flood fill applied at (%d, %d)", row, col), "", nil)
}

// setBrushTerrain sets the current brush terrain type
func setBrushTerrain(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 1 {
		return createEditorResponse(false, "", "Missing terrain type argument", nil)
	}
	
	terrainType := args[0].Int()
	
	err := globalEditor.SetBrushTerrain(terrainType)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to set brush terrain: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Brush terrain set to type %d", terrainType), "", nil)
}

// setBrushSize sets the current brush size
func setBrushSize(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	if len(args) < 1 {
		return createEditorResponse(false, "", "Missing brush size argument", nil)
	}
	
	size := args[0].Int()
	
	err := globalEditor.SetBrushSize(size)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to set brush size: %v", err), nil)
	}
	
	hexCount := 1
	if size > 0 {
		hexCount = 1 + 6*size*(size+1)/2 // Formula for hex area
	}
	
	return createEditorResponse(true, fmt.Sprintf("Brush size set to %d (affects %d hexes)", size, hexCount), "", map[string]interface{}{
		"size":     size,
		"hexCount": hexCount,
	})
}

// undo undoes the last operation
func undo(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	err := globalEditor.Undo()
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Cannot undo: %v", err), nil)
	}
	
	return createEditorResponse(true, "Undo successful", "", map[string]interface{}{
		"canUndo": globalEditor.CanUndo(),
		"canRedo": globalEditor.CanRedo(),
	})
}

// redo redoes the last undone operation
func redo(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	err := globalEditor.Redo()
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Cannot redo: %v", err), nil)
	}
	
	return createEditorResponse(true, "Redo successful", "", map[string]interface{}{
		"canUndo": globalEditor.CanUndo(),
		"canRedo": globalEditor.CanRedo(),
	})
}

// canUndo checks if undo is available
func canUndo(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	return createEditorResponse(true, "", "", map[string]interface{}{
		"canUndo": globalEditor.CanUndo(),
	})
}

// canRedo checks if redo is available
func canRedo(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	return createEditorResponse(true, "", "", map[string]interface{}{
		"canRedo": globalEditor.CanRedo(),
	})
}

// getMapInfo returns information about the current map
func getMapInfo(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	info := globalEditor.GetMapInfo()
	if info == nil {
		return createEditorResponse(false, "", "No map loaded", nil)
	}
	
	return createEditorResponse(true, "Map info retrieved", "", map[string]interface{}{
		"filename":      info.Filename,
		"width":         info.Width,
		"height":        info.Height,
		"totalTiles":    info.TotalTiles,
		"terrainCounts": info.TerrainCounts,
		"modified":      info.Modified,
	})
}

// validateMap validates the current map
func validateMap(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	issues := globalEditor.ValidateMap()
	
	isValid := len(issues) == 0
	message := "Map is valid"
	if !isValid {
		message = fmt.Sprintf("Map has %d issue(s)", len(issues))
	}
	
	return createEditorResponse(true, message, "", map[string]interface{}{
		"valid":  isValid,
		"issues": issues,
	})
}

// renderMap renders the current map to a data URL
func renderMap(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	// Default dimensions
	width, height := 800, 600
	
	// Parse optional dimensions
	if len(args) >= 2 {
		width = args[0].Int()
		height = args[1].Int()
	}
	
	// Create a temporary game for rendering
	game, err := globalEditor.ExportToGame(2)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to export for rendering: %v", err), nil)
	}
	
	// Render to buffer
	buffer := weewar.NewBuffer(width, height)
	
	// Calculate tile sizes
	mapInfo := globalEditor.GetMapInfo()
	tileWidth := float64(width) / float64(mapInfo.Width)
	tileHeight := float64(height) / float64(mapInfo.Height)
	yIncrement := tileHeight * 0.75
	
	err = game.RenderToBuffer(buffer, tileWidth, tileHeight, yIncrement)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to render map: %v", err), nil)
	}
	
	// Convert buffer to base64 data URL
	dataURL, err := buffer.ToDataURL()
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to create data URL: %v", err), nil)
	}
	
	return createEditorResponse(true, "Map rendered successfully", "", map[string]interface{}{
		"dataURL": dataURL,
		"width":   width,
		"height":  height,
	})
}

// exportToGame exports the current map as a playable game
func exportToGame(this js.Value, args []js.Value) interface{} {
	if globalEditor == nil {
		return createEditorResponse(false, "", "Editor not initialized", nil)
	}
	
	playerCount := 2
	if len(args) >= 1 {
		playerCount = args[0].Int()
	}
	
	game, err := globalEditor.ExportToGame(playerCount)
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to export to game: %v", err), nil)
	}
	
	// Save the game data
	saveData, err := game.SaveGame()
	if err != nil {
		return createEditorResponse(false, "", fmt.Sprintf("Failed to serialize game: %v", err), nil)
	}
	
	return createEditorResponse(true, fmt.Sprintf("Map exported as %d-player game", playerCount), "", map[string]interface{}{
		"gameData":    string(saveData),
		"playerCount": playerCount,
		"size":        len(saveData),
	})
}

// getTerrainTypes returns available terrain types
func getTerrainTypes(this js.Value, args []js.Value) interface{} {
	// Get terrain data from the weewar package
	terrainTypes := []map[string]interface{}{
		{"id": 0, "name": "Unknown", "moveCost": 1, "defenseBonus": 0},
		{"id": 1, "name": "Grass", "moveCost": 1, "defenseBonus": 0},
		{"id": 2, "name": "Desert", "moveCost": 1, "defenseBonus": 0},
		{"id": 3, "name": "Water", "moveCost": 2, "defenseBonus": 0},
		{"id": 4, "name": "Mountain", "moveCost": 2, "defenseBonus": 10},
		{"id": 5, "name": "Rock", "moveCost": 3, "defenseBonus": 20},
	}
	
	return createEditorResponse(true, "Terrain types retrieved", "", map[string]interface{}{
		"terrainTypes": terrainTypes,
	})
}

// createEditorResponse creates a JavaScript-compatible response object
func createEditorResponse(success bool, message, error string, data interface{}) js.Value {
	response := EditorResponse{
		Success: success,
		Message: message,
		Error:   error,
		Data:    data,
	}
	
	// Convert to JS object
	responseBytes, _ := json.Marshal(response)
	return js.Global().Get("JSON").Call("parse", string(responseBytes))
}