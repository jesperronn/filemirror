# FileMirror Interface Examples

This document shows visual examples of the current interface and the proposed git workflow interface.

## Current Interface (modeSelect)

### WITHOUT Preview Window (previewMode = previewHidden)

```
┌───────────────────────────────────────────────────────────────────────────┐
│ FileMirror - File Synchronization Tool                                   │
│                                                                           │
│ FILE LIST: ↑/↓ or k/j: navigate • s: set source • SPACE: toggle target • │
│ ENTER: confirm sync • p/CTRL-P: preview plain • TAB: next • ?: help •    │
│ q: quit                                                                   │
│                                                                           │
│ ╭───────────────────────────────────────────────────────────────────────╮ │
│ │ PATH: /Users/jesper/src/filemirror                                    │ │
│ ╰───────────────────────────────────────────────────────────────────────╯ │
│ ╭───────────────────────────────────────────────────────────────────────╮ │
│ │ SEARCH: *.go                                                          │ │
│ ╰───────────────────────────────────────────────────────────────────────╯ │
│                                                                           │
│ Source: model.go                                                          │
│                                                                           │
│ ╭───────────────────────────────────────────────────────────────────────╮ │
│ │ FILE LIST                SIZE       MODIFIED                          │ │
│ │ ▶[S] model.go            12.5 KB    2025-01-15 14:30                  │ │
│ │  [T] scanner.go          3.2 KB     2025-01-15 14:25                  │ │
│ │  [T] git.go              2.1 KB     2025-01-15 14:20                  │ │
│ │  [ ] main.go             1.8 KB     2025-01-15 14:15                  │ │
│ │  [ ] fileops.go          4.5 KB     2025-01-15 14:10                  │ │
│ │  [ ] diff.go             2.8 KB     2025-01-15 14:05                  │ │
│ │                                                                        │ │
│ │ Showing 6 of 15 files | Targets: 2                                    │ │
│ ╰───────────────────────────────────────────────────────────────────────╯ │
└───────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Full-width file list
- Path and search inputs at top
- Source file indicator (S marker)
- Target selection with T markers
- Cursor navigation with arrow indicator (▶)

### WITH Preview Window (previewMode = previewPlain or previewDiff)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│ FileMirror - File Synchronization Tool                                             │
│                                                                                     │
│ FILE LIST: ↑/↓ or k/j: navigate • s: set source • SPACE: toggle target •           │
│ p/CTRL-P: preview diff • CTRL-U/D: scroll preview • ENTER: confirm sync •           │
│ TAB: next • ?: help • q: quit                                                       │
│                                                                                     │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ PATH: /Users/jesper/src/filemirror                                              │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ SEARCH: *.go                                                                    │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│                                                                                     │
│ Source: model.go                                                                    │
│                                                                                     │
│ ╭─────────────────────╮│ Preview (diff): model.go → scanner.go ──────────────────┐ │
│ │ FILE LIST          ││                                                           │ │
│ │ SIZE      MODIFIED ││ @@ Source: model.go → Target @@                          │ │
│ │ ▶[S] model.go      ││  package main                                             │ │
│ │  12.5 KB  14:30    ││                                                           │ │
│ │  [T] scanner.go    ││  import (                                                 │ │
│ │  3.2 KB   14:25    ││ -    "github.com/charmbracelet/bubbles/textinput"        │ │
│ │  [T] git.go        ││ +    "fmt"                                                │ │
│ │  2.1 KB   14:20    ││ +    "io/fs"                                              │ │
│ │  [ ] main.go       ││ +    "os"                                                 │ │
│ │  1.8 KB   14:15    ││ +    "path/filepath"                                      │ │
│ │  [ ] fileops.go    ││                                                           │ │
│ │  4.5 KB   14:10    ││  type FileInfo struct {                                   │ │
│ │  [ ] diff.go       ││      Path     string                                      │ │
│ │  2.8 KB   14:05    ││ -    Branch   string                                      │ │
│ │                    ││ +    Size     int64                                        │ │
│ │ Showing 6 of 15    ││ +    Modified time.Time                                   │ │
│ │ files | Targets: 2 ││  }                                                        │ │
│ ╰─────────────────────╯│                                                           │ │
│                        │ [1-15 of 142 lines] CTRL-U/D to scroll ──────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Split-screen layout (50/50)
- File list on left, preview on right
- Preview modes: plain text or diff view
- Colored diff output (green for additions, red for deletions)
- Scrollable preview with position indicator

## Proposed MERGED Interface (modeConfirm with Git Workflow)

### Confirmation Screen with Integrated Git Workflow

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│ FileMirror - Confirm Copy & Git Workflow                                           │
│                                                                                     │
│ FILE LIST: Review • GIT WORKFLOW: TAB to navigate • ENTER: copy & commit •         │
│ CTRL-G: toggle git • ESC: cancel • q: quit                                         │
│                                                                                     │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ PATH: /Users/jesper/src/filemirror                                              │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ SEARCH: *.go                                                                    │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│                                                                                     │
│ Source: model.go                                                                    │
│                                                                                     │
│ ╭───────────────────────────╮│ Git Workflow Configuration ──────────────────────┐ │
│ │ FILES TO SYNC            ││                                                   │ │
│ │                          ││ [✓] Create git commit                             │ │
│ │ Source:                  ││                                                   │ │
│ │ ▶ model.go               ││ Branch Name:                                      │ │
│ │   12.5 KB                ││ ╭───────────────────────────────────────────────╮ │ │
│ │                          ││ │ chore/filesync-model█                         │ │ │
│ │ Targets (2):             ││ ╰───────────────────────────────────────────────╯ │ │
│ │ → scanner.go             ││                                                   │ │
│ │   3.2 KB                 ││ Commit Message:                                   │ │
│ │ → git.go                 ││ ╭───────────────────────────────────────────────╮ │ │
│ │   2.1 KB                 ││ │ Chore: Update model.go                        │ │ │
│ │                          ││ │                                               │ │ │
│ │                          ││ │ Synchronized from model.go                    │ │ │
│ │                          ││ │ Targets: scanner.go, git.go                   │ │ │
│ │                          ││ │                                               │ │ │
│ │                          ││ ╰───────────────────────────────────────────────╯ │ │
│ │                          ││                                                   │ │
│ │                          ││ [ ] Push to origin after commit                  │ │
│ │                          ││                                                   │ │
│ │                          ││ Repository: /Users/jesper/src/filemirror         │ │
│ │                          ││                                                   │ │
│ │                          ││        ┌──────────────┐  ┌──────────┐            │ │
│ │                          ││        │ Copy & Commit│  │  Cancel  │            │ │
│ ╰───────────────────────────╯│        └──────────────┘  └──────────┘            │ │
│                              │                                                   │ │
│                              └───────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ TAB: next field • ENTER: confirm • ESC: cancel • CTRL-G: toggle git                │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- **Split-panel layout**: File list on left, git configuration on right
- **Git enabled checkbox**: Toggle git workflow on/off without leaving screen
- **Branch name input**: Single-line text input with default value
- **Commit message textarea**: Multi-line, scrollable textarea
- **Push toggle checkbox**: Enable/disable auto-push to origin
- **Repository info**: Shows detected git repository
- **Action buttons**: "Copy & Commit" performs both operations, "Cancel" returns to selection
- **Quick toggle**: CTRL-G to quickly enable/disable git from any field

### During Commit Execution

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│ Git Workflow - Committing Changes...                                                │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│ ✓ /Users/jesper/projects/app1                                                       │
│   ├─ Created worktree: /tmp/fmr-worktree-abc123                                    │
│   ├─ Created branch: chore/filesync-model                                          │
│   ├─ Committed changes (1 file)                                                    │
│   └─ Cleaned up worktree                                                           │
│                                                                                     │
│ ⟳ /Users/jesper/projects/app2 (in progress...)                                     │
│   ├─ Created worktree: /tmp/fmr-worktree-def456                                    │
│   └─ Creating branch: chore/filesync-model                                         │
│                                                                                     │
│ ⊘ /tmp/not-a-repo/model.go (skipped - not a git repository)                        │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Real-time progress indicator
- Tree-style display of git operations
- Different symbols for status:
  - ✓ Completed successfully
  - ⟳ In progress
  - ⊘ Skipped

### Success Summary (stays visible after program exits)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│ Git Workflow Complete!                                                              │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│ ✓ Successfully committed to 2 repositories:                                         │
│                                                                                     │
│   • /Users/jesper/projects/app1                                                     │
│     Branch: chore/filesync-model                                                    │
│     Files:  1 file committed                                                        │
│                                                                                     │
│   • /Users/jesper/projects/app2                                                     │
│     Branch: chore/filesync-model                                                    │
│     Files:  1 file committed                                                        │
│                                                                                     │
│ ✓ Pushed to origin: NO (you can push manually later)                               │
│                                                                                     │
│ ⊘ Skipped 1 target (not in git repository):                                        │
│   • /tmp/not-a-repo/model.go                                                        │
│                                                                                     │
│ Next steps:                                                                         │
│   • Review commits: git log -1 in each repository                                   │
│   • Push to remote: git push -u origin chore/filesync-model                        │
│   • Create pull requests on GitHub                                                  │
│                                                                                     │
│ Press any key to exit...                                                            │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Comprehensive summary of all operations
- Grouped by success/failure/skipped
- Actionable next steps for the user
- Stays visible in terminal after program exits
- Clear indication of what was/wasn't pushed

## Navigation Flow Comparison

### Current Flow
```
modeSelect → modeConfirm → [COPY] → quit
```

**Steps:**
1. Select source file and target files
2. Confirm the copy operation
3. Files are copied
4. Program quits immediately

### With Git Workflow Feature (MERGED UI)
```
modeSelect → modeConfirm (with git fields) → [COPY & COMMIT] → quit with summary
                              ↓
                    [git checkbox toggleable]
