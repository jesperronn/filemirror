# Development Guide

This document contains best practices and guidelines for developing FileMirror (fmr).

## Best Practices

### Code Changes

#### Keep Changes Small and Focused
- **Always make as small additions as possible**
- Each commit should address a single concern
- Break large features into smaller, reviewable chunks
- Small changes are easier to review, test, and debug

#### Consistency Updates
- **If you see other things you want to change for consistency, stop and ask**
- Don't mix refactoring with feature work
- Create separate commits for consistency improvements
- Discuss larger refactoring efforts before starting

#### Testing Requirements
- **Always add and update unit tests**
- Every new function should have corresponding tests
- Update existing tests when changing behavior
- Aim for high test coverage on critical paths
- Run `bin/test` and `bin/lint` before committing
- since tests run in CI (both unix and windows), ensure tests asserting filenames or paths work cross-platform

#### Documentation Requirements
- **Always keep documentation updated**
- Update README.md when adding user-facing features
- Update inline comments when changing function behavior
- Update IMPROVEMENTS.md when completing features or fixing bugs
- Document all public APIs with clear Go doc comments

### Go-Specific Best Practices

#### Code Style
- Follow standard Go formatting (`gofmt`)
- Use `golangci-lint` for static analysis
- Keep functions small and focused (< 50 lines when possible)
- Use meaningful variable names (avoid single letters except for short-lived variables)
- Add comments for exported functions, types, and constants

#### Error Handling
- Always check errors; never ignore them
- Wrap errors with context using `fmt.Errorf` with `%w`
- Return errors rather than panicking (except in Init functions)
- Use custom error types for specific error conditions

#### Testing
- Use table-driven tests for multiple test cases
- Test both success and error paths
- Use subtests (`t.Run`) for better organization
- Mock external dependencies (file system, git commands)
- Include edge cases and boundary conditions

#### Performance
- Profile before optimizing
- Avoid premature optimization
- Use benchmarks (`bin/test --bench`) for performance-critical code
- Consider memory allocations in hot paths

### Git Workflow

#### Commit Messages
- Use conventional commit format
- First line: brief summary (imperative mood)
- Blank line
- Detailed description of what and why
- Include ticket/issue references if applicable
- End with AI co-authorship attribution

