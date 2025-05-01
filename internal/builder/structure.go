package builder

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

type CreateStructureOptions struct {
	FS   FileSystemCreator
	Cfg  *config.Config
	Root string
}

func CreateStructure(opts CreateStructureOptions) error {
	fs := opts.FS
	structure := opts.Cfg.Structure
	root := opts.Root
	return createStructureNodes(fs, structure, root)
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
