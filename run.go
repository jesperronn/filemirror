package filemirror

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
)

// Build-time variables (set by ldflags during build).
// If not set (e.g., via `go install`), they are detected at runtime.
var (
	Version   = ""
	BuildTime = ""
	GitCommit = ""
)

func init() {
	// If version wasn't set at build time via ldflags,
	// try to get it from the module build info (e.g., from `go install`)
	if Version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			// For `go install github.com/.../fmr@v1.1.0`,
			// info.Main.Version will be "v1.1.0"
			if info.Main.Version != "(devel)" && info.Main.Version != "" {
				Version = info.Main.Version
			} else {
				Version = "dev"
			}
		} else {
			Version = "dev"
		}
	}

	// Set defaults for other fields if not provided
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if GitCommit == "" {
		GitCommit = "unknown"
	}
}

// Config holds the parsed command-line configuration
type Config struct {
	WorkDir      string
	InitialQuery string
	ShowHelp     bool
	ShowVersion  bool
}

// parseArgs parses command-line arguments and returns a Config
// Returns an error if arguments are invalid
func parseArgs(args []string) (Config, error) {
	var cfg Config

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help", "help":
			cfg.ShowHelp = true
			return cfg, nil
		case "-v", "--version", "version":
			cfg.ShowVersion = true
			return cfg, nil
		case "-p", "--path":
			if i+1 < len(args) {
				cfg.WorkDir = args[i+1]
				i++ // Skip next arg
			} else {
				return cfg, errors.New("--path requires a directory argument")
			}
		default:
			// If not a flag, treat as search pattern
			if cfg.InitialQuery == "" {
				cfg.InitialQuery = arg
			}
		}
	}

	return cfg, nil
}

// validateAndSetupWorkDir validates the working directory and changes to it
// Returns the absolute path if successful
func validateAndSetupWorkDir(workDir string) (string, error) {
	if workDir == "" {
		return "", nil
	}

	absPath, err := filepath.Abs(workDir)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", workDir, err)
	}

	if err := os.Chdir(absPath); err != nil {
		return "", fmt.Errorf("cannot change to directory %q: %w", absPath, err)
	}

	return absPath, nil
}

// printVersion prints version information to the writer
func printVersion(w io.Writer) {
	_, _ = fmt.Fprintf(w, "fmr version %s\n", Version) //nolint:errcheck // Error writing to stdout/stderr is not actionable
	if BuildTime != "unknown" || GitCommit != "unknown" {
		_, _ = fmt.Fprintf(w, "  Build time: %s\n", BuildTime) //nolint:errcheck // Error writing to stdout/stderr is not actionable
		_, _ = fmt.Fprintf(w, "  Git commit: %s\n", GitCommit) //nolint:errcheck // Error writing to stdout/stderr is not actionable
	}
}

// Run is the main entry point for the filemirror application
func Run() {
	code := RunWithArgs(os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		os.Exit(code)
	}
}

// RunWithArgs runs the application with provided arguments and writers
// Returns an exit code (0 for success, non-zero for errors)
func RunWithArgs(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err) //nolint:errcheck // Error writing to stderr is not actionable
		return 1
	}

	if cfg.ShowHelp {
		PrintHelpTo(stdout)
		return 0
	}

	if cfg.ShowVersion {
		printVersion(stdout)
		return 0
	}

	// Validate and setup working directory
	absPath, err := validateAndSetupWorkDir(cfg.WorkDir)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err) //nolint:errcheck // Error writing to stderr is not actionable
		return 1
	}

	// Use the absolute path if we changed directories
	workDir := cfg.WorkDir
	if absPath != "" {
		workDir = absPath
	}

	// Create the model with initial query and working directory
	m := InitialModel(cfg.InitialQuery, workDir)

	// Start the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err) //nolint:errcheck // Error writing to stderr is not actionable
		return 1
	}

	// Print summary if available
	if fm, ok := finalModel.(model); ok {
		if fm.exitSummary != "" {
			_, _ = fmt.Fprint(stdout, fm.exitSummary) //nolint:errcheck // Error writing to stdout is not actionable
		}
	}

	return 0
}

// PrintHelp displays the help message to stdout
func PrintHelp() {
	PrintHelpTo(os.Stdout)
}

// PrintHelpTo displays the help message to the specified writer
func PrintHelpTo(w io.Writer) {
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
	_, _ = fmt.Fprintln(w, help) //nolint:errcheck // Error writing to writer is not actionable
}
