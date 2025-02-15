// File: scripts/coverage_formatter.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// detailedCoverageRegex matches typical coverage detail lines from `go tool cover -func`.
// Example of a matched line:
//
//	github.com/.../file.go:31:     funcName                100.0%
var detailedCoverageRegex = regexp.MustCompile(`^([^:]+\.go):(\d+):(\s+)(\S+)(\s+)([0-9]+\.[0-9]+%)$`)

// fallbackCoverageRegex matches coverage percentages in lines that do not match
// the above pattern (e.g., "total: (statements) 70.0%").
var fallbackCoverageRegex = regexp.MustCompile(`([0-9]+\.[0-9]+%)`)

// Coverage thresholds for color-coding coverage percentages.
const (
	HighCoverageThreshold   = 80.0
	MediumCoverageThreshold = 50.0
)

// Color styles for different parts of a coverage line.
var (
	dirStyle     = color.New(color.FgWhite).Add(color.Bold)
	fileStyle    = color.New(color.FgCyan).Add(color.Bold)
	lineNumStyle = color.New(color.FgMagenta).Add(color.Bold)
	funcStyle    = color.New(color.FgHiBlue)
	colorHighCov = color.New(color.FgGreen)
	colorMidCov  = color.New(color.FgYellow)
	colorLowCov  = color.New(color.FgRed)
)

// main reads from stdin and prints out styled coverage lines to stdout.
// Usage example:
//
//	go tool cover -func=coverage.out | go run ./scripts/coverage_formatter.go
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		originalLine := scanner.Text()
		styledLine := styleCoverageLine(originalLine)
		fmt.Println(styledLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}

// styleCoverageLine determines whether a line is a "detailed coverage" line
// (i.e., matching our `detailedCoverageRegex`) or a fallback line
// (like "total: (statements) 70.0%"). It returns the colorized version
// while preserving all original spacing.
func styleCoverageLine(line string) string {
	// Attempt to match the "detailed coverage" regex.
	if matches := detailedCoverageRegex.FindStringSubmatch(line); matches != nil {
		// Matches structure:
		//   [1]: fullPath     => e.g. "github.com/.../file.go"
		//   [2]: lineNumber   => e.g. "31"
		//   [3]: spacingBeforeFunc => e.g. "     "
		//   [4]: funcName     => e.g. "init"
		//   [5]: spacingBeforeCoverage => e.g. "           "
		//   [6]: coverageStr  => e.g. "100.0%"

		fullPath := matches[1]
		lineNumber := matches[2]
		spacingBeforeFunc := matches[3]
		funcName := matches[4]
		spacingBeforeCoverage := matches[5]
		coverageString := matches[6]

		coloredFilePath := formatPathAndFile(fullPath)
		coloredLineNumber := lineNumStyle.Sprint(lineNumber)
		coloredFunction := funcStyle.Sprint(funcName)
		coloredCoverage := colorizeCoverage(coverageString)

		// Rebuild the line EXACTLY with the original spacing.
		// e.g.: "github.com/.../file.go:31:" + spacingBeforeFunc + funcName + spacingBeforeCoverage + coverage
		return fmt.Sprintf("%s:%s:%s%s%s%s",
			coloredFilePath,
			coloredLineNumber,
			spacingBeforeFunc,
			coloredFunction,
			spacingBeforeCoverage,
			coloredCoverage,
		)
	}

	// If the line does not match our detailed coverage pattern,
	// we look for coverage percentages (e.g., "70.0%") and colorize them.
	return colorizeCoverageInLine(line)
}

// formatPathAndFile splits a path like "dir/subdir/file.go", coloring the directory
// part differently from the file name. If there's no directory component, it just
// colors the file name. Example output might look like:
//
//	github.com/.../dir/  file.go
//
// with dir vs. file in different colors.
func formatPathAndFile(fullPath string) string {
	dir := filepath.Dir(fullPath)
	file := filepath.Base(fullPath)

	if dir == "." || dir == "" {
		return fileStyle.Sprint(file)
	}
	return dirStyle.Sprintf("%s/", dir) + fileStyle.Sprint(file)
}

// colorizeCoverageInLine scans a fallback line for coverage percentages
// (e.g., "70.0%") and colorizes them in-place.
func colorizeCoverageInLine(line string) string {
	return fallbackCoverageRegex.ReplaceAllStringFunc(line, func(match string) string {
		return colorizeCoverage(match)
	})
}

// colorizeCoverage picks a color based on numeric coverage thresholds, returning
// a styled string. For instance, "100.0%" might return green text, "60.0%" yellow, etc.
func colorizeCoverage(coverageStr string) string {
	rawNumber := strings.TrimSuffix(coverageStr, "%")
	coverageValue, parseErr := strconv.ParseFloat(rawNumber, 64)
	if parseErr != nil {
		return coverageStr // fallback: invalid numeric
	}
	switch {
	case coverageValue >= HighCoverageThreshold:
		return colorHighCov.Sprint(coverageStr)
	case coverageValue >= MediumCoverageThreshold:
		return colorMidCov.Sprint(coverageStr)
	default:
		return colorLowCov.Sprint(coverageStr)
	}
}
