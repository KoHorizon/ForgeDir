# ForgeDir

A command-line tool written in Go to generate project directory and file structures from YAML configuration files.

## ✨ Overview

Setting up the initial directory structure for a new project can be a repetitive task. ForgeDir automates this process by reading a simple YAML file that defines your desired folder and file layout, and then instantly creates that structure on your filesystem.

ForgeDir is currently under active development and serves as a practical project for learning Go, focusing initially on the core functionality of structure creation.

## 🚀 Features

* Define complex nested directory and file structures using a clean YAML format.
* Automatically create directories.
* Automatically create empty files.
* Fast and efficient structure generation using Go's standard library.
* *(Planned)* Language-agnostic templates and boilerplate generation.

## ⬇️ Getting Started

### Prerequisites

* Go (version 1.18 or higher recommended)

### Installation

You can install ForgeDir using `go install`:

```bash
go install [github.com/YourGitHubUsername/forgedir@latest](https://github.com/YourGitHubUsername/forgedir@latest) # Replace YourGitHubUsername
