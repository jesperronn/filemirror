# Test Coverage Analysis

**Date:** 15 October 2025  
**Overall Coverage:** 14.1%  
**Status:** ❌ Needs Significant Improvement

## Executive Summary

The project currently has minimal test coverage (14.1%), with critical business logic completely untested. The most complex and error-prone components—git workflow operations and TUI state management—have 0% coverage.

## Coverage by File

| File | Coverage | Status | Priority |
|------|----------|--------|----------|
| `gitworkflow.go` | **0.0%** | ❌ No tests exist | 🔴 CRITICAL |
| `model.go` | **~5%** | ❌ Minimal coverage | 🔴 CRITICAL |
| `run.go` | **0.0%** | ❌ No tests exist | 🟡 HIGH |
| `cmd/filemirror/main.go` | **0.0%** | ❌ No tests exist | 🟢 LOW |
| `cmd/fmr/main.go` | **0.0%** | ❌ No tests exist | 🟢 LOW |
| `fileops.go` | **69.2%** | ⚠️ Partial coverage | 🟢 LOW |
| `scanner.go` | **71.1%** | ⚠️ Partial coverage | 🟢 LOW |
| `git.go` | **55.6%** | ⚠️ Partial coverage | 🟡 MEDIUM |
| `branchname.go` | **100%** | ✅ Full coverage | ✅ DONE |

## Top 20 Methods Missing Tests

### 🔴 CRITICAL PRIORITY - Complex Business Logic (0% Coverage)

#### 1. `updateSelect` (model.go:165)
- **Complexity:** ~250 lines
- **Impact:** CRITICAL - Core TUI interaction handler
- **Why Critical:**
  - Manages all keyboard input in selection mode
  - Handles focus switching between path/search/list
  - Controls file selection, source marking, preview modes
  - Multiple nested switch statements with state transitions
  - Any bug here breaks the entire user experience
- **Testing Strategy:**
  - Table-driven tests with key sequences
  - Test all focus states and transitions
  - Mock bubbletea messages
  - Test edge cases (empty lists, boundary conditions)

#### 2. `updateConfirm` (model.go:420)
- **Complexity:** ~250 lines
- **Impact:** CRITICAL - Confirmation & Git workflow UI
- **Why Critical:**
  - Handles confirmation screen with multiple interactive elements
  - Complex focus management (textarea, textinput, checkboxes)
  - Git workflow integration and validation
  - Tab/Shift-Tab navigation with wraparound
  - Executes actual file operations and git commands
- **Testing Strategy:**
  - Test each focus state independently
  - Verify navigation between all fields
  - Test git workflow toggling
  - Mock git operations
  - Test commit/cancel actions

#### 3. `performGitWorkflow` (gitworkflow.go:152)
- **Complexity:** ~20 lines, high cyclomatic complexity
- **Impact:** HIGH - Orchestrates multi-repo git operations
- **Why Critical:**
  - Coordinates git operations across multiple repositories
  - Aggregates errors from multiple sources
  - Partial success handling (some repos succeed, others fail)
  - Critical for git workflow feature correctness
- **Testing Strategy:**
  - Mock git command execution
  - Test single repo success/failure
  - Test multiple repos with mixed results
  - Verify error aggregation
  - Test empty repo map

#### 4. `processRepo` (gitworkflow.go:170)
- **Complexity:** ~40 lines with cleanup
- **Impact:** HIGH - Single repo git operations
- **Why Critical:**
  - Manages worktree lifecycle (create, use, cleanup)
  - Deferred cleanup with error handling complexity
  - Multiple failure points (worktree, copy, commit, push)
  - Cleanup must happen even on error
  - Resource leak risk if cleanup fails
- **Testing Strategy:**
  - Test success path
  - Test each failure point independently
  - Verify cleanup happens on error
  - Mock file operations and git commands
  - Test push optional behavior

#### 5. `createWorktreeAndBranch` (gitworkflow.go:42)
- **Complexity:** ~25 lines
- **Impact:** HIGH - Git worktree management
- **Why Critical:**
  - Checks branch existence before creating
  - Conditional command execution (create vs checkout)
  - Generates unique worktree paths
  - Foundation for entire git workflow
- **Testing Strategy:**
  - Test new branch creation
  - Test existing branch checkout
  - Verify worktree path uniqueness
  - Mock git commands
  - Test error handling

#### 6. `commitChanges` (gitworkflow.go:68)
- **Complexity:** ~30 lines
- **Impact:** HIGH - Git commit operations
- **Why Critical:**
  - Stages multiple files individually
  - Relative path calculations (error-prone)
  - Handles "nothing to commit" special case
  - Multiple command executions with error paths
