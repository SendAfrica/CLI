package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

var senderIDsCmd = &cobra.Command{
	Use:   "sender-ids",
	Short: "Manage sender IDs",
}

var senderIDsRequirementsCmd = &cobra.Command{
	Use:   "requirements",
	Short: "Get sender ID submission requirements",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/sender-ids/requirements"})
		if err != nil {
			return err
		}
		var result api.SenderIDRequirements
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var senderIDsUsableCmd = &cobra.Command{
	Use:   "usable",
	Short: "List usable sender IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		params := map[string]string{}
		if provider, _ := cmd.Flags().GetString("provider"); provider != "" {
			params["provider"] = provider
		}
		data, err := c.Do(client.RequestOpts{
			Method:      "GET",
			Path:        "/v1/sender-ids/usable",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.SenderID
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var senderIDsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sender IDs registered for this account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/sender-ids"})
		if err != nil {
			return err
		}
		var result []api.SenderID
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var senderIDsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Request a new sender ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		req := api.CreateSenderIDRequest{}
		req.Name, _ = cmd.Flags().GetString("name")
		req.Country, _ = cmd.Flags().GetString("country")
		req.Purpose, _ = cmd.Flags().GetString("purpose")
		req.SampleMessage, _ = cmd.Flags().GetString("sample-message")
		docsJSON, _ := cmd.Flags().GetString("documents")
		if docsJSON != "" {
			if err := jsonUnmarshal([]byte(docsJSON), &req.Documents); err != nil {
				return err
			}
		}

		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/sender-ids",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.SenderID
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var senderIDsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a sender ID by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{
			Method: "GET",
			Path:   "/v1/sender-ids/" + args[0],
		})
		if err != nil {
			return err
		}
		var result api.SenderID
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

func init() {
	senderIDsUsableCmd.Flags().String("provider", "", "filter by provider: swala or africastalking")
	senderIDsCreateCmd.Flags().String("name", "", "sender ID name (3-11 alphanumeric chars)")
	senderIDsCreateCmd.Flags().String("country", "TZ", "country code")
	senderIDsCreateCmd.Flags().String("purpose", "", "purpose (e.g. Transactional, Promotional)")
	senderIDsCreateCmd.Flags().String("sample-message", "", "sample message body (50-500 chars)")
	senderIDsCreateCmd.Flags().String("documents", "", "JSON array of document objects")
	_ = senderIDsCreateCmd.MarkFlagRequired("name")
	_ = senderIDsCreateCmd.MarkFlagRequired("purpose")
	_ = senderIDsCreateCmd.MarkFlagRequired("sample-message")
}
