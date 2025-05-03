# ForgeDir

<!-- Badges -->

[![Go Report Card](https://goreportcard.com/badge/github.com/KoHorizon/ForgeDir)](https://goreportcard.com/report/github.com/KoHorizon/ForgeDir)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A flexible, YAML-driven CLI tool written in Go to automate the creation of project directory and file structures.

---

## ✨ Overview

Initializing a new project often involves manually creating nested folders and placeholder files. **ForgeDir** streamlines this setup by reading a simple YAML spec and instantly scaffolding the desired layout.

This tool is under active development and doubles as a practical learning project for Go developers focusing on building CLI apps and implementing clean architecture principles.

---

## 🚀 Features

* **YAML-driven**: Define any nested directory/file hierarchy in a human-readable YAML format.
* **Fast & Lightweight**: Leverages Go's standard library—no external dependencies for the core scaffold.
* **Extensible**: Built with a plugin-ready generator system for future language-agnostic boilerplate.
* **Cross-platform**: Works on Linux, macOS, and Windows.

---

## 🔧 Getting Started

ForgeDir aims to simplify scaffolding new projects by interpreting a YAML configuration and generating the necessary directory and file structure automatically. Installation and detailed usage instructions will be provided in a future update.

---

## 🛠️ Development

Directory structure follows Cobra conventions:

```
cmd/       # CLI entrypoints
internal/  # Core builder, config loader, and generator logic
```

Use `go run main.go` for rapid iteration without rebuilding.

---

## 🤝 Contributing

Contributing to ForgeDir is easy and encouraged under the project's MIT License. You don't need special permissions—just follow these steps:

1. **Fork the repository.**
2. **Create a feature branch:** `git checkout -b feature/my-feature`
3. **Make your changes.**
4. **Commit your changes:** `git commit -m "feat: add my feature"`
5. **Push to your branch:** `git push origin feature/my-feature`
6. **Open a Pull Request** against the `main` branch.

By contributing, you agree that your work will be licensed under the MIT License. No additional Contributor License Agreement (CLA) is required unless otherwise stated.

---

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
