package config

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/fatih/color"
	"reflect"
)

// Config is a basic struct to hold config data.
type Config struct {
	URL string `json:"url,omitempty"`
	PARSERULES struct {
		TITLE string `json:title,omitempty`
		METADESCRIPTION string `json:metaDescription,omitempty`
	}
}

// Load reads config from filePath and returns a Config struct.
// Currently just returns a placeholder.
func Load(filePath string) (*Config, error) {
	// In future, parse JSON from filePath
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cfg Config
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&cfg); err != nil {
        return nil, err
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