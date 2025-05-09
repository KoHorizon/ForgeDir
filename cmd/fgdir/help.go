// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help for fgdir (overide)",
	Run: func(cmd *cobra.Command, args []string) {
		// Print custom help layout here:
		fmt.Println(`
fgdir is a project-scaffolding CLI.

Usage:
  fgdir [command]

Available Commands:
  help               Show help about the tool
  init               Read a YAML spec and scaffold the project
  validate           Validate that a spec.yaml is well-formed
  list-templates     List the built-in templates (or those for a given language)
  clean              Clean up a previously scaffolded project
  version            Show the current version of the CLI

Use "fgdir [command] --help" for more information about a command.
        `)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