Example:
```
Add keyboard shortcut for cycling preview modes

Implemented p/CTRL-P to cycle through hidden → plain → diff → hidden
preview states. This improves UX by consolidating preview controls
into a single, discoverable key binding.

Changes:
- Added previewHidden constant to previewMode enum
- Updated key handlers for p and ctrl+p
- Modified hints to show next preview state

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

#### Commit Frequency
- Commit working code frequently
- Each commit should compile and pass tests
- Use `git add -p` for partial staging when needed
- Squash work-in-progress commits before pushing

### Development Workflow

1. **Before Starting**
   - Pull latest changes
   - Run tests: `bin/test`
   - Review IMPROVEMENTS.md for context

2. **During Development**
   - Write tests first (TDD) or alongside code
   - Run tests frequently: `bin/test --verbose`
   - Keep changes focused on one thing
   - Update documentation as you go

3. **Before Committing**
   - Run full test suite: `bin/test`
   - Run linters: `bin/lint`
   - Review your changes: `git diff`
   - Ensure all tests pass
   - Update relevant documentation

4. **After Completing Work**
   - Update IMPROVEMENTS.md (mark items complete, add new items)
   - Ensure README reflects current features
   - Consider if DEVELOPMENT.md needs updates

### Code Review Guidelines

#### For Authors
- Keep PRs small (< 400 lines when possible)
- Write clear PR descriptions
- Link to related issues
- Respond to feedback promptly
- Update based on review comments

#### For Reviewers
- Review within 24 hours if possible
- Focus on logic, not style (use linters for style)
- Be constructive and kind
- Ask questions rather than making demands
- Approve when tests pass and code makes sense

## Project Structure

```
.
├── bin/                    # Development scripts
│   ├── build              # Build script
│   ├── lint               # Linting script
│   ├── setup              # Setup script
│   └── test               # Test script
├── .claude/               # Claude Code configuration
├── fileops.go             # File operations (copy, atomic write)
├── fileops_test.go        # File operations tests
├── git.go                 # Git integration (branch detection)
├── main.go                # Entry point and CLI parsing
├── main_test.go           # Main function tests
├── model.go               # TUI model and view logic
├── model_test.go          # TUI model tests
├── scanner.go             # File scanning and filtering
├── scanner_test.go        # Scanner tests
├── DEVELOPMENT.md         # This file
├── IMPROVEMENTS.md        # Feature and bug tracking
└── README.md              # User documentation
```

## Testing Strategy

### Unit Tests
- Test individual functions in isolation
- Mock external dependencies (filesystem, git)
- Cover happy path and error cases
- Use table-driven tests for multiple scenarios

### Integration Tests
- Test interaction between components
- Use temporary directories for file operations
- Clean up resources in teardown

### TUI Tests
- Test model state transitions
- Test key handling logic
- Test view rendering (where practical)
- Mock bubbletea messages

## Common Development Tasks

### Adding a New Feature

1. Review IMPROVEMENTS.md to understand context
2. Write tests for the new feature first
3. Implement the feature with minimal changes
4. Update tests to ensure they pass
5. Update documentation (README, inline comments)
6. Run full test suite and linting
7. Update IMPROVEMENTS.md to mark feature complete
8. Commit with descriptive message

### Fixing a Bug

1. Add a failing test that reproduces the bug
2. Fix the bug with minimal changes
3. Ensure the test now passes
4. Check for similar issues elsewhere
5. Update documentation if behavior changed
6. Update IMPROVEMENTS.md to mark bug fixed
7. Commit with clear description of fix

### Refactoring Code

1. **Stop and discuss** if changes affect multiple areas
2. Ensure tests exist and pass before refactoring
3. Make refactoring changes without adding features
4. Ensure tests still pass after refactoring
5. Commit refactoring separately from feature work
6. Update documentation if interfaces changed

## Tools and Commands

### Building
```bash
# Development build
go build -o fmr .

# Using build script (recommended)
bin/build

# Cross-platform build
bin/build -o linux -a amd64
```

### Testing
```bash
# Run all tests
go test ./...

# Using test script (recommended)
bin/test

# Verbose output
bin/test --verbose

# With coverage
bin/test --coverage

# With race detector
bin/test --race
```

### Linting
```bash
# Run all linters
bin/lint

# Auto-fix issues
bin/lint --fix

# Manual linting
gofmt -w .
go vet ./...
golangci-lint run
```

### Running
```bash
# Run without building
go run . "*.go"

# Run built binary
./fmr --path ~/projects "config.yaml"
```

## AI Development Notes

### For Claude Code and Other AI Agents

When working on this project:

1. **Always read these guidelines** before making changes
2. **Keep changes minimal** - resist the urge to "improve" unrelated code
3. **Ask before refactoring** - consistency changes should be discussed
4. **Tests are mandatory** - no code without tests
5. **Documentation is mandatory** - update as you go, not later
6. **One thing per commit** - don't mix concerns
7. **Follow Go conventions** - idiomatic Go is preferred
8. **Review your own changes** - read the diff before committing
9. **Update IMPROVEMENTS.md** - track completed work
10. **Commit message quality matters** - future developers will read them

### Common Pitfalls to Avoid

- Don't refactor and add features in the same commit
- Don't skip tests because "it's a small change"
- Don't ignore linter warnings
- Don't mix formatting changes with logic changes
- Don't assume documentation is up to date - verify and update
- Don't make changes without understanding the context
- Don't optimize without profiling first

## Getting Help

- Check IMPROVEMENTS.md for known issues and planned features
- Review existing tests for patterns and examples
- Read Go documentation: https://golang.org/doc/
- Review bubbletea examples: https://github.com/charmbracelet/bubbletea/tree/master/examples
- Ask questions in discussions rather than making assumptions

## License

MIT - See LICENSE file for details
