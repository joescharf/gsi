#!/bin/zsh

# go-superinit.zsh - Initialize a Go project with best practices
# Author: Joe Scharf <joe@joescharf.com>

# Exit on error, undefined variables, and pipe failures
set -euo pipefail

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Default values
DEFAULT_AUTHOR="Joe Scharf joe@joescharf.com"
DEFAULT_MODULE_PREFIX="github.com/joescharf"
DRY_RUN=false
VERBOSE=false
SKIP_BMAD=false
SKIP_GIT=false
SKIP_DOCS=false
INIT_UI=false
PROJECT_NAME=""

# BMAD installation command - modify this to change the bmad installation
BMAD_INSTALL_CMD="bunx bmad-method@alpha install"

# Cleanup function for error handling
cleanup() {
    local exit_code=$?
    if [[ $exit_code -ne 0 ]]; then
        echo "${RED}✗ Script failed with exit code $exit_code${NC}" >&2
    fi
    exit $exit_code
}

trap cleanup EXIT INT TERM

# Logging functions
log_info() {
    echo "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo "${GREEN}✓${NC} $1"
}

log_warning() {
    echo "${YELLOW}⚠${NC} $1" >&2
}

log_error() {
    echo "${RED}✗${NC} $1" >&2
}

log_verbose() {
    if [[ "$VERBOSE" == true ]]; then
        echo "${BLUE}  →${NC} $1"
    fi
}

# Execute command with optional dry-run
execute() {
    local cmd="$1"
    local description="$2"

    log_info "$description"
    log_verbose "Command: $cmd"

    if [[ "$DRY_RUN" == true ]]; then
        log_warning "[DRY-RUN] Would execute: $cmd"
        return 0
    fi

    if eval "$cmd"; then
        log_success "$description - Done"
        return 0
    else
        log_error "$description - Failed"
        return 1
    fi
}

# Usage information
usage() {
    local script_name="${0:t}"
    cat << EOF
Usage: $script_name [OPTIONS] [PROJECT_NAME]

Initialize a new Go project with best practices and tooling.

Arguments:
    PROJECT_NAME        Name of the project (required, use '.' for current directory)

Options:
    -a, --author TEXT   Author name and email (default: $DEFAULT_AUTHOR)
    -m, --module PATH   Go module path (default: $DEFAULT_MODULE_PREFIX/PROJECT_NAME)
    -d, --dry-run       Show what would be done without executing
    -v, --verbose       Enable verbose output
    --skip-bmad         Skip BMAD method installation
    --skip-git          Skip git initialization and commit
    --skip-docs         Skip mkdocs-material documentation scaffolding
    --ui                Initialize a React/shadcn/Tailwind UI in ui/ subdirectory (requires bun)
    -h, --help          Show this help message

Examples:
    $script_name my-awesome-app
    $script_name --author "Jane Doe jane@example.com" my-app
    $script_name --module github.com/myorg/myapp --dry-run my-app
    $script_name --skip-bmad --skip-git my-app
    $script_name --skip-docs my-app
    $script_name --ui my-app
    $script_name .    # Initialize in current directory

EOF
}

# Validation functions
check_command() {
    local cmd="$1"
    if command -v "$cmd" &> /dev/null; then
        log_verbose "Found $cmd: $(command -v $cmd)"
        return 0
    else
        log_error "$cmd is not installed or not in PATH"
        return 1
    fi
}

