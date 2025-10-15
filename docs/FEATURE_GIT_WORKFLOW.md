# Feature Plan: Git Workflow Integration for FileMirror

## Overview
Add a post-copy git workflow step that creates commits in **target project repositories** after file synchronization.

## Current Flow Understanding
1. User selects **source file** (in workDir)
2. User selects **target files** (in same workDir)
3. User confirms → files are copied
4. App quits

## New Flow (Merged Confirmation + Git Workflow)
Extend `modeConfirm` to include git workflow configuration in a split-panel UI

```
modeSelect → modeConfirm (with git workflow fields) → [COPY & COMMIT] → quit
                                  ↓
                            (git commit checkbox can be toggled off)
```

## Implementation Plan

### 1. Extend Confirmation Mode State (model.go)

**New Constants (no new mode needed):**
```go
type confirmFocus int
const (
    focusConfirmButtons confirmFocus = iota  // existing: Yes/No buttons
    focusGitEnabled                          // NEW: toggle git commit on/off
    focusBranchName                          // NEW: branch name input
    focusCommitMsg                           // NEW: commit message textarea
    focusPushToggle                          // NEW: push to origin toggle
)
```

**Add to model struct:**
```go
type model struct {
    // ... existing fields ...

    // Git workflow fields (integrated into modeConfirm)
    gitEnabled       bool                // toggle to enable/disable git workflow
    branchNameInput  textinput.Model
    commitMsgInput   textarea.Model      // multi-line for commit message
    shouldPush       bool                // toggle for auto-push
    confirmFocus     confirmFocus        // which input/button is focused
    gitRepos         map[string][]string // repo path -> list of changed files
}
```

### 2. Modify Confirmation Flow (model.go)

**When transitioning to `modeConfirm`:**
```go
case "enter":
    if m.mode == modeSelect && m.hasValidSelection() {
        m.mode = modeConfirm
        m.initGitWorkflow()  // NEW: initialize git fields when entering confirm mode
        return m, nil
    }
```

**New function `initGitWorkflow()`:**
- Detect git repos for each target file
- Group files by repository root
- Initialize default branch name: `chore/filesync-{source-filename-part}`
- Initialize default commit message template
- Set `gitEnabled = true` by default (user can toggle off)
- Set focus to confirm buttons initially

**In `updateConfirm()` when user confirms (ENTER on "Copy & Commit" button):**
```go
case "enter":
    // Perform the copy operation
    err := m.copySourceToTargets()
    if err != nil {
        m.err = err
        return m, nil
    }

    // NEW: If git is enabled, perform git operations
    if m.gitEnabled {
        err = m.performGitWorkflow()  // commit to all repos
        if err != nil {
            m.err = err
            return m, nil
        }
    }

    // Show success summary and quit
    return m, tea.Quit
```

### 3. Extend Confirmation View with Git Panel (model.go)

**Modify `viewConfirm()` to use split-panel layout:**

Display layout (see `docs/INTERFACE_EXAMPLES.md` for full visual):
```
╭─────────────────────────╮│ Git Workflow Configuration ───────────────┐
│ FILES TO SYNC          ││                                           │
│                        ││ [✓] Create git commit                     │
│ Source:                ││                                           │
│ ▶ model.go             ││ Branch Name:                              │
│   12.5 KB              ││ ╭───────────────────────────────────────╮ │
│                        ││ │ chore/filesync-model                  │ │
│ Targets (2):           ││ ╰───────────────────────────────────────╯ │
│ → scanner.go           ││                                           │
│   3.2 KB               ││ Commit Message:                           │
│ → git.go               ││ ╭───────────────────────────────────────╮ │
│   2.1 KB               ││ │ Chore: Update model.go                │ │
│                        ││ │                                       │ │
│                        ││ │ Synchronized from model.go            │ │
│                        ││ │ Targets: scanner.go, git.go           │ │
│                        ││ ╰───────────────────────────────────────╯ │
│                        ││                                           │
│                        ││ [ ] Push to origin after commit          │
│                        ││                                           │
│                        ││ Repository: /Users/jesper/src/filemirror │
│                        ││                                           │
│                        ││        ┌──────────────┐  ┌──────────┐    │
│                        ││        │ Copy & Commit│  │  Cancel  │    │
╰─────────────────────────╯│        └──────────────┘  └──────────┘    │
                           └───────────────────────────────────────────┘

TAB: next field • ENTER: confirm • ESC: cancel • CTRL-G: toggle git
```

