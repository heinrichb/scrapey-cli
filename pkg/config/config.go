// File: pkg/config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/heinrichb/scrapey-cli/pkg/utils"
)

/*
Global Verbose flag.

This flag determines whether verbose output is enabled.
It is set in `main.go` and used throughout the application.
*/
var Verbose bool

/*
Config holds configuration data used by Scrapey CLI.

Fields:
  - URL: A struct containing the base URL and routes to scrape.
  - ParseRules: A struct containing parsing rules.
  - Storage: A struct defining how data is saved.
  - ScrapingOptions: Settings for crawling behavior.
  - DataFormatting: Options for cleaning extracted content.

Usage:

	The configuration is loaded from a JSON file to guide the crawler and parser.
*/
type Config struct {
	Version string `json:"version"`
	URL     struct {
		Base        string   `json:"base"`
		Routes      []string `json:"routes"`
		IncludeBase bool     `json:"includeBase"`
	} `json:"url"`
	ParseRules struct {
		Title           string `json:"title,omitempty"`
		MetaDescription string `json:"metaDescription,omitempty"`
		ArticleContent  string `json:"articleContent,omitempty"`
		Author          string `json:"author,omitempty"`
		DatePublished   string `json:"datePublished,omitempty"`
	} `json:"parseRules"`
	Storage struct {
		OutputFormats []string `json:"outputFormats"`
		SavePath      string   `json:"savePath"`
		FileName      string   `json:"fileName"`
	} `json:"storage"`
	ScrapingOptions struct {
		MaxDepth      int     `json:"maxDepth"`
		RateLimit     float64 `json:"rateLimit"`
		RetryAttempts int     `json:"retryAttempts"`
		UserAgent     string  `json:"userAgent"`
	} `json:"scrapingOptions"`
	DataFormatting struct {
		CleanWhitespace bool `json:"cleanWhitespace"`
		RemoveHTML      bool `json:"removeHTML"`
	} `json:"dataFormatting"`
}

/*
ApplyDefaults populates missing fields in the Config struct with default values.

Usage:

	cfg.ApplyDefaults()

Notes:
  - Ensures that a missing Base URL defaults to "https://example.com".
  - Sets default scraping and storage parameters.
  - Provides a sensible fallback for all configurable values.
*/
func (cfg *Config) ApplyDefaults() {
	if cfg.URL.Base == "" {
		cfg.URL.Base = "https://example.com"
	}
	if len(cfg.URL.Routes) == 0 {
		cfg.URL.Routes = []string{"/"}
	}
	if cfg.ScrapingOptions.MaxDepth == 0 {
		cfg.ScrapingOptions.MaxDepth = 2
	}
	if cfg.ScrapingOptions.RateLimit == 0 {
		cfg.ScrapingOptions.RateLimit = 1.5
	}
	if cfg.ScrapingOptions.RetryAttempts == 0 {
		cfg.ScrapingOptions.RetryAttempts = 3
	}
	if cfg.ScrapingOptions.UserAgent == "" {
		cfg.ScrapingOptions.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}
	if len(cfg.Storage.OutputFormats) == 0 {
		cfg.Storage.OutputFormats = []string{"json"}
	}
	if cfg.Storage.SavePath == "" {
		cfg.Storage.SavePath = "output/"
	}
	if cfg.Storage.FileName == "" {
		cfg.Storage.FileName = "scraped_data"
	}
}

/*
Load reads configuration data from the specified filePath.

Parameters:
  - filePath: The path to the JSON configuration file.

Returns:
  - A pointer to a Config struct containing the parsed configuration.
  - An error if the file does not exist, cannot be read, or if the JSON is invalid.

Usage:

	cfg, err := Load("configs/default.json")
	if err != nil {
	    // Handle error
	}
	// Use cfg to configure the application.
*/
func Load(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", filePath)
	}

	utils.PrintColored("Loaded config from: ", filePath, color.FgHiGreen)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
	}

	// Apply default values where necessary.
	cfg.ApplyDefaults()

	// **Verbose Mode: Print Non-Empty Fields**
	if Verbose {
		utils.PrintNonEmptyFields("", cfg)
	}

	return &cfg, nil
}

/*
OverrideConfig applies all values from the provided `overrides` struct to the existing configuration.

Parameters:
  - overrides: A partial Config struct containing the fields to override. All values provided—including zero values,
    empty slices, and non-struct fields (e.g., Version)—are applied exactly as given.

Usage:

	cfg.OverrideConfig(Config{
		Version: "v2.0",
		URL: struct {
			Base        string   `json:"base"`
			Routes      []string `json:"routes"`
			IncludeBase bool     `json:"includeBase"`
		}{
			Base: "https://example.org",
		},
		ScrapingOptions: struct {
			MaxDepth      int     `json:"maxDepth"`
			RateLimit     float64 `json:"rateLimit"`
			RetryAttempts int     `json:"retryAttempts"`
			UserAgent     string  `json:"userAgent"`
		}{
			MaxDepth: 5,
		},
	})

Notes:
  - All values provided in `overrides` are applied, regardless of whether they are non-zero.
  - Uses **reflection** to dynamically override values while maintaining type safety.
  - Both struct and non-struct fields are overridden.
*/
func (cfg *Config) OverrideConfig(overrides Config) {
	cfgValue := reflect.ValueOf(cfg).Elem()
	overridesValue := reflect.ValueOf(overrides)

	// Since every top-level field in Config is exported and settable,
	// we do not check for validity or settable status.
	for i := 0; i < overridesValue.NumField(); i++ {
		field := overridesValue.Type().Field(i)
		overrideField := overridesValue.Field(i)
		configField := cfgValue.FieldByName(field.Name)

		if overrideField.Kind() == reflect.Struct {
			// For struct fields, override every subfield.
			for j := 0; j < overrideField.NumField(); j++ {
				subField := overrideField.Type().Field(j)
				overrideSubField := overrideField.Field(j)
				configSubField := configField.FieldByName(subField.Name)
				utils.PrintColored(fmt.Sprintf("Overriding %s.%s: ", field.Name, subField.Name),
					fmt.Sprint(overrideSubField.Interface()), color.FgHiMagenta)
				configSubField.Set(overrideSubField)
			}
		} else {
			// For non-struct fields, override unconditionally.
			utils.PrintColored(fmt.Sprintf("Overriding %s: ", field.Name),
				fmt.Sprint(overrideField.Interface()), color.FgHiMagenta)
			configField.Set(overrideField)
		}
	}
}
