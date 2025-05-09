// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	outputDir string
)

// rootCmd is now just the top‐level command (no Run or RunE)
var rootCmd = &cobra.Command{
	Use:   "fgdir",
	Short: "Scaffold a project structure from your YAML spec",
	// no Run/RuneE here
}

func init() {
	// Persistent flags are available to all subcommands and
	// control global behavior of the CLI.
	// global flags
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "config.yaml",
		"path to the YAML project spec",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputDir, "output", "o", "./tmp/generated-structure",
		"directory where the project will be generated",
	)
	// register subcommands
	rootCmd.AddCommand(initCmd)
}

// Execute runs the CLI.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
