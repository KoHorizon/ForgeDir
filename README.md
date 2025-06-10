# ForgeDir

[![Go Report Card](https://goreportcard.com/badge/github.com/KoHorizon/ForgeDir)](https://goreportcard.com/report/github.com/KoHorizon/ForgeDir)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ForgeDir is a CLI tool written in Go that scaffolds project structures from simple YAML specifications. It generates both directory/file structures and language-specific boilerplate code using customizable templates. Designed for developers who want to standardize their project layouts and reduce repetitive setup work.

---

## Table of Contents

1. [Features](#features)
2. [Quick Start](#quick-start)
3. [Installation](#installation)
4. [Usage Examples](#usage-examples)
5. [Configuration (`config.yaml`)](#configuration-configyaml)
6. [Custom Templates](#custom-templates)
7. [Available Commands](#available-commands)
8. [Workflow: Design Your Perfect Setup](#workflow-design-your-perfect-setup)
9. [Contributing](#contributing)
10. [License](#license)

---

## Features

🚀 **Project Scaffolding**: Generate complete project structures from YAML specifications
🎨 **Custom Templates**: Use built-in templates or create your own template collections
🔧 **Language Agnostic**: Support for Go, Python, Rust, and easily extensible to any language
📁 **Flexible Paths**: Works with relative paths, absolute paths, and `~` home directory expansion
✅ **Validation**: Built-in config validation to catch errors before generation
🏗️ **Clean Architecture**: Separation between structure creation and template generation

---

## Quick Start

```bash
# Install ForgeDir
go install github.com/KoHorizon/ForgeDir/cmd/fgdir@latest

# Create a project from the example config
fgdir init config.yaml

# View available templates
fgdir list-templates

# Validate your configuration before generating
fgdir validate config.yaml
```

---

## Installation

### Prerequisites
- Go 1.18 or newer
- Git (for cloning and releases)

### Install via `go install` (Recommended)
```bash
go install github.com/KoHorizon/ForgeDir/cmd/fgdir@latest
```

### Build from Source
```bash
# Clone the repository
git clone https://github.com/KoHorizon/ForgeDir.git
cd ForgeDir

# Build the binary
go build -o fgdir ./cmd/fgdir

# Add to PATH or use directly
./fgdir --help
```

### Verify Installation
```bash
fgdir version
fgdir help
```

---

## Usage Examples

### Basic Project Generation
```bash
# Generate project in current directory
fgdir init config.yaml

# Generate in specific directory
fgdir init config.yaml --output ./my-new-project

# Use custom templates
fgdir init config.yaml --templates ~/my-templates
```

### Working with Templates
```bash
# List all supported languages
fgdir list-templates

# View templates for specific language
fgdir list-templates go
fgdir list-templates python

# List templates from custom directory
fgdir list-templates --templates ~/my-templates
```

### Configuration Validation
```bash
# Validate configuration file
fgdir validate config.yaml
fgdir validate my-project-spec.yaml

# Validate with custom templates
fgdir validate config.yaml --templates ~/my-templates
```

---

## Configuration (`config.yaml`)

ForgeDir uses YAML files to define your project structure. Here's the basic format:

```yaml
# Basic project configuration
projectName: my_awesome_project
language: go

# Define your directory and file structure
structure:
  - type: dir
    name: cmd
    children:
      - type: file
        name: main.go

  - type: dir
    name: internal
    children:
      - type: dir
        name: handlers
        children:
          - type: file
            name: user.go
          - type: file
            name: auth.go
      - type: file
        name: config.go

  - type: dir
    name: pkg
    children:
      - type: file
        name: utils.go

  - type: file
    name: README.md
  - type: file
    name: go.mod
```

### Configuration Options

- **`projectName`**: Name of your project (used in templates)
- **`language`**: Target language (`go`, `python`, `rust`, or your custom language)
- **`structure`**: Array of directories and files to create

### Structure Node Types

- **`dir`**: Creates a directory (can contain `children`)
- **`file`**: Creates a file (populated with template content if available)

---

## Custom Templates

ForgeDir's real power comes from customizable templates. You can use built-in templates or create your own template collections.

### Understanding Template Structure

Templates are organized by language in directories:
```
templates/
├── go/
│   ├── main.go.tmpl
│   ├── handler.go.tmpl
│   ├── service.go.tmpl
│   └── (default).tmpl
├── python/
│   ├── __init__.py.tmpl
│   └── app.py.tmpl
└── rust/
    ├── main.rs.tmpl
    └── lib.rs.tmpl
```

### Creating Custom Templates

1. **Create your template directory structure:**
```bash
mkdir -p ~/my-templates/go
mkdir -p ~/my-templates/python
```

2. **Create template files** (example for Go):
```bash
# ~/my-templates/go/main.go.tmpl
cat > ~/my-templates/go/main.go.tmpl << 'EOF'
package {{ .DirName }}

import (
    "fmt"
    "log"
)

// {{ .FileName }} - Custom main entry point
func main() {
    fmt.Println("Welcome to {{ .Language }} project!")
    log.Println("Generated by ForgeDir with custom templates")
}
EOF
```

3. **Use your custom templates:**
```bash
fgdir init config.yaml --templates ~/my-templates
```

### Template Variables

Templates have access to these variables:
- **`{{ .Language }}`**: The configured language (`go`, `python`, etc.)
- **`{{ .DirName }}`**: Name of the directory containing the file
- **`{{ .FileName }}`**: Base filename without extension

### Template Matching Rules

1. **Exact match**: `main.go.tmpl` matches `main.go` files
2. **Fallback**: `(default).tmpl` is used when no specific template exists
3. **No template**: Empty files are created if no template is found

### Path Flexibility

ForgeDir supports various path formats for template directories:

```bash
# Relative paths
fgdir init config.yaml --templates ./templates
fgdir init config.yaml --templates ../shared-templates

# Home directory paths
fgdir init config.yaml --templates ~/my-templates
fgdir init config.yaml --templates ~/Documents/project-templates

# Absolute paths
fgdir init config.yaml --templates /usr/local/share/company-templates
```

---

## Available Commands

### `fgdir init`
Generate a project from YAML specification.
```bash
fgdir init [config.yaml] [flags]

Flags:
  -c, --config string     Path to YAML project spec (default "config.yaml")
  -o, --output string     Output directory (default ".")
  -t, --templates string  Custom templates directory
```

### `fgdir validate`
Validate configuration without generating files.
```bash
fgdir validate [config.yaml] [flags]

Flags:
  -t, --templates string  Custom templates directory
```

### `fgdir list-templates`
List available templates.
```bash
fgdir list-templates [language] [flags]

Flags:
  -t, --templates string  Custom templates directory

Examples:
  fgdir list-templates                    # List all languages
  fgdir list-templates go                 # List Go templates
  fgdir list-templates --templates ~/my-templates
```

### `fgdir version`
Show version information.
```bash
fgdir version
```

---

## Workflow: Design Your Perfect Setup

ForgeDir is designed around the idea that developers have preferred project layouts and coding patterns. Here's how to maximize its value:

### 1. **Analyze Your Preferred Structure**

Think about your typical project layout. For example, a Go microservice might always have:
- `cmd/` for main applications
- `internal/` for private code
- `pkg/` for public libraries
- `configs/` for configuration files
- Standard files like `Makefile`, `Dockerfile`, `.gitignore`

### 2. **Create Your Master Configuration**

Design a `config.yaml` that represents your ideal project structure:

```yaml
projectName: microservice_template
language: go
structure:
  - type: dir
    name: cmd
    children:
      - type: dir
        name: server
        children:
          - type: file
            name: main.go
  - type: dir
    name: internal
    children:
      - type: dir
        name: handlers
      - type: dir
        name: services
      - type: dir
        name: models
      - type: file
        name: config.go
  - type: dir
    name: pkg
    children:
      - type: dir
        name: api
  - type: dir
    name: configs
  - type: file
    name: Makefile
  - type: file
    name: Dockerfile
  - type: file
    name: .gitignore
  - type: file
    name: README.md
```

### 3. **Build Your Template Collection**

Create templates that match your coding style and standards:

```bash
mkdir -p ~/my-company-templates/go

# Create templates for your common patterns
echo 'package {{ .DirName }}

// {{ .FileName }} handles HTTP requests for {{ .DirName }}
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}' > ~/my-company-templates/go/handler.go.tmpl
```

### 4. **Test and Refine**

```bash
# Test your setup
fgdir validate my-template.yaml --templates ~/my-company-templates
fgdir init my-template.yaml --templates ~/my-company-templates --output ./test-project

# Iterate and improve
```

### 5. **Standardize Across Projects**

Once perfected, use the same configuration for all similar projects:

```bash
# For new microservices
fgdir init microservice-template.yaml --templates ~/company-templates --output ./new-service

# For CLI tools
fgdir init cli-template.yaml --templates ~/company-templates --output ./new-cli

# For web APIs
fgdir init api-template.yaml --templates ~/company-templates --output ./new-api
```

### Pro Tips

- **Version control your templates**: Keep your template directories in Git
- **Share with your team**: Use a shared templates repository
- **Multiple configurations**: Create different `config.yaml` files for different project types
- **Template inheritance**: Start with built-in templates and customize gradually

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository
2. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Make your changes and commit:
   ```bash
   git commit -m "feat: add my feature"
   ```
4. Push to your branch:
   ```bash
   git push origin feature/my-feature
   ```
5. Open a Pull Request

By contributing, you agree that your work will be licensed under the MIT License.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
