// File: pkg/config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/heinrichb/scrapey-cli/pkg/utils"
)

// Config holds configuration data.
type Config struct {
	URL        string `json:"url,omitempty"`
	PARSERULES struct {
		TITLE           string `json:"title,omitempty"`
		METADESCRIPTION string `json:"metaDescription,omitempty"`
	}
}

// Load reads config from the specified filePath.
// Returns an error if the file does not exist or if the JSON is invalid.
func Load(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", filePath)
	}

	fmt.Printf("%s%s\n", color.New(color.FgHiGreen).Sprint("Loaded config from: "), filePath)

	// Read the file contents using os.ReadFile (replacing deprecated ioutil.ReadFile).
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal the JSON into a Config struct.
	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file: %v", err)
	}
	utils.PrintNonEmptyFields("", cfg)
	return &cfg, nil
}
