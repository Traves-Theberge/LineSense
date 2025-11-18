# LineSense TODO

Progress tracker for v1.0 release and beyond.

## v1.0 Release Checklist

### Testing & Quality Assurance

- [x] Safety filter unit tests (17/17 passing)
- [x] Shell integration test execution
  - [x] Run `./test_shell_integration.sh` and verify all pass
  - [ ] Test zsh integration manually
- [x] End-to-end test with real API
  - [x] Set up OPENROUTER_API_KEY
  - [x] Run `go test -v ./test_e2e.go`
  - [x] Verify suggest command works
  - [x] Verify explain command works
- [x] Additional unit tests (107 tests total, all passing ✓✓)
  - [x] Add tests for `internal/ai` module (66.1% coverage ✓)
  - [x] Add tests for `internal/config` module (84.8% coverage ✓)
  - [x] Add tests for context gathering functions (90.7% coverage ✓✓✓)
- [x] Test coverage reporting
  - [x] Run `go test -cover ./...`
  - [x] Config module: 84.8% coverage (>80% ✓✓)
  - [x] Core module: **90.7% coverage** (>90% ✓✓✓✓)
  - [x] AI module: 66.1% coverage (>60% ✓)
  - [x] Generate coverage report with `go test -coverprofile=coverage.out` ✓
  - [x] Generate HTML coverage report (`coverage.html`) ✓

### Installation & Distribution

- [x] Installation script ✓✓✓
  - [x] Create `install.sh` for automated setup ✓
  - [x] Handle prerequisites (Go, jq optional) ✓
  - [x] Auto-detect shell (bash/zsh) ✓
  - [x] Copy example configs ✓
  - [x] Add to PATH ✓
  - [x] Source shell integration ✓
- [x] Cross-platform builds (via GoReleaser) ✓✓✓
  - [x] Build for Linux (amd64, arm64) ✓
  - [x] Build for macOS (amd64, arm64) ✓
  - [x] Build for Windows (amd64, arm64) ✓
  - [ ] Test on each platform (requires release)
- [x] Package management (prepared) ✓
  - [x] GoReleaser config for Homebrew formula ✓
  - [ ] Create AUR package (Arch Linux)
  - [ ] Consider snap/flatpak (Linux)
  - [ ] Consider Debian/RPM packages
- [x] Distribution ✓✓✓
  - [x] GitHub Releases with binaries (via GoReleaser) ✓
  - [x] Release notes template (automated via GoReleaser) ✓
  - [x] Version tagging strategy (semantic versioning) ✓
  - [x] Checksums and signatures ✓

### CI/CD Pipeline

- [x] GitHub Actions setup ✓✓✓
  - [x] Automated testing on push/PR ✓
  - [x] Multi-platform test matrix (Linux, macOS, Windows) ✓
  - [x] Go versions matrix (1.21, 1.22, 1.23) ✓
  - [x] Shell integration tests in CI ✓
- [x] Code quality checks ✓✓✓
  - [x] Add `golangci-lint` configuration (.golangci.yml) ✓
  - [x] Run `gofmt` check ✓
  - [x] Run `go vet` ✓
  - [x] Static analysis with golangci-lint ✓
  - [x] Security scanning with Gosec ✓
- [x] Automated builds ✓✓✓
  - [x] Build binaries on release tag (via GoReleaser) ✓
  - [x] Upload artifacts to GitHub Releases ✓
  - [x] Generate checksums (SHA256) ✓
  - [x] Support for multiple architectures ✓
- [ ] Documentation deployment
  - [ ] Auto-generate docs site (optional)
  - [ ] Update GitHub Pages (optional)

### Documentation Polish

- [x] README improvements
  - [x] Add badges (build status, coverage, version) ✓
  - [x] Add Testing & Quality section ✓
  - [x] Add Contributing section ✓
  - [ ] Add demo GIF/video
  - [ ] Add table of contents (optional)
- [x] Add missing docs
  - [x] LICENSE file (MIT License) ✓
  - [x] CONTRIBUTING.md ✓
  - [x] CHANGELOG.md ✓
  - [x] TESTING.md ✓
  - [ ] CODE_OF_CONDUCT.md (optional)
- [ ] Documentation review
  - [ ] Fix any broken links
  - [ ] Check code examples work
  - [ ] Ensure consistency across docs
  - [ ] Spell check all docs

### Community & Project Setup

- [x] GitHub repository setup
  - [x] Issue templates (bug, feature request) ✓
  - [x] Pull request template ✓
  - [ ] GitHub Discussions enabled (requires repo)
  - [ ] Project board for tracking (optional)
