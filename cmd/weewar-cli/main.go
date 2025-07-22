package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

// Version information
const (
	Version = "1.0.0"
	Build   = "dev"
)

var exampleMapPath = weewar.DevDataPath("storage/maps/map123")

func main() {
	// Command line flags
	var (
		interactive = flag.Bool("interactive", false, "Start in interactive mode")
		loadFile    = flag.String("load", "", "Load game from file")
		saveFile    = flag.String("save", "", "Save game to file after commands")
		worldPath   = flag.String("world", "", fmt.Sprintf("Load world from storage directory (e.g., %s)", exampleMapPath))
		renderFile  = flag.String("render", "", "Render game to PNG file")
		width       = flag.Int("width", 800, "Render width in pixels")
		height      = flag.Int("height", 600, "Render height in pixels")
		batch       = flag.String("batch", "", "Execute commands from batch file")
		record      = flag.String("record", "", "Record session to file")
		autorender  = flag.Bool("autorender", false, "Auto-render game state after each command")
		maxrenders  = flag.Int("maxrenders", 1, "Maximum number of auto-rendered files to keep (0 disables rotation)")
		renderdir   = flag.String("renderdir", "/tmp/turnengine/autorenders", "Directory for auto-rendered files")
		version     = flag.Bool("version", false, "Show version information")
		help        = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// Show version
	if *version {
		fmt.Printf("WeeWar CLI v%s (build %s)\n", Version, Build)
		return
	}

	// Show help
	if *help {
		showHelp()
		return
	}

	// Validate flags - require either load or world
	if *loadFile == "" && *worldPath == "" && !*interactive {
		log.Fatalf("Must specify either --load <file> or --world <path> or --interactive")
	}

	// Create CLI instance
	var cli *SimpleCLI
	var game *weewar.Game

	// Initialize game
	if *loadFile != "" {
		// Load existing game
		fmt.Printf("Loading game from %s...\n", *loadFile)
		if err := loadGameFromFile(*loadFile, &game); err != nil {
			log.Fatalf("Failed to load game: %v", err)
		}
		fmt.Println("Game loaded successfully")
	} else if *worldPath != "" {
		// Load game from world storage
		fmt.Printf("Loading world from %s...\n", *worldPath)
		var err error
		game, err = createGameFromWorld(*worldPath)
		if err != nil {
			log.Fatalf("Failed to load world: %v", err)
		}
		fmt.Println("World loaded successfully")
	} else if *interactive {
		// Interactive mode without game - user must load a world
		fmt.Println("Interactive mode: use 'load <file>' or specify --world <path> to load a game")
	}

	// Create CLI
	cli = NewSimpleCLI(game)

	// Configure auto-rendering if enabled
	if *autorender {
		cli.EnableAutoRender(*renderdir, *maxrenders, *width, *height)
	}

	// Start recording if requested
	if *record != "" {
		cli.StartRecording()
		fmt.Printf("Recording session (move list will be serializable)\n")
	}

	// Execute batch commands if provided
	if *batch != "" {
		fmt.Printf("Executing batch commands from %s...\n", *batch)
		if err := executeBatchCommands(cli, *batch); err != nil {
			log.Fatalf("Batch execution failed: %v", err)
		}
		fmt.Println("Batch commands completed successfully")
	}

	// Execute remaining command line arguments as commands
	if len(flag.Args()) > 0 {
		for _, cmd := range flag.Args() {
			fmt.Printf("Executing: %s\n", cmd)
			result := cli.ExecuteCommand(cmd)
			fmt.Printf("Result: %s\n", result)
		}
	}

	// Save game if requested
	if *saveFile != "" {
		fmt.Printf("Saving game to %s...\n", *saveFile)
		if err := saveGameToFile(cli.GetGame(), *saveFile); err != nil {
			log.Fatalf("Failed to save game: %v", err)
		}
		fmt.Println("Game saved successfully")
	}

	// Render game if requested
	if *renderFile != "" {
		fmt.Printf("Rendering game to %s (%dx%d)...\n", *renderFile, *width, *height)
		if err := renderGameToFile(cli.GetGame(), *renderFile, *width, *height); err != nil {
			log.Fatalf("Failed to render game: %v", err)
		}
		fmt.Println("Game rendered successfully")
	}

	// Start interactive mode if requested
	if *interactive {
		startInteractiveMode(cli)
	}
}

// createGameFromWorld creates a game from a world stored in the storage directory
func createGameFromWorld(worldPath string) (*weewar.Game, error) {
	// Check if directory exists
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("world directory does not exist: %s", worldPath)
	}

	// Load world data from data.json
	dataPath := fmt.Sprintf("%s/data.json", worldPath)
	worldData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read world data: %w", err)
	}

	// Load metadata (optional, for display purposes)
	metadataPath := fmt.Sprintf("%s/metadata.json", worldPath)
	var metadata map[string]any
	if metadataBytes, err := os.ReadFile(metadataPath); err == nil {
		if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
			// Metadata loaded successfully
		}
	}

	// Parse the world data and create a game
	// The world data contains tiles in Q,R coordinate format
	world, err := loadWorldFromStorageJSON(worldData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse world data: %w", err)
	}

	log.Println("Found World Data: ", world.Map.TileAt(weewar.AxialCoord{Q: -1, R: -3}))

	// Create rules engine from data file
	rulesEngine, err := weewar.LoadRulesEngineFromFile(weewar.DevDataPath("data/rules-data.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to load rules engine: %w", err)
	}

	// Create game with the loaded world and rules engine (using seed 0 for now)
	game, err := weewar.NewGame(world, rulesEngine, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	if metadata != nil {
		fmt.Printf("Loaded world: %v\n", metadata["name"])
	}

	return game, nil
}

