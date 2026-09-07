// Package config manages CLI configuration: profiles, credentials, and API URL.
//
// SECURITY: This package deliberately does NOT load .env files. Credentials
// are read only from:
//  1. The config file (written with 0600 permissions).
//  2. Environment variables (SENDAFRICA_API_URL, SENDAFRIA_API_KEY, etc.).
//  3. CLI flags (--api-key, --token, --url).
//
// The godotenv / viper-watched-env-file pattern is intentionally avoided so
// that a stray .env in the working directory can never leak secrets into the
// CLI's process or onto the command line.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	EnvAPIURL       = "SENDAFRICA_API_URL"
	EnvAPIKey       = "SENDAFRICA_API_KEY"
	EnvJWTToken     = "SENDAFRICA_JWT_TOKEN"
	EnvProfile      = "SENDAFRICA_PROFILE"
	EnvOutputFormat = "SENDAFRICA_OUTPUT"

	DefaultAPIURL = "https://api.sendafrica.com"
	DefaultOutput = "table"

	configDirName = "sendafrica-cli"
	configFile    = "config.json"
)

// Profile holds the credentials and endpoint for a single SendAfrica account.
type Profile struct {
	APIURL   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	JWTToken string `json:"jwt_token"`
}

// Config is the on-disk representation of all profiles.
type Config struct {
	CurrentProfile string             `json:"current_profile"`
	Profiles       map[string]Profile `json:"profiles"`
}

// ResolvedConfig is the final, ready-to-use configuration after merging
// config file, environment, and flags.
type ResolvedConfig struct {
	APIURL       string
	APIKey       string
	JWTToken     string
	OutputFormat string
	ProfileName  string
}

func defaultConfig() *Config {
	return &Config{
		CurrentProfile: "default",
		Profiles:       make(map[string]Profile),
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", configDirName, configFile), nil
}

func ensureConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", configDirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("cannot create config directory %s: %w", dir, err)
	}
	return dir, nil
}

// Load reads the config file from disk. Returns a default (empty) config if
// the file does not exist. NEVER reads .env files.
func Load() (*Config, error) {
	p, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	if cfg.CurrentProfile == "" {
		cfg.CurrentProfile = "default"
	}
	return &cfg, nil
}

// Save writes the config file with 0600 permissions.
func Save(cfg *Config) error {
	dir, err := ensureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	path := filepath.Join(dir, configFile)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	return nil
}

// Resolve merges config-file profile, environment variables, and explicit
// flags into a single ResolvedConfig. Environment variables and flags always
// take precedence over the config file.
//
// Parameters:
//   - profileName: explicit profile name from --profile flag (empty = use current)
//   - apiURL:      explicit URL from --url flag (empty = fall back to env/config)
//   - apiKey:      explicit API key from --api-key flag (empty = fall back to env/config)
//   - jwtToken:    explicit JWT from --token flag (empty = fall back to env/config)
//   - format:      output format from --output flag (empty = fall back to env/default)
func Resolve(profileName, apiURL, apiKey, jwtToken, format string) (*ResolvedConfig, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	if profileName == "" {
		profileName = cfg.CurrentProfile
	}
	if profileName == "" {
		profileName = "default"
	}

	profile, exists := cfg.Profiles[profileName]
	if !exists {
		// No profile on disk — don't error; fall through with empty values
		// so env vars and flags can still provide credentials. The caller
		// (apiClient/ensureAuth) will reject if neither is present.
		profile = Profile{}
	}

	// Start with profile values, then overlay env vars, then flags.
	resolved := &ResolvedConfig{
		APIURL:       profile.APIURL,
		APIKey:       profile.APIKey,
		JWTToken:     profile.JWTToken,
		OutputFormat: DefaultOutput,
		ProfileName:  profileName,
	}

	if envURL := os.Getenv(EnvAPIURL); envURL != "" {
		resolved.APIURL = envURL
	}
	if envKey := os.Getenv(EnvAPIKey); envKey != "" {
		resolved.APIKey = envKey
	}
	if envToken := os.Getenv(EnvJWTToken); envToken != "" {
		resolved.JWTToken = envToken
	}
	if envFmt := os.Getenv(EnvOutputFormat); envFmt != "" {
		resolved.OutputFormat = envFmt
	}

	if apiURL != "" {
		resolved.APIURL = apiURL
	}
	if apiKey != "" {
		resolved.APIKey = apiKey
	}
	if jwtToken != "" {
		resolved.JWTToken = jwtToken
	}
	if format != "" {
		resolved.OutputFormat = format
	}

	if resolved.APIURL == "" {
		resolved.APIURL = DefaultAPIURL
	}

	return resolved, nil
}

// AddProfile creates or updates a profile in the config file.
func AddProfile(name string, p Profile) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.Profiles[name] = p
	// If the current profile doesn't exist in the map, switch to this one.
	if _, ok := cfg.Profiles[cfg.CurrentProfile]; !ok || cfg.CurrentProfile == "" {
		cfg.CurrentProfile = name
	}
	return Save(cfg)
}

// SetCurrentProfile switches the active profile.
func SetCurrentProfile(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	if _, ok := cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	cfg.CurrentProfile = name
	return Save(cfg)
}

// ListProfiles returns all profile names.
func ListProfiles() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	return names, nil
}

// DeleteProfile removes a profile from the config file.
func DeleteProfile(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	delete(cfg.Profiles, name)
	if cfg.CurrentProfile == name {
		cfg.CurrentProfile = ""
	}
	return Save(cfg)
}

// HasConfig returns true if a config file exists on disk.
func HasConfig() bool {
	p, err := configPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}
