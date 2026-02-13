# Releasing

gsi scaffolds a complete [goreleaser](https://goreleaser.com) configuration that automates building binaries for Linux, macOS, and Windows, Docker images, Homebrew distribution, and optional macOS code signing.

## Prerequisites

- [goreleaser](https://goreleaser.com/install/) installed
- A GitHub personal access token at `~/.config/goreleaser/github_token` (for real releases)
- A `HOMEBREW_TAP_TOKEN` with write access to your tap repository (if using Homebrew distribution)
- [Docker](https://docs.docker.com/get-docker/) running (for Docker image builds)

## Testing a Release

Use snapshot mode to verify everything works without publishing:

```bash
make release-snapshot
```

This runs:

```bash
goreleaser release --snapshot --clean --skip docker,homebrew
```

Check the output in `dist/`:

```bash
ls dist/
# myapp_Linux_x86_64.tar.gz
# myapp_Linux_arm64.tar.gz
# myapp_Darwin_all.zip
# myapp_Windows_x86_64.zip
# myapp_Windows_arm64.zip
# checksums.txt

# Verify the binary
./dist/myapp-macos_darwin_all/myapp version
```

## Creating a Release

### Via GitHub Actions (recommended)

The scaffolded release workflow uses manual dispatch:

1. Go to **Actions > Release** in your GitHub repo
2. Click **Run workflow**
3. Enter the tag (e.g., `v1.0.0`)
4. Click **Run workflow**

The workflow handles Go setup, Bun, QEMU, Docker Buildx, GHCR login, and GoReleaser.

### Locally

For signed local releases (macOS code signing):

```bash
make release-local
```

For standard releases:

```bash
make release
```

## What Gets Built

The scaffolded `.goreleaser.yml` configures:

| Artifact | Details |
|----------|---------|
| Linux binaries | amd64 and arm64, CGO disabled, tar.gz archives |
| macOS universal binary | Combined amd64+arm64, zip archive |
| Windows binaries | amd64 and arm64, zip archives |
| Docker images | Multi-platform via `dockers_v2`, pushed to GHCR |
| Checksums | SHA256 in `checksums.txt` |
| Changelog | Grouped by feat/fix/other, excludes docs/test/ci/chore commits |

### ldflags

All builds use:

```
-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
```

These target variables in `main.go` which passes them to `cmd.Execute(version, commit, date)`.

### Before Hooks

GoReleaser runs these before any build:

1. `go mod download`
2. `cd ui && bun install --frozen-lockfile`
3. `cd ui && bun run build`
4. Copy `ui/dist/*` into `internal/ui/dist/` for embedding

## macOS Code Signing

gsi scaffolds a `<project>_pycodesign.ini` template for macOS code signing and notarization. To enable:

1. Replace the placeholder cert fingerprints in `<project>_pycodesign.ini`
2. Uncomment the `hooks.post` line in `.goreleaser.yml` under `universal_binaries`
3. Switch from `brews:` to `homebrew_casks:` for distributing signed .pkg installers

The pycodesign integration:

- Signs the macOS universal binary with Developer ID Application cert
- Creates a notarized .pkg installer with Developer ID Installer cert
- Uses `uv run pycodesign.py` for the signing workflow

## CI/CD Workflows

gsi scaffolds three GitHub Actions workflows:

### Release (`.github/workflows/release.yml`)

Manual dispatch with:

- Go setup from `go.mod`
- Bun setup for UI builds
- QEMU for multi-arch Docker
- Docker Buildx for multi-arch builds
- GHCR login for container registry
- GoReleaser (latest version)

### CI (`.github/workflows/ci.yml`)

Runs on push to main and PRs:

- **test** job: Build UI, embed, run `go test` and `go vet`
- **lint** job: Build UI, embed, run golangci-lint

### Docs (`.github/workflows/docs.yml`)

Runs on push to main (docs/** paths) and manual dispatch:

- **build** job: `uv sync` + `mkdocs build` + upload pages artifact
- **deploy** job: Deploy to GitHub Pages

## Homebrew Setup

The scaffolded config uses `brews:` (Formula) by default, suitable for unsigned CLI tools:

1. Create a tap repository (e.g., `<owner>/homebrew-tap`)
2. Create a GitHub personal access token (classic) with `repo` scope
3. For local releases, set: `export HOMEBREW_TAP_TOKEN=ghp_...`
4. For CI, add `HOMEBREW_TAP_TOKEN` as a repository secret

For signed macOS apps, switch to `homebrew_casks:` with a `Casks/` directory.

After setup, users can install via:

```bash
brew install <owner>/tap/<project>
```

## Docker Setup

Docker images are built automatically via `dockers_v2`. The Dockerfile uses:

- Alpine 3.21 (pinned, not `latest`)
- Non-root user for security
- `ca-certificates` and `tzdata` packages
- `ARG TARGETPLATFORM` for multi-arch support
- Environment variable for DB path: `<PROJECT>_DB_PATH=/data/<project>.db`
- Exposed port 8080

The images are tagged as:

- `ghcr.io/<owner>/<project>:v<version>`
- `ghcr.io/<owner>/<project>:latest`

## GitHub Pages Setup

After creating your GitHub repo, enable Pages:

```bash
gh api repos/<owner>/<project>/pages -X POST --field build_type=workflow
gh repo edit --homepage "https://<owner>.github.io/<project>/"
```

gsi attempts this automatically during scaffold if the repo already exists on GitHub.
