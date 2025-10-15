package filemirror

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
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

type confirmFocus int

const (
	focusCopyButton confirmFocus = iota
	focusCancelButton
	focusGitEnabled
	focusBranchName
	focusCommitMsg
	focusPushToggle
)

type previewMode int

const (
	previewHidden previewMode = iota
	previewPlain
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
	workDir       string      // current working directory
	previewScroll int         // scroll position in preview
	previewMode   previewMode // hidden, plain, or diff mode
	showHelp      bool        // whether to show help overlay

	// Git workflow fields (integrated into modeConfirm)
	gitEnabled      bool
	branchNameInput textinput.Model
	commitMsgInput  textarea.Model
	shouldPush      bool
	confirmFocus    confirmFocus
	gitRepos        map[string][]string // repo path -> list of changed files
}

type scanCompleteMsg struct {
	files []FileInfo
	err   error
}

func InitialModel(initialQuery, initialPath string) model {
	// Search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search pattern (e.g., *.go, config.json)..."
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
		focus:         focusList, // Start with file list focused
		workDir:       workDir,
		previewScroll: 0,
		previewMode:   previewPlain, // Start with plain view (can be changed to previewHidden)
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

	return m, nil
}

func (m *model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input field updates FIRST when focused (before command keys)
	// This prevents keys like 's', 'k', 'j', etc. from being intercepted
	if m.focus == focusPath || m.focus == focusSearch {
		var cmd tea.Cmd
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+p":
			// Cycle through preview modes (even when in input fields)
			m.previewMode = (m.previewMode + 1) % 3 // Cycle: hidden -> plain -> diff -> hidden
			m.previewScroll = 0
			return m, nil
		case "pagedown", "ctrl+d", "end":
			// Scroll preview down (works in any focus when preview is visible)
			// end key is fn+down on MacBook
			if m.previewMode != previewHidden {
				m.previewScroll += 10
			}
			return m, nil
		case "pageup", "ctrl+u", "home":
			// Scroll preview up (works in any focus when preview is visible)
			// home key is fn+up on MacBook
			if m.previewMode != previewHidden {
				m.previewScroll -= 10
				if m.previewScroll < 0 {
					m.previewScroll = 0
				}
			}
			return m, nil
		case "tab":
			// Handle tab to switch focus
			switch m.focus {
			case focusPath:
				m.focus = focusSearch
				m.pathInput.Blur()
				m.searchInput.Focus()
			case focusSearch:
				m.focus = focusList
				m.searchInput.Blur()
			}
			return m, nil
		case "shift+tab":
			// Handle shift+tab to switch focus backwards
			switch m.focus {
			case focusPath:
				m.focus = focusList
				m.pathInput.Blur()
			case focusSearch:
				m.focus = focusPath
				m.searchInput.Blur()
				m.pathInput.Focus()
			}
			return m, nil
		case "enter":
			// Enter triggers reload and moves to next field
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

			// Move to next field
			switch m.focus {
			case focusPath:
				m.focus = focusSearch
				m.pathInput.Blur()
				m.searchInput.Focus()
			case focusSearch:
				m.focus = focusList
				m.searchInput.Blur()
			}

			return m, func() tea.Msg {
				files, err := scanFiles(m.workDir, m.searchInput.Value())
				return scanCompleteMsg{files: files, err: err}
			}
		case "ctrl+r":
			// Reload files
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
			return m, func() tea.Msg {
				files, err := scanFiles(m.workDir, m.searchInput.Value())
				return scanCompleteMsg{files: files, err: err}
			}
		default:
			// Let the input handle all other keys (including typing)
			if m.focus == focusSearch {
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterFiles()
				return m, tea.Batch(cmd, func() tea.Msg {
					files, err := scanFiles(m.workDir, m.searchInput.Value())
					return scanCompleteMsg{files: files, err: err}
				})
			} else if m.focus == focusPath {
				m.pathInput, cmd = m.pathInput.Update(msg)
				return m, cmd
			}
		}
	}

	// Handle file list and other commands when NOT in input mode
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?":
		// Toggle help overlay
		m.showHelp = !m.showHelp
		return m, nil

	case "esc":
		// Close help overlay if open
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

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

	case "p", "ctrl+p":
		// Cycle through preview modes (works in any focus mode)
		m.previewMode = (m.previewMode + 1) % 3 // Cycle: hidden -> plain -> diff -> hidden
		m.previewScroll = 0
		return m, nil

	case "pagedown", "ctrl+d", "end":
		// Scroll preview down (works in any focus when preview is visible)
		// end key is fn+down on MacBook
		if m.previewMode != previewHidden {
			m.previewScroll += 10
		}
		return m, nil

	case "pageup", "ctrl+u", "home":
		// Scroll preview up (works in any focus when preview is visible)
		// home key is fn+up on MacBook
		if m.previewMode != previewHidden {
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
			m.initGitWorkflow()
		}
	}

	return m, nil
}

