// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// These variables can be set at build time using ldflags
var (
	version = "" // Will be set by ldflags or detected from build info
	commit  = "unknown"
	date    = "unknown"
)

// getVersion returns the version, trying multiple sources
func getVersion() string {
	// If version was set via ldflags, use it
	if version != "" {
		return version
	}

	// Try to get version from build info (works with go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		// For tagged releases, go install sets this automatically
		if info.Main.Version != "(devel)" && info.Main.Version != "" {
			return info.Main.Version
		}
	}

	// Fallback for development builds
	return "dev"
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ForgeDir %s\n", getVersion())

		// Show additional info if available
		if commit != "unknown" || date != "unknown" {
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built:  %s\n", date)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
