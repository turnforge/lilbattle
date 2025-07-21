package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	weewar "github.com/panyam/turnengine/games/weewar/lib"
)

// MoveRecord represents a single recorded move
type MoveRecord struct {
	Command   string `json:"command"`
	Timestamp string `json:"timestamp"`
	Turn      int    `json:"turn"`
	Player    int    `json:"player"`
}

// MoveList represents a sequence of recorded moves
type MoveList struct {
	Moves []MoveRecord `json:"moves"`
}

// SimpleCLI is a thin wrapper over Game methods with minimal logic
type SimpleCLI struct {
	game          *weewar.Game
	selectedUnit  *weewar.Unit
	movableCoords []weewar.AxialCoord
	attackCoords  []weewar.AxialCoord
	recording     bool
	moveList      *MoveList
}

// NewSimpleCLI creates a new simplified CLI
func NewSimpleCLI(game *weewar.Game) *SimpleCLI {
	return &SimpleCLI{
		game:          game,
		selectedUnit:  nil,
		movableCoords: make([]weewar.AxialCoord, 0),
		attackCoords:  make([]weewar.AxialCoord, 0),
		recording:     false,
		moveList:      &MoveList{Moves: make([]MoveRecord, 0)},
	}
}

// ExecuteCommand processes a command and returns result
func (cli *SimpleCLI) ExecuteCommand(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return "Empty command"
	}

	// Record the command if recording is enabled
	if cli.recording && cli.game != nil {
		cli.recordMove(command)
	}

	parts := strings.Fields(command)
	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "move":
		return cli.handleMove(args)
	case "attack":
		return cli.handleAttack(args)
	case "select":
		return cli.handleSelect(args)
	case "end":
		return cli.handleEndTurn()
	case "status":
		return cli.handleStatus()
	case "map":
		return cli.handleMap()
	case "units":
		return cli.handleUnits()
	case "player":
		return cli.handlePlayer(args)
	case "record":
		return cli.handleRecord(args)
	case "replay":
		return cli.handleReplay(args)
	case "help":
		return cli.handleHelp()
	case "quit", "exit":
		return "quit"
	default:
		return fmt.Sprintf("Unknown command: %s. Type 'help' for available commands.", cmd)
	}
}

// handleMove processes move command: move <from> <to>
func (cli *SimpleCLI) handleMove(args []string) string {
	if len(args) != 2 {
		return "Usage: move <from> <to>\nExample: move A1 3,4"
	}

	fromTarget, err := ParsePositionOrUnit(cli.game, args[0])
	if err != nil {
		return fmt.Sprintf("Invalid from position: %v", err)
	}

	toTarget, err := ParsePositionOrUnit(cli.game, args[1])
	if err != nil {
		return fmt.Sprintf("Invalid to position: %v", err)
	}

	// For move, 'from' should be a unit, 'to' should be a coordinate
	if !fromTarget.IsUnit {
		return fmt.Sprintf("From position must be a unit (like A1), got coordinate: %s", fromTarget.Raw)
	}

	fromCoord := fromTarget.GetCoordinate()
	toCoord := toTarget.GetCoordinate()

	// Call the game's move method
	err = cli.game.MoveUnitAt(fromCoord, toCoord)
	if err != nil {
		return fmt.Sprintf("Move failed: %v", err)
	}

	return fmt.Sprintf("Moved %s from %s to %s",
		fromTarget.Raw, fromCoord.String(), toCoord.String())
}

// handleAttack processes attack command: attack <attacker> <target>
func (cli *SimpleCLI) handleAttack(args []string) string {
	if len(args) != 2 {
		return "Usage: attack <attacker> <target>\nExample: attack A1 B2"
	}

	attackerTarget, err := ParsePositionOrUnit(cli.game, args[0])
	if err != nil {
		return fmt.Sprintf("Invalid attacker position: %v", err)
	}

	targetTarget, err := ParsePositionOrUnit(cli.game, args[1])
	if err != nil {
		return fmt.Sprintf("Invalid target position: %v", err)
	}

	attackerCoord := attackerTarget.GetCoordinate()
	targetCoord := targetTarget.GetCoordinate()

	// Call the game's attack method
	result, err := cli.game.AttackUnitAt(attackerCoord, targetCoord)
	if err != nil {
		return fmt.Sprintf("Attack failed: %v", err)
	}

	if result != nil {
		return fmt.Sprintf("Attack successful! %s attacked %s. Defender took %d damage",
			attackerTarget.Raw, targetTarget.Raw, result.DefenderDamage)
	}

	return fmt.Sprintf("%s attacked %s", attackerTarget.Raw, targetTarget.Raw)
}