func (m *model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle textarea input when focused
	if m.confirmFocus == focusCommitMsg {
		var cmd tea.Cmd
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			// Go back to selection mode
			m.mode = modeSelect
			return m, nil
		case "tab":
			// Move to next field
			m.confirmFocus = focusPushToggle
			m.commitMsgInput.Blur()
			return m, nil
		case "shift+tab":
			// Move to previous field
			m.confirmFocus = focusBranchName
			m.commitMsgInput.Blur()
			m.branchNameInput.Focus()
			return m, nil
		default:
			// Let textarea handle the input
			m.commitMsgInput, cmd = m.commitMsgInput.Update(msg)
			return m, cmd
		}
	}

	// Handle branch name input when focused
	if m.confirmFocus == focusBranchName {
		var cmd tea.Cmd
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.mode = modeSelect
			return m, nil
		case "tab":
			m.confirmFocus = focusCommitMsg
			m.branchNameInput.Blur()
			m.commitMsgInput.Focus()
			return m, nil
		case "shift+tab":
			m.confirmFocus = focusGitEnabled
			m.branchNameInput.Blur()
			return m, nil
		case "enter":
			// Move to commit message
			m.confirmFocus = focusCommitMsg
			m.branchNameInput.Blur()
			m.commitMsgInput.Focus()
			return m, nil
		default:
			m.branchNameInput, cmd = m.branchNameInput.Update(msg)
			return m, cmd
		}
	}

	// Handle other keys
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		// Go back to selection mode
		m.mode = modeSelect
		return m, nil

	case "ctrl+g":
		// Quick toggle git enabled
		m.gitEnabled = !m.gitEnabled
		return m, nil

	case "tab":
		// Cycle focus forward
		switch m.confirmFocus {
		case focusCopyButton:
			m.confirmFocus = focusCancelButton
		case focusCancelButton:
			if m.gitEnabled {
				m.confirmFocus = focusGitEnabled
			} else {
				// When git disabled, cycle back to copy button
				m.confirmFocus = focusCopyButton
			}
		case focusGitEnabled:
			m.confirmFocus = focusBranchName
			m.branchNameInput.Focus()
		case focusBranchName:
			m.confirmFocus = focusCommitMsg
			m.branchNameInput.Blur()
			m.commitMsgInput.Focus()
		case focusCommitMsg:
			m.confirmFocus = focusPushToggle
			m.commitMsgInput.Blur()
		case focusPushToggle:
			m.confirmFocus = focusCopyButton
		}
		return m, nil

	case "shift+tab":
		// Cycle focus backward
		switch m.confirmFocus {
		case focusCopyButton:
			m.confirmFocus = focusPushToggle
		case focusCancelButton:
			m.confirmFocus = focusCopyButton
		case focusGitEnabled:
			m.confirmFocus = focusCancelButton
		case focusBranchName:
			m.confirmFocus = focusGitEnabled
			m.branchNameInput.Blur()
		case focusCommitMsg:
			m.confirmFocus = focusBranchName
			m.commitMsgInput.Blur()
			m.branchNameInput.Focus()
		case focusPushToggle:
			m.confirmFocus = focusCommitMsg
			m.commitMsgInput.Focus()
		}
		return m, nil

	case " ":
		// Toggle checkboxes
		switch m.confirmFocus {
		case focusGitEnabled:
			m.gitEnabled = !m.gitEnabled
		case focusPushToggle:
			m.shouldPush = !m.shouldPush
		}
		return m, nil

	case "enter":
		// Execute on copy button or cancel button
		if m.confirmFocus == focusCopyButton {
			// Perform the copy operation
			err := m.copySourceToTargets()
			if err != nil {
				m.err = err
				return m, nil
			}

			// Perform git workflow if enabled
			if m.gitEnabled && len(m.gitRepos) > 0 {
				branchName := m.branchNameInput.Value()
				commitMsg := m.commitMsgInput.Value()

				successRepos, errors := performGitWorkflow(m.gitRepos, branchName, commitMsg, m.shouldPush)

				// Handle errors
				if len(errors) > 0 {
					errMsg := "Git workflow errors:\n"
					for _, err := range errors {
						errMsg += fmt.Sprintf("- %v\n", err)
					}
					if len(successRepos) > 0 {
						errMsg += fmt.Sprintf("\nSuccessfully committed to %d repositories", len(successRepos))
					}
					m.err = fmt.Errorf("%s", errMsg)
					return m, nil
				}
			}

			return m, tea.Quit
		} else if m.confirmFocus == focusCancelButton {
			// Cancel and go back to selection
			m.mode = modeSelect
			return m, nil
		}
		return m, nil
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
		m.cursor = maxInt(0, len(m.filteredFiles)-1)
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
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
		// Also try matching against full path for patterns like "src/*.go"
		matched, err = filepath.Match(pattern, filePath)
		return err == nil && matched
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

	var baseView string
	switch m.mode {
	case modeSelect:
		baseView = m.viewSelect()
	case modeConfirm:
		baseView = m.viewConfirm()
	default:
		return ""
	}

	// Overlay help modal if active
	if m.showHelp {
		return m.renderHelpOverlay()
	}

	return baseView
}

