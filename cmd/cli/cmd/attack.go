package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
	"github.com/panyam/turnengine/games/weewar/services"
)

// attackCmd represents the attack command
var attackCmd = &cobra.Command{
	Use:   "attack <attacker> <target>",
	Short: "Attack a unit",
	Long: `Attack a target unit with your unit.
Positions can be unit IDs (like A1) or coordinates (like 3,4).
The <target> position can also be a direction: L, R, TL, TR, BL, BR.

Examples:
  ww attack A1 B2         Attack unit B2 with unit A1
  ww attack A1 TR         Attack top-right neighbor with A1
  ww attack 3,4 5,6       Attack position 5,6 with unit at 3,4
  ww attack A1 B2 --dryrun Preview attack outcome without saving`,
	Args: cobra.ExactArgs(2),
	RunE: runAttack,
}

func init() {
	rootCmd.AddCommand(attackCmd)
}

func runAttack(cmd *cobra.Command, args []string) error {
	attackerPos := args[0]
	targetPos := args[1]

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

	// Get runtime game for parsing positions
	rtGame, err := pc.Presenter.GamesService.GetRuntimeGame(
		pc.Presenter.GamesService.SingletonGame,
		pc.Presenter.GamesService.SingletonGameState)
	if err != nil {
		return fmt.Errorf("failed to get runtime game: %w", err)
	}

	// Parse attacker position
	attackerTarget, err := services.ParsePositionOrUnit(rtGame, attackerPos)
	if err != nil {
		return fmt.Errorf("invalid attacker position: %w", err)
	}
	attackerCoord := attackerTarget.GetCoordinate()

	// Parse target position with context (supports directions)
	targetTarget, err := services.ParsePositionOrUnitWithContext(rtGame, targetPos, &attackerCoord)
	if err != nil {
		return fmt.Errorf("invalid target position: %w", err)
	}
	targetCoord := targetTarget.GetCoordinate()

	if isVerbose() {
		fmt.Printf("[VERBOSE] Attacking from %s to %s\n", attackerCoord.String(), targetCoord.String())
	}

	// Two-click pattern: Click attacker to select, then click target to attack
	// Click 1: Select attacker on base-map layer
	_, err = pc.Presenter.SceneClicked(ctx, &v1.SceneClickedRequest{
		GameId: gameID,
		Q:      int32(attackerCoord.Q),
		R:      int32(attackerCoord.R),
		Layer:  "base-map",
	})
	if err != nil {
		return fmt.Errorf("failed to select attacker: %w", err)
	}

	if isVerbose() {
		fmt.Printf("[VERBOSE] Attacker selected at %s\n", attackerCoord.String())
	}

	// Click 2: Click target on movement-highlight layer to execute attack
	_, err = pc.Presenter.SceneClicked(ctx, &v1.SceneClickedRequest{
		GameId: gameID,
		Q:      int32(targetCoord.Q),
		R:      int32(targetCoord.R),
		Layer:  "movement-highlight",
	})
	if err != nil {
		return fmt.Errorf("failed to execute attack: %w", err)
	}

	if isVerbose() {
		fmt.Printf("[VERBOSE] Attack executed on %s\n", targetCoord.String())
	}

	// Save state unless in dryrun mode
	if err := savePresenterState(pc, isDryrun()); err != nil {
		return err
	}

	// Format output
	formatter := NewOutputFormatter()

	if formatter.JSON {
		data := map[string]interface{}{
			"game_id": gameID,
			"action":  "attack",
			"attacker": map[string]int{
				"q": attackerCoord.Q,
				"r": attackerCoord.R,
			},
			"target": map[string]int{
				"q": targetCoord.Q,
				"r": targetCoord.R,
			},
			"success": true,
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	var sb strings.Builder
	sb.WriteString("Attack: Success\n")
	sb.WriteString(fmt.Sprintf("  Attacked from %s to %s\n", attackerCoord.String(), targetCoord.String()))
	sb.WriteString(fmt.Sprintf("\nCurrent player: %d, Turn: %d\n",
		pc.GameState.State.CurrentPlayer, pc.GameState.State.TurnCounter))

	return formatter.PrintText(sb.String())
}
