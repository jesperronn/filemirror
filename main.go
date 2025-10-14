package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const version = "0.1.0"

func main() {
	// Handle command line flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help", "help":
			printHelp()
			return
		case "-v", "--version", "version":
			fmt.Printf("multiedit version %s\n", version)
			return
		}
	}

	// Get initial query from command line args
	initialQuery := ""
	if len(os.Args) > 1 {
		initialQuery = os.Args[1]
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
    multiedit [PATTERN]

DESCRIPTION:
    multiedit helps you quickly propagate changes from one source file to
    multiple target files. It provides an interactive interface to search,
    select, and synchronize file contents across your project.

ARGUMENTS:
    PATTERN    Optional file pattern to search for (e.g., "*.go" or "config.json")
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
    multiedit                    # Search all files
    multiedit "*.go"             # Search for Go files
    multiedit "config.json"      # Search for config.json files
    multiedit "component"        # Search for files containing "component"

FLAGS:
    -h, --help       Show this help message
    -v, --version    Show version information

REPOSITORY:
    https://github.com/jesper/multiedit
`
	fmt.Println(help)
}

