package filemirror

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %q, want %q", tt.size, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly ten!!", 10, "exactly..."},
		{"this is a very long string", 10, "this is..."},
		{"", 10, ""},
		{"test", 4, "test"},
		{"testing", 4, "t..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestFilterFiles(t *testing.T) {
	m := InitialModel("", ".")
	m.files = []FileInfo{
		{Path: "main.go"},
		{Path: "model.go"},
		{Path: "scanner.go"},
		{Path: "config.json"},
		{Path: "README.md"},
	}

	// Test filtering with no query
	m.searchInput.SetValue("")
	m.filterFiles()
	if len(m.filteredFiles) != len(m.files) {
		t.Errorf("Expected %d files with empty query, got %d",
			len(m.files), len(m.filteredFiles))
	}

	// Test filtering with query
	m.searchInput.SetValue("go")
	m.filterFiles()
	expectedCount := 3 // main.go, model.go, scanner.go
	if len(m.filteredFiles) != expectedCount {
		t.Errorf("Expected %d files with 'go' query, got %d",
			expectedCount, len(m.filteredFiles))
	}

	// Test filtering with specific query
	m.searchInput.SetValue("main")
	m.filterFiles()
	if len(m.filteredFiles) != 1 {
		t.Errorf("Expected 1 file with 'main' query, got %d", len(m.filteredFiles))
	}
	if len(m.filteredFiles) > 0 && m.filteredFiles[0].Path != "main.go" {
		t.Errorf("Expected 'main.go', got %q", m.filteredFiles[0].Path)
	}
}

