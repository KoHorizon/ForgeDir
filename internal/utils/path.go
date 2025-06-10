// internal/utils/path.go
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands ~ to home directory and resolves relative paths
func ExpandPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Handle ~ expansion
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[2:])
	} else if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = homeDir
	}

	// Convert to absolute path (handles relative paths like ./my-templates)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
