// internal/generator/template_source.go
package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateSource abstracts where templates come from
type TemplateSource interface {
	// ParseTemplates parses templates for a given language
	ParseTemplates(language string) (*template.Template, error)
	// ListLanguages returns available languages
	ListLanguages() ([]string, error)
	// ListTemplates returns template files for a language
	ListTemplates(language string) ([]string, error)
}

// EmbeddedTemplateSource uses embedded templates
type EmbeddedTemplateSource struct {
	fs embed.FS
}

func NewEmbeddedTemplateSource(fs embed.FS) *EmbeddedTemplateSource {
	return &EmbeddedTemplateSource{fs: fs}
}

func (e *EmbeddedTemplateSource) ParseTemplates(language string) (*template.Template, error) {
	patterns := []string{filepath.Join("templates", language, "*.tmpl")}
	return template.ParseFS(e.fs, patterns...)
}

func (e *EmbeddedTemplateSource) ListLanguages() ([]string, error) {
	entries, err := e.fs.ReadDir("templates")
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, entry := range entries {
		if entry.IsDir() {
			languages = append(languages, entry.Name())
		}
	}
	return languages, nil
}

func (e *EmbeddedTemplateSource) ListTemplates(language string) ([]string, error) {
	templateDir := filepath.Join("templates", language)
	entries, err := e.fs.ReadDir(templateDir)
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

// FileSystemTemplateSource uses filesystem templates
type FileSystemTemplateSource struct {
	baseDir string
}

func NewFileSystemTemplateSource(baseDir string) *FileSystemTemplateSource {
	return &FileSystemTemplateSource{baseDir: baseDir}
}

func (f *FileSystemTemplateSource) ParseTemplates(language string) (*template.Template, error) {
	langDir := filepath.Join(f.baseDir, language)

	// Check if language directory exists
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("language '%s' not found in %s", language, f.baseDir)
	}

	// Parse all .tmpl files in the language directory
	pattern := filepath.Join(langDir, "*.tmpl")
	return template.ParseGlob(pattern)
}

func (f *FileSystemTemplateSource) ListLanguages() ([]string, error) {
	entries, err := os.ReadDir(f.baseDir)
	if err != nil {
		return nil, fmt.Errorf("reading templates directory %s: %w", f.baseDir, err)
	}

	var languages []string
	for _, entry := range entries {
		if entry.IsDir() {
			languages = append(languages, entry.Name())
		}
	}
	return languages, nil
}

func (f *FileSystemTemplateSource) ListTemplates(language string) ([]string, error) {
	langDir := filepath.Join(f.baseDir, language)

	entries, err := os.ReadDir(langDir)
	if err != nil {
		return nil, fmt.Errorf("language '%s' not found in %s", language, f.baseDir)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			templates = append(templates, entry.Name())
		}
	}
	return templates, nil
}
