# SendAfrica CLI

A command-line interface for the [SendAfrica API](https://github.com/cameltech/sendafrica-backend). Supports all non-admin endpoints — admin routes (`/v1/admin/*`) are intentionally excluded.

Built with [Cobra](https://github.com/spf13/cobra) for command parsing and [Viper](https://github.com/spf13/viper) for configuration binding.

## Security

- **Never loads `.env` files.** Credentials are sourced only from:
  1. The config file (`~/.config/sendafrica-cli/config.json`, `0600` perms)
  2. Environment variables (`SENDAFRICA_API_URL`, `SENDAFRICA_API_KEY`, `SENDAFRICA_JWT_TOKEN`, `SENDAFRICA_OUTPUT`, `SENDAFRICA_PROFILE`)
  3. CLI flags (`--api-key`, `--token`, `--api-url`, `--output`, `--profile`)
- **Config file** is written with `0600` permissions and the config directory with `0700`.
- **No admin endpoints** — the `sendafrica config show` command masks secrets.
- **API key takes priority** over JWT when both are set. JWT is used only if no API key is present.
- **Secrets never leak to the command line.** Profile credentials are read from the config file, not passed as visible flags.

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

# Option 3: Profiles (for managing multiple accounts)
sendafrica config add-profile work --api-key SA-...
sendafrica config use work
sendafrica campaigns list
```

## Authentication Model

The SendAfrica API supports two credential types. Different endpoints accept different credentials:

| Endpoint | API Key | JWT |
|---|---|---|
| `POST /v1/auth/login` | — | — |
| `POST /v1/auth/register` | — | — |
| `POST /v1/sms/` | Yes | Yes |
| `GET /v1/sms/logs` | Yes | Yes |
| `GET /v1/credits/balance` | Yes | Yes |
| `GET /v1/me` | Yes | Yes |
| `POST /v1/auth/logout` | — | Yes |
| `GET /v1/auth/api-keys` | — | Yes |
| `POST /v1/auth/change-password` | — | Yes |
| `POST /v1/contact-lists` | Yes | Yes |
| `GET /v1/campaigns` | Yes | Yes |
| `POST /v1/payments` | Yes | Yes |

If both API key and JWT are configured, the CLI sends the API key (as `X-API-Key`). If only a JWT is configured, it sends `Authorization: Bearer <token>`.

## Credential Resolution Order

Credentials are resolved in this priority (highest wins):

1. **CLI flags** — `--api-key`, `--token`, `--api-url`
2. **Environment variables** — `SENDAFRICA_API_KEY`, `SENDAFRICA_JWT_TOKEN`, `SENDAFRICA_API_URL`
3. **Config profile** — the currently active profile in `~/.config/sendafrica-cli/config.json`

This means you can provide credentials the most convenient way for each context without changing global state.

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

> **Tip:** Use `sendafrica login --save` to persist the JWT token to your current profile. Add `--print-token` to output the raw token (useful for scripting or piping to `SENDAFRICA_JWT_TOKEN`).

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

> **Tip:** The support chat command supports a conversational flow. When the API responds with `confirmation_required`, re-run the command with `--confirm` to proceed. You can also load a conversation history from a file using `--chat-file` (JSON array or NDJSON format).

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

## Paginated Endpoints

List endpoints that return many results support `--page` and `--per-page` flags. The default is page 1 with 50 items per page.

```bash
sendafrica sms-logs --page 2 --per-page 100
sendafrica contacts list-contacts my-list-id --page 1 --per-page 200
sendafrica campaigns recipients campaign-id --status delivered
```

## Architecture

```
cmd/sendafrica/main.go          — entry point, calls cmd.Execute()
internal/
├── cmd/                        — Cobra commands + helpers
│   ├── root.go                 — root command, flag binding, OnInitialize
│   ├── auth.go                 — login, logout, register, password, phone, API keys
│   ├── sms.go                  — sms-send, sms-bulk, sms-logs, credits
│   ├── contacts.go             — contact lists, contacts, Google sync, CSV import/export
│   ├── campaigns.go            — campaign CRUD, scheduling, recipients
│   ├── payments.go             — payment initiation, vouchers
│   ├── senderids.go            — sender ID requirements, creation, listing
│   ├── notifications.go        — notification CRUD, read/unread
│   ├── public.go               — health, packages, templates, rates (no auth)
│   ├── support.go              — support chat with conversational flow
│   ├── config_commands.go      — profile management: add/use/list/delete/show
│   └── helpers.go              — shared utilities (decode, mask, parse, pagination)
├── api/                        — request/response types (no logic)
├── client/                     — HTTP client, envelope parsing, multipart
├── config/                     — profile CRUD, credential resolution
└── output/                     — table/JSON/YAML formatters
```

### Key Design Decisions

- **No `.env` loading** — credentials come from config file, env vars, or flags only, preventing accidental secret leakage from stray `.env` files.
- **Config file is `0600`** — the config directory is created with `0700` permissions.
- **Secret masking** — `config show` masks all credentials. The `mask()` function preserves the first 4 and last 4 characters.
- **API key priority** — when both API key and JWT are set, the API key is sent. This is because the API key is more specific to the CLI use case and is the primary auth method for SMS sending.
- **Table output via struct tags** — types use a `table:"name"` struct tag to control column headers in table output. Fields tagged `table:"-"` are hidden from tables.
- **Response envelope unwrapping** — the API wraps all responses in `{success, data, error}`. The client unwraps the `data` field automatically.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for full development guidelines.

### Prerequisites

- Go 1.25+
- Git

### Build

```bash
go build -o sendafrica ./cmd/sendafrica
```

### Run from source

```bash
go run ./cmd/sendafrica [command] [flags]
```

### Testing

```bash
go test ./...
```

### Lint and Vet

```bash
go vet ./...
```

## Troubleshooting

### "no credentials found"

The CLI requires either an API key or JWT token. Check your setup:

```bash
sendafrica config show          # shows resolved config (secrets masked)
sendafrica config list          # lists available profiles
```

Ensure at least one of these is set:
- `export SENDAFRIA_API_KEY=SA-...`
- `export SENDAFRIA_JWT_TOKEN=eyJ...`
- `sendafrica config add-profile default --api-key SA-...`

### "configuration not loaded"

The config file at `~/.config/sendafrica-cli/config.json` could not be read. Check file permissions and ensure the directory exists with `0700`:

```bash
mkdir -p ~/.config/sendafrica-cli
chmod 700 ~/.config/sendafrica-cli
```

### API key endpoints rejecting JWT

Some endpoints (e.g., `POST /v1/sms/`) require an API key specifically. If you get an authentication error, ensure you have an API key configured:

```bash
sendafrica config add-profile work --api-key SA-...
sendafrica config use work
```

### Wrong API URL

By default the CLI targets `https://api.sendafrica.com`. To use a different endpoint (e.g., staging):

```bash
export SENDAFRIA_API_URL=https://staging-api.sendafrica.com
# or
sendafrica --api-url https://staging-api.sendafrica.com sms-send ...
```

## License

See [LICENSE](LICENSE) for details.