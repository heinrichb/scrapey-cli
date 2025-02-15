// File: pkg/utils/printstruct.go

package utils

import (
	"reflect"

	"github.com/fatih/color"
)

/*
PrintNonEmptyFields dynamically traverses a struct and prints its non-empty string fields.

Parameters:
  - prefix: A string to prepend to the field name, used to represent nested struct hierarchy (e.g., "Parent.Child.").
  - v: The struct or pointer to a struct to be traversed and inspected.

Usage:

	This function is useful for dynamically inspecting and displaying configurations or other data structures
	where the fields may be optional, and only non-empty values are of interest.

Example:

	Given a struct:

	  type Config struct {
	      URL string
	      Nested struct {
	          Title string
	      }
	  }

	Calling PrintNonEmptyFields("", configInstance) will output something like:

	  URL: http://example.com
	  Nested.Title: Example Title

Notes:
  - This function relies on the reflect package and assumes that the input is a struct or a pointer to a struct.
  - Only string fields are checked for non-emptiness; other types are ignored.
  - Colored output is now handled by the PrintColored utility from this package.
*/
func PrintNonEmptyFields(prefix string, v interface{}) {
	val := reflect.ValueOf(v)

	// Handle pointers by obtaining the element value.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// Iterate over each field of the struct.
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// If the field is a nested struct, recursively print its non-empty fields.
		if field.Kind() == reflect.Struct {
			PrintNonEmptyFields(prefix+fieldName+".", field.Interface())
		} else if field.Kind() == reflect.String && field.String() != "" {
			// Use PrintColored to output the field name (with a colon) in high-intensity yellow,
			// followed by the field value in default formatting.
			PrintColored(prefix+fieldName+": ", field.String(), color.FgHiYellow)
		}
	}
}
