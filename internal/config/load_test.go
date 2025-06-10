package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

func TestLoadConfigFromYaml_Success(t *testing.T) {
	yamlContent := `
projectName: testproj
language: go
structure:
  - type: dir
    name: foo
    children:
      - type: file
        name: bar.go
`
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.yaml")
	if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("unable to write temp yaml: %v", err)
	}

	cfg, err := config.LoadConfigFromYaml(path)
	if err != nil {
		t.Fatalf("LoadConfigFromYaml returned unexpected error: %v", err)
	}

	if cfg.ProjectName != "testproj" {
		t.Errorf("expected ProjectName %q, got %q", "testproj", cfg.ProjectName)
	}
	if cfg.Language != "go" {
		t.Errorf("expected Language %q, got %q", "go", cfg.Language)
	}
	if len(cfg.Structure) != 1 {
		t.Fatalf("expected 1 structure node, got %d", len(cfg.Structure))
	}
	node := cfg.Structure[0]
	if node.Type != config.TypeDir || node.Name != "foo" {
		t.Errorf("unexpected node: %+v", node)
	}
	if len(node.Children) != 1 || node.Children[0].Name != "bar.go" {
		t.Errorf("unexpected children: %+v", node.Children)
	}
}

func TestLoadConfigFromYaml_FileNotFound(t *testing.T) {
	_, err := config.LoadConfigFromYaml("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
