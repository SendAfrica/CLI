package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
	"github.com/cameltech/sendafrica-cli/internal/config"
	"github.com/cameltech/sendafrica-cli/internal/output"
)

// publicClient creates a client for public endpoints that don't require auth.
func publicClient() *client.Client {
	apiURL := config.DefaultAPIURL
	if resolvedCfg != nil && resolvedCfg.APIURL != "" {
		apiURL = resolvedCfg.APIURL
	}
	return client.New(apiURL, "", "")
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API server health",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := publicClient()

		resp, env, err := c.DoWithResponse(client.RequestOpts{
			Method: "GET",
			Path:   "/health",
		})
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if env.Error != nil {
			return env.Error
		}

		var result map[string]interface{}
		if err := decodeData(env.Data, &result); err != nil {
			return err
		}

		if printer().Format() == output.FormatJSON {
			return printer().Print(result)
		}
		fmt.Printf("Status: %v\n", result["status"])
		if ts, ok := result["timestamp"]; ok {
			fmt.Printf("Timestamp: %v\n", ts)
		}
		return nil
	},
}

var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List active credit packages (public)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := publicClient()
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/packages"})
		if err != nil {
			return err
		}

		var result []api.Package
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List active SMS templates (public)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := publicClient()
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/templates"})
		if err != nil {
			return err
		}

		var result []api.SMSTemplate
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var ratesCmd = &cobra.Command{
	Use:   "rates",
	Short: "List supported-country SMS rate card (public)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := publicClient()

		country, _ := cmd.Flags().GetString("country")
		path := "/v1/rates"
		if country != "" {
			path = "/v1/rates/" + country
		}

		data, err := c.Do(client.RequestOpts{Method: "GET", Path: path})
		if err != nil {
			return err
		}

		if country != "" {
			var result api.Rate
			if err := decodeData(data, &result); err != nil {
				return err
			}
			return printer().Print(result)
		}

		var result []api.Rate
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

func init() {
	ratesCmd.Flags().String("country", "", "filter to a specific country code (e.g. TZ)")
}
