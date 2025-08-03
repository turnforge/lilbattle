package weewar

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// =============================================================================
// Core Types (from core.go)
// =============================================================================

// PlayerInfo contains game-specific information about a player
type PlayerInfo struct {
	Player   int32  `json:"playerID"` // 1-based player index
	Name     string `json:"name"`     // Player display name
	TeamID   int    `json:"teamID"`   // Which team this player belongs to
	IsActive bool   `json:"isActive"` // Whether player is still in the game
	Color    string `json:"color"`    // Player's display color
}

// TeamInfo contains information about a team
type TeamInfo struct {
	TeamID      int    `json:"teamID"`      // 0-based team index
	Name        string `json:"name"`        // Team display name
	Color       string `json:"color"`       // Team color
	IsActive    bool   `json:"isActive"`    // Whether team has active players
	PlayerCount int    `json:"playerCount"` // Number of players in this team
}

// Game represents the unified game state and implements GameInterface
type Game struct {
	// Pure game state
	World *World `json:"world"` // Contains the pure state (map, units, entities)

	// Game flow control
	CurrentPlayer int32      `json:"currentPlayer"` // 0-based player index
	TurnCounter   int32      `json:"turnCounter"`   // 1-based turn number
	Status        GameStatus `json:"status"`        // Game status

	// Player and team information
	Players []PlayerInfo `json:"players"` // Information about each player
	Teams   []TeamInfo   `json:"teams"`   // Information about each team

	// Game systems and configuration
	Seed int64 `json:"seed"` // Random seed for deterministic gameplay

	// Random number generator
	rng *rand.Rand `json:"-"` // RNG for deterministic gameplay

	// Asset management
	assetProvider AssetProvider `json:"-"` // Asset provider for tiles and units (interface for platform flexibility)

	// Rules engine for data-driven game mechanics
	rulesEngine *RulesEngine `json:"-"` // Rules engine for movement costs, combat, unit data

	// Game metadata
	CreatedAt    time.Time `json:"createdAt"`    // When game was created
	LastActionAt time.Time `json:"lastActionAt"` // When last action was taken

	// Internal state
	winner    int32 `json:"winner"`    // Winner player ID (-1 if no winner)
	hasWinner bool  `json:"hasWinner"` // Whether game has ended with winner
}

// =============================================================================
// Convenience methods to access World fields
// =============================================================================

// GetPlayerInfo returns information about a specific player
func (g *Game) GetPlayerInfo(playerID int) (*PlayerInfo, error) {
	if playerID < 0 || playerID >= len(g.Players) {
		return nil, fmt.Errorf("invalid player ID: %d", playerID)
	}
	return &g.Players[playerID], nil
}

// GetTeamInfo returns information about a specific team
func (g *Game) GetTeamInfo(teamID int) (*TeamInfo, error) {
	if teamID < 0 || teamID >= len(g.Teams) {
		return nil, fmt.Errorf("invalid team ID: %d", teamID)
	}
	return &g.Teams[teamID], nil
}

// GetPlayersOnTeam returns all players belonging to a specific team
func (g *Game) GetPlayersOnTeam(teamID int) []*PlayerInfo {
	var teamPlayers []*PlayerInfo
	for i := range g.Players {
		if g.Players[i].TeamID == teamID {
			teamPlayers = append(teamPlayers, &g.Players[i])
		}
	}
	return teamPlayers
}

// ArePlayersOnSameTeam checks if two players are on the same team
func (g *Game) ArePlayersOnSameTeam(playerID1, playerID2 int) bool {
	if playerID1 < 0 || playerID1 >= len(g.Players) ||
		playerID2 < 0 || playerID2 >= len(g.Players) {
		return false
	}
	return g.Players[playerID1].TeamID == g.Players[playerID2].TeamID
}

// =============================================================================
// Helper Functions
// =============================================================================

