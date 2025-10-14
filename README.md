# multiedit

A Go CLI tool for interactively finding, selecting, and synchronizing file content across multiple files in your project.

## Overview

`multiedit` helps you quickly propagate changes from one source file to multiple target files. This is particularly useful when you have the same file (e.g., configuration, utilities, components) in multiple locations and need to keep them in sync.

## Features

### Core Functionality

- **Interactive File Search**: Search for files by name pattern with real-time filtering
- **File Selection**: Browse and select files from search results
- **Source/Target Selection**: Mark one file as the source and multiple files as targets
- **Content Synchronization**: Copy content from the source file to all selected target files
- **File Preview**: View file contents with syntax highlighting before synchronizing

### Display Information

For each file, display:
- File path (relative to current directory)
- File size
- Last modified timestamp
- Git branch (if file is in a git repository)

### Interactive Controls

- Real-time search/filtering as you type
- Keyboard shortcuts for marking source and target files
- Multi-select capability for target files
- File content preview with syntax highlighting
- Automatic sorting by modification time (newest first)

## Requirements

### Functional Requirements

1. **File Discovery**
   - Search for files recursively up to 4 levels deep
   - Filter by filename pattern (glob/wildcard support)
   - Exclude common directories (e.g., node_modules)
   - Display results sorted by modification time

2. **Interactive Selection**
   - Present results in an interactive list
   - Allow marking one file as SOURCE
   - Allow marking multiple files as TARGETS
   - Show visual indication of marked files

3. **File Synchronization**
   - Copy content from source file to target file(s)
   - Handle errors gracefully (permissions, missing files, etc.)
   - Optionally show diff before copying
   - Confirm before overwriting

4. **Preview**
   - Display file contents with syntax highlighting
   - Show git branch for files in repositories
   - Show file metadata (size, modified date)

### Technical Requirements

1. **Terminal UI**
   - Interactive fuzzy finder interface
   - Keyboard-driven navigation
   - Real-time search filtering
   - Multi-column layout (path, size, date)

2. **File Operations**
   - Safe file copying with proper error handling
   - Preserve file permissions where appropriate
   - Atomic operations to prevent partial writes

3. **Git Integration**
   - Detect if file is in a git repository
   - Show current branch name
   - Handle files outside git repos gracefully

4. **Performance**
   - Handle large numbers of files efficiently
   - Lazy loading for large directories
   - Fast search and filtering

## Implementation Approach

The tool should be built as a single Go binary with the following components:

1. **File Scanner**: Recursively find files matching patterns
2. **Interactive UI**: Terminal-based interface for selection (using bubbletea or similar)
3. **File Operations**: Safe copying with error handling
4. **Git Integration**: Query repository information
5. **Syntax Highlighting**: Display file previews with highlighting

## Usage (Proposed)

```bash
# Search for all files
multiedit

# Search for files matching pattern
multiedit "*.go"

# Search for specific filename
multiedit "config.json"
```

## Workflow

1. Run `multiedit` with optional search pattern
2. Browse through matching files (filtered by search query)
3. Press hotkey to mark one file as SOURCE
4. Press hotkey to mark one or more files as TARGETS
5. Confirm to copy source content to all targets

## Dependencies (Suggested)

- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/alecthomas/chroma` - Syntax highlighting
- `github.com/go-git/go-git/v5` - Git integration

## Non-Functional Requirements

- Cross-platform (Linux, macOS, Windows)
- No external dependencies (single binary)
- Fast startup time
- Low memory footprint
- Clear error messages
