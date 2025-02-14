// File: pkg/utils/utils.go

package utils

import (
    "fmt"
    "github.com/fatih/color"
    "reflect"
)
// PrintNonEmptyFields dynamically traverses a struct and prints its non-empty string fields.
//
// Parameters:
// - prefix: A string to prepend to the field name, used to represent nested struct hierarchy (e.g., "Parent.Child.").
// - v: The struct or pointer to a struct to be traversed and inspected.
//
// Usage:
// This function is useful for dynamically inspecting and displaying configurations or other data structures
// where the fields may be optional, and only non-empty values are of interest.
//
// Example:
// Given a struct:
//   type Config struct {
//       URL string
//       Nested struct {
//           Title string
//       }
//   }
// Calling PrintNonEmptyFields("", configInstance) will output:
//   URL: http://example.com
//   Nested.Title: Example Title
//
// Notes:
// - This function relies on the reflect package and assumes that the input is a struct or a pointer to a struct.
// - Only string fields are checked for non-emptiness; other types are ignored.
func PrintNonEmptyFields(prefix string, v interface{}) {
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
			PrintNonEmptyFields(prefix+fieldName+".", field.Interface())
		} else if field.Kind() == reflect.String && field.String() != "" {
			fmt.Println(color.New(color.FgHiYellow).Sprint(prefix+fieldName+":"), field.String())
		}
	}
}