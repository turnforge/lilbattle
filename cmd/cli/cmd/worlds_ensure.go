package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turnforge/lilbattle/lib"
)

var (
	worldEnsureDataDir  string
	worldEnsureDataJSON string
	worldEnsureMetaJSON string
)

// worldsEnsureCmd exposes lib.EnsureWorldExists as a CLI so ops flows
// and shell scripts (e.g. tests/e2e's scripts/seed-worlds.sh) can
// provision worlds against a server without embedding curl+protojson
// gymnastics in the script.
//
// Behavior mirrors the library:
//   - world missing + no fixture   → error (nothing to seed)
//   - world missing + fixture      → CreateWorld from fixture
//   - world exists  + no fixture   → success (presence probe)
//   - world exists  + fixture      → content hash comparison; mismatch is an
//                                    error, never an overwrite
var worldsEnsureCmd = &cobra.Command{
	Use:   "ensure <worldId>",
	Short: "Ensure a world exists on the target server; optionally seed or validate content",
	Long: `Probe the target for <worldId>. Behavior depends on whether a fixture
was supplied:

  --data-dir <path>       Directory with metadata.json + data.json (FS layout).
  --data-json <path>      Explicit path to data.json (WorldData protojson).
                          Pair with --meta-json for a full fixture; if omitted,
                          --data-dir + auto-lookup is normally simpler.
  --meta-json <path>      Explicit path to metadata.json (World protojson).

Cases:
  - Missing on target, no fixture  → error (nothing to seed)
  - Missing on target, fixture set → CreateWorld from fixture
  - Present on target, no fixture  → success
  - Present on target, fixture set → hash-compare tiles+units; mismatch is
                                     a hard error, never an overwrite.

Examples:
  ww worlds ensure 7e5016a4                           # probe only
  ww worlds ensure 7e5016a4 --data-dir tests/e2e/fixtures/worlds/7e5016a4/
  ww worlds ensure 7e5016a4 --profile prod            # against prod`,
	Args: cobra.ExactArgs(1),
	RunE: runWorldsEnsure,
}

func init() {
	worldsCmd.AddCommand(worldsEnsureCmd)
	worldsEnsureCmd.Flags().StringVar(&worldEnsureDataDir, "data-dir", "",
		"Directory containing metadata.json + data.json (FS-backend layout)")
	worldsEnsureCmd.Flags().StringVar(&worldEnsureDataJSON, "data-json", "",
		"Path to a WorldData protojson file")
	worldsEnsureCmd.Flags().StringVar(&worldEnsureMetaJSON, "meta-json", "",
		"Path to a World metadata protojson file (pair with --data-json)")
}

func runWorldsEnsure(cmd *cobra.Command, args []string) error {
	worldID := args[0]

	client, _, err := getWorldsClient("")
	if err != nil {
		return err
	}

	fixture := buildFixtureFromFlags()

	ctx := context.Background()
	if err := lib.EnsureWorldExists(ctx, client, worldID, fixture); err != nil {
		return err
	}

	if isVerbose() {
		fmt.Printf("[VERBOSE] world %s ensured on target\n", worldID)
	}
	return nil
}

// buildFixtureFromFlags returns nil when no fixture flag was set, so
// EnsureWorldExists treats the call as a presence probe. Explicit paths
// override --data-dir when both are set (unlikely but not forbidden).
func buildFixtureFromFlags() *lib.WorldFixture {
	if worldEnsureDataDir == "" && worldEnsureDataJSON == "" && worldEnsureMetaJSON == "" {
		return nil
	}
	return &lib.WorldFixture{
		Dir:              worldEnsureDataDir,
		DataJSONPath:     worldEnsureDataJSON,
		MetadataJSONPath: worldEnsureMetaJSON,
	}
}
