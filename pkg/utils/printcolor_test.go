// File: pkg/utils/printcolor_test.go

package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// init forces ANSI color output during tests.
func init() {
	os.Setenv("TERM", "xterm-256color")
	color.NoColor = false
}

// captureStdout redirects os.Stdout during the execution of f() and returns the captured output.
func captureStdout(f func()) string {
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

// TestPrintColored exercises all branches of the PrintColored function in a table-driven test.
// Each case is documented to explain what branch of PrintColored is being hit.
func TestPrintColored(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		// expectedContains lists substrings that must appear in the output.
		expectedContains []string
		// expectEmpty indicates that no output should be produced.
		expectEmpty bool
	}{
		{
			name:        "No arguments: nothing should be printed",
			args:        []interface{}{},
			expectEmpty: true,
		},
		{
			name:        "Non-string first argument: invalid input produces no output",
			args:        []interface{}{123},
			expectEmpty: true,
		},
		{
			name:             "Single string: prints in default white",
			args:             []interface{}{"Just a test string"},
			expectedContains: []string{"Just a test string", "\x1b["},
		},
		{
			name:             "Two strings with individual color: prefix and secondary",
			args:             []interface{}{"Prefix: ", "Value", color.FgHiGreen},
			expectedContains: []string{"Prefix: ", "Value", "\x1b[92m"}, // \x1b[92m represents high-intensity green.
		},
		{
			name: "Dynamic mode with valid colors: multiple segments with corresponding colors",
			args: []interface{}{
				[]string{"Segment1 ", "Segment2 ", "Segment3"},
				[]color.Attribute{color.FgHiGreen, color.FgHiMagenta},
			},
			expectedContains: []string{
				"Segment1 ", "Segment2 ", "Segment3",
				"\x1b[92m", // ANSI code for high-intensity green.
				"\x1b[95m", // ANSI code for high-intensity magenta.
			},
		},
		{
			name: "Dynamic mode with invalid second argument: defaults to white",
			args: []interface{}{
				[]string{"Only segment"},
				123, // Invalid second argument; triggers default white.
			},
			expectedContains: []string{"Only segment", "\x1b["},
		},
		{
			name: "Mixed arguments with a slice of colors: unpacking color slice correctly",
			args: []interface{}{
				"Mixed: ",
				"Value",
				[]color.Attribute{color.FgHiYellow, color.FgHiBlue},
			},
			// fatih/color combines the two attributes into one ANSI sequence (\x1b[93;94m)
			expectedContains: []string{
				"Mixed: ", "Value",
				"\x1b[93;94m", // Combined ANSI sequence for high-intensity yellow and blue.
			},
		},
	}

	// Iterate over each test case.
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			output := captureStdout(func() {
				PrintColored(tc.args...)
			})
			if tc.expectEmpty {
				if output != "" {
					t.Errorf("Expected no output, but got: %q", output)
				}
				return
			}
			// Verify that each expected substring is present in the output.
			for _, substr := range tc.expectedContains {
				if !strings.Contains(output, substr) {
					t.Errorf("Test case %q: expected output to contain %q, but it did not. Full output: %q", tc.name, substr, output)
				}
			}
		})
	}
}

// TestFprintColoredEmptyAttrs directly tests FprintColored with an empty attribute slice.
// This exercise the branch where no color attributes are provided,
// which is never reached via PrintColored (because it always supplies a default).
func TestFprintColoredEmptyAttrs(t *testing.T) {
	var buf bytes.Buffer

	// Case 1: secondary is empty.
	FprintColored(&buf, "directTest", "")
	output := buf.String()
	if !strings.Contains(output, "directTest") {
		t.Errorf("Expected output to contain %q, got %q", "directTest", output)
	}
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Expected output to end with a newline, got %q", output)
	}

	// Reset buffer for next case.
	buf.Reset()

	// Case 2: secondary is non-empty.
	FprintColored(&buf, "directTest", "Extra")
	output = buf.String()
	if !strings.Contains(output, "directTest") || !strings.Contains(output, "Extra") {
		t.Errorf("Expected output to contain both %q and %q, got %q", "directTest", "Extra", output)
	}
}
