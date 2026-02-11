Perfect! I have all the information needed. Here's the comprehensive prompt for the AI code generator:

---

# Go SuperInit (gsi) - AI Code Generator Prompt

## Project Overview

Create a Go-based CLI application called "Go SuperInit" (command: `gsi`) that automates the setup of new Go projects according to user preferences. The tool streamlines project initialization by automatically creating directory structures, installing dependencies, configuring tools, and executing setup commands.

## Core Functionality

### 1. Command Line Interface

#### Main Command

```bash
gsi [project-name] [flags]
```

**Behavior:**

- If `project-name` is provided and directory doesn't exist: create it and initialize there
- If `project-name` is provided and directory exists with files: refuse to proceed with error message
- If `project-name` is not provided: use current directory and prompt user for project name interactively
- Validate that Go is installed and version matches the preferred version in config before proceeding

**Flags:**

- `--dry-run`: Show what would be done without executing any changes
- `--verbose`: Display detailed output during execution
- `--template <name>`: Use a specific named template from config (defaults to `default_template` in settings)

#### Subcommands

```bash
gsi config edit    # Open settings.yaml in user's default editor ($EDITOR or fallback)
gsi config show    # Display current configuration
gsi config reset   # Reset to default configuration
gsi version        # Display gsi version
```

### 2. Configuration Management

**Configuration File Location:** `~/.config/gsi/settings.yaml`

**First Run Behavior:**

- If configuration file doesn't exist, create it with sensible defaults
- Create the configuration directory if it doesn't exist
- No interactive prompts on first run

**Configuration Structure:**

```yaml
go_version: "1.23"

libraries:
  - name: "github.com/spf13/viper"
    version: "v1.18.2"
  - name: "github.com/spf13/cobra"
    version: "latest"
  - name: "github.com/jackc/pgx/v5"
    version: "v5.5.0"

tools:
  - name: "github.com/vektra/mockery/v2"
    version: "latest"
  - name: "github.com/goreleaser/goreleaser"
    version: "latest"
  - name: "github.com/golangci/golangci-lint/cmd/golangci-lint"
    version: "latest"

custom_commands:
  - name: "Initialize Claude Code"
    command: "claude code init"
    halt_on_failure: false
  - name: "Setup BMAD"
    command: "bmad setup"
    halt_on_failure: true

templates:
  default:
    type: "local"
    path: "~/.config/gsi/templates/standard"
  web-api:
    type: "git"
    url: "https://github.com/username/go-web-template"
    branch: "main"
    subdir: ""
  cli-tool:
    type: "git"
    url: "https://github.com/username/go-cli-template"

default_template: "default"
```

**Default Configuration Values:**

- Go version: Latest stable version
- Libraries: Empty list (user adds as needed)
- Tools: Common tools like `golangci-lint`
- Templates: A minimal "default" template with basic Go project structure (cmd/, internal/, pkg/, .gitignore, README.md)
- Custom commands: Empty list

### 3. Project Initialization Sequence

When `gsi myproject` is executed, perform operations in this order:

1. **Validate Environment**

   - Check Go is installed
   - Verify Go version matches `go_version` in config
   - If mismatch, halt with clear error message

2. **Create/Validate Project Directory**

   - Create directory if it doesn't exist
   - If directory exists and contains files, refuse to proceed
   - Change working directory to project directory

3. **Initialize Git Repository**

   - Run `git init`

4. **Copy Template Files**

   - Based on selected template (from `--template` flag or `default_template`)
   - **Local template:** Recursively copy all files and subdirectories from specified path
   - **Git template:** Clone repository, optionally checkout specific branch, extract files from subdir if specified, remove .git directory
   - Exclude files/patterns based on .gitignore-style rules (if .gsiignore exists in template, use it; otherwise use common patterns like `.git/`, `.DS_Store`, `node_modules/`)

5. **Process Template Variables**

   - In both file contents and filenames, replace template variables:
     - `{{.ProjectName}}`: The project name
     - `{{.GoVersion}}`: Go version from config
     - `{{.ModulePath}}`: Full module path (prompt user or derive from project name)
   - Support standard Go template syntax

6. **Run `go mod init`**

   - Initialize Go module: `go mod init <module-path>`

7. **Install Libraries**

   - For each library in config, run `go get <name>@<version>`
   - If version is "latest", use `@latest`
   - Otherwise use specific version (e.g., `@v1.18.2`)

8. **Install Tools**

   - For each tool in config:
     - Check if tool is already installed (check PATH)
     - If not installed, run `go install <name>@<version>`
     - If version is "latest", use `@latest`

9. **Execute Custom Commands**

   - For each custom command:
     - Execute in project directory
     - If `halt_on_failure` is true and command fails, stop initialization and display error
     - If `halt_on_failure` is false and command fails, log warning and continue
     - Capture and display command output

10. **Final Status**
    - Display success message with summary of what was created
    - List next steps for the user

### 4. Template System Details

**Local Templates:**

- Recursively copy all files and directories
- Process template variables in file contents and filenames
- Respect exclusion patterns

**Git Templates:**

