package generator

import (
	"embed"
	"fmt"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
)

// Generator is your interface for scaffolding.
type Generator interface {
	GetLanguage() string
	Generate(cfg *config.Config, root string) error
}

//go:embed templates/*/*.tmpl
var tmplFS embed.FS

// GeneratorFactory creates generators for available languages
type GeneratorFactory struct {
	fs builder.FileSystem
}

func NewGeneratorFactory(fs builder.FileSystem) *GeneratorFactory {
	return &GeneratorFactory{fs: fs}
}

// CreateAvailableGenerators scans embedded templates and creates generators
func (f *GeneratorFactory) CreateAvailableGenerators() ([]Generator, error) {
	entries, err := tmplFS.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("reading templates: %w", err)
	}

	var generators []Generator
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		lang := entry.Name()
		gen, err := NewGenericGenerator(lang, f.fs)
		if err != nil {
			return nil, fmt.Errorf("creating generator for %s: %w", lang, err)
		}
		generators = append(generators, gen)
	}

	return generators, nil
}
