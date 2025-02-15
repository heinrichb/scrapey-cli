// File: cmd/scrapeycli/main.go

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
*/
var (
	configPath string
	url        string
)

/*
init registers command-line flags for configuration.

It sets up two flags for the config file ("config" and its shorthand "c")
and a flag for the URL override.
*/
func init() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file (shorthand)")
	flag.StringVar(&url, "url", "", "URL to scrape (overrides config)")
}

/*
main is the entry point of Scrapey CLI.

It parses command-line flags, prints a welcome message, loads the configuration,
handles URL overrides, and prints confirmation messages for each step.
*/
func main() {
	// Parse CLI flags.
	flag.Parse()

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

	// If a URL is provided via the command line, override the configuration's base URL.
	if url != "" {
		utils.PrintColored("Overriding config with URL flag: ", url, color.FgHiMagenta)
		cfg.URL.Base = url
	}

	// Print confirmation of loaded config.
	utils.PrintColored("Loaded config from: ", configPath, color.FgHiGreen)

	// Indicate that initialization is complete by printing a success message in green.
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
