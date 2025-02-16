package config

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf strings.Builder
	io.Copy(&buf, r)
	os.Stdout = oldStdout
	return buf.String()
}

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		desc        string
		filename    string
		expectedErr bool
		setup       func(string)
	}{
		{
			"Missing config file",
			"nonexistent.json",
			true,
			nil,
		},
		{
			"Unreadable config file",
			"unreadable_config.json",
			true,
			func(name string) { os.Chmod(name, 0000); defer os.Chmod(name, 0644) },
		},
		{
			"Invalid JSON format",
			"invalid_config.json",
			true,
			func(name string) { os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"`), 0644) },
		},
		{
			"Valid JSON with verbose mode",
			"valid_config.json",
			false,
			func(name string) { os.WriteFile(name, []byte(`{"url": {"base": "http://example.org"}}`), 0644) },
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if c.setup != nil {
				tmpFile, _ := os.CreateTemp("", c.filename)
				defer os.Remove(tmpFile.Name())
				c.setup(tmpFile.Name())
				c.filename = tmpFile.Name()
			}

			_, err := Load(c.filename)
			if (err != nil) != c.expectedErr {
				t.Fatalf("Unexpected error state: %v", err)
			}
		})
	}
}

func TestOverrideWithCLI(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	cases := []struct {
		desc         string
		override     Config
		expectFunc   func(*Config) bool
		expectOutput string
	}{
		{
			"Override URL.Base",
			Config{URL: struct {
				Base        string   `json:"base"`
				Routes      []string `json:"routes"`
				IncludeBase bool     `json:"includeBase"`
			}{Base: "https://override.com"}},
			func(c *Config) bool { return c.URL.Base == "https://override.com" },
			"Overriding URL.Base: ",
		},
		{
			"Override non-empty slice",
			Config{Storage: struct {
				OutputFormats []string `json:"outputFormats"`
				SavePath      string   `json:"savePath"`
				FileName      string   `json:"fileName"`
			}{OutputFormats: []string{"csv"}}},
			func(c *Config) bool { return reflect.DeepEqual(c.Storage.OutputFormats, []string{"csv"}) },
			"Overriding Storage.OutputFormats: ",
		},
		{
			"Override boolean",
			Config{URL: struct {
				Base        string   `json:"base"`
				Routes      []string `json:"routes"`
				IncludeBase bool     `json:"includeBase"`
			}{IncludeBase: true}},
			func(c *Config) bool { return c.URL.IncludeBase },
			"Overriding URL.IncludeBase: ",
		},
		{
			"Override multiple values",
			Config{
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
			func(c *Config) bool { return c.URL.Base == "https://multiple.com" && c.ScrapingOptions.MaxDepth == 5 },
			"Overriding URL.Base: ",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			output := captureOutput(func() { cfg.OverrideWithCLI(c.override) })
			if !c.expectFunc(cfg) {
				t.Errorf("Expected override not applied")
			}
			if !strings.Contains(output, c.expectOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", c.expectOutput, output)
			}
		})
	}
}
