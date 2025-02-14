// File: pkg/config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"reflect"
)

// Config holds configuration data.
type Config struct {
	URL string `json:"url,omitempty"`
	PARSERULES struct {
		TITLE string `json:title,omitempty`
		METADESCRIPTION string `json:metaDescription,omitempty`
	}
}

// Load reads config from the specified filePath.
// Returns an error if the file does not exist or if the JSON is invalid.
func Load(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", filePath)
	}

	fmt.Printf("Loading config from %s\n", filePath)

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
	printNonEmptyFields("", cfg)
    return &cfg, nil
}

func printNonEmptyFields(prefix string, v interface{}) {
	val := reflect.ValueOf(v)

	// Handle pointers or nested structs
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// Iterate over fields
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			printNonEmptyFields(prefix+fieldName+".", field.Interface())
		} else if field.Kind() == reflect.String && field.String() != "" {
			fmt.Println(color.New(color.FgHiYellow).Sprint(prefix+fieldName+":"), field.String())
		}
	}
}