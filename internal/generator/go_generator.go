package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

//go:embed templates/go/*.tmpl
var goTemplates embed.FS

// tmpl holds parsed Go templates from the embedded filesystem.
var tmpl = template.Must(
	template.New("go").ParseFS(goTemplates, "templates/go/*.tmpl"),
)

// GoGenerator implements BoilerplateGenerator using external Go templates.
type GoGenerator struct{}

var _ BoilerplateGenerator = (*GoGenerator)(nil)

// GetLanguage returns the language this generator supports.
func (g *GoGenerator) GetLanguage() string { return "go" }

// init registers this generator during package initialization.
func init() {
	Register(&GoGenerator{})
}

// Generate scans for .go files under projectRoot and applies matching templates
// or falls back to a minimal stub if no template is defined.
func (g *GoGenerator) Generate(cfg *config.Config, projectRoot string) error {
	fmt.Printf("GoGenerator: Generating boilerplate for project '%s' in %s\n",
		cfg.ProjectName, projectRoot)

	return filepath.WalkDir(projectRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only process .go files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		// Build tmpl name, e.g. "main.go.tmpl"
		tmplName := d.Name() + ".tmpl"
		tpl := tmpl.Lookup(tmplName)

		// Prepare common data
		pkg := determinePackageName(path, projectRoot)
		name := strings.TrimSuffix(d.Name(), ".go")

		// If we don’t have a specific template, write the default stub
		if tpl == nil {
			defaultContent := fmt.Sprintf(
				"package %s\n\n// TODO: implement %s\n",
				pkg, name,
			)
			return os.WriteFile(path, []byte(defaultContent), 0644)
		}

		// Otherwise render via text/template
		data := struct {
			PackageName string
			FileName    string
		}{
			PackageName: pkg,
			FileName:    name,
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("error executing template %s: %w", tmplName, err)
		}

		return os.WriteFile(path, buf.Bytes(), 0644)
	})
}

// determinePackageName infers the Go package name based on directory structure.
func determinePackageName(filePath, root string) string {
	dir := filepath.Dir(filePath)
	rel, err := filepath.Rel(root, dir)
	if err != nil || rel == "." {
		return "main"
	}
	if strings.HasPrefix(rel, "cmd") {
		return "main"
	}
	return filepath.Base(dir)
}