validate_environment() {
    log_info "Validating environment..."

    local required_commands=("go")
    local optional_commands=("git" "bun" "uv")
    local missing_required=()
    local missing_optional=()

    # Check required commands
    for cmd in "${required_commands[@]}"; do
        if ! check_command "$cmd"; then
            missing_required+=("$cmd")
        fi
    done

    # Check optional commands
    for cmd in "${optional_commands[@]}"; do
        if ! check_command "$cmd"; then
            missing_optional+=("$cmd")
            log_warning "$cmd is not installed (optional)"
        fi
    done

    if [[ ${#missing_required[@]} -gt 0 ]]; then
        log_error "Missing required commands: ${missing_required[*]}"
        log_error "Please install them and try again"
        return 1
    fi

    # Show Go version
    local go_version=$(go version)
    log_verbose "Go version: $go_version"

    log_success "Environment validation complete"
    return 0
}

check_existing_state() {
    local project_dir="$1"
    local existing=()

    log_info "Checking existing state..."

    # Check if go.mod exists
    if [[ -f "$project_dir/go.mod" ]]; then
        existing+=("go.mod")
    fi

    # Check if .git exists
    if [[ -d "$project_dir/.git" ]]; then
        existing+=(".git/")
    fi

    # Check if cmd directory exists
    if [[ -d "$project_dir/cmd" ]]; then
        existing+=("cmd/")
    fi

    # Check if _bmad directory exists
    if [[ -d "$project_dir/_bmad" ]]; then
        existing+=("_bmad/")
    fi

    # Check if ui directory exists
    if [[ -d "$project_dir/ui" ]]; then
        existing+=("ui/")
    fi

    # Check if docs directory exists
    if [[ -d "$project_dir/docs" ]]; then
        existing+=("docs/")
    fi

    if [[ ${#existing[@]} -gt 0 ]]; then
        log_info "Found existing project files (will skip where appropriate):"
        for item in "${existing[@]}"; do
            log_verbose "  - $item"
        done
    else
        log_success "No existing project files found"
    fi

    return 0
}

install_tool_if_missing() {
    local tool_name="$1"
    local install_cmd="$2"
    local check_cmd="${3:-$tool_name}"

    if command -v "$check_cmd" &> /dev/null; then
        log_success "$tool_name is already installed"
        log_verbose "$check_cmd location: $(command -v $check_cmd)"
        return 0
    fi

    execute "$install_cmd" "Installing $tool_name"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -a|--author)
                AUTHOR="$2"
                shift 2
                ;;
            -m|--module)
                GO_MODULE_PATH="$2"
                shift 2
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --skip-bmad)
                SKIP_BMAD=true
                shift
                ;;
            --skip-git)
                SKIP_GIT=true
                shift
                ;;
            --skip-docs)
                SKIP_DOCS=true
                shift
                ;;
            --ui)
                INIT_UI=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            -*)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
            *)
                PROJECT_NAME="$1"
                shift
                ;;
        esac
    done
}

