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

// TestOverrideWithCLI tests the OverrideWithCLI function in a single function with multiple cases.
// We patch utils.PrintColored to capture its output in a global variable.
// In addition to our previous override cases, we add tests for empty-slice overrides
// and for when no override is applied.
func TestOverrideWithCLI(t *testing.T) {
	// Patch utils.PrintColored.
	var captured string
	patchColored := monkey.Patch(utils.PrintColored, func(a ...interface{}) {
		captured += fmt.Sprint(a...)
	})
	defer patchColored.Unpatch()

	// Define table test cases for OverrideWithCLI.
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
			expectFunc:   func(c *Config) bool { return c.URL.Base == "https://override.com" },
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
				// Pre-populate with a different slice.
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
			expectFunc: func(c *Config) bool {
				return c.URL.Base == "https://multiple.com" && c.ScrapingOptions.MaxDepth == 5
			},
			expectOutput: "Overriding URL.Base: ",
		},
		{
			desc: "Empty slice override does not trigger override",
			override: Config{
				Storage: struct {
					OutputFormats []string `json:"outputFormats"`
					SavePath      string   `json:"savePath"`
					FileName      string   `json:"fileName"`
				}{OutputFormats: []string{}}, // Empty slice; should be skipped.
			},
			preSetup: func(c *Config) {
				// Set a non-empty default to confirm it is not overridden.
				c.Storage.OutputFormats = []string{"json"}
			},
			expectFunc: func(c *Config) bool {
				// Expect no change.
				return reflect.DeepEqual(c.Storage.OutputFormats, []string{"json"})
			},
			expectOutput: "", // No override message expected.
		},
		{
			desc: "No override applied when all fields are zero",
			override: Config{
				URL: struct {
					Base        string   `json:"base"`
					Routes      []string `json:"routes"`
					IncludeBase bool     `json:"includeBase"`
				}{}, // all zero values
				Storage: struct {
					OutputFormats []string `json:"outputFormats"`
					SavePath      string   `json:"savePath"`
					FileName      string   `json:"fileName"`
				}{},
				ParseRules: struct {
					Title           string `json:"title,omitempty"`
					MetaDescription string `json:"metaDescription,omitempty"`
					ArticleContent  string `json:"articleContent,omitempty"`
					Author          string `json:"author,omitempty"`
					DatePublished   string `json:"datePublished,omitempty"`
				}{},
				ScrapingOptions: struct {
					MaxDepth      int     `json:"maxDepth"`
					RateLimit     float64 `json:"rateLimit"`
					RetryAttempts int     `json:"retryAttempts"`
					UserAgent     string  `json:"userAgent"`
				}{},
				DataFormatting: struct {
					CleanWhitespace bool `json:"cleanWhitespace"`
					RemoveHTML      bool `json:"removeHTML"`
				}{},
			},
			expectFunc: func(c *Config) bool {
				// Expect no changes: the defaults remain.
				return c.URL.Base != "" && len(c.Storage.OutputFormats) > 0
			},
			expectOutput: "", // No output expected.
		},
	}

	// Run test cases.
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			captured = "" // Reset captured output.
			// Create a fresh config with defaults applied.
			testCfg := &Config{}
			testCfg.ApplyDefaults()
			// If any pre-setup is needed, run it.
			if tc.preSetup != nil {
				tc.preSetup(testCfg)
			}

			// Directly call OverrideWithCLI.
			testCfg.OverrideWithCLI(tc.override)
			// Verify that the override was applied (or not applied) as expected.
			if !tc.expectFunc(testCfg) {
				t.Errorf("Expected override condition not met. Got %+v", testCfg)
			}
			// Verify that the patched PrintColored was called with the expected substring.
			if !strings.Contains(captured, tc.expectOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", tc.expectOutput, captured)
			}
		})
	}
}
