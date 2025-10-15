package filemirror

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	Path     string
	Size     int64
	Modified time.Time
	Branch   string
}

var excludeDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	".next":        true,
	"dist":         true,
	"build":        true,
	"target":       true,
	".cache":       true,
}

func scanFiles(workDir, pattern string) ([]FileInfo, error) {
	var files []FileInfo
	maxDepth := 4

	// Use workDir as the base directory
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Make sure workDir is absolute
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	err = filepath.WalkDir(absWorkDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files/dirs we can't access
		}

		// Calculate relative depth from the work directory
		relPath, err := filepath.Rel(absWorkDir, path)
		if err != nil {
			return nil
		}
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > maxDepth {
			return fs.SkipDir
		}

		// Skip excluded directories
		if d.IsDir() {
			if excludeDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}

		// Skip if not a regular file
		if !d.Type().IsRegular() {
			return nil
		}

		// Filter by pattern if provided
		if pattern != "" && !matchesPattern(d.Name(), pattern) {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil // Skip files we can't stat
		}

		// Get git branch for this file
		branch := getGitBranch(path)

		// Store relative path from work directory
		files = append(files, FileInfo{
			Path:     relPath,
			Size:     info.Size(),
			Modified: info.ModTime(),
			Branch:   branch,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	// Sort by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Modified.After(files[j].Modified)
	})

	return files, nil
}

func matchesPattern(filename, pattern string) bool {
	// Simple pattern matching
	pattern = strings.ToLower(pattern)
	filename = strings.ToLower(filename)

	// If pattern contains *, use glob matching
	if strings.Contains(pattern, "*") {
		matched, err := filepath.Match(pattern, filename)
		return err == nil && matched
	}

	// Otherwise, simple substring match
	return strings.Contains(filename, pattern)
}
