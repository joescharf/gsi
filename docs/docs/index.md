# gsi

**gsi** (Go Super Init) is a CLI tool that scaffolds production-ready Go projects with best practices and modern tooling baked in from the start.

Instead of manually setting up cobra, viper, makefiles, docker, goreleaser, docs, and testing infrastructure every time you start a new Go project, `gsi` generates all of it in seconds.

## Key Features

- **CLI framework** -- Cobra + Viper with version command and config management
- **Embedded web UI** -- Serve command with SPA routing and a placeholder UI ready for React/shadcn
- **Documentation** -- mkdocs-material site with live reload via `uv`
- **Release automation** -- Goreleaser config for binaries, Docker images, and Homebrew
- **Docker** -- Multi-platform Dockerfile and `.dockerignore`
- **Testing** -- Mockery config, coverage targets, race detection
- **Code quality** -- EditorConfig, golangci-lint, go vet, gofmt targets
- **Optional React UI** -- `--ui` flag scaffolds a React/shadcn/Tailwind frontend via bun
- **BMAD method** -- Optional agile framework scaffolding

## Prerequisites

- **Go 1.21+** -- [golang.org/dl](https://go.dev/dl/)
- **git** -- for version detection and initial commit
- **uv** (optional) -- for docs scaffolding ([astral.sh/uv](https://docs.astral.sh/uv/))
- **bun** (optional) -- for `--ui` React scaffolding ([bun.sh](https://bun.sh))
- **goreleaser** (optional) -- for release automation ([goreleaser.com](https://goreleaser.com))

## Quick Start

```bash
# Install
go install github.com/joescharf/gsi@latest

# Scaffold a new project
gsi my-awesome-app

# Build and run
cd my-awesome-app && make build && ./bin/my-awesome-app --help
```

## Quick Links

| Topic | Description |
|-------|-------------|
| [Getting Started](getting-started.md) | Installation methods and first project walkthrough |
| [CLI Reference](cli-reference.md) | All flags, subcommands, and usage examples |
| [What Gets Scaffolded](scaffolded-output.md) | Complete directory tree and file descriptions |
| [Configuration](configuration.md) | Viper config precedence and Makefile targets |
| [Releasing](releasing.md) | Goreleaser workflow for binaries, Docker, and Homebrew |
| [Contributing](contributing.md) | Development setup and how to add scaffold steps |
