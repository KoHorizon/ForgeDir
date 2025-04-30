package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KoHorizon/ForgeDir/internal/config"
)

type GoGenerator struct{}

var _ BoilerplateGenerator = (*GoGenerator)(nil)

// GetLanguage returns the language this generator supports.
func (g *GoGenerator) GetLanguage() string {
	return "go"
}

// Generate scans the generated structure and injects boilerplate content for each .go file.
func (g *GoGenerator) Generate(cfg *config.Config, projectRoot string) error {
	fmt.Printf("GoGenerator: Generating boilerplate for project '%s' in %s\n", cfg.ProjectName, projectRoot)

	// structureRoot := filepath.Join(projectRoot, "tmp", "generated-structure")
	goFiles, err := collectGoFiles(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to collect .go files: %w", err)
	}

	for _, file := range goFiles {
		pkg := determinePackageName(file, projectRoot)
		fileName := filepath.Base(file)

		content, err := generateContentForFile(pkg, fileName)
		if err != nil {
			return fmt.Errorf("failed to generate content for %s: %w", file, err)
		}

		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %w", file, err)
		}
	}

	return nil
}

// collectGoFiles recursively finds all .go files under a given root directory.
func collectGoFiles(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// determinePackageName infers the Go package name based on file's directory structure.
func determinePackageName(filePath, root string) string {
	dir := filepath.Dir(filePath)

	relPath, err := filepath.Rel(root, dir)
	if err != nil || relPath == "." {
		return "main"
	}

	if strings.HasPrefix(relPath, "cmd") {
		return "main"
	}

	return filepath.Base(dir)
}

// generateContentForFile chooses the correct template based on filename and package.
func generateContentForFile(pkgName, fileName string) (string, error) {
	if genFunc, ok := fileTemplateMap[fileName]; ok {
		return genFunc(pkgName), nil
	}

	// Fallback template
	return fmt.Sprintf(`package %s

// TODO: implement %s
`, pkgName, fileName), nil
}

// fileTemplateMap maps specific filenames to content generators.
var fileTemplateMap = map[string]func(pkg string) string{
	"main.go": func(pkg string) string {
		return fmt.Sprintf(`package %s

import "fmt"

func main() {
	fmt.Println("This is the main entry point!")
}
`, pkg)
	},

	"handler.go": func(pkg string) string {
		return fmt.Sprintf(`package %s

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from handler"))
}
`, pkg)
	},

	"service.go": func(pkg string) string {
		return fmt.Sprintf(`package %s

type Service struct {}

func (s *Service) DoSomething() string {
	return "Service logic here"
}
`, pkg)
	},

	"controller.go": func(pkg string) string {
		return fmt.Sprintf(`package %s

type Controller struct {}

func (c *Controller) HandleRequest() {
	// TODO: implement controller logic
}
`, pkg)
	},
}
