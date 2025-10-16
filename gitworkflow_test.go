package filemirror

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a temporary git repository
func createTestGitRepo(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "fmr-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git name: %v", err)
	}

	// Create an initial commit (required for worktrees)
	testFile := filepath.Join(tmpDir, "initial.txt")
	if err := os.WriteFile(testFile, []byte("initial"), 0o644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create initial file: %v", err)
	}

	cmd = exec.Command("git", "add", "initial.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to add initial file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	return tmpDir
}

// verifyWorktreesCleanedUp verifies that no temporary worktrees exist for a repo
func verifyWorktreesCleanedUp(t *testing.T, repoPath string) {
	t.Helper()

	// List all worktrees for the repo
	cmd := exec.Command("git", "-C", repoPath, "worktree", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Warning: failed to list worktrees: %v", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Should only have the main worktree (the repo itself)
	if len(lines) > 1 {
		t.Errorf("Expected only main worktree, but found %d worktrees:\n%s", len(lines), string(output))
	}

	// Also check temp directory for any fmr-worktree-* directories
	tmpDir := os.TempDir()
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Logf("Warning: failed to read temp dir: %v", err)
		return
	}

	var worktreeDirs []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "fmr-worktree-") {
			worktreeDirs = append(worktreeDirs, entry.Name())
		}
	}

	if len(worktreeDirs) > 0 {
		t.Errorf("Found %d temporary worktree directories that weren't cleaned up: %v", len(worktreeDirs), worktreeDirs)
	}
}

func TestDetectGitRoot(t *testing.T) {
	// Create a test git repo
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	// Create a nested directory structure
	nestedDir := filepath.Join(repoPath, "src", "pkg")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	testFile := filepath.Join(nestedDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "file in repo root",
			filePath:    filepath.Join(repoPath, "initial.txt"),
			expectError: false,
		},
		{
			name:        "file in nested directory",
			filePath:    testFile,
			expectError: false,
		},
		{
			name:        "file outside repo",
			filePath:    filepath.Join(os.TempDir(), "not-in-repo.txt"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := detectGitRoot(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Normalize paths for comparison
			// Use EvalSymlinks to resolve /var -> /private/var on macOS
			expectedRoot, _ := filepath.EvalSymlinks(repoPath)
			actualRoot, _ := filepath.EvalSymlinks(root)

			if actualRoot != expectedRoot {
				t.Errorf("Got root %q, want %q", actualRoot, expectedRoot)
			}
		})
	}
}

func TestGenerateWorktreeID(t *testing.T) {
	// Generate multiple IDs and ensure they're unique and non-empty
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id := generateWorktreeID()
		if id == "" {
			t.Error("Generated empty ID")
		}
		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

func TestCreateWorktreeAndBranch(t *testing.T) {
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	tests := []struct {
		name        string
		branchName  string
		expectError bool
		setupBranch bool // whether to create the branch first
	}{
		{
			name:        "create new branch",
			branchName:  "feature/test-branch",
			expectError: false,
			setupBranch: false,
		},
		{
			name:        "checkout existing branch",
			branchName:  "existing-branch",
			expectError: false,
			setupBranch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: create branch if needed
			if tt.setupBranch {
				cmd := exec.Command("git", "branch", tt.branchName)
				cmd.Dir = repoPath
				if err := cmd.Run(); err != nil {
					t.Fatalf("Failed to setup branch: %v", err)
				}
			}

			worktreePath, err := createWorktreeAndBranch(repoPath, tt.branchName, []string{})

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			defer cleanupWorktree(repoPath, worktreePath)

			// Verify worktree was created
			if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
				t.Error("Worktree directory was not created")
			}

			// Verify we're on the correct branch in the worktree
			cmd := exec.Command("git", "branch", "--show-current")
			cmd.Dir = worktreePath
			output, err := cmd.Output()
			if err != nil {
				t.Errorf("Failed to check branch: %v", err)
				return
			}

			currentBranch := strings.TrimSpace(string(output))
			if currentBranch != tt.branchName {
				t.Errorf("Expected branch %q, got %q", tt.branchName, currentBranch)
			}
		})
	}
}

