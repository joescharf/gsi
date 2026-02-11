package templates

import (
	"strings"
	"testing"
)

func TestRenderWithVariables(t *testing.T) {
	data := Data{
		ProjectName:  "myapp",
		GoModulePath: "github.com/example/myapp",
	}

	tests := []struct {
		template string
		contains []string
	}{
		{"cmd_serve.go.tmpl", []string{"github.com/example/myapp/internal/ui", "package cmd"}},
		{"mockery_yml.tmpl", []string{"github.com/example/myapp", "with-expecter: true"}},
		{"editorconfig.tmpl", []string{"root = true", "indent_style = tab"}},
		{"index_html.tmpl", []string{"myapp", "<title>myapp</title>"}},
		{"embed_go.tmpl", []string{"package ui", "//go:embed all:dist"}},
		{"makefile.tmpl", []string{"BINARY_NAME", "go build"}},
		{"mkdocs_yml.tmpl", []string{"myapp Documentation", "material"}},
		{"docs_gitignore.tmpl", []string{"site/", ".venv/"}},
		{"docs_index_md.tmpl", []string{"# myapp", "Getting Started"}},
		{"docs_getting_started_md.tmpl", []string{"go install github.com/example/myapp@latest", "myapp --help"}},
		{"docs_extra_css.tmpl", []string{".md-nav__item", "font-size"}},
		{"gitignore.tmpl", []string{"bin/", ".DS_Store", "vendor/"}},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
			result, err := Render(tt.template, data)
			if err != nil {
				t.Fatalf("Render(%s) failed: %v", tt.template, err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("Render(%s) missing %q in output:\n%s", tt.template, want, result)
				}
			}
		})
	}
}

func TestRenderMissingTemplate(t *testing.T) {
	_, err := Render("nonexistent.tmpl", Data{})
	if err == nil {
		t.Error("expected error for missing template, got nil")
	}
}
