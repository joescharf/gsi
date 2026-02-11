# What Gets Scaffolded

## Directory Structure

A default `gsi my-app` produces the following:

```
my-app/
├── _bmad/                  # BMAD agile method framework
├── cmd/
│   ├── root.go             # Cobra root command with viper config
│   ├── serve.go            # Embedded web UI server command
│   └── version.go          # Version subcommand (ldflags-injected)
├── docs/
│   ├── docs/
│   │   ├── stylesheets/
│   │   │   └── extra.css   # Custom mkdocs theme overrides
│   │   ├── getting-started.md
│   │   └── index.md
│   ├── .gitignore           # Ignores site/ and .venv/
│   ├── mkdocs.yml           # mkdocs-material configuration
│   └── pyproject.toml       # uv-managed Python dependencies
├── internal/
│   └── ui/
│       ├── dist/
│       │   └── index.html   # Placeholder UI (replaced by React build)
│       └── embed.go         # go:embed directive for static assets
├── .dockerignore            # Excludes non-essential files from Docker builds
├── .editorconfig            # Consistent editor settings
├── .gitignore               # Standard Go + docs + UI ignores
├── .goreleaser.yml          # Release automation config
├── .mockery.yml             # Mock generation config
├── Dockerfile               # Multi-platform Alpine image
├── Makefile                 # Build, test, lint, release, docs, UI targets
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── main.go                  # Entry point
```

## File Descriptions

### Core Application

| File | Purpose |
|------|---------|
| `main.go` | Entry point, calls `cmd.Execute()` |
| `cmd/root.go` | Root cobra command with viper flag bindings |
| `cmd/version.go` | Prints version/commit/date injected via ldflags |
| `cmd/serve.go` | Starts HTTP server serving the embedded UI |
| `go.mod` / `go.sum` | Go module and dependency management |

### Build & Release

| File | Purpose |
|------|---------|
| `Makefile` | Targets for build, test, lint, release, docs, and UI |
| `.goreleaser.yml` | Goreleaser v2 config: binaries, archives, Docker, changelog |
| `Dockerfile` | Alpine-based multi-platform image |
| `.dockerignore` | Keeps Docker context small |

### Embedded UI

| File | Purpose |
|------|---------|
| `internal/ui/embed.go` | `//go:embed all:dist` directive |
| `internal/ui/dist/index.html` | Default placeholder page |

### Documentation

| File | Purpose |
|------|---------|
| `docs/mkdocs.yml` | mkdocs-material site config |
| `docs/pyproject.toml` | Python deps managed by uv |
| `docs/docs/index.md` | Landing page |
| `docs/docs/getting-started.md` | Installation and usage guide |

### Code Quality

| File | Purpose |
|------|---------|
| `.editorconfig` | Indent style, charset, line endings |
| `.mockery.yml` | Mockery v2 interface mock config |
| `.gitignore` | Standard Go project ignores |

## Optional Outputs

### `--ui` Flag

When `--ui` is passed, gsi runs `bun init --react=shadcn ui` which creates a full React/shadcn/Tailwind app in `ui/`. The Makefile includes `ui-dev`, `ui-build`, `ui-embed`, and `ui-deps` targets.

### `--skip-bmad` Flag

Omits the `_bmad/` directory.

### `--skip-docs` Flag

Omits the `docs/` directory and related Makefile targets still exist but are guarded.

### `--only-docs` Flag

Only generates the `docs/` directory and its contents. Useful for adding documentation to an existing project.

## Idempotency

gsi is designed to be run multiple times safely:

- **Files that exist are skipped** -- `WriteTemplateFile` checks for file existence before writing
- **Directories that exist are reused** -- no error if the project dir already exists
- **Git repos with commits skip the initial commit** -- won't create duplicate commits
- **Dependencies already present are skipped** -- `cobra-cli`, `uv`, packages, etc.

This means you can run `gsi --only-docs .` on an existing project, or re-run `gsi .` to add missing files without overwriting customized ones.
