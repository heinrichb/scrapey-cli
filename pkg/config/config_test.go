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

// TestApplyDefaults tests the ApplyDefaults function to ensure that missing fields are set to default values.
// This test function uses multiple cases to verify that defaults are correctly applied.
func TestApplyDefaults(t *testing.T) {
	cases := []struct {
		desc     string
		setup    func(cfg *Config) // Optionally pre-set some fields.
		validate func(t *testing.T, cfg *Config)
	}{
		{
			desc: "All fields missing should be set to defaults",
			setup: func(cfg *Config) {
				// Leave all fields at their zero values.
			},
			validate: func(t *testing.T, cfg *Config) {
				// Check URL defaults.
				if cfg.URL.Base != "https://example.com" {
					t.Errorf("Expected URL.Base to be 'https://example.com', got '%s'", cfg.URL.Base)
				}
				if len(cfg.URL.Routes) != 1 || cfg.URL.Routes[0] != "/" {
					t.Errorf("Expected URL.Routes to be ['/'], got %v", cfg.URL.Routes)
				}
				// Check ScrapingOptions defaults.
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
				// Check Storage defaults.
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
				// Pre-set some fields.
				cfg.URL.Base = "https://preset.com"
				cfg.Storage.SavePath = "custom_output/"
			},
			validate: func(t *testing.T, cfg *Config) {
				// Pre-set values should be retained.
				if cfg.URL.Base != "https://preset.com" {
					t.Errorf("Expected URL.Base to be 'https://preset.com', got '%s'", cfg.URL.Base)
				}
				if cfg.Storage.SavePath != "custom_output/" {
					t.Errorf("Expected Storage.SavePath to be 'custom_output/', got '%s'", cfg.Storage.SavePath)
				}
				// Other fields should be set to defaults.
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
				// Set all fields explicitly.
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
				// Expect all pre-set fields to remain unchanged.
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
			// Allow test-specific pre-setup.
			if tc.setup != nil {
				tc.setup(cfg)
			}
			// Call ApplyDefaults.
			cfg.ApplyDefaults()
			// Validate that defaults have been applied as expected.
			tc.validate(t, cfg)
		})
	}
}

// TestLoad tests the Load function in a single function with multiple cases.
// We cover scenarios like missing file, unreadable file, invalid JSON,
// and valid JSON with verbose mode on/off.
func TestLoad(t *testing.T) {
	// Patch utils.PrintColored and utils.PrintNonEmptyFields.
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

	// Define table test cases for Load.
	cases := []struct {
		desc        string
		fileSetup   func(fileName string) // Setup the file (write contents, change permissions)
		verbose     bool                  // Set global Verbose before calling Load.
		expectErr   bool                  // Expect Load() to return an error.
		checkOutput func(t *testing.T, colored, nonEmpty string)
	}{
		{
			desc:      "Missing config file",
			fileSetup: nil, // Do not create the file so that it is missing.
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				// For a missing file, no printing should occur.
				if colored != "" {
					t.Errorf("Expected no colored output for missing file, got: %s", colored)
				}
			},
		},
		{
			desc: "Unreadable config file",
			fileSetup: func(name string) {
				// Create a file with valid JSON.
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
				// We'll patch os.ReadFile below to simulate a read error.
			},
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				// Expect that PrintColored is called.
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
			},
		},
		{
			desc: "Invalid JSON format",
			fileSetup: func(name string) {
				// Write invalid JSON.
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   false,
			expectErr: true,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				// Even with invalid JSON, colored output should be produced.
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
			},
		},
		{
			desc: "Valid JSON without verbose mode",
			fileSetup: func(name string) {
				// Write valid minimal JSON.
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   false,
			expectErr: false,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				// When verbose is false, only colored output is expected.
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
				// Write valid minimal JSON.
				if err := os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			},
			verbose:   true,
			expectErr: false,
			checkOutput: func(t *testing.T, colored, nonEmpty string) {
				// With verbose mode on, both colored and non-empty outputs should be present.
				if !strings.Contains(colored, "Loaded config from: ") {
					t.Errorf("Expected colored output, got: %s", colored)
				}
				if nonEmpty != "nonEmptyFieldsCalled" {
					t.Errorf("Expected non-empty output when verbose is true, got: %s", nonEmpty)
				}
			},
		},
	}

	// Run test cases.
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// Reset captured outputs.
			capturedColored = ""
			// Reset capturedNonEmpty inside the patch by re-patching.
			patchNonEmpty.Unpatch()
			patchNonEmpty = monkey.Patch(utils.PrintNonEmptyFields, func(prefix string, cfg interface{}) {
				capturedNonEmpty += "nonEmptyFieldsCalled"
			})
			defer patchNonEmpty.Unpatch()

			// Set the global Verbose flag as needed.
			Verbose = tc.verbose

			// Prepare file. If no setup, use a name that does not exist.
			var fileName string
			if tc.fileSetup != nil {
				tmpFile, err := os.CreateTemp("", "config_*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				fileName = tmpFile.Name()
				tmpFile.Close() // Close so that file can be manipulated.
				tc.fileSetup(fileName)
				// For cleanup and permission safety.
				os.Chmod(fileName, 0644)
				defer os.Remove(fileName)
			} else {
				fileName = "nonexistent_config.json"
			}

			// For the unreadable file test, patch os.ReadFile to simulate a read error.
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
				// Skip further checks if error was expected.
				return
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
			}

			// Ensure ApplyDefaults populated required fields (e.g. URL.Base).
			if cfg.URL.Base == "" {
				t.Errorf("Expected URL.Base to be set, got empty")
			}
			// Validate captured output.
			tc.checkOutput(t, capturedColored, capturedNonEmpty)
		})
	}
}

