package main

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func getGitBranch(filePath string) string {
	// Get the directory containing the file
	dir := filepath.Dir(filePath)

	// Run git command to get current branch
	cmd := exec.Command("git", "-C", dir, "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		// Not in a git repo or git not available
		return "-"
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "-"
	}

	return branch
}