- **Testing Strategy:**
  - Test single and multiple file staging
  - Test "nothing to commit" scenario
  - Test relative path edge cases
  - Verify error messages include file context
  - Mock git commands

#### 7. `viewSelect` (model.go:682)
- **Complexity:** ~200 lines
- **Impact:** HIGH - Main view rendering
- **Why Critical:**
  - Renders entire selection screen UI
  - Multiple panel layout calculations
  - Status indicators and file list
  - Preview panel integration
  - Styling and formatting logic
- **Testing Strategy:**
  - Test with empty file list
  - Test with selected files
  - Test with source file set
  - Verify output contains expected sections
  - Test different terminal sizes

#### 8. `renderPreview` (model.go:878)
- **Complexity:** ~130 lines
- **Impact:** MEDIUM-HIGH - Preview display
- **Why Critical:**
  - Handles both plain and diff preview modes
  - File reading with error handling
  - Scrolling and viewport calculations
  - Line truncation for long lines
  - Diff colorization in diff mode
- **Testing Strategy:**
  - Test plain mode rendering
  - Test diff mode rendering
  - Test with files of different sizes
  - Test scrolling boundaries
  - Test error handling (file not found, read errors)
  - Verify color codes in diff mode

#### 9. `initGitWorkflow` (model.go:1223)
- **Complexity:** ~50 lines
- **Impact:** MEDIUM-HIGH - Git workflow initialization
- **Why Critical:**
  - Initializes all git workflow UI components
  - Generates default branch names from filenames
  - Creates commit message templates
  - Detects git repos for all selected files
  - Sets default values for workflow
- **Testing Strategy:**
  - Test branch name generation with various filenames
  - Test commit message templating
  - Test git repo detection (git and non-git files)
  - Verify default values
  - Test with different file selections

#### 10. `generateDiff` (model.go:1313)
- **Complexity:** ~40 lines
- **Impact:** MEDIUM - Diff generation
- **Why Critical:**
  - Line-by-line file comparison
  - Handles files of different lengths
  - Generates unified diff format
  - Used in preview mode
- **Testing Strategy:**
  - Test identical files (no diff)
  - Test completely different files
  - Test files with shared lines
  - Test files of different lengths
  - Verify diff format (+/- prefixes)

### 🟡 MEDIUM PRIORITY - Infrastructure & Integration

#### 11. `copySourceToTargets` (model.go:1205)
- **Complexity:** ~20 lines
- **Impact:** MEDIUM - File copy operation
- **Why Important:**
  - Core functionality: copies source to all targets
  - Iterates through selected files
  - Error handling and propagation
  - Depends on copyFile (which has coverage)
- **Testing Strategy:**
  - Test with single target
  - Test with multiple targets
  - Test error propagation
  - Mock copyFile failures

#### 12. `pushBranch` (gitworkflow.go:97)
- **Complexity:** ~10 lines
- **Impact:** MEDIUM - Git push operation
- **Why Important:**
  - Pushes branch to remote with upstream tracking
  - Command execution and error handling
  - Part of git workflow feature
- **Testing Strategy:**
  - Test successful push
  - Test push failure (network, auth, etc.)
  - Verify upstream tracking flag
  - Mock git commands

#### 13. `copyFileToWorktree` (gitworkflow.go:122)
- **Complexity:** ~30 lines
- **Impact:** MEDIUM - Worktree file operations
- **Why Important:**
  - Copies files to worktree preserving structure
  - Relative path calculation from repo root
  - Directory creation as needed
  - Multiple error paths
- **Testing Strategy:**
  - Test with flat structure
  - Test with nested directories
  - Test relative path edge cases
  - Test directory creation
  - Verify file content and permissions

#### 14. `cleanupWorktree` (gitworkflow.go:108)
- **Complexity:** ~15 lines
- **Impact:** MEDIUM - Resource cleanup
- **Why Important:**
  - Removes git worktree after use
  - Fallback to manual cleanup
  - Prevents tmp directory buildup
  - Critical for resource management
- **Testing Strategy:**
  - Test successful git cleanup
  - Test git cleanup failure + manual success
  - Test both cleanup methods failing
  - Verify tmp directory removal

#### 15. `detectGitRoot` (gitworkflow.go:14)
- **Complexity:** ~15 lines
- **Impact:** MEDIUM - Git detection
- **Why Important:**
  - Finds git repo root for a file path
  - Path resolution and validation
  - Error handling for non-git directories
  - Foundation for git workflow
- **Testing Strategy:**
  - Test with file in git repo
  - Test with file outside git repo
  - Test with nested git repos
  - Test relative and absolute paths

