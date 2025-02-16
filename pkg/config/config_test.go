// File: pkg/config/config_test.go

package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/heinrichb/scrapey-cli/pkg/utils"
)

// Helper functions to easily create pointer values.
func ptrString(s string) *string    { return &s }
func ptrInt(i int) *int             { return &i }
func ptrFloat64(f float64) *float64 { return &f }
func ptrBool(b bool) *bool          { return &b }

// TestApplyDefaults tests the ApplyDefaults function to ensure that missing fields are set to default values.
func TestApplyDefaults(t *testing.T) {
	cases := []struct {
		desc     string
		setup    func(cfg *Config)
		validate func(t *testing.T, cfg *Config)
	}{
		{
			desc:  "All fields missing should be set to defaults",
			setup: func(cfg *Config) {},
			validate: func(t *testing.T, cfg *Config) {
				if cfg.URL.Base != "https://example.com" {
					t.Errorf("Expected URL.Base to be 'https://example.com', got '%s'", cfg.URL.Base)
				}
				if len(cfg.URL.Routes) != 1 || cfg.URL.Routes[0] != "/" {
					t.Errorf("Expected URL.Routes to be ['/'], got %v", cfg.URL.Routes)
				}
				if cfg.ScrapingOptions.MaxDepth != 2 {
					t.Errorf("Expected ScrapingOptions.MaxDepth to be 2, got %d", cfg.ScrapingOptions.MaxDepth)
				}
				if cfg.ScrapingOptions.RateLimit != 1.5 {
					t.Errorf("Expected ScrapingOptions.RateLimit to be 1.5, got %f", cfg.ScrapingOptions.RateLimit)
				}
				if cfg.ScrapingOptions.RetryAttempts != 3 {
					t.Errorf("Expected ScrapingOptions.RetryAttempts to be 3, got %d", cfg.ScrapingOptions.RetryAttempts)
				}
				expectedUA := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
				if cfg.ScrapingOptions.UserAgent != expectedUA {
					t.Errorf("Expected ScrapingOptions.UserAgent to be '%s', got '%s'", expectedUA, cfg.ScrapingOptions.UserAgent)
				}
				if len(cfg.Storage.OutputFormats) != 1 || cfg.Storage.OutputFormats[0] != "json" {
					t.Errorf("Expected Storage.OutputFormats to be ['json'], got %v", cfg.Storage.OutputFormats)
				}
				if cfg.Storage.SavePath != "output/" {
					t.Errorf("Expected Storage.SavePath to be 'output/', got '%s'", cfg.Storage.SavePath)
				}
				if cfg.Storage.FileName != "scraped_data" {
					t.Errorf("Expected Storage.FileName to be 'scraped_data', got '%s'", cfg.Storage.FileName)
				}
			},
		},
		{
			desc: "Pre-set fields remain unchanged and missing fields get defaults",
			setup: func(cfg *Config) {
				cfg.URL.Base = "https://preset.com"
				cfg.Storage.SavePath = "custom_output/"
			},
			validate: func(t *testing.T, cfg *Config) {
				if cfg.URL.Base != "https://preset.com" {
					t.Errorf("Expected URL.Base to be 'https://preset.com', got '%s'", cfg.URL.Base)
				}
				if cfg.Storage.SavePath != "custom_output/" {
					t.Errorf("Expected Storage.SavePath to be 'custom_output/', got '%s'", cfg.Storage.SavePath)
				}
				if len(cfg.URL.Routes) != 1 || cfg.URL.Routes[0] != "/" {
					t.Errorf("Expected URL.Routes to be ['/'], got %v", cfg.URL.Routes)
				}
				if cfg.ScrapingOptions.MaxDepth != 2 {
					t.Errorf("Expected ScrapingOptions.MaxDepth to be 2, got %d", cfg.ScrapingOptions.MaxDepth)
				}
				if len(cfg.Storage.OutputFormats) != 1 || cfg.Storage.OutputFormats[0] != "json" {
					t.Errorf("Expected Storage.OutputFormats to be ['json'], got %v", cfg.Storage.OutputFormats)
				}
				if cfg.Storage.FileName != "scraped_data" {
					t.Errorf("Expected Storage.FileName to be 'scraped_data', got '%s'", cfg.Storage.FileName)
				}
			},
		},
		{
			desc: "No change if all fields are pre-set",
			setup: func(cfg *Config) {
				cfg.URL.Base = "https://preset.com"
				cfg.URL.Routes = []string{"/preset"}
				cfg.ScrapingOptions.MaxDepth = 10
				cfg.ScrapingOptions.RateLimit = 3.0
				cfg.ScrapingOptions.RetryAttempts = 5
				cfg.ScrapingOptions.UserAgent = "CustomAgent"
				cfg.Storage.OutputFormats = []string{"xml"}
				cfg.Storage.SavePath = "preset_output/"
				cfg.Storage.FileName = "preset_data"
			},
			validate: func(t *testing.T, cfg *Config) {
				if cfg.URL.Base != "https://preset.com" {
					t.Errorf("Expected URL.Base to be 'https://preset.com', got '%s'", cfg.URL.Base)
				}
				if !reflect.DeepEqual(cfg.URL.Routes, []string{"/preset"}) {
					t.Errorf("Expected URL.Routes to be ['/preset'], got %v", cfg.URL.Routes)
				}
				if cfg.ScrapingOptions.MaxDepth != 10 {
					t.Errorf("Expected ScrapingOptions.MaxDepth to be 10, got %d", cfg.ScrapingOptions.MaxDepth)
				}
				if cfg.ScrapingOptions.RateLimit != 3.0 {
					t.Errorf("Expected ScrapingOptions.RateLimit to be 3.0, got %f", cfg.ScrapingOptions.RateLimit)
				}
				if cfg.ScrapingOptions.RetryAttempts != 5 {
					t.Errorf("Expected ScrapingOptions.RetryAttempts to be 5, got %d", cfg.ScrapingOptions.RetryAttempts)
				}
				if cfg.ScrapingOptions.UserAgent != "CustomAgent" {
					t.Errorf("Expected ScrapingOptions.UserAgent to be 'CustomAgent', got '%s'", cfg.ScrapingOptions.UserAgent)
				}
				if !reflect.DeepEqual(cfg.Storage.OutputFormats, []string{"xml"}) {
					t.Errorf("Expected Storage.OutputFormats to be ['xml'], got %v", cfg.Storage.OutputFormats)
				}
				if cfg.Storage.SavePath != "preset_output/" {
					t.Errorf("Expected Storage.SavePath to be 'preset_output/', got '%s'", cfg.Storage.SavePath)
				}
				if cfg.Storage.FileName != "preset_data" {
					t.Errorf("Expected Storage.FileName to be 'preset_data', got '%s'", cfg.Storage.FileName)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := &Config{}
			if tc.setup != nil {
				tc.setup(cfg)
			}
			cfg.ApplyDefaults()
			tc.validate(t, cfg)
		})
	}
}

// TestLoad tests the Load function with various file conditions.
func TestLoad(t *testing.T) {
	var capturedColored string
	patchColored := monkey.Patch(utils.PrintColored, func(a ...interface{}) {
		capturedColored += fmt.Sprint(a...)
	})
	defer patchColored.Unpatch()

	var capturedNonEmpty string
	patchNonEmpty := monkey.Patch(utils.PrintNonEmptyFields, func(prefix string, cfg interface{}) {
		capturedNonEmpty += "nonEmptyFieldsCalled"
	})
	defer patchNonEmpty.Unpatch()

	cases := []struct {
		desc        string
		fileSetup   func(fileName string)
		verbose     bool
		expectErr   bool
		checkOutput func(t *testing.T, colored, nonEmpty string)
	}{
		{
			desc:      "Missing config file",
			fileSetup: nil,
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				if colored != "" {
					t.Errorf("Expected no colored output for missing file, got: %s", colored)
				}
			},
		},
		{
			desc: "Unreadable config file",
			fileSetup: func(name string) {
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
			},
		},
		{
			desc: "Invalid JSON format",
			fileSetup: func(name string) {
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
			},
		},
		{
			desc: "Valid JSON without verbose mode",
			fileSetup: func(name string) {
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   false,
			expectErr: false,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
				if nonEmpty != "" {
					t.Errorf("Expected no non-empty output when verbose is false, got: %s", nonEmpty)
				}
			},
		},
		{
			desc: "Valid JSON with verbose mode",
			fileSetup: func(name string) {
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   true,
			expectErr: false,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
				if nonEmpty != "nonEmptyFieldsCalled" {
					t.Errorf("Expected non-empty output when verbose is true, got: %s", nonEmpty)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			capturedColored = ""
			patchNonEmpty.Unpatch()
			patchNonEmpty = monkey.Patch(utils.PrintNonEmptyFields, func(prefix string, cfg interface{}) {
				capturedNonEmpty += "nonEmptyFieldsCalled"
			})
			defer patchNonEmpty.Unpatch()
			Verbose = tc.verbose

			var fileName string
			if tc.fileSetup != nil {
				tmpFile, err := os.CreateTemp("", "config_*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				fileName = tmpFile.Name()
				tmpFile.Close()
				tc.fileSetup(fileName)
				os.Chmod(fileName, 0644)
				defer os.Remove(fileName)
			} else {
				fileName = "nonexistent_config.json"
			}

			if tc.desc == "Unreadable config file" {
				patchReadFile := monkey.Patch(os.ReadFile, func(name string) ([]byte, error) {
					return nil, fmt.Errorf("simulated read error")
				})
				defer patchReadFile.Unpatch()
			}

			cfg, err := Load(fileName)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
			}
			if cfg.URL.Base == "" {
				t.Errorf("Expected URL.Base to be set, got empty")
			}
			tc.checkOutput(t, capturedColored, capturedNonEmpty)
		})
	}
}

// TestOverrideConfigFull tests the new OverrideConfig function using the ConfigOverride type.
// It creates a base config, applies a full override and verifies that all fields have been updated accordingly.
func TestOverrideConfigFull(t *testing.T) {
	var captured string
	patchColored := monkey.Patch(utils.PrintColored, func(a ...interface{}) {
		captured += fmt.Sprint(a...)
	})
	defer patchColored.Unpatch()

	// Create a base config with default values.
	base := &Config{}
	base.ApplyDefaults()

	// Create an override with non-nil pointers for every field.
	overrides := ConfigOverride{
		Version: ptrString("v2.0"),
		URL: &struct {
			Base        *string   `json:"base"`
			Routes      *[]string `json:"routes"`
			IncludeBase *bool     `json:"includeBase"`
		}{
			Base:        ptrString("https://override.com"),
			Routes:      &[]string{"/new", "/extra"},
			IncludeBase: ptrBool(true),
		},
		ParseRules: &struct {
			Title           *string `json:"title,omitempty"`
			MetaDescription *string `json:"metaDescription,omitempty"`
			ArticleContent  *string `json:"articleContent,omitempty"`
			Author          *string `json:"author,omitempty"`
			DatePublished   *string `json:"datePublished,omitempty"`
		}{
			Title:           ptrString("New Title"),
			MetaDescription: ptrString("New Meta"),
			ArticleContent:  ptrString("New Content"),
			Author:          ptrString("New Author"),
			DatePublished:   ptrString("2022-01-01"),
		},
		Storage: &struct {
			OutputFormats *[]string `json:"outputFormats"`
			SavePath      *string   `json:"savePath"`
			FileName      *string   `json:"fileName"`
		}{
			OutputFormats: &[]string{"csv"},
			SavePath:      ptrString("new_output/"),
			FileName:      ptrString("new_data"),
		},
		ScrapingOptions: &struct {
			MaxDepth      *int     `json:"maxDepth"`
			RateLimit     *float64 `json:"rateLimit"`
			RetryAttempts *int     `json:"retryAttempts"`
			UserAgent     *string  `json:"userAgent"`
		}{
			MaxDepth:      ptrInt(5),
			RateLimit:     ptrFloat64(2.0),
			RetryAttempts: ptrInt(4),
			UserAgent:     ptrString("OverrideAgent"),
		},
		DataFormatting: &struct {
			CleanWhitespace *bool `json:"cleanWhitespace"`
			RemoveHTML      *bool `json:"removeHTML"`
		}{
			CleanWhitespace: ptrBool(true),
			RemoveHTML:      ptrBool(true),
		},
	}

	// Apply the override.
	base.OverrideConfig(overrides)

	// Verify that each field has been updated.
	if base.Version != "v2.0" {
		t.Errorf("Expected Version to be 'v2.0', got '%s'", base.Version)
	}
	if base.URL.Base != "https://override.com" {
		t.Errorf("Expected URL.Base to be 'https://override.com', got '%s'", base.URL.Base)
	}
	if !reflect.DeepEqual(base.URL.Routes, []string{"/new", "/extra"}) {
		t.Errorf("Expected URL.Routes to be ['/new', '/extra'], got %v", base.URL.Routes)
	}
	if !base.URL.IncludeBase {
		t.Errorf("Expected URL.IncludeBase to be true")
	}
	if base.ParseRules.Title != "New Title" {
		t.Errorf("Expected ParseRules.Title to be 'New Title', got '%s'", base.ParseRules.Title)
	}
	if base.ParseRules.MetaDescription != "New Meta" {
		t.Errorf("Expected ParseRules.MetaDescription to be 'New Meta', got '%s'", base.ParseRules.MetaDescription)
	}
	if base.ParseRules.ArticleContent != "New Content" {
		t.Errorf("Expected ParseRules.ArticleContent to be 'New Content', got '%s'", base.ParseRules.ArticleContent)
	}
	if base.ParseRules.Author != "New Author" {
		t.Errorf("Expected ParseRules.Author to be 'New Author', got '%s'", base.ParseRules.Author)
	}
	if base.ParseRules.DatePublished != "2022-01-01" {
		t.Errorf("Expected ParseRules.DatePublished to be '2022-01-01', got '%s'", base.ParseRules.DatePublished)
	}
	if !reflect.DeepEqual(base.Storage.OutputFormats, []string{"csv"}) {
		t.Errorf("Expected Storage.OutputFormats to be ['csv'], got %v", base.Storage.OutputFormats)
	}
	if base.Storage.SavePath != "new_output/" {
		t.Errorf("Expected Storage.SavePath to be 'new_output/', got '%s'", base.Storage.SavePath)
	}
	if base.Storage.FileName != "new_data" {
		t.Errorf("Expected Storage.FileName to be 'new_data', got '%s'", base.Storage.FileName)
	}
	if base.ScrapingOptions.MaxDepth != 5 {
		t.Errorf("Expected ScrapingOptions.MaxDepth to be 5, got %d", base.ScrapingOptions.MaxDepth)
	}
	if base.ScrapingOptions.RateLimit != 2.0 {
		t.Errorf("Expected ScrapingOptions.RateLimit to be 2.0, got %f", base.ScrapingOptions.RateLimit)
	}
	if base.ScrapingOptions.RetryAttempts != 4 {
		t.Errorf("Expected ScrapingOptions.RetryAttempts to be 4, got %d", base.ScrapingOptions.RetryAttempts)
	}
	if base.ScrapingOptions.UserAgent != "OverrideAgent" {
		t.Errorf("Expected ScrapingOptions.UserAgent to be 'OverrideAgent', got '%s'", base.ScrapingOptions.UserAgent)
	}
	if !base.DataFormatting.CleanWhitespace {
		t.Errorf("Expected DataFormatting.CleanWhitespace to be true")
	}
	if !base.DataFormatting.RemoveHTML {
		t.Errorf("Expected DataFormatting.RemoveHTML to be true")
	}

	// Optionally, you can verify that PrintColored was called for each overridden field.
	expectedSubstrs := []string{
		"Overriding Version: v2.0",
		"Overriding URL.Base: https://override.com",
		"Overriding URL.Routes: [",
		"Overriding URL.IncludeBase: true",
		"Overriding ParseRules.Title: New Title",
		"Overriding ParseRules.MetaDescription: New Meta",
		"Overriding ParseRules.ArticleContent: New Content",
		"Overriding ParseRules.Author: New Author",
		"Overriding ParseRules.DatePublished: 2022-01-01",
		"Overriding Storage.OutputFormats: [",
		"Overriding Storage.SavePath: new_output/",
		"Overriding Storage.FileName: new_data",
		"Overriding ScrapingOptions.MaxDepth: 5",
		"Overriding ScrapingOptions.RateLimit: 2",
		"Overriding ScrapingOptions.RetryAttempts: 4",
		"Overriding ScrapingOptions.UserAgent: OverrideAgent",
		"Overriding DataFormatting.CleanWhitespace: true",
		"Overriding DataFormatting.RemoveHTML: true",
	}
	for _, substr := range expectedSubstrs {
		if !strings.Contains(captured, substr) {
			t.Errorf("Expected output to contain '%s', got '%s'", substr, captured)
		}
	}
}

// TestOverrideConfigNil tests that passing a ConfigOverride with all nil values does not change the config.
func TestOverrideConfigNil(t *testing.T) {
	var captured string
	patchColored := monkey.Patch(utils.PrintColored, func(a ...interface{}) {
		captured += fmt.Sprint(a...)
	})
	defer patchColored.Unpatch()

	// Create a base config with default values.
	base := &Config{}
	base.ApplyDefaults()

	// Create an override with all nil pointers.
	overrides := ConfigOverride{}

	// Apply the override.
	base.OverrideConfig(overrides)

	// Verify that no fields have changed (i.e. remain equal to their defaults).
	defaultConfig := &Config{}
	defaultConfig.ApplyDefaults()

	if !reflect.DeepEqual(base, defaultConfig) {
		t.Errorf("Expected config to remain unchanged when overrides are nil. Got %+v, expected %+v", base, defaultConfig)
	}

	// Since nothing is overridden, captured output should be empty.
	if captured != "" {
		t.Errorf("Expected no output from PrintColored when no overrides are applied, got '%s'", captured)
	}
}