// loadWorldFromStorageJSON parses world data from storage JSON format
func loadWorldFromStorageJSON(jsonData []byte) (*weewar.World, error) {
	var storageData struct {
		Tiles map[string]struct {
			Q        int `json:"q"`
			R        int `json:"r"`
			TileType int `json:"tile_type"`
			Player   int `json:"player"`
		} `json:"tiles"`
		MapUnits []struct {
			Q        int `json:"q"`
			R        int `json:"r"`
			Player   int `json:"player"`
			UnitType int `json:"unit_type"`
		} `json:"map_units"`
	}

	// Parse JSON
	if err := json.Unmarshal(jsonData, &storageData); err != nil {
		return nil, fmt.Errorf("failed to parse storage JSON: %w", err)
	}

	// Create map from tiles
	gameMap := weewar.NewMapRect(10, 10) // Initial size, will expand as needed
	maxPlayers := 0

	// Add tiles
	for _, tileData := range storageData.Tiles {
		coord := weewar.AxialCoord{Q: tileData.Q, R: tileData.R}
		tile := weewar.NewTile(coord, tileData.TileType)
		tile.Player = tileData.Player // Set the player ownership
		gameMap.AddTile(tile)

		if tileData.Player > maxPlayers {
			maxPlayers = tileData.Player
		}
	}

	// Determine number of players (add 1 since player IDs are 0-based)
	if maxPlayers == 0 {
		maxPlayers = 2 // Default to 2 players
	} else {
		maxPlayers++
	}

	// Create world
	world, err := weewar.NewWorld(maxPlayers, gameMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create world: %w", err)
	}

	// Add units
	for _, unitData := range storageData.MapUnits {
		coord := weewar.AxialCoord{Q: unitData.Q, R: unitData.R}
		unit := weewar.NewUnit(unitData.UnitType, unitData.Player)
		unit.SetPosition(coord)
		world.AddUnit(unit)
	}

	return world, nil
}

// loadGameFromFile loads a game from a file
func loadGameFromFile(filename string, game **weewar.Game) error {
	saveData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	loadedGame, err := weewar.LoadGame(saveData)
	if err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	*game = loadedGame
	return nil
}

