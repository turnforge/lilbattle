package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tilesCmd represents the tiles command
var tilesCmd = &cobra.Command{
	Use:   "tiles",
	Short: "List all tiles in the game",
	Long: `Display all tiles grouped by player, showing their position, health,
and remaining movement points.

Examples:
  ww tiles
  ww tiles --json`,
	RunE: runTiles,
}

func init() {
	rootCmd.AddCommand(tilesCmd)
}

func runTiles(cmd *cobra.Command, args []string) error {
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
		tiles := []map[string]any{}
		if pc.GameState.State.WorldData != nil {
			for _, tile := range pc.GameState.State.WorldData.TilesMap {
				if tile != nil {
					tiles = append(tiles, map[string]any{
						"player":    tile.Player,
						"shortcut":  tile.Shortcut,
						"q":         tile.Q,
						"r":         tile.R,
						"tile_type": tile.TileType,
					})
				}
			}
		}

		data := map[string]any{
			"game_id": gameID,
			"tiles":   tiles,
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	text := FormatTiles(pc, pc.GameState.State)
	return formatter.PrintText(text)
}
