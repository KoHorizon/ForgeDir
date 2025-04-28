package main

import (
	"fmt"
	"log"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

func main() {
	configName := "cmd/config.yaml"

	cfg, err := config.LoadConfigFromYaml(configName)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configName, err)
	}

	for _, node := range cfg.Structure {
		fmt.Printf("- Type: %s, Name: %s, Children: %d\n", node.Type, node.Name, len(node.Children))
	}
}
