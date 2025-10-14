# multiedit

A Go CLI tool for interactively finding, selecting, and synchronizing file content across multiple files in your project.

## Overview

`multiedit` helps you quickly propagate changes from one source file to multiple target files. This is particularly useful when you have the same file (e.g., configuration, utilities, components) in multiple locations and need to keep them in sync.

## Installation

### From Source

```bash
git clone https://github.com/jesper/multiedit
cd multiedit
go build -o multiedit .
```

### Quick Start

```bash
# Show help
./multiedit --help

# Search all files in current directory
./multiedit

# Search for specific file pattern
./multiedit "*.go"

# Search for files containing "config"
./multiedit "config"

# Search in a different directory
./multiedit --path ~/projects "*.go"
```

## Features

- **Interactive File Search**: Real-time filtering with pattern matching (wildcards supported)
- **Multi-File Sync**: Copy content from one source file to multiple targets at once
- **File Metadata Display**: Shows path, size, modification time, and git branch
- **Safe Operations**: Atomic file copying that preserves permissions
- **Smart Filtering**: Excludes common directories (node_modules, .git, vendor, etc.)
- **Keyboard-Driven UI**: Fast navigation with vim-style keybindings
- **Git Integration**: Shows current branch for files in git repositories

## Usage

```bash
multiedit [OPTIONS] [PATTERN]
```

### Options

- `-p, --path PATH`: Change to directory PATH before searching
  - Supports both absolute and relative paths
  - Allows you to search files in any directory without changing your current location
  - Example: `multiedit --path ~/projects "*.go"`

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
| Type | Search/filter files in real-time |
| `↑`/`↓` or `k`/`j` | Navigate through file list |
| `CTRL-S` | Mark current file as SOURCE |
| `CTRL-E` | Toggle current file as TARGET (can mark multiple) |
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
# 1. Start multiedit searching for config files
./multiedit "config.json"

# 2. In the interactive UI:
#    - Browse through the list of config.json files
#    - Press CTRL-S on the "correct" config file to mark it as source
#    - Navigate to other config files you want to update
#    - Press CTRL-E on each one to mark them as targets
#    - Press Enter to review your selection

# 3. Confirm the sync operation:
#    - Review source and target files
#    - Press 'y' to confirm and sync
#    - All target files now have the same content as the source!
```

## Use Cases

- **Sync configuration files** across multiple microservices
- **Update component files** in different feature branches
- **Propagate utility functions** across multiple packages
- **Copy templates** to multiple locations in a monorepo
- **Standardize linting configs** across projects

## Examples

### Search all Go files in current directory
```bash
./multiedit "*.go"
```

### Find all package.json files
```bash
./multiedit "package.json"
```

### Search for files containing "component"
```bash
./multiedit "component"
```

### Search in a specific directory
```bash
./multiedit --path /path/to/project "*.go"
```

### Search in parent directory
```bash
./multiedit --path .. "config.json"
```

### Search in home directory subdirectory
```bash
./multiedit -p ~/Documents/projects "*.md"
```

### Search all files (no filter)
```bash
./multiedit
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
- Git (optional, for branch detection)

### Build

```bash
go build -o multiedit .
```

### Test

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT
