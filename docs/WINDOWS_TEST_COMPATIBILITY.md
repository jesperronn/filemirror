# Windows Test Compatibility

This document describes the Windows compatibility considerations for the test suite.

## Tests That Skip on Windows

### File Permission Tests

**Test**: `TestCopyFilePreservesPermissions` (fileops_test.go)

**Reason**: Windows handles file permissions differently than Unix-like systems. Unix permission bits (like 0o755) don't translate directly to Windows ACLs (Access Control Lists).

**Status**: ✅ Automatically skipped on Windows with:
```go
if runtime.GOOS == "windows" {
    t.Skip("Skipping permissions test on Windows (file permissions work differently)")
}
```

## Tests That Should Work on Windows

### Git Tests (git_test.go)

**Requirements**: Git must be installed and available in PATH

**Status**: ✅ Should work on Windows
- All git commands use `git` directly (no shell-specific features)
- File paths use `filepath.Join()` which handles Windows path separators
- Temporary directories are cross-platform

**Tests**:
- `TestGetGitBranch` - Tests git branch detection
- `TestGetGitBranchDetachedHead` - Tests detached HEAD state
- `TestGetGitBranchEmptyRepo` - Tests empty repository handling

### File Operations Tests (fileops_test.go)

**Status**: ✅ Should work on Windows

All file operation tests use cross-platform APIs:
- `os.WriteFile()` - Cross-platform
- `os.MkdirTemp()` - Cross-platform
- `filepath.Join()` - Handles Windows path separators
- Octal permissions (0o644, 0o755) - Ignored on Windows but don't cause errors

**Tests**:
- `TestCopyFile` - Basic file copying
- `TestCopyFileNonExistentSource` - Error handling
- `TestCopyFileToNonExistentDirectory` - Error handling
- `TestCopyFileWithLargeContent` - Large file handling (1MB)
- `TestCopyFileEmptyFile` - Empty file handling
- `TestCopyFileOverwritesExisting` - Overwrite behavior

### Scanner Tests (scanner_test.go)

**Status**: ✅ Should work on Windows

Scanner tests use cross-platform directory traversal:
- `filepath.WalkDir()` - Cross-platform
- `os.PathSeparator` - Platform-specific separator
- Directory exclusions work on all platforms

**Tests**:
- `TestScanFiles` - Basic file scanning
- `TestScanFilesExcludesNodeModules` - Directory exclusion
- `TestScanFilesDeepDirectory` - Deep directory handling
- `TestScanFilesAllExcludedDirs` - All excluded directories
- `TestScanFilesEmptyDirectory` - Empty directory handling
- `TestScanFilesWithComplexPattern` - Pattern matching
- `TestScanFilesSorted` - File sorting
- `TestScanFilesRelativePaths` - Relative path handling

### GitWorkflow Tests (gitworkflow_test.go)

**Status**: ✅ Should work on Windows

**Requirements**: Git must be installed

All git operations use cross-platform APIs and git commands:
- Worktree creation uses git commands
- File paths use `filepath` package
- Temporary directories are cross-platform

**Tests**:
- `TestDetectGitRoot` - Git root detection
- `TestGenerateWorktreeID` - Worktree ID generation
- `TestCreateWorktreeAndBranch` - Worktree creation
- `TestCommitChanges` - Git commit operations
- `TestCleanupWorktree` - Worktree cleanup
- `TestCopyFileToWorktree` - File copying to worktree

### Model Tests (model_test.go)

**Status**: ✅ Should work on Windows

All model tests are platform-independent:
- String operations
- Pattern matching
- Data structure manipulation

## CI/CD Considerations

### GitHub Actions

If running tests in GitHub Actions on Windows, ensure:

1. **Git is available**: GitHub Actions Windows runners come with Git pre-installed
2. **Go is configured**: Use `actions/setup-go@v4` 
3. **Test command**: Use `go test ./...` (works on all platforms)

### Example GitHub Actions Configuration

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}
      - run: go test -v ./...
      - run: go test -race ./...
```

## Known Platform Differences

### File Paths

✅ **Handled automatically** by Go's `filepath` package:
- Unix: `/home/user/file.txt`
- Windows: `C:\Users\user\file.txt`

### Path Separators

✅ **Handled automatically**:
- `filepath.Join()` uses correct separator
- `os.PathSeparator` provides platform-specific separator

### Temporary Directories

✅ **Handled automatically** by `os.MkdirTemp()`:
- Unix: `/tmp/...`
- Windows: `C:\Users\...\AppData\Local\Temp\...`

### File Permissions

⚠️ **Different behavior**:
- Unix: Octal permissions (0o644, 0o755) are respected
- Windows: Permissions are ignored (ACLs used instead)
- **Solution**: Skip permission-specific tests on Windows

### Git Behavior

✅ **Generally consistent** across platforms:
- Branch detection works the same
- Commit operations work the same
- Worktrees work the same
- Line endings may differ (CRLF vs LF) but git handles this

## Troubleshooting

### Git Not Found on Windows

**Error**: `exec: "git": executable file not found in %PATH%`

**Solutions**:
1. Install Git for Windows: https://git-scm.com/download/win
2. Ensure git is in PATH
3. Restart terminal/IDE after installation

### Permission Denied Errors

**Error**: Permission denied when creating/modifying files

**Solutions**:
1. Run tests as administrator (usually not necessary)
2. Check antivirus isn't blocking file operations
3. Ensure test files are in a writable directory

### Path Length Issues (Windows)

**Error**: Path length exceeds 260 characters

**Solutions**:
1. Enable long paths: `git config --system core.longpaths true`
2. Or use shorter temporary directory paths
3. Tests use `t.TempDir()` which should handle this

## Summary

✅ **24 of 25 tests** work cross-platform
⚠️ **1 test** skipped on Windows (file permissions)

The test suite is designed to be cross-platform compatible and should run successfully on Windows, macOS, and Linux with Git installed.
