package scaffold

import (
	"bytes"
	"strings"
	"testing"

	"github.com/joescharf/gsi/internal/logger"
)

func testLogger() (*logger.Logger, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	return &logger.Logger{Verbose: true, Stdout: stdout, Stderr: stderr}, stdout, stderr
}

func TestExecuteDryRun(t *testing.T) {
	log, _, stderr := testLogger()
	exec := &Executor{DryRun: true, Logger: log, Dir: t.TempDir()}

	err := exec.Execute("echo hello", "Print hello")
	if err != nil {
		t.Fatalf("dry-run execute should not error: %v", err)
	}
	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message in stderr, got %q", stderr.String())
	}
}

func TestExecuteRealCommand(t *testing.T) {
	log, _, _ := testLogger()
	exec := &Executor{DryRun: false, Logger: log, Dir: t.TempDir()}

	err := exec.Execute("true", "Run true")
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestExecuteFailingCommand(t *testing.T) {
	log, _, _ := testLogger()
	exec := &Executor{DryRun: false, Logger: log, Dir: t.TempDir()}

	err := exec.Execute("false", "Run false")
	if err == nil {
		t.Fatal("expected error for failing command")
	}
}
