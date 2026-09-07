# SendAfrica CLI

A command-line interface for the [SendAfrica API](https://github.com/cameltech/sendafrica-backend). Supports all non-admin endpoints — admin routes (`/v1/admin/*`) are intentionally excluded.

## Security

- **Never loads `.env` files.** Credentials are sourced only from:
  1. The config file (`~/.config/sendafrica-cli/config.json`, `0600` perms)
  2. Environment variables (`SENDAFRICA_API_URL`, `SENDAFRICA_API_KEY`, `SENDAFRICA_JWT_TOKEN`)
  3. CLI flags (`--api-key`, `--token`, `--api-url`)
- **Config file** is written with `0600` permissions and the config directory with `0700`.
- **No admin endpoints** — the `sendafrica config show` command masks secrets.

## Installation

```bash
go install github.com/cameltech/sendafrica-cli/cmd/sendafrica@latest
```

## Quick Start

```bash
# Option 1: API key from environment
export SENDAFRIA_API_KEY=SA-...
sendafrica sms-send --to "+255712345678" --message "Hello from CLI" -o json

# Option 2: JWT from environment
export SENDAFRIA_JWT_TOKEN=eyJ...
sendafrica credits-balance
sendafrica sms-logs

# Option 3: Profiles
sendafrica config add-profile work --api-key SA-...
sendafrica config use work
sendafrica campaigns list
```

## Commands

### Public (no auth required)

| Command | Endpoint |
|---|---|
| `sendafrica health` | `GET /health` |
| `sendafrica packages` | `GET /v1/packages` |
| `sendafrica templates` | `GET /v1/templates` |
| `sendafrica rates` | `GET /v1/rates` |
| `sendafrica rates --country TZ` | `GET /v1/rates/TZ` |

### Authentication

| Command | Endpoint | Auth |
|---|---|---|
| `sendafrica register --email ... --password ... --name ...` | `POST /v1/auth/register` | None |
| `sendafrica login --email ... --password ...` | `POST /v1/auth/login` | None |
| `sendafrica logout` | `POST /v1/auth/logout` | JWT |
| `sendafrica verify-email --email ... --otp ...` | `POST /v1/auth/verify-email` | None |
| `sendafrica send-verification-email --email ...` | `POST /v1/auth/send-verification-email` | None |
| `sendafrica reset-password --email ...` | `POST /v1/auth/reset-password` | None |
| `sendafrica reset-password-confirm --email ... --otp ... --password ...` | `POST /v1/auth/reset-password-confirm` | None |
| `sendafrica me` | `GET /v1/auth/me` | JWT/API Key |
| `sendafrica me --update --name ...` | `PUT /v1/auth/me` | JWT/API Key |
| `sendafrica change-password --current-password ... --new-password ...` | `POST /v1/auth/change-password` | JWT |
| `sendafrica verify-phone --phone ... --otp ...` | `POST /v1/auth/verify-phone` | JWT |
| `sendafrica send-phone-otp --phone ...` | `POST /v1/auth/send-phone-otp` | JWT |

### API Keys (JWT only)

| Command | Endpoint |
|---|---|
| `sendafrica api-keys list` | `GET /v1/auth/api-keys` |
| `sendafrica api-keys create [name]` | `POST /v1/auth/api-keys` |
| `sendafrica api-keys get [keyId]` | `GET /v1/auth/api-keys/{keyId}` |
| `sendafrica api-keys delete [keyId]` | `DELETE /v1/auth/api-keys/{keyId}` |

### SMS

| Command | Endpoint | Auth |
|---|---|---|
| `sendafrica sms-send --to ... --message ...` | `POST /v1/sms/` | API Key |
| `sendafrica sms-bulk --to ... --message ...` | `POST /v1/sms/bulk` | JWT/API Key |
| `sendafrica sms-logs` | `GET /v1/sms/logs` | JWT/API Key |

> **Note:** `POST /v1/sms/` requires API Key; use `sms-bulk` or an API key for `sms-send`.

### Credits

| Command | Endpoint |
|---|---|
| `sendafrica credits-balance` | `GET /v1/credits/balance` |
| `sendafrica credits-history` | `GET /v1/credits/history` |

### Contacts

| Command | Endpoint |
|---|---|
| `sendafrica contacts list` | `GET /v1/contact-lists` |
| `sendafrica contacts create [name]` | `POST /v1/contact-lists` |
| `sendafrica contacts update [listId] --name ...` | `PUT /v1/contact-lists/{listId}` |
| `sendafrica contacts delete [listId]` | `DELETE /v1/contact-lists/{listId}` |
| `sendafrica contacts duplicate-check [listId]` | `GET /v1/contact-lists/{listId}/duplicate-check` |
| `sendafrica contacts import [listId] [file.csv]` | `POST /v1/contact-lists/{listId}/import` |
| `sendafrica contacts list-contacts [listId]` | `GET /v1/contact-lists/{listId}/contacts` |
| `sendafrica contacts add-contact [listId] [phones]` | `POST /v1/contact-lists/{listId}/contacts` |
| `sendafrica contacts get-contact [listId] [contactId]` | `GET /v1/contact-lists/{listId}/contacts/{contactId}` |
| `sendafrica contacts update-contact [listId] [contactId]` | `PUT /v1/contact-lists/{listId}/contacts/{contactId}` |
| `sendafrica contacts delete-contact [listId] [contactId]` | `DELETE /v1/contact-lists/{listId}/contacts/{contactId}` |
| `sendafrica contacts add-phone [listId] [contactId] [phone]` | `POST /v1/contact-lists/{listId}/contacts/{contactId}/phones` |
| `sendafrica contacts delete-phone [listId] [contactId] [phoneId]` | `DELETE /v1/contact-lists/{listId}/contacts/{contactId}/phones/{phoneId}` |
| `sendafrica contacts export [listId]` | `GET /v1/contact-lists/{listId}/contacts/export` |
| `sendafrica contacts google-status` | `GET /v1/contact-lists/google/status` |
| `sendafrica contacts google-sync` | `POST /v1/contact-lists/google/sync` |
| `sendafrica contacts google-disconnect` | `DELETE /v1/contact-lists/google/disconnect` |

### Campaigns

| Command | Endpoint |
|---|---|
| `sendafrica campaigns list` | `GET /v1/campaigns` |
| `sendafrica campaigns create --name ... --message ...` | `POST /v1/campaigns` |
| `sendafrica campaigns get [id]` | `GET /v1/campaigns/{id}` |
| `sendafrica campaigns update [id] --name ... --message ...` | `PATCH /v1/campaigns/{id}` |
| `sendafrica campaigns cancel [id]` | `POST /v1/campaigns/{id}/cancel` |
| `sendafrica campaigns delete [id]` | `DELETE /v1/campaigns/{id}` |
| `sendafrica campaigns schedule [id] [listId]` | `POST /v1/campaigns/{id}/schedule` |
| `sendafrica campaigns recipients [id]` | `GET /v1/campaigns/{id}/recipients` |

### Payments

| Command | Endpoint |
|---|---|
| `sendafrica payments initiate [packageId]` | `POST /v1/payments/` |

### Vouchers

| Command | Endpoint |
|---|---|
| `sendafrica vouchers rate` | `GET /v1/vouchers/rate` |
| `sendafrica vouchers purchase [amount]` | `POST /v1/vouchers/` |

### Sender IDs

| Command | Endpoint |
|---|---|
| `sendafrica sender-ids requirements` | `GET /v1/sender-ids/requirements` |
| `sendafrica sender-ids usable` | `GET /v1/sender-ids/usable` |
| `sendafrica sender-ids list` | `GET /v1/sender-ids` |
| `sendafrica sender-ids create --name ... --purpose ... --sample-message ...` | `POST /v1/sender-ids` |
| `sendafrica sender-ids get [id]` | `GET /v1/sender-ids/{id}` |

### Notifications

| Command | Endpoint |
|---|---|
| `sendafrica notifications list` | `GET /v1/notifications` |
| `sendafrica notifications get [id]` | `GET /v1/notifications/{id}` |
| `sendafrica notifications unread-count` | `GET /v1/notifications/unread-count` |
| `sendafrica notifications mark-read [id]` | `PATCH /v1/notifications/{id}/read` |
| `sendafrica notifications read-all` | `POST /v1/notifications/read-all` |

### Support

| Command | Endpoint |
|---|---|
| `sendafrica support-chat --message "..."` | `POST /v1/support/chat` |

### Config Management

| Command | Description |
|---|---|
| `sendafrica config add-profile [name] --api-key ... --api-url ...` | Add or update a profile |
| `sendafrica config use [name]` | Switch active profile |
| `sendafrica config list` | List profiles |
| `sendafrica config delete [name]` | Delete a profile |
| `sendafrica config show` | Show resolved config (secrets masked) |

## Output Formats

All commands support `--output` (`-o`):

- `table` (default) — human-readable tables and key/value pairs
- `json` — JSON output for scripting
- `yaml` — YAML output

```bash
sendafrica credits-balance -o json
sendafrica sms-logs -o json --page 1 --per-page 100
```
