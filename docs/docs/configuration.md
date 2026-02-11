# Configuration

## Config Precedence

gsi uses [Viper](https://github.com/spf13/viper) for configuration management. Values are resolved in this order (highest priority first):

1. **CLI flags** -- `--module`, `--author`, etc.
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
| `skip-bmad` | `false` |
| `skip-git` | `false` |
| `skip-docs` | `false` |
| `only-docs` | `false` |
| `ui` | `false` |

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
LDFLAGS := -ldflags "-X $(MODULE)/cmd.version=$(VERSION) -X $(MODULE)/cmd.commit=$(COMMIT) -X $(MODULE)/cmd.date=$(BUILD_DATE)"
```

These correspond to variables in `cmd/version.go`:

```go
var (
    version = "dev"
    commit  = "unknown"
    date    = "unknown"
)
```

The version is derived from `git describe --tags --always --dirty`, so tag your releases (e.g., `git tag v0.1.0`) to get meaningful version strings.