// handleSelect processes select command: select <unit>
func (cli *SimpleCLI) handleSelect(args []string) string {
	if len(args) != 1 {
		return "Usage: select <unit>\nExample: select A1"
	}

	target, err := ParsePositionOrUnit(cli.game, args[0])
	if err != nil {
		return fmt.Sprintf("Invalid unit: %v", err)
	}

	if !target.IsUnit {
		return fmt.Sprintf("Select target must be a unit (like A1), got coordinate: %s", target.Raw)
	}

	cli.selectedUnit = target.GetUnit()

	// Get movement and attack options for the selected unit
	cli.movableCoords = make([]weewar.AxialCoord, 0)
	cli.attackCoords = make([]weewar.AxialCoord, 0)

	// Get movement options
	moveOptions, err := cli.game.GetUnitMovementOptions(cli.selectedUnit)
	if err == nil {
		for _, option := range moveOptions {
			cli.movableCoords = append(cli.movableCoords, option.Coord)
		}
	}

	// Get attack options
	attackOptions, err := cli.game.GetUnitAttackOptions(cli.selectedUnit)
	if err == nil {
		cli.attackCoords = attackOptions
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Selected %s at %s\n", target.Raw, target.GetCoordinate().String()))

	if len(cli.movableCoords) > 0 {
		result.WriteString(fmt.Sprintf("Can move to %d positions: ", len(cli.movableCoords)))
		for i, coord := range cli.movableCoords {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(coord.String())
			if i >= 9 { // Limit output
				result.WriteString(fmt.Sprintf(" and %d more...", len(cli.movableCoords)-10))
				break
			}
		}
		result.WriteString("\n")
	} else {
		result.WriteString("No valid moves available\n")
	}

	if len(cli.attackCoords) > 0 {
		result.WriteString(fmt.Sprintf("Can attack %d positions: ", len(cli.attackCoords)))
		for i, coord := range cli.attackCoords {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(coord.String())
			if i >= 9 { // Limit output
				result.WriteString(fmt.Sprintf(" and %d more...", len(cli.attackCoords)-10))
				break
			}
		}
		result.WriteString("\n")
	} else {
		result.WriteString("No valid attacks available\n")
	}

	return result.String()
}

// handleEndTurn processes end turn command
func (cli *SimpleCLI) handleEndTurn() string {
	err := cli.game.EndTurn()
	if err != nil {
		return fmt.Sprintf("End turn failed: %v", err)
	}

	cli.selectedUnit = nil // Clear selection on turn end
	return fmt.Sprintf("Turn ended. Current player: %d", cli.game.CurrentPlayer)
}

// handleStatus shows current game status
func (cli *SimpleCLI) handleStatus() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Turn: %d\n", cli.game.TurnCounter))
	result.WriteString(fmt.Sprintf("Current Player: %d\n", cli.game.CurrentPlayer))
	result.WriteString(fmt.Sprintf("Game Status: %s\n", cli.game.Status))

	// Show unit counts per player
	for playerID, units := range cli.game.World.UnitsByPlayer {
		result.WriteString(fmt.Sprintf("Player %c: %d units\n",
			'A'+playerID, len(units)))
	}

	if cli.selectedUnit != nil {
		unitID := cli.game.GetUnitID(cli.selectedUnit)
		result.WriteString(fmt.Sprintf("Selected: %s at %s\n",
			unitID, cli.selectedUnit.Coord.String()))
	}

	return result.String()
}

// handleMap shows a simple map representation
func (cli *SimpleCLI) handleMap() string {
	if cli.game.World.Map == nil {
		return "No map loaded"
	}

	rows, cols := cli.game.World.GetMapSizeRect()
	return fmt.Sprintf("Map size: %dx%d\nUse render command to generate PNG visualization", rows, cols)
}

// handleUnits shows all units
func (cli *SimpleCLI) handleUnits() string {
	var result strings.Builder

	for playerID, units := range cli.game.World.UnitsByPlayer {
		if len(units) == 0 {
			continue
		}

		playerLetter := string(rune('A' + playerID))
		result.WriteString(fmt.Sprintf("Player %s units:\n", playerLetter))

		for i, unit := range units {
			if unit == nil {
				continue
			}
			unitID := fmt.Sprintf("%s%d", playerLetter, i+1)
			result.WriteString(fmt.Sprintf("  %s: Type %d at %s (HP: %d)\n",
				unitID, unit.UnitType, unit.Coord.String(), unit.AvailableHealth))
		}
	}

	if result.Len() == 0 {
		return "No units found"
	}

	return result.String()
}

// handlePlayer shows player information
func (cli *SimpleCLI) handlePlayer(args []string) string {
	playerID := cli.game.CurrentPlayer
	if len(args) > 0 {
		// Parse player ID if provided
		if len(args[0]) == 1 {
			playerLetter := strings.ToUpper(args[0])[0]
			if playerLetter >= 'A' && playerLetter <= 'Z' {
				playerID = int(playerLetter - 'A')
			}
		}
	}

	if playerID >= len(cli.game.World.UnitsByPlayer) {
		return fmt.Sprintf("Player %d does not exist", playerID)
	}

	units := cli.game.World.UnitsByPlayer[playerID]
	playerLetter := string(rune('A' + playerID))

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Player %s:\n", playerLetter))
	result.WriteString(fmt.Sprintf("  Units: %d\n", len(units)))

	if playerID == cli.game.CurrentPlayer {
		result.WriteString("  Status: Active (current turn)\n")
	} else {
		result.WriteString("  Status: Waiting\n")
	}

	return result.String()
}

