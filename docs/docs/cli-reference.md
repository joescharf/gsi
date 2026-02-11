# CLI Reference

## Usage

```
gsi [flags] [project-name]
```

The project name argument is required. Use `.` to initialize in the current directory.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--author` | `-a` | `"Joe Scharf joe@joescharf.com"` | Author name and email |
| `--module` | `-m` | `github.com/joescharf/<project>` | Go module path |
| `--dry-run` | `-d` | `false` | Show what would be done without executing |
| `--verbose` | `-v` | `false` | Enable verbose output |
| `--ui` | | `false` | Initialize a React/shadcn/Tailwind UI in `ui/` |
| `--skip-bmad` | | `false` | Skip BMAD method installation |
| `--skip-git` | | `false` | Skip git initialization and commit |
| `--skip-docs` | | `false` | Skip mkdocs-material documentation scaffolding |
| `--only-docs` | | `false` | Only add docs scaffolding (skip everything else) |

!!! note
    `--only-docs` and `--skip-docs` are mutually exclusive.

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
gsi --skip-bmad --skip-docs --skip-git my-app
```
