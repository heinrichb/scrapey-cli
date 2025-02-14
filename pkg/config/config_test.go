// File: pkg/config/config_test.go

package config

import (
	"os"
	"testing"
)

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

	// Attempt to load the config from the temporary file.
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	if cfg.URL != "http://example.org" {
		t.Errorf("Expected URL 'http://example.org', got '%s'", cfg.URL)
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	// Attempt to load a config from a non-existent file.
	_, err := Load("nonexistent_file.json")
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON using os.CreateTemp.
	tmpFile, err := os.CreateTemp("", "invalid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{"url": "http://example.org"` // Missing closing brace.
	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Attempt to load the config from the temporary file.
	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}
}
