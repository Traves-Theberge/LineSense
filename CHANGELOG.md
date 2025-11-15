# Changelog

All notable changes to LineSense will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-11-14

### Added
- **CI/CD Infrastructure**
  - GitHub Actions CI workflow (test on multiple Go versions and platforms)
  - GitHub Actions release workflow with GoReleaser
  - golangci-lint configuration for code quality
  - Security scanning with Gosec
  - Shell integration tests in CI
  - Automated cross-platform builds (Linux, macOS, Windows on amd64/arm64)
- **Installation & Distribution**
  - Automated installation script (`install.sh`)
  - GoReleaser configuration for professional releases
  - Support for Homebrew tap (prepared for future)
  - Docker image support (prepared for future)
  - Checksums and release artifacts
- **Testing & Quality**
  - Comprehensive test suite with 107 tests
  - Test coverage: Core 90.7%, Config 84.8%, AI 66.1%
  - Testing guide with best practices ([docs/TESTING.md](docs/TESTING.md))
- **Documentation & Community**
  - Complete documentation suite (INSTALLATION, CONFIGURATION, SECURITY, API, TESTING, CI_CD)
  - CI/CD and release process guide ([docs/CI_CD.md](docs/CI_CD.md))
  - MIT License
  - Contributing guidelines ([CONTRIBUTING.md](CONTRIBUTING.md))
  - GitHub issue templates (bug report, feature request)
  - GitHub pull request template
  - Comprehensive CHANGELOG

### Changed
- Updated README with CI badge and Go Report Card
- Enhanced README with automated installation instructions
- Improved installation documentation with multiple methods

## [0.2.5] - 2024-11-14

### Added
- Configuration management commands (`config init`, `config set-key`, `config set-model`, `config show`)
- Interactive setup wizard for first-time users
- Secure API key storage in shell RC files
- API key masking in output display
- File permission management (0600 for sensitive files)
- Auto-detection of user's shell (bash/zsh)

### Security
- API keys now stored in environment variables instead of config files
- Proper file permissions automatically enforced
- API key masking in all output

## [0.2.0] - 2024-11-14

### Added
- Shell integration scripts for bash and zsh
- Readline/ZLE widget bindings for autocompletion
- Color-coded risk indicators in shell output
- JSON parsing with jq (with grep/sed fallback)
- Configurable keybindings via environment variables
- Safety filter system with risk classification
- Built-in dangerous command detection (rm -rf, dd, mkfs, etc.)
- Configurable command denylists
- Three-tier risk classification (low/medium/high)
- Comprehensive safety filter tests (17 test cases)

### Changed
- Improved shell integration with smart JSON handling
- Enhanced output formatting for better readability

## [0.1.0] - 2024-11-14

### Added
- Core CLI infrastructure
- Configuration system with TOML support
- XDG Base Directory specification compliance
- OpenRouter API integration
- Context gathering (git, history, environment)
- AI-powered command suggestions
- AI-powered command explanations
- `suggest` command for getting command suggestions
- `explain` command for understanding commands
- Git repository detection and status parsing
- Shell history parsing (bash and zsh)
- Environment variable filtering
- Provider profile system
- Multiple model support
- Temperature and token configuration
- End-to-end test framework

### Security
- Environment variable filtering to exclude sensitive data
- Risk assessment for suggested commands

## [0.0.1] - 2024-11-13

### Added
- Initial project structure
- Product Requirements Document (PRD)
- Basic Go module setup

---

## Version Naming

- **Major version** (X.0.0) - Breaking changes
- **Minor version** (0.X.0) - New features, backward compatible
- **Patch version** (0.0.X) - Bug fixes, backward compatible

## Links

- [Compare v0.2.5...HEAD](https://github.com/Traves-Theberge/LineSense/compare/v0.2.5...HEAD)
- [Compare v0.2.0...v0.2.5](https://github.com/Traves-Theberge/LineSense/compare/v0.2.0...v0.2.5)
- [All Releases](https://github.com/Traves-Theberge/LineSense/releases)
