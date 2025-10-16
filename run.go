package filemirror

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Build-time variables (set by ldflags)
var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// Run is the main entry point for the filemirror application
func Run() {
	var workDir string
	var initialQuery string

	// Parse command line arguments
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help", "help":
			PrintHelp()
			return
		case "-v", "--version", "version":
			fmt.Printf("fmr version %s\n", Version)
			if BuildTime != "unknown" || GitCommit != "unknown" {
				fmt.Printf("  Build time: %s\n", BuildTime)
				fmt.Printf("  Git commit: %s\n", GitCommit)
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
	m := InitialModel(initialQuery, workDir)

	// Start the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print summary if available
	if fm, ok := finalModel.(model); ok {
		if fm.exitSummary != "" {
			fmt.Print(fm.exitSummary)
		}
	}
}

// PrintHelp displays the help message
func PrintHelp() {
	help := `fmr (FileMirror) - Interactive file synchronization tool

USAGE:
    fmr [OPTIONS] [PATTERN]

DESCRIPTION:
    FileMirror helps you quickly propagate changes from one source file to
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
    TAB            Cycle focus forward: Path → Search → File List → Path
    Shift+TAB      Cycle focus backward: Path ← Search ← File List
    CTRL-R         Reload files from current path (when on Path/Search)
    p / CTRL-P     Cycle preview modes: hidden → plain → diff → hidden
    PgUp/PgDn      Scroll preview (or CTRL-U/CTRL-D)
    ?              Show help overlay with all shortcuts
    Type           Edit focused input (Path or Search)
    ↑/↓ or k/j     Navigate through file list (when List is focused)
    s              Mark current file as SOURCE (when List is focused)
    Space          Toggle current file as TARGET (when List is focused)
    Enter          Proceed to confirmation (requires source + targets)
    y              Confirm and execute sync operation
    n / Esc        Cancel operation and return to selection
    q / CTRL-C     Quit the program

WORKFLOW:
    1. Run fmr with optional path and search pattern
    2. Path input is focused - type to edit or press TAB to move to Search
    3. Edit the path to navigate to different directories (press CTRL-R to reload)
    4. Press TAB to move to Search, type pattern to filter files
    5. Press TAB to move to File List, use ↑/↓ or k/j to navigate
    6. Press 's' to mark source, Space to toggle targets
    7. Press Enter to review, then 'y' to confirm and sync

FEATURES:
    - Interactive path editing - change directories without leaving the app
    - Real-time file filtering with glob pattern support (*.go, *.java, etc.)
    - Live file preview panel - see file contents before syncing
    - Diff preview mode - compare target files against source with colored diff
    - Searches up to 4 directory levels deep
    - Excludes common directories (node_modules, .git, vendor, etc.)
    - Shows file metadata (size, modified time, git branch)
    - Safe atomic file operations preserving permissions
    - Split-screen layout with scrollable preview

EXAMPLES:
    fmr                           # Start in current directory
    fmr "*.go"                    # Start with Go files filter
    fmr --path ~/projects "*.go"  # Start in specific directory
    fmr -p /tmp config.json       # Start in /tmp, filter config.json

    Once running:
    - Path is focused by default - type to edit
    - Press TAB to move to Search field
    - Type pattern to filter, press TAB to move to file list
    - Use ↑/↓ or k/j to navigate, 's' for source, Space for targets
    - Press Enter to confirm sync, TAB to return to Path input

REPOSITORY:
    https://github.com/jesperronn/filemirror-fmr
`
	fmt.Println(help)
}
