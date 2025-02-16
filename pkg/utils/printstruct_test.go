// File: pkg/utils/printstruct_test.go

package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput redirects os.Stdout during the execution of f() and returns the captured output.
// (Renamed from captureStdout to avoid conflict with the same helper in printcolor_test.go.)
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Define sample structs for testing PrintNonEmptyFields.

// SimpleStruct contains a string field and a non-string field. Only the string field should be printed.
type SimpleStruct struct {
	Name string
	Age  int
}

// NestedStruct is used for nesting within another struct.
type NestedStruct struct {
	Title   string
	Comment string
}

// ComplexStruct demonstrates nested structures. Only non-empty string fields should be printed.
type ComplexStruct struct {
	URL    string
	Nested NestedStruct
	Empty  string
	Other  int
}

// TestPrintNonEmptyFields verifies that PrintNonEmptyFields correctly traverses structs
// (or pointers to structs) and prints non-empty string fields with appropriate prefixes.
func TestPrintNonEmptyFields(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string // Expected substrings to be found in the printed output.
	}{
		{
			name:     "SimpleStruct with non-empty Name",
			input:    SimpleStruct{Name: "Alice", Age: 30},
			expected: []string{"Name:", "Alice"},
		},
		{
			name:     "SimpleStruct with empty Name",
			input:    SimpleStruct{Name: "", Age: 40},
			expected: []string{}, // No output expected because Name is empty.
		},
		{
			name: "ComplexStruct with nested non-empty fields",
			input: ComplexStruct{
				URL: "http://example.com",
				Nested: NestedStruct{
					Title:   "Example Title",
					Comment: "",
				},
				Empty: "",
				Other: 10,
			},
			expected: []string{
				"URL:", "http://example.com",
				"Nested.Title:", "Example Title",
			},
		},
		{
			name:     "Pointer to SimpleStruct with non-empty Name",
			input:    &SimpleStruct{Name: "Bob", Age: 25},
			expected: []string{"Name:", "Bob"},
		},
	}

	// Iterate over each test case.
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintNonEmptyFields("", tc.input)
			})
			// If no output is expected, verify that the captured output is empty.
			if len(tc.expected) == 0 && strings.TrimSpace(output) != "" {
				t.Errorf("Test case %q: expected no output, got %q", tc.name, output)
			}
			// Otherwise, verify that each expected substring is present in the output.
			for _, substr := range tc.expected {
				if !strings.Contains(output, substr) {
					t.Errorf("Test case %q: expected output to contain %q, but it did not. Full output: %q", tc.name, substr, output)
				}
			}
		})
	}
}
