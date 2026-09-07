package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Manage notifications",
}

var notificationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications for the current account",
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
			Path:        "/v1/notifications",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.Notification
		if _, err := decodePaginated(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var notificationsUnreadCountCmd = &cobra.Command{
	Use:   "unread-count",
	Short: "Get unread notification count",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/notifications/unread-count"})
		if err != nil {
			return err
		}
		var result map[string]int
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var notificationsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single notification",
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
			Path:   "/v1/notifications/" + args[0],
		})
		if err != nil {
			return err
		}
		var result api.Notification
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var notificationsMarkReadCmd = &cobra.Command{
	Use:   "mark-read [id]",
	Short: "Mark a notification as read",
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
			Method: "PATCH",
			Path:   "/v1/notifications/" + args[0] + "/read",
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var notificationsReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all notifications as read",
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
			Path:   "/v1/notifications/read-all",
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

func init() {
	notificationsListCmd.Flags().Int("page", 1, "page number")
	notificationsListCmd.Flags().Int("per-page", 50, "items per page")
}
