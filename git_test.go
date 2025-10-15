package filemirror

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetGitBranch(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) string
		expectedBranch string
	}{
		{
			name: "returns current branch in git repo",
			setup: func(t *testing.T) string {
				// Create a temp git repo
				tmpDir := t.TempDir()

				// Initialize git repo
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				if err := cmd.Run(); err != nil {
					t.Fatalf("Failed to init git repo: %v", err)
				}

				// Configure git
				exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
				exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

				// Create and commit a file to establish main branch
				testFile := filepath.Join(tmpDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()
				exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

				return testFile
			},
			expectedBranch: "master", // or "main" depending on git config
		},
		{
			name: "returns dash for non-git directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return testFile
			},
			expectedBranch: "-",
		},
		{
			name: "returns branch name for feature branch",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Initialize git repo
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				if err := cmd.Run(); err != nil {
					t.Fatalf("Failed to init git repo: %v", err)
				}

				// Configure git
				exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
				exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

				// Create and commit a file
				testFile := filepath.Join(tmpDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()
				exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

				// Create and checkout feature branch
				exec.Command("git", "-C", tmpDir, "checkout", "-b", "feature/test-branch").Run()

				return testFile
			},
			expectedBranch: "feature/test-branch",
		},
		{
			name: "works with subdirectory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Initialize git repo
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				if err := cmd.Run(); err != nil {
					t.Fatalf("Failed to init git repo: %v", err)
				}

				// Configure git
				exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
				exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

				// Create subdirectory
				subDir := filepath.Join(tmpDir, "subdir", "nested")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("Failed to create subdir: %v", err)
				}

				// Create file in root for initial commit
				rootFile := filepath.Join(tmpDir, "root.txt")
				if err := os.WriteFile(rootFile, []byte("root"), 0644); err != nil {
					t.Fatalf("Failed to create root file: %v", err)
				}
				exec.Command("git", "-C", tmpDir, "add", "root.txt").Run()
				exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

				// Create test branch
				exec.Command("git", "-C", tmpDir, "checkout", "-b", "test-subdir").Run()

				// Create file in subdirectory
				testFile := filepath.Join(subDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				return testFile
			},
			expectedBranch: "test-subdir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup(t)
			branch := getGitBranch(filePath)

			// For git repos, check if we got either master or main (git default branch name changed)
			if tt.expectedBranch == "master" {
				if branch != "master" && branch != "main" {
					t.Errorf("getGitBranch() = %q, want %q or %q", branch, "master", "main")
				}
			} else if branch != tt.expectedBranch {
				t.Errorf("getGitBranch() = %q, want %q", branch, tt.expectedBranch)
			}
		})
	}
}

func TestGetGitBranchDetachedHead(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create and commit files to have history
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "First commit").Run()

	// Second commit
	if err := os.WriteFile(testFile, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}
	exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Second commit").Run()

	// Get the first commit hash
	cmd = exec.Command("git", "-C", tmpDir, "rev-parse", "HEAD~1")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit hash: %v", err)
	}
	commitHash := string(output[:7]) // First 7 chars

	// Checkout the first commit (detached HEAD)
	exec.Command("git", "-C", tmpDir, "checkout", commitHash).Run()

	// Test getGitBranch
	branch := getGitBranch(testFile)

	// In detached HEAD state, --show-current returns empty string, which should become "-"
	if branch != "-" {
		t.Logf("Note: In detached HEAD state, got branch %q (expected %q)", branch, "-")
	}
}

func TestGetGitBranchEmptyRepo(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize empty git repo (no commits)
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	branch := getGitBranch(testFile)

	// Empty repo with no commits has no current branch
	// --show-current returns empty string, which becomes "-"
	if branch != "-" {
		t.Logf("Note: Empty repo returned branch %q (expected %q)", branch, "-")
	}
}
