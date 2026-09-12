package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cameltech/sendafrica-cli/internal/api"
	"github.com/cameltech/sendafrica-cli/internal/client"
	"github.com/cameltech/sendafrica-cli/internal/config"
	"github.com/cameltech/sendafrica-cli/internal/output"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with email and password (saves JWT to profile with --save)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		// Login doesn't require existing auth.
		apiURL := config.DefaultAPIURL
		if resolvedCfg != nil && resolvedCfg.APIURL != "" {
			apiURL = resolvedCfg.APIURL
		}
		apiURLFlag, _ := cmd.Flags().GetString("api-url")
		if apiURLFlag != "" {
			apiURL = apiURLFlag
		}

		c := client.New(apiURL, "", "")

		req := api.LoginRequest{Email: email, Password: password}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/login",
			Body:   req,
		})
		if err != nil {
			return err
		}

		var resp api.LoginResponse
		if err := decodeData(data, &resp); err != nil {
			return err
		}

		// Save JWT to the current profile if --save is set.
		if cmd.Flag("save").Changed {
			profileName := "default"
			if resolvedCfg != nil {
				profileName = resolvedCfg.ProfileName
			}
			cfg, loadErr := config.Load()
			if loadErr != nil {
				return fmt.Errorf("loading config: %w", loadErr)
			}
			if p, ok := cfg.Profiles[profileName]; ok {
				p.JWTToken = resp.AccessToken
				p.RefreshToken = resp.RefreshToken
				if p.APIURL == "" {
					p.APIURL = apiURL
				}
				cfg.Profiles[profileName] = p
			} else {
				cfg.Profiles[profileName] = config.Profile{
					APIURL:       apiURL,
					JWTToken:     resp.AccessToken,
					RefreshToken: resp.RefreshToken,
				}
			}
			if saveErr := config.Save(cfg); saveErr != nil {
				return fmt.Errorf("saving config: %w", saveErr)
			}
			fmt.Fprintln(os.Stderr, "JWT token saved to profile.")
		}

		if cmd.Flag("print-token").Changed {
			return printer().Print(map[string]interface{}{
				"access_token": resp.AccessToken,
				"token_type":   resp.TokenType,
				"expires_in":   resp.ExpiresIn,
			})
		}

		fmt.Fprintf(os.Stderr, "Logged in successfully. Expires in %d seconds.\n", resp.ExpiresIn)
		if printer().Format() == output.FormatJSON {
			return printer().Print(resp)
		}
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out (blacklist current JWT)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		_, err = c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/logout",
			Body:   map[string]interface{}{},
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, "Logged out successfully.")
		return nil
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the saved JWT access token",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg == nil || resolvedCfg.ProfileName == "" {
			return fmt.Errorf("no active profile")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		profile := cfg.Profiles[resolvedCfg.ProfileName]
		refreshToken, _ := cmd.Flags().GetString("refresh-token")
		if refreshToken == "" {
			refreshToken = profile.RefreshToken
		}
		if refreshToken == "" {
			return fmt.Errorf("no refresh token available; log in again with --save")
		}
		c := client.New(resolvedCfg.APIURL, "", "")
		data, err := c.Do(client.RequestOpts{Method: "POST", Path: "/v1/auth/refresh", Body: api.RefreshTokenRequest{RefreshToken: refreshToken}})
		if err != nil {
			return err
		}
		var resp api.RefreshTokenResponse
		if err := decodeData(data, &resp); err != nil {
			return err
		}
		profile.JWTToken = resp.AccessToken
		if resp.RefreshToken != "" {
			profile.RefreshToken = resp.RefreshToken
		}
		cfg.Profiles[resolvedCfg.ProfileName] = profile
		if err := config.Save(cfg); err != nil {
			return err
		}
		return printer().Print(resp)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")
		companyName, _ := cmd.Flags().GetString("company-name")
		phone, _ := cmd.Flags().GetString("phone")

		c := client.New(config.DefaultAPIURL, "", "")
		req := api.RegisterRequest{Email: email, Password: password, FirstName: firstName, LastName: lastName, CompanyName: companyName, Phone: phone}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/register",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var verifyEmailCmd = &cobra.Command{
	Use:   "verify-email",
	Short: "Verify email with OTP",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		otp, _ := cmd.Flags().GetString("otp")

		c := client.New(config.DefaultAPIURL, "", "")
		req := api.VerifyEmailRequest{Email: email, OTP: otp}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/verify-email",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var sendVerificationEmailCmd = &cobra.Command{
	Use:   "send-verification-email",
	Short: "Send verification email",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")

		c := client.New(config.DefaultAPIURL, "", "")
		req := api.SendVerificationEmailRequest{Email: email}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/send-verification-email",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Request password reset email",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")

		c := client.New(config.DefaultAPIURL, "", "")
		req := api.ResetPasswordRequest{Email: email}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/reset-password",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var resetPasswordConfirmCmd = &cobra.Command{
	Use:   "reset-password-confirm",
	Short: "Confirm password reset with OTP",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		otp, _ := cmd.Flags().GetString("otp")
		password, _ := cmd.Flags().GetString("password")

		c := client.New(config.DefaultAPIURL, "", "")
		req := api.ResetPasswordConfirmRequest{
			Email:    email,
			OTP:      otp,
			Password: password,
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/reset-password-confirm",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Get or update current user profile",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		if cmd.Flag("update").Changed {
			name, _ := cmd.Flags().GetString("name")
			req := api.UpdateProfileRequest{Name: name}
			data, err := c.Do(client.RequestOpts{
				Method: "PUT",
				Path:   "/v1/auth/me",
				Body:   req,
			})
			if err != nil {
				return err
			}
			return printer().Print(rawMessage(data))
		}

		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/auth/me"})
		if err != nil {
			return err
		}

		var result api.CurrentUser
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var changePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Change password (requires JWT, not API key)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("change-password requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		current, _ := cmd.Flags().GetString("current-password")
		newPass, _ := cmd.Flags().GetString("new-password")
		req := api.ChangePasswordRequest{
			CurrentPassword: current,
			NewPassword:     newPass,
		}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/change-password",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var verifyPhoneCmd = &cobra.Command{
	Use:   "verify-phone",
	Short: "Verify phone number with OTP (requires JWT)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("verify-phone requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		phone, _ := cmd.Flags().GetString("phone")
		otp, _ := cmd.Flags().GetString("otp")
		req := api.VerifyPhoneRequest{Phone: phone, OTP: otp}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/verify-phone",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

var sendPhoneOTPCmd = &cobra.Command{
	Use:   "send-phone-otp",
	Short: "Send OTP to phone for verification (requires JWT)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("send-phone-otp requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}

		phone, _ := cmd.Flags().GetString("phone")
		req := api.SendPhoneOTPRequest{Phone: phone}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/send-phone-otp",
			Body:   req,
		})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