// handleHelp shows available commands
func (cli *SimpleCLI) handleHelp() string {
	return `Available commands:
  move <from> <to>     - Move unit (e.g. "move A1 3,4" or "move A1 r5,6")
  attack <att> <tgt>   - Attack target (e.g. "attack A1 B2" or "attack 2,3 C1")  
  select <unit>        - Select unit and show options (e.g. "select A1")
  end                  - End current player's turn
  status               - Show game status
  map                  - Show map information
  units                - Show all units
  player [ID]          - Show player information (e.g. "player A")
  record <cmd>         - Recording commands (start/stop/clear/show)
  replay               - Show current move list as JSON
  help                 - Show this help
  quit                 - Exit game

Position formats:
  - Unit ID: A1, B12, C2 (Player letter + unit number)
  - Q,R coordinate: 3,4 or -1,2
  - Row/col coordinate: r4,5 (prefix with 'r')

Examples:
  move A1 3,4          # Move unit A1 to Q/R coordinate 3,4
  attack r4,5 B2       # Attack unit B2 with unit at row/col 4,5
  select C1            # Select unit C1 and show movement/attack options
  record start         # Start recording moves
  record show          # Show recorded moves`
}

// GetSelectedUnit returns currently selected unit (for rendering)
func (cli *SimpleCLI) GetSelectedUnit() *weewar.Unit {
	return cli.selectedUnit
}

// GetMovableCoords returns highlighted movement coordinates (for rendering)
func (cli *SimpleCLI) GetMovableCoords() []weewar.AxialCoord {
	return cli.movableCoords
}

// GetAttackCoords returns highlighted attack coordinates (for rendering)
func (cli *SimpleCLI) GetAttackCoords() []weewar.AxialCoord {
	return cli.attackCoords
}

// recordMove adds a command to the move list if recording is enabled
func (cli *SimpleCLI) recordMove(command string) {
	if cli.moveList == nil {
		cli.moveList = &MoveList{Moves: make([]MoveRecord, 0)}
	}

	record := MoveRecord{
		Command:   command,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Turn:      cli.game.TurnCounter,
		Player:    cli.game.CurrentPlayer,
	}

	cli.moveList.Moves = append(cli.moveList.Moves, record)
}

// handleRecord processes record command: record start/stop/save/load
func (cli *SimpleCLI) handleRecord(args []string) string {
	if len(args) == 0 {
		if cli.recording {
			return fmt.Sprintf("Recording active (%d moves recorded)", len(cli.moveList.Moves))
		}
		return "Recording inactive. Use 'record start' to begin recording."
	}

	switch strings.ToLower(args[0]) {
	case "start":
		cli.recording = true
		cli.moveList = &MoveList{Moves: make([]MoveRecord, 0)}
		return "Recording started"

	case "stop":
		cli.recording = false
		return fmt.Sprintf("Recording stopped (%d moves recorded)", len(cli.moveList.Moves))

	case "clear":
		cli.moveList = &MoveList{Moves: make([]MoveRecord, 0)}
		return "Move list cleared"

	case "show":
		if len(cli.moveList.Moves) == 0 {
			return "No moves recorded"
		}
		var result strings.Builder
		result.WriteString(fmt.Sprintf("Recorded moves (%d):\n", len(cli.moveList.Moves)))
		for i, move := range cli.moveList.Moves {
			result.WriteString(fmt.Sprintf("  %d. Turn %d, Player %d: %s\n",
				i+1, move.Turn, move.Player, move.Command))
		}
		return result.String()

	default:
		return "Usage: record start/stop/clear/show"
	}
}

// handleReplay processes replay command: replay <move_list_json>
func (cli *SimpleCLI) handleReplay(args []string) string {
	if len(args) == 0 {
		return "Usage: replay <move_list_json>"
	}

	// For now, just return the current move list as JSON
	// In a full implementation, this would parse the JSON and replay moves
	if len(cli.moveList.Moves) == 0 {
		return "No moves to replay"
	}

	jsonData, err := json.Marshal(cli.moveList)
	if err != nil {
		return fmt.Sprintf("Failed to serialize move list: %v", err)
	}

	return fmt.Sprintf("Move list JSON:\n%s", string(jsonData))
}

// GetMoveList returns the current move list for serialization
func (cli *SimpleCLI) GetMoveList() *MoveList {
	return cli.moveList
}

// LoadMoveList sets the move list from serialized data
func (cli *SimpleCLI) LoadMoveList(moveList *MoveList) {
	cli.moveList = moveList
}

// StartRecording begins recording moves
func (cli *SimpleCLI) StartRecording() {
	cli.recording = true
	if cli.moveList == nil {
		cli.moveList = &MoveList{Moves: make([]MoveRecord, 0)}
	}
}

// StopRecording stops recording moves
func (cli *SimpleCLI) StopRecording() {
	cli.recording = false
}

// GetGame returns the current game instance
func (cli *SimpleCLI) GetGame() *weewar.Game {
	return cli.game
}
