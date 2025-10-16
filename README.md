# FileMirror (fmr)

[![CI](https://github.com/jesperronn/filemirror-fmr/workflows/CI/badge.svg)](https://github.com/jesperronn/filemirror-fmr/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jesperronn/filemirror-fmr)](https://goreportcard.com/report/github.com/jesperronn/filemirror-fmr)
[![codecov](https://codecov.io/gh/jesperronn/filemirror-fmr/branch/main/graph/badge.svg)](https://codecov.io/gh/jesperronn/filemirror-fmr)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jesperronn/filemirror-fmr)](https://github.com/jesperronn/filemirror-fmr)

FileMirror is a repo-aware CLI for mirroring a canonical file across multiple directories and repositories.

Features: interactive TUI, diff preview, per-repo worktree creation to stage & commit changes, atomic writes, and dry-run mode.

> **Note:** FileMirror provides two identical binaries: `filemirror` and `fmr`. Use whichever you prefer - `fmr` is just a shorter alias.

**Quickstart:**
```bash
fmr --path ~/projects "config.yaml"
```

Designed for multi-repo workflows where safe, auditable propagation of canonical files is required.

## Installation

### Option 1: Using go install (Recommended)

Install directly to your `$GOPATH/bin`:

```bash
# Install the shorter fmr binary (recommended)
go install github.com/jesperronn/filemirror-fmr/cmd/fmr@latest

# Or install the full filemirror binary
go install github.com/jesperronn/filemirror-fmr/cmd/filemirror@latest
```

Both binaries are functionally identical - `fmr` is just a shorter alias for `filemirror`.
This will make the binary available globally in your terminal.

### Option 2: From Source with Build Script

```bash
git clone https://github.com/jesperronn/filemirror-fmr
cd filemirror
bin/build --install
```

The `--install` flag builds and installs both binaries to `$GOPATH/bin`.
You can also build a specific binary:

```bash
bin/build --binary fmr --install     # Build only fmr
bin/build --binary filemirror --install  # Build only filemirror
```

### Option 3: Manual Build

```bash
git clone https://github.com/jesperronn/filemirror-fmr
cd filemirror

# Build both binaries
go build -o fmr ./cmd/fmr
go build -o filemirror ./cmd/filemirror

# Or build just one
go build -o fmr ./cmd/fmr
```

Then optionally move the binary to your PATH:
```bash
sudo mv fmr /usr/local/bin/
# or
mv fmr ~/bin/  # if ~/bin is in your PATH
```

### Quick Start

```bash
# Show help (works with either binary)
./fmr --help

# Search all files in current directory
./fmr

# Search for specific file pattern
./fmr "*.go"

# Search for files containing "config"
./fmr "config"

# Search in a different directory
./fmr --path ~/projects "*.go"
```

## Features

- **Interactive TUI**: Split-screen view with live file preview and diff mode
- **Diff Preview Mode**: Compare target files against source with colored diff view
- **Repo-Aware**: Detects git repositories and shows current branch
- **Per-Repo Worktree**: Create worktrees for staging and committing changes per repository
- **Atomic Writes**: Safe file operations that preserve permissions
- **Dry-Run Mode**: Preview operations without making changes
- **Interactive Path Navigation**: Change directories without leaving the app using TAB and CTRL-R
- **Dual Input Fields**: Separate, editable Path and Search inputs for maximum flexibility
- **Glob Pattern Support**: Filter files with wildcards (`*.go`, `*.java`, etc.) or substring matching
- **Multi-File Sync**: Copy content from one source file to multiple targets at once
- **Keyboard-Driven UI**: Fast navigation with vim-style keybindings

## Usage

```bash
fmr [OPTIONS] [PATTERN]
```

### Options

- `-p, --path PATH`: Change to directory PATH before searching
  - Supports both absolute and relative paths
  - Allows you to search files in any directory without changing your current location
  - Example: `fmr --path ~/projects "*.go"`

- `-h, --help`: Show help message
- `-v, --version`: Show version information

### Arguments

- `PATTERN` (optional): File pattern to search for
  - Wildcards: `*.go`, `*.json`, `*.md`
  - Substring: `config`, `component`, `test`
  - If omitted: shows all files in directory tree (up to 4 levels deep)

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `TAB` | Cycle focus forward: Path → Search → File List → Path |
| `Shift+TAB` | Cycle focus backward: Path ← Search ← File List |
| `CTRL-R` | Reload files from the current path (Path/Search focus) |
| `p` / `CTRL-P` | Cycle preview modes: hidden → plain → diff → hidden |
| `PgUp`/`PgDn` | Scroll preview panel (or `CTRL-U`/`CTRL-D`) |
| `?` | Show help overlay with all shortcuts |
| Type | Edit the focused input (Path or Search) |
| `↑`/`↓` or `k`/`j` | Navigate through file list (when List is focused) |
| `s` | Mark current file as SOURCE (when List is focused) |
| `Space` | Toggle current file as TARGET (when List is focused) |
| `Enter` | Proceed to confirmation (requires source + targets) |
| `y` | Confirm and execute sync operation |
| `n` / `Esc` | Cancel operation and return to selection |
| `q` / `CTRL-C` | Quit the program |

### File Display

Each file entry shows:
- **Path**: Relative to current directory
- **Size**: Human-readable format (B, KB, MB, GB)
- **Modified**: Last modification timestamp
- **Branch**: Git branch (or `-` if not in a git repo)

Files are sorted by modification time (newest first).

## Workflow Example

Here's a typical workflow for synchronizing configuration files across multiple directories:

```bash
# 1. Start FileMirror (optionally in a specific directory)
./fmr --path ~/projects "config.json"

# 2. In the interactive UI, you'll see two input fields:
#    Path:   /Users/you/projects        (editable)
#    Search: config.json                 (editable)

# 3. Navigate and edit:
#    - Press TAB to cycle focus: Path → Search → File List
#    - Edit the Path field to navigate to different directories
#    - Press CTRL-R to reload files from the new path
#    - Edit the Search field to filter files by pattern
#    - Press TAB again to focus on the file list
#    - Use ↑/↓ or k/j to navigate through files

# 4. Mark files and review diffs:
#    - Press 's' on the "correct" config file to mark it as source
#    - Press 'p' or 'CTRL-P' to cycle through preview modes (hidden/plain/diff)
#    - Press Space on target files you want to update (can mark multiple)
#    - Press Enter to review your selection

# 5. Confirm the sync:
#    - Review source and target files
#    - Press 'y' to confirm and sync
#    - All target files now have the same content as the source!
```

## Git Workflow Integration

FileMirror includes built-in git workflow support to automatically commit synced files across multiple repositories.

### How It Works

After selecting your source and target files, the confirmation screen provides git workflow options:

1. **Automatic Repository Detection**: FileMirror detects which target files are in git repositories
2. **Branch Creation**: Creates a new branch (or reuses an existing compatible branch) in each repository
3. **Worktree Safety**: Uses git worktrees to avoid disrupting your working directory
4. **Commit & Push**: Stages, commits, and optionally pushes changes to remote

### Git Workflow Screen

When you press `Enter` after selecting files, you'll see:

```
╭─────────────────────────╮│ Git Workflow Configuration ──────────┐
│ FILES TO SYNC          ││                                       │
│ Source: config.yaml    ││ [✓] Create git commit                 │
│ Targets: (3 files)     ││                                       │
│                        ││ Branch Name:                          │
│                        ││ ╭────────────────────────────────────╮│
│                        ││ │ chore/filesync-config              ││
│                        ││ ╰────────────────────────────────────╯│
│                        ││                                       │
│                        ││ Commit Message:                       │
│                        ││ ╭────────────────────────────────────╮│
│                        ││ │ chore: Sync config.yaml            ││
│                        ││ │ Synchronized from config.yaml      ││
│                        ││ ╰────────────────────────────────────╯│
│                        ││                                       │
│                        ││ [ ] Push to origin after commit      │
│                        ││                                       │
│                        ││ Repository: 2 git repos detected     │
│                        ││ ✓ project-api                        │
│                        ││ ✓ project-web                        │
╰─────────────────────────╯│                                       │
                           │    ┌──────────────┐  ┌──────────┐   │
                           │    │ Copy & Commit│  │  Cancel  │   │
                           └────┴──────────────┴──┴──────────┴───┘
```

### Keyboard Shortcuts (Git Workflow)

| Key | Action |
|-----|--------|
| `TAB` / `Shift+TAB` | Navigate between fields |
| `CTRL-G` | Toggle git workflow on/off |
| `SPACE` | Toggle checkboxes (git enabled, push) |
| Type | Edit branch name or commit message |
| `ENTER` | Execute copy & commit (when on button) |
| `ESC` | Cancel and return to file list |

### Branch Reuse Logic

FileMirror intelligently handles existing branches:

- **New Branch**: Creates the branch if it doesn't exist
- **Safe Reuse**: Reuses existing branch if it only modified the same files
- **Conflict Prevention**: Shows error if branch exists with different file changes

### Workflow Example with Git

```bash
# 1. Start FileMirror
./fmr --path ~/projects "*.yaml"

# 2. Select files in the TUI:
#    - Press 's' on canonical config.yaml
#    - Press SPACE on target files in different repos
#    - Press ENTER to open confirmation

# 3. In Git Workflow screen:
#    - Git workflow is auto-enabled if repos detected
#    - Edit branch name: "chore/sync-config"
#    - Edit commit message as needed
#    - Check "Push to origin" if desired
#    - Press TAB to focus "Copy & Commit" button
#    - Press ENTER to execute

# 4. FileMirror will:
#    - Copy files to targets
#    - Create worktrees for each repository
#    - Commit changes on new branches
#    - Push to origin (if enabled)
#    - Display summary in terminal

# Terminal output after completion:
File sync completed successfully!

Git Workflow Summary:

✓ Successfully committed to 2 repository(ies):
  - /Users/you/projects/project-api
  - /Users/you/projects/project-web

Branch: chore/sync-config
Push: NO (you can push manually later)

Next steps:
  - Review commits: git log -1 (in each repository)
  - Push to remote: git push -u origin chore/sync-config
  - Create pull requests on GitHub/GitLab
```

### Configuration

**Default Settings:**
- **Git Workflow**: Enabled if git repos detected
- **Branch Naming**: `chore/filesync-{filename}`
- **Push to Origin**: Disabled (safer default)

**Branch Name Convention:**
- Uses `chore/*` prefix for file synchronization commits
- Automatically normalizes filenames to git-safe branch names
- Example: `config.yaml` → `chore/filesync-config`

### Requirements for Git Workflow

- Git installed and available in PATH
- Target files in git repositories
- Write permission to create branches and commits
- (Optional) Configured remote for pushing

### Disable Git Workflow

If you only want to copy files without git operations:

1. In confirmation screen, press `CTRL-G` to toggle git workflow off
2. Or uncheck `[✓] Create git commit` checkbox with `SPACE`
3. Button changes to "Copy files" instead of "Copy & Commit"

## Use Cases

- **Sync configuration files** across multiple microservices or directories
- **Navigate project structures** and find files across different paths
- **Update component files** in different feature branches with automatic commits
- **Propagate utility functions** across multiple packages with git audit trail
- **Copy templates** to multiple locations in a monorepo and commit atomically
- **Standardize linting configs** across projects with tracked changes
- **Mirror canonical files** across multiple repositories with automatic branching and PR workflow

## Examples

### Search all Go files in current directory
```bash
./fmr "*.go"
```

### Find all package.json files
```bash
./fmr "package.json"
```

### Search for files containing "component"
```bash
./fmr "component"
```

### Search in a specific directory
```bash
./fmr --path /path/to/project "*.go"
```

### Search in parent directory
```bash
./fmr --path .. "config.json"
```

### Search in home directory subdirectory
```bash
./fmr -p ~/Documents/projects "*.md"
```

### Search all files (no filter)
```bash
./fmr
```

### Interactive navigation (once running)
```
1. Start: ./fmr
2. Path input is focused by default - type to edit
3. Press TAB to move to Search input
4. Type a pattern (e.g., *.ts)
5. Press TAB to move to File List
6. Use ↑/↓ or k/j to navigate files
7. Press 's' to mark source, Space to mark targets
8. Press 'p' to toggle preview mode
9. Press Enter to confirm, or TAB/Shift+TAB to change focus
```

## How It Works

1. **Scanning**: Recursively scans up to 4 directory levels
2. **Filtering**: Excludes common directories (node_modules, .git, vendor, dist, build, target, .cache)
3. **Pattern Matching**: Supports wildcards (`*.ext`) and substring matching
4. **Git Detection**: Runs `git branch --show-current` for each file's directory
5. **Safe Copying**: Uses atomic write operations (write to temp file, then rename)
6. **Permission Preservation**: Copies file permissions from source to targets

## Technical Details

### Dependencies

- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/lipgloss` - Terminal styling

### Requirements

- Go 1.21 or later
- Git (optional, for branch detection and repo features)

### Development Scripts

The project includes several utility scripts in the `bin/` directory:

#### bin/setup - Development Environment Setup

```bash
# Basic setup
bin/setup

# Full setup with optional tools
bin/setup --full
```

This script:
- Checks for required tools (Go, Git)
- Installs Go dependencies
- Verifies dependencies
- Optionally installs development tools (golangci-lint)
- Runs initial tests and build

#### bin/build - Build Binary

```bash
# Simple build
bin/build

# Build with options
bin/build --clean          # Clean and build
bin/build --verbose        # Verbose output
bin/build --install        # Install to $GOPATH/bin

# Cross-compile
bin/build -o linux -a amd64   # Build for Linux AMD64
bin/build -o windows -a amd64 # Build for Windows AMD64
```

The build script includes:
- Automatic version detection from git tags
- Build time and commit hash embedding
- Pre-build testing and linting
- Code formatting checks
- Binary verification
- Cross-compilation support

#### bin/test - Run Tests

```bash
# Run all tests
bin/test

# Run with options
bin/test --verbose         # Verbose output
bin/test --coverage        # Generate coverage report
bin/test --race            # Run with race detector
bin/test --bench           # Run benchmarks
```

#### bin/lint - Linting and Code Quality

```bash
# Run all linters
bin/lint

# Auto-fix issues
bin/lint --fix

# Verbose output
bin/lint --verbose
```

This script runs:
- `gofmt` for formatting checks
- `go vet` for static analysis
- `golangci-lint` (if installed)
- Checks for TODO/FIXME comments

### Manual Commands

You can also use standard Go commands directly:

```bash
# Build
go build -o fmr .

# Test
go test -v ./...

# Format
gofmt -w .

# Vet
go vet ./...
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT
