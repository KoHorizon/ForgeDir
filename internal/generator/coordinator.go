package generator

import (
	"errors"
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

type Coordinator struct {
	LanguageBoilerplate map[string]BoilerplateGenerator
}

func NewCoordinator(languages map[string]BoilerplateGenerator) *Coordinator {
	return &Coordinator{
		LanguageBoilerplate: languages,
	}
}

func (c *Coordinator) RunBoilerplateGeneration(cfg *config.Config, projectRoot string) error {
	generator, ok := c.LanguageBoilerplate[cfg.Language]
	if !ok {
		return errors.New("no boilerplate generator found")
	}
	err := generator.Generate(cfg, projectRoot)
	if err != nil {
		return fmt.Errorf("boilerplate generation failed for %s: %w", cfg.Language, err)
	}
	return nil
}
