package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage payment orders",
}

var paymentsInitiateCmd = &cobra.Command{
	Use:   "initiate [packageId]",
	Short: "Initiate a package top-up payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		req := api.InitiatePaymentRequest{
			PackageID: args[0],
		}
		if method, _ := cmd.Flags().GetString("method"); method != "" {
			// method is informational; Snippe is the supported provider
			_ = method
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/payments/",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.PaymentResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

func init() {
	paymentsInitiateCmd.Flags().String("method", "", "payment method (e.g. snippe)")
}

// --- Vouchers ---

var vouchersCmd = &cobra.Command{
	Use:   "vouchers",
	Short: "Manage voucher-based credit top-ups",
}

var vouchersRateCmd = &cobra.Command{
	Use:   "rate",
	Short: "Get voucher rate tiers",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/vouchers/rate"})
		if err != nil {
			return err
		}
		var result api.VoucherRateResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var vouchersPurchaseCmd = &cobra.Command{
	Use:   "purchase [amount]",
	Short: "Purchase credits via voucher",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		var amount int
		if _, err := fmt.Sscanf(args[0], "%d", &amount); err != nil {
			return fmt.Errorf("invalid amount: %s", args[0])
		}

		req := api.PurchaseVoucherRequest{
			Amount:   amount,
			Currency: "TZS",
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/vouchers/",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.PurchaseVoucherResponse
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}
