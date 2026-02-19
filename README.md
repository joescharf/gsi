# gsi

A Go CLI tool that scaffolds new Go projects with best practices and tooling.

## What it creates

- Go module with [Cobra](https://github.com/spf13/cobra) CLI + [Viper](https://github.com/spf13/viper) config
- `main.go` with version vars passed to `cmd.Execute(version, commit, date)`
- `version` command using build vars from `main.*` ldflags
- `config` command with `init`, `edit`, `check` subcommands + `internal/config/` package
- Viper config file discovery, env var support, and `SetDefaults()` wiring
- [Mockery](https://github.com/vektra/mockery) configuration
- `.editorconfig`
- [BMAD method](https://github.com/bmad-method/bmad-method) framework (requires `npx`/Node.js)
- [mkdocs-material](https://squidfunk.github.io/mkdocs-material/) documentation site in `docs/` (requires `uv`)
- [Goreleaser](https://goreleaser.com) config for Linux, macOS, and Windows binaries, Docker images, and Homebrew
- Dockerfile (Alpine 3.21, non-root user, multi-platform) and `.dockerignore`
- GitHub Actions workflows: release (manual dispatch), CI (test + lint), docs (GitHub Pages)
- macOS universal binary with optional code signing via pycodesign
- Makefile with build, test, lint, release, release-local, docs, and UI targets
- Git repository with `.gitignore` and initial commit
- GitHub Pages configuration with setup instructions
- Optional: React/shadcn/Tailwind UI in `ui/` with `publicPath: "/"` SPA fix (requires `bun`)

## Prerequisites

- **Required:** `go`
- **Optional:** `git`, `npx`/Node.js (for BMAD), `bun` (for UI), `uv` (for docs), `gh` (for GitHub Pages setup)

## Installation

```sh
go install github.com/joescharf/gsi@latest
```

Or via Homebrew:

```sh
brew install joescharf/tap/gsi
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

# With UI (includes SPA routing fix via build.ts)
gsi --ui my-app

# Without config scaffolding
gsi --no-config my-app

# Initialize in current directory
gsi .
```

## Release infrastructure

Scaffolded projects get a complete release pipeline:

- **3-platform builds:** Linux, macOS (universal binary), Windows -- all amd64 + arm64
- **ldflags:** `-s -w -X main.version=... -X main.commit=... -X main.date=...`
- **CI workflow:** test + lint jobs with Go, Bun, and embedded UI
- **Docs workflow:** GitHub Pages deployment via mkdocs-material
- **Release workflow:** Manual dispatch with QEMU, Docker Buildx, GHCR login
- **Docker images:** Multi-arch via `dockers_v2` pushed to GHCR
- **Homebrew:** Formula distribution (or Cask for signed apps)
- **Code signing:** pycodesign config template for macOS notarization

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
  scripts/
    scrape.sh              # Screenshot capture pipeline
    shots.yaml             # shot-scraper page configuration
    add_browser_frame.py   # Adds macOS browser frames to screenshots
```

The mkdocs config includes `site_url`, `repo_url`, `repo_name`, and `edit_uri` pointed at your GitHub repo for "Edit on GitHub" links.

To serve the docs locally:

```sh
cd docs/ && uv run mkdocs serve
```

### Screenshot pipeline

The `docs/scripts/` directory contains a screenshot capture pipeline for documentation images. It uses [shot-scraper](https://shot-scraper.datasette.io/) to capture pages from a running local server, adds macOS-style browser frames via a Python script, and compresses the results with [imageoptim-cli](https://github.com/JamieMason/ImageOptim-CLI):

```sh
cd docs/scripts/ && bash scrape.sh
```

Edit `shots.yaml` to configure which pages to capture. The pipeline outputs to `docs/docs/img/`.

If `uv` is not installed, docs scaffolding is skipped automatically with a warning. Use `--no-docs` to opt out explicitly.

## Idempotency

The tool is idempotent -- it skips steps that have already been completed (e.g., existing `go.mod`, `cmd/`, `docs/`, `.git/`). The `main.go` and `cmd/root.go` files are overwritten after cobra-cli init to switch to the `main.*` ldflags pattern.