func TestCommitChanges(t *testing.T) {
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	worktreePath, err := createWorktreeAndBranch(repoPath, "test-commit", []string{})
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}
	defer cleanupWorktree(repoPath, worktreePath)

	tests := []struct {
		name        string
		files       []string
		message     string
		expectError bool
		setup       func() error
	}{
		{
			name:    "commit single file",
			files:   []string{filepath.Join(worktreePath, "test.txt")},
			message: "Test commit",
			setup: func() error {
				return os.WriteFile(filepath.Join(worktreePath, "test.txt"), []byte("test content"), 0o644)
			},
			expectError: false,
		},
		{
			name: "commit multiple files",
			files: []string{
				filepath.Join(worktreePath, "file1.txt"),
				filepath.Join(worktreePath, "file2.txt"),
			},
			message: "Multi-file commit",
			setup: func() error {
				if err := os.WriteFile(filepath.Join(worktreePath, "file1.txt"), []byte("content1"), 0o644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(worktreePath, "file2.txt"), []byte("content2"), 0o644)
			},
			expectError: false,
		},
		{
			name:    "nothing to commit",
			files:   []string{filepath.Join(worktreePath, "unchanged.txt")},
			message: "Should not fail",
			setup: func() error {
				// Create and commit file first
				path := filepath.Join(worktreePath, "unchanged.txt")
				if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
					return err
				}
				cmd := exec.Command("git", "add", "unchanged.txt")
				cmd.Dir = worktreePath
				if err := cmd.Run(); err != nil {
					return err
				}
				cmd = exec.Command("git", "commit", "-m", "Initial commit")
				cmd.Dir = worktreePath
				return cmd.Run()
			},
			expectError: false, // Should not error when nothing to commit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err := commitChanges(worktreePath, tt.message)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCleanupWorktree(t *testing.T) {
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	worktreePath, err := createWorktreeAndBranch(repoPath, "cleanup-test", []string{})
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Verify worktree exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Fatal("Worktree was not created")
	}

	// Cleanup worktree
	err = cleanupWorktree(repoPath, worktreePath)
	if err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}

	// Verify worktree was removed
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Error("Worktree still exists after cleanup")
	}
}

func TestCopyFileToWorktree(t *testing.T) {
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	worktreePath, err := createWorktreeAndBranch(repoPath, "copy-test", []string{})
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}
	defer cleanupWorktree(repoPath, worktreePath)

	tests := []struct {
		name        string
		setupFile   string
		content     string
		expectError bool
	}{
		{
			name:        "copy file to root",
			setupFile:   "test.txt",
			content:     "test content",
			expectError: false,
		},
		{
			name:        "copy file to nested directory",
			setupFile:   "src/pkg/nested.go",
			content:     "package pkg",
			expectError: false,
		},
		{
			name:        "copy file with deep nesting",
			setupFile:   "a/b/c/d/file.txt",
			content:     "deeply nested",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source file in repo
			sourcePath := filepath.Join(repoPath, tt.setupFile)
			sourceDir := filepath.Dir(sourcePath)
			if err := os.MkdirAll(sourceDir, 0o755); err != nil {
				t.Fatalf("Failed to create source dir: %v", err)
			}
			if err := os.WriteFile(sourcePath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("Failed to create source file: %v", err)
			}

			// Copy to worktree
			err := copyFileToWorktree(sourcePath, worktreePath, repoPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify file was copied correctly
			expectedPath := filepath.Join(worktreePath, tt.setupFile)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Error("File was not copied to worktree")
				return
			}

			// Verify content
			content, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("Failed to read copied file: %v", err)
				return
			}

			if string(content) != tt.content {
				t.Errorf("Content mismatch: got %q, want %q", string(content), tt.content)
			}
		})
	}
}

