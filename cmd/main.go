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

	cfg, err := config.LoadConfigFromYaml(configName)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configName, err)
	}

	err = builder.CreateStructure(cfg, projectRoot)
	if err != nil {
		log.Fatalf("Failed to create strucure : %v", err)
	}

	// Declare and initialize the map with the instance directly
	goGenerator := &generator.GoGenerator{}
	availableGenerators := map[string]generator.BoilerplateGenerator{
		goGenerator.GetLanguage(): goGenerator,
	}
	// Inject all the generator as depedency for Coodinator
	coordinator := generator.NewCoordinator(availableGenerators)

	// Start of boilerplate generation
	fmt.Printf("Running boilerplate generation for language '%s'...\n", cfg.Language)
	err = coordinator.RunBoilerplateGeneration(cfg, projectRoot)
	if err != nil {
		// Log a fatal error if boilerplate generation fails
		log.Fatalf("Boilerplate generation failed: %v", err)
	}
	fmt.Println("Boilerplate generation completed successfully.")
	fmt.Println("ForgeDir finished project generation.")
}
