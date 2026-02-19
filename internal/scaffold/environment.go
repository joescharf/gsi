package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joescharf/gsi/internal/logger"
)

// CheckCommand returns true if the named command exists on PATH.
func CheckCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ValidateEnvironment checks that required tools are available, matching the shell
// script's validate_environment() logic.
func ValidateEnvironment(cfg *Config, log *logger.Logger) error {
	if cfg.OnlyDocs {
		log.Info("Validating environment (docs-only mode)...")
		if !CheckCommand("uv") {
			log.Error("uv is required for docs scaffolding but is not installed")
			log.Error("Install uv: https://docs.astral.sh/uv/")
			return fmt.Errorf("uv is required for docs scaffolding")
		}
		log.Success("Environment validation complete")
		return nil
	}

	log.Info("Validating environment...")

	// Required
	required := []string{"go"}
	var missingRequired []string
	for _, cmd := range required {
		if !CheckCommand(cmd) {
			log.Error(cmd + " is not installed or not in PATH")
			missingRequired = append(missingRequired, cmd)
		} else {
			log.VerboseMsg("Found " + cmd)
		}
	}

	if len(missingRequired) > 0 {
		return fmt.Errorf("missing required commands: %v", missingRequired)
	}

	// Optional — only check/warn if the relevant capability is enabled
	if cfg.IsEnabled(CapGit) {
		if !CheckCommand("git") {
			log.Warning("git is not installed — auto-disabling git capability")
			cfg.Disable(CapGit)
		} else {
			log.VerboseMsg("Found git")
		}
	}

	if cfg.IsEnabled(CapBmad) {
		if !CheckCommand("npx") {
			log.Warning("npx is not installed — auto-disabling bmad capability")
			cfg.Disable(CapBmad)
		} else {
			log.VerboseMsg("Found npx")
		}
	}

	if cfg.IsEnabled(CapUI) {
		if !CheckCommand("bun") {
			log.Warning("bun is not installed (optional)")
			// UI has a hard requirement validated separately in Run()
		} else {
			log.VerboseMsg("Found bun")
		}
	}

	if cfg.IsEnabled(CapDocs) {
		if !CheckCommand("uv") {
			log.Warning("uv is not installed — auto-disabling docs capability")
			cfg.Disable(CapDocs)
		} else {
			log.VerboseMsg("Found uv")
		}
	}

	log.Success("Environment validation complete")
	return nil
}

// CheckExistingState inspects the project directory and returns a list of existing artifacts.
func CheckExistingState(dir string, log *logger.Logger) []string {
	log.Info("Checking existing state...")

	checks := []struct {
		path  string
		label string
		isDir bool
	}{
		{"go.mod", "go.mod", false},
		{".git", ".git/", true},
		{"cmd", "cmd/", true},
		{"_bmad", "_bmad/", true},
		{"ui", "ui/", true},
		{"docs", "docs/", true},
	}

	var existing []string
	for _, c := range checks {
		full := filepath.Join(dir, c.path)
		info, err := os.Stat(full)
		if err != nil {
			continue
		}
		if c.isDir && info.IsDir() {
			existing = append(existing, c.label)
		} else if !c.isDir && !info.IsDir() {
			existing = append(existing, c.label)
		}
	}

	if len(existing) > 0 {
		log.Info("Found existing project files (will skip where appropriate):")
		for _, item := range existing {
			log.VerboseMsg("  - " + item)
		}
	} else {
		log.Success("No existing project files found")
	}

	return existing
}
