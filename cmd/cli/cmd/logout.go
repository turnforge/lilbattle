package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout [profile]",
	Short: "Remove a profile and its stored credentials",
	Long: `Remove a profile and all its stored credentials.

If no profile is specified, removes the current active profile.
This will prompt for confirmation unless --confirm=false is set.

Examples:
  ww logout                        # Remove current profile (prompts)
  ww logout prod                   # Remove 'prod' profile
  ww logout --confirm=false        # Remove without prompting`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	store, err := getProfileStore()
	if err != nil {
		return fmt.Errorf("failed to initialize profile store: %w", err)
	}

	formatter := NewOutputFormatter()

	var profileNameArg string
	if len(args) > 0 {
		profileNameArg = args[0]
	} else {
		// Use current profile
		profileNameArg, err = store.GetCurrentProfile()
		if err != nil {
			return fmt.Errorf("failed to get current profile: %w", err)
		}
		if profileNameArg == "" {
			if formatter.JSON {
				return formatter.PrintJSON(map[string]any{
					"error": "no profile selected",
				})
			}
			fmt.Println("No profile selected.")
			fmt.Println("Use 'ww logout <profile>' to specify a profile, or 'ww select <profile>' first.")
			fmt.Println("\nAvailable profiles:")

			profiles, _ := store.ListProfiles()
			if len(profiles) == 0 {
				fmt.Println("  (none)")
			} else {
				for _, name := range profiles {
					fmt.Printf("  %s\n", name)
				}
			}
			return nil
		}
	}

	// Verify profile exists
	profile, err := store.LoadProfile(profileNameArg)
	if err != nil {
		return fmt.Errorf("profile '%s' does not exist", profileNameArg)
	}

	// Confirm deletion
	if shouldConfirm() && !formatter.JSON {
		fmt.Printf("Remove profile '%s' (%s) and all its credentials? [y/N]: ", profileNameArg, profile.Host)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := store.DeleteProfile(profileNameArg); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	if formatter.JSON {
		return formatter.PrintJSON(map[string]any{
			"profile": profileNameArg,
			"removed": true,
		})
	}

	fmt.Printf("Logged out and removed profile '%s'.\n", profileNameArg)
	return nil
}
