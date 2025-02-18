package config

import (
	"encoding/json"
	"fmt"
	"os"

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
ConfigOverride represents a partial configuration used for overriding values.
All fields are pointers, so that nil indicates "no override" while a non-nil value,
even if zero, is used to override the corresponding Config field.
*/
type ConfigOverride struct {
	Version *string `json:"version"`
	URL     *struct {
		Base        *string   `json:"base"`
		Routes      *[]string `json:"routes"`
		IncludeBase *bool     `json:"includeBase"`
	} `json:"url"`
	ParseRules *struct {
		Title           *string `json:"title,omitempty"`
		MetaDescription *string `json:"metaDescription,omitempty"`
		ArticleContent  *string `json:"articleContent,omitempty"`
		Author          *string `json:"author,omitempty"`
		DatePublished   *string `json:"datePublished,omitempty"`
	} `json:"parseRules"`
	Storage *struct {
		OutputFormats *[]string `json:"outputFormats"`
		SavePath      *string   `json:"savePath"`
		FileName      *string   `json:"fileName"`
	} `json:"storage"`
	ScrapingOptions *struct {
		MaxDepth      *int     `json:"maxDepth"`
		RateLimit     *float64 `json:"rateLimit"`
		RetryAttempts *int     `json:"retryAttempts"`
		UserAgent     *string  `json:"userAgent"`
	} `json:"scrapingOptions"`
	DataFormatting *struct {
		CleanWhitespace *bool `json:"cleanWhitespace"`
		RemoveHTML      *bool `json:"removeHTML"`
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
OverrideConfig applies values from the provided `overrides` object to the existing configuration.
Only fields with non-nil pointers in the overrides object are applied; all other fields remain unchanged.

Parameters:
  - overrides: A ConfigOverride struct containing only the fields to override.
    A nil pointer indicates that no override should occur for that field.

Usage:

	cfg.OverrideConfig(ConfigOverride{
		URL: &struct {
			Base        *string   `json:"base"`
			Routes      *[]string `json:"routes"`
			IncludeBase *bool     `json:"includeBase"`
		}{
			Base: ptrString("https://example.org"),
		},
		ScrapingOptions: &struct {
			MaxDepth      *int     `json:"maxDepth"`
			RateLimit     *float64 `json:"rateLimit"`
			RetryAttempts *int     `json:"retryAttempts"`
			UserAgent     *string  `json:"userAgent"`
		}{
			MaxDepth: ptrInt(5),
		},
	})

Notes:
  - Only fields with non-nil pointers in `overrides` are applied.
  - This allows partial configuration overrides without unintentionally overwriting existing values.
  - Both struct and non-struct fields are overridden if provided.
*/
func (cfg *Config) OverrideConfig(overrides ConfigOverride) {
	// Override non-struct field: Version.
	if overrides.Version != nil {
		utils.PrintColored("Overriding Version: ", *overrides.Version, color.FgHiMagenta)
		cfg.Version = *overrides.Version
	}

	// Override URL fields.
	if overrides.URL != nil {
		if overrides.URL.Base != nil {
			utils.PrintColored("Overriding URL.Base: ", *overrides.URL.Base, color.FgHiMagenta)
			cfg.URL.Base = *overrides.URL.Base
		}
		if overrides.URL.Routes != nil {
			utils.PrintColored("Overriding URL.Routes: ", fmt.Sprint(*overrides.URL.Routes), color.FgHiMagenta)
			cfg.URL.Routes = *overrides.URL.Routes
		}
		if overrides.URL.IncludeBase != nil {
			utils.PrintColored("Overriding URL.IncludeBase: ", fmt.Sprint(*overrides.URL.IncludeBase), color.FgHiMagenta)
			cfg.URL.IncludeBase = *overrides.URL.IncludeBase
		}
	}

	// Override ParseRules fields.
	if overrides.ParseRules != nil {
		if overrides.ParseRules.Title != nil {
			utils.PrintColored("Overriding ParseRules.Title: ", *overrides.ParseRules.Title, color.FgHiMagenta)
			cfg.ParseRules.Title = *overrides.ParseRules.Title
		}
		if overrides.ParseRules.MetaDescription != nil {
			utils.PrintColored("Overriding ParseRules.MetaDescription: ", *overrides.ParseRules.MetaDescription, color.FgHiMagenta)
			cfg.ParseRules.MetaDescription = *overrides.ParseRules.MetaDescription
		}
		if overrides.ParseRules.ArticleContent != nil {
			utils.PrintColored("Overriding ParseRules.ArticleContent: ", *overrides.ParseRules.ArticleContent, color.FgHiMagenta)
			cfg.ParseRules.ArticleContent = *overrides.ParseRules.ArticleContent
		}
		if overrides.ParseRules.Author != nil {
			utils.PrintColored("Overriding ParseRules.Author: ", *overrides.ParseRules.Author, color.FgHiMagenta)
			cfg.ParseRules.Author = *overrides.ParseRules.Author
		}
		if overrides.ParseRules.DatePublished != nil {
			utils.PrintColored("Overriding ParseRules.DatePublished: ", *overrides.ParseRules.DatePublished, color.FgHiMagenta)
			cfg.ParseRules.DatePublished = *overrides.ParseRules.DatePublished
		}
	}

	// Override Storage fields.
	if overrides.Storage != nil {
		if overrides.Storage.OutputFormats != nil {
			utils.PrintColored("Overriding Storage.OutputFormats: ", fmt.Sprint(*overrides.Storage.OutputFormats), color.FgHiMagenta)
			cfg.Storage.OutputFormats = *overrides.Storage.OutputFormats
		}
		if overrides.Storage.SavePath != nil {
			utils.PrintColored("Overriding Storage.SavePath: ", *overrides.Storage.SavePath, color.FgHiMagenta)
			cfg.Storage.SavePath = *overrides.Storage.SavePath
		}
		if overrides.Storage.FileName != nil {
			utils.PrintColored("Overriding Storage.FileName: ", *overrides.Storage.FileName, color.FgHiMagenta)
			cfg.Storage.FileName = *overrides.Storage.FileName
		}
	}

	// Override ScrapingOptions fields.
	if overrides.ScrapingOptions != nil {
		if overrides.ScrapingOptions.MaxDepth != nil {
			utils.PrintColored("Overriding ScrapingOptions.MaxDepth: ", fmt.Sprint(*overrides.ScrapingOptions.MaxDepth), color.FgHiMagenta)
			cfg.ScrapingOptions.MaxDepth = *overrides.ScrapingOptions.MaxDepth
		}
		if overrides.ScrapingOptions.RateLimit != nil {
			utils.PrintColored("Overriding ScrapingOptions.RateLimit: ", fmt.Sprint(*overrides.ScrapingOptions.RateLimit), color.FgHiMagenta)
			cfg.ScrapingOptions.RateLimit = *overrides.ScrapingOptions.RateLimit
		}
		if overrides.ScrapingOptions.RetryAttempts != nil {
			utils.PrintColored("Overriding ScrapingOptions.RetryAttempts: ", fmt.Sprint(*overrides.ScrapingOptions.RetryAttempts), color.FgHiMagenta)
			cfg.ScrapingOptions.RetryAttempts = *overrides.ScrapingOptions.RetryAttempts
		}
		if overrides.ScrapingOptions.UserAgent != nil {
			utils.PrintColored("Overriding ScrapingOptions.UserAgent: ", *overrides.ScrapingOptions.UserAgent, color.FgHiMagenta)
			cfg.ScrapingOptions.UserAgent = *overrides.ScrapingOptions.UserAgent
		}
	}

	// Override DataFormatting fields.
	if overrides.DataFormatting != nil {
		if overrides.DataFormatting.CleanWhitespace != nil {
			utils.PrintColored("Overriding DataFormatting.CleanWhitespace: ", fmt.Sprint(*overrides.DataFormatting.CleanWhitespace), color.FgHiMagenta)
			cfg.DataFormatting.CleanWhitespace = *overrides.DataFormatting.CleanWhitespace
		}
		if overrides.DataFormatting.RemoveHTML != nil {
			utils.PrintColored("Overriding DataFormatting.RemoveHTML: ", fmt.Sprint(*overrides.DataFormatting.RemoveHTML), color.FgHiMagenta)
			cfg.DataFormatting.RemoveHTML = *overrides.DataFormatting.RemoveHTML
		}
	}
}
