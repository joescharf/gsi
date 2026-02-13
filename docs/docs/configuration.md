# Configuration

## gsi's Own Config

gsi uses [Viper](https://github.com/spf13/viper) for its own configuration management. Values are resolved in this order (highest priority first):

1. **CLI flags** -- `--module`, `--author`, `--no-docker`, etc.
2. **Environment variables** -- (if configured)
3. **Config file** -- (if configured)
4. **Defaults** -- built-in default values

### Defaults

| Setting | Default Value |
|---------|---------------|
| `author` | `Joe Scharf joe@joescharf.com` |
| `module` | `github.com/joescharf/<project-name>` |
| `dry-run` | `false` |
| `verbose` | `false` |
| `only-docs` | `false` |

### Capability Defaults

| Capability | Default |
|------------|---------|
| `bmad` | ON |
| `config` | ON |
| `git` | ON |
| `docs` | ON |
| `ui` | OFF |
| `goreleaser` | ON |
| `docker` | ON |
| `release` | ON |
| `mockery` | ON |
| `editorconfig` | ON |
| `makefile` | ON |

Each capability can be toggled with `--<name>` (enable) or `--no-<name>` (disable). See the [CLI Reference](cli-reference.md) for details.

## Scaffolded Config Management

When the `config` capability is enabled (default), gsi scaffolds a complete viper-based configuration system into the generated project.

### What Gets Generated

- **`internal/config/config.go`** -- Core config helpers
- **`cmd/config.go`** -- `config init`, `config edit`, `config check` subcommands
- **`cmd/config_init.go`** -- `initConfig()` function wired via `cobra.OnInitialize()`

### Config File Location

The scaffolded project uses `os.UserConfigDir()` for OS-appropriate config paths:

| OS | Config Directory |
|----|-----------------|
| Linux/BSD | `$XDG_CONFIG_HOME/<project>/` or `~/.config/<project>/` |
| macOS | `~/Library/Application Support/<project>/` |
| Windows | `%AppData%/<project>/` |

### Config Precedence (Scaffolded Project)

The scaffolded project follows standard viper precedence:

1. **CLI flags** -- highest priority
2. **Environment variables** -- prefixed with `<PROJECT>_` (hyphens replaced with underscores), dot-separated keys use `_` (e.g., `MYAPP_SERVER_PORT`)
3. **Config file** -- `config.yaml` in config dir or current directory
4. **Defaults** -- from `config.SetDefaults()`

### Config Subcommands

```bash
# Create config directory and write default config.yaml
my-app config init

# Open config file in $EDITOR
my-app config edit

# Display current config values and their sources
my-app config check
```

### Customizing Defaults

Edit `internal/config/config.go` to add your application's defaults:

```go
func SetDefaults() {
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("log.level", "info")
    // Add your defaults here
    viper.SetDefault("db.path", "myapp.db")
}
```

## Generated Makefile Targets

The scaffolded `Makefile` provides these target groups:

### App

| Target | Description |
|--------|-------------|
| `make build` | Build the Go binary to `bin/` |
| `make install` | Install binary to `$GOPATH/bin` |
| `make run` | Build and run |
| `make serve` | Start embedded web UI server |
| `make clean` | Remove build artifacts |
| `make tidy` | Run `go mod tidy` |
| `make test` | Run tests with race detection |
| `make test-cover` | Run tests with coverage report |
| `make lint` | Run golangci-lint (requires install) |
| `make vet` | Run `go vet` |
| `make fmt` | Run `gofmt` |
| `make mocks` | Generate mocks with mockery |

### Release

| Target | Description |
|--------|-------------|
| `make release` | Create a release with goreleaser |
| `make release-local` | Create a signed local release (macOS code-signing) |
| `make release-snapshot` | Create a snapshot release (no publish) |

### Docs

| Target | Description |
|--------|-------------|
| `make docs-serve` | Start mkdocs dev server |
| `make docs-build` | Build static docs site |
| `make docs-deps` | Install doc dependencies via uv |

### UI (when `--ui` is used)

| Target | Description |
|--------|-------------|
| `make ui-dev` | Start React dev server |
| `make ui-build` | Build UI for production |
| `make ui-embed` | Copy built UI into `internal/ui/dist/` |
| `make ui-deps` | Install UI dependencies |

### Aggregate

| Target | Description |
|--------|-------------|
| `make all` | Build all existing artifacts |
| `make deps` | Install all dependencies |
| `make dev` | Start all dev servers in parallel |
| `make help` | Show available targets |

## Ldflags

The Makefile injects version information at build time via ldflags:

```makefile
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"
```

These target variables in `main.go`, which passes them to `cmd.Execute(version, commit, date)`:

```go
// main.go
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    cmd.Execute(version, commit, date)
}
```

```go
// cmd/root.go
func Execute(version, commit, date string) {
    buildVersion = version
    buildCommit = commit
    buildDate = date
    // ...
}
```

The `-s -w` flags strip debug info and DWARF tables for smaller binaries. The version is derived from `git describe --tags --always --dirty`, so tag your releases (e.g., `git tag v0.1.0`) to get meaningful version strings.
