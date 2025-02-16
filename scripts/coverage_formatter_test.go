// File: scripts/coverage_formatter_test.go

package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput temporarily redirects os.Stdout during the execution of f()
// and returns the captured output.
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

// TestColorizeCoverage tests the colorizeCoverage function with various coverage strings.
func TestColorizeCoverage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"HighCoverage", "100.0%", colorHighCov.Sprint("100.0%")},
		{"MediumCoverage", "60.0%", colorMidCov.Sprint("60.0%")},
		{"LowCoverage", "40.0%", colorLowCov.Sprint("40.0%")},
		{"InvalidCoverage", "foo%", "foo%"},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := colorizeCoverage(tc.input)
			if got != tc.expected {
				t.Errorf("colorizeCoverage(%q) = %q; expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestColorizeCoverageInLine verifies that colorizeCoverageInLine replaces all coverage percentages
// in a fallback line with their appropriately colored versions.
func TestColorizeCoverageInLine(t *testing.T) {
	input := "total: (statements) 70.0%"
	expected := fallbackCoverageRegex.ReplaceAllStringFunc(input, func(match string) string {
		return colorizeCoverage(match)
	})
	got := colorizeCoverageInLine(input)
	if got != expected {
		t.Errorf("colorizeCoverageInLine(%q) = %q; expected %q", input, got, expected)
	}
}

// TestFormatPathAndFile verifies that formatPathAndFile correctly splits and colors the directory
// portion from the file name.
func TestFormatPathAndFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"WithDir", "github.com/foo/bar/file.go",
			dirStyle.Sprintf("%s/", filepath.Dir("github.com/foo/bar/file.go")) + fileStyle.Sprint(filepath.Base("github.com/foo/bar/file.go"))},
		{"NoDir", "file.go", fileStyle.Sprint("file.go")},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := formatPathAndFile(tc.input)
			if got != tc.expected {
				t.Errorf("formatPathAndFile(%q) = %q; expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestStyleCoverageLineDetailed tests styleCoverageLine using a detailed coverage line input that
// matches the detailedCoverageRegex pattern.
func TestStyleCoverageLineDetailed(t *testing.T) {
	// Construct a detailed coverage line.
	// Format: <fullPath>:<lineNumber>:<spaces><funcName><spaces><coverage>
	fullPath := "github.com/foo/bar/file.go"
	lineNumber := "31"
	spacingBeforeFunc := "     "
	funcName := "init"
	spacingBeforeCoverage := "           "
	coverageStr := "100.0%"
	input := fullPath + ":" + lineNumber + ":" + spacingBeforeFunc + funcName + spacingBeforeCoverage + coverageStr

	// Note: our styleCoverageLine rebuilds the line as:
	//   formatPathAndFile(fullPath) + ":" + lineNumStyle.Sprint(lineNumber) + ":" + spacingBeforeFunc +
	//   funcStyle.Sprint(funcName) + spacingBeforeCoverage + colorizeCoverage(coverageStr)
	expected := formatPathAndFile(fullPath) + ":" +
		lineNumStyle.Sprint(lineNumber) + ":" +
		spacingBeforeFunc +
		funcStyle.Sprint(funcName) +
		spacingBeforeCoverage +
		colorizeCoverage(coverageStr)

	got := styleCoverageLine(input)
	if got != expected {
		t.Errorf("styleCoverageLine(detailed) = %q; expected %q", got, expected)
	}
}

// TestStyleCoverageLineFallback tests styleCoverageLine with an input that does not match the detailed pattern.
// In that case, the function should call colorizeCoverageInLine on the entire line.
func TestStyleCoverageLineFallback(t *testing.T) {
	input := "total: (statements) 70.0%"
	expected := colorizeCoverageInLine(input)
	got := styleCoverageLine(input)
	if got != expected {
		t.Errorf("styleCoverageLine(fallback) = %q; expected %q", got, expected)
	}
}

// TestMainIntegration tests the main function by simulating stdin input and capturing stdout output.
// It feeds multiple lines (one detailed and one fallback) and verifies that each line is styled as expected.
func TestMainIntegration(t *testing.T) {
	inputLines := []string{
		"github.com/foo/bar/file.go:31:     init           100.0%",
		"total: (statements) 70.0%",
	}
	input := strings.Join(inputLines, "\n")

	// Create a temporary file and write our input to it.
	tmpFile, err := os.CreateTemp("", "coverage-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(input); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek in temp file: %v", err)
	}

	// Replace os.Stdin with our temp file.
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = tmpFile

	output := captureOutput(func() {
		main()
	})

	scanner := bufio.NewScanner(strings.NewReader(output))
	var outputLines []string
	for scanner.Scan() {
		outputLines = append(outputLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Error scanning output: %v", err)
	}
	if len(outputLines) != len(inputLines) {
		t.Fatalf("Expected %d output lines, got %d", len(inputLines), len(outputLines))
	}

	// Check detailed line.
	detailedExpected := styleCoverageLine(inputLines[0])
	if outputLines[0] != detailedExpected {
		t.Errorf("Main integration detailed line = %q; expected %q", outputLines[0], detailedExpected)
	}

	// Check fallback line.
	fallbackExpected := styleCoverageLine(inputLines[1])
	if outputLines[1] != fallbackExpected {
		t.Errorf("Main integration fallback line = %q; expected %q", outputLines[1], fallbackExpected)
	}
}
