package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joescharf/gsi/internal/logger"
	"github.com/joescharf/gsi/internal/templates"
)

// WriteTemplateFile renders a template and writes it to path, skipping if the file already exists.
func WriteTemplateFile(path, templateName string, data templates.Data, dryRun bool, log *logger.Logger) error {
	return writeTemplateFileWithMode(path, templateName, data, 0o644, dryRun, log)
}

// WriteExecutableTemplateFile renders a template and writes it to path with executable permissions (0o755),
// skipping if the file already exists.
func WriteExecutableTemplateFile(path, templateName string, data templates.Data, dryRun bool, log *logger.Logger) error {
	return writeTemplateFileWithMode(path, templateName, data, 0o755, dryRun, log)
}

// writeTemplateFileWithMode renders a template and writes it to path with the given file mode,
// skipping if the file already exists.
func writeTemplateFileWithMode(path, templateName string, data templates.Data, mode os.FileMode, dryRun bool, log *logger.Logger) error {
	if _, err := os.Stat(path); err == nil {
		log.Info(path + " already exists, skipping")
		return nil
	}

	if dryRun {
		log.Warning(fmt.Sprintf("[DRY-RUN] Would create %s", path))
		return nil
	}

	log.Info("Creating " + path)

	content, err := templates.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("rendering template %s: %w", templateName, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	log.Success("Created " + path)
	return nil
}

// WriteStaticFile writes static content to path, skipping if the file already exists.
func WriteStaticFile(path string, content []byte, dryRun bool, log *logger.Logger) error {
	if _, err := os.Stat(path); err == nil {
		log.Info(path + " already exists, skipping")
		return nil
	}

	if dryRun {
		log.Warning(fmt.Sprintf("[DRY-RUN] Would create %s", path))
		return nil
	}

	log.Info("Creating " + path)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	log.Success("Created " + path)
	return nil
}
