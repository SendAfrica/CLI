# Contributing to SendAfrica CLI

Thank you for your interest in contributing to the SendAfrica CLI! This document covers everything you need to get started.

## Prerequisites

- **Go 1.25+** — the project uses modern Go features available in 1.25+
- **Git** — for version control and pushing changes
- A SendAfrica API account (for testing against the live API)

## Getting Started

```bash
# Clone your fork
git clone git@github.com:<your-username>/CLI.git
cd CLI

# Install dependencies (Go modules handle this automatically)
go mod download

# Verify the build
go build ./...

# Run tests
go test ./...
```

## Development Workflow

1. **Create a branch** for your feature or fix:

   ```bash
   git checkout -b feat/new-endpoint
   ```

2. **Make changes** following the conventions in this section.

3. **Run tests and checks** before committing:

   ```bash
   go test ./...
   go vet ./...
   ```

4. **Commit and push**, then open a pull request.

## Project Structure

```
cmd/sendafrica/main.go          — CLI entry point
internal/
├── cmd/                        — All Cobra command definitions
├── api/                        — Request/response types (DTOs only, no logic)
├── client/                     — HTTP client with envelope parsing
├── config/                     — Profile management and credential resolution
└── output/                     — Table/JSON/YAML output formatters
```

### Where to Add Things

- **New API endpoint** → Add a command in `internal/cmd/`, add request/response types in `internal/api/types.go`, wire it up in `internal/cmd/root.go`.
- **New output format** → Add to `internal/output/formatter.go`.
- **Config change** → Modify `internal/config/config.go`.
- **Shared command logic** → Add to `internal/cmd/helpers.go`.

## Coding Conventions

### Go Style

- Run `gofmt` before committing. The project does not enforce a linter beyond `go vet`.
- Use `go vet` as the primary static analysis tool.
- **Never load `.env` files.** The config package deliberately avoids `godotenv`. Credentials come from env vars, the config file (`0600`), or CLI flags.
- **Mask secrets.** The `mask()` helper in `config_commands.go` should be used wherever a credential is printed.
- **Table output** uses struct tags: `json:"..."` for API mapping and `table:"..."` for column headers. Use `table:"-"` to hide a field from table output.

### Command Structure

Each command file follows this pattern:

1. Define the `cobra.Command` with `Use`, `Short`, `Args`, and `RunE`.
2. In `RunE`, call `requireAuth()` first (for auth-required commands), then `apiClient()`.
3. Use `decodeData()` or `decodePaginated()` to parse the response.
4. Use `printer().Print()` to output the result (supports `--output table|json|yaml`).
5. Register flags and mark required flags in an `init()` function at the bottom of the file.

### Naming

- HTTP paths use the SendAfrica API path as-is (e.g., `/v1/sms/logs`).
- Command names use kebab-case (e.g., `sms-send`, `credits-balance`).
- Go variable names use camelCase (e.g., `smsSendCmd`, `creditsBalanceCmd`).

## Security Considerations

- **Never commit credentials.** The `.gitignore` excludes `.env` files.
- **Config file permissions.** The config file is written with `0600` and the directory with `0700`. Tests verify this.
- **No admin endpoints.** The CLI intentionally supports only non-admin API routes.

## Testing

Tests use Go's standard `testing` package (`testing` + `httptest`). No external test framework is required.

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run a specific package
go test ./internal/config/...

# Run with race detection
go test -race ./...

# Benchmark (where available)
go test -bench=. ./...
```

### Writing Tests

- **Client tests** — use `httptest.NewServer` to mock API responses. See `internal/client/client_test.go`.
- **Config tests** — use `t.Setenv("HOME", t.TempDir())` to isolate config file operations. See `internal/config/config_test.go`.
- **Output tests** — use `Printer.WithWriter(&bytes.Buffer)` to capture output. See `internal/output/formatter_test.go`.

## Adding a New Command

1. Create or edit a file in `internal/cmd/` (e.g., `internal/cmd/myfeature.go`).
2. Define the command variable (e.g., `var myFeatureCmd = &cobra.Command{...}`).
3. In `RunE`:
   - Call `requireAuth()` if auth is needed.
   - Get a client via `apiClient()`.
   - Make the API call via `c.Do(client.RequestOpts{...})`.
   - Decode and print the result.
4. Add flags and `MarkFlagRequired` calls in `init()`.
5. Register the command in `internal/cmd/root.go` (either directly on `rootCmd` or as a subcommand of a parent like `contactsCmd`).

Example minimal command:

```go
var myFeatureCmd = &cobra.Command{
    Use:   "my-feature",
    Short: "Does something useful",
    Args:  cobra.NoArgs,
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
            Path:   "/v1/my-feature",
        })
        if err != nil {
            return err
        }
        return printer().Print(rawMessage(data))
    },
}

func init() {
    // Register flags if needed
}
```

## Pull Requests

When opening a pull request:

1. Ensure all tests pass (`go test ./...`).
2. Ensure `go vet` is clean.
3. Add tests for new functionality.
4. Update the README if you've added or changed a command.
5. Write a clear PR description explaining **what** and **why**.

## Getting Help

- If the CLI behaves unexpectedly, run `sendafrica config show` (secrets are masked) to verify your setup.
- File issues on the [GitHub repository](https://github.com/SendAfrica/CLI).
