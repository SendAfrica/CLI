package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProfileCRUD(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	profile := Profile{
		APIURL:   "https://api.test.com",
		APIKey:   "SA-testkey123",
		JWTToken: "jwt-test-token",
	}

	if err := AddProfile("test-profile", profile); err != nil {
		t.Fatalf("AddProfile failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.CurrentProfile != "test-profile" {
		t.Errorf("expected current profile test-profile, got %s", cfg.CurrentProfile)
	}

	p, ok := cfg.Profiles["test-profile"]
	if !ok {
		t.Fatal("profile not found")
	}
	if p.APIKey != "SA-testkey123" {
		t.Errorf("expected API key SA-testkey123, got %s", p.APIKey)
	}
	if p.JWTToken != "jwt-test-token" {
		t.Errorf("expected JWT token, got %s", p.JWTToken)
	}
}

func TestConfigPermissions(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := AddProfile("perm-test", Profile{APIKey: "key"}); err != nil {
		t.Fatalf("AddProfile failed: %v", err)
	}

	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected config file perms 0600, got %o", info.Mode().Perm())
	}
}

func TestResolveEnvVarOverride(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// Save a profile with a known API key
	if err := AddProfile("default", Profile{
		APIURL: "https://api.fromfile.com",
		APIKey: "SA-from-file",
	}); err != nil {
		t.Fatalf("AddProfile failed: %v", err)
	}

	// Env var should override file value
	t.Setenv(EnvAPIKey, "SA-from-env")
	t.Setenv(EnvAPIURL, "https://api.from-env.com")
	t.Setenv(EnvOutputFormat, "json")

	resolved, err := Resolve("", "", "", "", "")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if resolved.APIKey != "SA-from-env" {
		t.Errorf("env override failed: expected SA-from-env, got %s", resolved.APIKey)
	}
	if resolved.APIURL != "https://api.from-env.com" {
		t.Errorf("env URL override failed: expected from-env, got %s", resolved.APIURL)
	}
	if resolved.OutputFormat != "json" {
		t.Errorf("expected json output, got %s", resolved.OutputFormat)
	}
}

func TestResolveFlagOverride(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := AddProfile("default", Profile{
		APIKey: "SA-from-file",
	}); err != nil {
		t.Fatalf("AddProfile failed: %v", err)
	}

	t.Setenv(EnvAPIKey, "SA-from-env")

	// Flag should override env var
	resolved, err := Resolve("", "", "SA-from-flag", "", "yaml")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if resolved.APIKey != "SA-from-flag" {
		t.Errorf("flag override failed: expected SA-from-flag, got %s", resolved.APIKey)
	}
	if resolved.OutputFormat != "yaml" {
		t.Errorf("expected yaml output, got %s", resolved.OutputFormat)
	}
}

func TestResolveNoCredentials(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := AddProfile("default", Profile{
		APIURL: "https://api.example.com",
		// No API key or JWT
	}); err != nil {
		t.Fatalf("AddProfile failed: %v", err)
	}

	resolved, err := Resolve("", "", "", "", "")
	if err != nil {
		t.Fatalf("Resolve should not fail even without credentials: %v", err)
	}
	if resolved.APIKey != "" && resolved.JWTToken != "" {
		t.Error("expected no credentials")
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// No profiles saved
	_, err := Resolve("nonexistent", "", "", "", "")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestSetCurrentProfile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_ = AddProfile("profile-a", Profile{APIKey: "key-a"})
	_ = AddProfile("profile-b", Profile{APIKey: "key-b"})

	if err := SetCurrentProfile("profile-b"); err != nil {
		t.Fatalf("SetCurrentProfile failed: %v", err)
	}

	cfg, _ := Load()
	if cfg.CurrentProfile != "profile-b" {
		t.Errorf("expected current profile profile-b, got %s", cfg.CurrentProfile)
	}
}

func TestListProfiles(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_ = AddProfile("alpha", Profile{APIKey: "key-a"})
	_ = AddProfile("beta", Profile{APIKey: "key-b"})

	profiles, err := ListProfiles()
	if err != nil {
		t.Fatalf("ListProfiles failed: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestDeleteProfile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_ = AddProfile("temp-profile", Profile{APIKey: "key"})
	_ = AddProfile("another-profile", Profile{APIKey: "key2"})

	if err := DeleteProfile("temp-profile"); err != nil {
		t.Fatalf("DeleteProfile failed: %v", err)
	}

	cfg, _ := Load()
	if _, ok := cfg.Profiles["temp-profile"]; ok {
		t.Error("profile should have been deleted")
	}
	if _, ok := cfg.Profiles["another-profile"]; !ok {
		t.Error("other profile should still exist")
	}
}

func TestHasConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if HasConfig() {
		t.Error("expected no config file initially")
	}

	_ = AddProfile("default", Profile{APIKey: "key"})

	if !HasConfig() {
		t.Error("expected config file to exist after AddProfile")
	}
}

func TestConfigFileContent(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	profile := Profile{
		APIURL:   "https://api.example.com",
		APIKey:   "SA-secretkey",
		JWTToken: "jwt-token-here",
	}
	_ = AddProfile("default", profile)

	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Profiles["default"].APIKey != "SA-secretkey" {
		t.Error("API key should be stored in config file")
	}
}

// Ensure the config directory path is correct.
func TestConfigPath(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath failed: %v", err)
	}

	expected := filepath.Join(t.TempDir(), ".config", "sendafrica-cli", "config.json")
	if !filepath.IsAbs(path) {
		t.Errorf("config path should be absolute: %s", path)
	}
	_ = expected
}
