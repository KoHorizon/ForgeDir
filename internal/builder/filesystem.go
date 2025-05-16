package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileSystem interface {
	CreateFolder(path string, permission os.FileMode) error
	WriteFile(path string, content []byte, permission os.FileMode) error
}

type OSFileSystem struct{}

func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// Common Folder Permissions:
// 0755 → Owner: read/write/execute | Group: read/execute | Others: read/execute
// 0700 → Owner: read/write/execute | Group: no access   | Others: no access
// 0777 → Everyone: read/write/execute (!! usually not recommended for security !!)
//
// Common File Permissions (for reference):
// 0644 → Owner: read/write | Group: read | Others: read
// 0600 → Owner: read/write | Group: no access | Others: no access
// 0755 → Owner: read/write/execute | Group: read/execute | Others: read/execute (for executable files)
//
// Octal Permission Breakdown:
//
//	0 = --- (no permission)
//	1 = --x (execute)
//	2 = -w- (write)
//	3 = -wx (write and execute)
//	4 = r-- (read)
//	5 = r-x (read and execute)
//	6 = rw- (read and write)
//	7 = rwx (read, write, and execute)
//
// Format: 0<Owner><Group><Others> (each digit is a sum of r=4, w=2, x=1)
// Note: Permissions are based on standard UNIX file mode bits, stable across Go versions.
//
// CreateFolder creates a folder
func (o *OSFileSystem) CreateFolder(folderPath string, permission os.FileMode) error {
	err := os.MkdirAll(folderPath, permission)
	if err != nil {
		return err
	}
	fmt.Printf("Created folder: %s\n", folderPath)
	return nil
}

// WriteFile ensures the parent directory exists, then creates or truncates
// the file at path and writes the provided content with the specified permission.
func (o *OSFileSystem) WriteFile(path string, content []byte, permission os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultFolderPermission); err != nil {
		return fmt.Errorf("mkdir parent %s: %w", dir, err)
	}
	// Write file content
	if err := os.WriteFile(path, content, permission); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	fmt.Printf("Wrote file: %s\n", path)
	return nil
}

// Define default permissions as constant for clarity
// Default permission 0755 (Owner: rwx, Group/Others: r-x)
// Fefault permission 0644 (Owner: read/write, Group/Others: read only)
const (
	DefaultFolderPermission os.FileMode = 0755
	DefaultFilePermission   os.FileMode = 0644
)
