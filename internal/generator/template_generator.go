package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
	"github.com/KoHorizon/ForgeDir/internal/utils"
)

// GenericGenerator uses embedded templates for boilerplate generation.
type GenericGenerator struct {
	lang string
	tmpl *template.Template
	fs   builder.FileSystem
}

// NewGenericGenerator initializes a GenericGenerator for the given language.
func NewGenericGenerator(lang string, fs builder.FileSystem) (*GenericGenerator, error) {
	patterns := []string{filepath.Join("templates", lang, "*.tmpl")}
	parsed, err := template.ParseFS(tmplFS, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parsing templates for %q: %w", lang, err)
	}
	return &GenericGenerator{lang: lang, tmpl: parsed, fs: fs}, nil
}

// GetLanguage returns the generator's language.
func (g *GenericGenerator) GetLanguage() string {
	return g.lang
}

// Generate renders templates into existing files under root.
// Only iterate the files listed in your config.Structure
func (g *GenericGenerator) Generate(cfg *config.Config, root string) error {
	// 1. flatten your YAML tree into a list of relative file paths
	var files []string

	// Simple DFS with an inline stack of (nodes, basePath)
	stack := []struct {
		nodes []config.StructureNode
		base  string
	}{{cfg.Structure, ""}}

	for len(stack) > 0 {
		// pop
		frame := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, n := range frame.nodes {
			rel := filepath.Join(frame.base, n.Name)
			if n.Type == config.TypeDir {
				// push child directory
				stack = append(stack, struct {
					nodes []config.StructureNode
					base  string
				}{n.Children, rel})
			} else if n.Type == config.TypeFile {
				files = append(files, rel)
			}
		}
	}

	// 2. render templates only for the files your spec asked for
	for _, rel := range files {
		target := filepath.Join(root, rel)
		if err := g.generateFile(target, root); err != nil {
			return err
		}
	}
	return nil
}

// generateFile tries in order: file-specific then catch-all (_.tmpl).
// It passes generic data into the template, not Go-specific.
func (g *GenericGenerator) generateFile(path, root string) error {
	name := filepath.Base(path)

	// lookup order: specific then catch-all
	tpl := g.tmpl.Lookup(name + ".tmpl")
	if tpl == nil {
		tpl = g.tmpl.Lookup("(default).tmpl")
	}

	// Prepare generic template data
	// Language: the configured language
	// DirName: the name of the file's directory under root (or empty)
	// FileName: the base name without extension
	dir := filepath.Dir(path)
	rel, _ := filepath.Rel(root, dir)
	d := ""
	if rel != "." {
		// take only the last segment
		parts := strings.Split(rel, string(filepath.Separator))
		d = parts[len(parts)-1]
	}
	data := struct {
		Language string
		DirName  string
		FileName string
	}{
		Language: g.lang,
		DirName:  d,
		FileName: strings.TrimSuffix(name, filepath.Ext(name)),
	}

	var content []byte
	if tpl != nil {
		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, data); err != nil {
			return fmt.Errorf("executing template %q: %w", name, err)
		}
		content = buf.Bytes()
	} else {
		content = []byte(fmt.Sprintf("// no template for %s\n", name))
	}

	// Ensure directory exists and write file
	dirPath := filepath.Dir(path)
	if err := g.fs.CreateFolder(dirPath, builder.DefaultFolderPermission); err != nil {
		return fmt.Errorf("creating folder %q: %w", dirPath, err)
	}
	if err := g.fs.WriteFile(path, content, builder.DefaultFilePermission); err != nil {
		return fmt.Errorf("writing file %q: %w", path, err)
	}

	return nil
}

// GetTemplatesForLanguage returns a list of template files for the given language
func GetTemplatesForLanguage(language string) ([]string, error) {
	templateDir := filepath.Join("templates", language)

	entries, err := tmplFS.ReadDir(templateDir)
	if err != nil {
		return nil, fmt.Errorf("language '%s' not found", language)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			templates = append(templates, entry.Name())
		}
	}

	return templates, nil
}

// CreateTemplateSource creates the appropriate template source based on templatesDir
func CreateTemplateSource(customTemplatesDir string) (TemplateSource, error) {
	if customTemplatesDir == "" {
		// Use embedded templates
		return NewEmbeddedTemplateSource(tmplFS), nil
	}

	// Expand path (handles ~ and relative paths)
	expandedPath, err := utils.ExpandPath(customTemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("expanding templates path '%s': %w", customTemplatesDir, err)
	}

	// Validate custom directory exists
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("custom templates directory '%s' does not exist", expandedPath)
	}

	return NewFileSystemTemplateSource(expandedPath), nil
}