**Key UI elements:**
- Left panel: Read-only file list (source + targets)
- Right panel: Git configuration fields
- Git enabled checkbox at top of right panel
- When git checkbox is OFF, gray out branch/commit/push fields

### 4. Git Operations Module (new file: gitworkflow.go)

**Core Functions:**

```go
package main

import (
    "os/exec"
    "path/filepath"
)

// detectGitRoot finds the git repository root for a file
func detectGitRoot(filePath string) (string, error) {
    dir := filepath.Dir(filePath)
    cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
    output, err := cmd.Output()
    // ... handle error and return root path
}

// createWorktreeAndBranch creates a new git worktree with branch
func createWorktreeAndBranch(repoPath, branchName, worktreePath string) error {
    // git worktree add -b <branchName> <worktreePath>
}

// commitChanges stages and commits files in the worktree
func commitChanges(worktreePath, message string, files []string) error {
    // cd to worktree
    // git add <files>
    // git commit -m "<message>"
}

// pushBranch pushes the branch to origin
func pushBranch(worktreePath, branchName string) error {
    // git push -u origin <branchName>
}

// cleanupWorktree removes the worktree
func cleanupWorktree(repoPath, worktreePath string) error {
    // git worktree remove <worktreePath>
}

// groupFilesByRepo groups changed files by their git repository
func groupFilesByRepo(files []string) (map[string][]string, error) {
    repos := make(map[string][]string)
    for _, file := range files {
        root, err := detectGitRoot(file)
        if err != nil {
            // File not in git repo, skip or warn
            continue
        }
        repos[root] = append(repos[root], file)
    }
    return repos, nil
}
```

### 5. Handle Multiple Projects

**Strategy: Same branch/commit for all repos**

Since targets may be in different repos:
1. Group target files by git repository root
2. Use the **same branch name** for all repos
3. Use the **same commit message** for all repos
4. Apply workflow to each repo sequentially
5. Show progress/results per repo

**Alternative (future enhancement):** Allow per-repo customization

### 6. User Input Fields

#### Branch Name Input
- **Type:** Single-line text input
- **Default:** `chore/filesync-{source-filename-part}`
  - Example: `chore/filesync-model`
- **Validation:** Check branch name doesn't already exist
- **Convention:** Follows `feature/*`, `fix/*` pattern from AGENT_NOTES.md
  - Could default to `sync/*` as a new category

#### Commit Message Input
- **Type:** Multi-line textarea
- **Default template:**
  ```
  Chore: Update {filename} from {source}

  Synchronized from {source-path}
  Target files:
  - {target-file-1}
  - {target-file-2}
  ```
- **Placeholder variables:**
  - `{filename}` - source file name
  - `{source-path}` - relative path to source
  - `{target-file-N}` - list of target files

#### Push Toggle
- **Type:** Boolean checkbox/toggle
- **Default:** `false` (NO push)
- **Rationale:** Safer default, follows AGENT_NOTES.md about explicit push
- **Display:** `[ ] Push to origin after commit` / `[x] Push to origin after commit`

### 7. Keyboard Shortcuts

**In `modeConfirm` (merged with git workflow):**
- **TAB** - Cycle focus: buttons → git enabled checkbox → branch input → commit input → push toggle → back to buttons
- **Shift+TAB** - Reverse cycle
- **CTRL-G** - Quick toggle for git enabled checkbox (from any focus)
- **ENTER** - Execute copy and git workflow (if git enabled)
- **ESC** - Cancel and return to file selection
- **CTRL-C** / **q** - Quit immediately without operations
- **Arrow keys** - Navigate in textarea when focused
- **SPACE** - Toggle checkboxes (git enabled, push toggle) when focused

