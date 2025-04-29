package builder

import (
	"fmt"
	"os"
)

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
//   0 = --- (no permission)
//   1 = --x (execute)
//   2 = -w- (write)
//   3 = -wx (write and execute)
//   4 = r-- (read)
//   5 = r-x (read and execute)
//   6 = rw- (read and write)
//   7 = rwx (read, write, and execute)
//
// Format: 0<Owner><Group><Others> (each digit is a sum of r=4, w=2, x=1)
// Note: Permissions are based on standard UNIX file mode bits, stable across Go versions.

// createFolder creates a folder with default permission 0755 (Owner: rwx, Group/Others: r-x)
func createFolder(folderPath string) error {
	return createFolderWithPermission(folderPath, 0755)
}

// createFolderWithPermission creates a folder with custom permissions
func createFolderWithPermission(folderPath string, permission os.FileMode) error {
	err := os.MkdirAll(folderPath, permission)
	if err != nil {
		return err
	}
	fmt.Printf("Created folder: %s\n", folderPath)
	return nil
}

// createFile creates a file with default permission 0644 (Owner: read/write, Group/Others: read only)
func createFile(filePath string) error {
	return createFileWithPermission(filePath, 0644)
}

// createFileWithPermission creates a file with custom permission
// Flags used :
// - os.O_CREATE: Create the file if it doesn't exist
// - os.O_WRONLY: Open the file for write-only access
// - os.O_TRUNC:  Truncate the file if it already exists (clear contents)
func createFileWithPermission(filePath string, permission os.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permission)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Created file : %s\n", filePath)
	return nil
}
