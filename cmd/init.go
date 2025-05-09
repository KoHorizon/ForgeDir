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
		fs := builder.NewOSFileSystemCreator()
		sb := builder.NewStructureBuilder(fs)
		if err := sb.Build(cfg, outputDir); err != nil {
			return fmt.Errorf("creating structure: %w", err)
		}

		// 3. Generate boilerplate
		coord := generator.NewCoordinator(generator.All())
		fmt.Printf("Generating boilerplate for %q in %s …\n", cfg.Language, outputDir)
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