// showHelp displays help information
func showHelp() {
	fmt.Printf("WeeWar CLI v%s - Command Line Interface for WeeWar Games\n\n", Version)

	fmt.Println("USAGE:")
	fmt.Println("  weewar-cli [options] [commands...]")
	fmt.Println()

	fmt.Println("OPTIONS:")
	fmt.Println("  -interactive         Start in interactive mode")
	fmt.Printf("  -world PATH          Load world from storage directory (e.g., %s)", exampleMapPath)
	fmt.Println("  -load FILE           Load saved game from file")
	fmt.Println("  -save FILE           Save game to file after commands")
	fmt.Println("  -render FILE         Render game to PNG file")
	fmt.Println("  -width N             Render width in pixels (default: 800)")
	fmt.Println("  -height N            Render height in pixels (default: 600)")
	fmt.Println("  -batch FILE          Execute commands from batch file")
	fmt.Println("  -record FILE         Record session to file")
	fmt.Println("  -autorender          Auto-render game state after each command")
	fmt.Println("  -maxrenders N        Maximum number of auto-rendered files to keep (default: 10, 0 disables)")
	fmt.Println("  -renderdir DIR       Directory for auto-rendered files (default: /tmp/turnengine/autorenders)")
	fmt.Println("  -version             Show version information")
	fmt.Println("  -help                Show this help")
	fmt.Println()

	fmt.Println("GAME COMMANDS:")
	fmt.Println("  move A1 B2           Move unit from A1 to B2")
	fmt.Println("  attack A1 B2         Attack unit at B2 with unit at A1")
	fmt.Println("  status               Show current game status")
	fmt.Println("  map                  Display the game map")
	fmt.Println("  units                Show all units")
	fmt.Println("  player [N]           Show player information")
	fmt.Println("  end                  End current player's turn")
	fmt.Println("  help [command]       Show help for specific command")
	fmt.Println("  quit                 Exit the game")
	fmt.Println()

	fmt.Println("EXAMPLES:")
	fmt.Println("  # Load world and start interactive session")
	fmt.Printf("  weewar-cli -world %s -interactive", exampleMapPath)
	fmt.Println()
	fmt.Println("  # Load saved game and show status")
	fmt.Println("  weewar-cli -load mygame.json status")
	fmt.Println()
	fmt.Println("  # Load world, make moves, and save")
	fmt.Printf("  weewar-cli -world %s 'move A1 B2' 'end' -save mygame.json", exampleMapPath)
	fmt.Println()
	fmt.Println("  # Render game to PNG")
	fmt.Println("  weewar-cli -load mygame.json -render game.png")
	fmt.Println()
	fmt.Println("  # Execute batch commands")
	fmt.Println("  weewar-cli -new -batch commands.txt")
	fmt.Println()

	fmt.Println("BATCH FILE FORMAT:")
	fmt.Println("  # Comments start with #")
	fmt.Println("  move A1 B2")
	fmt.Println("  attack B2 C3")
	fmt.Println("  end")
	fmt.Println("  # Another comment")
	fmt.Println("  status")
	fmt.Println()

	fmt.Println("POSITION FORMAT:")
	fmt.Println("  Positions use chess notation: A1, B2, C3, etc.")
	fmt.Println("  Columns are A-Z, rows are 1-99")
	fmt.Println()

	fmt.Println("For more information, visit: https://github.com/panyam/turnengine")
}

// startInteractiveMode starts the REPL for interactive gameplay
func startInteractiveMode(cli *SimpleCLI) {
	fmt.Println("WeeWar CLI - Interactive Mode")
	fmt.Println("Type 'help' for available commands, 'quit' to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Show prompt
		fmt.Print("> ")

		// Read input
		if !scanner.Scan() {
			break // EOF or error
		}

		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}

		// Execute command
		result := cli.ExecuteCommand(command)

		// Check for quit
		if result == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Show result
		fmt.Println(result)
		fmt.Println()
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

// executeBatchCommands executes commands from a file (for backward compatibility)
func executeBatchCommands(cli *SimpleCLI, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open batch file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fmt.Printf("Executing line %d: %s\n", lineNum, line)
		result := cli.ExecuteCommand(line)

		if result == "quit" {
			break
		}

		fmt.Printf("Result: %s\n", result)
	}

	return scanner.Err()
}

// saveGameToFile saves the current game state to a file
func saveGameToFile(game *weewar.Game, filename string) error {
	if game == nil {
		return fmt.Errorf("no game to save")
	}

	saveData, err := game.SaveGame()
	if err != nil {
		return fmt.Errorf("failed to serialize game: %w", err)
	}

	err = os.WriteFile(filename, saveData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// renderGameToFile renders the game to a PNG file
func renderGameToFile(game *weewar.Game, filename string, width, height int) error {
	if game == nil {
		return fmt.Errorf("no game to render")
	}

	// Create buffer for PNG rendering
	buffer := weewar.NewBuffer(width, height)

	// Create buffer renderer
	renderer := weewar.NewBufferRenderer()

	// Create ViewState (empty for basic rendering - will be enhanced with overlays)
	viewState := &weewar.ViewState{
		ShowGrid:        true,
		ShowCoordinates: false,
		ZoomLevel:       1.0,
		CameraX:         0,
		CameraY:         0,
	}

	// Configure rendering options
	options := weewar.WorldRenderOptions{
		CanvasWidth:  width,
		CanvasHeight: height,
		TileWidth:    64,
		TileHeight:   64,
		YIncrement:   48,
		ShowGrid:     true,
		ShowUI:       true,
	}

	// Render the world to the buffer
	renderer.RenderWorld(game.World, viewState, buffer, options)

	// Save to PNG file
	return buffer.Save(filename)
}
