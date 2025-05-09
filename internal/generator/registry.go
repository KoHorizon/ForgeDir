package generator

import "github.com/KoHorizon/ForgeDir/internal/config"

// Generator is your interface for scaffolding.
type Generator interface {
	GetLanguage() string
	Generate(cfg *config.Config, root string) error
}

var registry = make(map[string]Generator)

// Register adds a Generator to the map by its Language key.
func Register(g Generator) {
	registry[g.GetLanguage()] = g
}

// All returns every registered Generator.
func All() []Generator {
	gens := make([]Generator, 0, len(registry))
	for _, g := range registry {
		gens = append(gens, g)
	}
	return gens
}
