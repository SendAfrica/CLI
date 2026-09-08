// Package cmd defines all Cobra commands for the SendAfrica CLI.
//
// Each subcommand file (auth.go, sms.go, contacts.go, etc.) implements a
// group of API endpoints. The rootCmd wires them together and resolves
// configuration (credentials, API URL, output format) before any command runs.
//
// Common helpers:
//   - requireAuth() / apiClient() — check credentials and build a client
//   - printer() — returns an output.Printer for the resolved format
//   - decodeData() / decodePaginated() — unwrap the API response envelope
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cameltech/sendafrica-cli/internal/client"
	"github.com/cameltech/sendafrica-cli/internal/config"
	"github.com/cameltech/sendafrica-cli/internal/output"
)

var rootCmd = &cobra.Command{
	Use:   "sendafrica",
	Short: "SendAfrica CLI — manage SMS, credits, contacts, campaigns, and more",
	Long: `SendAfrica CLI

A command-line interface for the SendAfrica API. All non-admin endpoints are
supported. Admin endpoints (/v1/admin/*) are intentionally excluded.

Authentication:
  - API key:   export SENDAFRIA_API_KEY=SA-...
  - JWT token: export SENDAFRIA_JWT_TOKEN=...
  - Or use flags: --api-key, --token
  - Or use profiles: sendafrica config add-profile

Note: The CLI never reads .env files. Credentials are only sourced from
the config file, environment variables, or flags.`,
}

var (
	flagAPIURL   string
	flagAPIKey   string
	flagJWTToken string
	flagProfile  string
	flagOutput   string
	resolvedCfg  *config.ResolvedConfig
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "", "output format: table|json|yaml")
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "SendAfrica API base URL")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&flagJWTToken, "token", "", "JWT access token for authentication")
	rootCmd.PersistentFlags().StringVar(&flagProfile, "profile", "", "config profile to use")

	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
	_ = viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	// Public endpoints
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(packagesCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(ratesCmd)

	// Auth endpoints
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(verifyEmailCmd)
	rootCmd.AddCommand(sendVerificationEmailCmd)
	rootCmd.AddCommand(resetPasswordCmd)
	rootCmd.AddCommand(resetPasswordConfirmCmd)
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(changePasswordCmd)
	rootCmd.AddCommand(verifyPhoneCmd)
	rootCmd.AddCommand(sendPhoneOTPCmd)
	rootCmd.AddCommand(apiKeysCmd)

	// SMS endpoints
	rootCmd.AddCommand(smsSendCmd)
	rootCmd.AddCommand(smsSendBulkCmd)
	rootCmd.AddCommand(smsLogsCmd)

	// Credits endpoints
	rootCmd.AddCommand(creditsBalanceCmd)
	rootCmd.AddCommand(creditsHistoryCmd)

	// Contacts endpoints
	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsCreateCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)
	contactsCmd.AddCommand(contactsDeleteCmd)
	contactsCmd.AddCommand(contactsDuplicateCheckCmd)
	contactsCmd.AddCommand(contactsImportCmd)
	contactsCmd.AddCommand(contactsListContactsCmd)
	contactsCmd.AddCommand(contactsCreateContactCmd)
	contactsCmd.AddCommand(contactsGetContactCmd)
	contactsCmd.AddCommand(contactsUpdateContactCmd)
	contactsCmd.AddCommand(contactsDeleteContactCmd)
	contactsCmd.AddCommand(contactsExportCmd)
	contactsCmd.AddCommand(contactsAddPhoneCmd)
	contactsCmd.AddCommand(contactsDeletePhoneCmd)
	contactsCmd.AddCommand(contactsGoogleStatusCmd)
	contactsCmd.AddCommand(contactsGoogleSyncCmd)
	contactsCmd.AddCommand(contactsGoogleDisconnectCmd)
	rootCmd.AddCommand(contactsCmd)

	// Campaigns endpoints
	campaignsCmd.AddCommand(campaignsListCmd)
	campaignsCmd.AddCommand(campaignsCreateCmd)
	campaignsCmd.AddCommand(campaignsGetCmd)
	campaignsCmd.AddCommand(campaignsUpdateCmd)
	campaignsCmd.AddCommand(campaignsCancelCmd)
	campaignsCmd.AddCommand(campaignsDeleteCmd)
	campaignsCmd.AddCommand(campaignsScheduleCmd)
	campaignsCmd.AddCommand(campaignsRecipientsCmd)
	rootCmd.AddCommand(campaignsCmd)

	// Payments endpoints
	paymentsCmd.AddCommand(paymentsInitiateCmd)
	rootCmd.AddCommand(paymentsCmd)

	// Vouchers endpoints
	vouchersCmd.AddCommand(vouchersRateCmd)
	vouchersCmd.AddCommand(vouchersPurchaseCmd)
	rootCmd.AddCommand(vouchersCmd)

	// Sender IDs endpoints
	senderIDsCmd.AddCommand(senderIDsRequirementsCmd)
	senderIDsCmd.AddCommand(senderIDsUsableCmd)
	senderIDsCmd.AddCommand(senderIDsListCmd)
	senderIDsCmd.AddCommand(senderIDsCreateCmd)
	senderIDsCmd.AddCommand(senderIDsGetCmd)
	rootCmd.AddCommand(senderIDsCmd)

	// Notifications endpoints
	notificationsCmd.AddCommand(notificationsListCmd)
	notificationsCmd.AddCommand(notificationsUnreadCountCmd)
	notificationsCmd.AddCommand(notificationsGetCmd)
	notificationsCmd.AddCommand(notificationsMarkReadCmd)
	notificationsCmd.AddCommand(notificationsReadAllCmd)
	rootCmd.AddCommand(notificationsCmd)

	// Support endpoints
	rootCmd.AddCommand(supportChatCmd)

	// Config management
	rootCmd.AddCommand(configCmd)
}