// TestOverrideConfig tests the OverrideConfig function in a single function with multiple cases.
// We patch utils.PrintColored to capture its output in a global variable.
// In addition to our previous override cases, we add tests for empty-slice overrides,
// for when no override is applied, and for non-struct fields (like the new Version field).
func TestOverrideConfig(t *testing.T) {
	// Patch utils.PrintColored to capture printed messages.
	var captured string
	patchColored := monkey.Patch(utils.PrintColored, func(a ...interface{}) {
		captured += fmt.Sprint(a...)
	})
	defer patchColored.Unpatch()

	// Define table test cases for OverrideConfig.
	cases := []struct {
		desc         string
		override     Config
		preSetup     func(*Config)      // Optionally modify the initial config.
		expectFunc   func(*Config) bool // Checks that the override was applied.
		expectOutput string             // Expected substring in the printed output.
	}{
		{
			desc: "Override URL.Base",
			override: Config{
				URL: struct {
					Base        string   `json:"base"`
					Routes      []string `json:"routes"`
					IncludeBase bool     `json:"includeBase"`
				}{Base: "https://override.com"},
			},
			preSetup: nil,
			expectFunc: func(c *Config) bool {
				return c.URL.Base == "https://override.com"
			},
			expectOutput: "Overriding URL.Base: ",
		},
		{
			desc: "Override non-empty slice",
			override: Config{
				Storage: struct {
					OutputFormats []string `json:"outputFormats"`
					SavePath      string   `json:"savePath"`
					FileName      string   `json:"fileName"`
				}{OutputFormats: []string{"csv"}},
			},
			preSetup: func(c *Config) {
				c.Storage.OutputFormats = []string{"json"}
			},
			expectFunc: func(c *Config) bool {
				return reflect.DeepEqual(c.Storage.OutputFormats, []string{"csv"})
			},
			expectOutput: "Overriding Storage.OutputFormats: ",
		},
		{
			desc: "Override boolean",
			override: Config{
				URL: struct {
					Base        string   `json:"base"`
					Routes      []string `json:"routes"`
					IncludeBase bool     `json:"includeBase"`
				}{IncludeBase: true},
			},
			preSetup:     nil,
			expectFunc:   func(c *Config) bool { return c.URL.IncludeBase },
			expectOutput: "Overriding URL.IncludeBase: ",
		},
		{
			desc: "Override multiple values",
			override: Config{
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
				}{MaxDepth: 5},
			},
			preSetup: nil,
			expectFunc: func(c *Config) bool {
				return c.URL.Base == "https://multiple.com" && c.ScrapingOptions.MaxDepth == 5
			},
			expectOutput: "Overriding URL.Base: ",
		},
		{
			desc: "Override empty slice (applies override)",
			override: Config{
				Storage: struct {
					OutputFormats []string `json:"outputFormats"`
					SavePath      string   `json:"savePath"`
					FileName      string   `json:"fileName"`
				}{OutputFormats: []string{}}, // Even empty slice should override.
			},
			preSetup: func(c *Config) {
				c.Storage.OutputFormats = []string{"json"}
			},
			expectFunc: func(c *Config) bool {
				// Expect the override to apply, resulting in an empty slice.
				return reflect.DeepEqual(c.Storage.OutputFormats, []string{})
			},
			expectOutput: "Overriding Storage.OutputFormats: ",
		},
		{
			desc: "Override non-struct field (Version)",
			override: Config{
				Version: "v2.0",
			},
			preSetup: func(c *Config) {
				c.Version = "v1.0"
			},
			expectFunc: func(c *Config) bool {
				// Expect the version to be overridden to "v2.0".
				return c.Version == "v2.0"
			},
			expectOutput: "Overriding Version: ",
		},
	}

	// Run test cases.
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			captured = "" // Reset captured output.
			// Create a fresh config with defaults applied.
			testCfg := &Config{}
			testCfg.ApplyDefaults()
			if tc.preSetup != nil {
				tc.preSetup(testCfg)
			}

			// Apply the override.
			testCfg.OverrideConfig(tc.override)
			// Verify that the override was applied.
			if !tc.expectFunc(testCfg) {
				t.Errorf("Expected override condition not met. Got %+v", testCfg)
			}
			// Verify that PrintColored was called with the expected message.
			if !strings.Contains(captured, tc.expectOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", tc.expectOutput, captured)
			}
		})
	}
}
