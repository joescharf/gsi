package scaffold

// Config holds all CLI flags and derived values for a scaffold run.
type Config struct {
	ProjectName  string
	Author       string
	GoModulePath string
	DryRun       bool
	Verbose      bool
	SkipBmad     bool
	SkipGit      bool
	SkipDocs     bool
	OnlyDocs     bool
	InitUI       bool

	// Derived — set during validation
	ProjectDir string
}
