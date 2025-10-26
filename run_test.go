package filemirror

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// chdirMutex ensures tests that use os.Chdir() don't run in parallel
var chdirMutex sync.Mutex

func TestPathFlagChangesDirectory(t *testing.T) {
	// Lock to prevent parallel execution with other tests that use os.Chdir
	chdirMutex.Lock()
	defer chdirMutex.Unlock()

	// Try to recover to a valid directory if current one is invalid
	if _, err := os.Getwd(); err != nil {
		// Current directory is invalid, try to change to a safe location
		if homeDir, err := os.UserHomeDir(); err == nil {
			_ = os.Chdir(homeDir)
		}
	}

	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "fmr-path-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Remember original directory before any other setup
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}

	defer os.RemoveAll(tmpDir)
	defer func() {
		_ = os.Chdir(origDir) // Best effort to restore directory
	}()

	// Create a test file in the temp directory
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

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
	// Lock to prevent parallel execution with other tests that use os.Chdir
	chdirMutex.Lock()
	defer chdirMutex.Unlock()

	// Try to recover to a valid directory if current one is invalid
	if _, err := os.Getwd(); err != nil {
		// Current directory is invalid, try to change to a safe location
		if homeDir, err := os.UserHomeDir(); err == nil {
			_ = os.Chdir(homeDir)
		}
	}

	invalidPath := "/this/path/does/not/exist/hopefully"

	// Save and restore current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(origDir)
	}()

	// Try to change to invalid directory
	err = os.Chdir(invalidPath)
	if err == nil {
		t.Error("Expected error when changing to invalid directory")
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantCfg     Config
		wantErr     bool
		errContains string
	}{
		{
			name: "no arguments",
			args: []string{},
			wantCfg: Config{
				WorkDir:      "",
				InitialQuery: "",
				ShowHelp:     false,
				ShowVersion:  false,
			},
			wantErr: false,
		},
		{
			name: "help flag short",
			args: []string{"-h"},
			wantCfg: Config{
				ShowHelp: true,
			},
			wantErr: false,
		},
		{
			name: "help flag long",
			args: []string{"--help"},
			wantCfg: Config{
				ShowHelp: true,
			},
			wantErr: false,
		},
		{
			name: "help command",
			args: []string{"help"},
			wantCfg: Config{
				ShowHelp: true,
			},
			wantErr: false,
		},
		{
			name: "version flag short",
			args: []string{"-v"},
			wantCfg: Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
		{
			name: "version flag long",
			args: []string{"--version"},
			wantCfg: Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
		{
			name: "version command",
			args: []string{"version"},
			wantCfg: Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
		{
			name: "path flag short with value",
			args: []string{"-p", "/tmp"},
			wantCfg: Config{
				WorkDir: "/tmp",
			},
			wantErr: false,
		},
		{
			name: "path flag long with value",
			args: []string{"--path", "/tmp"},
			wantCfg: Config{
				WorkDir: "/tmp",
			},
			wantErr: false,
		},
		{
			name:        "path flag without value",
			args:        []string{"-p"},
			wantErr:     true,
			errContains: "--path requires a directory argument",
		},
		{
			name: "initial query",
			args: []string{"*.go"},
			wantCfg: Config{
				InitialQuery: "*.go",
			},
			wantErr: false,
		},
		{
			name: "path and query",
			args: []string{"-p", "/tmp", "*.txt"},
			wantCfg: Config{
				WorkDir:      "/tmp",
				InitialQuery: "*.txt",
			},
			wantErr: false,
		},
		{
			name: "query and path",
			args: []string{"*.go", "-p", "/tmp"},
			wantCfg: Config{
				WorkDir:      "/tmp",
				InitialQuery: "*.go",
			},
			wantErr: false,
		},
		{
			name: "multiple non-flag args takes first as query",
			args: []string{"first", "second"},
			wantCfg: Config{
				InitialQuery: "first",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseArgs(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cfg.WorkDir != tt.wantCfg.WorkDir {
				t.Errorf("WorkDir = %q, want %q", cfg.WorkDir, tt.wantCfg.WorkDir)
			}
			if cfg.InitialQuery != tt.wantCfg.InitialQuery {
				t.Errorf("InitialQuery = %q, want %q", cfg.InitialQuery, tt.wantCfg.InitialQuery)
			}
			if cfg.ShowHelp != tt.wantCfg.ShowHelp {
				t.Errorf("ShowHelp = %v, want %v", cfg.ShowHelp, tt.wantCfg.ShowHelp)
			}
			if cfg.ShowVersion != tt.wantCfg.ShowVersion {
				t.Errorf("ShowVersion = %v, want %v", cfg.ShowVersion, tt.wantCfg.ShowVersion)
			}
		})
	}
}

func TestValidateAndSetupWorkDir(t *testing.T) {
	// Save and restore current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir)

	tests := []struct {
		name        string
		workDir     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "empty workDir",
			workDir: "",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			workDir: os.TempDir(),
			wantErr: false,
		},
		{
			name:    "valid relative path",
			workDir: ".",
			wantErr: false,
		},
		{
			name:        "nonexistent directory",
			workDir:     "/nonexistent/directory/that/does/not/exist",
			wantErr:     true,
			errContains: "cannot change to directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absPath, err := validateAndSetupWorkDir(tt.workDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.workDir == "" {
				if absPath != "" {
					t.Errorf("Expected empty path for empty workDir, got %q", absPath)
				}
			} else {
				if absPath == "" {
					t.Error("Expected non-empty absolute path")
				}
				if !filepath.IsAbs(absPath) {
					t.Errorf("Expected absolute path, got %q", absPath)
				}
			}

			// Restore original directory after each test
			os.Chdir(origDir)
		})
	}
}

