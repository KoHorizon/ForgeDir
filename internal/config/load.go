package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

func LoadConfigFromYaml(filename string) (*Config, error) {
	var cfg Config

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
