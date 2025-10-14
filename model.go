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
	focusPath inputFocus = iota
	focusSearch
	focusList
)

type previewMode int

const (
	previewPlain previewMode = iota
	previewDiff
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
	showPreview   bool   // whether to show file preview panel
	previewScroll int    // scroll position in preview
	previewMode   previewMode // plain or diff mode
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
		focus:         focusPath,
		workDir:       workDir,
		showPreview:   true,           // Show preview by default
		previewScroll: 0,
		previewMode:   previewPlain, // Start with plain view
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

	// Update the focused input (only if not on file list)
	if m.focus == focusSearch {
		m.searchInput, cmd = m.searchInput.Update(msg)
	} else if m.focus == focusPath {
		m.pathInput, cmd = m.pathInput.Update(msg)
	}
	return m, cmd
}

func (m *model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		// Cycle focus forward: path -> search -> file list -> path
		switch m.focus {
		case focusPath:
			m.focus = focusSearch
			m.pathInput.Blur()
			m.searchInput.Focus()
		case focusSearch:
			m.focus = focusList
			m.searchInput.Blur()
		case focusList:
			m.focus = focusPath
			m.pathInput.Focus()
		}
		return m, nil

	case "shift+tab":
		// Cycle focus backward: path <- search <- file list <- path
		switch m.focus {
		case focusPath:
			m.focus = focusList
			m.pathInput.Blur()
		case focusSearch:
			m.focus = focusPath
			m.searchInput.Blur()
			m.pathInput.Focus()
		case focusList:
			m.focus = focusSearch
			m.searchInput.Focus()
		}
		return m, nil

	case "p":
		// Toggle preview panel
		m.showPreview = !m.showPreview
		m.previewScroll = 0
		return m, nil

	case "d":
		// Toggle diff mode (only when source is selected)
		if m.sourceFile != nil {
			if m.previewMode == previewPlain {
				m.previewMode = previewDiff
			} else {
				m.previewMode = previewPlain
			}
			m.previewScroll = 0
		}
		return m, nil

	case "pagedown", "ctrl+d":
		// Scroll preview down
		if m.showPreview {
			m.previewScroll += 10
		}
		return m, nil

	case "pageup", "ctrl+u":
		// Scroll preview up
		if m.showPreview {
			m.previewScroll -= 10
			if m.previewScroll < 0 {
				m.previewScroll = 0
			}
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
		if m.focus == focusList {
			if m.cursor > 0 {
				m.cursor--
				m.adjustViewport()
			}
		}

	case "down", "j":
		if m.focus == focusList {
			if m.cursor < len(m.filteredFiles)-1 {
				m.cursor++
				m.adjustViewport()
			}
		}

	case "s":
		// Mark current file as source (when on file list)
		if m.focus == focusList && m.cursor < len(m.filteredFiles) {
			file := m.filteredFiles[m.cursor]
			m.sourceFile = &file
		}

	case " ": // Space
		// Toggle target selection (when on file list)
		if m.focus == focusList && m.cursor < len(m.filteredFiles) {
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
		m.resetCursorIfNeeded()
		return
	}

	m.filteredFiles = []FileInfo{}
	for _, file := range m.files {
		if matchesFilePattern(file.Path, query) {
			m.filteredFiles = append(m.filteredFiles, file)
		}
	}

	m.resetCursorIfNeeded()
}

func (m *model) resetCursorIfNeeded() {
	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredFiles) {
		m.cursor = max(0, len(m.filteredFiles)-1)
	}
	m.adjustViewport()
}

