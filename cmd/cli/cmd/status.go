package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current game status",
	Long: `Display the current game state including turn number, current player,
and game status.

Examples:
  ww status
  ww status --json`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Get game ID
	gameID, err := getGameID()
	if err != nil {
		return err
	}

	// Create presenter
	pc, err := createPresenter(gameID)
	if err != nil {
		return err
	}

	// Get state from panel
	if pc.GameState.State == nil {
		return fmt.Errorf("game state not initialized")
	}
	if pc.GameState.Game == nil {
		return fmt.Errorf("game metadata not initialized")
	}

	// Format output
	formatter := NewOutputFormatter()

	if formatter.JSON {
		// Build player info for JSON output
		players := []map[string]interface{}{}

		// Count units per player
		unitCounts := make(map[int32]int)
		if pc.GameState.State.WorldData != nil {
			for _, unit := range pc.GameState.State.WorldData.Units {
				if unit != nil {
					unitCounts[unit.Player]++
				}
			}
		}

		// Count tiles per player
		tileCounts := make(map[int32]int)
		if pc.GameState.State.WorldData != nil {
			for _, tile := range pc.GameState.State.WorldData.Tiles {
				if tile != nil && tile.Player > 0 {
					tileCounts[tile.Player]++
				}
			}
		}

		if pc.GameState.Game.Config != nil {
			for _, player := range pc.GameState.Game.Config.Players {
				players = append(players, map[string]interface{}{
					"player_id":   player.PlayerId,
					"player_type": player.PlayerType,
					"name":        player.Name,
					"coins":       player.Coins,
					"units":       unitCounts[player.PlayerId],
					"tiles":       tileCounts[player.PlayerId],
					"team_id":     player.TeamId,
					"is_active":   player.IsActive,
				})
			}
		}

		// JSON output
		data := map[string]interface{}{
			"game_id":        gameID,
			"game_name":      pc.GameState.Game.Name,
			"description":    pc.GameState.Game.Description,
			"turn":           pc.GameState.State.TurnCounter,
			"current_player": pc.GameState.State.CurrentPlayer,
			"status":         pc.GameState.State.Status.String(),
			"winning_player": pc.GameState.State.WinningPlayer,
			"players":        players,
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	text := FormatGameStatus(pc.GameState.Game, pc.GameState.State)
	return formatter.PrintText(text)
}
