# Changelog

All notable changes to LineSense will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.3] - 2025-11-15

### Added
- **Multiple Command Suggestions**
  - Now returns 3-5 alternative command suggestions instead of just one
  - Suggestions ordered from most likely to least likely
  - Better handling of typos and ambiguous inputs with intent interpretation
  - Each suggestion gets independent risk assessment

### Improved
- **Smarter AI Prompts**
  - Updated system prompt to request multiple suggestions
  - Added explicit guidance for typo/intent interpretation
  - Better structured response format (one command per line)

### Technical
- Rewrote `parseSuggestions()` to handle multiple line-separated commands
- Added robustness with numbered list stripping (1. 2. 3. etc.)
- Limited to max 5 suggestions to avoid overwhelming users

## [0.4.2] - 2025-11-15

### Changed
- **Beautiful Shell Integration**
  - Shell integration now uses the same beautiful UI as the CLI (`--format pretty`)
  - Removed custom JSON parsing and formatting code
  - Suggestions display full styled output with boxes, colors, and borders
  - Explanations show the same professional format as CLI
  - Simplified explain keybinding from Ctrl+X Ctrl+E to just Ctrl+X

### Improved
- **Better Typo Handling**
  - Shell integration now handles typos intelligently with intent-based suggestions
  - Multiple suggestion options displayed (not just auto-replacement)
  - All CLI features automatically available in shell integration

### Removed
- Removed `_linesense_parse_json` function (no longer needed)
- Removed jq dependency for shell integration
- Removed custom formatting code (delegates to CLI)

## [0.4.1] - 2025-11-15

### Fixed
- **Shell Integration Keybindings**
  - Fixed keybinding variable quoting that was stripping backslashes
  - Changed default suggest keybinding from `\C- ` to `\C-@` (proper readline notation for Ctrl+Space)
  - Keybindings now work correctly: Ctrl+Space for suggestions, Ctrl+X Ctrl+E for explanations
  - Removed nested quotes in parameter expansion to preserve control sequences

### Changed
- **Silent Shell Integration by Default**
  - Removed startup messages when shell integration loads
  - Clean, unobtrusive loading experience
  - Updated documentation to reflect silent loading as default behavior

### Documentation
- Updated README to document silent shell integration
- Updated INSTALLATION.md with correct shell integration paths
- Fixed shell script source paths to use `~/.config/linesense/shell/`
- Added "Silent Loading" section to installation docs

## [0.4.0] - 2025-11-14

### Added
- **Beautiful UI with Charm Libraries**
  - Integrated Charm Bubbles for animated loading spinner
  - Added Lipgloss for stunning styled output with colors
  - Rounded borders and boxes for professional appearance
  - Color-coded risk levels (green/yellow/red indicators)
  - Professional typography with proper spacing
  - Icons for different sections (💡 suggestions, 📖 explanations)
- **New Output Formatting**
  - `--format` flag to choose between 'pretty' (default) or 'json'
  - Pretty format with styled suggestions and explanations
  - Boxed content with visual hierarchy
  - Numbered suggestion lists with risk indicators
  - Dynamic terminal width detection and text wrapping
- **Loading Indicators**
  - Animated spinner while AI is processing
  - Context-aware messages ("Generating suggestions...", "Analyzing command...")
  - Graceful fallback in non-interactive environments
- **Comprehensive Testing**
  - End-to-end test suite (test_e2e.sh)
  - Shell integration validation (test_integration.sh)
  - Performance benchmarking
  - Safety filter testing

### Changed
- Pretty output is now the default format (was JSON)
- Terminal width automatically adjusts to screen size
- Text wraps properly within borders (min: 40, max: 100 chars)
- Removed old ANSI color codes in favor of Lipgloss styling
- Cleaner spinner messages without emojis

### Technical
- Added charmbracelet/bubbles dependency for spinner
- Added charmbracelet/lipgloss dependency for styling
- Added golang.org/x/term for terminal size detection
- Refactored output logic into separate ui.go file
- Improved error handling in non-TTY environments

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
