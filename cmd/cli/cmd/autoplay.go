package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/turnforge/lilbattle/services"
)

var (
	autoplaySeed     int64
	autoplayMaxTurns int
	autoplayMaxMoves int
)

var autoplayCmd = &cobra.Command{
	Use:   "autoplay",
	Short: "Drive the current game to completion via the presenter's NextMove policy",
	Long: `Drive the current game to completion using the presenter's NextMove policy
(random by default — see issue 167). On each iteration the driver asks the
presenter for the next move; nil signals "end turn." The loop stops when the
game flips Finished or the safety cap --max-turns is hit (whichever first).

Use --seed for reproducible runs.

The actual driver loop lives in services.RunAutoplay so it can be exercised
by tests without going through the CLI. This command just wires flags +
the configured GamesService.

Examples:
  ww autoplay                       Drive the current game to completion
  ww autoplay --seed 42             Reproducible run (same seed → same picks)
  ww autoplay --max-turns 100       Safety cap (default 200)`,
	RunE: runAutoplay,
}

func init() {
	rootCmd.AddCommand(autoplayCmd)
	autoplayCmd.Flags().Int64Var(&autoplaySeed, "seed", time.Now().UnixNano(),
		"RNG seed for deterministic picks (default: current time)")
	autoplayCmd.Flags().IntVar(&autoplayMaxTurns, "max-turns", 200,
		"Safety cap: abort with error after this many turn cycles to prevent runaway loops")
	autoplayCmd.Flags().IntVar(&autoplayMaxMoves, "moves", 0,
		"Stop normally after this many ProcessMoves calls (0 = run to completion). Useful for short bounded runs and replays.")
}

func runAutoplay(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	gc, err := GetGameContext()
	if err != nil {
		return err
	}

	formatter := NewOutputFormatter()
	formatter.PrintText(fmt.Sprintf("Autoplay starting on game %s (seed=%d, max-turns=%d)",
		gc.GameID, autoplaySeed, autoplayMaxTurns))

	resp, err := services.RunAutoplay(ctx, &services.RunAutoplayRequest{
		Svc:      gc.Service,
		GameID:   gc.GameID,
		Seed:     autoplaySeed,
		MaxTurns: autoplayMaxTurns,
		MaxMoves: autoplayMaxMoves,
	})
	if err != nil {
		return err
	}

	formatter.PrintText(fmt.Sprintf(
		"Autoplay finished — winner: player %d (turn %d, actions=%d, turns observed=%d)",
		resp.FinalState.WinningPlayer, resp.FinalState.TurnCounter,
		resp.ActionsApplied, resp.TurnsObserved))
	return nil
}
