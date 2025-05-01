package main

import (
	"fmt"
	"log"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/generator"
)

func main() {
	configName := "./config.yaml"
	projectRoot := "../tmp/generated-structure"

	// Load the project configuration from YAML
	cfg, err := config.LoadConfigFromYaml(configName)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configName, err)
	}

	// Build the project directory structure
	fs := builder.NewOSFileSystemCreator()
	opts := builder.CreateStructureOptions{
		FS:   fs,
		Cfg:  cfg,
		Root: projectRoot,
	}
	if err := builder.CreateStructure(opts); err != nil {
		log.Fatalf("Failed to create strucure : %v", err)
	}

	// Prepare and inject all available code generators into the coordinator
	availableGenerators := generator.All()
	coordinator := generator.NewCoordinator(availableGenerators)

	// Generate language-specific boilerplate
	fmt.Printf("Running boilerplate generation for language '%s'...\n", cfg.Language)
	err = coordinator.RunBoilerplateGeneration(cfg, projectRoot)
	if err != nil {
		// Log a fatal error if boilerplate generation fails
		log.Fatalf("Boilerplate generation failed: %v", err)
	}

	// Generate language-specific boilerplate
	fmt.Println("Boilerplate generation completed successfully.")
	fmt.Println("ForgeDir finished project generation.")
}
