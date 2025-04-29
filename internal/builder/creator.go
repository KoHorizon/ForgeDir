package builder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
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

// CreateFolder creates a folder with default permission 0755 (Owner: rwx, Group/Others: r-x)
func CreateFolder(folderPath string) error {
	return CreateFolderWithPermission(folderPath, 0755)
}

// CreateFolderWithPermission creates a folder with custom permissions
func CreateFolderWithPermission(folderPath string, permission os.FileMode) error {
	err := os.Mkdir(folderPath, permission)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}
	fmt.Printf("Created folder: %s\n", folderPath)
	return nil
}

// CreateFile creates a file with default permission 0644 (Owner: read/write, Group/Others: read only)
func CreateFile(filePath string) error {
	return CreateFileWithPermission(filePath, 0644)
}

// CreateFileWithPermission creates a file with custom permission
// Flags used :
// - os.O_CREATE: Create the file if it doesn't exist
// - os.O_WRONLY: Open the file for write-only access
// - os.O_TRUNC:  Truncate the file if it already exists (clear contents)
func CreateFileWithPermission(filePath string, permission os.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permission)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Created file : %s\n", filePath)
	return nil
}

func CreateStructure(cfg *config.Config, pathToCreate string) error {
	return CreateStructureNodes(cfg.Structure, pathToCreate)
}

func CreateStructureNodes(nodes []config.StructureNode, currentPath string) error {
	for _, node := range nodes {
		filePath := filepath.Join(currentPath, node.Name)
		if node.Type == "dir" {
			err := CreateFolder(filePath)
			if err != nil {
				return err
			}
			err = CreateStructureNodes(node.Children, filePath)
			if err != nil {
				return err
			}
		} else if node.Type == "file" {
			err := CreateFile(filePath)
			if err != nil {
				return err
			}

		}
	}

	return nil
}
