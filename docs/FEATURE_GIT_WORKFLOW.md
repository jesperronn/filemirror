# Feature Plan: Git Workflow Integration for FileMirror

## Overview
Add a post-copy git workflow step that creates commits in **target project repositories** after file synchronization.

## Current Flow Understanding
1. User selects **source file** (in workDir)
2. User selects **target files** (in same workDir)
3. User confirms → files are copied
4. App quits

## New Flow (After Copy Success)
Add a new mode: `modeGitWorkflow` between copy success and quit

```
modeSelect → modeConfirm → [COPY] → modeGitWorkflow → quit
                                  ↓
                            (if git workflow skipped/disabled)
```

## Implementation Plan

### 1. Add New Mode & State (model.go)

**New Constants:**
```go
const (
    modeSelect mode = iota
    modeConfirm
    modeGitWorkflow  // NEW
)

type gitFocus int
const (
    focusBranchName gitFocus = iota
    focusCommitMsg
    focusPushToggle
    focusActionButtons
)
```

**Add to model struct:**
```go
type model struct {
    // ... existing fields ...

    // Git workflow fields
    branchNameInput  textinput.Model
    commitMsgInput   textarea.Model      // multi-line for commit message
    shouldPush       bool                // toggle for auto-push
    gitFocus         gitFocus            // which input is focused
    copiedFiles      []string            // track what was successfully copied
    gitRepos         map[string][]string // repo path -> list of changed files
}
```

### 2. Modify Confirmation Flow (model.go)

**In `updateConfirm()` when user presses 'y':**
```go
case "y", "Y":
    // Perform the copy operation
    err := m.copySourceToTargets()
    if err != nil {
        m.err = err
        return m, nil
    }

    // NEW: Instead of tea.Quit, transition to git workflow
    m.mode = modeGitWorkflow
    m.initGitWorkflow()  // new function
    return m, nil
```

**New function `initGitWorkflow()`:**
- Detect git repos for each copied target file
- Group files by repository root
- Initialize default branch name: `chore/filesync-{source-filename-part}`
- Initialize default commit message template
- Set focus to branch name input

### 3. Create Git Workflow Screen (model.go)

**New function: `viewGitWorkflow()`**

Display layout:
```
╭─────────────────────────────────────────────────╮
│ Git Workflow (Optional)                         │
├─────────────────────────────────────────────────┤
│                                                 │
│ Branch Name:                                    │
│ ╭─────────────────────────────────────────────╮ │
│ │ chore/filesync-<filename-part>              │ │
│ ╰─────────────────────────────────────────────╯ │
│                                                 │
│ Commit Message:                                 │
│ ╭─────────────────────────────────────────────╮ │
│ │ Sync: Update model.go from source           │ │
│ │                                             │ │
│ │ Synchronized from ../source/model.go        │ │
│ │ - path/to/target1/model.go                  │ │
│ │ - path/to/target2/model.go                  │ │
│ ╰─────────────────────────────────────────────╯ │
│                                                 │
│ [ ] Push to origin after commit                │
│                                                 │
│ Target Repositories:                            │
│ ✓ /path/to/project1 (1 file)                   │
│ ✓ /path/to/project2 (2 files)                  │
│ ✗ /path/to/project3 (not a git repo)           │
│                                                 │
│ [Commit & Continue] [Skip]                      │
╰─────────────────────────────────────────────────╯

where <filename-part> is derived from source filename (e.g. model.go => sync/filesync-model)
```

TAB: next field • ENTER: commit • ESC: skip • q: quit
```

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

**In `modeGitWorkflow`:**
- **TAB** - Cycle focus: branch input → commit input → push toggle → action buttons
- **Shift+TAB** - Reverse cycle
- **ENTER** - Execute git workflow (commit, optionally push) -- IMPORTANT: No ENTER in textarea
- **ESC** - Skip git workflow entirely, quit app
- **CTRL-C** / **q** - Quit immediately without git operations
- **Arrow keys** - Navigate in textarea when focused

### 8. Execution Flow

**When user presses ENTER in `modeGitWorkflow`:**

For each repository with changed files:
- Create worktree in temp location (e.g., `/tmp/fmr-worktree-{random-id}`)
- Create new branch with specified name
- Copy changed files to worktree
- Stage and commit changes with user's commit message
- Optionally push to origin if toggle enabled
- Cleanup worktree (always happens via defer)

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
   - Add `modeGitWorkflow` constant
   - Add git workflow fields to model struct
   - Add `viewGitWorkflow()` function
   - Add `updateGitWorkflow()` function
   - Modify `updateConfirm()` to transition to git workflow

2. **`git.go`**
   - Extend with helper functions if needed

3. **`go.mod`**
   - Add `github.com/charmbracelet/bubbles/textarea` dependency

## Design Decisions to Confirm

### 1. Scope: Same branch/message for all repos?
**Proposed:** YES - Same branch name and commit message for all target repos
- **Rationale:** Simpler UX, common use case
- **Future:** Could add "Advanced Mode" for per-repo customization

### 2. Worktree location?
**Proposed:** Temp directory (`/tmp/fmr-worktree-{random-id}`)
- **Rationale:** Clean, doesn't pollute user's workspace
- **Cleanup:** Auto-cleanup on success/failure with defer

### 3. Auto-cleanup worktrees?
**Proposed:** YES - Always cleanup after commit/push
- **Rationale:** Worktrees are implementation detail, user shouldn't see them
- **Exception:** On error, cleanup but log location for debugging

### 4. GPG signing?
**Proposed:** Respect user's git config
- **Rationale:** Git will auto-sign if configured in target repos
- **No special handling needed**

### 5. Multiple file commits?
**Proposed:** Single commit per repo with all changed files
- **Rationale:** Cleaner history, files are related (synced together)
- **Commit message lists all target files**
**Desiscion:** This program only handles sync of single source file to multiple targets. Therefore it is safe to have one single commit per repo and it will always contain only one file.

### 6. Branch naming convention?
**Proposed:** use `chore/*` category 

### 7. Optional vs Required?
**Proposed:** Git workflow is OPTIONAL (can skip)
- ESC key skips workflow and quits
- Useful when user wants to commit manually later
- Or when targets aren't in git repos

### 8. Show success message?
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
- [ ] Add `modeGitWorkflow` mode
- [ ] Add git workflow fields to model
- [ ] Create basic `viewGitWorkflow()` UI
- [ ] Add branch name input
- [ ] Add commit message textarea
- [ ] Add skip/continue options

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
- [ ] Update help screen with new workflow

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
