# Releasing

gsi scaffolds a [goreleaser](https://goreleaser.com) configuration that automates building binaries, Docker images, and optionally Homebrew distribution.

## Prerequisites

- [goreleaser](https://goreleaser.com/install/) installed
- A GitHub personal access token at `~/.config/goreleaser/github_token` (for real releases)
- [Docker](https://docs.docker.com/get-docker/) running (for Docker image builds)

## Testing a Release

Use snapshot mode to verify everything works without publishing:

```bash
make release-snapshot
```

This runs:

```bash
goreleaser release --snapshot --clean --skip homebrew,docker
```

Check the output in `dist/`:

```bash
ls dist/
# myapp_0.0.1-devel_linux_amd64.tar.gz
# myapp_0.0.1-devel_darwin_all.zip
# checksums.txt

# Verify the binary
./dist/myapp-macos_darwin_all/myapp version
```

## Creating a Release

### 1. Tag the Release

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

### 2. Run Goreleaser

```bash
make release
```

This creates:

- **Linux binaries** (amd64, arm64) as tar.gz archives
- **macOS universal binary** as a zip archive
- **Docker images** pushed to `ghcr.io/<owner>/<project>` with `v<version>` and `latest` tags
- **GitHub draft release** with changelog, binaries, and checksums

### 3. Publish the Draft

Go to your GitHub repository's Releases page, review the draft, and publish it.

## What Gets Built

The scaffolded `.goreleaser.yml` configures:

| Artifact | Details |
|----------|---------|
| Linux binaries | amd64 and arm64, CGO disabled |
| macOS universal binary | Combined amd64 binary |
| Linux archives | tar.gz format |
| macOS archives | zip format |
| Docker images | Multi-platform via `dockers_v2` |
| Checksums | SHA256 in `checksums.txt` |
| Changelog | Grouped by feat/fix/other, excludes docs/test commits |

## Homebrew Setup

The Homebrew section is commented out by default. To enable it:

1. Create a tap repository (e.g., `<owner>/homebrew-tap`)
2. Uncomment the `homebrew_casks` section in `.goreleaser.yml`
3. Configure the repository owner, name, and description
4. Ensure your GitHub token has repo write access to the tap

After setup, users can install via:

```bash
brew install <owner>/tap/<project>
```

## Docker Setup

Docker images are built automatically. To push to GitHub Container Registry:

1. Authenticate: `echo $GITHUB_TOKEN | docker login ghcr.io -u <username> --password-stdin`
2. Ensure the `Dockerfile` is in your project root
3. Run `make release` (or use CI/CD)

The images are tagged as:

- `ghcr.io/<owner>/<project>:v<version>`
- `ghcr.io/<owner>/<project>:latest`

## CI/CD Integration

For GitHub Actions, create `.github/workflows/release.yml`:

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
          go-version: stable
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
