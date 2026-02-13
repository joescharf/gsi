package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testScaffolder(t *testing.T, dryRun bool) (*Scaffolder, *strings.Builder, *strings.Builder) {
	t.Helper()
	dir := t.TempDir()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}

	cfg := Config{
		ProjectName:  "testproj",
		Author:       "Test Author test@example.com",
		GoModulePath: "github.com/example/testproj",
		DryRun:       dryRun,
		Verbose:      true,
		ProjectDir:   dir,
		Capabilities: DefaultCapabilities(),
	}

	s := NewScaffolder(cfg)
	// Override logger writers to capture output in tests
	s.Logger.Stdout = stdout
	s.Logger.Stderr = stderr
	s.Executor.Logger = s.Logger
	return s, stdout, stderr
}

func TestStepInstallBmadSkipFlag(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapBmad] = false

	if err := s.stepInstallBmad(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping BMAD") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepInstallBmadExistingDir(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	os.MkdirAll(filepath.Join(s.Config.ProjectDir, "_bmad"), 0o755)

	if err := s.stepInstallBmad(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "already exists") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateServeCmdIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	// First call should create
	if err := s.stepGenerateServeCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "serve.go")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected serve.go to be created")
	}

	// Second call should skip
	if err := s.stepGenerateServeCmd(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateServeCmdDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateServeCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "serve.go")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

func TestStepInitGitSkipFlag(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapGit] = false

	if err := s.stepInitGit(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping git") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepPrintSummaryOnlyDocs(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.OnlyDocs = true
	s.stepPrintSummary()
	if !strings.Contains(stdout.String(), "docs-serve") {
		t.Errorf("expected docs-serve in summary, got %q", stdout.String())
	}
}

func TestStepGenerateGoreleaserIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateGoreleaser(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, ".goreleaser.yml")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected .goreleaser.yml to be created")
	}

	// Second call should skip
	if err := s.stepGenerateGoreleaser(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateGoreleaserDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateGoreleaser(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, ".goreleaser.yml")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

func TestStepGenerateDockerfileIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateDockerfile(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "Dockerfile")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected Dockerfile to be created")
	}

	// Second call should skip
	if err := s.stepGenerateDockerfile(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateDockerfileDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateDockerfile(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "Dockerfile")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

func TestStepGenerateDockerignoreIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateDockerignore(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, ".dockerignore")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected .dockerignore to be created")
	}

	// Second call should skip
	if err := s.stepGenerateDockerignore(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateDockerignoreDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateDockerignore(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, ".dockerignore")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

func TestStepPrintSummaryFull(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.stepPrintSummary()
	if !strings.Contains(stdout.String(), "make build") {
		t.Errorf("expected 'make build' in summary, got %q", stdout.String())
	}
}

// --- Capability guard tests ---

func TestStepGenerateGoreleaserDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapGoreleaser] = false

	if err := s.stepGenerateGoreleaser(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping goreleaser") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateDockerfileDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapDocker] = false

	if err := s.stepGenerateDockerfile(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping Dockerfile") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateDockerignoreDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapDocker] = false

	if err := s.stepGenerateDockerignore(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping .dockerignore") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateReleaseWorkflowDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapRelease] = false

	if err := s.stepGenerateReleaseWorkflow(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping release") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateMockeryDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapMockery] = false

	if err := s.stepGenerateMockeryConfig(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping mockery") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateEditorConfigDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapEditorconfig] = false

	if err := s.stepGenerateEditorConfig(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping editorconfig") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateMakefileDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapMakefile] = false

	if err := s.stepGenerateMakefile(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping Makefile") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepInitDocsDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapDocs] = false

	if err := s.stepInitDocs(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping docs") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

// --- Version command template tests ---

func TestStepGenerateVersionCmdIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateVersionCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "version.go")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected version.go to be created")
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "testproj") {
		t.Error("expected project name in version.go")
	}
	if !strings.Contains(string(content), "buildVersion") {
		t.Error("expected buildVersion reference in version.go")
	}

	// Second call should skip
	if err := s.stepGenerateVersionCmd(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateVersionCmdDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateVersionCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "version.go")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

// --- Config scaffold step tests ---

func TestStepGenerateConfigCmdIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateConfigCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "config.go")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected config.go to be created")
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "config init") {
		t.Error("expected 'config init' subcommand in config.go")
	}

	// Second call should skip
	if err := s.stepGenerateConfigCmd(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateConfigCmdDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapConfig] = false

	if err := s.stepGenerateConfigCmd(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping config command") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateConfigCmdDryRun(t *testing.T) {
	s, _, stderr := testScaffolder(t, true)

	if err := s.stepGenerateConfigCmd(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "config.go")
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not exist in dry-run mode")
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}
}

func TestStepGenerateConfigPkgIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateConfigPkg(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "internal", "config", "config.go")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected internal/config/config.go to be created")
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "SetDefaults") {
		t.Error("expected SetDefaults function in config.go")
	}
	if !strings.Contains(string(content), "viper") {
		t.Error("expected viper usage in config.go")
	}

	// Second call should skip
	if err := s.stepGenerateConfigPkg(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateConfigPkgDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapConfig] = false

	if err := s.stepGenerateConfigPkg(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping config package") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}

func TestStepGenerateConfigInitIdempotent(t *testing.T) {
	s, _, _ := testScaffolder(t, false)

	if err := s.stepGenerateConfigInit(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Config.ProjectDir, "cmd", "config_init.go")
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected cmd/config_init.go to be created")
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "initConfig") {
		t.Error("expected initConfig function in config_init.go")
	}

	// Second call should skip
	if err := s.stepGenerateConfigInit(); err != nil {
		t.Fatal(err)
	}
}

func TestStepGenerateConfigInitDisabled(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.Config.Capabilities[CapConfig] = false

	if err := s.stepGenerateConfigInit(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Skipping config init") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}
}
