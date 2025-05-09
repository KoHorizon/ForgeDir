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