// matchesFilePattern checks if a file path matches a pattern
// Supports both glob patterns (*.go) and substring matching
func matchesFilePattern(filePath, pattern string) bool {
	pattern = strings.ToLower(pattern)
	filePath = strings.ToLower(filePath)

	// Extract just the filename from the path for glob matching
	filename := filepath.Base(filePath)

	// If pattern contains *, use glob matching on filename
	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, filename)
		if matched {
			return true
		}
		// Also try matching against full path for patterns like "src/*.go"
		matched, _ = filepath.Match(pattern, filePath)
		return matched
	}

	// Otherwise, simple substring match on full path
	return strings.Contains(filePath, pattern)
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
	previewHint := "p: toggle preview"
	if !m.showPreview {
		previewHint = "p: show preview"
	}

	// Add diff toggle hint if source is selected
	diffHint := ""
	if m.sourceFile != nil && m.showPreview {
		if m.previewMode == previewPlain {
			diffHint = " | d: show diff"
		} else {
			diffHint = " | d: plain view"
		}
	}

	focusHint := "Path"
	switch m.focus {
	case focusPath:
		focusHint = "Path"
	case focusSearch:
		focusHint = "Search"
	case focusList:
		focusHint = "List"
	}
	instructions := ""
	if m.focus == focusList {
		instructions = fmt.Sprintf("TAB: cycle focus (%s) | %s%s | s: source | SPACE: target | ENTER: confirm | q: quit", focusHint, previewHint, diffHint)
	} else {
		instructions = fmt.Sprintf("TAB: cycle focus (%s) | CTRL-R: reload | %s%s | q: quit", focusHint, previewHint, diffHint)
	}
	b.WriteString(instructStyle.Render(instructions) + "\n\n")

	// Path input
	pathLabel := "  Path:   "
	if m.focus == focusPath {
		pathLabel = "→ Path:   "
	}
	b.WriteString(pathLabel + m.pathInput.View() + "\n")

	// Search input
	searchLabel := "  Search: "
	if m.focus == focusSearch {
		searchLabel = "→ Search: "
	}
	b.WriteString(searchLabel + m.searchInput.View() + "\n\n")

	// Source file indicator
	if m.sourceFile != nil {
		sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s", m.sourceFile.Path)) + "\n\n")
	}

	// Determine layout based on preview setting
	var fileListWidth int
	var previewContent string

	if m.showPreview {
		// Split screen: 50% file list, 50% preview
		fileListWidth = m.width / 2
		previewContent = m.renderPreview()
	} else {
		// Full width for file list
		fileListWidth = m.width
	}

	// File list header
	headerRowStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8"))
	pathWidth := min(fileListWidth-40, 50) // Adjust based on available space
	b.WriteString(headerRowStyle.Render(fmt.Sprintf("%-*s %-10s %-15s\n", pathWidth, "PATH", "SIZE", "MODIFIED")))

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

		pathDisplayWidth := pathWidth - 5 // Account for cursor and marker
		line := fmt.Sprintf("%s[%s] %-*s %-10s %-15s",
			cursor,
			marker,
			pathDisplayWidth,
			truncate(file.Path, pathDisplayWidth),
			formatSize(file.Size),
			file.Modified.Format("2006-01-02 15:04"),
		)

		b.WriteString(style.Render(line) + "\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	b.WriteString("\n" + footerStyle.Render(fmt.Sprintf("Showing %d of %d files | Targets: %d",
		len(m.filteredFiles), len(m.files), len(m.selected))))

	// If preview is enabled, combine file list and preview side by side
	if m.showPreview && previewContent != "" {
		fileList := b.String()
		return lipgloss.JoinHorizontal(lipgloss.Top, fileList, previewContent)
	}

	return b.String()
}

