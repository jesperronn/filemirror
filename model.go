package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeSelect mode = iota
	modeConfirm
)

type inputFocus int

const (
	focusSearch inputFocus = iota
	focusPath
)

type model struct {
	files         []FileInfo
	filteredFiles []FileInfo
	cursor        int
	selected      map[int]bool // target files
	sourceFile    *FileInfo
	searchInput   textinput.Model
	pathInput     textinput.Model
	err           error
	width         int
	height        int
	mode          mode
	viewport      int // for scrolling
	focus         inputFocus
	workDir       string // current working directory
}

type scanCompleteMsg struct {
	files []FileInfo
	err   error
}

func initialModel(initialQuery string, initialPath string) model {
	// Search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search pattern (e.g., *.go, config.json)..."
	searchInput.Focus()
	searchInput.CharLimit = 156
	searchInput.Width = 50
	searchInput.SetValue(initialQuery)

	// Path input
	pathInput := textinput.New()
	pathInput.Placeholder = "Working directory..."
	pathInput.CharLimit = 256
	pathInput.Width = 50

	// Get current working directory or use provided path
	workDir := initialPath
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			workDir = "."
		}
	}
	pathInput.SetValue(workDir)

	m := model{
		files:         []FileInfo{},
		filteredFiles: []FileInfo{},
		cursor:        0,
		selected:      make(map[int]bool),
		searchInput:   searchInput,
		pathInput:     pathInput,
		width:         80,
		height:        24,
		mode:          modeSelect,
		viewport:      0,
		focus:         focusSearch,
		workDir:       workDir,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			files, err := scanFiles(m.workDir, m.searchInput.Value())
			return scanCompleteMsg{files: files, err: err}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanCompleteMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.files = msg.files
		m.filterFiles()
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeSelect:
			return m.updateSelect(msg)
		case modeConfirm:
			return m.updateConfirm(msg)
		}
	}

	// Update the focused input
	if m.focus == focusSearch {
		m.searchInput, cmd = m.searchInput.Update(msg)
	} else {
		m.pathInput, cmd = m.pathInput.Update(msg)
	}
	return m, cmd
}

func (m *model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		// Toggle focus between search and path inputs
		if m.focus == focusSearch {
			m.focus = focusPath
			m.searchInput.Blur()
			m.pathInput.Focus()
		} else {
			m.focus = focusSearch
			m.pathInput.Blur()
			m.searchInput.Focus()
		}
		return m, nil

	case "ctrl+r":
		// Reload: change to the path and rescan files
		newPath := m.pathInput.Value()
		absPath, err := filepath.Abs(newPath)
		if err != nil {
			m.err = fmt.Errorf("invalid path: %w", err)
			return m, nil
		}

		if err := os.Chdir(absPath); err != nil {
			m.err = fmt.Errorf("cannot change to directory: %w", err)
			return m, nil
		}

		m.workDir = absPath
		m.err = nil

		// Rescan files in new directory
		return m, func() tea.Msg {
			files, err := scanFiles(m.workDir, m.searchInput.Value())
			return scanCompleteMsg{files: files, err: err}
		}

	case "up", "k":
		if m.focus == focusSearch {
			if m.cursor > 0 {
				m.cursor--
				m.adjustViewport()
			}
		}

	case "down", "j":
		if m.focus == focusSearch {
			if m.cursor < len(m.filteredFiles)-1 {
				m.cursor++
				m.adjustViewport()
			}
		}

	case "ctrl+s":
		// Mark current file as source
		if m.cursor < len(m.filteredFiles) {
			file := m.filteredFiles[m.cursor]
			m.sourceFile = &file
		}

	case "ctrl+e":
		// Toggle target selection
		if m.cursor < len(m.filteredFiles) {
			m.selected[m.cursor] = !m.selected[m.cursor]
		}

	case "enter":
		// Proceed to confirmation if we have source and targets
		if m.sourceFile != nil && len(m.selected) > 0 {
			m.mode = modeConfirm
		}

	default:
		// Update the focused input
		var cmd tea.Cmd
		if m.focus == focusSearch {
			m.searchInput, cmd = m.searchInput.Update(msg)

			// Re-filter on search query change
			m.filterFiles()
			// Trigger rescan if query changed
			return m, tea.Batch(cmd, func() tea.Msg {
				files, err := scanFiles(m.workDir, m.searchInput.Value())
				return scanCompleteMsg{files: files, err: err}
			})
		} else {
			m.pathInput, cmd = m.pathInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Perform the copy operation
		err := m.copySourceToTargets()
		if err != nil {
			m.err = err
			return m, nil
		}
		return m, tea.Quit

	case "n", "N", "esc":
		// Go back to selection mode
		m.mode = modeSelect
		return m, nil

	case "ctrl+c", "q":
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) filterFiles() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filteredFiles = m.files
		return
	}

	m.filteredFiles = []FileInfo{}
	for _, file := range m.files {
		if strings.Contains(strings.ToLower(file.Path), query) {
			m.filteredFiles = append(m.filteredFiles, file)
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredFiles) {
		m.cursor = max(0, len(m.filteredFiles)-1)
	}
	m.adjustViewport()
}

