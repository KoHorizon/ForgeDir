package main

import (
	"log"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
)

func main() {
	configName := "config.yaml"

	cfg, err := config.LoadConfigFromYaml(configName)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configName, err)
	}

	err = builder.CreateStructure(cfg, "../tmp/generated-structure")
	if err != nil {
		log.Fatalf("Failed to create strucure : %v", err)
	}
}
