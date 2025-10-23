package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFmrBinary(t *testing.T) {
	// Use a temporary directory for the test binary
	tmpDir := t.TempDir()
	binaryName := "test-fmr"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tmpDir, binaryName)

	// Test that the binary can be built
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build fmr binary: %v", err)
	}

	// Test --version flag
	cmd = exec.Command(binaryPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run fmr --version: %v", err)
	}

	if !strings.Contains(string(output), "fmr version") {
		t.Errorf("Expected version output, got: %s", output)
	}

	// Test --help flag
	cmd = exec.Command(binaryPath, "--help")
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run fmr --help: %v", err)
	}

	if !strings.Contains(string(output), "FileMirror") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

func TestMain(t *testing.T) {
	// This is a simple test to ensure the main function runs without panicking.
	// It doesn't test the full functionality of the application.
	go main()
}
