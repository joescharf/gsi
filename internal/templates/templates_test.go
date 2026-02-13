package templates

import (
	"strings"
	"testing"
)

func TestRenderWithVariables(t *testing.T) {
	data := Data{
		ProjectName:      "myapp",
		ProjectNameUpper: "MYAPP",
		GoModulePath:     "github.com/example/myapp",
		GoModuleOwner:    "example",
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
		{"makefile.tmpl", []string{"BINARY_NAME", "go build", "main.version", "HOMEBREW_TAP_TOKEN", "serve: all"}},
		{"mkdocs_yml.tmpl", []string{"myapp Documentation", "material", "example.github.io/myapp", "repo_url", "repo_name: example/myapp", "edit/main/docs/docs/"}},
		{"docs_gitignore.tmpl", []string{"site/", ".venv/"}},
		{"docs_index_md.tmpl", []string{"# myapp", "Getting Started"}},
		{"docs_getting_started_md.tmpl", []string{"go install github.com/example/myapp@latest", "myapp --help"}},
		{"docs_extra_css.tmpl", []string{".md-nav__item", "font-size"}},
		{"gitignore.tmpl", []string{"bin/", ".DS_Store", "vendor/"}},
		{"goreleaser_yml.tmpl", []string{"project_name: myapp", "main.version", "ghcr.io/example/myapp", "myapp-linux", "myapp-macos", "myapp-windows", "homebrew_casks:", "dockers_v2:"}},
		{"dockerfile.tmpl", []string{"alpine:3.21", "MYAPP_DB_PATH", "USER myapp", "COPY ${TARGETPLATFORM}/myapp", "ENTRYPOINT", "CMD [\"serve\"]"}},
		{"dockerignore.tmpl", []string{".git", "dist/"}},
		{"docs_scripts_scrape_sh.tmpl", []string{"#!/bin/bash", `--title "myapp"`, "shot-scraper"}},
		{"docs_scripts_shots_yaml.tmpl", []string{"myapp-dashboard.png", "localhost:8080"}},
		{"docs_scripts_add_browser_frame_py.tmpl", []string{"#!/usr/bin/env python3", "pillow", "SUPERSAMPLE_SCALE"}},
		{"github_release_yml.tmpl", []string{"go-version-file: go.mod", "oven-sh/setup-bun", "docker/setup-qemu-action", "docker/setup-buildx-action", "docker/login-action", "version: latest"}},
		{"github_ci_yml.tmpl", []string{"go-version-file: go.mod", "oven-sh/setup-bun", "go test", "go vet", "golangci-lint"}},
		{"github_docs_yml.tmpl", []string{"astral-sh/setup-uv", "mkdocs build", "upload-pages-artifact", "deploy-pages"}},
		{"main_go.tmpl", []string{"github.com/example/myapp/cmd", "cmd.Execute(version, commit, date)"}},
		{"cmd_root_go.tmpl", []string{"package cmd", "func Execute(version, commit, date string)", "buildVersion"}},
		{"cmd_version.go.tmpl", []string{"package cmd", "buildVersion", "buildCommit", "buildDate"}},
		{"build_ts.tmpl", []string{"bun-plugin-tailwind", `publicPath: "/"`, "Bun.build"}},
		{"pycodesign_ini.tmpl", []string{"application_id", "bundle_id = com.example.myapp", "myapp-macos_darwin_all/myapp"}},
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
