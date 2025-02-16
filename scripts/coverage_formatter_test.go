package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput captures stdout during the execution of f.
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

// alwaysErrorReader is a reader that always returns an error.
type alwaysErrorReader struct{}

func (r alwaysErrorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

// TestColorizeCoverage tests the colorizeCoverage function.
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := colorizeCoverage(tc.input)
			if got != tc.expected {
				t.Errorf("colorizeCoverage(%q) = %q; expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestColorizeCoverageInLine tests colorizeCoverageInLine.
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

// TestFormatPathAndFile tests formatPathAndFile.
func TestFormatPathAndFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"WithDir", "github.com/foo/bar/file.go",
			dirStyle.Sprintf("%s/", filepath.Dir("github.com/foo/bar/file.go")) +
				fileStyle.Sprint(filepath.Base("github.com/foo/bar/file.go"))},
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

// TestStyleCoverageLineDetailed tests styleCoverageLine with a detailed line.
func TestStyleCoverageLineDetailed(t *testing.T) {
	fullPath := "github.com/foo/bar/file.go"
	lineNumber := "31"
	spacingBeforeFunc := "     "
	funcName := "init"
	spacingBeforeCoverage := "           "
	coverageStr := "100.0%"
	input := fullPath + ":" + lineNumber + ":" + spacingBeforeFunc + funcName + spacingBeforeCoverage + coverageStr
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

// TestStyleCoverageLineFallback tests styleCoverageLine with fallback input.
func TestStyleCoverageLineFallback(t *testing.T) {
	input := "total: (statements) 70.0%"
	expected := colorizeCoverageInLine(input)
	got := styleCoverageLine(input)
	if got != expected {
		t.Errorf("styleCoverageLine(fallback) = %q; expected %q", got, expected)
	}
}

// TestRunError tests run() with an error using alwaysErrorReader.
func TestRunError(t *testing.T) {
	err := run(alwaysErrorReader{})
	if err == nil {
		t.Error("Expected run() to return error, got nil")
	}
}

// Save original values to restore after tests.
var (
	originalInputReader = inputReader
	originalExitFunc    = exitFunc
)

// TestMainNoError tests main() when run() succeeds.
func TestMainNoError(t *testing.T) {
	exitCalled := false
	exitCode := 0
	exitFunc = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { exitFunc = originalExitFunc }()

	inputReader = strings.NewReader("total: (statements) 70.0%\n")
	defer func() { inputReader = originalInputReader }()

	main()
	if exitCalled {
		t.Errorf("Expected main() not to call exitFunc, but it was called with code %d", exitCode)
	}
}

// TestMainWithError tests main() when run() returns an error.
func TestMainWithError(t *testing.T) {
	exitCalled := false
	exitCode := 0
	exitFunc = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { exitFunc = originalExitFunc }()

	inputReader = alwaysErrorReader{}
	defer func() { inputReader = originalInputReader }()

	main()
	if !exitCalled {
		t.Error("Expected main() to call exitFunc due to error, but it was not called")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// TestMainIntegration tests main() by overriding inputReader and capturing output.
func TestMainIntegration(t *testing.T) {
	inputLines := []string{
		"github.com/foo/bar/file.go:31:     init           100.0%",
		"total: (statements) 70.0%",
	}
	input := strings.Join(inputLines, "\n")
	inputReader = strings.NewReader(input)
	defer func() { inputReader = originalInputReader }()

	output := captureOutput(func() { main() })
	scanner := bufio.NewScanner(strings.NewReader(output))
	var outputLines []string
	for scanner.Scan() {
		outputLines = append(outputLines, scanner.Text())
	}
	if len(outputLines) != len(inputLines) {
		t.Fatalf("Expected %d output lines, got %d", len(inputLines), len(outputLines))
	}
	detailedExpected := styleCoverageLine(inputLines[0])
	if outputLines[0] != detailedExpected {
		t.Errorf("Main integration detailed line = %q; expected %q", outputLines[0], detailedExpected)
	}
	fallbackExpected := styleCoverageLine(inputLines[1])
	if outputLines[1] != fallbackExpected {
		t.Errorf("Main integration fallback line = %q; expected %q", outputLines[1], fallbackExpected)
	}
}
