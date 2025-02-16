// File: pkg/parser/parser_test.go

package parser

import "testing"

// TestParseHTML verifies that ParseHTML returns an empty map and no error
// regardless of the input provided.
func TestParseHTML(t *testing.T) {
	// Test with a non-empty HTML string.
	html := "<html><body><p>Hello, World!</p></body></html>"
	data, err := ParseHTML(html)
	if err != nil {
		t.Errorf("Expected no error for non-empty input, got %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty map for non-empty input, got %v", data)
	}

	// Test with an empty string.
	data, err = ParseHTML("")
	if err != nil {
		t.Errorf("Expected no error for empty input, got %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty map for empty input, got %v", data)
	}
}
