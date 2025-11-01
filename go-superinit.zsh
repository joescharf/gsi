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
PROJECT_NAME=""

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
    -h, --help          Show this help message

Examples:
    $script_name my-awesome-app
    $script_name --author "Jane Doe jane@example.com" my-app
    $script_name --module github.com/myorg/myapp --dry-run my-app
    $script_name --skip-bmad --skip-git my-app
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
    local optional_commands=("git" "npx")
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
    local issues=()

    log_info "Checking existing state..."

    # Check if go.mod exists
    if [[ -f "$project_dir/go.mod" ]]; then
        issues+=("go.mod already exists")
    fi

    # Check if .git exists
    if [[ -d "$project_dir/.git" ]] && [[ "$SKIP_GIT" == false ]]; then
        issues+=(".git directory already exists")
    fi

    # Check if cmd directory exists
    if [[ -d "$project_dir/cmd" ]]; then
        issues+=("cmd/ directory already exists")
    fi

    if [[ ${#issues[@]} -gt 0 ]]; then
        log_warning "Found existing project files:"
        for issue in "${issues[@]}"; do
            log_warning "  - $issue"
        done
        echo
        read "response?Continue anyway? (y/N) "
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            log_error "Aborted by user"
            return 1
        fi
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
            # Directory exists - check if it's empty
            if [[ -n "$(ls -A "$PROJECT_DIR" 2>/dev/null)" ]]; then
                log_error "Directory '$PROJECT_DIR' already exists and is not empty"
                log_error "Please use a different name or remove existing files"
                return 1
            else
                log_warning "Directory '$PROJECT_DIR' exists but is empty, using it"
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
    if [[ "$DRY_RUN" == true ]]; then
        echo "  ${YELLOW}Mode:          DRY-RUN${NC}"
    fi
    echo

    # Validate environment
    validate_environment || return 1

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
    if [[ "$SKIP_BMAD" == false ]] && command -v npx &> /dev/null; then
        execute "npx bmad-method install -f -i claude-code -d ./" \
                "Installing BMAD method framework"
    elif [[ "$SKIP_BMAD" == false ]]; then
        log_warning "Skipping BMAD installation (npx not found)"
    else
        log_info "Skipping BMAD installation (--skip-bmad flag)"
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
        log_warning "go.mod already exists, skipping go mod init"
    fi

    # Initialize Cobra CLI application
    if [[ ! -d "cmd" ]] || [[ "$DRY_RUN" == true ]]; then
        execute "cobra-cli init --viper --author \"$AUTHOR\" --config \$HOME/.config/$PROJECT_NAME" \
                "Creating Cobra CLI application structure"
    else
        log_warning "cmd/ directory already exists, skipping cobra-cli init"
    fi

    # Add version command
    if [[ ! -f "cmd/version.go" ]] || [[ "$DRY_RUN" == true ]]; then
        execute "cobra-cli add version" \
                "Adding version command"
    else
        log_warning "cmd/version.go already exists, skipping"
    fi

    # Tidy dependencies
    execute "go mod tidy" \
            "Tidying Go dependencies"

    # Initialize mockery
    if [[ ! -f ".mockery.yml" ]] || [[ "$DRY_RUN" == true ]]; then
        # Check if mockery is available
        if command -v mockery &> /dev/null; then
            execute "mockery init $PROJECT_NAME" \
                    "Initializing mockery for test mocks"
        else
            log_warning "mockery not found, skipping mockery init"
            log_info "Install with: go install github.com/vektra/mockery/v2@latest"
        fi
    else
        log_info ".mockery.yml already exists, skipping mockery init"
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

    # Git initialization
    if [[ "$SKIP_GIT" == false ]]; then
        # Initialize git repository
        if [[ ! -d ".git" ]] || [[ "$DRY_RUN" == true ]]; then
            execute "git init" \
                    "Initializing git repository"
        else
            log_warning ".git directory already exists, skipping git init"
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
                log_warning "Git repository already has commits, skipping initial commit"
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
    echo "  1. Review the generated code in cmd/"
    echo "  2. Update the project description in cmd/root.go"
    echo "  3. Run 'go build' to build your application"
    echo "  4. Run './$PROJECT_NAME --help' to see available commands"
    echo
}

# Run main function with all arguments
main "$@"
