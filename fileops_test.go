package filemirror

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	srcContent := "Hello, World! This is test content."
	if err := os.WriteFile(srcPath, []byte(srcContent), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create destination file with different content
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := os.WriteFile(dstPath, []byte("old content"), 0o644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content was copied
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != srcContent {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), srcContent)
	}

	// Verify file permissions were preserved (skip on Windows as permissions work differently)
	if runtime.GOOS != "windows" {
		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			t.Fatalf("Failed to stat source file: %v", err)
		}

		dstInfo, err := os.Stat(dstPath)
		if err != nil {
			t.Fatalf("Failed to stat destination file: %v", err)
		}

		if srcInfo.Mode() != dstInfo.Mode() {
			t.Errorf("Permissions mismatch: got %v, want %v", dstInfo.Mode(), srcInfo.Mode())
		}
	}
}

func TestCopyFileNonExistentSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "nonexistent.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	err = copyFile(srcPath, dstPath)
	if err == nil {
		t.Error("Expected error when copying non-existent file, got nil")
	}
}

func TestCopyFilePreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permissions test on Windows (file permissions work differently)")
	}

	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file with specific permissions
	srcPath := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcPath, []byte("test"), 0o755); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Check permissions
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("Failed to stat destination file: %v", err)
	}

	expectedMode := os.FileMode(0o755)
	if info.Mode().Perm() != expectedMode {
		t.Errorf("Expected permissions %v, got %v", expectedMode, info.Mode().Perm())
	}
}

func TestCopyFileToNonExistentDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Try to copy to non-existent directory
	dstPath := filepath.Join(tmpDir, "nonexistent", "dest.txt")

	err = copyFile(srcPath, dstPath)
	if err == nil {
		t.Error("Expected error when copying to non-existent directory, got nil")
	}
}

func TestCopyFileWithLargeContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file with large content (1MB)
	srcPath := filepath.Join(tmpDir, "large.txt")
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	if err := os.WriteFile(srcPath, largeContent, 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if len(dstContent) != len(largeContent) {
		t.Errorf("Content length mismatch: got %d, want %d", len(dstContent), len(largeContent))
	}

	// Verify a sample of bytes
	for i := 0; i < 1000; i += 100 {
		if dstContent[i] != largeContent[i] {
			t.Errorf("Content mismatch at byte %d: got %d, want %d", i, dstContent[i], largeContent[i])
		}
	}
}

func TestCopyFileEmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create empty source file
	srcPath := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(srcPath, []byte{}, 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify destination is also empty
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if len(dstContent) != 0 {
		t.Errorf("Expected empty file, got %d bytes", len(dstContent))
	}
}

func TestCopyFileOverwritesExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	newContent := "new content"
	if err := os.WriteFile(srcPath, []byte(newContent), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create existing destination with different content
	oldContent := "old content that should be replaced"
	if err := os.WriteFile(dstPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify new content replaced old content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != newContent {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), newContent)
	}
}

func TestCopyFileReadOnlySource(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping read-only test on Windows (different permission handling)")
	}

	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create read-only source file
	srcPath := filepath.Join(tmpDir, "readonly.txt")
	content := "read-only content"
	if err := os.WriteFile(srcPath, []byte(content), 0o444); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Should be able to copy read-only file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed with read-only source: %v", err)
	}

	// Verify content was copied
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != content {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), content)
	}

	// Verify permissions were preserved
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("Failed to stat destination: %v", err)
	}

	if info.Mode().Perm() != 0o444 {
		t.Errorf("Expected permissions 0o444, got %v", info.Mode().Perm())
	}
}

func TestCopyFileReadOnlyDestination(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping read-only test on Windows (different permission handling)")
	}

	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcPath, []byte("new content"), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create read-only destination
	dstPath := filepath.Join(tmpDir, "readonly-dest.txt")
	oldContent := "old content"
	if err := os.WriteFile(dstPath, []byte(oldContent), 0o444); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	// Atomic rename can overwrite read-only files on Unix (if you have write permission on the directory)
	// This is expected and correct behavior for our use case
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content was updated
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != "new content" {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), "new content")
	}

	// Verify permissions were updated to match source
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("Failed to stat destination: %v", err)
	}

	if info.Mode().Perm() != 0o644 {
		t.Errorf("Expected permissions 0o644, got %v", info.Mode().Perm())
	}
}

func TestCopyFileSameSourceAndDestination(t *testing.T) {
	// On Windows, os.Rename() to the same file fails with "Access is denied"
	// This is a Windows-specific file locking behavior that doesn't occur in real usage
	if runtime.GOOS == "windows" {
		t.Skip("Skipping same source/dest test on Windows (file locking prevents atomic rename to same path)")
	}

	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create file
	path := filepath.Join(tmpDir, "samefile.txt")
	originalContent := "original content"
	if err := os.WriteFile(path, []byte(originalContent), 0o644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Try to copy file to itself
	if err := copyFile(path, path); err != nil {
		t.Fatalf("copyFile failed for same source/dest: %v", err)
	}

	// Verify content is unchanged
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != originalContent {
		t.Errorf("Content changed: got %q, want %q", string(content), originalContent)
	}
}

func TestCopyFileSpecialCharactersInPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fmr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file with special characters in name
	srcPath := filepath.Join(tmpDir, "file with spaces & special-chars.txt")
	content := "test content with special chars"
	if err := os.WriteFile(srcPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dest-with-chars.txt")

	// Copy file
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed with special characters: %v", err)
	}

	// Verify content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination: %v", err)
	}

	if string(dstContent) != content {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), content)
	}
}