func (m model) viewSelect() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(headerStyle.Render("FileMirror - File Synchronization Tool") + "\n\n")

	// Context-sensitive keyboard hints
	// Use adaptive color: dark on light backgrounds, light on dark backgrounds
	// Always visible regardless of preview state
	instructStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}).
		MaxWidth(m.width).
		Inline(true)
	var hints string

	switch m.focus {
	case focusPath:
		pathHints := []string{"Type to edit", "ENTER: reload & next", "TAB: next", "CTRL-P: cycle preview"}
		if m.previewMode != previewHidden {
			pathHints = append(pathHints, "CTRL-U/D: scroll preview")
		}
		pathHints = append(pathHints, "CTRL-C: quit")
		hints = "PATH: " + strings.Join(pathHints, " • ")
	case focusSearch:
		searchHints := []string{"Type pattern (*.go, config)", "ENTER: reload & next", "TAB: next", "Shift+TAB: prev", "CTRL-P: cycle preview"}
		if m.previewMode != previewHidden {
			searchHints = append(searchHints, "CTRL-U/D: scroll preview")
		}
		searchHints = append(searchHints, "CTRL-C: quit")
		hints = "SEARCH: " + strings.Join(searchHints, " • ")
	case focusList:
		fileHints := []string{"↑/↓ or k/j: navigate", "s: set source", "SPACE: toggle target"}
		if m.sourceFile != nil && len(m.selected) > 0 {
			fileHints = append(fileHints, "ENTER: confirm sync")
		}
		// Show preview mode hint
		previewModeStr := map[previewMode]string{
			previewHidden: "preview plain",
			previewPlain:  "preview diff",
			previewDiff:   "hide preview",
		}[m.previewMode]
		fileHints = append(fileHints, fmt.Sprintf("p/CTRL-P: %s", previewModeStr))

		if m.previewMode != previewHidden {
			fileHints = append(fileHints, "CTRL-U/D: scroll preview")
		}
		fileHints = append(fileHints, "TAB: next", "?: help", "q: quit")
		hints = "FILE LIST: " + strings.Join(fileHints, " • ")
	}

	b.WriteString(instructStyle.Render(hints) + "\n\n")

	// Path input with border
	pathBorderColor := lipgloss.Color("240")
	pathLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	if m.focus == focusPath {
		pathBorderColor = lipgloss.Color("12") // Bright blue when focused
		pathLabelStyle = pathLabelStyle.Bold(true)
	}
	pathBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(pathBorderColor).
		Padding(0, 1).
		Width(m.width - 4)

	pathContent := pathLabelStyle.Render("PATH") + ": " + m.pathInput.View()
	b.WriteString(pathBox.Render(pathContent) + "\n")

	// Search input with border
	searchBorderColor := lipgloss.Color("240")
	searchLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	if m.focus == focusSearch {
		searchBorderColor = lipgloss.Color("12") // Bright blue when focused
		searchLabelStyle = searchLabelStyle.Bold(true)
	}
	searchBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(searchBorderColor).
		Padding(0, 1).
		Width(m.width - 4)

	searchContent := searchLabelStyle.Render("SEARCH") + ": " + m.searchInput.View()
	b.WriteString(searchBox.Render(searchContent) + "\n\n")

	// Source file indicator
	if m.sourceFile != nil {
		sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s", m.sourceFile.Path)) + "\n\n")
	}

	// Determine layout based on preview mode
	var fileListWidth int
	var previewContent string

	if m.previewMode != previewHidden {
		// Split screen: 50% file list, 50% preview
		fileListWidth = m.width / 2
		previewContent = m.renderPreview()
	} else {
		// Full width for file list
		fileListWidth = m.width
	}

	// File list with border
	var fileListContent strings.Builder

	// File list border
	listBorderColor := lipgloss.Color("240")
	if m.focus == focusList {
		listBorderColor = lipgloss.Color("12") // Bright blue when focused
	}

	// File list header
	headerRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	if m.focus == focusList {
		headerRowStyle = headerRowStyle.Bold(true)
	}
	pathWidth := minInt(fileListWidth-10, 50) // Adjust based on available space
	fileListContent.WriteString(headerRowStyle.Render(fmt.Sprintf("%-*s %-10s %-15s", pathWidth, "FILE LIST", "SIZE", "MODIFIED")) + "\n")

	// File list
	maxVisible := m.height - 16 // Adjusted for borders
	if maxVisible < 1 {
		maxVisible = 1
	}

	start := m.viewport
	end := minInt(start+maxVisible, len(m.filteredFiles))

	for i := start; i < end; i++ {
		file := m.filteredFiles[i]
		cursor := " "
		if m.cursor == i {
			cursor = "▶" // More visible arrow
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
			style = style.Background(lipgloss.Color("240")).Bold(true)
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

		fileListContent.WriteString(style.Render(line) + "\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	fileListContent.WriteString(footerStyle.Render(fmt.Sprintf("\nShowing %d of %d files | Targets: %d",
		len(m.filteredFiles), len(m.files), len(m.selected))))

	// Wrap file list in border
	listBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(listBorderColor).
		Padding(0, 1).
		Width(fileListWidth - 4)

	renderedFileList := listBox.Render(fileListContent.String())

	// If preview is enabled, combine file list and preview side by side
	var bottomSection string
	if m.previewMode != previewHidden {
		bottomSection = lipgloss.JoinHorizontal(lipgloss.Top, renderedFileList, previewContent)
	} else {
		bottomSection = renderedFileList
	}

	b.WriteString(bottomSection)
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
		headerTitle = fmt.Sprintf(" Preview (diff): %s → %s ", m.sourceFile.Path, currentFile.Path)
	} else {
		// Show plain file content
		lines = strings.Split(string(content), "\n")
		headerTitle = fmt.Sprintf(" Preview (plain): %s ", currentFile.Path)
	}

	// Calculate preview dimensions
	previewWidth := m.width / 2
	// Match file list height to prevent overflow when joined horizontally
	// File list uses m.height - 16, so preview should use same or less
	previewHeight := m.height - 16

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
	end := minInt(start+previewHeight, len(lines))

	for i := start; i < end; i++ {
		line := lines[i]
		// Truncate long lines
		if len(line) > previewWidth-3 {
			line = line[:previewWidth-6] + "..."
		}

		// Color diff lines if in diff mode
		if m.previewMode == previewDiff && m.sourceFile != nil {
			lineStyle := contentStyle
			if line != "" {
				switch line[0] {
				case '+':
					lineStyle = contentStyle.Foreground(lipgloss.Color("34")) // Darker green for additions (better contrast)
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
		scrollInfo = fmt.Sprintf(" [%d-%d of %d lines] CTRL-U/D to scroll ", start+1, end, len(lines))
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

	// Header
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(headerStyle.Render("FileMirror - Confirm Copy & Git Workflow") + "\n\n")

	// Instructions
	instructStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"})
	hints := "FILE LIST: Review • GIT WORKFLOW: TAB to navigate • ENTER: copy & commit • CTRL-G: toggle git • ESC: cancel • q: quit"
	b.WriteString(instructStyle.Render(hints) + "\n\n")

	// Path and search (for context)
	pathBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width - 4)
	b.WriteString(pathBox.Render(fmt.Sprintf("PATH: %s", m.workDir)) + "\n")

	searchBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(m.width - 4)
	b.WriteString(searchBox.Render(fmt.Sprintf("SEARCH: %s", m.searchInput.Value())) + "\n\n")

	// Source indicator
	if m.sourceFile != nil {
		sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s", m.sourceFile.Path)) + "\n\n")
	}

	// Split panel layout
	fileListWidth := m.width / 3
	gitPanelWidth := (m.width * 2) / 3

	// Left panel: File list
	var fileListContent strings.Builder
	fileListContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("FILES TO SYNC") + "\n\n")

	if m.sourceFile != nil {
		fileListContent.WriteString("Source:\n")
		fileListContent.WriteString(fmt.Sprintf("▶ %s\n", m.sourceFile.Path))
		fileListContent.WriteString(fmt.Sprintf("  %s\n\n", formatSize(m.sourceFile.Size)))
	}

	targetCount := 0
	fileListContent.WriteString(fmt.Sprintf("Targets (%d):\n", len(m.selected)))
	for idx := range m.selected {
		if idx < len(m.filteredFiles) {
			file := m.filteredFiles[idx]
			fileListContent.WriteString(fmt.Sprintf("→ %s\n", file.Path))
			fileListContent.WriteString(fmt.Sprintf("  %s\n", formatSize(file.Size)))
			targetCount++
		}
	}

	fileListBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(fileListWidth - 4).
		Height(m.height - 18)

	renderedFileList := fileListBox.Render(fileListContent.String())

	// Right panel: Git workflow configuration
	var gitPanelContent strings.Builder
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	gitPanelContent.WriteString(titleStyle.Render("Git Workflow Configuration") + "\n\n")

	// Git enabled checkbox
	gitCheckbox := "[ ]"
	if m.gitEnabled {
		gitCheckbox = "[✓]"
	}
	gitEnabledStyle := lipgloss.NewStyle()
	if m.confirmFocus == focusGitEnabled {
		gitEnabledStyle = gitEnabledStyle.Background(lipgloss.Color("240")).Bold(true)
	}
	gitPanelContent.WriteString(gitEnabledStyle.Render(fmt.Sprintf("%s Create git commit", gitCheckbox)) + "\n\n")

	// Branch name
	branchLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	if m.confirmFocus == focusBranchName {
		branchLabelStyle = branchLabelStyle.Bold(true)
	}
	gitPanelContent.WriteString(branchLabelStyle.Render("Branch Name:") + "\n")

	branchBorderColor := lipgloss.Color("240")
	if m.confirmFocus == focusBranchName {
		branchBorderColor = lipgloss.Color("12")
	}
	branchBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(branchBorderColor).
		Padding(0, 1).
		Width(gitPanelWidth - 8)
	gitPanelContent.WriteString(branchBox.Render(m.branchNameInput.View()) + "\n\n")

	// Commit message
	commitLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	if m.confirmFocus == focusCommitMsg {
		commitLabelStyle = commitLabelStyle.Bold(true)
	}
	gitPanelContent.WriteString(commitLabelStyle.Render("Commit Message:") + "\n")

	commitBorderColor := lipgloss.Color("240")
	if m.confirmFocus == focusCommitMsg {
		commitBorderColor = lipgloss.Color("12")
	}
	commitBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(commitBorderColor).
		Padding(0, 1).
		Width(gitPanelWidth - 8)
	gitPanelContent.WriteString(commitBox.Render(m.commitMsgInput.View()) + "\n\n")

	// Push toggle
	pushCheckbox := "[ ]"
	if m.shouldPush {
		pushCheckbox = "[✓]"
	}
	pushStyle := lipgloss.NewStyle()
	if m.confirmFocus == focusPushToggle {
		pushStyle = pushStyle.Background(lipgloss.Color("240")).Bold(true)
	}
	gitPanelContent.WriteString(pushStyle.Render(fmt.Sprintf("%s Push to origin after commit", pushCheckbox)) + "\n\n")

	// Repository info
	if len(m.gitRepos) > 0 {
		gitPanelContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("Repository: %d git repos detected", len(m.gitRepos))) + "\n")
		for repo := range m.gitRepos {
			gitPanelContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(fmt.Sprintf("✓ %s", filepath.Base(repo))) + "\n")
		}
	} else {
		gitPanelContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✗ No git repositories detected") + "\n")
	}

	gitPanelContent.WriteString("\n")

	// Action buttons
	copyButtonStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2).
		Foreground(lipgloss.Color("10"))
	cancelButtonStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2)

	if m.confirmFocus == focusCopyButton {
		copyButtonStyle = copyButtonStyle.Background(lipgloss.Color("240")).Bold(true)
	}
	if m.confirmFocus == focusCancelButton {
		cancelButtonStyle = cancelButtonStyle.Background(lipgloss.Color("240")).Bold(true)
	}

	copyButton := copyButtonStyle.Render("Copy & Commit")
	cancelButton := cancelButtonStyle.Render("Cancel")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, copyButton, "  ", cancelButton)
	gitPanelContent.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(gitPanelWidth - 4).Render(buttons))

	gitPanelBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(gitPanelWidth - 4).
		Height(m.height - 18)

	renderedGitPanel := gitPanelBox.Render(gitPanelContent.String())

	// Combine panels side by side
	splitView := lipgloss.JoinHorizontal(lipgloss.Top, renderedFileList, renderedGitPanel)
	b.WriteString(splitView + "\n\n")

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	b.WriteString(footerStyle.Render("TAB: next field • ENTER: confirm • ESC: cancel • CTRL-G: toggle git"))

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

