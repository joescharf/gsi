# Getting Started

## Installation

### Go Install (recommended)

```bash
go install github.com/joescharf/gsi@latest
```

### Homebrew

```bash
brew install joescharf/tap/gsi
```

### Docker

```bash
docker run --rm -v $(pwd):/workspace -w /workspace ghcr.io/joescharf/gsi:latest my-project
```

### From Source

```bash
git clone https://github.com/joescharf/gsi.git
cd gsi
make build
./bin/gsi --help
```

## Your First Project

### 1. Scaffold

```bash
gsi my-app
```

Expected output:

```
  Project Name:  my-app
  Project Dir:   /home/user/my-app
  Module Path:   github.com/joescharf/my-app
  Author:        Joe Scharf joe@joescharf.com

Starting project initialization...

  Installing BMAD method framework
  cobra-cli is already installed
  Initializing Go module
  Creating Cobra CLI application structure
  Adding version command
  Creating cmd/serve.go
  Creating .mockery.yml
  Creating .editorconfig
  Creating internal/ui/dist/index.html
  Creating internal/ui/embed.go
  Tidying Go dependencies
  Creating Makefile
  Creating .goreleaser.yml
  Creating Dockerfile
  Creating .dockerignore
  Initializing docs...
  Creating .gitignore
  Creating initial commit

Project initialization complete!
```

### 2. Build and Run

```bash
cd my-app
make build
./bin/my-app --help
```

### 3. Start Development

```bash
make run          # Build and run
make serve        # Start embedded web UI server
make docs-serve   # Start docs dev server (requires uv)
make test         # Run tests
make help         # See all available targets
```

## Common Flag Combinations

```bash
# Custom module path and author
gsi --module github.com/myorg/myapp --author "Jane Doe jane@example.com" my-app

# Dry run to preview what would be created
gsi --dry-run my-app

# Include React/shadcn UI
gsi --ui my-app

# Skip optional components
gsi --skip-bmad --skip-docs my-app

# Add docs to an existing project
gsi --only-docs my-app

# Initialize in current directory
gsi .
```

## Next Steps

- Review [CLI Reference](cli-reference.md) for all available flags
- See [What Gets Scaffolded](scaffolded-output.md) for a breakdown of generated files
- Read [Releasing](releasing.md) to publish your first release
