# CI/CD and Release Process

This document explains LineSense's continuous integration, continuous deployment, and release processes.

## Table of Contents

- [Overview](#overview)
- [Continuous Integration (CI)](#continuous-integration-ci)
- [Release Process](#release-process)
- [GoReleaser Configuration](#goreleaser-configuration)
- [Making a Release](#making-a-release)
- [Troubleshooting](#troubleshooting)

## Overview

LineSense uses GitHub Actions for automated testing, building, and releasing. The CI/CD pipeline ensures:

- ✅ All tests pass on multiple platforms and Go versions
- ✅ Code quality meets standards (formatting, linting, security)
- ✅ Binaries are built for all supported platforms
- ✅ Releases are automated and consistent
- ✅ Artifacts are properly checksummed and signed

## Continuous Integration (CI)

### CI Workflow

The CI workflow (`.github/workflows/ci.yml`) runs on:
- Every push to `main` or `develop` branches
- Every pull request to `main` or `develop` branches

### CI Jobs

#### 1. Test Job

Runs the test suite across multiple configurations:

**Matrix:**
- **Operating Systems:** Ubuntu (Linux), macOS
- **Go Versions:** 1.21, 1.22, 1.23

**Steps:**
1. Checkout code
2. Set up Go environment
3. Cache Go modules for faster builds
4. Download dependencies
5. Run tests with race detection: `go test -v -race -coverprofile=coverage.out ./...`
6. Generate coverage report
7. Upload coverage to Codecov (for main branch only)

**Example:**
```yaml
- name: Run tests
  run: go test -v -race -coverprofile=coverage.out ./...
```

#### 2. Lint Job

Runs golangci-lint to check code quality:

**Checks:**
- `errcheck` - Unchecked errors
- `gosimple` - Code simplifications
- `govet` - Go vet issues
- `staticcheck` - Static analysis
- `gosec` - Security issues
- `revive` - Code style
- And many more (see `.golangci.yml`)

**Example:**
```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
```

#### 3. Format Check Job

Verifies all Go code is properly formatted:

```bash
gofmt -l .  # Lists unformatted files
```

If any files are unformatted, the job fails with a diff showing needed changes.

#### 4. Go Vet Job

Runs `go vet` to find suspicious constructs:

```bash
go vet ./...
```

#### 5. Build Job

Builds binaries for all target platforms:

**Platforms:**
- Linux (Ubuntu)
- macOS
- Windows

**Steps:**
1. Build: `go build -v ./cmd/linesense`
2. Test binary: `./linesense --version` and `./linesense --help`

This ensures binaries are buildable and executable on all platforms.

#### 6. Shell Integration Tests

Tests shell integration scripts:

**Steps:**
1. Install zsh (Ubuntu already has bash)
2. Run `test_shell_integration.sh`
3. Verify bash and zsh scripts have valid syntax:
   ```bash
   bash -n shell/linesense.bash
   zsh -n shell/linesense.zsh
   ```

#### 7. Security Scan

Runs Gosec security scanner:

**What it checks:**
- Hardcoded credentials
- SQL injection vulnerabilities
- Command injection risks
- Weak cryptography
- File permission issues
- And more

**Output:** SARIF format uploaded to GitHub Security tab

### Viewing CI Results

1. **In Pull Requests:**
   - CI status appears at the bottom of the PR
   - Click "Details" next to each check to see logs

2. **In Actions Tab:**
   - Visit `https://github.com/Traves-Theberge/LineSense/actions`
   - Click on a workflow run to see all jobs
   - Click on a job to see detailed logs

3. **Badges:**
   - README shows CI badge: [![CI](https://github.com/Traves-Theberge/LineSense/workflows/CI/badge.svg)](https://github.com/Traves-Theberge/LineSense/actions)

## Release Process

### Release Workflow

The release workflow (`.github/workflows/release.yml`) runs when you push a version tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

### GoReleaser

LineSense uses [GoReleaser](https://goreleaser.com/) for professional release automation.

**What GoReleaser does:**
1. ✅ Builds binaries for multiple platforms and architectures
2. ✅ Creates archives (`.tar.gz` for Unix, `.zip` for Windows)
3. ✅ Generates checksums (SHA256)
4. ✅ Creates GitHub Release with release notes
5. ✅ Generates changelog from git commits
6. ✅ (Future) Publishes to Homebrew, Docker Hub

### Supported Platforms

GoReleaser builds for:

| OS      | Architectures        |
|---------|---------------------|
| Linux   | amd64, arm64, arm   |
| macOS   | amd64, arm64 (Apple Silicon) |
| Windows | amd64               |

**Note:** Windows arm64 builds are excluded as they're not commonly needed.

### Release Assets

Each release includes:

```
LineSense_v1.0.0_Linux_x86_64.tar.gz
LineSense_v1.0.0_Linux_arm64.tar.gz
LineSense_v1.0.0_Darwin_x86_64.tar.gz
LineSense_v1.0.0_Darwin_arm64.tar.gz
LineSense_v1.0.0_Windows_x86_64.zip
checksums.txt
```

Each archive contains:
- The `linesense` binary
- README.md
- LICENSE
- CHANGELOG.md
- Shell integration scripts (`shell/linesense.{bash,zsh}`)
- Example configuration files (`examples/*.toml`)

### Changelog Generation

GoReleaser automatically generates changelogs from commit messages:

**Commit format:**
```
<type>(<scope>): <description>

feat(ai): add support for custom temperature settings
fix(context): handle empty shell history gracefully
docs(readme): update installation instructions
```

**Commit types:**
- `feat:` - New features (appears under "Features")
- `fix:` - Bug fixes (appears under "Bug Fixes")
- `perf:` - Performance improvements
- `docs:` - Documentation only (excluded from changelog)
- `test:` - Tests only (excluded from changelog)
- `chore:` - Maintenance (excluded from changelog)

**Example generated changelog:**
```markdown
## Features
- **ai**: add support for custom temperature settings
- **config**: interactive setup wizard

## Bug Fixes
- **context**: handle empty shell history gracefully
- **safety**: fix nil pointer in IsBlocked
```

## GoReleaser Configuration

The `.goreleaser.yml` file configures the entire release process.

### Key Sections

#### Builds

```yaml
builds:
  - id: linesense
    binary: linesense
    main: ./cmd/linesense
    env:
      - CGO_ENABLED=0  # Static binaries
    ldflags:
      - -s -w  # Strip debug info
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
```

**ldflags explanation:**
- `-s -w` - Reduce binary size by stripping debug symbols
- `-X main.version=...` - Embed version info in binary

#### Archives

```yaml
archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    files:
      - README.md
      - LICENSE
      - shell/*
      - examples/*
```

Creates compressed archives with all necessary files.

#### Release Notes

```yaml
release:
  header: |
    ## LineSense {{ .Tag }}

    ### Installation
    ```bash
    curl -fsSL https://raw.githubusercontent.com/traves/LineSense/{{ .Tag }}/install.sh | bash
    ```
```

Adds installation instructions to every release.

#### Future Features

The config includes (but skips) support for:
- **Homebrew:** Automated formula updates
- **Docker:** Multi-arch container images
- **GPG Signing:** Cryptographic signatures

To enable these, set `skip: false` and configure the necessary tokens.

## Making a Release

### Prerequisites

1. **All tests passing:**
   ```bash
   go test ./...
   ```

2. **Update CHANGELOG.md:**
   ```markdown
   ## [1.0.0] - 2025-11-15

   ### Added
   - New feature X
   - New feature Y

   ### Fixed
   - Bug Z
   ```

3. **Update version references:**
   - Check `main.go` version string (if any)
   - Ensure README examples use latest version

### Release Steps

#### 1. Create and Push Tag

```bash
# For a new major release
git tag v1.0.0

# For a patch release
git tag v1.0.1

# For a prerelease
git tag v1.0.0-beta.1

# Push the tag
git push origin v1.0.0
```

**Version Naming:**
- `v1.0.0` - Major release
- `v1.1.0` - Minor release (new features)
- `v1.0.1` - Patch release (bug fixes)
- `v1.0.0-rc.1` - Release candidate
- `v1.0.0-beta.1` - Beta release

#### 2. GitHub Actions Runs Automatically

1. Tests run on all platforms
2. GoReleaser builds binaries
3. Archives are created
4. Checksums are generated
5. GitHub Release is created

**Monitor progress:**
```
https://github.com/Traves-Theberge/LineSense/actions
```

#### 3. Verify Release

1. Visit: `https://github.com/Traves-Theberge/LineSense/releases/latest`
2. Check:
   - ✅ All binaries present
   - ✅ Checksums file exists
   - ✅ Release notes look good
   - ✅ Archives are downloadable

#### 4. Test Installation

Download and test the release:

```bash
# Linux/macOS
wget https://github.com/Traves-Theberge/LineSense/releases/download/v1.0.0/LineSense_v1.0.0_Linux_x86_64.tar.gz
tar -xzf LineSense_v1.0.0_Linux_x86_64.tar.gz
./linesense --version

# Or use install script
curl -fsSL https://raw.githubusercontent.com/traves/LineSense/v1.0.0/install.sh | bash
```

#### 5. Announce Release

- Post on GitHub Discussions
- Share on social media
- Update project website (if applicable)
- Notify users in relevant communities

### Hotfix Release

If you need to fix a critical bug:

```bash
# Create fix branch
git checkout -b hotfix/critical-bug main

# Make fixes
git commit -m "fix: critical security issue"

# Merge to main
git checkout main
git merge hotfix/critical-bug

# Tag and release
git tag v1.0.1
git push origin main v1.0.1
```

## Testing GoReleaser Locally

Before making a release, test GoReleaser locally:

### Install GoReleaser

```bash
# macOS
brew install goreleaser

# Linux
go install github.com/goreleaser/goreleaser@latest

# Or download from https://goreleaser.com/install/
```

### Test Build

```bash
# Snapshot build (no tag required)
goreleaser release --snapshot --clean

# Check what would be released
goreleaser release --skip=publish --clean

# Build only, no release
goreleaser build --snapshot --clean
```

**Output directory:** `dist/`

```bash
$ ls dist/
LineSense_darwin_amd64/
LineSense_darwin_arm64/
LineSense_linux_amd64/
LineSense_linux_arm64/
LineSense_windows_amd64/
checksums.txt
```

### Validate Configuration

```bash
# Check .goreleaser.yml syntax
goreleaser check

# Show what would happen (dry run)
goreleaser release --skip=publish --skip=validate --clean
```

## Troubleshooting

### CI Failing

**Problem:** Tests fail on CI but pass locally

**Solutions:**
1. Check Go version: `go version` (CI uses 1.21, 1.22, 1.23)
2. Check for race conditions: `go test -race ./...`
3. Check OS-specific issues (test on Linux if you develop on macOS)
4. Review CI logs carefully for environment differences

**Problem:** golangci-lint fails

**Solutions:**
1. Run locally: `golangci-lint run`
2. Auto-fix issues: `golangci-lint run --fix`
3. Format code: `gofmt -w .`
4. Review `.golangci.yml` for disabled checks

**Problem:** Security scan fails

**Solutions:**
1. Run Gosec locally: `gosec ./...`
2. Review security issues carefully
3. Add exceptions to `.goreleaser.yml` if false positive (with justification)

### Release Issues

**Problem:** GoReleaser fails with "tag not found"

**Solution:**
```bash
git fetch --tags
git push origin --tags
```

**Problem:** Binary doesn't run on target platform

**Solution:**
1. Ensure `CGO_ENABLED=0` for static binaries
2. Test with Docker:
   ```bash
   docker run --rm -v $PWD/dist:/dist ubuntu /dist/LineSense_linux_amd64/linesense --version
   ```

**Problem:** Release notes are empty/incomplete

**Solution:**
1. Ensure CHANGELOG.md is updated
2. Use conventional commit messages
3. Check GoReleaser changelog config

**Problem:** Archives missing files

**Solution:**
1. Check `files:` section in `.goreleaser.yml`
2. Ensure paths are relative to project root
3. Test locally: `goreleaser build --snapshot`

### Getting Help

- **GoReleaser Docs:** https://goreleaser.com/
- **GitHub Actions Docs:** https://docs.github.com/en/actions
- **Project Issues:** https://github.com/Traves-Theberge/LineSense/issues
- **golangci-lint Docs:** https://golangci-lint.run/

## Best Practices

### Commit Messages

Use conventional commits for better changelogs:

```bash
feat(ai): add Anthropic Claude support
fix(config): prevent nil pointer in LoadConfig
docs(readme): update installation steps
test(core): add tests for git integration
chore(deps): update dependencies
```

### Versioning

Follow Semantic Versioning (semver):

- **Major (v2.0.0):** Breaking changes
- **Minor (v1.1.0):** New features, backward compatible
- **Patch (v1.0.1):** Bug fixes, backward compatible

### Testing Before Release

**Pre-release checklist:**
1. ✅ `go test ./...` passes
2. ✅ `golangci-lint run` passes
3. ✅ `goreleaser check` passes
4. ✅ `goreleaser release --snapshot` builds successfully
5. ✅ Manual testing on target platforms
6. ✅ CHANGELOG.md updated
7. ✅ Version references updated

### Release Frequency

- **Patch releases:** As needed for bug fixes
- **Minor releases:** Monthly or when significant features are ready
- **Major releases:** Rarely, only for breaking changes

### Release Communication

For each release:
1. Create detailed release notes
2. Highlight breaking changes prominently
3. Provide migration guides for major versions
4. Announce in project discussions
5. Update documentation

## Metrics and Monitoring

### CI Metrics

Track:
- **Build time:** Should be < 10 minutes
- **Test time:** Should be < 5 minutes
- **Success rate:** Aim for > 95%
- **Flaky tests:** Investigate and fix

### Release Metrics

Track:
- **Download counts:** GitHub provides stats
- **Popular platforms:** Which binaries are most downloaded
- **Installation method:** Script vs manual vs package manager

### Quality Metrics

Track:
- **Test coverage:** Maintain > 80%
- **Linter issues:** Keep at 0
- **Security issues:** Address immediately
- **Bug reports:** Track and prioritize

---

**Related Documentation:**
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Contribution guidelines
- [TESTING.md](TESTING.md) - Testing guide
- [CHANGELOG.md](../CHANGELOG.md) - Version history

**Last Updated:** 2025-11-14
