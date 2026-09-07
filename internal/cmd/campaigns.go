package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage SMS campaigns",
}

var campaignsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all campaigns",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		params := map[string]string{}
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		if page > 0 {
			params["page"] = fmt.Sprintf("%d", page)
		}
		if perPage > 0 {
			params["per_page"] = fmt.Sprintf("%d", perPage)
		}
		data, err := c.Do(client.RequestOpts{
			Method:      "GET",
			Path:        "/v1/campaigns",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.Campaign
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var campaignsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		message, _ := cmd.Flags().GetString("message")
		from, _ := cmd.Flags().GetString("from")

		req := api.CreateCampaignRequest{
			Name:    name,
			Message: message,
		}
		if from != "" {
			req.From = from
		}

		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/campaigns",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.Campaign
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var campaignsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a campaign by ID",
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
			Path:   "/v1/campaigns/" + args[0],
		})
		if err != nil {
			return err
		}
		var result api.Campaign
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var campaignsUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a campaign (name, message, or sender)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		req := api.UpdateCampaignRequest{}
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			req.Name = name
		}
		if msg, _ := cmd.Flags().GetString("message"); msg != "" {
			req.Message = msg
		}
		if from, _ := cmd.Flags().GetString("from"); from != "" {
			req.From = from
		}

		data, err := c.Do(client.RequestOpts{
			Method: "PATCH",
			Path:   "/v1/campaigns/" + args[0],
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var campaignsCancelCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Cancel a campaign",
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
			Method: "POST",
			Path:   "/v1/campaigns/" + args[0] + "/cancel",
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var campaignsScheduleCmd = &cobra.Command{
	Use:   "schedule [id] [listId]",
	Short: "Attach a contact list and schedule a campaign",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		req := api.AttachContactListRequest{ListID: args[1]}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/campaigns/" + args[0] + "/schedule",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var campaignsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a campaign",
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
			Method: "DELETE",
			Path:   "/v1/campaigns/" + args[0],
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var campaignsRecipientsCmd = &cobra.Command{
	Use:   "recipients [id]",
	Short: "List campaign recipient delivery statuses",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		params := map[string]string{}
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		status, _ := cmd.Flags().GetString("status")
		if page > 0 {
			params["page"] = fmt.Sprintf("%d", page)
		}
		if perPage > 0 {
			params["per_page"] = fmt.Sprintf("%d", perPage)
		}
		if status != "" {
			params["status"] = status
		}
		data, err := c.Do(client.RequestOpts{
			Method:      "GET",
			Path:        "/v1/campaigns/" + args[0] + "/recipients",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.CampaignRecipient
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

func init() {
	campaignsListCmd.Flags().Int("page", 1, "page number")
	campaignsListCmd.Flags().Int("per-page", 50, "items per page")

	campaignsCreateCmd.Flags().String("name", "", "campaign name")
	campaignsCreateCmd.Flags().String("message", "", "message body")
	campaignsCreateCmd.Flags().String("from", "", "sender ID")
	_ = campaignsCreateCmd.MarkFlagRequired("name")
	_ = campaignsCreateCmd.MarkFlagRequired("message")

	campaignsUpdateCmd.Flags().String("name", "", "new campaign name")
	campaignsUpdateCmd.Flags().String("message", "", "new message body")
	campaignsUpdateCmd.Flags().String("from", "", "new sender ID")

	campaignsRecipientsCmd.Flags().Int("page", 1, "page number")
	campaignsRecipientsCmd.Flags().Int("per-page", 50, "items per page")
	campaignsRecipientsCmd.Flags().String("status", "", "filter by status: sent|delivered|failed|pending")
}
