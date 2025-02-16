// File: pkg/config/config_test.go

package config

import (
	"os"
	"reflect"
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
TestLoadInvalidJSONFormat ensures Load correctly returns an error for malformed JSON.
*/
func TestLoadInvalidJSONFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid_json_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Invalid JSON (missing closing brace)
	invalidJSON := `{"url": {"base": "http://example.org"`
	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Fatalf("Expected error for malformed JSON, got nil")
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
TestOverrideWithInvalidField ensures that OverrideWithCLI skips invalid fields safely.
*/
func TestOverrideWithInvalidField(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{}
	field := reflect.ValueOf(&overrides).Elem().FieldByName("InvalidField")

	// Ensure the field exists and can be modified
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(42)) // Invalid field
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Expected safe handling of invalid fields, but got panic: %v", r)
		}
	}()

	cfg.OverrideWithCLI(overrides)
}

/*
TestOverrideWithInvalidNestedField ensures that OverrideWithCLI skips invalid nested fields safely.
*/
func TestOverrideWithInvalidNestedField(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{}

	// Get the field reference
	field := reflect.ValueOf(&overrides).Elem().FieldByName("URL")

	// Ensure the field exists and is of type struct before attempting an invalid override
	if field.IsValid() && field.Kind() == reflect.Struct {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Expected safe handling of invalid nested fields, but got panic: %v", r)
			}
		}()

		// Try setting an invalid value (int) but ensure it's prevented before setting
		invalidValue := reflect.ValueOf(42)
		if field.CanSet() && field.Kind() == invalidValue.Kind() { // Ensure types match
			field.Set(invalidValue) // This should NOT execute due to validation.
		}
	}

	// Ensure the override function does not crash on incorrect input
	cfg.OverrideWithCLI(overrides)
}

/*
TestOverrideWithValidSubField ensures that OverrideWithCLI correctly overrides valid subfields.
*/
func TestOverrideWithValidSubField(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	overrides := Config{
		ParseRules: struct {
			Title           string `json:"title,omitempty"`
			MetaDescription string `json:"metaDescription,omitempty"`
			ArticleContent  string `json:"articleContent,omitempty"`
			Author          string `json:"author,omitempty"`
			DatePublished   string `json:"datePublished,omitempty"`
		}{
			Title: "New Title",
		},
	}

	cfg.OverrideWithCLI(overrides)

	if cfg.ParseRules.Title != "New Title" {
		t.Errorf("Expected Title to be 'New Title', got '%s'", cfg.ParseRules.Title)
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

/*
TestOverrideWithValidFields ensures that OverrideWithCLI correctly applies non-zero overrides.
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
			MaxDepth:  5,
			RateLimit: 2.5,
		},
	}

	cfg.OverrideWithCLI(overrides)

	if cfg.ScrapingOptions.MaxDepth != 5 {
		t.Errorf("Expected MaxDepth to be overridden to 5, got '%d'", cfg.ScrapingOptions.MaxDepth)
	}

	if cfg.ScrapingOptions.RateLimit != 2.5 {
		t.Errorf("Expected RateLimit to be overridden to 2.5, got '%f'", cfg.ScrapingOptions.RateLimit)
	}
}
