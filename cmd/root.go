// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	outputDir    string
	templatesDir string // New: custom templates directory
)

// rootCmd is now just the top‐level command (no Run or RunE)
var rootCmd = &cobra.Command{
	Use:   "fgdir",
	Short: "Scaffold a project structure from your YAML spec",
}

func init() {
	// Disabling the built-in completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Removing the built-in help "sub" command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Add global flag for custom templates directory
	rootCmd.PersistentFlags().StringVarP(
		&templatesDir, "templates", "t", "",
		"path to custom templates directory (default: use built-in templates)",
	)
}

// Execute runs the CLI.
func Execute() {
	rootCmd.Execute()
}
