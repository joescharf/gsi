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
	s.Config.SkipBmad = true

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
	s.Config.SkipGit = true

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

func TestStepPrintSummaryFull(t *testing.T) {
	s, stdout, _ := testScaffolder(t, false)
	s.stepPrintSummary()
	if !strings.Contains(stdout.String(), "make build") {
		t.Errorf("expected 'make build' in summary, got %q", stdout.String())
	}
}
