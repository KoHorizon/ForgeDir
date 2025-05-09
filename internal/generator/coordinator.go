// coordinator.go
package generator

import (
	"errors"
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

type Coordinator struct {
	LanguageBoilerplate map[string]Generator
}

// Now take a slice of Generator (which matches registry.All())
func NewCoordinator(gens []Generator) *Coordinator {
	m := make(map[string]Generator, len(gens))
	for _, g := range gens {
		m[g.GetLanguage()] = g
	}
	return &Coordinator{LanguageBoilerplate: m}
}

func (c *Coordinator) RunBoilerplateGeneration(cfg *config.Config, projectRoot string) error {
	gen, ok := c.LanguageBoilerplate[cfg.Language]
	if !ok {
		return errors.New("no boilerplate generator found for " + cfg.Language)
	}
	if err := gen.Generate(cfg, projectRoot); err != nil {
		return fmt.Errorf("boilerplate generation failed for %s: %w", cfg.Language, err)
	}
	return nil
}
