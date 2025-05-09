// cmd/init.go
package cmd

import (
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/generator"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <spec.yaml>",
	Short: "Read a YAML spec and scaffold the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// allow overriding the config flag by positional arg
		cfgFile = args[0]

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
	// if needed, you can add init-specific flags here (but config/output are global)
	// e.g. initCmd.Flags().Bool("dry-run", false, "simulate only")
}
