package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type Node struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	Children []Node `yaml:"children, omitempty"`
}

type Config struct {
	ProjectName string `yaml:"projectName"`
	Structure   []Node `yaml:"structure"`
}

func readConfig(filename string) (*Config, error) {
	// Read the YAML filename
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}
	return &config, nil
}

func main() {

	c, err := readConfig("cmd/config.yaml")
	if err != nil {
		fmt.Printf("Name: %s\n", c.ProjectName)
	}
	fmt.Println(len(c.Structure))
	for _, node := range c.Structure {
		fmt.Printf("- Type: %s, Name: %s, Children: %d\n", node.Type, node.Name, len(node.Children))
	}
}
