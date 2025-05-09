// internal/builder/structure_builder_test.go
package builder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
)

type fakeFS struct {
	Folders []string
	Files   []string
}

func (f *fakeFS) CreateFolder(path string, perm os.FileMode) error {
	f.Folders = append(f.Folders, path)
	return nil
}

func (f *fakeFS) WriteFile(path string, content []byte, perm os.FileMode) error {
	f.Files = append(f.Files, path)
	return nil
}

func TestStructureBuilder_Build(t *testing.T) {
	fs := &fakeFS{}
	sb := builder.NewStructureBuilder(fs)

	// simple tree: root/a/ then file root/a/b.txt
	cfg := &config.Config{
		Structure: []config.StructureNode{
			{
				Type: config.TypeDir, Name: "a",
				Children: []config.StructureNode{
					{Type: config.TypeFile, Name: "b.txt"},
				},
			},
		},
	}
	root := "root"

	if err := sb.Build(cfg, root); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expectedDir := filepath.Join(root, "a")
	expectedFile := filepath.Join(root, "a", "b.txt")

	if len(fs.Folders) != 1 || fs.Folders[0] != expectedDir {
		t.Errorf("expected folder %q, got %v", expectedDir, fs.Folders)
	}
	if len(fs.Files) != 1 || fs.Files[0] != expectedFile {
		t.Errorf("expected file %q, got %v", expectedFile, fs.Files)
	}
}
