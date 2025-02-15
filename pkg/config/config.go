// File: pkg/config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/heinrichb/scrapey-cli/pkg/utils"
)

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
	URL struct {
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

Notes:
  - This function uses os.ReadFile to read the file.
  - It prints a confirmation message in high-intensity green using the PrintColored utility.
  - It then calls PrintNonEmptyFields from the utils package to display non-empty config fields.
*/
func Load(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", filePath)
	}

	// Print confirmation that the config was loaded, using our PrintColored utility.
	utils.PrintColored("Loaded config from: ", filePath, color.FgHiGreen)

	// Read file contents using os.ReadFile.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal JSON into a Config struct.
	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
	}

	// Print non-empty configuration fields using a utility function.
	utils.PrintNonEmptyFields("", cfg)
	return &cfg, nil
}
