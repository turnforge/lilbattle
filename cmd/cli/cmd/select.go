package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select <profile>",
	Short: "Set the active profile for subsequent commands",
	Long: `Set a profile as the active profile for all subsequent CLI commands.

The active profile's host and credentials will be used automatically
unless overridden by --server, --profile flags, or LILBATTLE_SERVER env var.

Examples:
  ww select prod           # Set 'prod' as the active profile
  ww select local-dev      # Switch to local development profile

Use 'ww profiles' to list all available profiles.`,
	Args: cobra.ExactArgs(1),
	RunE: runSelect,
}

func init() {
	rootCmd.AddCommand(selectCmd)
}

func runSelect(cmd *cobra.Command, args []string) error {
	profileNameArg := args[0]

	store, err := getProfileStore()
	if err != nil {
		return fmt.Errorf("failed to initialize profile store: %w", err)
	}

	formatter := NewOutputFormatter()

	// Verify profile exists
	profile, err := store.LoadProfile(profileNameArg)
	if err != nil {
		return fmt.Errorf("profile '%s' does not exist. Use 'ww login %s' to create it", profileNameArg, profileNameArg)
	}

	// Set as current profile
	if err := store.SetCurrentProfile(profileNameArg); err != nil {
		return fmt.Errorf("failed to set current profile: %w", err)
	}

	if formatter.JSON {
		return formatter.PrintJSON(map[string]any{
			"profile":  profileNameArg,
			"host":     profile.Host,
			"email":    profile.Email,
			"selected": true,
		})
	}

	fmt.Printf("Switched to profile '%s'\n", profileNameArg)
	fmt.Printf("  Host: %s\n", profile.Host)
	if profile.Email != "" {
		fmt.Printf("  Email: %s\n", profile.Email)
	}

	return nil
}
