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

	// Optional
	for _, cmd := range []string{"git", "bun", "uv"} {
		if !CheckCommand(cmd) {
			log.Warning(cmd + " is not installed (optional)")
		} else {
			log.VerboseMsg("Found " + cmd)
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
