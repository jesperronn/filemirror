# FileMirror (fmr) - Improvements & Roadmap

This document tracks planned improvements and feature ideas for FileMirror.

## High Priority Features

### 1. Resizable Preview Panel
- **Description**: Allow dynamic resizing of preview panel width
- **Implementation**: Press `CTRL-Left`/`CTRL-Right` to adjust width by 10% increments
- **Rationale**: Users may want more/less space for diff preview depending on content

### 2. Pattern Matching Feedback
- **Description**: Better visual feedback for pattern matching
- **Implementation**:
  - Show match count: "Found 15 files matching '*.go'"
  - Indicate mode (glob vs substring) in UI
  - Support multiple patterns: `*.go,*.md` or `*.{go,md}`
- **Rationale**: Clearer understanding of what files are being shown

### 3. Search History
- **Description**: Remember recent search patterns
- **Implementation**:
  - Use ↑/↓ in search field to cycle through history
  - Store in `~/.fmr_history`
  - Limit to last 50 patterns
- **Rationale**: Faster repeat searches

### 4. Path History and Bookmarks
- **Description**: Quick access to frequently used directories
- **Implementation**:
  - Remember recently used paths
  - Allow bookmarking with keyboard shortcut
  - Quick access with numbered shortcuts or fuzzy search
- **Rationale**: Navigate to common project directories faster

## Medium Priority Features

### 5. Configuration File Support
- **Description**: User-configurable settings
- **Implementation**: Support `.fmrrc` or `~/.config/fmr/config.yaml`
- **Configurable settings**:
  - Excluded directories (beyond defaults)
  - Default search depth (currently 4 levels)
  - Color schemes
  - Keyboard shortcuts
  - Default working directory
- **Rationale**: Personalize tool to workflow

### 6. Enhanced Diff Preview
- **Description**: Richer diff visualization before sync
- **Implementation**:
  - Side-by-side diff option (vs unified)
  - Syntax highlighting in diffs
  - Option to skip specific targets after reviewing diff
- **Rationale**: Better decision making before sync

### 7. Batch Operations
- **Description**: More flexible sync modes
- **Implementation**:
  - Multiple source files (merge/combine)
  - Different sync modes: replace (current), append, prepend
  - Section-based merging for configs
- **Rationale**: Support more complex sync scenarios

## Low Priority / Future

### 8. Performance Improvements
- **Lazy loading**: Virtualize file list for thousands of files
- **Parallel operations**: Copy to multiple targets in parallel
- **Cached scans**: Cache directory scans with invalidation

### 9. Advanced Git Features
- **Commit templates**: Save/load common commit message patterns
- **Conflict detection**: Warn if target files have uncommitted changes
- **Auto-PR creation**: Integrate with `gh` CLI to create pull requests

### 10. Remote Support
- **Description**: Sync files over SSH
- **Implementation**: Support paths like `user@host:/path/to/dir`
- **Rationale**: Multi-host configuration sync

## Recently Completed ✓

- [x] Pattern matching for glob patterns (2025-10-14)
- [x] File preview panel with diff mode (2025-10-14)
- [x] Keyboard hints visibility (2025-10-15)
- [x] CTRL+P preview toggle (2025-10-15)
- [x] Multi-state preview (hidden/plain/diff) (2025-10-15)
- [x] Error display in confirmation UI (2025-10-16)
- [x] Git workflow integration (2025-10-16)
- [x] Branch reuse validation logic (2025-10-16)
- [x] Enhanced exit summary (2025-10-16)

## Contributing

If you'd like to work on any of these items:
1. Check for open issues on GitHub
2. Create an issue if one doesn't exist
3. Discuss approach before implementing
4. Submit PR with tests and documentation

## Notes

- Features should be added incrementally
- Maintain backwards compatibility
- Keep the tool fast and lightweight
- Prioritize user experience
