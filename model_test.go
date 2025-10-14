package main

import (
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
	m := initialModel("", ".")
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
	if min(5, 10) != 5 {
		t.Errorf("min(5, 10) should be 5")
	}
	if min(10, 5) != 5 {
		t.Errorf("min(10, 5) should be 5")
	}
	if max(5, 10) != 10 {
		t.Errorf("max(5, 10) should be 10")
	}
	if max(10, 5) != 10 {
		t.Errorf("max(10, 5) should be 10")
	}
}