func TestPrintVersion(t *testing.T) {
	// Save original values
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit
	defer func() {
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	}()

	tests := []struct {
		name            string
		version         string
		buildTime       string
		gitCommit       string
		expectBuildTime bool
		expectGitCommit bool
	}{
		{
			name:            "version only",
			version:         "1.0.0",
			buildTime:       "unknown",
			gitCommit:       "unknown",
			expectBuildTime: false,
			expectGitCommit: false,
		},
		{
			name:            "version with build info",
			version:         "1.0.0",
			buildTime:       "2024-01-01",
			gitCommit:       "abc123",
			expectBuildTime: true,
			expectGitCommit: true,
		},
		{
			name:            "version with only build time",
			version:         "1.0.0",
			buildTime:       "2024-01-01",
			gitCommit:       "unknown",
			expectBuildTime: true,
			expectGitCommit: false,
		},
		{
			name:            "version with only git commit",
			version:         "1.0.0",
			buildTime:       "unknown",
			gitCommit:       "abc123",
			expectBuildTime: false,
			expectGitCommit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Version = tt.version
			BuildTime = tt.buildTime
			GitCommit = tt.gitCommit

			var buf bytes.Buffer
			printVersion(&buf)
			output := buf.String()

			if !strings.Contains(output, tt.version) {
				t.Errorf("Expected output to contain version %q, got %q", tt.version, output)
			}

			if tt.expectBuildTime {
				if !strings.Contains(output, "Build time") || !strings.Contains(output, tt.buildTime) {
					t.Errorf("Expected output to contain build time %q, got %q", tt.buildTime, output)
				}
			}

			if tt.expectGitCommit {
				if !strings.Contains(output, "Git commit") || !strings.Contains(output, tt.gitCommit) {
					t.Errorf("Expected output to contain git commit %q, got %q", tt.gitCommit, output)
				}
			}
		})
	}
}

func TestPrintHelpTo(t *testing.T) {
	var buf bytes.Buffer
	PrintHelpTo(&buf)
	output := buf.String()

	expectedStrings := []string{
		"fmr (FileMirror)",
		"USAGE:",
		"OPTIONS:",
		"-p, --path",
		"-h, --help",
		"-v, --version",
		"KEYBOARD SHORTCUTS:",
		"EXAMPLES:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help output to contain %q, but it wasn't found", expected)
		}
	}
}

func TestRunWithArgs(t *testing.T) {
	// Save and restore current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir)

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		expectStdout []string
		expectStderr []string
	}{
		{
			name:         "help flag",
			args:         []string{"-h"},
			wantExitCode: 0,
			expectStdout: []string{"fmr (FileMirror)", "USAGE:"},
		},
		{
			name:         "version flag",
			args:         []string{"-v"},
			wantExitCode: 0,
			expectStdout: []string{"fmr version"},
		},
		{
			name:         "invalid path flag",
			args:         []string{"-p"},
			wantExitCode: 1,
			expectStderr: []string{"--path requires a directory argument"},
		},
		{
			name:         "nonexistent directory",
			args:         []string{"-p", "/nonexistent/dir"},
			wantExitCode: 1,
			expectStderr: []string{"cannot change to directory"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			exitCode := RunWithArgs(tt.args, &stdout, &stderr)

			if exitCode != tt.wantExitCode {
				t.Errorf("Exit code = %d, want %d", exitCode, tt.wantExitCode)
			}

			stdoutStr := stdout.String()
			for _, expected := range tt.expectStdout {
				if !strings.Contains(stdoutStr, expected) {
					t.Errorf("Expected stdout to contain %q, got %q", expected, stdoutStr)
				}
			}

			stderrStr := stderr.String()
			for _, expected := range tt.expectStderr {
				if !strings.Contains(stderrStr, expected) {
					t.Errorf("Expected stderr to contain %q, got %q", expected, stderrStr)
				}
			}

			// Restore original directory after each test
			os.Chdir(origDir)
		})
	}
}
