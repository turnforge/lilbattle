package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami [profile]",
	Short: "Show current authentication status",
	Long: `Show the current profile and authentication status.

If no profile is specified, shows the current active profile.
Use --profile flag or specify a profile name to see status of a specific profile.

Examples:
  ww whoami              # Show current profile status
  ww whoami prod         # Show status of 'prod' profile`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func runWhoami(cmd *cobra.Command, args []string) error {
	store, err := getProfileStore()
	if err != nil {
		return fmt.Errorf("failed to initialize profile store: %w", err)
	}

	formatter := NewOutputFormatter()

	var profileNameArg string
	if len(args) > 0 {
		profileNameArg = args[0]
	} else {
		// Check --profile flag first
		profileNameArg = getProfileName()
		if profileNameArg == "" {
			// Then check current profile
			profileNameArg, err = store.GetCurrentProfile()
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}
		}
	}

	if profileNameArg == "" {
		if formatter.JSON {
			return formatter.PrintJSON(map[string]any{
				"authenticated":   false,
				"current_profile": "",
				"message":         "no profile selected",
			})
		}
		fmt.Println("No profile selected.")
		fmt.Println("\nUse 'ww login <profile> --host <url>' to create a profile.")
		fmt.Println("Use 'ww select <profile>' to set an active profile.")

		// Show available profiles
		profiles, _ := store.ListProfiles()
		if len(profiles) > 0 {
			fmt.Println("\nAvailable profiles:")
			for _, name := range profiles {
				fmt.Printf("  %s\n", name)
			}
		}
		return nil
	}

	profile, err := store.LoadProfile(profileNameArg)
	if err != nil {
		return fmt.Errorf("profile '%s' does not exist", profileNameArg)
	}

	creds, _ := store.LoadCredentials(profileNameArg)
	currentProfile, _ := store.GetCurrentProfile()

	if formatter.JSON {
		result := map[string]any{
			"profile": profileNameArg,
			"host":    profile.Host,
			"email":   profile.Email,
			"current": profileNameArg == currentProfile,
		}

		if creds != nil {
			result["authenticated"] = !creds.IsExpired()
			result["user_id"] = creds.UserID
			result["user_email"] = creds.UserEmail
			result["expires_at"] = creds.ExpiresAt
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

	if creds != nil {
		if creds.UserEmail != "" && creds.UserEmail != profile.Email {
			fmt.Printf("  User: %s\n", creds.UserEmail)
		}

		if creds.IsExpired() {
			fmt.Printf("  Status: EXPIRED (expired %s ago)\n", time.Since(creds.ExpiresAt).Round(time.Minute))
			fmt.Printf("\nUse 'ww login %s' to re-authenticate.\n", profileNameArg)
		} else {
			remaining := time.Until(creds.ExpiresAt)
			fmt.Printf("  Status: Authenticated (expires in %s)\n", remaining.Round(time.Minute))
		}
	} else {
		fmt.Printf("  Status: Not authenticated\n")
		fmt.Printf("\nUse 'ww login %s' to authenticate.\n", profileNameArg)
	}

	return nil
}
