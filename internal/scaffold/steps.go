package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joescharf/gsi/internal/templates"
)

func (s *Scaffolder) templateData() templates.Data {
	// Derive owner from module path (e.g., "github.com/joescharf/myapp" -> "joescharf")
	owner := ""
	parts := strings.Split(s.Config.GoModulePath, "/")
	if len(parts) >= 2 {
		owner = parts[1]
	}

	return templates.Data{
		ProjectName:   s.Config.ProjectName,
		GoModulePath:  s.Config.GoModulePath,
		GoModuleOwner: owner,
	}
}

// stepInstallBmad installs the BMAD method framework via bun.
func (s *Scaffolder) stepInstallBmad() error {
	if !s.Config.IsEnabled(CapBmad) {
		s.Logger.Info("Skipping BMAD installation (--no-bmad)")
		return nil
	}

	bmadDir := filepath.Join(s.Config.ProjectDir, "_bmad")
	if _, err := os.Stat(bmadDir); err == nil {
		s.Logger.Info("_bmad/ directory already exists, skipping BMAD installation")
		return nil
	}

	if !CheckCommand("bun") {
		s.Logger.Warning("Skipping BMAD installation (bun not found)")
		return nil
	}

	return s.Executor.Execute("bunx bmad-method@alpha install", "Installing BMAD method framework")
}

// stepInstallCobraCli installs cobra-cli if not already on PATH.
func (s *Scaffolder) stepInstallCobraCli() error {
	if CheckCommand("cobra-cli") {
		s.Logger.Success("cobra-cli is already installed")
		return nil
	}
	return s.Executor.Execute("go install github.com/spf13/cobra-cli@latest", "Installing cobra-cli")
}

// stepGoModInit initializes the Go module.
func (s *Scaffolder) stepGoModInit() error {
	gomod := filepath.Join(s.Config.ProjectDir, "go.mod")
	if _, err := os.Stat(gomod); err == nil && !s.Config.DryRun {
		s.Logger.Info("go.mod already exists, skipping go mod init")
		return nil
	}
	return s.Executor.Execute(
		fmt.Sprintf("go mod init %s", s.Config.GoModulePath),
		"Initializing Go module",
	)
}

// stepCobraInit runs cobra-cli init to scaffold the CLI structure.
func (s *Scaffolder) stepCobraInit() error {
	cmdDir := filepath.Join(s.Config.ProjectDir, "cmd")
	if _, err := os.Stat(cmdDir); err == nil && !s.Config.DryRun {
		s.Logger.Info("cmd/ directory already exists, skipping cobra-cli init")
		return nil
	}
	return s.Executor.Execute(
		fmt.Sprintf(`cobra-cli init --viper --author "%s" --config $HOME/.config/%s`,
			s.Config.Author, s.Config.ProjectName),
		"Creating Cobra CLI application structure",
	)
}

// stepGenerateVersionCmd writes cmd/version.go from template with ldflags build vars.
func (s *Scaffolder) stepGenerateVersionCmd() error {
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "cmd", "version.go"),
		"cmd_version.go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateServeCmd writes cmd/serve.go from template.
func (s *Scaffolder) stepGenerateServeCmd() error {
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "cmd", "serve.go"),
		"cmd_serve.go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateConfigCmd writes cmd/config.go from template.
