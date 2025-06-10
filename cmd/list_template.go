// Copyright © 2025 KoHorizon
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"sort"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/generator"
	"github.com/spf13/cobra"
)

var listTemplatesCmd = &cobra.Command{
	Use:                   "list-templates [language]",
	Short:                 "List the built-in templates (or those for a given language)",
	Args:                  cobra.MaximumNArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		templateSource, err := generator.CreateTemplateSource(templatesDir)
		if err != nil {
			fmt.Printf("❌ Error setting up templates: %v\n", err)
			return
		}

		// Create factory to get available generators
		fs := builder.NewOSFileSystem()
		factory := generator.NewGeneratorFactory(fs, templateSource)
		generators, err := factory.CreateAvailableGenerators()
		if err != nil {
			fmt.Printf("❌ Error loading templates: %v\n", err)
			return
		}

		// If no language specified, list all available languages
		if len(args) == 0 {
			listAllLanguages(generators)
			return
		}

		// List templates for specific language
		language := args[0]
		listTemplatesForLanguage(language, generators)
	},
}

// listAllLanguages shows all available languages
func listAllLanguages(generators []generator.Generator) {
	if len(generators) == 0 {
		fmt.Println("No templates available")
		return
	}

	fmt.Println("Available languages:")

	// Sort languages for consistent output
	var languages []string
	for _, gen := range generators {
		languages = append(languages, gen.GetLanguage())
	}
	sort.Strings(languages)

	for _, lang := range languages {
		fmt.Printf("  %s\n", lang)
	}

	fmt.Printf("\nUse 'fgdir list-templates <language>' to see templates for a specific language.\n")
}

// listTemplatesForLanguage shows templates for a specific language
func listTemplatesForLanguage(language string, generators []generator.Generator) {
	// Find the generator for this language
	var targetGen generator.Generator
	for _, gen := range generators {
		if gen.GetLanguage() == language {
			targetGen = gen
			break
		}
	}

	if targetGen == nil {
		fmt.Printf("❌ Language '%s' not supported\n", language)
		fmt.Println("\nAvailable languages:")
		for _, gen := range generators {
			fmt.Printf("  %s\n", gen.GetLanguage())
		}
		return
	}

	// Get templates for this language
	templates, err := getTemplatesForLanguage(language)
	if err != nil {
		fmt.Printf("❌ Error reading templates for '%s': %v\n", language, err)
		return
	}

	if len(templates) == 0 {
		fmt.Printf("No templates found for '%s'\n", language)
		return
	}

	fmt.Printf("Templates for '%s':\n", language)
	sort.Strings(templates)
	for _, tmpl := range templates {
		fmt.Printf("  %s\n", tmpl)
	}
}

// getTemplatesForLanguage reads template files from embedded filesystem
func getTemplatesForLanguage(language string) ([]string, error) {
	return generator.GetTemplatesForLanguage(language)
}

func init() {
	rootCmd.AddCommand(listTemplatesCmd)
}
