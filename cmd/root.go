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
}

func init() {
	// Disabling the built-in completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Romoving the built-in help “sub” command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Custom commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(helpCmd)
}

// Execute runs the CLI.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
