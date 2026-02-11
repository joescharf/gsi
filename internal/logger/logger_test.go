package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	var out bytes.Buffer
	l := &Logger{Stdout: &out, Stderr: &bytes.Buffer{}}
	l.Info("hello")
	if !strings.Contains(out.String(), "hello") {
		t.Errorf("expected output to contain 'hello', got %q", out.String())
	}
	if !strings.Contains(out.String(), "ℹ") {
		t.Errorf("expected info icon, got %q", out.String())
	}
}

func TestSuccess(t *testing.T) {
	var out bytes.Buffer
	l := &Logger{Stdout: &out, Stderr: &bytes.Buffer{}}
	l.Success("done")
	if !strings.Contains(out.String(), "done") {
		t.Errorf("expected output to contain 'done', got %q", out.String())
	}
	if !strings.Contains(out.String(), "✓") {
		t.Errorf("expected success icon, got %q", out.String())
	}
}

func TestWarningWritesToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	l := &Logger{Stdout: &stdout, Stderr: &stderr}
	l.Warning("watch out")
	if stdout.Len() != 0 {
		t.Errorf("expected nothing on stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "watch out") {
		t.Errorf("expected stderr to contain 'watch out', got %q", stderr.String())
	}
}

func TestErrorWritesToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	l := &Logger{Stdout: &stdout, Stderr: &stderr}
	l.Error("bad")
	if stdout.Len() != 0 {
		t.Errorf("expected nothing on stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "bad") {
		t.Errorf("expected stderr to contain 'bad', got %q", stderr.String())
	}
}

func TestVerboseSuppressed(t *testing.T) {
	var out bytes.Buffer
	l := &Logger{Verbose: false, Stdout: &out, Stderr: &bytes.Buffer{}}
	l.VerboseMsg("secret")
	if out.Len() != 0 {
		t.Errorf("expected no output when verbose is false, got %q", out.String())
	}
}

func TestVerboseEnabled(t *testing.T) {
	var out bytes.Buffer
	l := &Logger{Verbose: true, Stdout: &out, Stderr: &bytes.Buffer{}}
	l.VerboseMsg("detail")
	if !strings.Contains(out.String(), "detail") {
		t.Errorf("expected output to contain 'detail', got %q", out.String())
	}
}
