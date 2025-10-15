# Agent Notes

## Git Workflow Policy

**CRITICAL RULE**: Always work on a feature branch. NEVER commit directly to `main`.

### Workflow Steps:
1. **FIRST**: Check current branch with `git branch --show-current`
2. **BEFORE ANY WORK**: Create a new branch: `git checkout -b feature/description`
3. Make commits on the feature branch
4. **ALWAYS push the branch to origin**: `git push -u origin feature/description`
5. Create a pull request on GitHub
6. Merge via PR after CI passes

### Critical Requirements:
- ⚠️ **NEVER EVER commit directly to `main`** - this is the most important rule
- **ALWAYS create a new branch BEFORE starting any work**
- **Push the branch to origin immediately after committing**
- This ensures the user can create a pull request on GitHub
- Do NOT wait for the user to ask - pushing is part of the standard workflow

### Exceptions:
- Only exception: If you're already on a feature branch and making amendments to the previous commit

### Branch Management Rules:
- **ONLY checkout a new branch for NEW features/bugs/tasks**
- **STAY on the current branch when fixing/improving existing work**
- If already working on a branch (e.g., `feature/ci`), continue on that branch
- Do NOT create new branches for related fixes - add commits to the existing branch
- Check current branch with `git branch --show-current` before creating new branches

### Branch Naming Convention:
- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation changes
- `refactor/*` - Code refactoring
- `ci/*` - CI/CD changes

## Git Commit Policy

**IMPORTANT**: All commits must be GPG signed.

- Git signing is enabled via `git config commit.gpgsign true`
- All future commits will be automatically signed
- When using git commit commands, the `-S` flag is no longer needed as it's configured globally for this repository

## Commit Message Format

All commits should include:
1. Clear, descriptive subject line
2. Detailed body explaining the changes
3. Co-authored-by trailer:
   ```
   🤖 Generated with <name of the agent>

   Co-Authored-By: <co-author name> <co-author email>
   ```

   where

   - `<name of the agent>` is the name of the AI agent used (e.g., Claude Code)
   - `<co-author name>`, `<co-author email>` is the email of the AI agent
     Example: Claude <noreply@anthropic.com>

### After Committing:
**ALWAYS run `git push` immediately after committing.**
This is a mandatory step to enable the user to create a pull request on GitHub.

## Repository Information

- Main branch: `main`
- Module path: `github.com/jesperronn/filemirror-fmr`
- Binary names: `fmr` and `filemirror` (identical binaries, fmr is shorter alias)
