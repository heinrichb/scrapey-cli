// File: pkg/config/config_test.go

package config

import (
	"os"
	"reflect"
	"testing"
)

/*
TestLoadMissingFile ensures that loading a non-existent config file returns an error.
*/
func TestLoadMissingFile(t *testing.T) {
	_, err := Load("nonexistent.json")
	if err == nil {
		t.Fatalf("Expected error for missing config file, but got nil")
	}
}

/*
TestLoadUnreadableFile ensures that loading an unreadable file returns an error.
*/
func TestLoadUnreadableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "unreadable_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Make the file unreadable
	os.Chmod(tmpFile.Name(), 0000)
	defer os.Chmod(tmpFile.Name(), 0644) // Restore permissions after test

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for unreadable file, but got nil")
	}
}

/*
TestLoadInvalidJSONFormat ensures that loading a file with invalid JSON format returns an error.
*/
func TestLoadInvalidJSONFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{"url": {"base": "http://example.org"` // Missing closing brace
	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid JSON format, but got nil")
	}
}

/*
TestLoadVerboseMode ensures that verbose mode prints additional output.
*/
func TestLoadVerboseMode(t *testing.T) {
	Verbose = true
	defer func() { Verbose = false }()

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

	_, err = Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected successful load with verbose mode, but got error: %v", err)
	}
}

/*
TestOverrideWithInvalidField ensures that OverrideWithCLI correctly skips invalid fields.
*/
func TestOverrideWithInvalidField(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{}
	field := reflect.ValueOf(&overrides).Elem().FieldByName("InvalidField")

	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(42))
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Expected safe handling of invalid fields, but got panic: %v", r)
		}
	}()

	cfg.OverrideWithCLI(overrides)
}

/*
TestOverrideWithEmptySlices ensures that OverrideWithCLI skips overriding fields with empty slices.
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
		t.Errorf("Expected routes to remain unchanged, but they were overridden with an empty slice.")
	}
}

/*
TestOverrideWithValidFields ensures that OverrideWithCLI correctly applies valid overrides.
*/
func TestOverrideWithValidFields(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{
		ScrapingOptions: struct {
			MaxDepth      int     `json:"maxDepth"`
			RateLimit     float64 `json:"rateLimit"`
			RetryAttempts int     `json:"retryAttempts"`
			UserAgent     string  `json:"userAgent"`
		}{
			MaxDepth:  10,
			RateLimit: 3.5,
		},
	}

	cfg.OverrideWithCLI(overrides)

	if cfg.ScrapingOptions.MaxDepth != 10 {
		t.Errorf("Expected MaxDepth to be overridden to 10, got %d", cfg.ScrapingOptions.MaxDepth)
	}

	if cfg.ScrapingOptions.RateLimit != 3.5 {
		t.Errorf("Expected RateLimit to be overridden to 3.5, got %f", cfg.ScrapingOptions.RateLimit)
	}
}
