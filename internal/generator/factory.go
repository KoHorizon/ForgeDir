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
	fs             builder.FileSystem
	templateSource TemplateSource
}

func NewGeneratorFactory(fs builder.FileSystem, templateSource TemplateSource) *GeneratorFactory {
	return &GeneratorFactory{
		fs:             fs,
		templateSource: templateSource,
	}
}

// CreateAvailableGenerators scans available templates and creates generators
func (f *GeneratorFactory) CreateAvailableGenerators() ([]Generator, error) {
	languages, err := f.templateSource.ListLanguages()
	if err != nil {
		return nil, fmt.Errorf("reading languages: %w", err)
	}

	var generators []Generator
	for _, lang := range languages {
		gen, err := f.createGeneratorForLanguage(lang)
		if err != nil {
			return nil, fmt.Errorf("creating generator for %s: %w", lang, err)
		}
		generators = append(generators, gen)
	}

	return generators, nil
}

func (f *GeneratorFactory) createGeneratorForLanguage(language string) (Generator, error) {
	tmpl, err := f.templateSource.ParseTemplates(language)
	if err != nil {
		return nil, fmt.Errorf("parsing templates for %s: %w", language, err)
	}

	return &GenericGenerator{
		lang: language,
		tmpl: tmpl,
		fs:   f.fs,
	}, nil
}

// GetTemplatesForLanguage returns templates for a specific language
func (f *GeneratorFactory) GetTemplatesForLanguage(language string) ([]string, error) {
	return f.templateSource.ListTemplates(language)
}
