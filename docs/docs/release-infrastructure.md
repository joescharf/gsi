# Release Infrastructure

This document specifies the complete release infrastructure that gsi scaffolds for new projects. It covers GitHub Actions CI/CD, GoReleaser configuration, Docker images, Homebrew distribution, and GitHub Pages documentation deployment.

The reference implementation is the [pm](https://github.com/joescharf/pm) project, with patterns drawn from [fdsn](https://github.com/joescharf/fdsn).

## Overview

gsi scaffolds a full release pipeline with five components:

| Component | File(s) | Trigger | Artifacts |
|-----------|---------|---------|-----------|
| **CI** | `.github/workflows/ci.yml` | Push to main, PRs | Test + lint results |
| **Release** | `.github/workflows/release.yml` | `v*` tag push | Binaries, Docker, Homebrew |
| **Docs** | `.github/workflows/docs.yml` | `docs/**` changes on main | GitHub Pages site |
| **GoReleaser** | `.goreleaser.yml` | Called by release workflow | Archives, checksums, images, formula |
| **Docker** | `Dockerfile` | Called by GoReleaser | Multi-arch GHCR images |

```
v* tag push ──► release.yml ──► goreleaser ──┬── binaries (linux/darwin/windows × amd64/arm64)
                                              ├── archives (tar.gz, zip for windows)
                                              ├── checksums.txt
                                              ├── Docker images → ghcr.io/<owner>/<project>
                                              ├── Homebrew cask → <owner>/homebrew-tap
                                              └── GitHub Release with changelog

push to main ──► ci.yml ──┬── test (go test -v -race -count=1)
                          └── lint (golangci-lint)

docs/** push ──► docs.yml ──► mkdocs build ──► GitHub Pages
```

## GoReleaser Configuration

**File:** `.goreleaser.yml`

**Template:** `internal/templates/files/goreleaser_yml.tmpl`

### Key design decisions

- **Single build matrix** — one `builds:` entry covers all six platform combinations (linux/darwin/windows × amd64/arm64). The old template had two separate builds (linux-only + darwin-amd64-only with broken universal_binaries).
- **ldflags target `main.*`** — variables are `main.version`, `main.commit`, `main.date` (in `main.go`), not `cmd.*`. The `-s -w` flags strip debug symbols for smaller binaries.
- **Before hooks build the UI** — `bun install --frozen-lockfile`, `bun run build`, then copy `ui/dist/*` into `internal/ui/dist/`. This ensures embedded UI is fresh for every release.
- **`homebrew_casks:`** — the non-deprecated goreleaser v2 key (replaces `brews:`). Pushes to `<owner>/homebrew-tap` repo using `HOMEBREW_TAP_TOKEN`.
- **`dockers_v2:`** — the non-deprecated goreleaser v2 key (replaces `dockers:` + `docker_manifests:`). Automatically handles multi-arch builds. The Dockerfile must use `ARG TARGETPLATFORM` / `COPY ${TARGETPLATFORM}/<binary>`.
- **No `release.draft: true`** — releases publish immediately. Drafts add friction with no benefit when using tag-triggered CI.

### Template content

```yaml
version: 2

before:
  hooks:
    - sh -c "cd ui && bun install --frozen-lockfile"
    - sh -c "cd ui && bun run build"
    - sh -c "rm -rf internal/ui/dist/assets && cp -r ui/dist/* internal/ui/dist/"

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

archives:
  - formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats:
          - zip
    files:
      - README.md

checksum:
  name_template: "checksums.txt"

snapshot:
  version_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc
  groups:
    - title: Features
      regexp: "^.*feat[(\\w)]*:+.*$"
      order: 0
    - title: "Bug fixes"
      regexp: "^.*fix[(\\w)]*:+.*$"
      order: 1
    - title: Others
      order: 999
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^ci:"
      - "^chore:"

homebrew_casks:
  - name: <project>
    repository:
      owner: <owner>
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    directory: Casks
    homepage: https://github.com/<owner>/<project>
    description: "<project description>"

dockers_v2:
  - images:
      - "ghcr.io/<owner>/<project>"
    tags:
      - "v{{ .Version }}"
      - latest
```

### Go template escaping

The goreleaser YAML contains its own `{{ }}` template expressions. In gsi's Go templates, these must be escaped:

| GoReleaser syntax | Go template escaping |
|-------------------|---------------------|
| `{{.Version}}` | `{{"{{"}} .Version {{"}}"}}` |
| `{{ incpatch .Version }}` | `{{"{{"}} incpatch .Version {{"}}"}}` |
| `{{ .Env.HOMEBREW_TAP_TOKEN }}` | `{{"{{"}} .Env.HOMEBREW_TAP_TOKEN {{"}}"}}` |

Go template variables (`{{.ProjectName}}`, `{{.GoModuleOwner}}`) render normally.

## Dockerfile

**File:** `Dockerfile`

**Template:** `internal/templates/files/dockerfile.tmpl`

### Key design decisions

- **Pinned Alpine version** (`alpine:3.21`) — reproducible builds, not `alpine:latest`.
- **Non-root user** — creates a dedicated system user/group for the application.
- **`/data` directory** — owned by the app user, used for SQLite databases. The `DB_PATH` env var points here.
- **`tzdata`** — required for `time.LoadLocation()` in Go binaries built with `CGO_ENABLED=0`.
- **`TARGETPLATFORM` ARG** — required by goreleaser's `dockers_v2` to copy the correct platform-specific binary.
- **`ENTRYPOINT` + `CMD`** — binary as entrypoint, `serve` as default command.

### Template content

```dockerfile
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S <project> && adduser -S <project> -G <project>
RUN mkdir -p /data && chown <project>:<project> /data
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/<project> /usr/local/bin/<project>
USER <project>
ENV <PROJECT_UPPER>_DB_PATH=/data/<project>.db
EXPOSE 8080
ENTRYPOINT ["<project>"]
CMD ["serve"]
```

!!! note
    The `ENV` line uses an uppercase project name prefix (e.g., `PM_DB_PATH`, `FDSN_DB_PATH`). The template should derive this from `ProjectName`.

## GitHub Actions Workflows

### CI Workflow

**File:** `.github/workflows/ci.yml`

**Template:** `internal/templates/files/github_ci_yml.tmpl` *(new)*

**Triggers:** Push to `main`, pull requests targeting `main`.

Two parallel jobs:

- **test** — Builds embedded UI, runs `go test -v -race -count=1 ./...` and `go vet ./...`
- **lint** — Builds embedded UI, runs `golangci-lint` via the official action

Both jobs set up Go (from `go.mod`) and Bun, then build and embed the UI before running checks. This ensures the `internal/ui/dist/` embed directory is populated so Go compilation succeeds.

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: oven-sh/setup-bun@v2

      - name: Install UI dependencies
        run: cd ui && bun install --frozen-lockfile

      - name: Build UI
        run: cd ui && bun run build

      - name: Embed UI
        run: rm -rf internal/ui/dist/assets && cp -r ui/dist/* internal/ui/dist/

      - name: Run tests
        run: go test -v -race -count=1 ./...

      - name: Run vet
        run: go vet ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: oven-sh/setup-bun@v2

      - name: Install UI dependencies
        run: cd ui && bun install --frozen-lockfile

      - name: Build UI
        run: cd ui && bun run build

      - name: Embed UI
        run: rm -rf internal/ui/dist/assets && cp -r ui/dist/* internal/ui/dist/

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: "v2.9"
```

!!! note
    The UI build steps should only be included when the `ui` capability is relevant. Since gsi always scaffolds the `internal/ui/embed.go` with `//go:embed`, the UI build is always needed for `go test` to compile. If a project has no `ui/` directory, the before hooks in goreleaser and CI steps would need adjustment.

### Release Workflow

**File:** `.github/workflows/release.yml`

**Template:** `internal/templates/files/github_release_yml.tmpl` *(update existing)*

**Triggers:** Push of `v*` tags.

Changes from the old template:

| Aspect | Old | New |
|--------|-----|-----|
| Go version | `go-version: stable` | `go-version-file: go.mod` |
| Bun setup | missing | `oven-sh/setup-bun@v2` |
| QEMU | missing | `docker/setup-qemu-action@v3` |
| Docker Buildx | missing | `docker/setup-buildx-action@v3` |
| GHCR login | missing | `docker/login-action@v3` |
| GoReleaser version | `"~> v2"` | `latest` |

The old template was minimal (just Go + goreleaser). The new template includes everything needed for Docker multi-arch builds and GHCR publishing.

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write
  packages: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: oven-sh/setup-bun@v2

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

### Docs Workflow

**File:** `.github/workflows/docs.yml`

**Template:** `internal/templates/files/github_docs_yml.tmpl` *(new)*

**Triggers:** Push to `main` with changes in `docs/**`, or manual dispatch.

Two jobs: **build** (uv sync + mkdocs build + upload pages artifact) and **deploy** (deploy to GitHub Pages).

```yaml
name: Docs

on:
  push:
    branches: [main]
    paths:
      - "docs/**"
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: astral-sh/setup-uv@v5

      - name: Install dependencies
        run: cd docs && uv sync

      - name: Build docs
        run: cd docs && uv run mkdocs build

      - uses: actions/upload-pages-artifact@v3
        with:
          path: docs/site

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

## MkDocs Configuration

**File:** `docs/mkdocs.yml`

**Template:** `internal/templates/files/mkdocs_yml.tmpl` *(update existing)*

Changes from the old template:

| Field | Old | New |
|-------|-----|-----|
| `site_url` | `https://<project>.example.com` | `https://<owner>.github.io/<project>/` |
| `repo_url` | *(missing)* | `https://github.com/<owner>/<project>` |
| `repo_name` | *(missing)* | `<owner>/<project>` |
| `edit_uri` | `""` | `edit/main/docs/docs/` |

The `repo_url` adds a GitHub link to the docs header. The `edit_uri` enables "Edit on GitHub" links on each page. The `site_url` points to the actual GitHub Pages URL where docs will be deployed.

## Docs .gitignore

**File:** `docs/.gitignore`

**Template:** `internal/templates/files/docs_gitignore.tmpl`

The current template is already correct:

```
site/
.venv/
__pycache__/
.cache/
```

## Makefile

**File:** `Makefile`

**Template:** `internal/templates/files/makefile.tmpl` *(update existing)*

### ldflags fix

The ldflags line must target `main.*`, not `$(MODULE)/cmd.*`:

```makefile
# Old (broken):
LDFLAGS := -ldflags "-X $(MODULE)/cmd.version=$(VERSION) ..."

# New (correct):
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"
```

The `-s -w` flags strip debug info, matching goreleaser's ldflags.

### serve target

Change from `serve: build` to `serve: all` so the UI is built before serving:

```makefile
# Old:
serve: build ## Start the embedded web UI server

# New:
serve: all ## Start the embedded web UI server
```

## Scaffolding Changes Summary

### New templates to create

| Template file | Output path | Gated by |
|---------------|-------------|----------|
| `github_ci_yml.tmpl` | `.github/workflows/ci.yml` | `CapRelease` (same gate as release workflow) |
| `github_docs_yml.tmpl` | `.github/workflows/docs.yml` | `CapDocs` |

### Existing templates to update

| Template file | Changes |
|---------------|---------|
| `goreleaser_yml.tmpl` | Single build matrix, correct ldflags, `homebrew_casks:`, `dockers_v2:` (new format), before hooks for UI build, remove `universal_binaries` and `release.draft` |
| `dockerfile.tmpl` | Pinned alpine, non-root user, tzdata, `/data` dir, `TARGETPLATFORM`, env var, expose port, default CMD |
| `github_release_yml.tmpl` | Add Bun setup, QEMU, Docker Buildx, GHCR login, `go-version-file` |
| `makefile.tmpl` | Fix ldflags to target `main.*`, add `-s -w`, change serve dependency to `all` |
| `mkdocs_yml.tmpl` | Fix `site_url`, add `repo_url`, `repo_name`, `edit_uri` |

### New scaffold steps to add

In `internal/scaffold/steps.go`:

```go
// stepGenerateCIWorkflow writes .github/workflows/ci.yml from template.
func (s *Scaffolder) stepGenerateCIWorkflow() error {
	if !s.Config.IsEnabled(CapRelease) {
		s.Logger.Info("Skipping CI workflow (--no-release)")
		return nil
	}
	dir := filepath.Join(s.Config.ProjectDir, ".github", "workflows")
	if !s.Config.DryRun {
		os.MkdirAll(dir, 0o755)
	}
	return WriteTemplateFile(
		filepath.Join(dir, "ci.yml"),
		"github_ci_yml.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}

// stepGenerateDocsWorkflow writes .github/workflows/docs.yml from template.
func (s *Scaffolder) stepGenerateDocsWorkflow() error {
	if !s.Config.IsEnabled(CapDocs) {
		s.Logger.Info("Skipping docs workflow (--no-docs)")
		return nil
	}
	dir := filepath.Join(s.Config.ProjectDir, ".github", "workflows")
	if !s.Config.DryRun {
		os.MkdirAll(dir, 0o755)
	}
	return WriteTemplateFile(
		filepath.Join(dir, "docs.yml"),
		"github_docs_yml.tmpl",
		s.templateData(),
		s.Config.DryRun,
		s.Logger,
	)
}
```

### Step registration in `scaffold.go`

Add the new steps to the `steps` slice in `Run()`, after `stepGenerateReleaseWorkflow`:

```go
steps := []func() error{
    // ... existing steps ...
    s.stepGenerateReleaseWorkflow,
    s.stepGenerateCIWorkflow,        // new
    s.stepGenerateDocsWorkflow,      // new
}
```

### Template Data additions

The `templates.Data` struct may need a new field for the uppercase project name used in Docker ENV vars (e.g., `PM_DB_PATH`). Options:

1. Add `ProjectNameUpper string` to `templates.Data` and compute it in `templateData()`
2. Use a Go template function: `{{ .ProjectName | upper }}`
3. Omit the ENV line from the Dockerfile template (projects can add it manually)

Option 1 is simplest and most explicit.

## Post-Scaffold Setup (Manual)

After scaffolding, users must complete these one-time GitHub configuration steps:

### 1. HOMEBREW_TAP_TOKEN secret

Create a GitHub PAT (classic) with `repo` scope that has write access to the `<owner>/homebrew-tap` repository. Add it as a repository secret named `HOMEBREW_TAP_TOKEN` in the project's Settings > Secrets and variables > Actions.

### 2. GitHub Pages source

In the project's Settings > Pages, set the Build and deployment Source to **GitHub Actions**.

### 3. GHCR package visibility

After the first release, check the package visibility at `ghcr.io/<owner>/<project>` and make it public if desired.

### 4. Homebrew tap repository

Ensure `<owner>/homebrew-tap` exists on GitHub with a `Casks/` directory. GoReleaser will push `.rb` cask files there on release.

## Local Development

### Testing a release locally

```bash
make release-snapshot
```

This runs `goreleaser release --snapshot --clean --skip docker,homebrew` and produces artifacts in `dist/`.

### Verifying goreleaser config

```bash
goreleaser check
```

### Building docs locally

```bash
make docs-serve    # live-reload dev server
make docs-build    # static build to docs/site/
```

## Verification Checklist

After scaffolding a new project, verify:

- [ ] `goreleaser check` passes with no warnings
- [ ] `goreleaser release --snapshot --clean --skip docker,homebrew` builds all 6 binaries
- [ ] `make docs-build` produces `docs/site/`
- [ ] `.github/workflows/ci.yml` exists and has test + lint jobs
- [ ] `.github/workflows/release.yml` includes QEMU, Buildx, GHCR login steps
- [ ] `.github/workflows/docs.yml` exists with pages deployment
- [ ] Dockerfile uses pinned Alpine, non-root user, and `TARGETPLATFORM`
- [ ] ldflags in both Makefile and .goreleaser.yml target `main.*` (not `cmd.*`)
