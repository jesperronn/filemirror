package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		filename string
		pattern  string
		expected bool
	}{
		{"main.go", "*.go", true},
		{"main.go", "main", true},
		{"main.go", "test", false},
		{"config.json", "*.json", true},
		{"config.json", "config", true},
		{"README.md", "*.md", true},
		{"test.txt", "*.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename+"_"+tt.pattern, func(t *testing.T) {
			result := matchesPattern(tt.filename, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v",
					tt.filename, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestScanFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldDir) // Best effort to restore directory
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create test files
	testFiles := []string{
		"test1.go",
		"test2.go",
		"config.json",
		"README.md",
		"subdir/test3.go",
	}

	for _, file := range testFiles {
		dir := filepath.Dir(file)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
		if err := os.WriteFile(file, []byte("test content"), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test scanning all files
	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(files))
	}

	// Test scanning with pattern
	files, err = scanFiles(tmpDir, "*.go")
	if err != nil {
		t.Fatalf("scanFiles with pattern failed: %v", err)
	}

	expectedGoFiles := 3
	if len(files) != expectedGoFiles {
		t.Errorf("Expected %d .go files, got %d", expectedGoFiles, len(files))
	}
}

func TestScanFilesExcludesNodeModules(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldDir) // Best effort to restore directory
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create files in node_modules (should be excluded)
	if err := os.MkdirAll("node_modules/package", 0o755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	if err := os.WriteFile("node_modules/package/index.js", []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create file in node_modules: %v", err)
	}

	// Create a regular file
	if err := os.WriteFile("app.js", []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	// Should only find app.js, not the file in node_modules
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d (should exclude node_modules)", len(files))
	}

	if len(files) > 0 && files[0].Path != "app.js" {
		t.Errorf("Expected to find app.js, got %s", files[0].Path)
	}
}
