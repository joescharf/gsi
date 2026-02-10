# go-superinit

A Zsh script that scaffolds new Go CLI projects with best practices and tooling.

## What it creates

- Go module with [Cobra](https://github.com/spf13/cobra) CLI + [Viper](https://github.com/spf13/viper) config
- `version` command
- [Mockery](https://github.com/vektra/mockery) configuration
- `.editorconfig`
- [BMAD method](https://github.com/bmad-method/bmad-method) framework (requires `bun`)
- [mkdocs-material](https://squidfunk.github.io/mkdocs-material/) documentation site in `docs/` (requires `uv`)
- Git repository with `.gitignore` and initial commit
- Optional: React/shadcn/Tailwind UI in `ui/` (requires `bun`)

## Prerequisites

- **Required:** `go`
- **Optional:** `git`, `bun` (for BMAD + UI), `uv` (for docs)

## Usage

```sh
./go-superinit.zsh [OPTIONS] [PROJECT_NAME]
```

### Options

| Flag | Description |
|------|-------------|
| `-a, --author TEXT` | Author name and email |
| `-m, --module PATH` | Go module path |
| `-d, --dry-run` | Show what would be done without executing |
| `-v, --verbose` | Enable verbose output |
| `--skip-bmad` | Skip BMAD method installation |
| `--skip-git` | Skip git initialization and commit |
| `--skip-docs` | Skip mkdocs-material documentation scaffolding |
| `--ui` | Initialize a React/shadcn/Tailwind UI in `ui/` |

### Examples

```sh
# Basic project
./go-superinit.zsh my-app

# Custom module path, dry-run
./go-superinit.zsh --module github.com/myorg/myapp --dry-run my-app

# Skip docs and BMAD
./go-superinit.zsh --skip-docs --skip-bmad my-app

# With UI
./go-superinit.zsh --ui my-app

# Initialize in current directory
./go-superinit.zsh .
```

## Docs scaffolding

By default, `go-superinit` creates a `docs/` directory with a [mkdocs-material](https://squidfunk.github.io/mkdocs-material/) documentation site managed by [uv](https://docs.astral.sh/uv/):

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

If `uv` is not installed, docs scaffolding is skipped automatically with a warning. Use `--skip-docs` to opt out explicitly.

## Idempotency

The script is idempotent -- it skips steps that have already been completed (e.g., existing `go.mod`, `cmd/`, `docs/`, `.git/`).
