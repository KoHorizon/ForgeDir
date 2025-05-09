package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
)

//go:embed templates/*/*.tmpl
var tmplFS embed.FS

// GenericGenerator uses embedded templates for boilerplate generation.
type GenericGenerator struct {
	lang string
	tmpl *template.Template
	fs   builder.FileSystemCreator
}

// NewGenericGenerator initializes a GenericGenerator for the given language.
func NewGenericGenerator(lang string, fs builder.FileSystemCreator) (*GenericGenerator, error) {
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
func (g *GenericGenerator) Generate(cfg *config.Config, root string) error {
	return filepath.WalkDir(root, g.walkFunc(root))
}

func (g *GenericGenerator) walkFunc(root string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		return g.generateFile(path, root)
	}
}

// generateFile tries in order: file-specific and then catch-all (_.tmpl).
func (g *GenericGenerator) generateFile(path, root string) error {
	name := filepath.Base(path)

	// lookup order: specific then catch-all
	tpl := g.tmpl.Lookup(name + ".tmpl")
	if tpl == nil {
		tpl = g.tmpl.Lookup("(default).tmpl")
	}

	// Prepare template data
	data := struct {
		PackageName string
		FileName    string
	}{
		PackageName: determinePackageName(path, root),
		FileName:    strings.TrimSuffix(name, filepath.Ext(name)),
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
	dir := filepath.Dir(path)
	if err := g.fs.CreateFolder(dir, builder.DefaultFolderPermission); err != nil {
		return fmt.Errorf("creating folder %q: %w", dir, err)
	}
	if err := g.fs.WriteFile(path, content, builder.DefaultFilePermission); err != nil {
		return fmt.Errorf("writing file %q: %w", path, err)
	}

	return nil
}

func init() {
	fsCreator := builder.NewOSFileSystemCreator()
	entries, err := tmplFS.ReadDir("templates")
	if err != nil {
		panic(fmt.Errorf("reading templates: %w", err))
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		lang := e.Name()
		gen, err := NewGenericGenerator(lang, fsCreator)
		if err != nil {
			panic(fmt.Errorf("loading %q templates: %w", lang, err))
		}
		Register(gen)
	}
}

// determinePackageName infers the Go package name based on directory structure.
func determinePackageName(filePath, root string) string {
	dir := filepath.Dir(filePath)
	rel, err := filepath.Rel(root, dir)
	if err != nil || rel == "." || strings.HasPrefix(rel, "cmd") {
		return "main"
	}
	return filepath.Base(dir)
}