// renderPreview renders the file preview panel
func (m model) renderPreview() string {
	if len(m.filteredFiles) == 0 || m.cursor >= len(m.filteredFiles) {
		return m.renderEmptyPreview()
	}

	currentFile := m.filteredFiles[m.cursor]
	filePath := filepath.Join(m.workDir, currentFile.Path)

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return m.renderPreviewError(fmt.Sprintf("Error reading file: %v", err))
	}

	// Determine what to show based on preview mode
	var lines []string
	var headerTitle string

	if m.previewMode == previewDiff && m.sourceFile != nil {
		// Show diff against source file
		sourceFilePath := filepath.Join(m.workDir, m.sourceFile.Path)
		sourceContent, err := os.ReadFile(sourceFilePath)
		if err != nil {
			return m.renderPreviewError(fmt.Sprintf("Error reading source file: %v", err))
		}

		// Generate diff
		lines = m.generateDiff(string(sourceContent), string(content))
		headerTitle = fmt.Sprintf(" Diff: %s → %s ", m.sourceFile.Path, currentFile.Path)
	} else {
		// Show plain file content
		lines = strings.Split(string(content), "\n")
		headerTitle = fmt.Sprintf(" Preview: %s ", currentFile.Path)
	}

	// Calculate preview dimensions
	previewWidth := m.width / 2
	previewHeight := m.height - 10

	var b strings.Builder

	// Preview header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Width(previewWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true)

	b.WriteString(headerStyle.Render(headerTitle) + "\n")

	// Preview content
	contentStyle := lipgloss.NewStyle().
		Width(previewWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true).
		PaddingLeft(1)

	// Render visible lines
	start := m.previewScroll
	end := min(start+previewHeight, len(lines))

	for i := start; i < end; i++ {
		line := lines[i]
		// Truncate long lines
		if len(line) > previewWidth-3 {
			line = line[:previewWidth-6] + "..."
		}

		// Color diff lines if in diff mode
		if m.previewMode == previewDiff && m.sourceFile != nil {
			lineStyle := contentStyle
			if len(line) > 0 {
				switch line[0] {
				case '+':
					lineStyle = contentStyle.Foreground(lipgloss.Color("10")) // Green for additions
				case '-':
					lineStyle = contentStyle.Foreground(lipgloss.Color("9")) // Red for deletions
				case '@':
					lineStyle = contentStyle.Foreground(lipgloss.Color("12")) // Blue for context markers
				}
			}
			b.WriteString(lineStyle.Render(line) + "\n")
		} else {
			b.WriteString(contentStyle.Render(line) + "\n")
		}
	}

	// Fill remaining space if content is short
	for i := end - start; i < previewHeight; i++ {
		b.WriteString(contentStyle.Render("") + "\n")
	}

	// Preview footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(previewWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true)

	scrollInfo := ""
	if len(lines) > previewHeight {
		scrollInfo = fmt.Sprintf(" [%d-%d of %d lines] PgUp/PgDn to scroll ", start+1, end, len(lines))
	} else {
		scrollInfo = fmt.Sprintf(" [%d lines] ", len(lines))
	}
	b.WriteString(footerStyle.Render(scrollInfo))

	return b.String()
}

func (m model) renderEmptyPreview() string {
	previewWidth := m.width / 2
	style := lipgloss.NewStyle().
		Width(previewWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true).
		Foreground(lipgloss.Color("240")).
		Padding(1)

	return style.Render("No file selected")
}

func (m model) renderPreviewError(errMsg string) string {
	previewWidth := m.width / 2
	style := lipgloss.NewStyle().
		Width(previewWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true).
		Foreground(lipgloss.Color("9")).
		Padding(1)

	return style.Render(errMsg)
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

// generateDiff generates a simple unified diff between two strings
func (m model) generateDiff(source, target string) []string {
	sourceLines := strings.Split(source, "\n")
	targetLines := strings.Split(target, "\n")

	var result []string
	result = append(result, fmt.Sprintf("@@ Source: %s → Target @@", m.sourceFile.Path))

	// Simple line-by-line comparison
	maxLen := max(len(sourceLines), len(targetLines))

	for i := 0; i < maxLen; i++ {
		var sourceLine, targetLine string

		if i < len(sourceLines) {
			sourceLine = sourceLines[i]
		}
		if i < len(targetLines) {
			targetLine = targetLines[i]
		}

		// If lines are different, show both
		if sourceLine != targetLine {
			if i < len(sourceLines) {
				result = append(result, fmt.Sprintf("-%s", sourceLine))
			}
			if i < len(targetLines) {
				result = append(result, fmt.Sprintf("+%s", targetLine))
			}
		} else {
			// Lines are the same, show context
			result = append(result, fmt.Sprintf(" %s", sourceLine))
		}
	}

	return result
}
