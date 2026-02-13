# What Gets Scaffolded

## Directory Structure

A default `gsi my-app` produces the following:

```
my-app/
├── _bmad/                  # BMAD agile method framework
├── cmd/
│   ├── config.go           # Config command (init/edit/check subcommands)
│   ├── config_init.go      # Viper config file discovery wiring
│   ├── root.go             # Cobra root command with Execute(version, commit, date)
│   ├── serve.go            # Embedded web UI server command
│   └── version.go          # Version subcommand (uses buildVersion from root)
├── docs/
│   ├── docs/
│   │   ├── stylesheets/
│   │   │   └── extra.css   # Custom mkdocs theme overrides
│   │   ├── getting-started.md
│   │   └── index.md
│   ├── scripts/
│   │   ├── add_browser_frame.py  # Adds macOS browser frames to screenshots
│   │   ├── scrape.sh             # Screenshot capture pipeline
│   │   └── shots.yaml            # shot-scraper page configuration
│   ├── .gitignore           # Ignores site/, .venv/, scripts/img*
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
│       ├── ci.yml           # CI workflow (test + lint)
│       ├── docs.yml         # Docs deployment to GitHub Pages
│       └── release.yml      # Release workflow (manual dispatch)
├── .gitignore               # Standard Go + docs + UI ignores
├── .goreleaser.yml          # Release automation (3-platform, Docker, Homebrew)
├── .mockery.yml             # Mock generation config
├── Dockerfile               # Multi-platform Alpine 3.21 image (non-root user)
├── Makefile                 # Build, test, lint, release, release-local, docs, UI targets
├── my-app_pycodesign.ini    # macOS code signing config template
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── main.go                  # Entry point with version vars, calls cmd.Execute()
```

With `--ui`:

```
ui/
├── build.ts               # Bun build script with publicPath: "/" (SPA fix)
├── package.json           # Build script set to "bun run build.ts"
├── src/
│   └── ...                # React/shadcn/Tailwind app
└── ...
```

## File Descriptions

### Core Application

| File | Purpose |
|------|---------|
| `main.go` | Entry point with `version`, `commit`, `date` vars; calls `cmd.Execute(version, commit, date)` |
| `cmd/root.go` | Root cobra command with `Execute(version, commit, date string)` that stores build metadata |
| `cmd/version.go` | Prints version/commit/date from `buildVersion`/`buildCommit`/`buildDate` set by Execute |
| `cmd/serve.go` | Starts HTTP server serving the embedded UI |
| `cmd/config.go` | `config init`, `config edit`, `config check` subcommands |
| `cmd/config_init.go` | Wires `initConfig()` via `cobra.OnInitialize` for viper config file discovery |
| `internal/config/config.go` | `ConfigDir()`, `DefaultConfigFile()`, `SetDefaults()`, `SaveConfig()`, `InitViper()` helpers |
| `go.mod` / `go.sum` | Go module and dependency management |

### Build & Release

| File | Purpose |
|------|---------|
| `Makefile` | Targets for build, test, lint, release, release-local, docs, and UI |
| `.goreleaser.yml` | Goreleaser v2: 3-platform builds (Linux/macOS/Windows), archives, Docker, Homebrew, changelog |
| `.github/workflows/release.yml` | Manual dispatch release with QEMU, Buildx, GHCR login, GoReleaser |
| `.github/workflows/ci.yml` | Push/PR CI: test + lint jobs with bun UI embed |
| `.github/workflows/docs.yml` | GitHub Pages deployment for mkdocs-material docs |
| `Dockerfile` | Alpine 3.21, non-root user, tzdata, TARGETPLATFORM, env var for DB path |
| `.dockerignore` | Keeps Docker context small |
| `<project>_pycodesign.ini` | macOS code signing config template (Developer ID certs) |

### Embedded UI

| File | Purpose |
|------|---------|
| `internal/ui/embed.go` | `//go:embed all:dist` directive |
| `internal/ui/dist/index.html` | Default placeholder page |
| `ui/build.ts` | Bun build script with `publicPath: "/"` for SPA routing (when `--ui`) |

### Documentation

| File | Purpose |
|------|---------|
| `docs/mkdocs.yml` | mkdocs-material config with `site_url`, `repo_url`, `repo_name`, `edit_uri` |
| `docs/pyproject.toml` | Python deps managed by uv |
| `docs/docs/index.md` | Landing page |
| `docs/docs/getting-started.md` | Installation and usage guide |
| `docs/scripts/scrape.sh` | Screenshot capture pipeline (health check, shot-scraper, browser frames, compress) |
| `docs/scripts/shots.yaml` | shot-scraper page configuration (starter with dashboard + detail page) |
| `docs/scripts/add_browser_frame.py` | Adds macOS-style browser frames to screenshots using Pillow |

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
| `docs` | `docs/` and all contents, `.github/workflows/docs.yml` | ON |
| `ui` | `ui/` (React/shadcn/Tailwind app with `build.ts`) | OFF |
| `goreleaser` | `.goreleaser.yml` | ON |
| `docker` | `Dockerfile`, `.dockerignore` | ON |
| `release` | `.github/workflows/release.yml`, `.github/workflows/ci.yml`, `<project>_pycodesign.ini` | ON |
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
- **main.go and cmd/root.go are always overwritten** -- these are regenerated from templates to ensure the `main.*` ldflags pattern

This means you can run `gsi --only-docs .` on an existing project, or re-run `gsi .` to add missing files without overwriting customized ones.
