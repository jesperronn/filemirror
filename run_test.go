package filemirror

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

var exitCalled bool

func mockExit(code int) {
	exitCalled = true
}

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

func TestPrintHelp(t *testing.T) {
	// Keep old stdout
	old := os.Stdout
	// Create a new pipe
	r, w, _ := os.Pipe()
	// Set stdout to the new pipe
	os.Stdout = w

	PrintHelp()

	// Close the writer
	w.Close()
	// Restore old stdout
	os.Stdout = old

	// Read the output
	var buf strings.Builder
	io.Copy(&buf, r)

	// Check the output
	if !strings.Contains(buf.String(), "fmr (FileMirror) - Interactive file synchronization tool") {
		t.Errorf("Expected help message to contain 'fmr (FileMirror) - Interactive file synchronization tool', but it didn't")
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "version flag",
			args: []string{"-v"},
			want: "fmr version 0.1.0",
		},
		{
			name: "version long flag",
			args: []string{"--version"},
			want: "fmr version 0.1.0",
		},
		{
			name: "help flag",
			args: []string{"-h"},
			want: "fmr (FileMirror) - Interactive file synchronization tool",
		},
		{
			name: "help long flag",
			args: []string{"--help"},
			want: "fmr (FileMirror) - Interactive file synchronization tool",
		},
		{
			name: "path flag with valid path",
			args: []string{"-p", "."},
			want: "", // No output expected, but the chdir should succeed
		},
		{
			name: "path flag with invalid path",
			args: []string{"-p", "/invalid/path"},
			want: "Error: cannot change to directory",
		},
		{
			name: "initial query",
			args: []string{"my-query"},
			want: "", // No output expected, but the query should be set
		},
		{
			name: "exit summary",
			args: []string{},
			want: "my exit summary",
		},
	}

	// Mock newProgram
	oldNewProgram := newProgram
	newProgram = func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		if mm, ok := m.(model); ok {
			mm.exitSummary = "my exit summary"
			return tea.NewProgram(mm)
		}
		return tea.NewProgram(m)
	}
	defer func() { newProgram = oldNewProgram }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock osExit
			oldOsExit := osExit
			defer func() { osExit = oldOsExit }()
			exitCalled = false
			osExit = mockExit

			// Keep old stdout
			old := os.Stdout
			// Create a new pipe
			r, w, _ := os.Pipe()
			// Set stdout to the new pipe
			os.Stdout = w

			// Set os.Args
			os.Args = append([]string{"fmr"}, tt.args...)

			Run()

			// Close the writer
			w.Close()
			// Restore old stdout
			os.Stdout = old

			// Read the output
			var buf strings.Builder
			io.Copy(&buf, r)

			// Check the output
			if tt.want != "" && !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected output to contain %q, but got %q", tt.want, buf.String())
			}
		})
	}
}
