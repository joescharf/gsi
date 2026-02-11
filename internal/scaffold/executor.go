package scaffold

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joescharf/gsi/internal/logger"
)

// Executor runs shell commands with dry-run support.
type Executor struct {
	DryRun bool
	Logger *logger.Logger
	Dir    string // working directory for commands
}

// Execute runs a shell command via sh -c. It mirrors the shell script's execute() function.
func (e *Executor) Execute(command, description string) error {
	e.Logger.Info(description)
	e.Logger.VerboseMsg("Command: " + command)

	if e.DryRun {
		e.Logger.Warning(fmt.Sprintf("[DRY-RUN] Would execute: %s", command))
		return nil
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = e.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		e.Logger.Error(description + " - Failed")
		return fmt.Errorf("%s: %w", description, err)
	}

	e.Logger.Success(description + " - Done")
	return nil
}

// RunCommand runs a command directly (not via shell).
func (e *Executor) RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = e.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandQuiet runs a command suppressing all output. Used for existence checks.
func (e *Executor) RunCommandQuiet(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = e.Dir
	return cmd.Run()
}
