package scaffold

// Capability name constants.
const (
	CapBmad         = "bmad"
	CapConfig       = "config"
	CapGit          = "git"
	CapDocs         = "docs"
	CapUI           = "ui"
	CapGoreleaser   = "goreleaser"
	CapDocker       = "docker"
	CapRelease      = "release"
	CapMockery      = "mockery"
	CapEditorconfig = "editorconfig"
	CapMakefile     = "makefile"
)

// DefaultCapabilities returns the default enabled/disabled state for each capability.
func DefaultCapabilities() map[string]bool {
	return map[string]bool{
		CapBmad:         true,
		CapConfig:       true,
		CapGit:          true,
		CapDocs:         true,
		CapUI:           false,
		CapGoreleaser:   true,
		CapDocker:       true,
		CapRelease:      true,
		CapMockery:      true,
		CapEditorconfig: true,
		CapMakefile:     true,
	}
}

// Config holds all CLI flags and derived values for a scaffold run.
type Config struct {
	ProjectName  string
	Author       string
	GoModulePath string
	DryRun       bool
	Verbose      bool
	OnlyDocs     bool
	Capabilities map[string]bool

	// Derived — set during validation
	ProjectDir string
}

// IsEnabled returns whether the named capability is enabled.
func (c *Config) IsEnabled(name string) bool {
	enabled, ok := c.Capabilities[name]
	return ok && enabled
}

// Disable turns off a capability at runtime (e.g., when a soft dependency is missing).
func (c *Config) Disable(name string) {
	c.Capabilities[name] = false
}
