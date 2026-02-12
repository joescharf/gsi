package scaffold

import "testing"

func TestDefaultCapabilities(t *testing.T) {
	caps := DefaultCapabilities()

	// All expected capabilities should be present
	expected := []string{
		CapBmad, CapConfig, CapGit, CapDocs, CapUI,
		CapGoreleaser, CapDocker, CapRelease,
		CapMockery, CapEditorconfig, CapMakefile,
	}
	for _, name := range expected {
		if _, ok := caps[name]; !ok {
			t.Errorf("expected capability %q in defaults", name)
		}
	}

	// UI should default to OFF, others ON
	if caps[CapUI] != false {
		t.Error("expected CapUI to default to false")
	}
	for _, name := range []string{CapBmad, CapConfig, CapGit, CapDocs, CapGoreleaser, CapDocker, CapRelease, CapMockery, CapEditorconfig, CapMakefile} {
		if caps[name] != true {
			t.Errorf("expected %q to default to true", name)
		}
	}
}

func TestIsEnabled(t *testing.T) {
	cfg := Config{
		Capabilities: map[string]bool{
			"foo": true,
			"bar": false,
		},
	}

	if !cfg.IsEnabled("foo") {
		t.Error("expected foo to be enabled")
	}
	if cfg.IsEnabled("bar") {
		t.Error("expected bar to be disabled")
	}
	if cfg.IsEnabled("nonexistent") {
		t.Error("expected nonexistent capability to be disabled")
	}
}

func TestDisable(t *testing.T) {
	cfg := Config{
		Capabilities: map[string]bool{
			"foo": true,
		},
	}

	if !cfg.IsEnabled("foo") {
		t.Fatal("precondition: foo should be enabled")
	}

	cfg.Disable("foo")

	if cfg.IsEnabled("foo") {
		t.Error("expected foo to be disabled after Disable()")
	}
}

func TestDisableNonexistent(t *testing.T) {
	cfg := Config{
		Capabilities: map[string]bool{},
	}

	// Should not panic
	cfg.Disable("nonexistent")

	if cfg.IsEnabled("nonexistent") {
		t.Error("expected nonexistent to be disabled")
	}
}