# Main execution
main() {
    # Parse arguments
    parse_args "$@"

    # Require project name
    if [[ -z "$PROJECT_NAME" ]]; then
        log_error "Project name is required"
        echo
        usage
        return 1
    fi

    # Validate project name format (allow alphanumeric, hyphens, underscores, slashes, and dots)
    if [[ ! "$PROJECT_NAME" =~ ^[a-zA-Z0-9_/.-]+$ ]]; then
        log_error "Invalid project name: must contain only letters, numbers, hyphens, underscores, dots, and slashes"
        return 1
    fi

    # Determine if using current directory or creating new one
    if [[ "$PROJECT_NAME" == "." ]] || [[ "$PROJECT_NAME" == "./" ]]; then
        # Use current directory
        PROJECT_DIR="$(pwd)"
        PROJECT_NAME="$(basename "$PROJECT_DIR")"
        log_info "Initializing in current directory"
        log_verbose "Project directory: $PROJECT_DIR"
    else
        # Create new directory with project name
        # Resolve to absolute path
        if [[ "$PROJECT_NAME" = /* ]]; then
            PROJECT_DIR="$PROJECT_NAME"
        else
            PROJECT_DIR="$(pwd)/$PROJECT_NAME"
        fi

        # Extract just the name (last component) for module path
        PROJECT_NAME="$(basename "$PROJECT_DIR")"

        # Check if directory already exists
        if [[ -d "$PROJECT_DIR" ]]; then
            # Directory exists - use it (idempotent behavior)
            if [[ -n "$(ls -A "$PROJECT_DIR" 2>/dev/null)" ]]; then
                log_info "Directory '$PROJECT_DIR' already exists, continuing with initialization"
            else
                log_info "Directory '$PROJECT_DIR' exists but is empty, using it"
            fi
        else
            # Create the directory
            if [[ "$DRY_RUN" == false ]]; then
                log_info "Creating project directory: $PROJECT_DIR"
                if ! mkdir -p "$PROJECT_DIR"; then
                    log_error "Failed to create directory: $PROJECT_DIR"
                    return 1
                fi
                log_success "Created project directory"
            else
                log_warning "[DRY-RUN] Would create directory: $PROJECT_DIR"
            fi
        fi

        # Change into project directory
        if [[ "$DRY_RUN" == false ]]; then
            if ! cd "$PROJECT_DIR"; then
                log_error "Failed to change to directory: $PROJECT_DIR"
                return 1
            fi
            log_verbose "Changed to project directory: $PROJECT_DIR"
        else
            log_warning "[DRY-RUN] Would change to directory: $PROJECT_DIR"
        fi
    fi

    # Set defaults
    AUTHOR="${AUTHOR:-$DEFAULT_AUTHOR}"
    GO_MODULE_PATH="${GO_MODULE_PATH:-$DEFAULT_MODULE_PREFIX/$PROJECT_NAME}"

    # Show configuration
    echo
    log_info "Configuration:"
    echo "  Project Name:  $PROJECT_NAME"
    echo "  Project Dir:   $PROJECT_DIR"
    echo "  Module Path:   $GO_MODULE_PATH"
    echo "  Author:        $AUTHOR"
    echo "  Config Dir:    $HOME/.config/$PROJECT_NAME"
    echo "  Init UI:       $INIT_UI"
    echo "  Skip Docs:     $SKIP_DOCS"
    if [[ "$DRY_RUN" == true ]]; then
        echo "  ${YELLOW}Mode:          DRY-RUN${NC}"
    fi
    echo

    # Validate environment
    validate_environment || return 1

    # Validate bun is available if --ui is requested
    if [[ "$INIT_UI" == true ]]; then
        if ! command -v bun &> /dev/null; then
            log_error "bun is required for UI initialization (--ui flag) but is not installed"
            log_error "Install bun: https://bun.sh"
            return 1
        fi
    fi

    # Validate uv is available for docs scaffolding
    if [[ "$SKIP_DOCS" == false ]]; then
        if ! command -v uv &> /dev/null; then
            log_warning "uv is not installed, skipping docs scaffolding (install: https://docs.astral.sh/uv/)"
            SKIP_DOCS=true
        fi
    fi

    # Check existing state (in dry-run, check PROJECT_DIR; otherwise check current dir since we already cd'd)
    if [[ "$DRY_RUN" == false ]]; then
        check_existing_state "$PWD" || return 1
    else
        check_existing_state "$PROJECT_DIR" || return 1
    fi

    echo
    log_info "Starting project initialization..."
    echo

    # Install BMAD method framework
    if [[ "$SKIP_BMAD" == true ]]; then
        log_info "Skipping BMAD installation (--skip-bmad flag)"
    elif [[ -d "_bmad" ]]; then
        log_info "_bmad/ directory already exists, skipping BMAD installation"
    elif command -v bun &> /dev/null; then
        execute "$BMAD_INSTALL_CMD" \
                "Installing BMAD method framework"
    else
        log_warning "Skipping BMAD installation (bun not found)"
    fi

    # Install cobra-cli if needed
    install_tool_if_missing "cobra-cli" \
                           "go install github.com/spf13/cobra-cli@latest" \
                           "cobra-cli"

    # Initialize Go module
    if [[ ! -f "go.mod" ]] || [[ "$DRY_RUN" == true ]]; then
        execute "go mod init $GO_MODULE_PATH" \
                "Initializing Go module"
    else
        log_info "go.mod already exists, skipping go mod init"
    fi

    # Initialize Cobra CLI application
    if [[ ! -d "cmd" ]] || [[ "$DRY_RUN" == true ]]; then
        execute "cobra-cli init --viper --author \"$AUTHOR\" --config \$HOME/.config/$PROJECT_NAME" \
                "Creating Cobra CLI application structure"
    else
        log_info "cmd/ directory already exists, skipping cobra-cli init"
    fi

    # Add version command
    if [[ ! -f "cmd/version.go" ]] || [[ "$DRY_RUN" == true ]]; then
        execute "cobra-cli add version" \
                "Adding version command"
    else
        log_info "cmd/version.go already exists, skipping"
    fi

    # Create serve command
    if [[ ! -f "cmd/serve.go" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating cmd/serve.go"
            mkdir -p cmd
            cat > cmd/serve.go << GOEOF
package cmd

import (
	"fmt"
	"net/http"

	"${GO_MODULE_PATH}/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the embedded web UI server",
	Long:  "Start an HTTP server that serves the embedded web UI.\nBy default it listens on port 8080. Use --port to change it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		port := viper.GetInt("port")

		handler, err := ui.Handler()
		if err != nil {
			return fmt.Errorf("failed to initialize UI handler: %w", err)
		}

		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Serving UI at http://localhost%s\n", addr)
		return http.ListenAndServe(addr, handler)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntP("port", "p", 8080, "port to listen on")
	viper.SetDefault("port", 8080)
	_ = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}
GOEOF
            log_success "Created cmd/serve.go"
        else
            log_warning "[DRY-RUN] Would create cmd/serve.go"
        fi
    else
        log_info "cmd/serve.go already exists, skipping"
    fi

    # Create mockery configuration
    if [[ ! -f ".mockery.yml" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating .mockery.yml configuration"
            cat > .mockery.yml << EOF
with-expecter: true
packages:
  $GO_MODULE_PATH:
    config:
      recursive: true
      dir: "{{.InterfaceDir}}/mocks"
      mockname: "Mock{{.InterfaceName}}"
      outpkg: "mocks"
EOF
            log_success "Created .mockery.yml"
        else
            log_warning "[DRY-RUN] Would create .mockery.yml"
        fi
    else
        log_info ".mockery.yml already exists, skipping"
    fi

    # Create .editorconfig
    if [[ ! -f ".editorconfig" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating .editorconfig"
            cat > .editorconfig << 'EOF'
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[{*.go,Makefile,.gitmodules,go.mod,go.sum}]
indent_style = tab

[*.md]
indent_style = tab
trim_trailing_whitespace = false

[*.{yml,yaml,json}]
indent_style = space
indent_size = 2

[*.{js,jsx,ts,tsx,css,less,sass,scss,vue,py}]
indent_style = space
indent_size = 4
EOF
            log_success "Created .editorconfig"
        else
            log_warning "[DRY-RUN] Would create .editorconfig"
        fi
    else
        log_info ".editorconfig already exists, skipping"
    fi

    # Create default placeholder UI page
    if [[ ! -f "internal/ui/dist/index.html" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating internal/ui/dist/index.html"
            mkdir -p internal/ui/dist
            cat > internal/ui/dist/index.html << EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${PROJECT_NAME}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: #0f172a;
            color: #e2e8f0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            text-align: center;
            max-width: 480px;
            padding: 2rem;
        }
        h1 {
            font-size: 2rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
            color: #f8fafc;
        }
        .status {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            background: #1e293b;
            border: 1px solid #334155;
            border-radius: 9999px;
            padding: 0.5rem 1rem;
            margin: 1.5rem 0;
            font-size: 0.875rem;
        }
        .dot {
            width: 8px;
            height: 8px;
            background: #22c55e;
            border-radius: 50%;
            animation: pulse 2s ease-in-out infinite;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.4; }
        }
        .hint {
            color: #94a3b8;
            font-size: 0.875rem;
            line-height: 1.6;
            margin-top: 1.5rem;
        }
        code {
            background: #1e293b;
            border: 1px solid #334155;
            border-radius: 4px;
            padding: 0.15rem 0.4rem;
            font-size: 0.8rem;
            color: #7dd3fc;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>${PROJECT_NAME}</h1>
        <div class="status">
            <span class="dot"></span>
            Server is running
        </div>
        <p class="hint">
            This is the default placeholder page.<br>
            To embed your React UI, run:<br>
            <code>make ui-build && make ui-embed</code>
        </p>
    </div>
</body>
</html>
EOF
            log_success "Created internal/ui/dist/index.html"
        else
            log_warning "[DRY-RUN] Would create internal/ui/dist/index.html"
        fi
    else
        log_info "internal/ui/dist/index.html already exists, skipping"
    fi

    # Create embed.go for serving the UI
    if [[ ! -f "internal/ui/embed.go" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating internal/ui/embed.go"
            mkdir -p internal/ui
            cat > internal/ui/embed.go << 'GOEOF'
package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:dist
var distFS embed.FS

// DistFS returns the embedded dist/ filesystem with the "dist" prefix stripped.
func DistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}

// Handler returns an http.Handler that serves the embedded UI with SPA fallback.
// Static files are served directly. Paths without a file extension are treated as
// client-side routes and served index.html. Missing assets return 404.
func Handler() (http.Handler, error) {
	sub, err := DistFS()
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServerFS(sub)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path
		p := path.Clean(r.URL.Path)
		if p == "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Strip leading slash for fs operations
		p = strings.TrimPrefix(p, "/")

		// Check if the file exists in the embedded FS
		_, err := fs.Stat(sub, p)
		if err == nil {
			// File exists, serve it directly
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist — check if it looks like a static asset
		if strings.Contains(p, ".") {
			// Has extension (e.g. .js, .css, .png) — genuine missing asset
			http.NotFound(w, r)
			return
		}

		// No extension — treat as SPA client-side route, serve index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}), nil
}
GOEOF
            log_success "Created internal/ui/embed.go"
        else
            log_warning "[DRY-RUN] Would create internal/ui/embed.go"
        fi
    else
        log_info "internal/ui/embed.go already exists, skipping"
    fi

    # Tidy dependencies (after all Go files are generated)
    execute "go mod tidy" \
            "Tidying Go dependencies"

    # Create Makefile
    if [[ ! -f "Makefile" ]] || [[ "$DRY_RUN" == true ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating Makefile"
            cat > Makefile << 'MAKEEOF'
# Makefile
BINARY_NAME := $(shell basename $(CURDIR))
MODULE := $(shell head -1 go.mod | awk '{print $$2}')
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

# Conditionally include UI and docs targets if their directories exist
ALL_TARGETS := build
$(if $(wildcard ui/package.json),$(eval ALL_TARGETS += ui-build ui-embed))
$(if $(wildcard docs/mkdocs.yml),$(eval ALL_TARGETS += docs-build))

.DEFAULT_GOAL := all

##@ App
.PHONY: build install run serve clean tidy test lint vet fmt mocks

build: ## Build the Go binary
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

install: ## Install the binary to $GOPATH/bin
	go install $(LDFLAGS) .

run: build ## Build and run the binary
	./bin/$(BINARY_NAME)

serve: build ## Start the embedded web UI server
	./bin/$(BINARY_NAME) serve

clean: ## Remove build artifacts
	rm -rf bin/
	rm -f coverage.out

tidy: ## Run go mod tidy
	go mod tidy

test: ## Run tests
	go test -v -race -count=1 ./...

test-cover: ## Run tests with coverage
	go test -v -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: vet ## Run golangci-lint
	@which golangci-lint > /dev/null 2>&1 || { echo "Install golangci-lint: https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Run gofmt
	gofmt -s -w .

mocks: ## Generate mocks with mockery
	@which mockery > /dev/null 2>&1 || { echo "Install mockery: go install github.com/vektra/mockery/v2@latest"; exit 1; }
	mockery

##@ Docs (mkdocs-material via uv)
.PHONY: docs-serve docs-build docs-deps

docs-serve: ## Serve docs locally (requires uv + docs/ directory)
	@[ -d docs ] && [ -f docs/mkdocs.yml ] || { echo "No docs/ directory with mkdocs.yml found."; exit 1; }
	cd docs && uv run mkdocs serve

docs-build: ## Build docs site (requires uv + docs/ directory)
	@[ -d docs ] && [ -f docs/mkdocs.yml ] || { echo "No docs/ directory with mkdocs.yml found."; exit 1; }
	cd docs && uv run mkdocs build

docs-deps: ## Install doc dependencies (requires uv + docs/ directory)
	@[ -d docs ] && [ -f docs/pyproject.toml ] || { echo "No docs/ directory with pyproject.toml found."; exit 1; }
	cd docs && uv sync

##@ UI (React/shadcn via bun)
.PHONY: ui-dev ui-build ui-embed ui-deps

ui-dev: ## Start UI dev server (requires bun + ui/ directory)
	@[ -d ui ] && [ -f ui/package.json ] || { echo "No ui/ directory found. Re-run go-superinit with --ui to create one."; exit 1; }
	cd ui && bun dev

ui-build: ## Build UI for production (requires bun + ui/ directory)
	@[ -d ui ] && [ -f ui/package.json ] || { echo "No ui/ directory found. Re-run go-superinit with --ui to create one."; exit 1; }
	cd ui && bun run build

ui-embed: ## Copy built UI into internal/ui/dist for embedding
	@[ -d ui/dist ] || { echo "No ui/dist/ directory found. Run 'make ui-build' first."; exit 1; }
	rm -rf internal/ui/dist/*
	cp -r ui/dist/* internal/ui/dist/

ui-deps: ## Install UI dependencies (requires bun + ui/ directory)
	@[ -d ui ] && [ -f ui/package.json ] || { echo "No ui/ directory found. Re-run go-superinit with --ui to create one."; exit 1; }
	cd ui && bun install

##@ All
.PHONY: all deps dev

all: $(ALL_TARGETS) ## Build all existing artifacts (app + UI + docs)

deps: tidy ## Install all dependencies
	@[ -d docs ] && [ -f docs/pyproject.toml ] && (cd docs && uv sync) || true
	@[ -d ui ] && [ -f ui/package.json ] && (cd ui && bun install) || true

dev: ## Start all dev servers (app + docs + UI) in parallel
	@echo "Starting dev servers..."
	@$(MAKE) -j3 run docs-serve ui-dev 2>/dev/null || $(MAKE) run

##@ Help
.PHONY: help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
MAKEEOF
            log_success "Created Makefile"
        else
            log_warning "[DRY-RUN] Would create Makefile"
        fi
    else
        log_info "Makefile already exists, skipping"
    fi

    # Initialize mkdocs-material documentation
    if [[ "$SKIP_DOCS" == true ]]; then
        log_info "Skipping docs scaffolding (--skip-docs flag)"
    elif [[ -d "docs" ]] && [[ "$DRY_RUN" == false ]]; then
        log_info "docs/ directory already exists, skipping docs scaffolding"
    else
        # Initialize uv project in docs/
        execute "uv init --name ${PROJECT_NAME}-docs docs" \
                "Initializing uv project in docs/"

        # Add mkdocs-material dependencies
        execute "cd docs && uv add mkdocs-material 'mkdocs-git-revision-date-localized-plugin>=1.4' && cd .." \
                "Adding mkdocs-material dependencies"

        # Remove uv init scaffolding (placeholder files and nested .git)
        if [[ "$DRY_RUN" == false ]]; then
            rm -rf docs/.git
            rm -f docs/hello.py docs/main.py docs/README.md
        else
            log_warning "[DRY-RUN] Would remove uv init scaffolding (docs/.git, docs/hello.py, docs/main.py, docs/README.md)"
        fi

        # Write mkdocs.yml
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating docs/mkdocs.yml"
            cat > docs/mkdocs.yml << EOF
site_name: ${PROJECT_NAME} Documentation
site_url: https://${PROJECT_NAME}.example.com
site_description: "${PROJECT_NAME} documentation"
edit_uri: ""

extra_css:
  - stylesheets/extra.css

theme:
  name: material
  language: en
  features:
    - search.suggest
    - search.highlight
    - search.share
    - navigation.indexes
    - navigation.instant
    - navigation.instant.prefetch
    - navigation.instant.progress
    - content.code.copy
  palette:
    - media: "(prefers-color-scheme: light)"
      scheme: default
      toggle:
        icon: material/lightbulb-outline
        name: Switch to light mode
    - media: "(prefers-color-scheme: dark)"
      scheme: slate
      primary: light blue
      accent: indigo
      toggle:
        icon: material/lightbulb
        name: Switch to dark mode

plugins:
  - git-revision-date-localized
  - search

nav:
  - Welcome: index.md
  - Getting Started: getting-started.md

markdown_extensions:
  # Python Markdown
  - admonition
  - meta
  - footnotes
  - attr_list
  - def_list
  - toc:
      permalink: true

  # Python Markdown Extensions
  - pymdownx.highlight:
      anchor_linenums: true
      auto_title: false
  - pymdownx.inlinehilite
  - pymdownx.details
  - pymdownx.tilde
  - pymdownx.superfences:
      custom_fences:
        - name: mermaid
          class: mermaid
          format: !!python/name:pymdownx.superfences.fence_code_format
  - pymdownx.tasklist:
      custom_checkbox: true
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.emoji:
      emoji_index: !!python/name:material.extensions.emoji.twemoji
      emoji_generator: !!python/name:material.extensions.emoji.to_svg
EOF
            log_success "Created docs/mkdocs.yml"
        else
            log_warning "[DRY-RUN] Would create docs/mkdocs.yml"
        fi

        # Write docs/.gitignore
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating docs/.gitignore"
            cat > docs/.gitignore << 'EOF'
site/
.venv/
__pycache__/
.cache/
EOF
            log_success "Created docs/.gitignore"
        else
            log_warning "[DRY-RUN] Would create docs/.gitignore"
        fi

        # Create docs/docs/ content directory
        if [[ "$DRY_RUN" == false ]]; then
            mkdir -p docs/docs/stylesheets
        else
            log_warning "[DRY-RUN] Would create docs/docs/stylesheets/"
        fi

        # Write docs/docs/index.md
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating docs/docs/index.md"
            cat > docs/docs/index.md << EOF
# ${PROJECT_NAME}

Welcome to the **${PROJECT_NAME}** documentation.

## Quick Links

| Topic | Description |
|-------|-------------|
| [Getting Started](getting-started.md) | Installation and first steps |
EOF
            log_success "Created docs/docs/index.md"
        else
            log_warning "[DRY-RUN] Would create docs/docs/index.md"
        fi

        # Write docs/docs/getting-started.md
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating docs/docs/getting-started.md"
            cat > docs/docs/getting-started.md << EOF
# Getting Started

## Installation

\`\`\`bash
go install ${GO_MODULE_PATH}@latest
\`\`\`

## Usage

\`\`\`bash
${PROJECT_NAME} --help
\`\`\`
EOF
            log_success "Created docs/docs/getting-started.md"
        else
            log_warning "[DRY-RUN] Would create docs/docs/getting-started.md"
        fi

        # Write docs/docs/stylesheets/extra.css
        if [[ "$DRY_RUN" == false ]]; then
            log_info "Creating docs/docs/stylesheets/extra.css"
            cat > docs/docs/stylesheets/extra.css << 'EOF'
/* Compact navigation */
.md-nav__item {
  padding: 0.05rem 0;
}

/* Code block font size */
.md-typeset code {
  font-size: 0.8rem;
}

.md-typeset pre code {
  font-size: 0.8rem;
}
EOF
            log_success "Created docs/docs/stylesheets/extra.css"
        else
            log_warning "[DRY-RUN] Would create docs/docs/stylesheets/extra.css"
        fi
    fi

    # Initialize UI with React/shadcn/Tailwind
    if [[ "$INIT_UI" == true ]]; then
        if [[ -d "ui" ]] && [[ "$DRY_RUN" == false ]]; then
            log_info "ui/ directory already exists, skipping UI initialization"
        else
            execute "bun init --react=shadcn ui" \
                    "Initializing React/shadcn/Tailwind UI in ui/"
        fi
    fi

    # Git initialization
    if [[ "$SKIP_GIT" == false ]]; then
        # Initialize git repository
        if [[ ! -d ".git" ]] || [[ "$DRY_RUN" == true ]]; then
            execute "git init" \
                    "Initializing git repository"
        else
            log_info ".git directory already exists, skipping git init"
        fi

        # Create .gitignore
        if [[ ! -f ".gitignore" ]] || [[ "$DRY_RUN" == true ]]; then
            if [[ "$DRY_RUN" == false ]]; then
                log_info "Creating .gitignore for Go projects"
                cat > .gitignore << 'EOF'
# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Dependencies
vendor/

# Test files
*.test
coverage.out
*.prof

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Build artifacts
dist/
build/

# Environment
.env
.env.local

# UI (bun/React)
ui/node_modules/
internal/ui/dist/assets/

# Docs (mkdocs-material)
docs/.venv/
docs/site/
EOF
                log_success "Created .gitignore"
            else
                log_warning "[DRY-RUN] Would create .gitignore"
            fi
        else
            log_info ".gitignore already exists, skipping"
        fi

        # Initial commit
        if [[ "$DRY_RUN" == false ]]; then
            if git rev-parse HEAD &> /dev/null; then
                log_info "Git repository already has commits, skipping initial commit"
            else
                execute "git add ." \
                        "Staging files for initial commit"
                execute "git commit -m 'initial commit'" \
                        "Creating initial commit"
            fi
        else
            log_warning "[DRY-RUN] Would create initial commit"
        fi
    else
        log_info "Skipping git initialization (--skip-git flag)"
    fi

    echo
    log_success "Project initialization complete! 🎉"
    echo
    log_info "Next steps:"
    local step=1
    echo "  $step. Review the generated code in cmd/"
    ((step++))
    echo "  $step. Update the project description in cmd/root.go"
    ((step++))
    echo "  $step. Run 'make build' to build your application"
    ((step++))
    echo "  $step. Run 'make run' or './bin/$PROJECT_NAME --help' to see available commands"
    ((step++))
    echo "  $step. Run 'make serve' to start the embedded web UI server"
    if [[ "$SKIP_DOCS" == false ]]; then
        ((step++))
        echo "  $step. Run 'make docs-serve' to start the docs dev server"
    fi
    if [[ "$INIT_UI" == true ]]; then
        ((step++))
        echo "  $step. Run 'make ui-dev' to start the React dev server"
    fi
    ((step++))
    echo "  $step. Run 'make help' to see all available targets"
    echo
}

# Run main function with all arguments
main "$@"
