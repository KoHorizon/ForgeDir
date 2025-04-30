package builder

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

func CreateStructure(cfg *config.Config, pathToCreate string) error {
	realFSCreator := NewOSFileSystemCreator() // Use the constructor
	return createStructureNodes(realFSCreator, cfg.Structure, pathToCreate)
}

func createStructureNodes(fsCreator FileSystemCreator, nodes []config.StructureNode, currentPath string) error {
	for _, node := range nodes {
		filePath := filepath.Join(currentPath, node.Name)
		var err error

		switch node.Type {
		case config.TypeDir:
			err = fsCreator.CreateFolder(filePath, DefaultFolderPermission)
			if err == nil {
				err = createStructureNodes(fsCreator, node.Children, filePath)
			}
		case config.TypeFile:
			err = fsCreator.CreateFile(filePath, DefaultFilePermission)
		default:
			err = fmt.Errorf("unknown node type: %s (name: %s)", node.Type, node.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
