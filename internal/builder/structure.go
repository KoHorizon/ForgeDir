package builder

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

func CreateStructure(cfg *config.Config, pathToCreate string) error {
	return createStructureNodes(cfg.Structure, pathToCreate)
}

func createStructureNodes(nodes []config.StructureNode, currentPath string) error {
	for _, node := range nodes {
		filePath := filepath.Join(currentPath, node.Name)
		var err error

		switch node.Type {
		case config.TypeDir:
			err = createFolder(filePath)
			if err == nil {
				err = createStructureNodes(node.Children, filePath)
			}
		case config.TypeFile:
			err = createFile(filePath)
		default:
			err = fmt.Errorf("unknown node type: %s (name: %s)", node.Type, node.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