func (m *model) adjustViewport() {
	maxVisible := m.height - 10 // Reserve space for header/footer
	if maxVisible < 1 {
		maxVisible = 1
	}

	if m.cursor < m.viewport {
		m.viewport = m.cursor
	} else if m.cursor >= m.viewport+maxVisible {
		m.viewport = m.cursor - maxVisible + 1
	}
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.\n", m.err)
	}

	switch m.mode {
	case modeSelect:
		return m.viewSelect()
	case modeConfirm:
		return m.viewConfirm()
	}

	return ""
}

func (m model) viewSelect() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(headerStyle.Render("MultiEdit - File Synchronization Tool") + "\n\n")

	// Instructions
	instructStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	b.WriteString(instructStyle.Render("TAB: switch input | CTRL-R: reload | CTRL-S: mark source | CTRL-E: toggle target | ENTER: sync | q: quit") + "\n\n")

	// Path input
	pathLabel := "Path:   "
	if m.focus == focusPath {
		pathLabel = "> Path: "
	}
	b.WriteString(pathLabel + m.pathInput.View() + "\n")

	// Search input
	searchLabel := "Search: "
	if m.focus == focusSearch {
		searchLabel = "> Search: "
	}
	b.WriteString(searchLabel + m.searchInput.View() + "\n\n")

	// Source file indicator
	if m.sourceFile != nil {
		sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s", m.sourceFile.Path)) + "\n\n")
	}

	// File list header
	headerRowStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8"))
	b.WriteString(headerRowStyle.Render(fmt.Sprintf("%-50s %-15s %-20s %s\n", "PATH", "SIZE", "MODIFIED", "BRANCH")))

	// File list
	maxVisible := m.height - 10
	if maxVisible < 1 {
		maxVisible = 1
	}

	start := m.viewport
	end := min(start+maxVisible, len(m.filteredFiles))

	for i := start; i < end; i++ {
		file := m.filteredFiles[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		marker := " "
		if m.selected[i] {
			marker = "T" // Target
		}
		if m.sourceFile != nil && m.sourceFile.Path == file.Path {
			marker = "S" // Source
		}

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = style.Background(lipgloss.Color("240"))
		}
		if m.selected[i] {
			style = style.Foreground(lipgloss.Color("11"))
		}
		if m.sourceFile != nil && m.sourceFile.Path == file.Path {
			style = style.Foreground(lipgloss.Color("10"))
		}

		line := fmt.Sprintf("%s[%s] %-47s %-15s %-20s %s",
			cursor,
			marker,
			truncate(file.Path, 47),
			formatSize(file.Size),
			file.Modified.Format("2006-01-02 15:04"),
			file.Branch,
		)

		b.WriteString(style.Render(line) + "\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	b.WriteString("\n" + footerStyle.Render(fmt.Sprintf("Showing %d of %d files | Targets: %d",
		len(m.filteredFiles), len(m.files), len(m.selected))))

	return b.String()
}

func (m model) viewConfirm() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	b.WriteString(headerStyle.Render("Confirm Synchronization") + "\n\n")

	if m.sourceFile != nil {
		sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s\n", m.sourceFile.Path)))
	}

	b.WriteString("\nTargets:\n")
	targetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	for idx := range m.selected {
		if idx < len(m.filteredFiles) {
			b.WriteString(targetStyle.Render(fmt.Sprintf("  - %s\n", m.filteredFiles[idx].Path)))
		}
	}

	b.WriteString("\n")
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	b.WriteString(warnStyle.Render("This will overwrite the target files!") + "\n\n")

	b.WriteString("Do you want to proceed? (y/n): ")

	return b.String()
}

func (m *model) copySourceToTargets() error {
	if m.sourceFile == nil {
		return fmt.Errorf("no source file selected")
	}

	for idx := range m.selected {
		if idx < len(m.filteredFiles) {
			target := m.filteredFiles[idx]
			if err := copyFile(m.sourceFile.Path, target.Path); err != nil {
				return fmt.Errorf("failed to copy to %s: %w", target.Path, err)
			}
		}
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
