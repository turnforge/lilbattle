package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile parent command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
	Long: `Manage authentication profiles.

Subcommands:
  show    Show details of a profile
  delete  Delete a profile`,
}

// profileShowCmd represents the profile show command
var profileShowCmd = &cobra.Command{
	Use:   "show [profile]",
	Short: "Show profile details",
	Long: `Show detailed information about a profile.

If no profile is specified, shows the current active profile.

Examples:
  ww profile show prod     # Show details of 'prod' profile
  ww profile show          # Show current profile details`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProfileShow,
}

// profileDeleteCmd represents the profile delete command
var profileDeleteCmd = &cobra.Command{
	Use:   "delete <profile>",
	Short: "Delete a profile",
	Long: `Delete a profile and all its stored credentials.

This will prompt for confirmation unless --confirm=false is set.

Examples:
  ww profile delete prod              # Delete 'prod' profile (prompts)
  ww profile delete prod --confirm=false  # Delete without prompting`,
	Args: cobra.ExactArgs(1),
	RunE: runProfileDelete,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileDeleteCmd)
}

func runProfileShow(cmd *cobra.Command, args []string) error {
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
			fmt.Println("Use 'ww profile show <profile>' or 'ww select <profile>' first.")
			return nil
		}
	}

	profile, err := store.LoadProfile(profileNameArg)
	if err != nil {
		return fmt.Errorf("profile '%s' does not exist", profileNameArg)
	}

	creds, _ := store.LoadCredentials(profileNameArg)
	currentProfile, _ := store.GetCurrentProfile()

	if formatter.JSON {
		result := map[string]any{
			"name":    profileNameArg,
			"host":    profile.Host,
			"email":   profile.Email,
			"current": profileNameArg == currentProfile,
		}

		if profile.Password != "" {
			result["password_saved"] = true
		}

		if creds != nil {
			result["authenticated"] = !creds.IsExpired()
			result["user_id"] = creds.UserID
			result["user_email"] = creds.UserEmail
			result["expires_at"] = creds.ExpiresAt
			result["created_at"] = creds.CreatedAt
		} else {
			result["authenticated"] = false
		}

		return formatter.PrintJSON(result)
	}

	fmt.Printf("Profile: %s", profileNameArg)
	if profileNameArg == currentProfile {
		fmt.Printf(" (current)")
	}
	fmt.Println()

	fmt.Printf("  Host: %s\n", profile.Host)
	if profile.Email != "" {
		fmt.Printf("  Email: %s\n", profile.Email)
	}
	if profile.Password != "" {
		fmt.Printf("  Password: saved\n")
	}

	if creds != nil {
		fmt.Println()
		fmt.Printf("  Authentication:\n")
		if creds.UserID != "" {
			fmt.Printf("    User ID: %s\n", creds.UserID)
		}
		if creds.UserEmail != "" {
			fmt.Printf("    User Email: %s\n", creds.UserEmail)
		}

		if creds.IsExpired() {
			fmt.Printf("    Status: EXPIRED (expired %s ago)\n", time.Since(creds.ExpiresAt).Round(time.Minute))
		} else {
			remaining := time.Until(creds.ExpiresAt)
			fmt.Printf("    Status: Valid (expires in %s)\n", remaining.Round(time.Minute))
		}
		fmt.Printf("    Created: %s\n", creds.CreatedAt.Format(time.RFC3339))
	} else {
		fmt.Println()
		fmt.Printf("  Authentication: not authenticated\n")
	}

	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	profileNameArg := args[0]

	store, err := getProfileStore()
	if err != nil {
		return fmt.Errorf("failed to initialize profile store: %w", err)
	}

	formatter := NewOutputFormatter()

	// Verify profile exists
	_, err = store.LoadProfile(profileNameArg)
	if err != nil {
		return fmt.Errorf("profile '%s' does not exist", profileNameArg)
	}

	// Confirm deletion
	if shouldConfirm() && !formatter.JSON {
		fmt.Printf("Delete profile '%s' and all its credentials? [y/N]: ", profileNameArg)
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
			"deleted": true,
		})
	}

	fmt.Printf("Profile '%s' deleted.\n", profileNameArg)
	return nil
}
