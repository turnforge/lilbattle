package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// profilesCmd represents the profiles command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List all configured profiles",
	Long: `List all configured profiles with their authentication status.

The current active profile is marked with an asterisk (*).

Examples:
  ww profiles              # List all profiles
  ww profiles --json       # Output as JSON`,
	RunE: runProfiles,
}

func init() {
	rootCmd.AddCommand(profilesCmd)
}

func runProfiles(cmd *cobra.Command, args []string) error {
	store, err := getProfileStore()
	if err != nil {
		return fmt.Errorf("failed to initialize profile store: %w", err)
	}

	formatter := NewOutputFormatter()

	profiles, err := store.ListProfiles()
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	currentProfile, _ := store.GetCurrentProfile()

	if len(profiles) == 0 {
		if formatter.JSON {
			return formatter.PrintJSON(map[string]any{
				"profiles":        []any{},
				"current_profile": "",
			})
		}
		fmt.Println("No profiles configured.")
		fmt.Println("Use 'ww login <profile> --host <url>' to create a profile.")
		return nil
	}

	if formatter.JSON {
		profileList := make([]map[string]any, 0, len(profiles))
		for _, name := range profiles {
			profile, _ := store.LoadProfile(name)
			creds, _ := store.LoadCredentials(name)

			entry := map[string]any{
				"name":    name,
				"current": name == currentProfile,
			}

			if profile != nil {
				entry["host"] = profile.Host
				entry["email"] = profile.Email
			}

			if creds != nil {
				entry["authenticated"] = !creds.IsExpired()
				entry["expires_at"] = creds.ExpiresAt
			} else {
				entry["authenticated"] = false
			}

			profileList = append(profileList, entry)
		}
		return formatter.PrintJSON(map[string]any{
			"profiles":        profileList,
			"current_profile": currentProfile,
		})
	}

	fmt.Println("Configured profiles:")
	for _, name := range profiles {
		profile, _ := store.LoadProfile(name)
		creds, _ := store.LoadCredentials(name)

		marker := "  "
		if name == currentProfile {
			marker = "* "
		}

		status := "no credentials"
		if creds != nil {
			if creds.IsExpired() {
				status = "expired"
			} else {
				remaining := time.Until(creds.ExpiresAt)
				if remaining < 24*time.Hour {
					status = fmt.Sprintf("expires in %s", remaining.Round(time.Minute))
				} else {
					status = "valid"
				}
			}
		}

		host := "not set"
		if profile != nil && profile.Host != "" {
			host = profile.Host
		}

		fmt.Printf("%s%s\n", marker, name)
		fmt.Printf("    Host: %s\n", host)
		fmt.Printf("    Status: %s\n", status)
	}

	if currentProfile != "" {
		fmt.Printf("\n* = current profile\n")
	} else {
		fmt.Printf("\nNo profile selected. Use 'ww select <profile>' to set one.\n")
	}

	return nil
}
