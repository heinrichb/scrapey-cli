package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/heinrichb/scrapey-cli/pkg/config"
)

var (
    configPath string
    url        string
)

func init() {
    // Two flags for config (-config and -c) so user can choose either
    flag.StringVar(&configPath, "config", "", "Path to config file")
    flag.StringVar(&configPath, "c", "", "Path to config file (shorthand)")

    // Flag for URL
    flag.StringVar(&url, "url", "", "URL to scrape (overrides config)")
}

func main() {
    // Parse CLI flags
    flag.Parse()

    // Print a colored welcome message
    color.Cyan("Welcome to Scrapey CLI!")

    // If no config path is provided, we can optionally default to something:
    if configPath == "" {
        configPath = "configs/default.json"
    }

	// Show a colored message for loading config
	fmt.Printf("%s%s\n", color.New(color.FgHiYellow).Sprint("Loading config from: "), configPath)
    // Attempt to load config
    cfg, err := config.Load(configPath)
    if err != nil {
        color.Red("Failed to load config: %v", err)
        os.Exit(1)
    }

    // If user supplied a URL, override config
    if url != "" {
		fmt.Printf("%s%s\n", color.New(color.FgHiMagenta).Sprint("Overriding config with URL flag: "), url)
		cfg.URL = url
    }

	// Another colored message to confirm it's loaded
	fmt.Printf("%s%s\n", color.New(color.FgHiGreen).Sprint("Loaded config from: "), configPath)

	// Indicate successful finish for now
	color.Green("Scrapey CLI initialization complete.")
}
