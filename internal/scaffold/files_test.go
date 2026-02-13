package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joescharf/gsi/internal/templates"
)

func TestWriteTemplateFileCreatesFile(t *testing.T) {
	log, _, _ := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	data := templates.Data{ProjectName: "myapp", GoModulePath: "github.com/example/myapp"}
	err := WriteTemplateFile(path, "editorconfig.tmpl", data, false, log)
	if err != nil {
		t.Fatalf("WriteTemplateFile failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if !strings.Contains(string(content), "root = true") {
		t.Errorf("expected editorconfig content, got %q", string(content))
	}
}

func TestWriteTemplateFileIdempotent(t *testing.T) {
	log, stdout, _ := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	// Write original content
	if err := os.WriteFile(path, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}

	data := templates.Data{ProjectName: "myapp"}
	err := WriteTemplateFile(path, "editorconfig.tmpl", data, false, log)
	if err != nil {
		t.Fatalf("WriteTemplateFile failed: %v", err)
	}

	// Should have skipped
	if !strings.Contains(stdout.String(), "already exists") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}

	// Content should be unchanged
	content, _ := os.ReadFile(path)
	if string(content) != "original" {
		t.Errorf("file should not have been overwritten, got %q", string(content))
	}
}

func TestWriteTemplateFileDryRun(t *testing.T) {
	log, _, stderr := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "test.txt")

	data := templates.Data{ProjectName: "myapp"}
	err := WriteTemplateFile(path, "editorconfig.tmpl", data, true, log)
	if err != nil {
		t.Fatalf("WriteTemplateFile dry-run failed: %v", err)
	}

	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}

	// File should not exist
	if _, err := os.Stat(path); err == nil {
		t.Error("expected file to not exist in dry-run mode")
	}
}

func TestWriteExecutableTemplateFileCreatesFile(t *testing.T) {
	log, _, _ := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sh")

	data := templates.Data{ProjectName: "myapp", GoModulePath: "github.com/example/myapp"}
	err := WriteExecutableTemplateFile(path, "editorconfig.tmpl", data, false, log)
	if err != nil {
		t.Fatalf("WriteExecutableTemplateFile failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat written file: %v", err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Errorf("expected 0755 permissions, got %o", info.Mode().Perm())
	}
}

func TestWriteExecutableTemplateFileIdempotent(t *testing.T) {
	log, stdout, _ := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sh")

	// Write original content
	if err := os.WriteFile(path, []byte("original"), 0o755); err != nil {
		t.Fatal(err)
	}

	data := templates.Data{ProjectName: "myapp"}
	err := WriteExecutableTemplateFile(path, "editorconfig.tmpl", data, false, log)
	if err != nil {
		t.Fatalf("WriteExecutableTemplateFile failed: %v", err)
	}

	// Should have skipped
	if !strings.Contains(stdout.String(), "already exists") {
		t.Errorf("expected skip message, got %q", stdout.String())
	}

	// Content should be unchanged
	content, _ := os.ReadFile(path)
	if string(content) != "original" {
		t.Errorf("file should not have been overwritten, got %q", string(content))
	}
}

func TestWriteExecutableTemplateFileDryRun(t *testing.T) {
	log, _, stderr := testLogger()
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "test.sh")

	data := templates.Data{ProjectName: "myapp"}
	err := WriteExecutableTemplateFile(path, "editorconfig.tmpl", data, true, log)
	if err != nil {
		t.Fatalf("WriteExecutableTemplateFile dry-run failed: %v", err)
	}

	if !strings.Contains(stderr.String(), "[DRY-RUN]") {
		t.Errorf("expected dry-run message, got %q", stderr.String())
	}

	// File should not exist
	if _, err := os.Stat(path); err == nil {
		t.Error("expected file to not exist in dry-run mode")
	}
}
