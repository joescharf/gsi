# Contributing

## Development Setup

```bash
# Clone the repository
git clone https://github.com/joescharf/gsi.git
cd gsi

# Build
make build

# Run tests
make test

# Run linter
make lint
```

## Project Structure

```
gsi/
├── cmd/                        # Cobra commands
│   ├── root.go                 # Root command, capability flags, viper bindings
│   ├── serve.go                # Serve command for embedded UI
│   └── version.go              # Version command (ldflags-injected)
├── internal/
│   ├── logger/                 # Structured logger with color output
│   ├── scaffold/
│   │   ├── scaffold.go         # Main orchestrator (step sequencing)
│   │   ├── steps.go            # Individual scaffold step methods
│   │   ├── steps_test.go       # Step tests (idempotency, dry-run, guards)
│   │   ├── config.go           # Config struct, capability constants, IsEnabled/Disable
│   │   ├── config_test.go      # Config unit tests (DefaultCapabilities, IsEnabled, Disable)
│   │   ├── environment.go      # Environment validation (capability-aware)
│   │   └── files.go            # WriteTemplateFile / WriteStaticFile helpers
│   ├── templates/
│   │   ├── templates.go        # Template rendering engine
│   │   ├── templates_test.go   # Template render tests
│   │   └── files/              # Embedded template files (*.tmpl)
│   └── ui/                     # Embedded UI assets
├── docs/                       # mkdocs-material documentation
├── .goreleaser.yml             # Release configuration
├── Dockerfile                  # Docker image definition
├── Makefile                    # Build targets
├── main.go                     # Entry point
└── go.mod                      # Module definition
```

## How to Add a New Scaffold Step

Adding a new file to the scaffold output follows a pattern:

### 1. Create the Template

Add a new `.tmpl` file in `internal/templates/files/`:

```
internal/templates/files/my_config.tmpl
```

Available template variables (from `templates.Data`):

- `{{.ProjectName}}` -- project name (e.g., `my-app`)
- `{{.GoModulePath}}` -- full module path (e.g., `github.com/user/my-app`)
- `{{.GoModuleOwner}}` -- GitHub owner (e.g., `user`)

### 2. Add the Step Method

In `internal/scaffold/steps.go`, add a method with a capability guard:

```go
func (s *Scaffolder) stepGenerateMyConfig() error {
    if !s.Config.IsEnabled(CapMyThing) {
        s.Logger.Info("Skipping my config (--no-mything)")
        return nil
    }
    return WriteTemplateFile(
        filepath.Join(s.Config.ProjectDir, ".myconfig.yml"),
        "my_config.tmpl",
        s.templateData(),
        s.Config.DryRun,
        s.Logger,
    )
}
```

`WriteTemplateFile` handles:

- Skipping if file already exists (idempotency)
- Dry-run logging
- Directory creation
- Template rendering

### 3. Add the Capability (if new)

If adding a new toggleable capability:

1. Add a constant in `internal/scaffold/config.go`:
    ```go
    CapMyThing = "mything"
    ```
2. Add to `DefaultCapabilities()` with the desired default (ON or OFF)
3. Add to the `capabilities` slice in `cmd/root.go`:
    ```go
    {"mything", true, "My thing description"},
    ```

### 4. Register the Step

In `internal/scaffold/scaffold.go`, add the step to the `steps` slice:

```go
steps := []func() error{
    // ... existing steps ...
    s.stepGenerateMyConfig,
}
```

### 5. Add Tests

Add tests in `internal/scaffold/steps_test.go`:

```go
func TestStepGenerateMyConfigIdempotent(t *testing.T) {
    s, _, _ := testScaffolder(t, false)
    if err := s.stepGenerateMyConfig(); err != nil {
        t.Fatal(err)
    }
    // Verify file exists
    path := filepath.Join(s.Config.ProjectDir, ".myconfig.yml")
    if _, err := os.Stat(path); err != nil {
        t.Fatal("expected file to be created")
    }
    // Second call should skip (idempotent)
    if err := s.stepGenerateMyConfig(); err != nil {
        t.Fatal(err)
    }
}

func TestStepGenerateMyConfigDryRun(t *testing.T) {
    s, _, stderr := testScaffolder(t, true)
    if err := s.stepGenerateMyConfig(); err != nil {
        t.Fatal(err)
    }
    path := filepath.Join(s.Config.ProjectDir, ".myconfig.yml")
    if _, err := os.Stat(path); err == nil {
        t.Error("file should not exist in dry-run mode")
    }
    if !strings.Contains(stderr.String(), "[DRY-RUN]") {
        t.Errorf("expected dry-run message")
    }
}

func TestStepGenerateMyConfigDisabled(t *testing.T) {
    s, stdout, _ := testScaffolder(t, false)
    s.Config.Capabilities[CapMyThing] = false
    if err := s.stepGenerateMyConfig(); err != nil {
        t.Fatal(err)
    }
    if !strings.Contains(stdout.String(), "Skipping") {
        t.Errorf("expected skip message")
    }
}
```

And a template render test in `internal/templates/templates_test.go`:

```go
{"my_config.tmpl", []string{"expected-content"}},
```

## Capability Architecture

The capability system uses a `map[string]bool` pattern:

- **`Config.Capabilities`** -- maps capability names to enabled/disabled state
- **`Config.IsEnabled(name)`** -- checks if a capability is ON
- **`Config.Disable(name)`** -- turns off a capability at runtime (e.g., missing soft dep)
- **`DefaultCapabilities()`** -- returns the default state for all capabilities

Flags are registered in `cmd/root.go` as `--<name>` (visible) and `--no-<name>` (hidden). The `--no-<name>` flag takes precedence if both are set.

Soft dependency auto-disable happens in `ValidateEnvironment()` -- if a tool like `bun` or `uv` is missing, the corresponding capability is disabled automatically with a warning.

## Code Style

- Run `make fmt` before committing
- Run `make vet` and `make lint` to catch issues
- Tests: `make test` (with race detector)
- Follow existing patterns for consistency

## Template Escaping

When templates generate files that themselves use `{{ }}` syntax (like goreleaser configs), escape the inner template variables:

```
{{"{{"}} .Version {{"}}"}}
```

This renders as `{{ .Version }}` in the output file, while `{{.ProjectName}}` is replaced at scaffold time.
