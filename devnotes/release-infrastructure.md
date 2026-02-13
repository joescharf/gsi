# Release Infrastructure, Code Signing, GitHub Pages, and SPA Fix

*2026-02-13T18:25:09Z*

This release brings gsi's scaffold output to parity with the pm/fdsn reference implementations. Four issues were resolved: updated release infrastructure templates, macOS code signing support, GitHub Pages automation, and an SPA blank-page-on-refresh fix.

## 1. main.* ldflags Pattern

Scaffolded projects now use the main.* ldflags pattern instead of cmd.* -- version vars live in main.go and are passed to cmd.Execute(version, commit, date). This matches the pm project's architecture and avoids coupling goreleaser ldflags to internal package paths.

New templates `main_go.tmpl` and `cmd_root_go.tmpl` are generated after cobra-cli init, overwriting the cobra-cli defaults. This uses the new `OverwriteTemplateFile()` helper in files.go which doesn't skip existing files.

## 2. Three-Platform Builds

The goreleaser template now produces separate build IDs for linux, macos, and windows (all amd64+arm64), with per-OS archive formats (tar.gz for Linux, zip for macOS/Windows). The name_template uses `title .Os` for proper casing.

```bash
go run . --dry-run /tmp/gsi-demo 2>&1 | grep -E '(Creating|Skipping|workflow|pycodesign|GitHub|main.go|root.go)'
```

```output
[0;34mℹ[0m Creating Cobra CLI application structure
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/main.go
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/cmd/root.go
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/.github/workflows/release.yml
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/.github/workflows/ci.yml
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/.github/workflows/docs.yml
[1;33m⚠[0m [DRY-RUN] Would create /tmp/gsi-demo/gsi-demo_pycodesign.ini
[1;33m⚠[0m GitHub repo joescharf/gsi-demo not found, skipping Pages configuration
  gh api repos/joescharf/gsi-demo/pages -X POST --field build_type=workflow
  2. Update the project description in cmd/root.go
[0;34mℹ[0m GitHub Setup:
  Run 'gh repo create joescharf/gsi-demo --public --source=.' to create the GitHub repo
  Run 'gh api repos/joescharf/gsi-demo/pages -X POST --field build_type=workflow' to enable GitHub Pages
```

## 3. CI and Docs Workflows

Two new workflow templates are scaffolded:

- **ci.yml** -- test + lint jobs on push/PR, with bun UI build and embed steps
- **docs.yml** -- GitHub Pages deployment via mkdocs-material (build + deploy jobs)

These are gated by CapRelease and CapDocs respectively.

## 4. macOS Code Signing

A `pycodesign_ini.tmpl` template is scaffolded for each project with placeholder cert fingerprints. The goreleaser template includes a commented-out post-hook under universal_binaries for signing.

gsi itself now has its own `gsi_pycodesign.ini` with Scharfnado LLC certs, and its `.goreleaser.yml` is updated with windows builds, signing post-hook, and split archives.

## 5. GitHub Pages Automation

A new `stepConfigureGitHubPages()` step attempts to enable Pages via `gh api` during scaffold. If the repo doesn't exist yet (common for new projects), it gracefully falls back to printing setup instructions in the summary.

The summary now includes a "GitHub Setup" section with the exact `gh` commands needed.

## 6. SPA Blank Page Fix

The root cause of blank pages on browser refresh was that `bun init --react=shadcn` creates a React app whose default build uses relative asset paths. When the browser refreshes on a client-side route like `/settings`, the server returns `index.html` but the browser tries to load assets relative to `/settings/` instead of `/`.

The fix: a new `build_ts.tmpl` template that sets `publicPath: "/"` in the Bun build config. The `stepInitUI()` now writes this `build.ts` after `bun init` and updates `package.json` to use `"bun run build.ts"` as the build script.

## 7. Updated Existing Templates

Several existing templates were updated to match the pm reference:

- **goreleaser_yml.tmpl** -- Full rewrite: 3 builds, `-s -w` strip flags, UI before-hooks, brews + dockers_v2
- **dockerfile.tmpl** -- Alpine 3.21 (pinned), non-root user, tzdata, `PROJECT_DB_PATH` env var, `CMD ["serve"]`
- **github_release_yml.tmpl** -- Manual dispatch, bun setup, QEMU, Buildx, GHCR login, GoReleaser latest
- **makefile.tmpl** -- main.* ldflags with -s -w, `serve: all`, `release-local` target
- **mkdocs_yml.tmpl** -- `site_url` on github.io, `repo_url`, `repo_name`, `edit_uri`
- **cmd_version.go.tmpl** -- Uses `buildVersion`/`buildCommit`/`buildDate` from Execute()

## Verification

```bash
go test -race -count=1 ./... 2>&1 | tail -8
```

```output
?   	github.com/joescharf/gsi	[no test files]
?   	github.com/joescharf/gsi/cmd	[no test files]
ok  	github.com/joescharf/gsi/internal/logger	1.279s
ok  	github.com/joescharf/gsi/internal/scaffold	1.645s
ok  	github.com/joescharf/gsi/internal/templates	1.433s
```

```bash
goreleaser check 2>&1
```

```output
  • checking                                  path=.goreleaser.yml
  • brews is being phased out in favor of homebrew_casks, check https://goreleaser.com/deprecations#brews for more info
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
```

```bash
go build -o /dev/null . && echo 'Binary compiles cleanly'
```

```output
Binary compiles cleanly
```

## Files Changed

### New files
| File | Purpose |
|------|---------|
| `internal/templates/files/github_ci_yml.tmpl` | CI workflow template |
| `internal/templates/files/github_docs_yml.tmpl` | Docs workflow template |
| `internal/templates/files/main_go.tmpl` | main.go with version vars |
| `internal/templates/files/cmd_root_go.tmpl` | Root cmd with Execute(v,c,d) |
| `internal/templates/files/pycodesign_ini.tmpl` | Pycodesign config template |
| `internal/templates/files/build_ts.tmpl` | Bun build script with publicPath |
| `gsi_pycodesign.ini` | gsi's own pycodesign config |
| `.github/workflows/release.yml` | gsi's own release workflow |

### Modified files
| File | Changes |
|------|---------|
| `internal/templates/templates.go` | Added `ProjectNameUpper` to Data |
| `internal/templates/templates_test.go` | Added 8 new template tests, updated assertions |
| `internal/scaffold/steps.go` | 6 new steps, updated stepInitUI + stepPrintSummary |
| `internal/scaffold/scaffold.go` | Registered new steps |
| `internal/scaffold/files.go` | Added OverwriteTemplateFile() |
| `internal/templates/files/goreleaser_yml.tmpl` | Full rewrite (3-platform) |
| `internal/templates/files/dockerfile.tmpl` | Alpine 3.21, non-root, ENV |
| `internal/templates/files/github_release_yml.tmpl` | Bun/QEMU/Buildx/GHCR |
| `internal/templates/files/makefile.tmpl` | main.* ldflags, release-local |
| `internal/templates/files/mkdocs_yml.tmpl` | github.io URLs, repo fields |
| `internal/templates/files/cmd_version.go.tmpl` | Uses buildVersion vars |
| `.goreleaser.yml` | Windows build, signing, split archives |
| `Makefile` | main.* ldflags, release-local |
| `main.go` | Version vars, cmd.Execute(v,c,d) |
| `cmd/root.go` | Execute(v,c,d), buildVersion vars |
| `cmd/version.go` | Uses buildVersion |