// initGitWorkflow initializes git workflow fields when entering confirm mode
func (m *model) initGitWorkflow() {
	// Initialize branch name input
	m.branchNameInput = textinput.New()
	m.branchNameInput.Placeholder = "Branch name..."
	m.branchNameInput.CharLimit = 100
	m.branchNameInput.Width = 50

	// Generate default branch name from source filename
	sourceName := "filesync"
	if m.sourceFile != nil {
		sourceName = normalizeBranchName(m.sourceFile.Path)
	}
	m.branchNameInput.SetValue(fmt.Sprintf("chore/filesync-%s", sourceName))

	// Initialize commit message textarea
	m.commitMsgInput = textarea.New()
	m.commitMsgInput.Placeholder = "Commit message..."
	m.commitMsgInput.CharLimit = 1000
	m.commitMsgInput.SetWidth(50)
	m.commitMsgInput.SetHeight(5)

	// Generate default commit message
	targetFiles := []string{}
	for idx := range m.selected {
		if idx < len(m.filteredFiles) {
			targetFiles = append(targetFiles, m.filteredFiles[idx].Path)
		}
	}

	commitMsg := fmt.Sprintf("chore: Sync %s from source\n\nSynchronized from %s\nTarget files:\n",
		filepath.Base(m.sourceFile.Path),
		m.sourceFile.Path)
	for _, target := range targetFiles {
		commitMsg += fmt.Sprintf("- %s\n", target)
	}
	m.commitMsgInput.SetValue(commitMsg)

	// Detect git repos for target files
	m.gitRepos = make(map[string][]string)
	for idx := range m.selected {
		if idx < len(m.filteredFiles) {
			targetPath := filepath.Join(m.workDir, m.filteredFiles[idx].Path)
			root, err := detectGitRoot(targetPath)
			if err == nil {
				m.gitRepos[root] = append(m.gitRepos[root], targetPath)
			}
		}
	}

	// Enable git by default if we have git repos
	m.gitEnabled = len(m.gitRepos) > 0
	m.shouldPush = false             // Safer default
	m.confirmFocus = focusCopyButton // Start on copy button
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
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
	maxLen := maxInt(len(sourceLines), len(targetLines))

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

// renderHelpOverlay renders the help modal overlay
func (m model) renderHelpOverlay() string {
	helpContent := `KEYBOARD SHORTCUTS

NAVIGATION
  TAB             Cycle focus forward: Path → Search → File List → Path
  Shift+TAB       Cycle focus backward

PATH INPUT
  Type            Edit working directory path
  CTRL-R          Reload files from current path

SEARCH INPUT
  Type            Filter files by pattern
                  Examples: *.go, config.json, component

FILE LIST
  ↑ / ↓           Navigate up/down
  k / j           Navigate up/down (vim-style)
  s               Mark current file as SOURCE
  SPACE           Toggle current file as TARGET
  ENTER           Proceed to confirmation (requires source + targets)

PREVIEW PANEL
  p / CTRL-P      Cycle preview modes: hidden → plain → diff → hidden
  CTRL-U / CTRL-D Scroll preview up/down
  Fn+↑ / Fn+↓     Scroll preview up/down (MacBook)
  PgUp / PgDn     Scroll preview up/down (if available)

GENERAL
  ?               Toggle this help screen
  q / CTRL-C      Quit program
  ESC             Close help / Cancel operation

WORKFLOW
  1. Navigate to FILE LIST (press TAB if needed)
  2. Use ↑/↓ or k/j to find your source file
  3. Press 's' to mark it as source
  4. Navigate to target files and press SPACE to select them
  5. Press ENTER to review, then 'y' to confirm sync

Press ESC or ? to close this help`

	// Calculate modal dimensions
	modalWidth := minInt(m.width-4, 80)
	modalHeight := minInt(m.height-4, 35)

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(modalWidth).
		MaxWidth(modalWidth)

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	// Content style
	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// Build modal content
	var modalContent strings.Builder
	modalContent.WriteString(titleStyle.Render("FILEMIRROR HELP") + "\n\n")
	modalContent.WriteString(contentStyle.Render(helpContent))

	modal := modalStyle.Render(modalContent.String())

	// Center the modal on screen
	verticalPadding := (m.height - modalHeight) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	var result strings.Builder
	for i := 0; i < verticalPadding; i++ {
		result.WriteString("\n")
	}

	// Center horizontally
	horizontalPadding := (m.width - modalWidth) / 2
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}

	modalLines := strings.Split(modal, "\n")
	for _, line := range modalLines {
		result.WriteString(strings.Repeat(" ", horizontalPadding) + line + "\n")
	}

	return result.String()
}
