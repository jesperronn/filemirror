package filemirror

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close source file: %w", closeErr)
		}
	}()

	// Get source file info to preserve permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create temporary file in the same directory as destination
	tmpFile, err := os.CreateTemp(filepath.Dir(dst), ".fmr-tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath) // Best effort cleanup, ignore error
	}()

	// Copy content
	_, err = io.Copy(tmpFile, sourceFile)
	if err != nil {
		_ = tmpFile.Close() // Best effort close, ignore error
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Close temp file before rename
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Set permissions to match source
	if err := os.Chmod(tmpPath, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, dst); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
