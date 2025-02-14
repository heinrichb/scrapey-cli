// File: cmd/scrapeycli/main_test.go
package main

import (
	"flag"
	"os/exec"
	"strings"
	"testing"
)

// runMainCommand is a helper function that executes main.go with the provided arguments.
// It sets the working directory to the project root (two levels up from cmd/scrapeycli/)
// and returns the combined output along with any error.
func runMainCommand(_ *testing.T, args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "./cmd/scrapeycli/main.go"}, args...)...)
	cmd.Dir = "../.." // Set working directory to project root.
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// This test verifies that all necessary command-line flags are properly registered.
// The application depends on these flags to receive configuration input and URL overrides.
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

// This test executes the main program with a valid configuration file.
// It verifies that the program initializes correctly by checking for the welcome message in the output.
func TestMainExecution(t *testing.T) {
	output, err := runMainCommand(t, "--config", "configs/default.json")
	if err != nil {
		t.Fatalf("Failed to run main.go: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Welcome to Scrapey CLI!") {
		t.Errorf("Expected welcome message not found in output.\nOutput: %s", output)
	}
}

// This test simulates a configuration load failure by providing a non-existent config file.
// It confirms that the program correctly handles the error by exiting with a status code of 1.
func TestMainConfigFailure(t *testing.T) {
	_, err := runMainCommand(t, "--config", "nonexistent.json")
	if err == nil {
		t.Fatalf("Expected failure due to config load error, but got success")
	}

	// Verify that the error is of type *exec.ExitError and that the exit code is 1.
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Error was not of type *exec.ExitError: %v", err)
	}
}

// This test runs the main program with both a valid configuration file and a URL override parameter.
// It checks that the output includes the URL override message and that the overridden URL is displayed,
// confirming that the URL override branch in the code is executed.
func TestURLOverride(t *testing.T) {
	output, err := runMainCommand(t, "--config", "configs/default.json", "--url", "https://example.org")
	if err != nil {
		t.Fatalf("Failed to run main.go with URL override: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "Overriding config with URL flag:") {
		t.Errorf("Expected URL override message not found in output.\nOutput: %s", output)
	}
	if !strings.Contains(output, "https://example.org") {
		t.Errorf("Expected overridden URL not found in output.\nOutput: %s", output)
	}
}
