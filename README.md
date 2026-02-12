# gsi

A Go CLI tool that scaffolds new Go projects with best practices and tooling.

## What it creates

- Go module with [Cobra](https://github.com/spf13/cobra) CLI + [Viper](https://github.com/spf13/viper) config
- `version` command with ldflags build vars
- `config` command with `init`, `edit`, `check` subcommands + `internal/config/` package
- Viper config file discovery, env var support, and `SetDefaults()` wiring
- [Mockery](https://github.com/vektra/mockery) configuration
- `.editorconfig`
- [BMAD method](https://github.com/bmad-method/bmad-method) framework (requires `bun`)
- [mkdocs-material](https://squidfunk.github.io/mkdocs-material/) documentation site in `docs/` (requires `uv`)
- [Goreleaser](https://goreleaser.com) config for binaries, Docker images, and Homebrew
- Dockerfile and `.dockerignore`
- GitHub Actions release workflow
- Makefile with build, test, lint, release, docs, and UI targets
- Git repository with `.gitignore` and initial commit
- Optional: React/shadcn/Tailwind UI in `ui/` (requires `bun`)

## Prerequisites

- **Required:** `go`
- **Optional:** `git`, `bun` (for BMAD + UI), `uv` (for docs)

## Installation

```sh
go install github.com/joescharf/gsi@latest
```

## Usage

```sh
gsi [OPTIONS] [PROJECT_NAME]
```

### Capability Flags

Every toggleable capability has a `--<name>` flag (to enable) and a hidden `--no-<name>` flag (to disable). Most capabilities default to ON; `ui` defaults to OFF.

| Flag | Default | Description |
|------|---------|-------------|
| `--bmad` / `--no-bmad` | ON | BMAD method framework installation |
| `--config` / `--no-config` | ON | Viper config management scaffolding |
| `--git` / `--no-git` | ON | Git initialization and initial commit |
| `--docs` / `--no-docs` | ON | mkdocs-material documentation scaffolding |
| `--ui` / `--no-ui` | OFF | React/shadcn/Tailwind UI in `ui/` |
| `--goreleaser` / `--no-goreleaser` | ON | GoReleaser configuration |
| `--docker` / `--no-docker` | ON | Dockerfile and .dockerignore |
| `--release` / `--no-release` | ON | GitHub Actions release workflow |
| `--mockery` / `--no-mockery` | ON | Mockery configuration |
| `--editorconfig` / `--no-editorconfig` | ON | EditorConfig file |
| `--makefile` / `--no-makefile` | ON | Makefile with common targets |

### Other Flags

| Flag | Description |
|------|-------------|
| `-a, --author TEXT` | Author name and email |
| `-m, --module PATH` | Go module path |
| `-d, --dry-run` | Show what would be done without executing |
| `-v, --verbose` | Enable verbose output |
| `--only-docs` | Only add docs scaffolding (skip everything else) |

### Examples

```sh
# Basic project
gsi my-app

# Custom module path, dry-run
gsi --module github.com/myorg/myapp --dry-run my-app

# Skip docs and BMAD
gsi --no-docs --no-bmad my-app

# Minimal: no docker, no release, no goreleaser
gsi --no-docker --no-release --no-goreleaser my-app

# With UI
gsi --ui my-app

# Without config scaffolding
gsi --no-config my-app

# Initialize in current directory
gsi .
```

## Docs scaffolding

By default, `gsi` creates a `docs/` directory with a [mkdocs-material](https://squidfunk.github.io/mkdocs-material/) documentation site managed by [uv](https://docs.astral.sh/uv/):

```
docs/
  mkdocs.yml           # Material theme, dark/light mode, mermaid, admonitions
  pyproject.toml       # uv-managed Python project
  uv.lock
  .gitignore
  docs/
    index.md           # Welcome page
    getting-started.md # Starter page with install/usage
    stylesheets/
      extra.css        # Compact nav + smaller code font
```

To serve the docs locally:

```sh
cd docs/ && uv run mkdocs serve
```

If `uv` is not installed, docs scaffolding is skipped automatically with a warning. Use `--no-docs` to opt out explicitly.

## Idempotency

The tool is idempotent -- it skips steps that have already been completed (e.g., existing `go.mod`, `cmd/`, `docs/`, `.git/`).
