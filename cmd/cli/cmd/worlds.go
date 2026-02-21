package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	"github.com/turnforge/lilbattle/lib"
	"github.com/turnforge/lilbattle/services/connectclient"
	"github.com/turnforge/lilbattle/web/assets/themes"
)

// worldsCmd is the parent command for world operations
var worldsCmd = &cobra.Command{
	Use:   "worlds",
	Short: "Manage worlds on a server",
	Long: `List, inspect, and render worlds from a configured server profile.

Examples:
  ww worlds list                  # list worlds on active profile
  ww worlds list --profile prod   # list worlds on prod profile
  ww worlds get aruba             # get world details
  ww worlds get prod:aruba        # get world from specific profile
  ww worlds show aruba            # render world map inline`,
}

// worldsListCmd lists all worlds
var worldsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worlds on the server",
	Long: `List all worlds with summary info (ID, name, tile/unit counts, difficulty).

Examples:
  ww worlds list
  ww worlds list --profile prod
  ww worlds list --json`,
	RunE: runWorldsList,
}

// worldsGetCmd shows detailed world metadata
var worldsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Show detailed metadata for a world",
	Long: `Show detailed metadata for a specific world including tiles, units,
description, tags, and difficulty.

The world ID can include a profile prefix (profile:id).

Examples:
  ww worlds get aruba
  ww worlds get prod:aruba
  ww worlds get --profile prod aruba
  ww worlds get aruba --json`,
	Args: cobra.ExactArgs(1),
	RunE: runWorldsGet,
}

// worldsShowCmd renders a world map
var worldsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Render world map in the terminal",
	Long: `Render a world's map as an inline image in the terminal.
Requires iTerm2, Kitty, or another terminal with inline image support.

The world ID can include a profile prefix (profile:id).

Examples:
  ww worlds show aruba
  ww worlds show prod:aruba
  ww worlds show aruba -o map.png
  ww worlds show aruba --labels --tile-labels`,
	Args: cobra.ExactArgs(1),
	RunE: runWorldsShow,
}

var (
	worldShowLabels     bool
	worldShowTileLabels bool
	worldShowOutput     string
)

func init() {
	rootCmd.AddCommand(worldsCmd)
	worldsCmd.AddCommand(worldsListCmd)
	worldsCmd.AddCommand(worldsGetCmd)
	worldsCmd.AddCommand(worldsShowCmd)

	worldsShowCmd.Flags().BoolVar(&worldShowLabels, "labels", true, "Show unit labels")
	worldsShowCmd.Flags().BoolVar(&worldShowTileLabels, "tile-labels", true, "Show tile labels")
	defaultOutput := os.Getenv("LILBATTLE_MAP_OUTPUT")
	worldsShowCmd.Flags().StringVarP(&worldShowOutput, "output", "o", defaultOutput, "Save image to file instead of displaying")
}

// getWorldsClient creates a ConnectWorldsClient for the active profile or specified spec.
// If spec contains a ":" (profile:worldId format), it parses the profile and returns
// the client along with the extracted world ID.
// Otherwise it uses the active profile's server URL.
func getWorldsClient(spec string) (*connectclient.ConnectWorldsClient, string, error) {
	// Check if spec contains profile:worldId format
	if spec != "" && !strings.HasPrefix(spec, "http") && strings.Contains(spec, ":") {
		ws, err := parseWorldSpec(spec)
		if err != nil {
			return nil, "", err
		}
		client := connectclient.NewConnectWorldsClientWithAuth(ws.APIEndpoint(), ws.Token)
		return client, ws.WorldID, nil
	}

	// Use active profile
	serverURL := getServerURL()
	if serverURL == "" {
		return nil, "", fmt.Errorf("no server configured (set --profile, --server, or LILBATTLE_SERVER)")
	}

	token := GetTokenForProfile(getProfileName())
	apiURL := GetAPIEndpoint(serverURL)

	if isVerbose() {
		fmt.Printf("[VERBOSE] Connecting to: %s (auth: %v)\n", apiURL, token != "")
	}

	client := connectclient.NewConnectWorldsClientWithAuth(apiURL, token)
	return client, spec, nil
}

func runWorldsList(cmd *cobra.Command, args []string) error {
	client, _, err := getWorldsClient("")
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := client.ListWorlds(ctx, &v1.ListWorldsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list worlds: %w", err)
	}

	formatter := NewOutputFormatter()

	if formatter.JSON {
		items := []map[string]any{}
		for _, w := range resp.Items {
			items = append(items, worldSummaryMap(w))
		}
		return formatter.PrintJSON(map[string]any{
			"worlds": items,
			"total":  len(items),
		})
	}

	if len(resp.Items) == 0 {
		fmt.Println("No worlds found.")
		return nil
	}

	// Table-style output
	fmt.Printf("%-20s %-30s %-12s\n", "ID", "NAME", "DIFFICULTY")
	fmt.Println(strings.Repeat("-", 64))
	for _, w := range resp.Items {
		difficulty := w.Difficulty
		if difficulty == "" {
			difficulty = "-"
		}
		fmt.Printf("%-20s %-30s %-12s\n",
			truncate(w.Id, 20),
			truncate(w.Name, 30),
			difficulty,
		)
	}
	fmt.Printf("\n%d world(s)\n", len(resp.Items))

	return nil
}

