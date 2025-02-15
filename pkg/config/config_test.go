// File: pkg/config/config_test.go

package config

import (
	"os"
	"testing"
)

/*
TestLoadValidConfig creates a temporary file with valid JSON configuration data,
writes a valid Base URL value, and attempts to load the configuration using Load.
It verifies that the returned Config object contains the expected Base URL.
*/
func TestLoadValidConfig(t *testing.T) {
	// Create a temporary file with valid JSON using os.CreateTemp.
	tmpFile, err := os.CreateTemp("", "valid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"url": {"base": "http://example.org", "routes": ["/test"], "includeBase": true}}`
	if _, err := tmpFile.Write([]byte(validJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close the file so it can be read by Load.

	// Load the configuration from the temporary file.
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	// Check that the Base URL in the configuration matches the expected value.
	if cfg.URL.Base != "http://example.org" {
		t.Errorf("Expected Base URL 'http://example.org', got '%s'", cfg.URL.Base)
	}

	// Check that the IncludeBase field is set correctly.
	if !cfg.URL.IncludeBase {
		t.Errorf("Expected IncludeBase to be true, got false")
	}

	// Ensure at least one route exists.
	if len(cfg.URL.Routes) == 0 || cfg.URL.Routes[0] != "/test" {
		t.Errorf("Expected routes to include '/test', got %v", cfg.URL.Routes)
	}
}

/*
TestLoadNonexistentFile attempts to load a configuration from a non-existent file
and verifies that Load returns an error.
*/
func TestLoadNonexistentFile(t *testing.T) {
	// Attempt to load a config from a non-existent file.
	_, err := Load("nonexistent_file.json")
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}

/*
TestLoadInvalidJSON creates a temporary file with invalid JSON content
(missing a closing brace) and attempts to load the configuration.
The test confirms that an error is returned due to invalid JSON.
*/
func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON using os.CreateTemp.
	tmpFile, err := os.CreateTemp("", "invalid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{"url": {"base": "http://example.org"`
	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Attempt to load the configuration from the temporary file.
	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}
}
