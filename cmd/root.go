// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/generator"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	outputDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fgdir",
	Short: "Scaffold a project structure from your YAML spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Load
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

		// 3. Wire generators & run boilerplate
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
	// Persistent flags are available to all subcommands and
	// control global behavior of the CLI.
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c",
		"config.yaml",
		"path to the YAML project spec",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputDir,
		"output", "o",
		"./tmp/generated-structure",
		"directory where the project will be generated",
	)

	// Local flags apply only to the root command.
	// (You can remove this if you’re not actually using “toggle”.)
	rootCmd.Flags().BoolP(
		"toggle", "t",
		false,
		"help message for toggle",
	)
}

// Execute runs the CLI.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