func runWorldsGet(cmd *cobra.Command, args []string) error {
	client, worldID, err := getWorldsClient(args[0])
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := client.GetWorld(ctx, &v1.GetWorldRequest{Id: worldID})
	if err != nil {
		return fmt.Errorf("failed to get world: %w", err)
	}

	formatter := NewOutputFormatter()

	tileCount := 0
	unitCount := 0
	if resp.WorldData != nil {
		tileCount = len(resp.WorldData.TilesMap)
		unitCount = len(resp.WorldData.UnitsMap)
	}

	if formatter.JSON {
		data := map[string]any{
			"id":          resp.World.Id,
			"name":        resp.World.Name,
			"description": resp.World.Description,
			"difficulty":  resp.World.Difficulty,
			"tags":        resp.World.Tags,
			"tiles":       tileCount,
			"units":       unitCount,
			"version":     int64(0),
		}
		if resp.WorldData != nil {
			data["version"] = resp.WorldData.Version
		}
		if resp.World.CreatedAt != nil {
			data["created_at"] = resp.World.CreatedAt.AsTime().Format("2006-01-02 15:04:05")
		}
		if resp.World.UpdatedAt != nil {
			data["updated_at"] = resp.World.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")
		}
		return formatter.PrintJSON(data)
	}

	// Text output
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("World: %s\n", resp.World.Name))
	sb.WriteString(fmt.Sprintf("  ID:          %s\n", resp.World.Id))
	if resp.World.Description != "" {
		sb.WriteString(fmt.Sprintf("  Description: %s\n", resp.World.Description))
	}
	if resp.World.Difficulty != "" {
		sb.WriteString(fmt.Sprintf("  Difficulty:  %s\n", resp.World.Difficulty))
	}
	if len(resp.World.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("  Tags:        %s\n", strings.Join(resp.World.Tags, ", ")))
	}
	sb.WriteString(fmt.Sprintf("  Tiles:       %d\n", tileCount))
	sb.WriteString(fmt.Sprintf("  Units:       %d\n", unitCount))
	if resp.WorldData != nil {
		sb.WriteString(fmt.Sprintf("  Version:     %d\n", resp.WorldData.Version))
	}
	if resp.World.CreatedAt != nil {
		sb.WriteString(fmt.Sprintf("  Created:     %s\n", resp.World.CreatedAt.AsTime().Format("2006-01-02 15:04:05")))
	}
	if resp.World.UpdatedAt != nil {
		sb.WriteString(fmt.Sprintf("  Updated:     %s\n", resp.World.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")))
	}

	// Show player breakdown if there are units
	if resp.WorldData != nil && len(resp.WorldData.UnitsMap) > 0 {
		playerUnits := make(map[int32]int)
		for _, u := range resp.WorldData.UnitsMap {
			if u != nil {
				playerUnits[u.Player]++
			}
		}
		sb.WriteString("\n  Players:\n")
		for p := int32(1); p <= 8; p++ {
			if count, ok := playerUnits[p]; ok {
				sb.WriteString(fmt.Sprintf("    Player %d: %d unit(s)\n", p, count))
			}
		}
	}

	return formatter.PrintText(sb.String())
}

func runWorldsShow(cmd *cobra.Command, args []string) error {
	client, worldID, err := getWorldsClient(args[0])
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := client.GetWorld(ctx, &v1.GetWorldRequest{Id: worldID})
	if err != nil {
		return fmt.Errorf("failed to get world: %w", err)
	}

	if resp.WorldData == nil {
		return fmt.Errorf("world %s has no data", worldID)
	}

	// Ensure map data is migrated
	lib.MigrateWorldData(resp.WorldData)

	// Create theme and renderer
	theme := themes.NewDefaultTheme(lib.DefaultRulesEngine().GetCityTerrains())
	renderer, err := themes.NewPNGWorldRenderer(theme)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	options := lib.DefaultRenderOptions()
	options.ShowUnitLabels = worldShowLabels
	options.ShowTileLabels = worldShowTileLabels

	pngData, _, err := renderer.Render(resp.WorldData.TilesMap, resp.WorldData.UnitsMap, options)
	if err != nil {
		return fmt.Errorf("failed to render map: %w", err)
	}

	// Save to file or display inline
	if worldShowOutput != "" {
		if err := os.WriteFile(worldShowOutput, pngData, 0644); err != nil {
			return fmt.Errorf("failed to write image to %s: %w", worldShowOutput, err)
		}
		fmt.Printf("Map saved to %s\n", worldShowOutput)
		return nil
	}

	// Print world name before inline image
	fmt.Printf("%s (%s)\n", resp.World.Name, resp.World.Id)

	// iTerm2 inline image
	encoded := base64.StdEncoding.EncodeToString(pngData)
	fmt.Printf("\033]1337;File=inline=1;preserveAspectRatio=1:%s\a", encoded)
	fmt.Println()

	return nil
}

// worldSummaryMap converts a World proto to a summary map for JSON output
func worldSummaryMap(w *v1.World) map[string]any {
	m := map[string]any{
		"id":         w.Id,
		"name":       w.Name,
		"difficulty": w.Difficulty,
	}
	if w.Description != "" {
		m["description"] = w.Description
	}
	return m
}

// truncate shortens a string to maxLen, adding "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