func TestMinMax(t *testing.T) {
	// Test minInt
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"5 and 10", 5, 10, 5},
		{"10 and 5", 10, 5, 5},
		{"equal values", 7, 7, 7},
		{"negative and positive", -5, 10, -5},
		{"both negative", -10, -5, -10},
		{"zero and positive", 0, 5, 0},
		{"zero and negative", 0, -5, -5},
	}

	for _, tt := range tests {
		t.Run("minInt_"+tt.name, func(t *testing.T) {
			result := minInt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("minInt(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}

	// Test maxInt
	maxTests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"5 and 10", 5, 10, 10},
		{"10 and 5", 10, 5, 10},
		{"equal values", 7, 7, 7},
		{"negative and positive", -5, 10, 10},
		{"both negative", -10, -5, -5},
		{"zero and positive", 0, 5, 5},
		{"zero and negative", 0, -5, 0},
	}

	for _, tt := range maxTests {
		t.Run("maxInt_"+tt.name, func(t *testing.T) {
			result := maxInt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("maxInt(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestGenerateDiff(t *testing.T) {
	m := InitialModel("", ".")
	m.sourceFile = &FileInfo{Path: "source.txt"}

	tests := []struct {
		name           string
		source         string
		target         string
		expectContains []string
	}{
		{
			name:   "identical files",
			source: "line1\nline2\nline3",
			target: "line1\nline2\nline3",
			expectContains: []string{
				" line1",
				" line2",
				" line3",
			},
		},
		{
			name:   "completely different",
			source: "old1\nold2",
			target: "new1\nnew2",
			expectContains: []string{
				"-old1",
				"+new1",
				"-old2",
				"+new2",
			},
		},
		{
			name:   "line added",
			source: "line1\nline2",
			target: "line1\nline2\nline3",
			expectContains: []string{
				" line1",
				" line2",
				"+line3",
			},
		},
		{
			name:   "line removed",
			source: "line1\nline2\nline3",
			target: "line1\nline3",
			expectContains: []string{
				" line1",
				"-line2",
				"+line3",
			},
		},
		{
			name:   "empty source",
			source: "",
			target: "new line",
			expectContains: []string{
				"+new line",
			},
		},
		{
			name:   "empty target",
			source: "old line",
			target: "",
			expectContains: []string{
				"-old line",
			},
		},
		{
			name:   "mixed changes",
			source: "keep1\nchange\nkeep2",
			target: "keep1\nmodified\nkeep2",
			expectContains: []string{
				" keep1",
				"-change",
				"+modified",
				" keep2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := m.generateDiff(tt.source, tt.target)

			// Check that expected patterns are in the diff
			for _, expected := range tt.expectContains {
				found := false
				for _, line := range diff {
					if line == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected diff to contain %q, but it wasn't found.\nFull diff: %v", expected, diff)
				}
			}

			// Verify header is present
			if len(diff) > 0 && !strings.Contains(diff[0], "Source:") {
				t.Errorf("Expected diff header with 'Source:', got: %q", diff[0])
			}
		})
	}
}

func TestCopySourceToTargets(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "fmr-copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "test content for copying"
	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create target files
	target1 := filepath.Join(tmpDir, "target1.txt")
	target2 := filepath.Join(tmpDir, "target2.txt")
	if err := os.WriteFile(target1, []byte("old content 1"), 0o644); err != nil {
		t.Fatalf("Failed to create target1: %v", err)
	}
	if err := os.WriteFile(target2, []byte("old content 2"), 0o644); err != nil {
		t.Fatalf("Failed to create target2: %v", err)
	}

	tests := []struct {
		name        string
		sourceFile  *FileInfo
		selected    map[int]bool
		files       []FileInfo
		expectError bool
	}{
		{
			name: "copy to single target",
			sourceFile: &FileInfo{
				Path: sourceFile,
			},
			selected: map[int]bool{0: true},
			files: []FileInfo{
				{Path: target1},
			},
			expectError: false,
		},
		{
			name: "copy to multiple targets",
			sourceFile: &FileInfo{
				Path: sourceFile,
			},
			selected: map[int]bool{0: true, 1: true},
			files: []FileInfo{
				{Path: target1},
				{Path: target2},
			},
			expectError: false,
		},
		{
			name: "no targets selected",
			sourceFile: &FileInfo{
				Path: sourceFile,
			},
			selected:    map[int]bool{},
			files:       []FileInfo{{Path: target1}},
			expectError: false, // Should succeed with no operation
		},
		{
			name:        "no source file",
			sourceFile:  nil,
			selected:    map[int]bool{0: true},
			files:       []FileInfo{{Path: target1}},
			expectError: true,
		},
		{
			name: "selected index out of bounds",
			sourceFile: &FileInfo{
				Path: sourceFile,
			},
			selected: map[int]bool{0: true, 5: true}, // 5 is out of bounds
			files: []FileInfo{
				{Path: target1},
			},
			expectError: false, // Should succeed by skipping invalid index
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := InitialModel("", tmpDir)
			m.sourceFile = tt.sourceFile
			m.selected = tt.selected
			m.filteredFiles = tt.files

			err := m.copySourceToTargets()

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify files were copied correctly
			if !tt.expectError && tt.sourceFile != nil {
				for idx, selected := range tt.selected {
					if selected && idx < len(tt.files) {
						content, err := os.ReadFile(tt.files[idx].Path)
						if err != nil {
							t.Errorf("Failed to read target file: %v", err)
							continue
						}
						if string(content) != sourceContent {
							t.Errorf("Target file content = %q, want %q", string(content), sourceContent)
						}
					}
				}
			}
		})
	}
}

func TestMatchesFilePattern(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		pattern  string
		want     bool
	}{
		{
			name:     "exact match",
			filePath: "/path/to/test.go",
			pattern:  "test.go",
			want:     true,
		},
		{
			name:     "contains match",
			filePath: "/path/to/mytest.go",
			pattern:  "test",
			want:     true,
		},
		{
			name:     "case insensitive match",
			filePath: "/path/to/TEST.go",
			pattern:  "test",
			want:     true,
		},
		{
			name:     "no match",
			filePath: "/path/to/file.go",
			pattern:  "test",
			want:     false,
		},
		{
			name:     "path contains pattern",
			filePath: "/test/path/to/file.go",
			pattern:  "test",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilePattern(tt.filePath, tt.pattern)
			if got != tt.want {
				t.Errorf("matchesFilePattern(%q, %q) = %v, want %v",
					tt.filePath, tt.pattern, got, tt.want)
			}
		})
	}
}
