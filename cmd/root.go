package cmd

import (
	"fmt"
	"os"

	"github.com/joescharf/gsi/internal/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gsi [project-name]",
	Short: "Initialize a Go project with best practices and tooling",
	Long: `gsi scaffolds a new Go project with cobra, viper,
mkdocs-material documentation, an embedded web UI, mockery, editorconfig,
and optional React/shadcn/Tailwind frontend.

Examples:
  gsi my-awesome-app
  gsi --author "Jane Doe jane@example.com" my-app
  gsi --module github.com/myorg/myapp --dry-run my-app
  gsi --skip-bmad --skip-git my-app
  gsi --only-docs my-app
  gsi --ui my-app
  gsi .    # Initialize in current directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("project name is required (use '.' for current directory)")
		}

		cfg := scaffold.Config{
			ProjectName:  args[0],
			Author:       viper.GetString("author"),
			GoModulePath: viper.GetString("module"),
			DryRun:       viper.GetBool("dry-run"),
			Verbose:      viper.GetBool("verbose"),
			SkipBmad:     viper.GetBool("skip-bmad"),
			SkipGit:      viper.GetBool("skip-git"),
			SkipDocs:     viper.GetBool("skip-docs"),
			OnlyDocs:     viper.GetBool("only-docs"),
			InitUI:       viper.GetBool("ui"),
		}

		return scaffold.NewScaffolder(cfg).Run()
	},
}

// Execute is the CLI entry point called by main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("author", "a", "Joe Scharf joe@joescharf.com", "Author name and email")
	rootCmd.Flags().StringP("module", "m", "", "Go module path (default: github.com/joescharf/<project>)")
	rootCmd.Flags().BoolP("dry-run", "d", false, "Show what would be done without executing")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().Bool("skip-bmad", false, "Skip BMAD method installation")
	rootCmd.Flags().Bool("skip-git", false, "Skip git initialization and commit")
	rootCmd.Flags().Bool("skip-docs", false, "Skip mkdocs-material documentation scaffolding")
	rootCmd.Flags().Bool("only-docs", false, "Only add docs scaffolding (skip everything else)")
	rootCmd.Flags().Bool("ui", false, "Initialize a React/shadcn/Tailwind UI in ui/ subdirectory")

	// Bind all flags to viper
	viper.BindPFlag("author", rootCmd.Flags().Lookup("author"))
	viper.BindPFlag("module", rootCmd.Flags().Lookup("module"))
	viper.BindPFlag("dry-run", rootCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	viper.BindPFlag("skip-bmad", rootCmd.Flags().Lookup("skip-bmad"))
	viper.BindPFlag("skip-git", rootCmd.Flags().Lookup("skip-git"))
	viper.BindPFlag("skip-docs", rootCmd.Flags().Lookup("skip-docs"))
	viper.BindPFlag("only-docs", rootCmd.Flags().Lookup("only-docs"))
	viper.BindPFlag("ui", rootCmd.Flags().Lookup("ui"))

	// Set defaults via viper
	viper.SetDefault("author", "Joe Scharf joe@joescharf.com")
}
