// internal/generator/coordinator_test.go
package generator_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/generator"
)

type dummyGen struct {
	lang   string
	called bool
	err    error
}

func (d *dummyGen) GetLanguage() string { return d.lang }
func (d *dummyGen) Generate(cfg *config.Config, root string) error {
	d.called = true
	return d.err
}

func TestRunBoilerplate_Success(t *testing.T) {
	gen := &dummyGen{lang: "go"}
	coord := generator.NewCoordinator([]generator.Generator{gen})

	err := coord.RunBoilerplateGeneration(&config.Config{Language: "go"}, "/some/root")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !gen.called {
		t.Error("expected generator.Generate to be called")
	}
}

func TestRunBoilerplate_NoGenerator(t *testing.T) {
	coord := generator.NewCoordinator(nil)
	err := coord.RunBoilerplateGeneration(&config.Config{Language: "js"}, "/any")
	if err == nil || !strings.Contains(err.Error(), "no boilerplate generator found") {
		t.Errorf("expected no-generator error, got %v", err)
	}
}

func TestRunBoilerplate_GeneratorError(t *testing.T) {
	gen := &dummyGen{lang: "py", err: errors.New("boom")}
	coord := generator.NewCoordinator([]generator.Generator{gen})

	err := coord.RunBoilerplateGeneration(&config.Config{Language: "py"}, "/root")
	if err == nil || !strings.Contains(err.Error(), "boilerplate generation failed") {
		t.Errorf("expected generate-error, got %v", err)
	}
}
