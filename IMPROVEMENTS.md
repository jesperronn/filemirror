# Multiedit - Improvements & Bug Fixes

This document tracks known bugs, planned improvements, and feature ideas for the multiedit project.

## Known Bugs

### Fixed ✓

1. **~~File list doesn't show filtered results by pattern~~** ✓ FIXED
   - **Description**: When entering a pattern like `*.java` in the search input, the file list doesn't update to show only matching files
   - **Status**: Fixed - Added `matchesFilePattern()` function that properly handles glob patterns
   - **Implementation**: Now supports both glob patterns (`*.go`) and substring matching
   - **Location**: `model.go` (matchesFilePattern function)
   - **Date Fixed**: 2025-10-14

2. **~~Missing file preview panel~~** ✓ FIXED
   - **Description**: No right-side preview panel showing file contents
   - **Status**: Fixed - Added split-screen preview panel
   - **Features**:
     - Toggle preview with `p` key
     - Scroll with PgUp/PgDn or Ctrl+U/Ctrl+D
     - Shows file contents in real-time
     - Displays line count and scroll position
     - Adjusts layout to 50/50 split when enabled
   - **Location**: `model.go` (renderPreview functions)
   - **Date Fixed**: 2025-10-14

### High Priority

### Medium Priority

3. **Pattern matching inconsistency**
   - **Description**: Pattern matching uses both glob (filepath.Match) and substring matching, can be confusing
   - **Current**: If pattern contains `*`, use glob; otherwise substring
   - **Improvement**: Clearer indication of which mode is active, or support both simultaneously
   - **Location**: `scanner.go:81-94` (matchesPattern function)

4. **No visual diff before sync**
   - **Description**: Users cannot see what will change before confirming sync
   - **Expected**: Show diff between source and target files
   - **Current**: Only shows file paths in confirmation screen
   - **Location**: `model.go:310-336` (viewConfirm function)

## Feature Improvements

### UI/UX Enhancements

1. **Add file preview panel (right side)**
   - Show file contents with syntax highlighting
   - Display when a file is selected/focused
   - Could use existing libraries:
     - `github.com/alecthomas/chroma` for syntax highlighting
     - `github.com/charmbracelet/glamour` for markdown rendering
   - Update layout to split screen: file list (left) | preview (right)
   - Add toggle key (e.g., `p` or `F2`) to show/hide preview

2. **Improve pattern matching feedback**
   - Show pattern matching mode in UI (glob vs substring)
   - Display match count: "Found 15 files matching '*.go'"
   - Highlight matching parts of filenames
   - Support multiple patterns: `*.go,*.md` or `*.{go,md}`

3. **Better visual indicators**
   - Color-code files by type/extension
   - Show file icons (if terminal supports)
   - Add progress bar for file operations
   - Highlight modified vs unmodified files

4. **Enhanced confirmation screen**
   - Show file size comparison (source vs targets)
   - Display modification timestamps
   - Preview first few lines of changes
   - Add option to review each target individually

### Functionality Enhancements

5. **Diff preview before sync**
   - Show unified diff between source and each target
   - Use colors to highlight additions/deletions
   - Allow reviewing diffs before confirming
   - Option to skip specific targets after seeing diff

6. **Search history**
   - Remember recent search patterns
   - Use ↑/↓ in search field to cycle through history
   - Store in `~/.multiedit_history` or similar

7. **Path history and bookmarks**
   - Remember recently used paths
   - Allow bookmarking frequently used directories
   - Quick access with numbered shortcuts

8. **Batch operations**
   - Select multiple source files
   - Merge content from multiple sources
   - Support different sync modes:
     - Replace (current behavior)
     - Append
     - Prepend
     - Merge sections

9. **Undo functionality**
   - Keep backup of original files before sync
   - Allow undo of last sync operation
   - Store backups in `.multiedit/backups/`

10. **Advanced filtering**
    - Filter by file size
    - Filter by modification time
    - Filter by git status (modified, untracked, etc.)
    - Combine multiple filters

### Performance Improvements

11. **Lazy loading for large file lists**
    - Don't load all files at once if there are thousands
    - Paginate or virtualize the file list
    - Show loading indicator during scan

12. **Cached file scans**
    - Cache scan results for recently visited directories
    - Invalidate cache on file system changes
    - Make scans faster when returning to same directory

13. **Parallel file operations**
    - Copy to multiple targets in parallel
    - Show progress for each operation
    - Handle errors gracefully per target

### Developer Experience

14. **Configuration file support**
    - Support `.multieditrc` or `.multiedit.json`
    - Configure:
      - Excluded directories
      - Default search depth
      - Color schemes
      - Keyboard shortcuts
      - Default working directory

15. **Plugin system**
    - Allow custom file processors
    - Pre-sync and post-sync hooks
    - Custom syntax highlighters

16. **Logging and debugging**
    - Add `--debug` flag for verbose logging
    - Log file operations to help troubleshoot
    - Show scan statistics (files scanned, time taken, etc.)

## Nice-to-Have Features

17. **Template support**
    - Define templates for common file patterns
    - Variable substitution when copying
    - Save and reuse sync configurations

18. **Remote file support**
    - Sync files over SSH
    - Support remote paths: `user@host:/path/to/dir`
    - Integrate with SCP/SFTP

19. **Git integration enhancements**
    - Show git diff for files
    - Respect .gitignore
    - Offer to commit after sync
    - Stage synced files automatically

20. **Dry-run mode**
    - Show what would be synced without doing it
    - Useful for testing patterns and selections
    - Add `--dry-run` flag

21. **Watch mode**
    - Continuously watch source file for changes
    - Auto-sync to targets when source changes
    - Useful for development workflows

## Testing Improvements

22. **Integration tests**
    - Test full user workflows
    - Test TUI interactions
    - Mock file system operations

23. **Benchmark tests**
    - Test performance with large file sets
    - Measure scan time improvements
    - Profile memory usage

24. **Platform-specific tests**
    - Test on Linux, macOS, Windows
    - Test path handling differences
    - Test terminal compatibility

## Documentation Improvements

25. **Video/GIF demos**
    - Create animated demo showing workflow
    - Add to README
    - Show different use cases

26. **Man page**
    - Create proper man page
    - Install with `man multiedit`
    - Document all flags and shortcuts

27. **Cookbook/Recipes**
    - Common use cases with examples
    - Tips and tricks
    - Integration with other tools

28. **Architecture documentation**
    - Document code structure
    - Explain key design decisions
    - Contributing guidelines

## Priority Ranking

### Recently Completed ✓
- [x] Bug #1: Pattern matching for glob patterns (2025-10-14)
- [x] Bug #2: File preview panel (2025-10-14)

### Must Fix (P0)
- None currently

### Should Have (P1)
- [ ] Feature #5: Diff preview before sync
- [x] ~~Feature #1: File preview panel implementation~~ (COMPLETED 2025-10-14)
- [ ] Feature #9: Undo functionality

### Nice to Have (P2)
- [ ] Feature #6: Search history
- [ ] Feature #7: Path history
- [ ] Feature #14: Configuration file
- [ ] Feature #10: Advanced filtering

### Future (P3)
- [ ] All other features listed above

## Notes

- Features should be added incrementally
- Maintain backwards compatibility
- Keep the tool fast and lightweight
- Prioritize user experience
- Document all new features

## Contributing

If you'd like to work on any of these items:
1. Check if there's an open issue
2. Create an issue if one doesn't exist
3. Discuss the approach before implementing
4. Submit a PR with tests and documentation
