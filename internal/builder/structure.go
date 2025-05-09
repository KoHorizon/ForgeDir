package builder

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

// StructureBuilder builds a project scaffold.
type StructureBuilder struct {
	fs FileSystemCreator
}

func NewStructureBuilder(fs FileSystemCreator) *StructureBuilder {
	return &StructureBuilder{
		fs: fs,
	}
}

// Build uses the Config's Structure tree to instantiate folders & files under root.
func (b *StructureBuilder) Build(cfg *config.Config, root string) error {
	return b.createNodes(cfg.Structure, root)
}

// createNodes is the recursive guts of Build.
func (b *StructureBuilder) createNodes(nodes []config.StructureNode, currPath string) error {
	for _, node := range nodes {
		target := filepath.Join(currPath, node.Name)
		switch node.Type {
		case config.TypeDir:
			if err := b.fs.CreateFolder(target, DefaultFolderPermission); err != nil {
				return fmt.Errorf("mkdir %q: %w", target, err)
			}
			if err := b.createNodes(node.Children, target); err != nil {
				return err
			}
		case config.TypeFile:
			if err := b.fs.WriteFile(target, []byte{}, DefaultFilePermission); err != nil {
				return fmt.Errorf("touch %q: %w", target, err)
			}
		default:
			return fmt.Errorf("unknown node type %q for %q", node.Type, node.Name)
		}
	}
	return nil
}
