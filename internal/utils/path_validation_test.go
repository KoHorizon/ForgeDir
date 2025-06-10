package utils

import (
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		// Valid cases
		{
			name:        "simple filename",
			input:       "main.go",
			expectError: false,
		},
		{
			name:        "simple directory name",
			input:       "src",
			expectError: false,
		},
		{
			name:        "filename with underscore",
			input:       "test_file.txt",
			expectError: false,
		},
		{
			name:        "filename with dash",
			input:       "my-file.json",
			expectError: false,
		},
		{
			name:        "filename with dots",
			input:       "config.dev.yaml",
			expectError: false,
		},

		// Invalid cases - path traversal
		{
			name:        "parent directory traversal",
			input:       "../",
			expectError: true,
			errorMsg:    "path traversal sequences (..) are not allowed",
		},
		{
			name:        "parent directory in filename",
			input:       "../evil.txt",
			expectError: true,
			errorMsg:    "path traversal sequences (..) are not allowed",
		},
		{
			name:        "nested parent directory",
			input:       "../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal sequences (..) are not allowed",
		},
		{
			name:        "hidden parent directory",
			input:       "normal..evil",
			expectError: true,
			errorMsg:    "path traversal sequences (..) are not allowed",
		},

		// Invalid cases - absolute paths
		{
			name:        "absolute unix path",
			input:       "/etc/passwd",
			expectError: true,
			errorMsg:    "absolute paths are not allowed",
		},
		{
			name:        "absolute windows path",
			input:       "C:\\Windows\\System32",
			expectError: true,
			errorMsg:    "absolute paths are not allowed",
		},

		// Invalid cases - path separators
		{
			name:        "forward slash in name",
			input:       "dir/file.txt",
			expectError: true,
			errorMsg:    "path separators are not allowed in names",
		},
		{
			name:        "backslash in name",
			input:       "dir\\file.txt",
			expectError: true,
			errorMsg:    "path separators are not allowed in names",
		},

		// Invalid cases - empty/whitespace
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorMsg:    "path name cannot be empty",
		},
		{
			name:        "only whitespace",
			input:       "   ",
			expectError: true,
			errorMsg:    "path name cannot be empty",
		},

		// Invalid cases - reserved characters
		{
			name:        "less than symbol",
			input:       "file<name.txt",
			expectError: true,
			errorMsg:    "invalid character '<' in path name",
		},
		{
			name:        "pipe symbol",
			input:       "file|name.txt",
			expectError: true,
			errorMsg:    "invalid character '|' in path name",
		},
		{
			name:        "question mark",
			input:       "file?.txt",
			expectError: true,
			errorMsg:    "invalid character '?' in path name",
		},

		// Invalid cases - reserved names (Windows)
		{
			name:        "CON reserved name",
			input:       "CON",
			expectError: true,
			errorMsg:    "reserved filename not allowed",
		},
		{
			name:        "con lowercase",
			input:       "con",
			expectError: true,
			errorMsg:    "reserved filename not allowed",
		},
		{
			name:        "COM1 reserved",
			input:       "COM1",
			expectError: true,
			errorMsg:    "reserved filename not allowed",
		},
		{
			name:        "AUX reserved",
			input:       "aux",
			expectError: true,
			errorMsg:    "reserved filename not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error for input %q, but got none", tt.input)
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for input %q, but got: %v", tt.input, err)
				}
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		projectRoot string
		targetPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid path within project",
			projectRoot: tempDir,
			targetPath:  filepath.Join(tempDir, "src", "main.go"),
			expectError: false,
		},
		{
			name:        "path exactly at project root",
			projectRoot: tempDir,
			targetPath:  tempDir,
			expectError: false,
		},
		{
			name:        "attempt to escape project root",
			projectRoot: tempDir,
			targetPath:  filepath.Join(tempDir, "..", "evil.txt"),
			expectError: true,
			errorMsg:    "path outside project root not allowed",
		},
		{
			name:        "complex escape attempt",
			projectRoot: tempDir,
			targetPath:  filepath.Join(tempDir, "src", "..", "..", "evil.txt"),
			expectError: true,
			errorMsg:    "path outside project root not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizePath(tt.projectRoot, tt.targetPath)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error for targetPath %q, but got none", tt.targetPath)
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for targetPath %q, but got: %v", tt.targetPath, err)
				}
				if result == "" {
					t.Error("expected non-empty result for valid path")
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 || s[len(s)-len(substr):] == substr ||
			s[:len(substr)] == substr ||
			indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
