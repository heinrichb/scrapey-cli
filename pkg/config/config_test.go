// File: pkg/config/config_test.go

package config

import (
	"os"
	"reflect"
	"testing"
)

/*
TestLoadValidConfig ensures that loading a valid config file works as expected.
*/
func TestLoadValidConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "valid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"url": {"base": "http://example.org", "routes": ["/test"], "includeBase": true}}`
	if _, err := tmpFile.Write([]byte(validJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected valid config, got error: %v", err)
	}

	expectedBaseURL := "http://example.org"
	expectedRoute := "/test"

	if cfg.URL.Base != expectedBaseURL {
		t.Errorf("Expected Base URL '%s', got '%s'", expectedBaseURL, cfg.URL.Base)
	}

	if !cfg.URL.IncludeBase {
		t.Errorf("Expected IncludeBase to be true, got false")
	}

	if len(cfg.URL.Routes) == 0 || cfg.URL.Routes[0] != expectedRoute {
		t.Errorf("Expected routes to include '%s', got %v", expectedRoute, cfg.URL.Routes)
	}
}

/*
TestLoadEmptyFile ensures that loading an empty config file does not cause an unexpected crash.
*/
func TestLoadEmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "empty_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Close() // Empty file

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected an error for empty config file, got nil")
	}
}

/*
TestLoadInvalidJSON ensures that a badly formatted JSON file correctly returns an error.
*/
func TestLoadInvalidJSON(t *testing.T) {
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

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}
}

/*
TestApplyDefaults verifies that ApplyDefaults correctly sets default values.
*/
func TestApplyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	expectedBaseURL := "https://example.com"
	expectedOutputFormat := "json"

	if cfg.URL.Base != expectedBaseURL {
		t.Errorf("Expected default Base URL '%s', got '%s'", expectedBaseURL, cfg.URL.Base)
	}
	if len(cfg.Storage.OutputFormats) == 0 || cfg.Storage.OutputFormats[0] != expectedOutputFormat {
		t.Errorf("Expected default output format '%s', got %v", expectedOutputFormat, cfg.Storage.OutputFormats)
	}
}

/*
TestOverrideWithCLI ensures that OverrideWithCLI dynamically updates config values.
*/
func TestOverrideWithCLI(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	newBaseURL := "https://cli-example.com"
	newRoutes := []string{"/cli-route1", "/cli-route2"}
	newMaxDepth := 10
	newUserAgent := "Custom CLI UserAgent"

	overrides := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Base:        newBaseURL,
			Routes:      newRoutes,
			IncludeBase: true,
		},
		ScrapingOptions: struct {
			MaxDepth      int     `json:"maxDepth"`
			RateLimit     float64 `json:"rateLimit"`
			RetryAttempts int     `json:"retryAttempts"`
			UserAgent     string  `json:"userAgent"`
		}{
			MaxDepth:  newMaxDepth,
			UserAgent: newUserAgent,
		},
	}

	cfg.OverrideWithCLI(overrides)

	if cfg.URL.Base != newBaseURL {
		t.Errorf("Expected Base URL to be overridden to '%s', got '%s'", newBaseURL, cfg.URL.Base)
	}

	if len(cfg.URL.Routes) != len(newRoutes) || cfg.URL.Routes[0] != newRoutes[0] || cfg.URL.Routes[1] != newRoutes[1] {
		t.Errorf("Expected Routes to be '%v', got '%v'", newRoutes, cfg.URL.Routes)
	}

	if !cfg.URL.IncludeBase {
		t.Errorf("Expected IncludeBase to be true, got false")
	}

	if cfg.ScrapingOptions.MaxDepth != newMaxDepth {
		t.Errorf("Expected MaxDepth to be overridden to %d, got '%d'", newMaxDepth, cfg.ScrapingOptions.MaxDepth)
	}

	if cfg.ScrapingOptions.UserAgent != newUserAgent {
		t.Errorf("Expected UserAgent to be '%s', got '%s'", newUserAgent, cfg.ScrapingOptions.UserAgent)
	}
}

/*
TestOverrideWithEmptyCLI ensures that OverrideWithCLI does nothing if no values are provided.
*/
func TestOverrideWithEmptyCLI(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	// Make a deep copy of the original config for comparison.
	originalConfig := *cfg

	cfg.OverrideWithCLI(Config{}) // Pass an empty config override

	// Use reflect.DeepEqual to compare struct contents.
	if !reflect.DeepEqual(*cfg, originalConfig) {
		t.Errorf("Expected config to remain unchanged when empty overrides are applied.\nExpected: %+v\nGot: %+v", originalConfig, *cfg)
	}
}
