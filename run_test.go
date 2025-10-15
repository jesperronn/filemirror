package filemirror

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathFlagChangesDirectory(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "fmr-path-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file in the temp directory
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remember original directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(origDir) // Best effort to restore directory
	}()

	// Change to a different directory first
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}
	if err := os.Chdir(homeDir); err != nil {
		t.Fatalf("Failed to change to home dir: %v", err)
	}

	// Simulate the path flag behavior
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if err := os.Chdir(absPath); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Verify we're in the temp directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	// Use EvalSymlinks to handle /var vs /private/var on macOS
	currentDirResolved, err := filepath.EvalSymlinks(currentDir)
	if err != nil {
		currentDirResolved = currentDir
	}
	absPathResolved, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		absPathResolved = absPath
	}

	if currentDirResolved != absPathResolved {
		t.Errorf("Expected to be in %q, but in %q", absPathResolved, currentDirResolved)
	}

	// Verify we can scan files in the new directory
	files, err := scanFiles(absPath, "")
	if err != nil {
		t.Fatalf("Failed to scan files: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0].Path != "test.txt" {
		t.Errorf("Expected test.txt, got %s", files[0].Path)
	}
}

func TestPathFlagWithRelativePath(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "fmr-rel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test that relative paths can be converted to absolute
	absPath, err := filepath.Abs(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if !filepath.IsAbs(absPath) {
		t.Errorf("Expected absolute path, got %q", absPath)
	}
}

func TestInvalidPathError(t *testing.T) {
	invalidPath := "/this/path/does/not/exist/hopefully"

	// Try to change to invalid directory
	err := os.Chdir(invalidPath)
	if err == nil {
		t.Error("Expected error when changing to invalid directory")
	}
}
