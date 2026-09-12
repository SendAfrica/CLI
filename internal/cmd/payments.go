package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
)

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
			Provider: "manual",
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

var vouchersSendOTPCmd = &cobra.Command{Use: "send-otp [phone]", Short: "Send an OTP to a declared payer phone", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}
	c, err := apiClient()
	if err != nil {
		return err
	}
	data, err := c.Do(client.RequestOpts{Method: "POST", Path: "/v1/vouchers/otp/send", Body: api.DeclaredPhoneOTPRequest{Phone: args[0]}})
	if err != nil {
		return err
	}
	return printer().Print(rawMessage(data))
}}

var vouchersVerifyOTPCmd = &cobra.Command{Use: "verify-otp [phone] [otp]", Short: "Verify a declared payer phone OTP", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuth(); err != nil {
		return err
	}
	c, err := apiClient()
	if err != nil {
		return err
	}
	data, err := c.Do(client.RequestOpts{Method: "POST", Path: "/v1/vouchers/otp/verify", Body: api.DeclaredPhoneOTPRequest{Phone: args[0], OTP: args[1]}})
	if err != nil {
		return err
	}
	return printer().Print(rawMessage(data))
}}
