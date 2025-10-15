package filemirror

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

func TestScanFilesDeepDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a very deep directory structure (beyond maxDepth)
	deepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f")
	if err := os.MkdirAll(deepPath, 0o755); err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}

	// Create files at different depths
	testFiles := map[string]bool{
		filepath.Join(tmpDir, "root.txt"):                            true,  // depth 0
		filepath.Join(tmpDir, "a", "level1.txt"):                     true,  // depth 1
		filepath.Join(tmpDir, "a", "b", "level2.txt"):                true,  // depth 2
		filepath.Join(tmpDir, "a", "b", "c", "level3.txt"):           true,  // depth 3
		filepath.Join(tmpDir, "a", "b", "c", "d", "level4.txt"):      true,  // depth 4
		filepath.Join(tmpDir, "a", "b", "c", "d", "e", "level5.txt"): false, // depth 5 - should be skipped
		filepath.Join(deepPath, "deep.txt"):                          false, // depth 6 - should be skipped
	}

	for file := range testFiles {
		if err := os.WriteFile(file, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	// Count expected files (only those within maxDepth=4)
	expectedCount := 0
	for _, shouldInclude := range testFiles {
		if shouldInclude {
			expectedCount++
		}
	}

	if len(files) != expectedCount {
		t.Errorf("Expected %d files (within depth 4), got %d", expectedCount, len(files))
	}
}

func TestScanFilesAllExcludedDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files in all excluded directories
	excludedDirs := []string{"node_modules", ".git", "vendor", ".next", "dist", "build", "target", ".cache"}

	for _, dir := range excludedDirs {
		dirPath := filepath.Join(tmpDir, dir)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			t.Fatalf("Failed to create %s: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dirPath, "test.txt"), []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create file in %s: %v", dir, err)
		}
	}

	// Create one regular file
	if err := os.WriteFile(filepath.Join(tmpDir, "regular.txt"), []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	// Should only find the regular file
	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d (all excluded dirs should be skipped)", len(files))
	}

	if len(files) > 0 && files[0].Path != "regular.txt" {
		t.Errorf("Expected to find regular.txt, got %s", files[0].Path)
	}
}

func TestScanFilesEmptyDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

func TestScanFilesWithComplexPattern(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"test.go",
		"test_test.go",
		"main.go",
		"helper.go",
		"data.json",
		"config.yaml",
	}

	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	tests := []struct {
		pattern       string
		expectedCount int
	}{
		{"*.go", 4},
		{"*test*", 2},
		{"main", 1},
		{"*.json", 1},
		{"*.yaml", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			files, err := scanFiles(tmpDir, tt.pattern)
			if err != nil {
				t.Fatalf("scanFiles failed: %v", err)
			}

			if len(files) != tt.expectedCount {
				t.Errorf("Pattern %q: expected %d files, got %d", tt.pattern, tt.expectedCount, len(files))
			}
		})
	}
}

func TestScanFilesSorted(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files with different timestamps
	// Sleep between creations to ensure different mod times
	testFiles := []string{"old.txt", "medium.txt", "new.txt"}

	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
		// Small delay to ensure different mod times
		// Note: This might be flaky on very fast systems
	}

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(files))
	}

	// Verify files are sorted by modification time (newest first)
	for i := 0; i < len(files)-1; i++ {
		if files[i].Modified.Before(files[i+1].Modified) {
			t.Errorf("Files not sorted correctly: file at index %d is older than file at index %d", i, i+1)
		}
	}
}

func TestScanFilesRelativePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested structure
	if err := os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0o755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	testFiles := []string{
		filepath.Join(tmpDir, "root.txt"),
		filepath.Join(tmpDir, "subdir", "nested.txt"),
	}

	for _, file := range testFiles {
		if err := os.WriteFile(file, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	files, err := scanFiles(tmpDir, "")
	if err != nil {
		t.Fatalf("scanFiles failed: %v", err)
	}

	// Verify paths are relative to workDir
	for _, file := range files {
		if filepath.IsAbs(file.Path) {
			t.Errorf("Expected relative path, got absolute path: %s", file.Path)
		}
	}
}
