// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/utils"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:                   "validate [spec.yaml]",
	Short:                 "Validate that a spec.yaml is well-formed",
	Args:                  cobra.MaximumNArgs(1),
	DisableFlagsInUseLine: true, // Hide "[flags]" in usage
	SilenceUsage:          true, // Don't show usage on errors
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine config file path
		configPath := "config.yaml" // default
		if len(args) == 1 {
			configPath = args[0]
		}

		// Try to load the config
		cfg, err := config.LoadConfigFromYaml(configPath)
		if err != nil {
			fmt.Printf("❌ Validation failed: %v\n", err)
			return fmt.Errorf("invalid configuration")
		}

		// Perform additional validation
		if err := validateConfig(cfg); err != nil {
			fmt.Printf("❌ Validation failed: %v\n", err)
			return fmt.Errorf("configuration validation errors")
		}

		// Success!
		fmt.Printf("✅ Configuration '%s' is valid\n", configPath)
		fmt.Printf("   Project: %s\n", cfg.ProjectName)
		fmt.Printf("   Language: %s\n", cfg.Language)
		fmt.Printf("   Structure nodes: %d\n", countNodes(cfg.Structure))

		return nil
	},
}

// validateConfig performs business logic validation beyond YAML parsing
func validateConfig(cfg *config.Config) error {
	// Check required fields
	if cfg.ProjectName == "" {
		return fmt.Errorf("projectName is required")
	}

	if cfg.Language == "" {
		return fmt.Errorf("language is required")
	}

	if len(cfg.Structure) == 0 {
		return fmt.Errorf("structure cannot be empty")
	}

	// Validate structure nodes (including path security)
	return validateStructureNodes(cfg.Structure, "")
}

// validateStructureNodes validates each node in the structure tree
func validateStructureNodes(nodes []config.StructureNode, path string) error {
	for _, node := range nodes {
		currentPath := path + "/" + node.Name

		// Validate node type
		if node.Type != config.TypeDir && node.Type != config.TypeFile {
			return fmt.Errorf("invalid type '%s' at %s (must be 'dir' or 'file')", node.Type, currentPath)
		}

		// Validate node name
		if node.Name == "" {
			return fmt.Errorf("name is required at %s", currentPath)
		}

		// Security validation: Check for path traversal and other security issues
		if err := utils.ValidatePath(node.Name); err != nil {
			return fmt.Errorf("security validation failed for '%s' at %s: %w", node.Name, currentPath, err)
		}

		// Files cannot have children
		if node.Type == config.TypeFile && len(node.Children) > 0 {
			return fmt.Errorf("file '%s' cannot have children", currentPath)
		}

		// Recursively validate children
		if len(node.Children) > 0 {
			if err := validateStructureNodes(node.Children, currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// countNodes counts total nodes in the structure tree
func countNodes(nodes []config.StructureNode) int {
	count := len(nodes)
	for _, node := range nodes {
		count += countNodes(node.Children)
	}
	return count
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
