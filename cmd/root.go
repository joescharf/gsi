package cmd

import (
	"fmt"
	"os"

	"github.com/joescharf/gsi/internal/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// capabilityDef describes a scaffold capability flag.
type capabilityDef struct {
	name         string
	defaultValue bool
	description  string
}

// capabilities lists all toggleable scaffold capabilities.
var capabilities = []capabilityDef{
	{"bmad", true, "BMAD method framework installation"},
	{"config", true, "Viper config management scaffolding"},
	{"git", true, "Git initialization and initial commit"},
	{"docs", true, "mkdocs-material documentation scaffolding"},
	{"ui", false, "React/shadcn/Tailwind UI in ui/ subdirectory"},
	{"goreleaser", true, "GoReleaser configuration"},
	{"docker", true, "Dockerfile and .dockerignore"},
	{"release", true, "GitHub Actions release workflow"},
	{"mockery", true, "Mockery configuration"},
	{"editorconfig", true, "EditorConfig file"},
	{"makefile", true, "Makefile with common targets"},
}

var rootCmd = &cobra.Command{
	Use:   "gsi [project-name]",
	Short: "Initialize a Go project with best practices and tooling",
	Long: `gsi scaffolds a new Go project with cobra, viper,
mkdocs-material documentation, an embedded web UI, mockery, editorconfig,
and optional React/shadcn/Tailwind frontend.

Each capability can be toggled with --<name> / --no-<name> flags.
Defaults: most capabilities ON, ui OFF.

Examples:
  gsi my-awesome-app
  gsi --author "Jane Doe jane@example.com" my-app
  gsi --module github.com/myorg/myapp --dry-run my-app
  gsi --no-bmad --no-git my-app
  gsi --no-docker --no-release my-app
  gsi --only-docs my-app
  gsi --ui my-app
  gsi .    # Initialize in current directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("project name is required (use '.' for current directory)")
		}

		// Build capabilities map from defaults, then apply flag overrides
		caps := scaffold.DefaultCapabilities()
		for _, cap := range capabilities {
			noFlag := "no-" + cap.name
			// --no-<name> takes precedence if explicitly set
			if cmd.Flags().Changed(noFlag) {
				noVal, _ := cmd.Flags().GetBool(noFlag)
				caps[cap.name] = !noVal
			} else if cmd.Flags().Changed(cap.name) {
				val, _ := cmd.Flags().GetBool(cap.name)
				caps[cap.name] = val
			}
		}

		cfg := scaffold.Config{
			ProjectName:  args[0],
			Author:       viper.GetString("author"),
			GoModulePath: viper.GetString("module"),
			DryRun:       viper.GetBool("dry-run"),
			Verbose:      viper.GetBool("verbose"),
			OnlyDocs:     viper.GetBool("only-docs"),
			Capabilities: caps,
		}

		return scaffold.NewScaffolder(cfg).Run()
	},
}

var (
	buildVersion string
	buildCommit  string
	buildDate    string
)

// Execute is the CLI entry point called by main.
func Execute(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("author", "a", "Joe Scharf joe@joescharf.com", "Author name and email")
	rootCmd.Flags().StringP("module", "m", "", "Go module path (default: github.com/joescharf/<project>)")
	rootCmd.Flags().BoolP("dry-run", "d", false, "Show what would be done without executing")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().Bool("only-docs", false, "Only add docs scaffolding (skip everything else)")

	// Register capability flags: --<name> and hidden --no-<name>
	for _, cap := range capabilities {
		rootCmd.Flags().Bool(cap.name, cap.defaultValue, cap.description)
		rootCmd.Flags().Bool("no-"+cap.name, !cap.defaultValue, "Disable "+cap.description)
		_ = rootCmd.Flags().MarkHidden("no-" + cap.name)
	}

	// Bind non-capability flags to viper
	viper.BindPFlag("author", rootCmd.Flags().Lookup("author"))
	viper.BindPFlag("module", rootCmd.Flags().Lookup("module"))
	viper.BindPFlag("dry-run", rootCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	viper.BindPFlag("only-docs", rootCmd.Flags().Lookup("only-docs"))

	// Set defaults via viper
	viper.SetDefault("author", "Joe Scharf joe@joescharf.com")
}
