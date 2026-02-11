package logger

import (
	"fmt"
	"io"
	"os"
)

const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[0;34m"
	colorReset  = "\033[0m"
)

// Logger provides colored, leveled output matching the shell script's style.
type Logger struct {
	Verbose bool
	Stdout  io.Writer
	Stderr  io.Writer
}

// New returns a Logger that writes to os.Stdout and os.Stderr.
func New(verbose bool) *Logger {
	return &Logger{
		Verbose: verbose,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

func (l *Logger) Info(msg string) {
	fmt.Fprintf(l.Stdout, "%sℹ%s %s\n", colorBlue, colorReset, msg)
}

func (l *Logger) Success(msg string) {
	fmt.Fprintf(l.Stdout, "%s✓%s %s\n", colorGreen, colorReset, msg)
}

func (l *Logger) Warning(msg string) {
	fmt.Fprintf(l.Stderr, "%s⚠%s %s\n", colorYellow, colorReset, msg)
}

func (l *Logger) Error(msg string) {
	fmt.Fprintf(l.Stderr, "%s✗%s %s\n", colorRed, colorReset, msg)
}

func (l *Logger) VerboseMsg(msg string) {
	if l.Verbose {
		fmt.Fprintf(l.Stdout, "%s  →%s %s\n", colorBlue, colorReset, msg)
	}
}

// Plain prints a line without any icon prefix (for config display, etc.).
func (l *Logger) Plain(msg string) {
	fmt.Fprintln(l.Stdout, msg)
}