// initializeStartingUnits initializes stats for units already in the World
func (g *Game) initializeStartingUnits() error {
	// Get unit stats from rules engine (required)
	if g.rulesEngine == nil {
		return fmt.Errorf("rules engine not set - required for unit initialization")
	}

	// Initialize stats for existing units in the world
	fmt.Printf("initializeStartingUnits: PlayerCount=%d, TurnCounter=%d\n", g.World.PlayerCount(), g.TurnCounter)
	for playerID := int32(1); playerID <= g.World.PlayerCount(); playerID++ {
		fmt.Printf("initializeStartingUnits: Player %d has %d units\n", playerID, len(g.World.unitsByPlayer[playerID]))
		for _, unit := range g.World.unitsByPlayer[playerID] {
			// Get unit data from rules engine
			fmt.Printf("initializeStartingUnits: Processing unit type %d, current DistanceLeft=%d\n", unit.UnitType, unit.DistanceLeft)
			unitData, err := g.rulesEngine.GetUnitData(unit.UnitType)
			if err != nil {
				return fmt.Errorf("failed to get unit data for type %d: %w", unit.UnitType, err)
			}

			// Initialize unit stats from rules data
			unit.AvailableHealth = unitData.Health
			unit.DistanceLeft = unitData.MovementPoints
			unit.TurnCounter = g.TurnCounter
			fmt.Printf("initializeStartingUnits: Set unit DistanceLeft to %d (from rules engine MovementPoints %d)\n", unit.DistanceLeft, unitData.MovementPoints)
		}
	}

	return nil
}

// resetPlayerUnits resets movement and actions for a player's units
func (g *Game) resetPlayerUnits(playerID int32) error {
	if playerID <= 0 || playerID > g.World.PlayerCount() {
		return fmt.Errorf("invalid player ID: %d", playerID)
	}

	if g.rulesEngine == nil {
		return fmt.Errorf("rules engine not set - required for unit reset")
	}

	fmt.Printf("resetPlayerUnits: Resetting player %d units, PlayerCount=%d, TurnCounter=%d\n", playerID, g.World.PlayerCount(), g.TurnCounter)
	for _, unit := range g.World.unitsByPlayer[playerID] {
		// Get unit data from rules engine
		unitData, err := g.rulesEngine.GetUnitData(unit.UnitType)
		if err != nil {
			return fmt.Errorf("failed to get unit data for type %d: %w", unit.UnitType, err)
		}

		// Reset movement points from rules data
		fmt.Printf("resetPlayerUnits: Unit type %d, current DistanceLeft=%d, rules MovementPoints=%d\n", 
			unit.UnitType, unit.DistanceLeft, unitData.MovementPoints)
		unit.DistanceLeft = int32(unitData.MovementPoints)
		unit.TurnCounter = int32(g.TurnCounter)
		fmt.Printf("resetPlayerUnits: Set unit DistanceLeft to %d\n", unit.DistanceLeft)
	}

	return nil
}

// checkVictoryConditions checks if any player has won
func (g *Game) checkVictoryConditions() (winner int32, hasWinner bool) {
	// Simple victory condition: last player with units wins
	playersWithUnits := 0
	lastPlayerWithUnits := int32(-1)

	for playerID := int32(1); playerID <= g.World.PlayerCount(); playerID++ {
		if len(g.World.unitsByPlayer[playerID]) > 0 {
			playersWithUnits++
			lastPlayerWithUnits = playerID
		}
	}

	if playersWithUnits == 1 {
		return lastPlayerWithUnits, true
	}

	return -1, false
}

// validateGameState validates the current game state
func (g *Game) validateGameState() error {
	if g.World == nil {
		return fmt.Errorf("game has no world")
	}

	if g.CurrentPlayer < 0 || g.CurrentPlayer > g.World.PlayerCount() {
		return fmt.Errorf("invalid current player: %d", g.CurrentPlayer)
	}

	if g.TurnCounter < 1 {
		return fmt.Errorf("invalid turn counter: %d", g.TurnCounter)
	}

	if int32(len(g.World.unitsByPlayer)) != g.World.PlayerCount() {
		return fmt.Errorf("units array length (%d) doesn't match player count (%d)", len(g.World.unitsByPlayer), g.World.PlayerCount())
	}

	return nil
}

// GetUnitID generates a unique identifier for a unit in the format PN
// where P is the player letter (A-Z) and N is the unit number for that player
func (g *Game) GetUnitID(unit *v1.Unit) string {
	if unit == nil {
		return ""
	}

	// This method was only used for cli - we can come back to this when needed
	panic("to be deprecated")
	// return unit.unitID
}

// GetAssetManager returns the current AssetManager instance (legacy compatibility)
func (g *Game) GetAssetManager() *AssetManager {
	// Try to cast the AssetProvider to *AssetManager for backward compatibility
	if am, ok := g.assetProvider.(*AssetManager); ok {
		return am
	}
	return nil
}

// SetAssetManager sets the AssetManager instance for tile and unit rendering (legacy compatibility)
func (g *Game) SetAssetManager(assetManager *AssetManager) {
	g.assetProvider = assetManager
}

