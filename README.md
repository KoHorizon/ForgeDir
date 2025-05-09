# ForgeDir

[![Go Report Card](https://goreportcard.com/badge/github.com/KoHorizon/ForgeDir)](https://goreportcard.com/report/github.com/KoHorizon/ForgeDir)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ForgeDir is a CLI tool written in Go that scaffolds a project structure from a simple YAML specification, generating both folders/files and language-specific boilerplate templates. It uses Cobra for the CLI interface and follows clean-architecture principles to separate concerns between configuration loading, filesystem operations, and template generation.

---

## Table of Contents

1. [Features](#features)
2. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
3. [Usage](#usage)
4. [Specification (`spec.yaml`)](#specification-specyaml)
5. [Contributing](#contributing)
6. [License](#license)

---

## Features

* Scaffold project directories & empty files from a YAML spec
* Generate language-specific boilerplate via Go `embed` and `text/template`
* Built with Go, using Cobra for an idiomatic CLI
* Extensible: add new language templates by dropping `*.tmpl` files
* Supports user overrides of templates via a local or global `templates/` directory

---

## Getting Started

ForgeDir aims to simplify scaffolding new projects by interpreting a YAML configuration and generating the necessary directory and file structure automatically.

### Prerequisites

- Go 1.18 or newer installed
- Git (for cloning and tagging releases)

### Installation

#### Rapid iteration with `go run`

```bash
# From your project root:
go run ./cmd/fgdir --help
go run ./cmd/fgdir init config.yaml
```

#### Build a standalone binary

```bash
# From your project root:
go build -o fgdir ./cmd/fgdir

# Run it locally:
./fgdir help
./fgdir init config.yaml
```

#### Install via go install

For Go 1.17+ users, simply run:

```bash
go install github.com/KoHorizon/ForgeDir/cmd/fgdir@latest
```

This will drop the `fgdir` binary into:

```
$(go env GOPATH)/bin
```

Make sure that directory is on your PATH:

```bash
# bash / zsh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Now you can invoke `fgdir` directly:

```bash
fgdir version
fgdir help
fgdir init config.yaml
```

---

## Usage

```bash
# Show help:
fgdir help

# Scaffold a project from your YAML spec:
fgdir init [flags] <spec.yaml>
```

Common flags for `init`:

- `-c, --config <path>` path to the YAML project spec (default: `config.yaml`)
- `-o, --output <path>` directory where the project will be generated (default: current directory)

---

## Specification (`spec.yaml`)

Define your project spec in YAML, for example:

```yaml
# config.yaml
projectName: my_project
language: go
structure:
  - type: dir # Define a directory
    name: cmd
    children:
      - type: file
        name: main.go
  - type: dir
    name: internal
    children:
      - type: file
        name: core.go
  - type: dir
    name: pkg
    children:
      - type: file
        name: handler.go
      - type: dir
        name: api
        children:
          - type: file
            name: handlers.go
```

Use `fgdir init --help` to see available flags.

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
5. Open a Pull Request against the main (or dev) branch

By contributing, you agree that your work will be licensed under the MIT License.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