#### 16. `Update` (model.go:137)
- **Complexity:** ~30 lines
- **Impact:** MEDIUM - Main update dispatcher
- **Why Important:**
  - Routes all messages to appropriate handlers
  - Window size handling
  - Scan completion handling
  - Mode switching logic
- **Testing Strategy:**
  - Test WindowSizeMsg handling
  - Test scanCompleteMsg with success/error
  - Test KeyMsg routing to both modes
  - Test unknown message types

#### 17. `View` (model.go:659)
- **Complexity:** ~20 lines
- **Impact:** MEDIUM - Main view dispatcher
- **Why Important:**
  - Routes to appropriate view based on mode
  - Error display
  - Help overlay integration
  - Top-level view logic
- **Testing Strategy:**
  - Test error display
  - Test mode switching (select, confirm)
  - Test help overlay rendering
  - Test default mode

#### 18. `viewConfirm` (model.go:1020)
- **Complexity:** ~180 lines
- **Impact:** MEDIUM-LOW - Confirmation view
- **Why Important:**
  - Renders entire confirmation screen
  - Split panel layout (file list + git config)
  - Git workflow UI components
  - Multiple input field rendering
- **Testing Strategy:**
  - Test with git enabled/disabled
  - Test with different file selections
  - Verify all UI elements present
  - Test different terminal sizes

### 🟢 LOWER PRIORITY - Simpler Functions

#### 19. `generateWorktreeID` (gitworkflow.go:32)
- **Complexity:** ~10 lines
- **Impact:** LOW - ID generation
- **Why Test:**
  - Random ID generation
  - Fallback to PID on error
  - Simple but should verify both paths
- **Testing Strategy:**
  - Test successful random generation
  - Test fallback logic
  - Verify ID format

#### 20. `Init` (model.go:127)
- **Complexity:** ~10 lines
- **Impact:** LOW - Initialization
- **Why Test:**
  - Bubbletea Init method
  - Returns batch of commands
  - Important for startup flow
- **Testing Strategy:**
  - Test command batch creation
  - Verify initial scan triggered
  - Test with empty/populated initial state

## Detailed Function Coverage Report

```
Function                                            Coverage    Lines
─────────────────────────────────────────────────────────────────────
cmd/filemirror/main.go:main                         0.0%        5
cmd/fmr/main.go:main                                0.0%        5
gitworkflow.go:detectGitRoot                        0.0%        14
gitworkflow.go:generateWorktreeID                   0.0%        10
gitworkflow.go:createWorktreeAndBranch              0.0%        26
gitworkflow.go:commitChanges                        0.0%        29
gitworkflow.go:pushBranch                           0.0%        11
gitworkflow.go:cleanupWorktree                      0.0%        14
gitworkflow.go:copyFileToWorktree                   0.0%        30
gitworkflow.go:performGitWorkflow                   0.0%        18
gitworkflow.go:processRepo                          0.0%        35
model.go:Init                                       0.0%        10
model.go:Update                                     0.0%        28
model.go:updateSelect                               0.0%        255
model.go:updateConfirm                              0.0%        239
model.go:View                                       0.0%        23
model.go:viewSelect                                 0.0%        196
model.go:renderPreview                              0.0%        130
model.go:renderEmptyPreview                         0.0%        13
model.go:renderPreviewError                         0.0%        13
model.go:viewConfirm                                0.0%        185
model.go:copySourceToTargets                        0.0%        18
model.go:initGitWorkflow                            0.0%        53
model.go:generateDiff                               0.0%        38
model.go:renderHelpOverlay                          0.0%        61
run.go:Run                                          0.0%        61
run.go:PrintHelp                                    0.0%        75
git.go:getGitBranch                                 55.6%       18
fileops.go:copyFile                                 69.2%       65
scanner.go:scanFiles                                71.1%       84
model.go:InitialModel                               77.8%       45
model.go:matchesFilePattern                         50.0%       24
model.go:resetCursorIfNeeded                        66.7%       6
model.go:adjustViewport                             57.1%       14
model.go:filterFiles                                100.0%      13
model.go:truncate                                   100.0%      7
model.go:formatSize                                 100.0%      13
model.go:minInt                                     100.0%      7
model.go:maxInt                                     100.0%      7
branchname.go:normalizeBranchName                   100.0%      17
scanner.go:matchesPattern                           100.0%      11
```

## Testing Recommendations

### Phase 1: Critical Business Logic (Week 1)
**Goal: Get core functionality tested**

