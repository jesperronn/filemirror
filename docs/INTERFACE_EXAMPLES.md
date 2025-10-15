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

## Proposed NEW Interface (modeGitWorkflow)

### Git Workflow Screen (After Confirmation)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│ Git Workflow (Optional)                                                             │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│ Files successfully copied! Configure git commit settings below.                     │
│                                                                                     │
│ Branch Name:                                                                        │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ chore/filesync-model█                                                           │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│                                                                                     │
│ Commit Message: (multi-line, use arrow keys to navigate)                           │
│ ╭─────────────────────────────────────────────────────────────────────────────────╮ │
│ │ Chore: Update model.go from source                                              │ │
│ │                                                                                 │ │
│ │ Synchronized from /Users/jesper/src/filemirror/model.go                        │ │
│ │ Target files:                                                                   │ │
│ │ - /Users/jesper/projects/app1/model.go                                          │ │
│ │ - /Users/jesper/projects/app2/model.go                                          │ │
│ │                                                                                 │ │
│ ╰─────────────────────────────────────────────────────────────────────────────────╯ │
│                                                                                     │
│ [ ] Push to origin after commit                                                    │
│                                                                                     │
│ Target Repositories:                                                                │
│ ✓ /Users/jesper/projects/app1 (1 file, branch: main)                               │
│ ✓ /Users/jesper/projects/app2 (1 file, branch: main)                               │
│ ✗ /tmp/not-a-repo/model.go (not in a git repository)                               │
│                                                                                     │
│                    ┌────────────────────┐  ┌──────┐                                │
│                    │ Commit & Continue  │  │ Skip │                                │
│                    └────────────────────┘  └──────┘                                │
│                                                                                     │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ TAB: next field • ENTER: commit & continue • ESC: skip workflow • q: quit          │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

**Key Features:**
- Branch name input (single-line)
- Commit message textarea (multi-line, scrollable)
- Push toggle checkbox
- Repository status list showing:
  - ✓ Valid git repositories
  - ✗ Non-git targets that will be skipped
- Action buttons for commit or skip
- Optional workflow (can press ESC to skip entirely)

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

### With Git Workflow Feature
```
modeSelect → modeConfirm → [COPY] → modeGitWorkflow → [GIT OPS] → quit with summary
                                           ↓
                                      [ESC to skip]
```

**Steps:**
1. Select source file and target files
2. Confirm the copy operation
3. Files are copied successfully
4. **NEW**: Git workflow screen appears
5. **NEW**: User can configure branch name, commit message, and push option
6. **NEW**: Commits are created in target repositories using worktrees
7. **NEW**: Summary stays visible in terminal after exit

## Key Differences Summary

### Current State
- **Modes**: 2 (modeSelect, modeConfirm)
- **Input fields**: 3 (Path, Search, File List)
- **Preview**: Toggleable side-by-side view
- **After confirmation**: Immediate quit
- **Git integration**: None

### After Git Workflow Feature
- **Modes**: 3 (modeSelect, modeConfirm, **modeGitWorkflow**)
- **Total input fields**: 8 (existing 3 + branch name + commit message + push toggle + 2 buttons)
- **Preview**: Same as before
- **After confirmation**: Optional git workflow
- **Git integration**: Full support for multi-repo commits
  - Automatic worktree creation
  - Branch creation in target repos
  - Commit with custom message
  - Optional push to origin
  - Automatic cleanup

## Design Principles

1. **Optional workflow**: Git workflow can be skipped with ESC key
2. **Multi-repository support**: Handles targets across different repos
3. **Safety first**: Uses worktrees to avoid disrupting working directories
4. **Clear feedback**: Progress indicators and persistent summary
5. **Actionable output**: Summary includes next steps for the user
6. **Consistent styling**: Maintains FileMirror's visual language
