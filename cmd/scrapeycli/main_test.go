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
  - t: The current testing context (not used directly, but conforms to typical test helper function signatures).
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
The application depends on these flags for configuration input and URL overrides.

Checks:
  - "config" and "c" flags
  - "url" flag
*/
func TestFlagRegistration(t *testing.T) {
	if f := flag.Lookup("config"); f == nil {
		t.Error("Expected flag 'config' to be registered")
	}
	if f := flag.Lookup("c"); f == nil {
		t.Error("Expected shorthand flag 'c' to be registered")
	}
	if f := flag.Lookup("url"); f == nil {
		t.Error("Expected flag 'url' to be registered")
	}
}

/*
TestMainExecution runs the main program with a valid configuration file and checks for the expected output.
*/
func TestMainExecution(t *testing.T) {
	output, err := runMainCommand(t, "--config", "configs/default.json")
	if err != nil {
		t.Fatalf("Failed to run main.go: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Welcome to Scrapey CLI!") {
		t.Errorf("Expected welcome message not found in output.\nOutput: %s", output)
	}

	if !strings.Contains(output, "Base URL: https://example.com") {
		t.Errorf("Expected base URL output not found.\nOutput: %s", output)
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

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Error was not of type *exec.ExitError: %v", err)
	}
}

/*
TestURLOverride verifies that specifying a URL via CLI correctly overrides the Base URL.
*/
func TestURLOverride(t *testing.T) {
	output, err := runMainCommand(t, "--config", "configs/default.json", "--url", "https://example.org")
	if err != nil {
		t.Fatalf("Failed to run main.go with URL override: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "Overriding config with URL flag:") {
		t.Errorf("Expected URL override message not found in output.\nOutput: %s", output)
	}
	if !strings.Contains(output, "Base URL: https://example.org") {
		t.Errorf("Expected overridden URL not found in output.\nOutput: %s", output)
	}
}