- Clone the repository to a temporary location
- Checkout specified branch if provided
- If `subdir` is specified, only copy files from that subdirectory
- Remove .git directory after copying
- Process template variables same as local templates
- If clone fails or URL is unreachable, halt with error

**Template Variables:**
Available variables for templating:

- `{{.ProjectName}}`: Project name
- `{{.GoVersion}}`: Go version from config
- `{{.ModulePath}}`: Go module path
- Use Go's `text/template` package for processing

**Exclusion Patterns:**

- Check for `.gsiignore` file in template root
- If not found, use default exclusions: `.git/`, `.DS_Store`, `*.swp`, `node_modules/`, `.idea/`
- Parse using gitignore-style pattern matching

### 5. Error Handling

**Validation Errors:**

- Go not installed: "Error: Go is not installed. Please install Go and try again."
- Go version mismatch: "Error: Go version {{installed}} does not match required version {{required}}"
- Project directory exists with files: "Error: Directory '{{name}}' already exists and is not empty. Please use a different name or remove existing files."

**Template Errors:**

- Template not found: "Error: Template '{{name}}' not found in configuration"
- Local path doesn't exist: "Error: Template path '{{path}}' does not exist"
- Git clone fails: "Error: Failed to clone template from '{{url}}': {{error}}"

**Execution Errors:**

- Library installation fails: Log warning, continue unless critical
- Tool installation fails: Log warning, continue
- Custom command fails with `halt_on_failure: true`: Stop and display error
- Custom command fails with `halt_on_failure: false`: Log warning, continue

### 6. Output and Logging

**Standard Output:**

- Show progress for each major step
- Use clear, concise messages
- Display summary at end

**Verbose Mode (`--verbose`):**

- Show detailed output for each operation
- Display full command output
- Show file operations (copies, template processing)

**Dry-run Mode (`--dry-run`):**

- Display all operations that would be performed
- Show resolved template variables
- List files that would be created/copied
- Show commands that would be executed
- Do not make any actual changes

## Technical Requirements

### Implementation Language

- Go 1.21 or higher

### Recommended Libraries

- **CLI Framework:** `github.com/spf13/cobra` for command structure
- **Configuration:** `github.com/spf13/viper` for reading YAML config
- **Templating:** Go's standard `text/template` package
- **Git Operations:** `github.com/go-git/go-git/v5` for cloning repositories
- **YAML:** `gopkg.in/yaml.v3` for parsing configuration

### Code Structure

```
gsi/
├── cmd/
│   └── gsi/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go        # Configuration loading and management
│   │   └── defaults.go      # Default configuration values
│   ├── project/
│   │   ├── initializer.go   # Main initialization logic
│   │   ├── validator.go     # Environment validation
│   │   └── template.go      # Template processing
│   ├── installer/
│   │   ├── libraries.go     # Library installation
│   │   ├── tools.go         # Tool installation
│   │   └── commands.go      # Custom command execution
│   └── cli/
│       ├── root.go          # Root command
│       ├── init.go          # Init command (default)
│       ├── config.go        # Config subcommands
│       └── version.go       # Version command
├── pkg/
│   └── template/
│       ├── processor.go     # Template variable processing
│       └── copier.go        # File copying with exclusions
├── .gitignore
├── .goreleaser.yaml
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### Key Implementation Details

1. **Configuration Management:**

   - Use viper to read YAML configuration
   - Implement config file creation with defaults on first run
   - Support expanding `~` in paths to user's home directory

2. **Template Processing:**

   - Use Go's `text/template` for variable substitution
   - Process both file contents and filenames
   - Handle binary files gracefully (skip template processing for non-text files)

3. **Git Operations:**

   - Use go-git for cloning repositories
   - Implement shallow clone for efficiency
   - Clean up temporary directories after copying

4. **Command Execution:**

   - Use `os/exec` for running external commands
   - Capture stdout and stderr
   - Implement proper error handling and logging

5. **Path Handling:**
   - Support absolute and relative paths
   - Expand `~` to user home directory
   - Handle cross-platform path separators

## Testing Requirements

Include unit tests for:

- Configuration loading and validation
- Template variable processing
- File copying with exclusions
- Command execution logic
- Error handling scenarios

## Documentation

Include in README.md:

- Installation instructions
- Quick start guide
- Configuration file format and options
- Template creation guide
- Examples of common use cases
- Troubleshooting section

## Additional Features for Initial Release

1. **Version Command:** Display tool version and Go version
2. **Config Validation:** Validate settings.yaml on load and provide helpful error messages
3. **Progress Indicators:** Show progress for long-running operations (cloning, installing)
4. **Colored Output:** Use colored terminal output for better readability (success=green, error=red, warning=yellow)

## Success Criteria

The generated application should:

1. Successfully initialize a new Go project with all specified components
2. Handle all error cases gracefully with clear messages
3. Support both local and git-based templates
4. Process template variables correctly
5. Install libraries and tools as configured
6. Execute custom commands with proper error handling
7. Provide a smooth user experience with clear progress indicators
8. Include comprehensive documentation
9. Be maintainable and well-structured for future enhancements

---

This prompt provides all necessary details for an AI code generator to create a fully functional Go SuperInit application that meets your requirements.
