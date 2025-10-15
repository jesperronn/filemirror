# FileMirror Dual Binary Implementation

## Summary

Successfully implemented dual binary support for FileMirror, providing both `filemirror` and `fmr` binaries that are functionally identical. Users can install either binary based on their preference.

## Changes Made

### 1. Project Restructure
- **Changed root package from `main` to `filemirror`** in all source files:
  - `model.go`, `scanner.go`, `fileops.go`, `git.go`
  - All test files: `model_test.go`, `scanner_test.go`, `fileops_test.go`, `main_test.go`

- **Created new `run.go`** in the root package containing:
  - `Run()` - Main entry point function (exported)
  - `PrintHelp()` - Help message function (exported)
  - `Version`, `BuildTime`, `GitCommit` - Build-time variables (exported)

- **Created `cmd/` structure**:
  ```
  cmd/
  ├── filemirror/
  │   ├── main.go       # Entry point that calls filemirror.Run()
  │   └── main_test.go  # Tests for filemirror binary
  └── fmr/
      ├── main.go       # Entry point that calls filemirror.Run()
      └── main_test.go  # Tests for fmr binary
  ```

### 2. CI/CD Updates

#### GitHub Actions CI (`.github/workflows/ci.yml`)
- Updated build job to build both `filemirror` and `fmr` binaries
- Updated artifact naming to include both binaries
- Build commands now target `./cmd/filemirror` and `./cmd/fmr`

#### GitHub Actions Release (`.github/workflows/release.yml`)
- Updated to build both binaries for all platforms:
  - Linux (AMD64, ARM64)
  - macOS (Intel, Apple Silicon)
  - Windows (AMD64)
- Updated ldflags to reference `github.com/jesperronn/filemirror-fmr.Version` instead of `main.version`
- Updated release notes with installation instructions for both binaries
- Updated asset list to include all 10 binaries (5 platforms × 2 binaries)

### 3. Build Script Updates (`bin/build`)
- Added support for building both binaries simultaneously
- Added `--binary` flag to build specific binary (`filemirror` or `fmr`)
- Updated ldflags to use correct package path
- Enhanced loop to build each binary separately
- Improved output to show which binary is being built
- Installation now handles both binaries

### 4. Documentation Updates (`README.md`)
- Added note explaining both binaries are functionally identical
- Updated installation instructions for both:
  - `go install github.com/jesperronn/filemirror-fmr/cmd/fmr@latest`
  - `go install github.com/jesperronn/filemirror-fmr/cmd/filemirror@latest`
- Updated build script examples to show both options
- Added examples for building specific binaries

### 5. Testing
- Created test files for both binaries in `cmd/filemirror/main_test.go` and `cmd/fmr/main_test.go`
- Tests verify:
  - Binary can be built successfully
  - `--version` flag works correctly
  - `--help` flag works correctly
- All tests pass: `go test ./... -v`

## Installation

### For Users

**Using go install (Recommended):**
```bash
# Install the shorter fmr binary
go install github.com/jesperronn/filemirror-fmr/cmd/fmr@latest

# Or install the full filemirror binary
go install github.com/jesperronn/filemirror-fmr/cmd/filemirror@latest
```

**From source:**
```bash
git clone https://github.com/jesperronn/filemirror-fmr
cd filemirror
bin/build --install
```

### For Developers

**Build both binaries:**
```bash
bin/build
# or
go build -o filemirror ./cmd/filemirror
go build -o fmr ./cmd/fmr
```

**Build specific binary:**
```bash
bin/build --binary fmr
bin/build --binary filemirror
```

**Run tests:**
```bash
go test ./...
```

## Technical Details

### Package Structure
- **Root package (`filemirror`)**: Contains all core logic, models, and utilities
- **cmd/filemirror**: Thin wrapper that imports and calls `filemirror.Run()`
- **cmd/fmr**: Identical wrapper with shorter name

### Build Configuration
- Both binaries use the same ldflags for versioning
- Package path for version injection: `github.com/jesperronn/filemirror-fmr`
- Exported variables: `Version`, `BuildTime`, `GitCommit`

### Benefits
✅ Clean architecture following Go best practices
✅ No shell configuration needed
✅ Works seamlessly with `go install`
✅ Minimal code duplication (just 7 lines per binary)
✅ Easy to maintain
✅ Both binaries included in releases
✅ Users can choose their preferred name

## Files Changed
- Modified: `model.go`, `scanner.go`, `fileops.go`, `git.go` (package rename)
- Modified: All test files (package rename)
- Created: `run.go` (new main logic)
- Created: `cmd/filemirror/main.go`, `cmd/filemirror/main_test.go`
- Created: `cmd/fmr/main.go`, `cmd/fmr/main_test.go`
- Modified: `.github/workflows/ci.yml`
- Modified: `.github/workflows/release.yml`
- Modified: `bin/build`
- Modified: `README.md`
- Removed: `main.go` (logic moved to `run.go`)

## Verification

All functionality verified:
- ✅ Both binaries build successfully
- ✅ Both binaries show correct version info
- ✅ Both binaries display help text
- ✅ All tests pass
- ✅ Build script handles both binaries
- ✅ CI/CD workflows updated
- ✅ Documentation updated
