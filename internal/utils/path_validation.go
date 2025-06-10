package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath ensures a path component is safe for use within a project directory.
// It prevents path traversal attacks by rejecting paths that:
// - Are empty or just whitespace
// - Are absolute paths (including Windows drive paths)
// - Contain ".." segments
// - Contain invalid characters for filenames
// - Contain path separators (should be handled by structure)
func ValidatePath(name string) error {
	// Remove leading/trailing whitespace
	name = strings.TrimSpace(name)

	// Check for empty names
	if name == "" {
		return fmt.Errorf("path name cannot be empty")
	}

	// Check for absolute paths first (before separator check)
	// This catches both Unix (/path) and Windows (C:\path) absolute paths
	if filepath.IsAbs(name) {
		return fmt.Errorf("absolute paths are not allowed: %s", name)
	}

	// Additional check for Windows-style paths that might not be caught by filepath.IsAbs
	// on non-Windows systems (like drive letters)
	if len(name) >= 2 && name[1] == ':' {
		return fmt.Errorf("absolute paths are not allowed: %s", name)
	}

	// Check for path traversal attempts
	if strings.Contains(name, "..") {
		return fmt.Errorf("path traversal sequences (..) are not allowed: %s", name)
	}

	// Check for directory separators (should be handled by structure, not individual names)
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("path separators are not allowed in names: %s", name)
	}

	// Check for invalid filename characters (Windows + common problematic chars)
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("invalid character '%s' in path name: %s", char, name)
		}
	}

	// Check for reserved names (Windows reserved filenames)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	upperName := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upperName == reserved {
			return fmt.Errorf("reserved filename not allowed: %s", name)
		}
	}

	return nil
}

// SanitizePath creates a safe version of a path by cleaning it and ensuring it stays within bounds.
// This is used as a secondary safety measure after validation.
func SanitizePath(projectRoot, targetPath string) (string, error) {
	// Clean the path to resolve any . or .. components
	cleanPath := filepath.Clean(targetPath)

	// Convert both paths to absolute for comparison
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute project root: %w", err)
	}

	absTarget, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute target path: %w", err)
	}

	// Check if the target path is within the project root
	relPath, err := filepath.Rel(absRoot, absTarget)
	if err != nil {
		return "", fmt.Errorf("failed to calculate relative path: %w", err)
	}

	// If the relative path starts with "..", it's outside the project root
	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("path outside project root not allowed: %s", targetPath)
	}

	return absTarget, nil
}
