// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/generator"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [spec.yaml]",
	Short: "Read a YAML spec and scaffold the project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			cfgFile = args[0]
		}
		outputDir, _ = filepath.Abs(outputDir)

		// 1. Load config
		cfg, err := config.LoadConfigFromYaml(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config %q: %w", cfgFile, err)
		}

		// 2. Build file tree
		fs := builder.NewOSFileSystem()
		sb := builder.NewStructureBuilder(fs)
		if err := sb.Build(cfg, outputDir); err != nil {
			return fmt.Errorf("creating structure: %w", err)
		}

		// 3. Generate boilerplate
		fmt.Printf("Generating boilerplate for %q in %s …\n", cfg.Language, outputDir)
		templateSource, err := generator.CreateTemplateSource(templatesDir)
		if err != nil {
			return fmt.Errorf("setting up templates: %w", err)
		}
		factory := generator.NewGeneratorFactory(fs, templateSource)
		generators, err := factory.CreateAvailableGenerators()
		if err != nil {
			return fmt.Errorf("creating generators: %w", err)
		}

		coord := generator.NewCoordinator(generators)
		if err := coord.RunBoilerplateGeneration(cfg, outputDir); err != nil {
			return fmt.Errorf("boilerplate generation failed: %w", err)
		}

		fmt.Println("✅ ForgeDir finished project generation.")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(
		&cfgFile, "config", "c", "config.yaml",
		"path to the YAML project spec",
	)
	initCmd.Flags().StringVarP(
		&outputDir, "output", "o", ".",
		"directory where the project will be generated (default is current directory)",
	)

	rootCmd.AddCommand(initCmd)
}
