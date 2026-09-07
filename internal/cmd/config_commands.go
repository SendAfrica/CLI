package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration and profiles",
}

var configAddProfileCmd = &cobra.Command{
	Use:   "add-profile [name]",
	Short: "Add or update a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		apiKey, _ := cmd.Flags().GetString("api-key")
		jwt, _ := cmd.Flags().GetString("token")
		apiURL, _ := cmd.Flags().GetString("api-url")

		if apiKey == "" && jwt == "" {
			return fmt.Errorf("at least one of --api-key or --token is required")
		}

		profile := config.Profile{
			APIKey:   apiKey,
			JWTToken: jwt,
		}
		if apiURL != "" {
			profile.APIURL = apiURL
		} else {
			profile.APIURL = config.DefaultAPIURL
		}

		if err := config.AddProfile(name, profile); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Profile %q saved (api_url=%s).\n", name, profile.APIURL)
		return nil
	},
}

var configUseProfileCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch the active profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := config.SetCurrentProfile(name); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Active profile set to %q.\n", name)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Current profile: %s\n", cfg.CurrentProfile)
		for name, p := range cfg.Profiles {
			marker := " "
			if name == cfg.CurrentProfile {
				marker = "*"
			}
			authMethod := "api_key"
			if p.APIKey == "" && p.JWTToken != "" {
				authMethod = "jwt"
			} else if p.APIKey != "" && p.JWTToken != "" {
				authMethod = "both"
			}
			fmt.Fprintf(os.Stderr, "  %s %-20s %s  (%s)\n", marker, name, p.APIURL, authMethod)
		}
		return nil
	},
}

var configDeleteProfileCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := config.DeleteProfile(name); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Profile %q deleted.\n", name)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show resolved configuration (secrets masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := map[string]string{
			"api_url":       resolvedCfg.APIURL,
			"profile":       resolvedCfg.ProfileName,
			"output_format": resolvedCfg.OutputFormat,
			"api_key":       mask(resolvedCfg.APIKey),
			"jwt_token":     mask(resolvedCfg.JWTToken),
		}
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		return nil
	},
}

func mask(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func init() {
	configAddProfileCmd.Flags().String("api-key", "", "API key for this profile")
	configAddProfileCmd.Flags().String("token", "", "JWT token for this profile")
	configAddProfileCmd.Flags().String("api-url", "", "API URL (default: "+config.DefaultAPIURL+")")

	configCmd.AddCommand(configAddProfileCmd)
	configCmd.AddCommand(configUseProfileCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configDeleteProfileCmd)
	configCmd.AddCommand(configShowCmd)
}
