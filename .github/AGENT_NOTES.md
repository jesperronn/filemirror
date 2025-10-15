# Agent Notes

## Git Workflow Policy

**IMPORTANT**: Always work in pull requests on separate branches. Never commit directly to `main`.

### Workflow Steps:
1. Create a new branch for each feature/fix: `git checkout -b feature/description`
2. Make commits on the feature branch
3. **ALWAYS push the branch to origin**: `git push -u origin feature/description`
4. Create a pull request on GitHub
5. Merge via PR after CI passes

### Critical Requirements:
- **ALWAYS create commits on a new branch** (never on `main`)
- **ALWAYS push the branch to origin immediately after committing**
- This ensures the user can create a pull request on GitHub
- Do NOT wait for the user to ask - pushing is part of the standard workflow

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
- Signing key: 88EAD15913F2AD92314C4F86B557F2DD740B2A3C
- All future commits will be automatically signed
- When using git commit commands, the `-S` flag is no longer needed as it's configured globally for this repository

## Commit Message Format

All commits should include:
1. Clear, descriptive subject line
2. Detailed body explaining the changes
3. Co-authored-by trailer:
   ```
   🤖 Generated with [Claude Code](https://claude.com/claude-code)

   Co-Authored-By: Claude <noreply@anthropic.com>
   ```

### After Committing:
**ALWAYS run `git push` immediately after committing.**
This is a mandatory step to enable the user to create a pull request on GitHub.

## Repository Information

- Main branch: `main`
- Module path: `github.com/jesperronn/filemirror`
- Binary name: `fmr`
