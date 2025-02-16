// File: cmd/scrapeycli/main_test.go

package main

import (
	"flag"
	"os/exec"
	"strings"
	"testing"
)

/*
runMainCommand is a helper function that executes main.go with the provided arguments.
It sets the working directory to the project root (two levels up from cmd/scrapeycli/)
and returns the combined output along with any error.

Parameters:
  - t: The current testing context.
  - args: A variadic list of arguments to be passed to the go run command.

Usage:

	output, err := runMainCommand(t, "--config", "configs/default.json")
*/
func runMainCommand(_ *testing.T, args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "./cmd/scrapeycli/main.go"}, args...)...)
	cmd.Dir = "../.." // Set working directory to project root.
	output, err := cmd.CombinedOutput()
	return string(output), err
}

/*
TestFlagRegistration verifies that all necessary command-line flags are properly registered.
The application depends on these flags for configuration input and CLI overrides.
*/
func TestFlagRegistration(t *testing.T) {
	expectedFlags := []string{"config", "c", "url", "maxDepth", "rateLimit"}
	for _, flagName := range expectedFlags {
		if f := flag.Lookup(flagName); f == nil {
			t.Errorf("Expected flag '%s' to be registered", flagName)
		}
	}
}

/*
TestMainExecution runs the main program with a valid configuration file
and ensures it initializes correctly.
*/
func TestMainExecution(t *testing.T) {
	output, err := runMainCommand(t)
	if err != nil {
		t.Fatalf("Failed to run main.go: %v\nOutput: %s", err, output)
	}

	// Define expected phrases used multiple times
	requiredPhrases := []string{
		"Welcome to Scrapey CLI!",
		"Scrapey CLI initialization complete.",
		"Base URL: https://example.com",
	}

	// Validate presence of required phrases
	for _, phrase := range requiredPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain '%s'.\nOutput: %s", phrase, output)
		}
	}
}

/*
TestMainConfigFailure simulates a config load failure by specifying a non-existent file.
It checks that the program exits with code 1, indicating proper error handling.
*/
func TestMainConfigFailure(t *testing.T) {
	_, err := runMainCommand(t, "--config", "nonexistent.json")
	if err == nil {
		t.Fatalf("Expected failure due to config load error, but got success")
	}

	// Validate correct exit behavior
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Error was not of type *exec.ExitError: %v", err)
	}
}

/*
TestCLIOverrides verifies that CLI arguments correctly override the configuration.

It ensures that:
  - The base URL can be overridden.
  - Scraping depth (maxDepth) can be overridden.
  - Rate limit can be overridden.

The test **does not rely on exact print statements** to avoid fragility.
*/
func TestCLIOverrides(t *testing.T) {
	// CLI argument values (used multiple times)
	newBaseURL := "https://cli-example.com"
	newMaxDepth := "10"
	newRateLimit := "2.5"

	// Run command
	output, err := runMainCommand(t,
		"--url", newBaseURL,
		"--maxDepth", newMaxDepth,
		"--rateLimit", newRateLimit,
	)
	if err != nil {
		t.Fatalf("Failed to run main.go with CLI overrides: %v\nOutput: %s", err, output)
	}

	// Expected CLI override outputs (used multiple times)
	expectedOutputs := map[string]string{
		"Base URL: ":                  newBaseURL,
		"ScrapingOptions.MaxDepth: ":  newMaxDepth,
		"ScrapingOptions.RateLimit: ": newRateLimit,
	}

	// Validate overrides dynamically
	for key, expected := range expectedOutputs {
		if !strings.Contains(output, key+expected) {
			t.Errorf("Expected override '%s%s' not found in output.\nOutput: %s", key, expected, output)
		}
	}
}