// GetAssetProvider returns the current AssetProvider instance
func (g *Game) GetAssetProvider() AssetProvider {
	return g.assetProvider
}

// SetAssetProvider sets the AssetProvider instance for tile and unit rendering
func (g *Game) SetAssetProvider(provider AssetProvider) {
	g.assetProvider = provider
}

// GetRulesEngine returns the current RulesEngine instance
func (g *Game) GetRulesEngine() *RulesEngine {
	return g.rulesEngine
}

// SetRulesEngine sets the RulesEngine instance for data-driven game mechanics
func (g *Game) SetRulesEngine(rulesEngine *RulesEngine) {
	g.rulesEngine = rulesEngine
}

// NewGame creates a new game instance with the specified parameters
func NewGame(world *World, rulesEngine *RulesEngine, seed int64) (*Game, error) {
	// Validate parameters
	if rulesEngine == nil {
		return nil, fmt.Errorf("rules engine is required")
	}

	// Create the game struct
	game := &Game{
		World:         world,
		Seed:          seed,
		CurrentPlayer: 1,
		TurnCounter:   1,
		Status:        GameStatusPlaying,
		winner:        -1,
		hasWinner:     false,
		CreatedAt:     time.Now(),
		LastActionAt:  time.Now(),
		rng:           rand.New(rand.NewSource(seed)),
		assetProvider: NewAssetManager("data"),
		rulesEngine:   rulesEngine,
	}

	// Initialize units storage for compatibility (will be migrated)

	// privateMap is already assigned in the struct initialization above

	// Initialize starting units (simplified for now)
	// TODO: Replace with actual unit placement from map data
	if err := game.initializeStartingUnits(); err != nil {
		return nil, fmt.Errorf("failed to initialize starting units: %w", err)
	}

	return game, nil
}

// LoadGame restores a game from saved JSON data
func LoadGame(saveData []byte) (*Game, error) {
	var game Game
	if err := json.Unmarshal(saveData, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game data: %w", err)
	}

	// Restore transient state
	game.rng = rand.New(rand.NewSource(game.Seed))
	game.assetProvider = NewAssetManager("data")
	game.rulesEngine = nil // Will be set by caller

	// Note: Neighbor connections are no longer stored, calculated on-demand

	// Validate loaded game state
	if err := game.validateGameState(); err != nil {
		return nil, fmt.Errorf("invalid saved game state: %w", err)
	}

	return &game, nil
}

// =============================================================================
// GameController Interface Implementation
// =============================================================================

// LoadGame restores game from saved state (interface method)
func (g *Game) LoadGame(saveData []byte) (*Game, error) {
	return LoadGame(saveData)
}

// SaveGame serializes current game state
func (g *Game) SaveGame() ([]byte, error) {
	// Update last action time
	g.LastActionAt = time.Now()

	// Serialize to JSON
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize game state: %w", err)
	}

	return data, nil
}

// GetCurrentPlayer returns active player ID
func (g *Game) GetCurrentPlayer() int32 {
	return g.CurrentPlayer
}

// GetTurnNumber returns current turn count
func (g *Game) GetTurnNumber() int32 {
	return g.TurnCounter
}

// GetGameStatus returns current game state
func (g *Game) GetGameStatus() GameStatus {
	return g.Status
}

// GetWinner returns winning player if game ended
func (g *Game) GetWinner() (int32, bool) {
	return g.winner, g.hasWinner
}

// =============================================================================
// UnitInterface Interface Implementation
// =============================================================================

// GetUnitsForPlayer returns all units owned by player
func (g *Game) GetUnitsForPlayer(playerID int) []*v1.Unit {
	if playerID < 0 || playerID >= len(g.World.unitsByPlayer) {
		return nil
	}

	// Return a copy to prevent external modification
	units := make([]*v1.Unit, len(g.World.unitsByPlayer[playerID]))
	copy(units, g.World.unitsByPlayer[playerID])
	return units
}

// GetUnitTypeName returns the display name for a unit type
func (g *Game) GetUnitTypeName(unitType int) string {
	if g.assetProvider != nil {
		// Try to get unit data from JSON if asset provider is loaded
		if am, ok := g.assetProvider.(*AssetManager); ok {
			if err := am.LoadGameData(); err == nil {
				if unitData, err := am.GetUnitData(unitType); err == nil {
					return unitData.Name
				}
			}
		}
	}

	// Fallback to generic name
	return fmt.Sprintf("Unit Type %d", unitType)
}
