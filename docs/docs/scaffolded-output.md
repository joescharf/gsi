# What Gets Scaffolded

## Directory Structure

A default `gsi my-app` produces the following:

```
my-app/
├── _bmad/                  # BMAD agile method framework
├── cmd/
│   ├── config.go           # Config command (init/edit/check subcommands)
│   ├── config_init.go      # Viper config file discovery wiring
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
│   ├── config/
│   │   └── config.go        # Viper helpers (ConfigDir, SetDefaults, SaveConfig)
│   └── ui/
│       ├── dist/
│       │   └── index.html   # Placeholder UI (replaced by React build)
│       └── embed.go         # go:embed directive for static assets
├── .dockerignore            # Excludes non-essential files from Docker builds
├── .editorconfig            # Consistent editor settings
├── .github/
│   └── workflows/
│       └── release.yml      # GitHub Actions release workflow
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
| `cmd/config.go` | `config init`, `config edit`, `config check` subcommands |
| `cmd/config_init.go` | Wires `initConfig()` via `cobra.OnInitialize` for viper config file discovery |
| `internal/config/config.go` | `ConfigDir()`, `DefaultConfigFile()`, `SetDefaults()`, `SaveConfig()`, `InitViper()` helpers |
| `go.mod` / `go.sum` | Go module and dependency management |

### Build & Release

| File | Purpose |
|------|---------|
| `Makefile` | Targets for build, test, lint, release, docs, and UI |
| `.goreleaser.yml` | Goreleaser v2 config: binaries, archives, Docker, changelog |
| `.github/workflows/release.yml` | GitHub Actions workflow triggered on version tags |
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

## Capability-Gated Outputs

Each of the following can be toggled with `--<name>` / `--no-<name>` flags:

| Capability | Files | Default |
|------------|-------|---------|
| `bmad` | `_bmad/` | ON |
| `config` | `cmd/config.go`, `cmd/config_init.go`, `internal/config/config.go` | ON |
| `git` | `.git/`, `.gitignore`, initial commit | ON |
| `docs` | `docs/` and all contents | ON |
| `ui` | `ui/` (React/shadcn/Tailwind app) | OFF |
| `goreleaser` | `.goreleaser.yml` | ON |
| `docker` | `Dockerfile`, `.dockerignore` | ON |
| `release` | `.github/workflows/release.yml` | ON |
| `mockery` | `.mockery.yml` | ON |
| `editorconfig` | `.editorconfig` | ON |
| `makefile` | `Makefile` | ON |

### `--only-docs` Flag

Only generates the `docs/` directory and its contents. Useful for adding documentation to an existing project.

## Config Management (--config)

When enabled (default), gsi scaffolds a complete viper-based config system:

- **`cmd/config.go`** -- Three subcommands:
    - `config init` -- Creates config directory and writes default `config.yaml`
    - `config edit` -- Opens config file in `$EDITOR`
    - `config check` -- Displays current config values and sources
- **`cmd/config_init.go`** -- Registers `initConfig()` via `cobra.OnInitialize()` to set up viper config file discovery at startup
- **`internal/config/config.go`** -- Helper package with:
    - `ConfigDir()` -- OS-appropriate config directory via `os.UserConfigDir()`
    - `DefaultConfigFile()` -- Default config file path
    - `SetDefaults()` -- Centralized `viper.SetDefault()` calls
    - `SaveConfig()` -- Wrapper for `viper.WriteConfigAs()`
    - `InitViper()` -- Full viper setup (config file search, env vars, defaults)

## Idempotency

gsi is designed to be run multiple times safely:

- **Files that exist are skipped** -- `WriteTemplateFile` checks for file existence before writing
- **Directories that exist are reused** -- no error if the project dir already exists
- **Git repos with commits skip the initial commit** -- won't create duplicate commits
- **Dependencies already present are skipped** -- `cobra-cli`, `uv`, packages, etc.

This means you can run `gsi --only-docs .` on an existing project, or re-run `gsi .` to add missing files without overwriting customized ones.
