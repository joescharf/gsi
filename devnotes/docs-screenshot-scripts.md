# Docs Screenshot Scripts

*2026-02-13T17:39:55Z*

Added a docs screenshot pipeline to the gsi scaffold. When `CapDocs` is enabled, new projects now get a `docs/scripts/` directory with three files:

- **scrape.sh** — Orchestrates the full pipeline: health-checks the local server, captures screenshots with shot-scraper, adds macOS browser frames, compresses with imageoptim-cli, and copies to the docs site.
- **shots.yaml** — Starter shot-scraper configuration with two example pages (dashboard + detail). Uses the project name in output filenames.
- **add_browser_frame.py** — Python script (runs via uv) that adds macOS-style title bars with traffic-light buttons, rounded corners, and drop shadows to screenshots using Pillow.

## Implementation Details

**No new capability flag** — The scripts are part of docs tooling, gated by the existing `CapDocs` capability.

**New Go function**: `WriteExecutableTemplateFile` in `internal/scaffold/files.go` writes templates with `0o755` permissions (used for `scrape.sh`). Refactored `WriteTemplateFile` to share a private `writeTemplateFileWithMode` helper.

**Template variables**: `scrape.sh` and `shots.yaml` use `{{.ProjectName}}` for the browser frame title and output filenames. `add_browser_frame.py` is a pass-through (no Go template variables).

**Docs .gitignore** updated to ignore `scripts/img/` and `scripts/img-raw/` (intermediate screenshot artifacts).

```bash
go test ./internal/scaffold/ -run 'TestWriteExecutable' -v -count=1 2>&1
```

```output
=== RUN   TestWriteExecutableTemplateFileCreatesFile
--- PASS: TestWriteExecutableTemplateFileCreatesFile (0.00s)
=== RUN   TestWriteExecutableTemplateFileIdempotent
--- PASS: TestWriteExecutableTemplateFileIdempotent (0.00s)
=== RUN   TestWriteExecutableTemplateFileDryRun
--- PASS: TestWriteExecutableTemplateFileDryRun (0.00s)
PASS
ok  	github.com/joescharf/gsi/internal/scaffold	0.176s
```

```bash
go test ./internal/templates/ -run 'TestRenderWithVariables/docs_scripts' -v -count=1 2>&1
```

```output
=== RUN   TestRenderWithVariables
=== RUN   TestRenderWithVariables/docs_scripts_scrape_sh.tmpl
=== RUN   TestRenderWithVariables/docs_scripts_shots_yaml.tmpl
=== RUN   TestRenderWithVariables/docs_scripts_add_browser_frame_py.tmpl
--- PASS: TestRenderWithVariables (0.00s)
    --- PASS: TestRenderWithVariables/docs_scripts_scrape_sh.tmpl (0.00s)
    --- PASS: TestRenderWithVariables/docs_scripts_shots_yaml.tmpl (0.00s)
    --- PASS: TestRenderWithVariables/docs_scripts_add_browser_frame_py.tmpl (0.00s)
PASS
ok  	github.com/joescharf/gsi/internal/templates	0.150s
```

```bash
cd /tmp && rm -rf mytest && /Users/joescharf/app/gsi/gsi mytest --dry-run 2>&1 | grep -E 'scripts/(scrape|shots|add_browser)'
```

```output
```

```bash
cd /tmp && rm -rf mytest && /Users/joescharf/app/gsi/gsi mytest --dry-run 2>&1 | grep -E 'scripts/(scrape|shots|add_browser)'
```

```output
```

```bash
cd /tmp && rm -rf mytest && /Users/joescharf/app/gsi/gsi mytest --dry-run 2>&1 | grep -E 'scripts/(scrape|shots|add_browser)'
```

```output
```

```bash
cd /tmp && rm -rf mytest && /Users/joescharf/app/gsi/gsi mytest --dry-run 2>&1 | grep scripts
```

```output
```

```bash
cd /tmp && rm -rf mytest && /Users/joescharf/app/gsi/gsi mytest --dry-run 2>&1 | grep -E 'scripts/(scrape|shots|add_browser)'
```

```output
[1;33m⚠[0m [DRY-RUN] Would create /tmp/mytest/docs/scripts/scrape.sh
[1;33m⚠[0m [DRY-RUN] Would create /tmp/mytest/docs/scripts/shots.yaml
[1;33m⚠[0m [DRY-RUN] Would create /tmp/mytest/docs/scripts/add_browser_frame.py
```

## Files Changed

| File | Action |
|------|--------|
| `internal/scaffold/files.go` | Added `WriteExecutableTemplateFile` + `writeTemplateFileWithMode` helper |
| `internal/scaffold/files_test.go` | 3 new tests for executable template writing |
| `internal/templates/files/docs_scripts_add_browser_frame_py.tmpl` | New template (Python, browser frames) |
| `internal/templates/files/docs_scripts_scrape_sh.tmpl` | New template (bash, pipeline orchestrator) |
| `internal/templates/files/docs_scripts_shots_yaml.tmpl` | New template (YAML, shot-scraper config) |
| `internal/templates/files/docs_gitignore.tmpl` | Added `scripts/img/` and `scripts/img-raw/` |
| `internal/templates/templates_test.go` | 3 new template render test entries |
| `internal/scaffold/steps.go` | Wired scripts into `stepInitDocs()` |
| `docs/docs/scaffolded-output.md` | Updated directory tree and file descriptions |
| `README.md` | Updated docs section with screenshot pipeline |
