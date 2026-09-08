# AGENTS.md

This file provides guidance for AI agents working in this repository.

## Build & Test Commands

```bash
go build ./...        # Build all packages
go vet ./...          # Run static analysis
go test ./...         # Run all tests
go test -v ./...      # Run tests with verbose output
go test -race ./...   # Run tests with race detection
```

## Project Structure

The SendAfrica CLI is a Go application using [Cobra](https://github.com/spf13/cobra) for commands and [Viper](https://github.com/spf13/viper) for flag binding.

```
cmd/sendafrica/main.go     — CLI entry point
internal/
├── cmd/                  — Cobra commands + shared helpers
├── api/                  — Request/response types (DTOs)
├── client/               — HTTP client, envelope parsing, multipart
├── config/               — Profile CRUD, credential resolution
└── output/               — Table/JSON/YAML formatters
```

## Key Patterns

- **Auth check**: `requireAuth()` → `apiClient()` before each API call.
- **Response decode**: `decodeData()` for single objects, `decodePaginated()` for lists with `{items: [...]}` envelope, `decodeListOrPaginated()` when the API may return either.
- **Output**: `printer().Print(result)` handles `--output table|json|yaml` automatically.
- **Secrets**: Never print raw credentials. Use `mask()` from `config_commands.go`. Config file is written with `0600` perms.
- **No `.env`**: The config package deliberately does NOT load environment files. Credentials come from env vars, config file, or flags only.

## Adding a New Command

1. Define a `*cobra.Command` in the relevant file under `internal/cmd/`.
2. Call `requireAuth()` + `apiClient()` in `RunE` (skip if public endpoint).
3. Make the API call: `c.Do(client.RequestOpts{Method, Path, Body, QueryParams})`.
4. Decode: `decodeData(data, &result)` or `decodePaginated(data, &result)`.
5. Print: `printer().Print(result)`.
6. Register flags in `init()`.
7. Wire up in `internal/cmd/root.go`.

## Testing

- **Client tests**: Use `httptest.NewServer` to mock API responses.
- **Config tests**: Use `t.Setenv("HOME", t.TempDir())` to isolate config file operations.
- **Output tests**: Use `Printer.WithWriter(&bytes.Buffer)` to capture output.

## Common Tasks

### Add a new API endpoint command
1. Add request/response types to `internal/api/types.go`.
2. Add the command to the appropriate file in `internal/cmd/`.
3. Register it in `internal/cmd/root.go`.
4. Update `README.md` with the new command.
5. Add tests if applicable.

### Update credential resolution
Modify `internal/config/config.go` `Resolve()` function. The override priority is: flags > env vars > config file.

### Change output format behavior
Modify `internal/output/formatter.go`. New formats should be added as `Format*` constants.
