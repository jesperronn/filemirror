package filemirror

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// detectGitRoot finds the git repository root for a file
func detectGitRoot(filePath string) (string, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(absPath)
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// groupFilesByRepo groups changed files by their git repository
func groupFilesByRepo(files []string) (map[string][]string, error) {
	repos := make(map[string][]string)

	for _, file := range files {
		root, err := detectGitRoot(file)
		if err != nil {
			// File not in git repo, skip
			continue
		}

		// Store absolute path of the file
		absPath, err := filepath.Abs(file)
		if err != nil {
			continue
		}

		repos[root] = append(repos[root], absPath)
	}

	return repos, nil
}

// generateWorktreeID generates a random ID for worktree paths
func generateWorktreeID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// createWorktreeAndBranch creates a new git worktree with branch
func createWorktreeAndBranch(repoPath, branchName string) (string, error) {
	// Generate unique worktree path
	worktreePath := filepath.Join(os.TempDir(), fmt.Sprintf("fmr-worktree-%s", generateWorktreeID()))

	// Check if branch already exists
	checkCmd := exec.Command("git", "-C", repoPath, "rev-parse", "--verify", branchName)
	branchExists := checkCmd.Run() == nil

	var cmd *exec.Cmd
	if branchExists {
		// Branch exists, check it out in the worktree
		cmd = exec.Command("git", "-C", repoPath, "worktree", "add", worktreePath, branchName)
	} else {
		// Create new branch in worktree
		cmd = exec.Command("git", "-C", repoPath, "worktree", "add", "-b", branchName, worktreePath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create worktree: %w\n%s", err, string(output))
	}

	return worktreePath, nil
}

// commitChanges stages and commits files in the worktree
func commitChanges(worktreePath, message string, files []string) error {
	// Stage files
	for _, file := range files {
		// Get the path relative to the worktree
		relPath, err := filepath.Rel(worktreePath, file)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		addCmd := exec.Command("git", "-C", worktreePath, "add", relPath)
		if output, err := addCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to stage file %s: %w\n%s", relPath, err, string(output))
		}
	}

	// Commit changes
	commitCmd := exec.Command("git", "-C", worktreePath, "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		// Check if it's a "nothing to commit" error
		if strings.Contains(string(output), "nothing to commit") {
			return nil // Not an error, just nothing changed
		}
		return fmt.Errorf("failed to commit: %w\n%s", err, string(output))
	}

	return nil
}

// pushBranch pushes the branch to origin
func pushBranch(worktreePath, branchName string) error {
	cmd := exec.Command("git", "-C", worktreePath, "push", "-u", "origin", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push: %w\n%s", err, string(output))
	}

	return nil
}

// cleanupWorktree removes the worktree
func cleanupWorktree(repoPath, worktreePath string) error {
	// Remove worktree
	cmd := exec.Command("git", "-C", repoPath, "worktree", "remove", worktreePath, "--force")
	if output, err := cmd.CombinedOutput(); err != nil {
		// Try to remove directory manually if git worktree remove fails
		if rmErr := os.RemoveAll(worktreePath); rmErr != nil {
			return fmt.Errorf("failed to cleanup worktree: %w\n%s", err, string(output))
		}
	}

	return nil
}

// copyFileToWorktree copies a file to the worktree preserving its structure
func copyFileToWorktree(sourceFilePath, worktreePath, repoRoot string) error {
	// Get the path relative to the repo root
	relPath, err := filepath.Rel(repoRoot, sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Target path in worktree
	targetPath := filepath.Join(worktreePath, relPath)

	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	input, err := os.ReadFile(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}

	if err := os.WriteFile(targetPath, input, 0644); err != nil {
		return fmt.Errorf("failed to write target: %w", err)
	}

	return nil
}

// performGitWorkflow executes the complete git workflow for changed files
func performGitWorkflow(repos map[string][]string, branchName, commitMessage string, shouldPush bool) ([]string, []error) {
	var successRepos []string
	var errors []error

	for repoPath, files := range repos {
		// Create worktree
		worktreePath, err := createWorktreeAndBranch(repoPath, branchName)
		if err != nil {
			errors = append(errors, fmt.Errorf("repo %s: %w", repoPath, err))
			continue
		}

		// Ensure cleanup happens
		defer cleanupWorktree(repoPath, worktreePath)

		// Copy files to worktree
		for _, file := range files {
			if err := copyFileToWorktree(file, worktreePath, repoPath); err != nil {
				errors = append(errors, fmt.Errorf("repo %s: %w", repoPath, err))
				continue
			}
		}

		// Commit changes
		if err := commitChanges(worktreePath, commitMessage, files); err != nil {
			errors = append(errors, fmt.Errorf("repo %s: %w", repoPath, err))
			continue
		}

		// Push if requested
		if shouldPush {
			if err := pushBranch(worktreePath, branchName); err != nil {
				errors = append(errors, fmt.Errorf("repo %s (push failed): %w", repoPath, err))
				// Don't skip this repo, commit was successful
			}
		}

		successRepos = append(successRepos, repoPath)
	}

	return successRepos, errors
}
