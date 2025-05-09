# ForgeDir

[![Go Report Card](https://goreportcard.com/badge/github.com/KoHorizon/ForgeDir)](https://goreportcard.com/report/github.com/KoHorizon/ForgeDir)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ForgeDir is a CLI tool written in Go that scaffolds a project structure from a simple YAML specification, generating both folders/files and language‑specific boilerplate templates. It uses Cobra for the CLI interface and follows clean‑architecture principles to separate concerns between configuration loading, filesystem operations, and template generation.

---

## Table of Contents

1. [Features](#features)
2. [Getting Started](#getting-started)
   * [Prerequisites](#prerequisites)
   * [Installation](#installation)
3. [Usage](#usage)
   * [Commands](#commands)
4. [Specification (`spec.yaml`)](#specification-specyaml)
5. [Custom Templates](#custom-templates)
6. [Publishing & Releases](#publishing--releases)
7. [Versioning](#versioning)
8. [Testing](#testing)
9. [Contributing](#contributing)
10. [License](#license)

---

## Features

* Scaffold project directories & empty files from a YAML spec
* Generate language‑specific boilerplate via Go `embed` and `text/template`
* Built with Go, using Cobra for an idiomatic CLI
* Extensible: add new language templates by dropping `*.tmpl` files
* Supports user overrides of templates via a local or global `templates/` directory

---

## Getting Started

### Prerequisites

* Go 1.18 or newer installed
* Git (for cloning and tagging releases)

### Installation

Choose one of the following methods:

#### 1. Rapid iteration with `go run`

```bash
# From your project root (with main.go):
go run main.go --help
go run main.go init config.yaml
```

#### 2. Build a standalone binary

```bash
# Build locally:
go build -o fgdir main.go

# Run:
./fgdir --help
./fgdir init config.yaml
```

#### 3. Install into your \$GOBIN

```bash
# Assuming module path github.com/yourname/forgedir:
go install github.com/yourname/forgedir@latest

# Now you can run:
fgdir --help
fgdir init config.yaml
```

---

## Usage

Once installed, `fgdir` supports the following commands:

```bash
# Show help:
fgdir --help

# Scaffold a project from your YAML spec:
fgdir init [flags] --templates <path> <spec.yaml>

# Validate a spec without scaffolding:
fgdir validate <spec.yaml>

# List available templates (built-in or for a specific language):
fgdir list-templates [--lang go|js|py]

# Clean up a generated project:
fgdir clean <target-dir>

# Print version:
fgdir version
```

Common global flags:

* `-c, --config <path>`
* `-o, --output <path>`

Example:

```bash
fgdir init -c myspec.yaml -o ./outdir
```

---

## Specification (`spec.yaml`)

Your project spec is a YAML file, for example:

```yaml
project:
  name: "my-service"
  language: "go"

structures:
  - path: "cmd/myservice"
  - path: "internal/server"
  - path: "templates"

files:
  - path: "main.go"
  - path: "Dockerfile"
```

Use `fgdir validate spec.yaml` to check for schema errors before scaffolding.

---

## Custom Templates

ForgeDir ships with built‑in templates in `templates/<lang>/` embedded into the binary. To override or extend them:

1. Create a local `./templates/<lang>/` directory beside your spec.
2. Copy any default `*.tmpl` files you wish to customize into that folder.
3. Run `fgdir init --templates ./templates spec.yaml`.

The lookup order is:

1. Files in your `--templates` folder (OSFS overrides)
2. Built‑in embedded templates (go\:embed defaults)

Missing files still fall back to defaults, so you only need to override the ones you care about.

---

## Publishing & Releases

1. Push your code to a public repo (e.g. GitHub) with module path in `go.mod`:

   ```go
   module github.com/yourname/forgedir
   ```

2. Tag a release:

   ```bash
   ```

git tag v0.1.0
git push origin v0.1.0

````

3. Instruct users to install via:

   ```bash
go install github.com/yourname/forgedir@v0.1.0
````

4. (Optional) Use [goreleaser](https://goreleaser.com/) in CI to publish pre‑built binaries for Linux/Mac/Windows on GitHub Releases.
5. (Optional) Provide Homebrew/Scoop/Apt formulas so users can install via package managers.

---

## Versioning

We follow [Semantic Versioning](https://semver.org/):

* **v0.x**: initial development, breaking changes allowed
* **v1.0.0**: stable release, backwards compatibility guaranteed

Tag incremental releases (v0.2.0, v0.3.0, …) as you add features and tests. Reserve `v1.0.0` for an API‑stable milestone.

---

## Testing

We recommend writing unit tests for:

* `config` package (loading & validation of specs)
* `builder` package (folder/file creation logic)
* `generator` package (template rendering)
* CLI commands (using Cobra’s `ExecuteC` in tests)

Run tests with:

```bash
go test ./... -cover
```

Aim for high coverage on core packages before tagging each release.

---

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/your-feature`)
3. Write tests and code
4. Ensure all tests pass
5. Open a Pull Request with a clear description of changes

Please adhere to Go idioms and existing code style.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
