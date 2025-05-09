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
4. [Specification (`spec.yaml`)](#specification-specyaml)
5. [License](#license)

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

#### Rapid iteration with `go run`

```bash
# From your project root:
go run main.go --help
go run main.go init config.yaml
```

#### Build a standalone binary

```bash
# Build locally:
go build -o fgdir main.go

# Run:
./fgdir help
./fgdir init config.yaml
```

#### Install into your \$GOBIN

```bash
# Assuming module path github.com/yourname/forgedir:
go install github.com/yourname/forgedir@latest

# Now you can run:
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

* `-c, --config <path>`   path to the YAML project spec (default: `config.yaml`)
* `-o, --output <path>`   directory where the project will be generated (default: current directory)

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
      - type: file # Define a file
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

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
