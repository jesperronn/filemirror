package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

const version = "0.1.0"

func main() {
	var workDir string
	var initialQuery string

	// Parse command line arguments
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help", "help":
			printHelp()
			return
		case "-v", "--version", "version":
			fmt.Printf("multiedit version %s\n", version)
			return
		case "-p", "--path":
			if i+1 < len(args) {
				workDir = args[i+1]
				i++ // Skip next arg
			} else {
				fmt.Fprintf(os.Stderr, "Error: --path requires a directory argument\n")
				os.Exit(1)
			}
		default:
			// If not a flag, treat as search pattern
			if initialQuery == "" {
				initialQuery = arg
			}
		}
	}

	// Change to specified directory if provided
	if workDir != "" {
		absPath, err := filepath.Abs(workDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid path %q: %v\n", workDir, err)
			os.Exit(1)
		}

		if err := os.Chdir(absPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot change to directory %q: %v\n", absPath, err)
			os.Exit(1)
		}
	}

	// Create the model
	m := initialModel(initialQuery)

	// Start the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	help := `multiedit - Interactive file synchronization tool

USAGE:
    multiedit [OPTIONS] [PATTERN]

DESCRIPTION:
    multiedit helps you quickly propagate changes from one source file to
    multiple target files. It provides an interactive interface to search,
    select, and synchronize file contents across your project.

OPTIONS:
    -p, --path PATH    Change to directory PATH before searching
                       Supports both absolute and relative paths
    -h, --help         Show this help message
    -v, --version      Show version information

ARGUMENTS:
    PATTERN            Optional file pattern to search for (e.g., "*.go" or "config.json")
                       If omitted, shows all files in the current directory tree

KEYBOARD SHORTCUTS:
    Type           Search/filter files in real-time
    ↑/↓ or k/j     Navigate through file list
    CTRL-S         Mark current file as SOURCE
    CTRL-E         Toggle current file as TARGET (can mark multiple)
    Enter          Proceed to confirmation (requires source + targets)
    y              Confirm and execute sync operation
    n / Esc        Cancel operation and return to selection
    q / CTRL-C     Quit the program

WORKFLOW:
    1. Run multiedit with optional search pattern
    2. Browse and filter files using search
    3. Mark one file as SOURCE with CTRL-S
    4. Mark one or more files as TARGETS with CTRL-E
    5. Press Enter to review selection
    6. Confirm with 'y' to copy source content to all targets

FEATURES:
    - Searches up to 4 directory levels deep
    - Excludes common directories (node_modules, .git, vendor, etc.)
    - Shows file metadata (size, modified time, git branch)
    - Safe atomic file operations preserving permissions
    - Pattern matching with wildcards (*.go, *.json, etc.)

EXAMPLES:
    multiedit                           # Search all files in current directory
    multiedit "*.go"                    # Search for Go files
    multiedit --path ~/projects "*.go"  # Search in specific directory
    multiedit -p /tmp config.json       # Search config.json in /tmp
    multiedit --path ../backend         # Search all files in ../backend

REPOSITORY:
    https://github.com/jesper/multiedit
`
	fmt.Println(help)
}

