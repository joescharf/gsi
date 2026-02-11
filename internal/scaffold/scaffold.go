package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/joescharf/gsi/internal/logger"
)

// Scaffolder holds all state needed to run the scaffold steps.
type Scaffolder struct {
	Config   Config
	Logger   *logger.Logger
	Executor *Executor
}

// NewScaffolder creates a Scaffolder from the given Config.
func NewScaffolder(cfg Config) *Scaffolder {
	log := logger.New(cfg.Verbose)
	return &Scaffolder{
		Config: cfg,
		Logger: log,
		Executor: &Executor{
			DryRun: cfg.DryRun,
			Logger: log,
			Dir:    cfg.ProjectDir,
		},
	}
}

var validProjectName = regexp.MustCompile(`^[a-zA-Z0-9_/.\-]+$`)

// Run is the main orchestrator that sequences all scaffold steps.
func (s *Scaffolder) Run() error {
	cfg := &s.Config

	// Resolve project name and directory
	if cfg.ProjectName == "." || cfg.ProjectName == "./" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}
		cfg.ProjectDir = dir
		cfg.ProjectName = filepath.Base(dir)
		s.Logger.Info("Initializing in current directory")
		s.Logger.VerboseMsg("Project directory: " + cfg.ProjectDir)
	} else {
		if !validProjectName.MatchString(cfg.ProjectName) {
			return fmt.Errorf("invalid project name: must contain only letters, numbers, hyphens, underscores, dots, and slashes")
		}

		if filepath.IsAbs(cfg.ProjectName) {
			cfg.ProjectDir = cfg.ProjectName
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting current directory: %w", err)
			}
			cfg.ProjectDir = filepath.Join(cwd, cfg.ProjectName)
		}
		cfg.ProjectName = filepath.Base(cfg.ProjectDir)

		// Create or reuse directory
		info, err := os.Stat(cfg.ProjectDir)
		if err == nil && info.IsDir() {
			s.Logger.Info(fmt.Sprintf("Directory '%s' already exists, continuing with initialization", cfg.ProjectDir))
		} else if os.IsNotExist(err) {
			if !cfg.DryRun {
				s.Logger.Info("Creating project directory: " + cfg.ProjectDir)
				if err := os.MkdirAll(cfg.ProjectDir, 0o755); err != nil {
					return fmt.Errorf("creating directory: %w", err)
				}
				s.Logger.Success("Created project directory")
			} else {
				s.Logger.Warning("[DRY-RUN] Would create directory: " + cfg.ProjectDir)
			}
		} else if err != nil {
			return fmt.Errorf("checking directory: %w", err)
		}
	}

	// Update executor working directory
	s.Executor.Dir = cfg.ProjectDir

	// Set defaults for module path
	if cfg.GoModulePath == "" {
		cfg.GoModulePath = "github.com/joescharf/" + cfg.ProjectName
	}

	// Display configuration
	s.Logger.Plain("")
	s.Logger.Info("Configuration:")
	s.Logger.Plain("  Project Name:  " + cfg.ProjectName)
	s.Logger.Plain("  Project Dir:   " + cfg.ProjectDir)
	s.Logger.Plain("  Module Path:   " + cfg.GoModulePath)
	s.Logger.Plain("  Author:        " + cfg.Author)
	home, _ := os.UserHomeDir()
	s.Logger.Plain("  Config Dir:    " + filepath.Join(home, ".config", cfg.ProjectName))
	s.Logger.Plain(fmt.Sprintf("  Init UI:       %v", cfg.InitUI))
	s.Logger.Plain(fmt.Sprintf("  Skip Docs:     %v", cfg.SkipDocs))
	s.Logger.Plain(fmt.Sprintf("  Only Docs:     %v", cfg.OnlyDocs))
	if cfg.DryRun {
		s.Logger.Plain("  \033[1;33mMode:          DRY-RUN\033[0m")
	}
	s.Logger.Plain("")

	// Validate mutually exclusive flags
	if cfg.OnlyDocs && cfg.SkipDocs {
		return fmt.Errorf("--only-docs and --skip-docs are mutually exclusive")
	}

	// Validate environment
	if err := ValidateEnvironment(cfg, s.Logger); err != nil {
		return err
	}

	// Validate bun if --ui
	if cfg.InitUI && !cfg.OnlyDocs {
		if !CheckCommand("bun") {
			s.Logger.Error("bun is required for UI initialization (--ui flag) but is not installed")
			s.Logger.Error("Install bun: https://bun.sh")
			return fmt.Errorf("bun is required for --ui flag")
		}
	}

	// Auto-skip docs if uv is missing (non-only-docs mode)
	if !cfg.SkipDocs && !cfg.OnlyDocs {
		if !CheckCommand("uv") {
			s.Logger.Warning("uv is not installed, skipping docs scaffolding (install: https://docs.astral.sh/uv/)")
			cfg.SkipDocs = true
		}
	}

	// Check existing state
	CheckExistingState(cfg.ProjectDir, s.Logger)

	s.Logger.Plain("")
	s.Logger.Info("Starting project initialization...")
	s.Logger.Plain("")

	// Run steps — respect --only-docs guard
	if !cfg.OnlyDocs {
		steps := []func() error{
			s.stepInstallBmad,
			s.stepInstallCobraCli,
			s.stepGoModInit,
			s.stepCobraInit,
			s.stepAddVersionCmd,
			s.stepGenerateServeCmd,
			s.stepGenerateMockeryConfig,
			s.stepGenerateEditorConfig,
			s.stepGenerateUIPlaceholder,
			s.stepGenerateEmbedGo,
			s.stepGoModTidy,
			s.stepGenerateMakefile,
		}
		for _, step := range steps {
			if err := step(); err != nil {
				return err
			}
		}
	}

	// Docs (runs in both normal and --only-docs mode)
	if err := s.stepInitDocs(); err != nil {
		return err
	}

	if !cfg.OnlyDocs {
		// UI init
		if err := s.stepInitUI(); err != nil {
			return err
		}

		// Git init
		if err := s.stepInitGit(); err != nil {
			return err
		}
	}

	s.stepPrintSummary()
	return nil
}
