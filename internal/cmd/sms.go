package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

var smsSendCmd = &cobra.Command{
	Use:   "sms-send",
	Short: "Send a single SMS (API key or JWT)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		to, _ := cmd.Flags().GetString("to")
		message, _ := cmd.Flags().GetString("message")
		from, _ := cmd.Flags().GetString("from")

		req := api.SMSSendRequest{To: to, Message: message}
		if from != "" {
			req.From = from
		}

		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/sms/",
			Body:   req,
		})
		if err != nil {
			return err
		}

		var result api.SMSResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var smsSendBulkCmd = &cobra.Command{
	Use:   "sms-bulk",
	Short: "Send a bulk SMS (JWT or API key)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		toStr, _ := cmd.Flags().GetString("to")
		message, _ := cmd.Flags().GetString("message")
		from, _ := cmd.Flags().GetString("from")

		req := api.SMSBulkSendRequest{
			To:      parseStringSlice(toStr),
			Message: message,
		}
		if from != "" {
			req.From = from
		}

		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/sms/bulk",
			Body:   req,
		})
		if err != nil {
			return err
		}

		var result api.SMSResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var smsLogsCmd = &cobra.Command{
	Use:   "sms-logs",
	Short: "List SMS message logs",
	Args:  cobra.NoArgs,
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
			Path:        "/v1/sms/logs",
			QueryParams: params,
		})
		if err != nil {
			return err
		}

		var result []api.SMSLog
		if _, err := decodePaginated(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

// Credits commands

var creditsBalanceCmd = &cobra.Command{
	Use:   "credits-balance",
	Short: "Get current credit balance",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/credits/balance"})
		if err != nil {
			return err
		}
		var result api.CreditBalance
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var creditsHistoryCmd = &cobra.Command{
	Use:   "credits-history",
	Short: "List credit transactions",
	Args:  cobra.NoArgs,
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
			Path:        "/v1/credits/history",
			QueryParams: params,
		})
		if err != nil {
			return err
		}
		var result []api.Transaction
		if _, err := decodePaginated(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

func init() {
	smsSendCmd.Flags().String("to", "", "recipient phone number(s), comma-separated")
	smsSendCmd.Flags().String("message", "", "message body")
	smsSendCmd.Flags().String("from", "", "sender ID")
	_ = smsSendCmd.MarkFlagRequired("to")
	_ = smsSendCmd.MarkFlagRequired("message")

	smsSendBulkCmd.Flags().String("to", "", "recipient phone numbers, comma-separated")
	smsSendBulkCmd.Flags().String("message", "", "message body")
	smsSendBulkCmd.Flags().String("from", "", "sender ID")
	_ = smsSendBulkCmd.MarkFlagRequired("to")
	_ = smsSendBulkCmd.MarkFlagRequired("message")

	smsLogsCmd.Flags().Int("page", 1, "page number")
	smsLogsCmd.Flags().Int("per-page", 50, "items per page")

	creditsHistoryCmd.Flags().Int("page", 1, "page number")
	creditsHistoryCmd.Flags().Int("per-page", 50, "items per page")
}
