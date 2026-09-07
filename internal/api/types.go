// Package api defines the request/response types for the SendAfrica API.
package api

import "encoding/json"

// PaginatedResponse wraps list endpoints that return {items, page, per_page, total, total_pages}.
type PaginatedResponse struct {
	Items      json.RawMessage `json:"items"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
	Count      int             `json:"count,omitempty"`
}

// ExtractItems unmarshals the Items field into the provided slice.
func (p PaginatedResponse) ExtractItems(target interface{}) error {
	if len(p.Items) == 0 {
		return nil
	}
	return json.Unmarshal(p.Items, target)
}

// Auth types

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordConfirmRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type SendPhoneOTPRequest struct {
	Phone string `json:"phone"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type CurrentUser struct {
	ID         string `json:"id" table:"id"`
	Email      string `json:"email" table:"email"`
	Name       string `json:"name" table:"name"`
	Phone      string `json:"phone,omitempty" table:"phone"`
	IsVerified bool   `json:"is_verified" table:"is_verified"`
	IsAdmin    bool   `json:"is_admin,omitempty" table:"is_admin"`
	CreatedAt  string `json:"created_at" table:"created_at"`
}

type APIKey struct {
	ID        string `json:"id" table:"id"`
	Name      string `json:"name" table:"name"`
	Key       string `json:"key,omitempty" table:"key"`
	KeyPrefix string `json:"key_prefix,omitempty" table:"key_prefix"`
	CreatedAt string `json:"created_at" table:"created_at"`
	LastUsed  string `json:"last_used,omitempty" table:"last_used"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

// SMS types

type SMSSendRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}

type SMSResponse struct {
	ID          string `json:"id" table:"id"`
	Status      string `json:"status" table:"status"`
	To          string `json:"to_phone" table:"to"`
	From        string `json:"from_id" table:"from"`
	Message     string `json:"message,omitempty" table:"message"`
	Cost        int    `json:"credits_used,omitempty" table:"cost"`
	GatewayCode string `json:"gateway_code,omitempty" table:"gateway_code"`
	CreatedAt   string `json:"created_at" table:"created_at"`
}

type SMSBulkSendRequest struct {
	To       []string `json:"to"`
	Message  string   `json:"message"`
	From     string   `json:"from,omitempty"`
	Schedule string   `json:"schedule,omitempty"`
}

type SMSLog struct {
	ID           string `json:"id" table:"id"`
	Status       string `json:"status" table:"status"`
	To           string `json:"to_phone" table:"to"`
	From         string `json:"from_id" table:"from"`
	Message      string `json:"message" table:"message"`
	Cost         int    `json:"credits_used" table:"cost"`
	DeliveredAt  string `json:"delivered_at,omitempty" table:"delivered_at"`
	GatewayCode  string `json:"gateway_code,omitempty" table:"gateway_code"`
	CreatedAt    string `json:"created_at" table:"created_at"`
}

type SMSLogsResponse struct {
	Items []SMSLog `json:"items"`
	Count int      `json:"count"`
}

// Credits types

type CreditBalance struct {
	Balance   int    `json:"balance" table:"balance"`
	Currency  string `json:"currency,omitempty" table:"currency"`
	AccountID string `json:"account_id,omitempty" table:"account_id"`
}

type Transaction struct {
	ID          string `json:"id" table:"id"`
	Type        string `json:"type" table:"type"`
	Amount      int    `json:"amount" table:"amount"`
	Balance     int    `json:"balance_after" table:"balance"`
	Status      string `json:"status,omitempty" table:"status"`
	Description string `json:"description,omitempty" table:"description"`
	CreatedAt   string `json:"created_at" table:"created_at"`
}

type TransactionHistory struct {
	Items []Transaction `json:"items"`
	Count int           `json:"count"`
}

// Public types

type Package struct {
	ID          string `json:"id" table:"id"`
	Name        string `json:"name" table:"name"`
	Description string `json:"description,omitempty" table:"description"`
	Credits     int    `json:"credits" table:"credits"`
	Price       int    `json:"price" table:"price"`
	Currency    string `json:"currency" table:"currency"`
	IsActive    bool   `json:"is_active" table:"is_active"`
	SortOrder   int    `json:"sort_order,omitempty" table:"sort_order"`
}

type SMSTemplate struct {
	ID       string `json:"id" table:"id"`
	Name     string `json:"name" table:"name"`
	Content  string `json:"content" table:"content"`
	Category string `json:"category,omitempty" table:"category"`
}

type Rate struct {
	Country     string `json:"name" table:"country"`
	Iso2        string `json:"iso2" table:"iso2"`
	CallingCode string `json:"dial_code" table:"calling_code"`
	RateTZS     int    `json:"rate_tzs" table:"rate_tzs"`
}

// Contacts types

type ContactList struct {
	ID           string `json:"id" table:"id"`
	Name         string `json:"name" table:"name"`
	ContactCount int    `json:"contact_count,omitempty" table:"contact_count"`
	Source       string `json:"source,omitempty" table:"source"`
	CreatedAt    string `json:"created_at" table:"created_at"`
	UpdatedAt    string `json:"updated_at,omitempty" table:"updated_at"`
}

type CreateContactListRequest struct {
	Name string `json:"name"`
}

type UpdateContactListRequest struct {
	Name string `json:"name"`
}

type DuplicateCheckResponse struct {
	DuplicateCount    int `json:"duplicate_count" table:"duplicate_count"`
	ListContactCount  int `json:"list_contact_count" table:"list_contact_count"`
}

type Contact struct {
	ID        string         `json:"id" table:"id"`
	ListID    string         `json:"list_id,omitempty" table:"list_id"`
	Name      string         `json:"name,omitempty" table:"name"`
	Phones    []ContactPhone `json:"phones,omitempty" table:"-"`
	CreatedAt string         `json:"created_at" table:"created_at"`
	UpdatedAt string         `json:"updated_at,omitempty" table:"updated_at"`
}

type ContactPhone struct {
	ID        string `json:"id" table:"id"`
	Phone     string `json:"phone" table:"phone"`
	Label     string `json:"label,omitempty" table:"label"`
	IsPrimary bool   `json:"is_primary,omitempty" table:"is_primary"`
}

type CreateContactRequest struct {
	Name   string   `json:"name,omitempty"`
	Phones []string `json:"phones"`
}

type UpdateContactRequest struct {
	Name   string   `json:"name,omitempty"`
	Phones []string `json:"phones,omitempty"`
}

type AddContactPhoneRequest struct {
	Phone string `json:"phone"`
	Label string `json:"label,omitempty"`
}

type ContactImportResult struct {
	RowsImported int `json:"rows_imported" table:"rows_imported"`
	RowsSkipped  int `json:"rows_skipped,omitempty" table:"rows_skipped"`
}

type GoogleContactsStatus struct {
	Connected      bool   `json:"connected" table:"connected"`
	ConnectedAt    string `json:"connected_at,omitempty" table:"connected_at"`
	LastSyncedAt   string `json:"last_synced_at,omitempty" table:"last_synced_at"`
	LastSyncStatus string `json:"last_sync_status,omitempty" table:"last_sync_status"`
	LastSyncError  string `json:"last_sync_error,omitempty" table:"last_sync_error"`
}

// Campaigns types

type Campaign struct {
	ID          string `json:"id" table:"id"`
	Name        string `json:"name" table:"name"`
	Status      string `json:"status" table:"status"`
	Message     string `json:"message,omitempty" table:"message"`
	From        string `json:"from,omitempty" table:"from"`
	ScheduledAt string `json:"scheduled_at,omitempty" table:"scheduled_at"`
	Estimate    struct {
		Cost     int    `json:"cost"`
		Count    int    `json:"count"`
		Currency string `json:"currency"`
	} `json:"estimate,omitempty" table:"-"`
	CreatedAt string `json:"created_at" table:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" table:"updated_at"`
}

type CreateCampaignRequest struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}

