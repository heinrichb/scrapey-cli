// File: pkg/config/config_test.go

package config

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

/*
captureOutput captures the output written to os.Stdout during the execution of f.
This helper is used to verify the PrintColored output without trying to override the function.
*/
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout
	return buf.String()
}

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

	output := captureOutput(func() {
		cfg.OverrideWithCLI(overrides)
	})
	if len(cfg.URL.Routes) == 0 {
		t.Errorf("Expected routes to remain unchanged, but they were overridden with an empty slice.")
	}
	if output != "" {
		t.Errorf("Expected no output when skipping empty slice, got %s", output)
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

/*
TestOverrideWithCLI_OverrideURLBase tests that a non-zero string override for URL.Base is applied,
and that the corresponding PrintColored output is produced.
*/
func TestOverrideWithCLI_OverrideURLBase(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	override := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Base: "https://override.com",
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.URL.Base != "https://override.com" {
		t.Errorf("Expected URL.Base to be overridden to https://override.com, got %s", cfg.URL.Base)
	}
	if !strings.Contains(output, "Overriding URL.Base: ") {
		t.Errorf("Expected output to contain 'Overriding URL.Base: ', got %s", output)
	}
}

/*
TestOverrideWithCLI_SkipZeroOverride tests that a zero value in the override does not change an existing value.
*/
func TestOverrideWithCLI_SkipZeroOverride(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	originalBase := cfg.URL.Base
	override := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Base: "", // zero value; should not override
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.URL.Base != originalBase {
		t.Errorf("Expected URL.Base to remain %s, got %s", originalBase, cfg.URL.Base)
	}
	if output != "" {
		t.Errorf("Expected no output for zero override, got %s", output)
	}
}

/*
TestOverrideWithCLI_OverrideNonEmptySlice tests that a non-empty slice override for Storage.OutputFormats is applied,
and that the expected PrintColored output is produced.
*/
func TestOverrideWithCLI_OverrideNonEmptySlice(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	override := Config{
		Storage: struct {
			OutputFormats []string `json:"outputFormats"`
			SavePath      string   `json:"savePath"`
			FileName      string   `json:"fileName"`
		}{
			OutputFormats: []string{"csv"},
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if !reflect.DeepEqual(cfg.Storage.OutputFormats, []string{"csv"}) {
		t.Errorf("Expected Storage.OutputFormats to be overridden to [csv], got %v", cfg.Storage.OutputFormats)
	}
	if !strings.Contains(output, "Overriding Storage.OutputFormats: ") {
		t.Errorf("Expected output to contain 'Overriding Storage.OutputFormats: ', got %s", output)
	}
}

/*
TestOverrideWithCLI_OverrideBoolean tests that a boolean override (for URL.IncludeBase) is applied when non-zero.
*/
func TestOverrideWithCLI_OverrideBoolean(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	override := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			IncludeBase: true,
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.URL.IncludeBase != true {
		t.Errorf("Expected URL.IncludeBase to be overridden to true, got %v", cfg.URL.IncludeBase)
	}
	if !strings.Contains(output, "Overriding URL.IncludeBase: ") {
		t.Errorf("Expected output to contain 'Overriding URL.IncludeBase: ', got %s", output)
	}
}

/*
TestOverrideWithCLI_OverrideMultiple tests that multiple fields across different sections are correctly overridden.
It also verifies that the output contains all expected override messages.
*/
func TestOverrideWithCLI_OverrideMultiple(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	override := Config{
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Base:        "https://multiple.com",
			Routes:      []string{"/new"},
			IncludeBase: true,
		},
		ScrapingOptions: struct {
			MaxDepth      int     `json:"maxDepth"`
			RateLimit     float64 `json:"rateLimit"`
			RetryAttempts int     `json:"retryAttempts"`
			UserAgent     string  `json:"userAgent"`
		}{
			MaxDepth:  5,
			UserAgent: "CustomAgent",
		},
		Storage: struct {
			OutputFormats []string `json:"outputFormats"`
			SavePath      string   `json:"savePath"`
			FileName      string   `json:"fileName"`
		}{
			SavePath: "custom_output/",
		},
		DataFormatting: struct {
			CleanWhitespace bool `json:"cleanWhitespace"`
			RemoveHTML      bool `json:"removeHTML"`
		}{
			CleanWhitespace: true,
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.URL.Base != "https://multiple.com" {
		t.Errorf("Expected URL.Base to be 'https://multiple.com', got %s", cfg.URL.Base)
	}
	if !reflect.DeepEqual(cfg.URL.Routes, []string{"/new"}) {
		t.Errorf("Expected URL.Routes to be ['/new'], got %v", cfg.URL.Routes)
	}
	if cfg.URL.IncludeBase != true {
		t.Errorf("Expected URL.IncludeBase to be true, got %v", cfg.URL.IncludeBase)
	}
	if cfg.ScrapingOptions.MaxDepth != 5 {
		t.Errorf("Expected ScrapingOptions.MaxDepth to be 5, got %d", cfg.ScrapingOptions.MaxDepth)
	}
	if cfg.ScrapingOptions.UserAgent != "CustomAgent" {
		t.Errorf("Expected ScrapingOptions.UserAgent to be 'CustomAgent', got %s", cfg.ScrapingOptions.UserAgent)
	}
	if cfg.Storage.SavePath != "custom_output/" {
		t.Errorf("Expected Storage.SavePath to be 'custom_output/', got %s", cfg.Storage.SavePath)
	}
	if cfg.DataFormatting.CleanWhitespace != true {
		t.Errorf("Expected DataFormatting.CleanWhitespace to be true, got %v", cfg.DataFormatting.CleanWhitespace)
	}
	expectedSubstrings := []string{
		"Overriding URL.Base: ",
		"Overriding URL.Routes: ",
		"Overriding URL.IncludeBase: ",
		"Overriding ScrapingOptions.MaxDepth: ",
		"Overriding ScrapingOptions.UserAgent: ",
		"Overriding Storage.SavePath: ",
		"Overriding DataFormatting.CleanWhitespace: ",
	}
	for _, substr := range expectedSubstrings {
		if !strings.Contains(output, substr) {
			t.Errorf("Expected output to contain '%s', got %s", substr, output)
		}
	}
}

/*
TestOverrideWithCLI_OverrideParseRules tests that non-zero overrides in the ParseRules struct are applied.
*/
func TestOverrideWithCLI_OverrideParseRules(t *testing.T) {
	cfg := &Config{}
	// Set initial values for ParseRules.
	cfg.ParseRules.Title = "Old Title"
	cfg.ParseRules.MetaDescription = "Old Desc"
	cfg.ParseRules.ArticleContent = "Old Content"
	cfg.ParseRules.Author = "Old Author"
	cfg.ParseRules.DatePublished = "Old Date"

	// Provide overrides.
	override := Config{
		ParseRules: struct {
			Title           string `json:"title,omitempty"`
			MetaDescription string `json:"metaDescription,omitempty"`
			ArticleContent  string `json:"articleContent,omitempty"`
			Author          string `json:"author,omitempty"`
			DatePublished   string `json:"datePublished,omitempty"`
		}{
			Title:           "New Title",
			MetaDescription: "New Desc",
			ArticleContent:  "New Content",
			Author:          "New Author",
			DatePublished:   "2025-02-15",
		},
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.ParseRules.Title != "New Title" {
		t.Errorf("Expected ParseRules.Title to be 'New Title', got '%s'", cfg.ParseRules.Title)
	}
	if cfg.ParseRules.MetaDescription != "New Desc" {
		t.Errorf("Expected ParseRules.MetaDescription to be 'New Desc', got '%s'", cfg.ParseRules.MetaDescription)
	}
	if cfg.ParseRules.ArticleContent != "New Content" {
		t.Errorf("Expected ParseRules.ArticleContent to be 'New Content', got '%s'", cfg.ParseRules.ArticleContent)
	}
	if cfg.ParseRules.Author != "New Author" {
		t.Errorf("Expected ParseRules.Author to be 'New Author', got '%s'", cfg.ParseRules.Author)
	}
	if cfg.ParseRules.DatePublished != "2025-02-15" {
		t.Errorf("Expected ParseRules.DatePublished to be '2025-02-15', got '%s'", cfg.ParseRules.DatePublished)
	}
	expectedSubstrings := []string{
		"Overriding ParseRules.Title: ",
		"Overriding ParseRules.MetaDescription: ",
		"Overriding ParseRules.ArticleContent: ",
		"Overriding ParseRules.Author: ",
		"Overriding ParseRules.DatePublished: ",
	}
	for _, substr := range expectedSubstrings {
		if !strings.Contains(output, substr) {
			t.Errorf("Expected output to contain '%s', got %s", substr, output)
		}
	}
}

/*
TestOverrideWithCLI_OverrideNonStruct tests that a non-struct field (TestDummy) is overridden.
This test exercises the else branch in OverrideWithCLI.
*/
func TestOverrideWithCLI_OverrideNonStruct(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()
	override := Config{
		TestDummy: 123,
	}
	output := captureOutput(func() {
		cfg.OverrideWithCLI(override)
	})
	if cfg.TestDummy != 123 {
		t.Errorf("Expected TestDummy to be overridden to 123, got %d", cfg.TestDummy)
	}
	if !strings.Contains(output, "Overriding TestDummy: ") {
		t.Errorf("Expected output to contain 'Overriding TestDummy: ', got %s", output)
	}
}