### 8. Execution Flow

**When user presses ENTER in `modeConfirm` (on "Copy & Commit" button):**

1. **Copy phase**: Copy source file to all target files
2. **Git phase** (only if `gitEnabled == true):
   - For each repository with changed files:
     - Create worktree in temp location (e.g., `/tmp/fmr-worktree-{random-id}`)
     - Create new branch with specified name
     - Copy changed files to worktree
     - Stage and commit changes with user's commit message
     - Optionally push to origin if `shouldPush == true`
     - Cleanup worktree (always happens via defer)
3. **Summary phase**: Display results and exit

### 9. Error Handling

**Scenarios:**

1. **Target file not in git repo**
   - Show warning in UI: `✗ /path/to/project3 (not a git repo)`
   - Skip that file, continue with others

2. **Branch already exists**
   - Detect before creating worktree
   - if the only file changed in the branch is the same, then allow and reuse the branch. Overwrite the existing branch with the new contents
   - in other scenarios, show error, add to error list that can be shown in the end

3. **Git operations fail**
   - Show error message in UI
   - Options: [Retry] [Skip Repo] [Skip All]

4. **Worktree already exists**
   - Auto-cleanup stale worktrees
   - Or suggest different branch name

5. **Push fails (no remote, auth issues)**
   - Show error but mark commit as successful
   - Inform user to push manually

### 10. GPG Signing

**From AGENT_NOTES.md:**
  -- do nothing special, just respect user's git config

**Implementation:**
- Respect user's git config (no special handling needed)
- Git will automatically sign if configured
- If signing fails, show error and abort commit

### 11. Git Worktree Approach

**Recommendation:** Use worktree approach for safety

## Files to Create/Modify

### New Files:
1. **`gitworkflow.go`** - Git operations (worktree, commit, push)
2. **`FEATURE_GIT_WORKFLOW.md`** - This document

### Modified Files:
1. **`model.go`**
   - Add `confirmFocus` type with git-related focus states
   - Add git workflow fields to model struct (`gitEnabled`, `branchNameInput`, etc.)
   - Modify `viewConfirm()` to show split-panel layout with git fields
   - Modify `updateConfirm()` to handle git field navigation and execution
   - Add `initGitWorkflow()` to initialize git fields when entering confirm mode

2. **`git.go`**
   - Extend with helper functions if needed

3. **`go.mod`**
   - Add `github.com/charmbracelet/bubbles/textarea` dependency

## Design Decisions to Confirm

### 1. UI Architecture: Separate mode vs merged confirmation?
**Decision:** MERGED - Git workflow integrated into `modeConfirm` with split-panel UI
- **Rationale:** Single decision point, faster workflow, better context
- **Implementation:** Left panel shows files, right panel shows git configuration
- **Toggle:** `gitEnabled` checkbox allows skipping git operations

### 2. Scope: Same branch/message for all repos?
**Proposed:** YES - Same branch name and commit message for all target repos
- **Rationale:** Simpler UX, common use case
- **Future:** Could add "Advanced Mode" for per-repo customization

### 3. Worktree location?
**Proposed:** Temp directory (`/tmp/fmr-worktree-{random-id}`)
- **Rationale:** Clean, doesn't pollute user's workspace
- **Cleanup:** Auto-cleanup on success/failure with defer

### 4. Auto-cleanup worktrees?
**Proposed:** YES - Always cleanup after commit/push
- **Rationale:** Worktrees are implementation detail, user shouldn't see them
- **Exception:** On error, cleanup but log location for debugging

### 5. GPG signing?
**Proposed:** Respect user's git config
- **Rationale:** Git will auto-sign if configured in target repos
- **No special handling needed**

### 6. Multiple file commits?
**Proposed:** Single commit per repo with all changed files
- **Rationale:** Cleaner history, files are related (synced together)
- **Commit message lists all target files**
**Desiscion:** This program only handles sync of single source file to multiple targets. Therefore it is safe to have one single commit per repo and it will always contain only one file.

### 7. Branch naming convention?
**Proposed:** use `chore/*` category 

### 8. Optional vs Required?
**Decision:** Git workflow is OPTIONAL via checkbox toggle
- `gitEnabled` checkbox (default: ON) can be toggled off
- When disabled, only file copy is performed
- Faster than separate mode approach (no need to transition screens)

### 9. Show success message?
**Proposed:** YES - Show summary before quit
```
Git Workflow Complete!

✓ Committed to 2 repositories:
  - /path/to/project1 (branch: sync/model.go-...)
  - /path/to/project2 (branch: sync/model.go-...)

✓ Pushed to origin: YES

Press any key to exit...

The summary must stay in the terminal after exit for user reference.
```

## Implementation Phases

### Phase 1: Basic Structure
- [ ] Add `confirmFocus` type with git-related states
- [ ] Add git workflow fields to model (`gitEnabled`, `branchNameInput`, etc.)
- [ ] Modify `viewConfirm()` to show split-panel layout
- [ ] Add git enabled checkbox
- [ ] Add branch name input
- [ ] Add commit message textarea
- [ ] Add push toggle checkbox

### Phase 2: Git Operations
- [ ] Create `gitworkflow.go`
- [ ] Implement `detectGitRoot()`
- [ ] Implement `groupFilesByRepo()`
- [ ] Implement `createWorktreeAndBranch()`
- [ ] Implement `commitChanges()`
- [ ] Add error handling

### Phase 3: Push & Cleanup
- [ ] Add push toggle UI
- [ ] Implement `pushBranch()`
- [ ] Implement `cleanupWorktree()`
- [ ] Add defer cleanup
- [ ] Test error scenarios

### Phase 4: Polish
- [ ] Add success/failure messages
- [ ] Improve error messages
- [ ] Add branch name validation
- [ ] Add commit message templates
- [ ] Update help screen with new keyboard shortcuts
- [ ] Add CTRL-G quick toggle for git enabled

### Phase 5: Documentation
- [ ] Update README with git workflow feature
- [ ] Update help overlay (`?` key)
- [ ] Add examples to documentation
- [ ] Update AGENT_NOTES.md with `sync/*` convention

## Testing Scenarios

1. **Single repo, single file**
   - Source and target in same repo

2. **Multiple repos, multiple files**
   - Targets across 3 different repos

3. **Mixed: git + non-git targets**
   - Some targets in git repos, some not

4. **Branch already exists**
   - Target repo already has `sync/model.go-...` branch

5. **Git signing enabled**
   - Ensure commits are properly signed

6. **Push with no remote**
   - Target repo has no origin configured

7. **User skips workflow (ESC)**
   - Files copied but no git operations

8. **Errors during commit**
   - Ensure worktrees are cleaned up

## Future Enhancements

1. **Per-repo customization**
   - Allow different branch names per repo
   - Allow different commit messages per repo

2. **Template library**
   - Save/load commit message templates
   - Common patterns for different sync scenarios

3. **Conflict detection**
   - Check if target files have uncommitted changes
   - Warn before overwriting

4. **History/Undo**
   - Track previous sync operations
   - Allow reverting a sync

5. **Auto-PR creation**
   - Optionally create GitHub/GitLab PRs after push
   - Use `gh` CLI or API

6. **Dry-run mode**
   - Show what would be committed without actually doing it
   - Review changes before confirming

## Open Questions

1. Should we store git workflow preferences (branch template, auto-push) in a config file?
2. Should we validate commit message format (conventional commits)?
3. Should we support multiple commit strategies (one per file vs all in one)?
4. Should we integrate with GitHub CLI (`gh`) for PR creation?
5. Should we show a diff preview before committing?