1. **Create `gitworkflow_test.go`**
   - Priority functions: `performGitWorkflow`, `processRepo`, `createWorktreeAndBranch`, `commitChanges`
   - Use table-driven tests
   - Mock `exec.Command` with test helpers
   - Target: 70% coverage for gitworkflow.go

2. **Expand `model_test.go`**
   - Add tests for `updateSelect` (most complex method)
   - Test key handling with mock messages
   - Test focus transitions
   - Test file selection logic
   - Target: 40% coverage for model.go

### Phase 2: Integration & Edge Cases (Week 2)
**Goal: Test interactions and error paths**

3. **Git Workflow Integration Tests**
   - Test complete workflow end-to-end
   - Test error scenarios (cleanup, push failures)
   - Test multi-repo operations
   - Verify cleanup happens on errors

4. **TUI State Machine Tests**
   - Test `updateConfirm` with all focus states
   - Test mode transitions (select → confirm → select)
   - Test git workflow enable/disable
   - Test tab navigation completeness

### Phase 3: View Layer & Polish (Week 3)
**Goal: Test rendering and edge cases**

5. **View Rendering Tests**
   - Test `viewSelect` and `viewConfirm` with various states
   - Test `renderPreview` in both modes
   - Test with different terminal sizes
   - Verify output formatting

6. **Edge Cases & Error Handling**
   - Test with empty file lists
   - Test with very long file paths
   - Test with non-existent files
   - Test with permission errors
   - Test concurrent operations

### Testing Tools & Techniques

#### Mocking Git Commands
```go
// Use testable command execution pattern
type Commander interface {
    CombinedOutput() ([]byte, error)
}

// In tests, provide mock implementation
type mockCommander struct {
    output []byte
    err    error
}
```

#### Mocking File System
```go
// Use afero or similar for filesystem abstraction
// Or use real filesystem with t.TempDir()
func TestGitWorkflow(t *testing.T) {
    tmpDir := t.TempDir()
    // Create test files in tmpDir
}
```

#### Testing Bubbletea Models
```go
func TestUpdateSelect(t *testing.T) {
    tests := []struct {
        name     string
        initial  model
        key      string
        expected model
    }{
        {
            name:    "up arrow moves cursor up",
            initial: modelWithCursor(5),
            key:     "up",
            expected: modelWithCursor(4),
        },
    }
    // Run tests...
}
```

#### Table-Driven Tests
```go
func TestCommitChanges(t *testing.T) {
    tests := []struct {
        name    string
        files   []string
        wantErr bool
        errMsg  string
    }{
        {"single file", []string{"test.go"}, false, ""},
        {"multiple files", []string{"a.go", "b.go"}, false, ""},
        {"nothing to commit", []string{}, false, ""},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Coverage Goals

### Short Term (1 Month)
- **Overall:** 50%+
- **gitworkflow.go:** 70%+
- **model.go:** 40%+
- **Critical functions:** 80%+

### Medium Term (3 Months)
- **Overall:** 70%+
- **gitworkflow.go:** 85%+
- **model.go:** 60%+
- **All critical functions:** 90%+

### Long Term (6 Months)
- **Overall:** 80%+
- **All files:** 70%+
- **Critical paths:** 95%+

## Benefits of Improved Coverage

### Immediate
- **Catch regressions** before they reach users
- **Refactor confidently** knowing tests will catch breaks
- **Document behavior** through test examples
- **Faster debugging** with failing tests pinpointing issues

### Long Term
- **Easier onboarding** for new contributors
- **Higher code quality** through test-driven thinking
- **Fewer production bugs** through comprehensive testing
- **Better architecture** from designing for testability

## Continuous Improvement

### CI/CD Integration
- Add coverage checks to GitHub Actions
- Fail build if coverage drops below threshold
- Generate and publish coverage reports
- Track coverage trends over time

### Developer Workflow
- Run tests before committing (pre-commit hook)
- Include coverage in PR reviews
- Require tests for new features
- Update tests when fixing bugs

### Monitoring
- Track coverage percentage over time
- Identify coverage gaps in new code
- Set coverage targets per file/package
- Review uncovered code regularly

## Next Steps

1. **Prioritize testing** the top 10 methods from this list
2. **Set up test infrastructure** for mocking git commands
3. **Write tests incrementally** (don't try to do everything at once)
4. **Run coverage regularly** to track progress
5. **Update this document** as coverage improves

---

**Generated:** 15 October 2025  
**Tool:** `go test -coverprofile=tmp/coverage.txt -covermode=atomic ./...`  
**Report:** `go tool cover -func=tmp/coverage.txt`  
**HTML:** `go tool cover -html=tmp/coverage.txt -o tmp/coverage.html`
