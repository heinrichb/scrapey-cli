// File: cmd/scrapeycli/main_test.go
package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMainExecution(t *testing.T) {
	// Set the working directory to the project root.
	cmd := exec.Command("go", "run", "./cmd/scrapeycli/main.go", "--config", "configs/default.json")
	cmd.Dir = "../.."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run main.go: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "Welcome to Scrapey CLI!") {
		t.Errorf("Expected welcome message not found in output.\nOutput: %s", outStr)
	}
}

func TestMainConfigFailure(t *testing.T) {
	// Set the working directory to the project root.
	cmd := exec.Command("go", "run", "./cmd/scrapeycli/main.go", "--config", "nonexistent.json")
	cmd.Dir = "../.."

	output, err := cmd.CombinedOutput()
	// The command should fail and exit with code 1.
	if err == nil {
		t.Fatalf("Expected failure due to config load error, but got success.\nOutput: %s", string(output))
	}

	// Check if the error is an ExitError and verify that the exit code is 1.
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("Error was not of type *exec.ExitError: %v", err)
	}
}
