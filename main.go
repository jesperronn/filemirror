package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Build-time variables (set by ldflags)
var (
	version   = "0.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

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
			if buildTime != "unknown" || gitCommit != "unknown" {
				fmt.Printf("  Build time: %s\n", buildTime)
				fmt.Printf("  Git commit: %s\n", gitCommit)
			}
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

	// Create the model with initial query and working directory
	m := initialModel(initialQuery, workDir)

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
    TAB            Switch between Path and Search inputs
    CTRL-R         Reload files from current path
    p              Toggle file preview panel
    PgUp/PgDn      Scroll preview (or CTRL-U/CTRL-D)
    Type           Edit focused input (Path or Search)
    ↑/↓ or k/j     Navigate through file list (when Search is focused)
    CTRL-S         Mark current file as SOURCE
    CTRL-E         Toggle current file as TARGET (can mark multiple)
    Enter          Proceed to confirmation (requires source + targets)
    y              Confirm and execute sync operation
    n / Esc        Cancel operation and return to selection
    q / CTRL-C     Quit the program

WORKFLOW:
    1. Run multiedit with optional path and search pattern
    2. Use TAB to switch between Path and Search inputs
    3. Edit the path to navigate to different directories (press CTRL-R to reload)
    4. Type in Search to filter files by pattern
    5. Navigate with ↑/↓ or k/j, mark source (CTRL-S) and targets (CTRL-E)
    6. Press Enter to review, then 'y' to confirm and sync

FEATURES:
    - Interactive path editing - change directories without leaving the app
    - Real-time file filtering with glob pattern support (*.go, *.java, etc.)
    - Live file preview panel - see file contents before syncing
    - Searches up to 4 directory levels deep
    - Excludes common directories (node_modules, .git, vendor, etc.)
    - Shows file metadata (size, modified time, git branch)
    - Safe atomic file operations preserving permissions
    - Split-screen layout with scrollable preview

EXAMPLES:
    multiedit                           # Start in current directory
    multiedit "*.go"                    # Start with Go files filter
    multiedit --path ~/projects "*.go"  # Start in specific directory
    multiedit -p /tmp config.json       # Start in /tmp, filter config.json

    Once running:
    - Press TAB to edit the path field
    - Type a new path and press CTRL-R to navigate there
    - Press TAB again to go back to search
    - Edit search pattern to filter files

REPOSITORY:
    https://github.com/jesper/multiedit
`
	fmt.Println(help)
}

