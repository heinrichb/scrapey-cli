// File: pkg/config/config_test.go
package config

import (
	"os"
	"testing"
)

// TestLoadValidConfig creates a temporary file containing valid JSON configuration data,
// writes a valid URL value, and then attempts to load the configuration using Load.
// The test verifies that the returned Config object contains the expected URL.
func TestLoadValidConfig(t *testing.T) {
	// Create a temporary file with valid JSON using os.CreateTemp.
	tmpFile, err := os.CreateTemp("", "valid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"url": "http://example.org"}`
	if _, err := tmpFile.Write([]byte(validJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close the file so it can be read by Load.

	// Load the configuration from the temporary file.
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	// Check that the URL in the configuration matches the expected value.
	if cfg.URL != "http://example.org" {
		t.Errorf("Expected URL 'http://example.org', got '%s'", cfg.URL)
	}
}

// TestLoadNonexistentFile attempts to load a configuration from a file path that does not exist,
// and verifies that Load returns an error.
func TestLoadNonexistentFile(t *testing.T) {
	// Attempt to load a config from a non-existent file.
	_, err := Load("nonexistent_file.json")
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}

// TestLoadInvalidJSON creates a temporary file with invalid JSON content,
// then attempts to load the configuration. The test confirms that an error is returned due to invalid JSON.
func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON using os.CreateTemp.
	tmpFile, err := os.CreateTemp("", "invalid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid JSON (missing closing brace) into the temporary file.
	invalidJSON := `{"url": "http://example.org"`
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
