package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/heinrichb/scrapey-cli/pkg/config"
	"github.com/heinrichb/scrapey-cli/pkg/utils"
)

/*
Global variables for storing command-line arguments.

- configPath: The path to the configuration file.
- url: The URL to be scraped, which may override the URL in the config.
- maxDepth: Overrides the scraping depth if set.
- rateLimit: Overrides the request rate limit.
- verbose: Enables verbose output.
*/
var (
	configPath string
	url        string
	maxDepth   int
	rateLimit  float64
	verbose    bool
)

/*
init registers command-line flags for configuration.

It sets up flags for:
- The config file ("config" and its shorthand "c").
- URL override.
- Scraping depth override.
- Rate limit override.
- Verbose output ("verbose" and its shorthand "v").
*/
func init() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file (shorthand)")
	flag.StringVar(&url, "url", "", "URL to scrape (overrides config)")
	flag.IntVar(&maxDepth, "maxDepth", 0, "Override max crawl depth")
	flag.Float64Var(&rateLimit, "rateLimit", 0, "Override request rate limit (seconds)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
}

// Helper functions to create pointers for literal values.
func ptrString(s string) *string    { return &s }
func ptrInt(i int) *int             { return &i }
func ptrFloat64(f float64) *float64 { return &f }

/*
main is the entry point of Scrapey CLI.

It parses command-line flags, prints a welcome message, loads the configuration,
applies CLI overrides using a ConfigOverride object, and prints confirmation messages.
*/
func main() {
	// Parse CLI flags.
	flag.Parse()

	// Store the verbose flag in global state.
	config.Verbose = verbose

	// Print a welcome message in cyan using our PrintColored utility.
	utils.PrintColored("Welcome to Scrapey CLI!", "", color.FgCyan)

	// Default to "configs/default.json" if no config path is provided.
	if configPath == "" {
		configPath = "configs/default.json"
	}

	// Attempt to load the configuration from the specified file.
	cfg, err := config.Load(configPath)
	if err != nil {
		// If loading fails, print an error message in red and exit.
		utils.PrintColored("Failed to load config: ", err.Error(), color.FgRed)
		os.Exit(1)
	}

	// Construct a partial ConfigOverride struct for CLI overrides.
	cliOverrides := config.ConfigOverride{}

	// Apply URL override if provided.
	if url != "" {
		cliOverrides.URL = &struct {
			Base        *string   `json:"base"`
			Routes      *[]string `json:"routes"`
			IncludeBase *bool     `json:"includeBase"`
		}{
			Base: ptrString(url),
		}
	}

	// Apply maxDepth override if provided.
	if maxDepth > 0 {
		if cliOverrides.ScrapingOptions == nil {
			cliOverrides.ScrapingOptions = &struct {
				MaxDepth      *int     `json:"maxDepth"`
				RateLimit     *float64 `json:"rateLimit"`
				RetryAttempts *int     `json:"retryAttempts"`
				UserAgent     *string  `json:"userAgent"`
			}{}
		}
		cliOverrides.ScrapingOptions.MaxDepth = ptrInt(maxDepth)
	}

	// Apply rateLimit override if provided.
	if rateLimit > 0 {
		if cliOverrides.ScrapingOptions == nil {
			cliOverrides.ScrapingOptions = &struct {
				MaxDepth      *int     `json:"maxDepth"`
				RateLimit     *float64 `json:"rateLimit"`
				RetryAttempts *int     `json:"retryAttempts"`
				UserAgent     *string  `json:"userAgent"`
			}{}
		}
		cliOverrides.ScrapingOptions.RateLimit = ptrFloat64(rateLimit)
	}

	// Apply all CLI overrides dynamically.
	cfg.OverrideConfig(cliOverrides)

	// Print confirmation of loaded config.
	utils.PrintColored("Scrapey CLI initialization complete.", "", color.FgGreen)

	// Print which routes will be scraped.
	utils.PrintColored("Base URL: ", cfg.URL.Base, color.FgYellow)
	if cfg.URL.IncludeBase {
		utils.PrintColored("Including base URL in scraping.", "", color.FgGreen)
	}
	for _, route := range cfg.URL.Routes {
		utils.PrintColored("Scraping route: ", route, color.FgHiBlue)
	}
}
