package filemirror

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// normalizeBranchName converts a filename into a git-safe branch name component.
// Only lowercase letters [a-z] and hyphens [-] are allowed.
// All other characters are converted to hyphens, and multiple consecutive hyphens
// are collapsed to a single hyphen.
func normalizeBranchName(filename string) string {
	// Remove file extension
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	// If empty after removing extension, use the base name
	if nameWithoutExt == "" {
		nameWithoutExt = base
	}

	// Convert to lowercase
	normalized := strings.ToLower(nameWithoutExt)

	// Replace any character that is not a-z or 0-9 with a hyphen
	// We'll allow digits initially, then remove them in a second pass if needed
	var result strings.Builder
	for _, r := range normalized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			// Unicode letters/digits -> convert to hyphen
			result.WriteRune('-')
		} else {
			// Any other character (space, dot, etc.) -> hyphen
			result.WriteRune('-')
		}
	}

	// Remove digits (only a-z and - allowed per requirement)
	digitRemoved := regexp.MustCompile(`[0-9]`).ReplaceAllString(result.String(), "-")

	// Collapse multiple consecutive hyphens to single hyphen
	collapsed := regexp.MustCompile(`-+`).ReplaceAllString(digitRemoved, "-")

	// Trim leading and trailing hyphens
	trimmed := strings.Trim(collapsed, "-")

	// If empty after all processing, use a fallback
	if trimmed == "" {
		return "file"
	}

	return trimmed
}
