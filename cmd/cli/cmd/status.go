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

	// Format output
	formatter := NewOutputFormatter()

	if formatter.JSON {
		// JSON output
		data := map[string]interface{}{
			"game_id":        gameID,
			"turn":           pc.GameState.State.TurnCounter,
			"current_player": pc.GameState.State.CurrentPlayer,
			"status":         pc.GameState.State.Status.String(),
			"winning_player": pc.GameState.State.WinningPlayer,
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	text := FormatGameStatus(pc.GameState.State)
	return formatter.PrintText(text)
}