func (s *Scaffolder) stepGenerateConfigCmd() error {
	if !s.Config.IsEnabled(CapConfig) {
		s.Logger.Info("Skipping config command scaffolding (--no-config)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "cmd", "config.go"),
		"cmd_config.go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateConfigPkg writes internal/config/config.go from template.
func (s *Scaffolder) stepGenerateConfigPkg() error {
	if !s.Config.IsEnabled(CapConfig) {
		s.Logger.Info("Skipping config package scaffolding (--no-config)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "internal", "config", "config.go"),
		"config_go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateConfigInit writes cmd/config_init.go from template.
func (s *Scaffolder) stepGenerateConfigInit() error {
	if !s.Config.IsEnabled(CapConfig) {
		s.Logger.Info("Skipping config init scaffolding (--no-config)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "cmd", "config_init.go"),
		"cmd_config_init.go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateMockeryConfig writes .mockery.yml from template.
func (s *Scaffolder) stepGenerateMockeryConfig() error {
	if !s.Config.IsEnabled(CapMockery) {
		s.Logger.Info("Skipping mockery config (--no-mockery)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, ".mockery.yml"),
		"mockery_yml.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateEditorConfig writes .editorconfig from template.
func (s *Scaffolder) stepGenerateEditorConfig() error {
	if !s.Config.IsEnabled(CapEditorconfig) {
		s.Logger.Info("Skipping editorconfig (--no-editorconfig)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, ".editorconfig"),
		"editorconfig.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateUIPlaceholder writes internal/ui/dist/index.html from template.
func (s *Scaffolder) stepGenerateUIPlaceholder() error {
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "internal", "ui", "dist", "index.html"),
		"index_html.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateEmbedGo writes internal/ui/embed.go from template.
func (s *Scaffolder) stepGenerateEmbedGo() error {
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "internal", "ui", "embed.go"),
		"embed_go.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGoModTidy runs go mod tidy.
func (s *Scaffolder) stepGoModTidy() error {
	return s.Executor.Execute("go mod tidy", "Tidying Go dependencies")
}

// stepGenerateMakefile writes the Makefile from template.
func (s *Scaffolder) stepGenerateMakefile() error {
	if !s.Config.IsEnabled(CapMakefile) {
		s.Logger.Info("Skipping Makefile (--no-makefile)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "Makefile"),
		"makefile.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateGoreleaser writes .goreleaser.yml from template.
func (s *Scaffolder) stepGenerateGoreleaser() error {
	if !s.Config.IsEnabled(CapGoreleaser) {
		s.Logger.Info("Skipping goreleaser config (--no-goreleaser)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, ".goreleaser.yml"),
		"goreleaser_yml.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateDockerfile writes Dockerfile from template.
func (s *Scaffolder) stepGenerateDockerfile() error {
	if !s.Config.IsEnabled(CapDocker) {
		s.Logger.Info("Skipping Dockerfile (--no-docker)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, "Dockerfile"),
		"dockerfile.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateReleaseWorkflow writes .github/workflows/release.yml from template.
func (s *Scaffolder) stepGenerateReleaseWorkflow() error {
	if !s.Config.IsEnabled(CapRelease) {
		s.Logger.Info("Skipping release workflow (--no-release)")
		return nil
	}
	dir := filepath.Join(s.Config.ProjectDir, ".github", "workflows")
	if !s.Config.DryRun {
		os.MkdirAll(dir, 0o755)
	}
	return WriteTemplateFile(
		filepath.Join(dir, "release.yml"),
		"github_release_yml.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateDockerignore writes .dockerignore from template.
func (s *Scaffolder) stepGenerateDockerignore() error {
	if !s.Config.IsEnabled(CapDocker) {
		s.Logger.Info("Skipping .dockerignore (--no-docker)")
		return nil
	}
	return WriteTemplateFile(
		filepath.Join(s.Config.ProjectDir, ".dockerignore"),
		"dockerignore.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepInitDocs scaffolds mkdocs-material documentation.
func (s *Scaffolder) stepInitDocs() error {
	if !s.Config.IsEnabled(CapDocs) {
		s.Logger.Info("Skipping docs scaffolding (--no-docs)")
		return nil
	}

	// Auto-skip if uv is missing
	if !CheckCommand("uv") {
		s.Logger.Warning("uv is not installed, skipping docs scaffolding (install: https://docs.astral.sh/uv/)")
		s.Config.Disable(CapDocs)
		return nil
	}

	dir := s.Config.ProjectDir
	data := s.templateData()

	// Initialize uv project in docs/
	pyproject := filepath.Join(dir, "docs", "pyproject.toml")
	if _, err := os.Stat(pyproject); os.IsNotExist(err) || s.Config.DryRun {
		if err := s.Executor.Execute(
			fmt.Sprintf("uv init --name %s-docs docs", s.Config.ProjectName),
			"Initializing uv project in docs/",
		); err != nil {
			return err
		}

		// Remove uv init scaffolding
		if !s.Config.DryRun {
			for _, f := range []string{
				filepath.Join(dir, "docs", ".git"),
				filepath.Join(dir, "docs", "hello.py"),
				filepath.Join(dir, "docs", "main.py"),
				filepath.Join(dir, "docs", "README.md"),
			} {
				os.RemoveAll(f)
			}
		} else {
			s.Logger.Warning("[DRY-RUN] Would remove uv init scaffolding (docs/.git, docs/hello.py, docs/main.py, docs/README.md)")
		}
	} else {
		s.Logger.Info("docs/pyproject.toml already exists, skipping uv init")
	}

	// Add mkdocs-material dependencies
	if err := s.addDocsDeps(pyproject); err != nil {
		return err
	}

	// Write mkdocs.yml
	if err := WriteTemplateFile(
		filepath.Join(dir, "docs", "mkdocs.yml"),
		"mkdocs_yml.tmpl", data, s.Config.DryRun, s.Logger,
	); err != nil {
		return err
	}

	// Write docs/.gitignore
	if err := WriteTemplateFile(
		filepath.Join(dir, "docs", ".gitignore"),
		"docs_gitignore.tmpl", data, s.Config.DryRun, s.Logger,
	); err != nil {
		return err
	}

	// Create docs/docs/stylesheets directory
	if !s.Config.DryRun {
		os.MkdirAll(filepath.Join(dir, "docs", "docs", "stylesheets"), 0o755)
	} else {
		s.Logger.Warning("[DRY-RUN] Would create docs/docs/stylesheets/")
	}

	// Write docs/docs/index.md
	if err := WriteTemplateFile(
		filepath.Join(dir, "docs", "docs", "index.md"),
		"docs_index_md.tmpl", data, s.Config.DryRun, s.Logger,
	); err != nil {
		return err
	}

	// Write docs/docs/getting-started.md
	if err := WriteTemplateFile(
		filepath.Join(dir, "docs", "docs", "getting-started.md"),
		"docs_getting_started_md.tmpl", data, s.Config.DryRun, s.Logger,
	); err != nil {
		return err
	}

	// Write docs/docs/stylesheets/extra.css
	return WriteTemplateFile(
		filepath.Join(dir, "docs", "docs", "stylesheets", "extra.css"),
		"docs_extra_css.tmpl", data, s.Config.DryRun, s.Logger,
	)
}

// addDocsDeps adds mkdocs-material to the docs pyproject.toml if not already present.
func (s *Scaffolder) addDocsDeps(pyproject string) error {
	// Check if already has mkdocs-material
	if !s.Config.DryRun {
		content, err := os.ReadFile(pyproject)
		if err == nil {
			if strings.Contains(string(content), "mkdocs-material") {
				s.Logger.Info("mkdocs-material already in docs/pyproject.toml, skipping")
				return nil
			}
		}
	}

	return s.Executor.Execute(
		"cd docs && uv add mkdocs-material 'mkdocs-git-revision-date-localized-plugin>=1.4' && cd ..",
		"Adding mkdocs-material dependencies",
	)
}

// stepInitUI initializes a React/shadcn UI in ui/.
func (s *Scaffolder) stepInitUI() error {
	if !s.Config.IsEnabled(CapUI) {
		return nil
	}

	uiDir := filepath.Join(s.Config.ProjectDir, "ui")
	if _, err := os.Stat(uiDir); err == nil && !s.Config.DryRun {
		s.Logger.Info("ui/ directory already exists, skipping UI initialization")
		return nil
	}

	return s.Executor.Execute("bun init --react=shadcn ui", "Initializing React/shadcn/Tailwind UI in ui/")
}

// stepInitGit initializes git, creates .gitignore, and makes an initial commit.
func (s *Scaffolder) stepInitGit() error {
	if !s.Config.IsEnabled(CapGit) {
		s.Logger.Info("Skipping git initialization (--no-git)")
		return nil
	}

	dir := s.Config.ProjectDir

	// git init
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) || s.Config.DryRun {
		if err := s.Executor.Execute("git init", "Initializing git repository"); err != nil {
			return err
		}
	} else {
		s.Logger.Info(".git directory already exists, skipping git init")
	}

	// .gitignore
	if err := WriteTemplateFile(
		filepath.Join(dir, ".gitignore"),
		"gitignore.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	); err != nil {
		return err
	}

	// Initial commit
	if s.Config.DryRun {
		s.Logger.Warning("[DRY-RUN] Would create initial commit")
		return nil
	}

	// Check if repo already has commits
	if s.Executor.RunCommandQuiet("git", "rev-parse", "HEAD") == nil {
		s.Logger.Info("Git repository already has commits, skipping initial commit")
		return nil
	}

	if err := s.Executor.Execute("git add .", "Staging files for initial commit"); err != nil {
		return err
	}
	return s.Executor.Execute("git commit -m 'initial commit'", "Creating initial commit")
}

// stepPrintSummary prints the "Next steps" summary.
func (s *Scaffolder) stepPrintSummary() {
	s.Logger.Plain("")
	s.Logger.Success("Project initialization complete!")
	s.Logger.Plain("")
	s.Logger.Info("Next steps:")

	step := 1
	if s.Config.OnlyDocs {
		s.Logger.Plain(fmt.Sprintf("  %d. Run 'make docs-serve' to start the docs dev server", step))
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Edit docs in docs/docs/", step))
	} else {
		s.Logger.Plain(fmt.Sprintf("  %d. Review the generated code in cmd/", step))
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Update the project description in cmd/root.go", step))
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Run 'make build' to build your application", step))
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Run 'make run' or './bin/%s --help' to see available commands", step, s.Config.ProjectName))
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Run 'make serve' to start the embedded web UI server", step))

		if s.Config.IsEnabled(CapDocs) {
			step++
			s.Logger.Plain(fmt.Sprintf("  %d. Run 'make docs-serve' to start the docs dev server", step))
		}
		if s.Config.IsEnabled(CapUI) {
			step++
			s.Logger.Plain(fmt.Sprintf("  %d. Run 'make ui-dev' to start the React dev server", step))
		}
		step++
		s.Logger.Plain(fmt.Sprintf("  %d. Run 'make help' to see all available targets", step))
	}
	s.Logger.Plain("")
}
