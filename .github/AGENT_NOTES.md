# Agent Notes

## Git Workflow Policy

**IMPORTANT**: Always work in pull requests on separate branches. Never commit directly to `main`.

### Workflow Steps:
1. Create a new branch for each feature/fix: `git checkout -b feature/description`
2. Make commits on the feature branch
3. Push the branch to origin: `git push -u origin feature/description`
4. Create a pull request on GitHub
5. Merge via PR after CI passes

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

## Repository Information

- Main branch: `main`
- Module path: `github.com/jesperronn/filemirror`
- Binary name: `fmr`
