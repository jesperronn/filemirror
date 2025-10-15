# Agent Notes

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