// API Key management commands

var apiKeysCmd = &cobra.Command{
	Use:   "api-keys",
	Short: "Manage API keys (requires JWT)",
}

var apiKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("api-key management requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/auth/api-keys"})
		if err != nil {
			return err
		}
		var result []api.APIKey
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var apiKeysCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("api-key management requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		req := api.CreateAPIKeyRequest{Name: args[0]}
		data, err := c.Do(client.RequestOpts{
			Method: "POST",
			Path:   "/v1/auth/api-keys",
			Body:   req,
		})
		if err != nil {
			return err
		}
		var result api.APIKey
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var apiKeysGetCmd = &cobra.Command{
	Use:   "get [keyId]",
	Short: "Get an API key by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("api-key management requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "GET", Path: "/v1/auth/api-keys/" + args[0]})
		if err != nil {
			return err
		}
		var result api.APIKey
		if err := decodeData(data, &result); err != nil {
			return err
		}
		return printer().Print(result)
	},
}

var apiKeysDeleteCmd = &cobra.Command{
	Use:   "delete [keyId]",
	Short: "Delete an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.APIKey != "" && resolvedCfg.JWTToken == "" {
			return fmt.Errorf("api-key management requires a JWT token, not an API key")
		}
		c, err := apiClient()
		if err != nil {
			return err
		}
		data, err := c.Do(client.RequestOpts{Method: "DELETE", Path: "/v1/auth/api-keys/" + args[0]})
		if err != nil {
			return err
		}
		return printer().Print(rawMessage(data))
	},
}