type UpdateCampaignRequest struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	From    string `json:"from,omitempty"`
}

type AttachContactListRequest struct {
	ListID string `json:"list_id"`
}

type CampaignRecipient struct {
	ID     string `json:"id" table:"id"`
	Phone  string `json:"phone" table:"phone"`
	Status string `json:"status" table:"status"`
	Error  string `json:"error,omitempty" table:"error"`
}

type CampaignRecipients struct {
	Items []CampaignRecipient `json:"items"`
	Count int                 `json:"count"`
}

// Payments types

type InitiatePaymentRequest struct {
	PackageID string `json:"package_id"`
}

type PaymentResponse struct {
	ID        string `json:"id" table:"id"`
	Status    string `json:"status" table:"status"`
	Amount    int    `json:"amount" table:"amount"`
	Currency  string `json:"currency" table:"currency"`
	Method    string `json:"method,omitempty" table:"method"`
	PackageID string `json:"package_id,omitempty" table:"package_id"`
	CreatedAt string `json:"created_at" table:"created_at"`
}

// Vouchers types

type VoucherRateResponse struct {
	MinAmountTZS int            `json:"min_amount_tzs" table:"min_amount_tzs"`
	Tiers        []VoucherTier `json:"tiers" table:"-"`
}

type VoucherTier struct {
	MaxAmountTZS      int `json:"max_amount_tzs" table:"max_amount_tzs"`
	RateTZSPerCredit  int `json:"rate_tzs_per_credit" table:"rate_tzs_per_credit"`
}

type PurchaseVoucherRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

type PurchaseVoucherResponse struct {
	ID        string `json:"id" table:"id"`
	Amount    int    `json:"amount" table:"amount"`
	Currency  string `json:"currency" table:"currency"`
	Credits   int    `json:"credits" table:"credits"`
	Status    string `json:"status" table:"status"`
	Method    string `json:"method,omitempty" table:"method"`
	CreatedAt string `json:"created_at" table:"created_at"`
}

// Sender IDs types

type SenderIDRequirements struct {
	Countries []SenderIDCountry  `json:"countries,omitempty"`
	Documents []SenderIDDocument `json:"documents,omitempty"`
	Purposes  []string           `json:"purposes,omitempty"`
}

type SenderIDCountry struct {
	UID         string `json:"uid" table:"uid"`
	Name        string `json:"name" table:"name"`
	Iso2        string `json:"iso2" table:"iso2"`
	CallingCode string `json:"calling_code" table:"calling_code"`
}

type SenderIDDocument struct {
	UID             string   `json:"uid" table:"uid"`
	Name            string   `json:"name" table:"name"`
	IsRequired      bool     `json:"is_required" table:"is_required"`
	AcceptedFormats []string `json:"accepted_formats,omitempty" table:"accepted_formats"`
	MaxSizeKB       int      `json:"max_size_kb,omitempty" table:"max_size_kb"`
}

type SenderID struct {
	ID             string `json:"id,omitempty" table:"id"`
	Name           string `json:"name" table:"name"`
	Country        string `json:"country,omitempty" table:"country"`
	Purpose        string `json:"purpose,omitempty" table:"purpose"`
	SampleMessage  string `json:"sample_message,omitempty" table:"sample_message"`
	Status         string `json:"status,omitempty" table:"status"`
	IsUsable       bool   `json:"is_usable,omitempty" table:"is_usable"`
	RejectionReason string `json:"rejection_reason,omitempty" table:"rejection_reason"`
	CreditsCharged int    `json:"credits_charged,omitempty" table:"credits_charged"`
	SubmittedAt    string `json:"submitted_at,omitempty" table:"submitted_at"`
	CreatedAt      string `json:"created_at,omitempty" table:"created_at"`
	UpdatedAt      string `json:"updated_at,omitempty" table:"updated_at"`
	// Fields from /usable endpoint
	IsDefault bool   `json:"is_default,omitempty" table:"is_default"`
	Type      string `json:"type,omitempty" table:"type"`
	Provider  string `json:"provider,omitempty" table:"provider"`
	Description string `json:"description,omitempty" table:"description"`
}

type CreateSenderIDRequest struct {
	Name          string                  `json:"name"`
	Country       string                  `json:"country"`
	Purpose       string                  `json:"purpose"`
	SampleMessage string                  `json:"sample_message"`
	Documents     []SenderIDDocumentInput `json:"documents,omitempty"`
}

type SenderIDDocumentInput struct {
	RequirementUID string `json:"requirement_uid"`
	Filename       string `json:"filename"`
	ContentBase64  string `json:"content_base64"`
}

// Notifications types

type Notification struct {
	ID        string `json:"id" table:"id"`
	Title     string `json:"title" table:"title"`
	Message   string `json:"message,omitempty" table:"message"`
	IsRead    bool   `json:"is_read" table:"is_read"`
	Type      string `json:"type,omitempty" table:"type"`
	CreatedAt string `json:"created_at" table:"created_at"`
}

// Support types

type SupportChatRequest struct {
	SessionID        string               `json:"session_id,omitempty"`
	Messages         []SupportChatMessage `json:"messages"`
	UserConfirmation bool                 `json:"user_confirmation"`
}

type SupportChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SupportChatResponse struct {
	SessionID            string  `json:"session_id" table:"session_id"`
	Status               string  `json:"status" table:"status"`
	Reply                string  `json:"reply" table:"reply"`
	ConfirmationRequired *string `json:"confirmation_required" table:"confirmation_required"`
	Iterations           int     `json:"iterations" table:"iterations"`
}
