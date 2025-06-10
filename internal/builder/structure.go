package builder

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/utils"
)

// StructureBuilder builds a project scaffold.
type StructureBuilder struct {
	fs          FileSystem
	projectRoot string // Store the absolute project root for validation
}

func NewStructureBuilder(fs FileSystem) *StructureBuilder {
	return &StructureBuilder{
		fs: fs,
	}
}

// Build uses the Config's Structure tree to instantiate folders & files under root.
func (b *StructureBuilder) Build(cfg *config.Config, root string) error {
	// Store absolute project root for path validation
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to get absolute project root: %w", err)
	}
	b.projectRoot = absRoot

	// Validate the entire structure before creating anything
	if err := b.validateStructure(cfg.Structure, ""); err != nil {
		return fmt.Errorf("structure validation failed: %w", err)
	}

	return b.createNodes(cfg.Structure, root)
}

// validateStructure performs a pre-flight validation of all path names in the structure
func (b *StructureBuilder) validateStructure(nodes []config.StructureNode, currentPath string) error {
	for _, node := range nodes {
		// Validate the node name itself
		if err := utils.ValidatePath(node.Name); err != nil {
			return fmt.Errorf("invalid name '%s' at path '%s': %w", node.Name, currentPath, err)
		}

		// Build the target path for this node
		targetPath := filepath.Join(currentPath, node.Name)

		// Validate that the resulting path would be safe
		fullPath := filepath.Join(b.projectRoot, targetPath)
		if _, err := utils.SanitizePath(b.projectRoot, fullPath); err != nil {
			return fmt.Errorf("unsafe path '%s': %w", targetPath, err)
		}

		// Recursively validate children
		if len(node.Children) > 0 {
			if err := b.validateStructure(node.Children, targetPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// createNodes is the recursive guts of Build.
func (b *StructureBuilder) createNodes(nodes []config.StructureNode, currPath string) error {
	for _, node := range nodes {
		// Validate the individual path component (redundant but defensive)
		if err := utils.ValidatePath(node.Name); err != nil {
			return fmt.Errorf("invalid path component '%s': %w", node.Name, err)
		}

		target := filepath.Join(currPath, node.Name)

		// Additional safety check: ensure target is within project root
		safeTarget, err := utils.SanitizePath(b.projectRoot, target)
		if err != nil {
			return fmt.Errorf("path safety check failed for '%s': %w", target, err)
		}

		switch node.Type {
		case config.TypeDir:
			if err := b.fs.CreateFolder(safeTarget, DefaultFolderPermission); err != nil {
				return fmt.Errorf("mkdir %q: %w", safeTarget, err)
			}
			if err := b.createNodes(node.Children, target); err != nil {
				return err
			}
		case config.TypeFile:
			if err := b.fs.WriteFile(safeTarget, []byte{}, DefaultFilePermission); err != nil {
				return fmt.Errorf("touch %q: %w", safeTarget, err)
			}
		default:
			return fmt.Errorf("unknown node type %q for %q", node.Type, node.Name)
		}
	}
	return nil
}