// initConfig runs before any command. It loads the config profile and
// resolves env-var/flag overrides. It NEVER loads .env files.
func initConfig() {
	viper.SetEnvPrefix("sendafrica")
	viper.AutomaticEnv()

	resolved, err := config.Resolve(
		flagProfile,
		flagAPIURL,
		flagAPIKey,
		flagJWTToken,
		flagOutput,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning:", err)
		return
	}
	resolvedCfg = resolved
}

// apiClient returns an authenticated SendAfrica API client from resolved config.
func apiClient() (*client.Client, error) {
	if resolvedCfg == nil {
		return nil, fmt.Errorf("configuration not loaded; set SENDAFRIA_API_KEY or use --api-key")
	}
	if resolvedCfg.APIKey == "" && resolvedCfg.JWTToken == "" {
		return nil, fmt.Errorf("no credentials found; set SENDAFRIA_API_KEY or SENDAFRIA_JWT_TOKEN, or use --api-key/--token flags")
	}
	return client.New(resolvedCfg.APIURL, resolvedCfg.APIKey, resolvedCfg.JWTToken), nil
}

// printer returns an output.Printer configured with the resolved output format.
func printer() *output.Printer {
	format := output.FormatTable
	if resolvedCfg != nil && resolvedCfg.OutputFormat != "" {
		format = resolvedCfg.OutputFormat
	}
	return output.New(format)
}

// ensureAuth checks that credentials are present, returning an error otherwise.
func ensureAuth() error {
	if resolvedCfg == nil {
		return fmt.Errorf("configuration not loaded")
	}
	if resolvedCfg.APIKey == "" && resolvedCfg.JWTToken == "" {
		return fmt.Errorf("no credentials found; set SENDAFRIA_API_KEY or SENDAFRIA_JWT_TOKEN, or use --api-key/--token flags")
	}
	return nil
}
