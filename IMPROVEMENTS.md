# Multiedit - Improvements & Bug Fixes

This document tracks known bugs, planned improvements, and feature ideas for the multiedit project.

## Known Bugs

### High Priority

None currently

### Medium Priority

4. possibility to resize previe panel for instance pressing ctrl+left/right to adjust width of preview panel. Adjust with 10% increments/decrements.


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

1. **Git integrations**
    files that are git tracked and modified should have the option to  be updated in a git commit after sync.

    - after sync, detect if any files were git tracked and modified.
    -  if so, prompt user to commit changes.
    -  if user agrees, ensure the synced files are staged in a separate commit in each repo.

    So you need to ask for commit message, and branch name. When given commit message and branch name. create separate worktrees for each repo, stage the files, and commit, and push if possible.

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
- [x] Bug #3: Keyboard hints visibility (2025-10-15)
- [x] Bug #4: CTRL+P support for preview toggle (2025-10-15)
- [x] Bug #5: Multi-state preview toggle - hidden/plain/diff cycling (2025-10-15)

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
