package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1/models"
)

// endturnCmd represents the endturn command
var endturnCmd = &cobra.Command{
	Use:   "endturn",
	Short: "End the current player's turn",
	Long: `End the current player's turn and advance to the next player.
All units for the new player will be reset with full movement points.

Examples:
  ww endturn
  ww endturn --dryrun    Preview turn transition without saving`,
	RunE: runEndTurn,
}

func init() {
	rootCmd.AddCommand(endturnCmd)
}

func runEndTurn(cmd *cobra.Command, args []string) error {
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

	ctx := context.Background()

	// Get current state before ending turn
	previousPlayer := pc.GameState.State.CurrentPlayer
	previousTurn := pc.GameState.State.TurnCounter

	if isVerbose() {
		fmt.Printf("[VERBOSE] Ending turn for player %d (turn %d)\n", previousPlayer, previousTurn)
	}

	// Call EndTurnButtonClicked on presenter
	_, err = pc.Presenter.EndTurnButtonClicked(ctx, &v1.EndTurnButtonClickedRequest{
		GameId: gameID,
	})
	if err != nil {
		return fmt.Errorf("failed to end turn: %w", err)
	}

	// Get new state after ending turn
	newPlayer := pc.GameState.State.CurrentPlayer
	newTurn := pc.GameState.State.TurnCounter

	if isVerbose() {
		fmt.Printf("[VERBOSE] Turn ended. Now player %d's turn (turn %d)\n", newPlayer, newTurn)
	}

	// Save state unless in dryrun mode
	if err := savePresenterState(pc, isDryrun()); err != nil {
		return err
	}

	// Format output
	formatter := NewOutputFormatter()

	if formatter.JSON {
		data := map[string]interface{}{
			"game_id":         gameID,
			"action":          "endturn",
			"previous_player": previousPlayer,
			"previous_turn":   previousTurn,
			"current_player":  newPlayer,
			"current_turn":    newTurn,
			"success":         true,
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	var sb strings.Builder
	sb.WriteString("End Turn: Success\n")
	sb.WriteString(fmt.Sprintf("  Turn ended for player %d\n", previousPlayer))
	sb.WriteString(fmt.Sprintf("  Now player %d's turn (turn %d)\n", newPlayer, newTurn))

	if pc.GameState.State.WinningPlayer != 0 {
		sb.WriteString(fmt.Sprintf("\nGame Over! Winner: Player %d\n", pc.GameState.State.WinningPlayer))
	}

	return formatter.PrintText(sb.String())
}
