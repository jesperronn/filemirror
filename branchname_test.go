package filemirror

import (
	"runtime"
	"testing"
)

func TestNormalizeBranchName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "model.go",
			expected: "model",
		},
		{
			name:     "filename with spaces",
			input:    "my model file.go",
			expected: "my-model-file",
		},
		{
			name:     "filename with dots",
			input:    "config.prod.json",
			expected: "config-prod",
		},
		{
			name:     "filename with unicode",
			input:    "mödel.go",
			expected: "m-del",
		},
		{
			name:     "filename with mixed case",
			input:    "MyModel.go",
			expected: "mymodel",
		},
		{
			name:     "filename with underscores",
			input:    "my_model_file.go",
			expected: "my-model-file",
		},
		{
			name:     "filename with numbers",
			input:    "model123.go",
			expected: "model",
		},
		{
			name:     "filename with special chars",
			input:    "model@#$%file.go",
			expected: "model-file",
		},
		{
			name:     "filename with hyphens already",
			input:    "my-model-file.go",
			expected: "my-model-file",
		},
		{
			name:     "filename with multiple consecutive spaces",
			input:    "my   model   file.go",
			expected: "my-model-file",
		},
		{
			name:     "filename with trailing dots",
			input:    "model...go",
			expected: "model",
		},
		{
			name:     "filename starting with dot",
			input:    ".gitignore",
			expected: "gitignore",
		},
		{
			name:     "filename with only extension",
			input:    ".go",
			expected: "go",
		},
		{
			name:     "filename without extension",
			input:    "Makefile",
			expected: "makefile",
		},
		{
			name:     "filename with emoji",
			input:    "model🚀.go",
			expected: "model",
		},
		{
			name:     "filename with chinese characters",
			input:    "模型.go",
			expected: "file",
		},
		{
			name:     "filename with parentheses",
			input:    "model(1).go",
			expected: "model",
		},
		{
			name:     "filename with brackets",
			input:    "model[test].go",
			expected: "model-test",
		},
		{
			name:     "complex filename",
			input:    "My_Model File (v2.1) [FINAL].go",
			expected: "my-model-file-v-final",
		},
		{
			name:     "filename with leading/trailing hyphens after normalization",
			input:    "_model_.go",
			expected: "model",
		},
		{
			name:     "only special characters",
			input:    "@#$%.go",
			expected: "file",
		},
		{
			name:     "path with directory",
			input:    "src/models/user.go",
			expected: "user",
		},
		{
			name:     "windows path on unix",
			input:    "C:\\Users\\model.go",
			expected: "c-users-model", // filepath.Base on Unix treats backslashes as regular chars
		},
		{
			name:     "filename with ampersand",
			input:    "model&file.go",
			expected: "model-file",
		},
		{
			name:     "filename with plus",
			input:    "model+file.go",
			expected: "model-file",
		},
		{
			name:     "filename with equals",
			input:    "model=file.go",
			expected: "model-file",
		},
		{
			name:     "multiple extensions",
			input:    "archive.tar.gz",
			expected: "archive-tar",
		},
		{
			name:     "filename with colon (invalid on windows)",
			input:    "model:file.go",
			expected: "model-file",
		},
		{
			name:     "filename with quotes",
			input:    `model"file".go`,
			expected: "model-file",
		},
		{
			name:     "filename with single quotes",
			input:    "model'file'.go",
			expected: "model-file",
		},
		{
			name:     "filename with backticks",
			input:    "model`file`.go",
			expected: "model-file",
		},
		{
			name:     "filename with tilde",
			input:    "~model.go",
			expected: "model",
		},
		{
			name:     "filename with exclamation",
			input:    "model!.go",
			expected: "model",
		},
		{
			name:     "all uppercase",
			input:    "README.MD",
			expected: "readme",
		},
		{
			name:     "camelCase",
			input:    "myModelFile.go",
			expected: "mymodelfile",
		},
		{
			name:     "PascalCase",
			input:    "MyModelFile.go",
			expected: "mymodelfile",
		},
		{
			name:     "kebab-case already",
			input:    "my-model-file.go",
			expected: "my-model-file",
		},
		{
			name:     "snake_case",
			input:    "my_model_file.go",
			expected: "my-model-file",
		},
		{
			name:     "single character",
			input:    "a.go",
			expected: "a",
		},
		{
			name:     "single special character",
			input:    "$.go",
			expected: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip the "windows path on unix" test on Windows, since filepath.Base
			// handles backslashes as path separators on Windows
			if tt.name == "windows path on unix" && runtime.GOOS == "windows" {
				t.Skip("Skipping Unix-specific path test on Windows")
			}

			result := normalizeBranchName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeBranchName(%q) = %q, want %q", tt.input, result, tt.expected)
			}

			// Verify result only contains a-z and -
			for _, r := range result {
				if !((r >= 'a' && r <= 'z') || r == '-') {
					t.Errorf("normalizeBranchName(%q) = %q contains invalid character: %c", tt.input, result, r)
				}
			}

			// Verify no leading/trailing hyphens
			if result != "" {
				if result[0] == '-' || result[len(result)-1] == '-' {
					t.Errorf("normalizeBranchName(%q) = %q has leading or trailing hyphen", tt.input, result)
				}
			}

			// Verify no consecutive hyphens
			for i := 0; i < len(result)-1; i++ {
				if result[i] == '-' && result[i+1] == '-' {
					t.Errorf("normalizeBranchName(%q) = %q has consecutive hyphens", tt.input, result)
				}
			}
		})
	}
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expectErr  bool
	}{
		{"valid branch name", "feature/new-branch", false},
		{"valid branch name with numbers", "feature/new-branch-123", false},
		{"empty branch name", " ", true},
		{"branch name starting with dot", ".feature/new-branch", true},
		{"branch name ending with dot", "feature/new-branch.", true},
		{"branch name ending with .lock", "feature/new-branch.lock", true},
		{"branch name containing ..", "feature/../new-branch", true},
		{"branch name containing space", "feature/new branch", true},
		{"branch name containing ~", "feature/~new-branch", true},
		{"branch name containing ^", "feature/^new-branch", true},
		{"branch name containing :", "feature/:new-branch", true},
		{"branch name containing ?", "feature/?new-branch", true},
		{"branch name containing *", "feature/*new-branch", true},
		{"branch name containing [", "feature/[new-branch", true},
		        {"branch name containing \\", `feature\\new-branch`, true},
				{"branch name starting with /", "/feature/new-branch", true},		{"branch name ending with /", "feature/new-branch/", true},
		{"branch name containing //", "feature//new-branch", true},
		{"branch name is @", "@", true},
		{"branch name containing @{", "feature@{new-branch", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBranchName(tt.branchName)
			if (err != nil) != tt.expectErr {
				t.Errorf("validateBranchName(%q) error = %v, expectErr %v", tt.branchName, err, tt.expectErr)
			}
		})
	}
}