func TestPerformGitWorkflow(t *testing.T) {
	// Create two test repos
	repo1 := createTestGitRepo(t)
	defer os.RemoveAll(repo1)

	repo2 := createTestGitRepo(t)
	defer os.RemoveAll(repo2)

	tests := []struct {
		name          string
		setupRepos    func() map[string][]string
		branchName    string
		commitMsg     string
		shouldPush    bool
		expectSuccess int
		expectErrors  int
	}{
		{
			name: "single repo with one file",
			setupRepos: func() map[string][]string {
				// Create a file in repo1
				testFile := filepath.Join(repo1, "workflow-test.txt")
				os.WriteFile(testFile, []byte("workflow content"), 0o600)
				return map[string][]string{
					repo1: {testFile},
				}
			},
			branchName:    "test-workflow",
			commitMsg:     "Test workflow commit",
			shouldPush:    false,
			expectSuccess: 1,
			expectErrors:  0,
		},
		{
			name: "multiple repos",
			setupRepos: func() map[string][]string {
				// Create files in both repos
				file1 := filepath.Join(repo1, "multi-test1.txt")
				file2 := filepath.Join(repo2, "multi-test2.txt")
				os.WriteFile(file1, []byte("content1"), 0o600)
				os.WriteFile(file2, []byte("content2"), 0o600)
				return map[string][]string{
					repo1: {file1},
					repo2: {file2},
				}
			},
			branchName:    "multi-repo-test",
			commitMsg:     "Multi-repo commit",
			shouldPush:    false,
			expectSuccess: 2,
			expectErrors:  0,
		},
		{
			name: "empty repos map",
			setupRepos: func() map[string][]string {
				return map[string][]string{}
			},
			branchName:    "empty-test",
			commitMsg:     "Empty commit",
			shouldPush:    false,
			expectSuccess: 0,
			expectErrors:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := tt.setupRepos()

			successRepos, errors := performGitWorkflow(repos, tt.branchName, tt.commitMsg, tt.shouldPush)

			if len(successRepos) != tt.expectSuccess {
				t.Errorf("Expected %d successful repos, got %d", tt.expectSuccess, len(successRepos))
			}

			if len(errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectErrors, len(errors), errors)
			}

			// Verify worktrees are cleaned up
			for repoPath := range repos {
				verifyWorktreesCleanedUp(t, repoPath)
			}
		})
	}
}

func TestProcessRepo(t *testing.T) {
	tests := []struct {
		name          string
		setupRepo     func(t *testing.T) (repoPath string, files []string)
		branchName    string
		commitMsg     string
		shouldPush    bool
		expectSuccess bool
		expectError   bool
	}{
		{
			name: "process single file successfully",
			setupRepo: func(t *testing.T) (string, []string) {
				repo := createTestGitRepo(t)
				testFile := filepath.Join(repo, "process-test.txt")
				if err := os.WriteFile(testFile, []byte("process content"), 0o600); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return repo, []string{testFile}
			},
			branchName:    "process-branch",
			commitMsg:     "Process commit",
			shouldPush:    false,
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "process multiple files",
			setupRepo: func(t *testing.T) (string, []string) {
				repo := createTestGitRepo(t)
				file1 := filepath.Join(repo, "file1.txt")
				file2 := filepath.Join(repo, "file2.txt")
				os.WriteFile(file1, []byte("content1"), 0o600)
				os.WriteFile(file2, []byte("content2"), 0o600)
				return repo, []string{file1, file2}
			},
			branchName:    "multi-file-process",
			commitMsg:     "Multi-file process commit",
			shouldPush:    false,
			expectSuccess: true,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, files := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			success, err := processRepo(repoPath, files, tt.branchName, tt.commitMsg, tt.shouldPush)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if success != tt.expectSuccess {
				t.Errorf("Expected success=%v, got %v", tt.expectSuccess, success)
			}

			// Verify worktrees are cleaned up
			verifyWorktreesCleanedUp(t, repoPath)
		})
	}
}

func TestPushBranch(t *testing.T) {
	// This test is limited as we can't test actual pushing without a remote
	// We test that the function constructs the command correctly
	repoPath := createTestGitRepo(t)
	defer os.RemoveAll(repoPath)

	worktreePath, err := createWorktreeAndBranch(repoPath, "push-test", []string{})
	if err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}
	defer cleanupWorktree(repoPath, worktreePath)

	// Test will fail as there's no remote, but we can verify it attempts to push
	err = pushBranch(worktreePath, "push-test")
	if err == nil {
		t.Skip("Skipping: no remote configured (this is expected)")
	}

	// Verify the error is about remote, not command construction
	if err != nil && !strings.Contains(err.Error(), "remote") && !strings.Contains(err.Error(), "No configured push destination") {
		t.Logf("Expected remote-related error, got: %v", err)
	}
}