- [x] Community files
  - [x] Contributing guidelines ✓
  - [x] License (MIT) ✓
  - [ ] Code of Conduct (optional)
  - [ ] Security policy (SECURITY.md exists)
- [ ] Project metadata
  - [x] Choose OSI-approved license (MIT) ✓
  - [ ] Add project description (when repo created)
  - [ ] Add topics/tags to repo (when repo created)
  - [ ] Set up GitHub Sponsors (optional)

### Pre-Release Testing

- [ ] Manual testing checklist
  - [ ] Fresh install on clean system
  - [ ] Test `config init` flow
  - [ ] Test `config set-key` flow
  - [ ] Test suggest with various prompts
  - [ ] Test explain with various commands
  - [ ] Test safety filters with dangerous commands
  - [ ] Test shell integration in bash
  - [ ] Test shell integration in zsh
- [ ] Beta testing
  - [ ] Get 3-5 users to test
  - [ ] Collect feedback
  - [ ] Fix critical bugs
  - [ ] Update docs based on confusion points

### Release Preparation

- [ ] Version bump
  - [x] Update version in `main.go` (updated to 0.4.0) ✓
  - [ ] Update CHANGELOG.md
  - [ ] Tag release (v1.0.0)
- [ ] Release checklist
  - [ ] All tests passing
  - [ ] Documentation complete
  - [ ] Binaries built for all platforms
  - [ ] Release notes written
  - [ ] Security review complete
- [ ] Launch
  - [ ] Create GitHub Release
  - [ ] Post on relevant communities (r/golang, Hacker News, etc.)
  - [ ] Tweet/social media announcement
  - [ ] Update project website (if applicable)

---

## Phase 3: Advanced Features (Post v1.0)

### Usage Logging & Learning

- [ ] Basic logging infrastructure
  - [ ] SQLite database for local storage
  - [ ] Track accepted suggestions
  - [ ] Track rejected suggestions
  - [ ] Privacy controls and opt-out
- [ ] Analytics
  - [ ] Most common commands
  - [ ] Suggestion accuracy rate
  - [ ] Usage statistics
  - [ ] Export analytics data
- [ ] Personalization
  - [ ] Learn from user patterns
  - [ ] Boost similar suggestions
  - [ ] Context-aware weighting
  - [ ] User preference learning

### Enhanced UX Features

- [ ] Multi-suggestion support
  - [ ] Show top 3-5 suggestions
  - [ ] Arrow key navigation
  - [ ] Preview mode
  - [ ] Selection UIho
- [ ] Better error handling
  - [ ] Retry logic for API failures
  - [ ] Fallback suggestions
  - [ ] Actionable error messages
  - [ ] Network timeout handling
- [ ] Offline mode
  - [ ] Cache frequent suggestions
  - [ ] Local pattern matching
  - [ ] Suggestion history replay
  - [ ] Works without internet

### Developer Experience

- [ ] Developer documentation
  - [ ] Comprehensive CONTRIBUTING.md
  - [ ] Development setup guide
  - [ ] Architecture documentation
  - [ ] API design docs
- [ ] Development tools
  - [ ] Mock API for testing
  - [ ] Debugging mode
  - [ ] Verbose logging option
  - [ ] Performance profiling
- [ ] Plugin system
  - [ ] Custom context providers
  - [ ] Custom safety filters
  - [ ] Custom suggestion sources
  - [ ] Plugin API documentation

---

## Known Issues

- [ ] API key not visible in environment initially (requires shell reload)
- [ ] Shell integration requires manual PATH setup
- [ ] No Windows native support yet (WSL only)
- [ ] Large command histories slow down context gathering

---

## Ideas / Future Considerations

- [ ] Local model support (Ollama integration)
- [ ] Custom prompt templates
- [ ] Command aliases integration
- [ ] History search integration
- [ ] Syntax highlighting in shell
- [ ] Autocomplete for flags/options
- [ ] Integration with other shells (fish, nushell)
- [ ] Web UI for configuration
- [ ] VS Code extension
- [ ] Mobile app (Termux)

---

**Last Updated**: 2025-11-14
**Current Focus**: v1.0 Release - Pre-Release Testing & Manual Verification

**Recent Completions**:
- ✅ CI/CD Infrastructure (GitHub Actions, GoReleaser)
- ✅ Installation automation (install.sh)
- ✅ Comprehensive testing suite (107 tests, 90.7% core coverage)
- ✅ Complete documentation suite
- ✅ Community infrastructure (CONTRIBUTING, issue templates, etc.)
