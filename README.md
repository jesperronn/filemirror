# FileMirror (fmr)

[![CI](https://github.com/jesperronn/filemirror-fmr/workflows/CI/badge.svg)](https://github.com/jesperronn/filemirror-fmr/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jesperronn/filemirror-fmr)](https://goreportcard.com/report/github.com/jesperronn/filemirror-fmr)
[![codecov](https://codecov.io/gh/jesperronn/filemirror-fmr/branch/main/graph/badge.svg)](https://codecov.io/gh/jesperronn/filemirror-fmr)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A repo-aware CLI for syncing canonical files across multiple directories and repositories with automatic git workflow support.

**Interactive TUI • Diff preview • Per-repo git commits • Atomic writes**

## Quick Start

```bash
# Install
go install github.com/jesperronn/filemirror-fmr/cmd/fmr@latest

# Run (searches for files in current directory)
fmr "*.yaml"

# Or search in specific directory
fmr --path ~/projects "config"
```

**In the TUI:**
- Press `s` to mark source file
- Press `SPACE` to mark target files
- Press `p` to preview diffs
- Press `ENTER` to confirm → configure git workflow → sync & commit

## Interface

```
┌─────────────────────────────────────────────────────────────────────────┐
│ FileMirror - File Synchronization Tool                                 │
│                                                                         │
│ ↑/↓: navigate • s: source • SPACE: target • p: preview • ENTER: sync   │
│                                                                         │
│ PATH: /Users/you/projects                                              │
│ SEARCH: *.yaml                                                         │
│                                                                         │
│ Source: config.yaml                                                    │
│                                                                         │
│ ╭────────────────────╮│ Preview (diff): config.yaml → api/config ────┐│
│ │ FILE LIST         ││                                              ││
│ │ ▶[S] config.yaml  ││ @@ Diff Preview @@                          ││
│ │  [T] api/config   ││ - old_value: false                          ││
│ │  [T] web/config   ││ + new_value: true                           ││
│ │  [ ] test.yaml    ││                                              ││
│ │                   ││ Press CTRL-U/D to scroll                    ││
│ ╰────────────────────╯│                                              ││
│                       └──────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────┘
```

## Features

- **Interactive TUI** - Split-screen with live diff preview
- **Git Workflow** - Automatic branch creation, commit, and optional push per repository
- **Multi-Repo Support** - Sync files across different git repositories
- **Safe Operations** - Atomic writes, worktree isolation, permission preservation
- **Smart Navigation** - TAB between path/search/files, vim keys (hjkl), inline editing

## Installation

### Go Install (Recommended)
```bash
go install github.com/jesperronn/filemirror-fmr/cmd/fmr@latest
```

### From Source
```bash
git clone https://github.com/jesperronn/filemirror-fmr
cd filemirror-fmr
go build -o fmr ./cmd/fmr
sudo mv fmr /usr/local/bin/
```

## Usage

```bash
fmr [OPTIONS] [PATTERN]
```

**Options:**
- `-p, --path PATH` - Start in directory PATH
- `-h, --help` - Show help
- `-v, --version` - Show version

**Pattern Examples:**
- `*.go` - All Go files
- `config` - Files containing "config"
- _(empty)_ - All files (4 levels deep)

## Keyboard Shortcuts

### File Selection
| Key | Action |
|-----|--------|
| `↑`/`↓` or `k`/`j` | Navigate files |
| `s` | Mark as source |
| `SPACE` | Toggle target |
| `ENTER` | Proceed to confirmation |

### View & Navigation
| Key | Action |
|-----|--------|
| `TAB` / `Shift+TAB` | Cycle focus: Path → Search → Files |
| `p` / `CTRL-P` | Toggle preview: hidden → plain → diff |
| `CTRL-U` / `CTRL-D` | Scroll preview |
| `CTRL-R` | Reload files |
| `?` | Help overlay |

### Git Workflow (Confirmation Screen)
| Key | Action |
|-----|--------|
| `TAB` | Navigate fields (branch, commit msg, push) |
| `CTRL-G` | Toggle git on/off |
| `SPACE` | Toggle checkboxes |
| `ENTER` | Execute copy & commit |
| `ESC` | Cancel |

## Git Workflow

After selecting files, the confirmation screen provides git integration:

```
╭──────────────────╮│ Git Workflow Configuration ──────────────┐
│ FILES TO SYNC   ││                                          │
│ Source: cfg.yaml││ [✓] Create git commit                    │
│ Targets: (2)    ││                                          │
│                 ││ Branch: chore/filesync-cfg               │
│                 ││ ╭──────────────────────────────────────╮ │
│                 ││ │ chore/filesync-cfg                   │ │
│                 ││ ╰──────────────────────────────────────╯ │
│                 ││                                          │
│                 ││ Commit Message:                          │
│                 ││ ╭──────────────────────────────────────╮ │
│                 ││ │ chore: Sync cfg.yaml                 │ │
│                 ││ │ Synchronized from cfg.yaml           │ │
│                 ││ ╰──────────────────────────────────────╯ │
│                 ││                                          │
│                 ││ [ ] Push to origin after commit         │
│                 ││                                          │
│                 ││ Repository: 2 git repos detected        │
╰──────────────────╯│    ┌────────────┐  ┌────────┐         │
                    │    │Copy&Commit │  │ Cancel │         │
                    └────┴────────────┴──┴────────┴─────────┘
```

**How it works:**
1. Detects git repos for target files
2. Creates isolated worktree per repository
3. Creates branch (or safely reuses existing)
4. Commits synced file with custom message
5. Optionally pushes to origin
6. Cleans up worktrees automatically

**Branch reuse:** If branch exists with only the same file modified, it's reused. Otherwise, an error is shown.

**Default settings:**
- Git workflow: ON (if repos detected)
- Branch naming: `chore/filesync-{filename}`
- Push to origin: OFF (manual push safer)

**Toggle off:** Press `CTRL-G` or uncheck `Create git commit` to copy files only (no git operations).

## Workflow Example

```bash
# 1. Start FileMirror
fmr --path ~/projects "*.yaml"

# 2. In TUI:
#    - Press 's' on canonical config.yaml
#    - Press SPACE on target files in different repos
#    - Press 'p' to preview diffs
#    - Press ENTER

# 3. In Git Workflow screen:
#    - Edit branch name if needed
#    - Edit commit message
#    - Check "Push to origin" if desired
#    - Press ENTER on "Copy & Commit"

# 4. FileMirror executes:
#    ✓ Copies files
#    ✓ Creates worktrees per repo
#    ✓ Commits on new branches
#    ✓ Pushes if enabled
#    ✓ Shows summary

# Terminal output after completion:
Git Workflow Complete!

✓ Successfully committed to 2 repository(ies):
  - /Users/you/projects/project-api
  - /Users/you/projects/project-web

Branch: chore/sync-config
Push: NO (you can push manually later)

Next steps:
  - Review commits: git log -1 (in each repository)
  - Push to remote: git push -u origin chore/sync-config
  - Create pull requests
```

## Use Cases

- Sync configuration files across microservices with automatic commits
- Update component files in different feature branches
- Propagate utility functions across packages with git audit trail
- Mirror canonical files across repositories with branching & PR workflow
- Standardize linting configs with tracked changes

## How It Works

1. **Scanning** - Recursively scans up to 4 directory levels
2. **Filtering** - Excludes: node_modules, .git, vendor, dist, build, target, .cache
3. **Pattern Matching** - Wildcards (`*.ext`) or substring matching
4. **Git Detection** - Detects repo root for each file
5. **Safe Copying** - Atomic writes (temp file → rename)
6. **Git Worktrees** - Isolated branch work without disrupting working directory
7. **Cleanup** - Automatic worktree removal (even on errors)

## Development

### Setup
```bash
bin/setup          # Install dependencies
bin/setup --full   # Include dev tools (golangci-lint)
```

### Build
```bash
bin/build                    # Build binary
bin/build --install          # Install to $GOPATH/bin
bin/build -o linux -a amd64  # Cross-compile
```

### Test
```bash
bin/test                # Run tests
bin/test --coverage     # With coverage report
bin/test --race         # Race detector
```

### Lint
```bash
bin/lint           # Run linters
bin/lint --fix     # Auto-fix issues
```

## Requirements

- Go 1.21+
- Git (for branch detection and git workflow)

## Contributing

Contributions welcome! Please submit issues or pull requests.

## License

MIT
