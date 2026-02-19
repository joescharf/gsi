# CLI Reference

## Usage

```
gsi [flags] [project-name]
```

The project name argument is required. Use `.` to initialize in the current directory.

## Capability Flags

Every toggleable capability has a `--<name>` flag (to enable) and a hidden `--no-<name>` flag (to disable). The `--no-<name>` flag takes precedence if both are set.

| Flag | Default | Description |
|------|---------|-------------|
| `--bmad` / `--no-bmad` | ON | BMAD method framework installation |
| `--config` / `--no-config` | ON | Viper config management scaffolding |
| `--git` / `--no-git` | ON | Git initialization and initial commit |
| `--docs` / `--no-docs` | ON | mkdocs-material documentation scaffolding |
| `--ui` / `--no-ui` | OFF | React/shadcn/Tailwind UI in `ui/` subdirectory |
| `--goreleaser` / `--no-goreleaser` | ON | GoReleaser configuration |
| `--docker` / `--no-docker` | ON | Dockerfile and .dockerignore |
| `--release` / `--no-release` | ON | GitHub Actions release workflow |
| `--mockery` / `--no-mockery` | ON | Mockery configuration |
| `--editorconfig` / `--no-editorconfig` | ON | EditorConfig file |
| `--makefile` / `--no-makefile` | ON | Makefile with common targets |

Capabilities with missing soft dependencies are auto-disabled at runtime:

- `bmad` is auto-disabled if `npx` is not found
- `docs` is auto-disabled if `uv` is not found
- `git` is auto-disabled if `git` is not found

The `ui` capability has a **hard** dependency on `bun` -- gsi will error if `--ui` is set and `bun` is missing.

## Other Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--author` | `-a` | `"Joe Scharf joe@joescharf.com"` | Author name and email |
| `--module` | `-m` | `github.com/joescharf/<project>` | Go module path |
| `--dry-run` | `-d` | `false` | Show what would be done without executing |
| `--verbose` | `-v` | `false` | Enable verbose output |
| `--only-docs` | | `false` | Only add docs scaffolding (skip everything else) |

!!! note
    `--only-docs` and `--no-docs` are mutually exclusive.

## Subcommands

### `gsi version`

Print version, commit hash, and build date.

```bash
$ gsi version
gsi 0.1.0 (commit: abc1234, built: 2026-01-15T10:30:00Z)
```

### `gsi serve`

Start the embedded web UI server (available in scaffolded projects, not in gsi itself).

## Examples

```bash
# Basic scaffold
gsi my-app

# Full customization
gsi \
  --module github.com/acme/widget \
  --author "Jane Doe jane@acme.com" \
  --ui \
  --verbose \
  my-widget

# Preview without creating files
gsi --dry-run --verbose my-app

# Add docs to existing Go project
gsi --only-docs .

# Minimal scaffold (no BMAD, no docs, no git)
gsi --no-bmad --no-docs --no-git my-app

# Skip Docker and release workflow
gsi --no-docker --no-release my-app

# Skip config management scaffolding
gsi --no-config my-app

# Kitchen sink: disable most optional capabilities
gsi --no-bmad --no-docs --no-docker --no-release --no-mockery --no-editorconfig my-app
```