```

**Steps:**
1. Select source file and target files
2. **MERGED**: Confirm screen shows files AND git workflow configuration
3. **NEW**: User can toggle git checkbox, configure branch name, commit message, and push option
4. **NEW**: Press ENTER to copy files and (optionally) create git commits
5. **NEW**: Summary stays visible in terminal after exit

**Key Difference:** No separate git workflow mode - everything happens in one enhanced confirmation screen

## Key Differences Summary

### Current State
- **Modes**: 2 (modeSelect, modeConfirm)
- **Input fields**: 3 (Path, Search, File List)
- **Confirmation view**: Simple yes/no dialog
- **Preview**: Toggleable side-by-side view (in modeSelect only)
- **After confirmation**: Immediate quit
- **Git integration**: None

### After Git Workflow Feature (MERGED DESIGN)
- **Modes**: 2 (modeSelect, **enhanced modeConfirm**)
- **Total input fields**: 7 (path + search in select mode, then git enabled checkbox + branch name + commit message + push toggle + 2 buttons in confirm mode)
- **Confirmation view**: Split-panel with files on left, git config on right
- **Preview**: Toggleable in modeSelect (hidden in modeConfirm to make room for git fields)
- **After confirmation**: Copy and optional git commit in single operation
- **Git integration**: Full support for multi-repo commits
  - Toggleable via checkbox (default: ON)
  - Automatic worktree creation
  - Branch creation in target repos
  - Commit with custom message
  - Optional push to origin
  - Automatic cleanup
- **Benefits of merged approach**:
  - Single decision point (review files + configure git together)
  - Faster workflow (no mode transition)
  - Better context (see files while writing commit message)
  - Easy to disable git (just toggle checkbox)

## Design Principles

1. **Merged UI**: Git workflow integrated into confirmation screen (no separate mode)
2. **Optional workflow**: Git can be toggled off via checkbox or CTRL-G shortcut
3. **Single decision point**: Review files and configure git in one screen
4. **Multi-repository support**: Handles targets across different repos
5. **Safety first**: Uses worktrees to avoid disrupting working directories
6. **Clear feedback**: Progress indicators and persistent summary
7. **Actionable output**: Summary includes next steps for the user
8. **Consistent styling**: Maintains FileMirror's visual language with split-panel approach