func init() {
	loginCmd.Flags().String("email", "", "account email")
	loginCmd.Flags().String("password", "", "account password")
	loginCmd.Flags().Bool("print-token", false, "print the JWT token to stdout")
	loginCmd.Flags().Bool("save", false, "save JWT to current profile")
	loginCmd.Flags().String("api-url", "", "API URL override")
	refreshCmd.Flags().String("refresh-token", "", "refresh token override")

	registerCmd.Flags().String("email", "", "account email")
	registerCmd.Flags().String("password", "", "account password")
	registerCmd.Flags().String("first-name", "", "first name")
	registerCmd.Flags().String("last-name", "", "last name")
	registerCmd.Flags().String("company-name", "", "company name")
	registerCmd.Flags().String("phone", "", "phone number in E.164 format")

	verifyEmailCmd.Flags().String("email", "", "account email")
	verifyEmailCmd.Flags().String("otp", "", "verification OTP")

	sendVerificationEmailCmd.Flags().String("email", "", "account email")

	resetPasswordCmd.Flags().String("email", "", "account email")

	resetPasswordConfirmCmd.Flags().String("email", "", "account email")
	resetPasswordConfirmCmd.Flags().String("otp", "", "reset OTP")
	resetPasswordConfirmCmd.Flags().String("password", "", "new password")

	meCmd.Flags().Bool("update", false, "update profile instead of reading")
	meCmd.Flags().String("name", "", "new name")

	changePasswordCmd.Flags().String("current-password", "", "current password")
	changePasswordCmd.Flags().String("new-password", "", "new password")

	verifyPhoneCmd.Flags().String("phone", "", "phone number in E.164 format")
	verifyPhoneCmd.Flags().String("otp", "", "phone verification OTP")

	sendPhoneOTPCmd.Flags().String("phone", "", "phone number in E.164 format")

	// Mark required flags
	_ = loginCmd.MarkFlagRequired("email")
	_ = loginCmd.MarkFlagRequired("password")
	_ = registerCmd.MarkFlagRequired("email")
	_ = registerCmd.MarkFlagRequired("password")
	_ = registerCmd.MarkFlagRequired("first-name")
	_ = registerCmd.MarkFlagRequired("last-name")
	_ = verifyEmailCmd.MarkFlagRequired("email")
	_ = verifyEmailCmd.MarkFlagRequired("otp")
	_ = sendVerificationEmailCmd.MarkFlagRequired("email")
	_ = resetPasswordCmd.MarkFlagRequired("email")
	_ = resetPasswordConfirmCmd.MarkFlagRequired("email")
	_ = resetPasswordConfirmCmd.MarkFlagRequired("otp")
	_ = resetPasswordConfirmCmd.MarkFlagRequired("password")
	_ = changePasswordCmd.MarkFlagRequired("current-password")
	_ = changePasswordCmd.MarkFlagRequired("new-password")
	_ = verifyPhoneCmd.MarkFlagRequired("phone")
	_ = verifyPhoneCmd.MarkFlagRequired("otp")
	_ = sendPhoneOTPCmd.MarkFlagRequired("phone")

	apiKeysCmd.AddCommand(apiKeysListCmd)
	apiKeysCmd.AddCommand(apiKeysCreateCmd)
	apiKeysCmd.AddCommand(apiKeysGetCmd)
	apiKeysCmd.AddCommand(apiKeysDeleteCmd)
}
