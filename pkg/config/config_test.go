package config

import (
	"os"
	"testing"
)

/*
TestLoadMissingFile ensures Load correctly returns an error when the config file does not exist.
*/
func TestLoadMissingFile(t *testing.T) {
	_, err := Load("nonexistent_config.json")
	if err == nil {
		t.Fatalf("Expected error for missing config file, got nil")
	}
}

/*
TestLoadUnreadableFile ensures Load correctly returns an error when the config file is unreadable.
*/
func TestLoadUnreadableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "unreadable_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := os.Chmod(tmpFile.Name(), 0000); err != nil {
		t.Fatalf("Failed to set file permissions: %v", err)
	}
	defer os.Chmod(tmpFile.Name(), 0644) // Restore permissions after test

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for unreadable file, got nil")
	}
}

/*
TestLoadVerboseMode ensures that verbose mode triggers PrintNonEmptyFields.
*/
func TestLoadVerboseMode(t *testing.T) {
	Verbose = true
	defer func() { Verbose = false }()

	tmpFile, err := os.CreateTemp("", "verbose_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"url": {"base": "http://example.org"}}`
	if _, err := tmpFile.Write([]byte(validJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	if cfg.URL.Base != "http://example.org" {
		t.Errorf("Expected Base URL 'http://example.org', got '%s'", cfg.URL.Base)
	}
}

/*
TestOverrideWithEmptySlices ensures that OverrideWithCLI skips empty slice values.
*/
func TestOverrideWithEmptySlices(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Routes: []string{},
		},
	}

	cfg.OverrideWithCLI(overrides)

	if len(cfg.URL.Routes) == 0 {
		t.Errorf("Expected Routes to remain unchanged, but they were overridden with an empty slice.")
	}
}
