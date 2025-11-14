# LineSense TODO

Progress tracker for v1.0 release and beyond.

## v1.0 Release Checklist

### Testing & Quality Assurance

- [x] Safety filter unit tests (17/17 passing)
- [x] Shell integration test execution
  - [x] Run `./test_shell_integration.sh` and verify all pass
  - [ ] Test bash integration manually
  - [ ] Test zsh integration manually
- [ ] End-to-end test with real API
  - [ ] Set up OPENROUTER_API_KEY
  - [ ] Run `go test -v ./test_e2e.go`
  - [ ] Verify suggest command works
  - [ ] Verify explain command works
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

- [ ] Installation script
  - [ ] Create `install.sh` for automated setup
  - [ ] Handle prerequisites (Go, jq optional)
  - [ ] Auto-detect shell (bash/zsh)
  - [ ] Copy example configs
  - [ ] Add to PATH
  - [ ] Source shell integration
- [ ] Cross-platform builds
  - [ ] Build for Linux (amd64, arm64)
  - [ ] Build for macOS (amd64, arm64)
  - [ ] Build for Windows (amd64)
  - [ ] Test on each platform
- [ ] Package management
  - [ ] Create Homebrew formula (macOS)
  - [ ] Create AUR package (Arch Linux)
  - [ ] Consider snap/flatpak (Linux)
  - [ ] Consider Debian/RPM packages
- [ ] Distribution
  - [ ] GitHub Releases with binaries
  - [ ] Release notes template
  - [ ] Version tagging strategy

### CI/CD Pipeline

- [ ] GitHub Actions setup
  - [ ] Automated testing on push/PR
  - [ ] Multi-platform test matrix (Linux, macOS)
  - [ ] Go versions matrix (1.21, 1.22, latest)
- [ ] Code quality checks
  - [ ] Add `golangci-lint` configuration
  - [ ] Run `gofmt` check
  - [ ] Run `go vet`
  - [ ] Static analysis with `staticcheck`
- [ ] Automated builds
  - [ ] Build binaries on release tag
  - [ ] Upload artifacts to GitHub Releases
  - [ ] Generate checksums (SHA256)
- [ ] Documentation deployment
  - [ ] Auto-generate docs site (optional)
  - [ ] Update GitHub Pages (optional)

### Documentation Polish

- [ ] README improvements
  - [ ] Add badges (build status, coverage, version)
  - [ ] Add demo GIF/video
  - [ ] Add "Star this repo" section
  - [ ] Add table of contents
- [ ] Add missing docs
  - [ ] LICENSE file (choose license)
  - [ ] CONTRIBUTING.md
  - [ ] CODE_OF_CONDUCT.md
  - [ ] CHANGELOG.md
- [ ] Documentation review
  - [ ] Fix any broken links
  - [ ] Check code examples work
  - [ ] Ensure consistency across docs
  - [ ] Spell check all docs

### Community & Project Setup

- [ ] GitHub repository setup
  - [ ] Issue templates (bug, feature request, question)
  - [ ] Pull request template
  - [ ] GitHub Discussions enabled
  - [ ] Project board for tracking
- [ ] Community files
  - [ ] Code of Conduct
  - [ ] Contributing guidelines
  - [ ] Security policy
- [ ] Project metadata
  - [ ] Choose OSI-approved license
  - [ ] Add project description
  - [ ] Add topics/tags to repo
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
  - [ ] Update version in `main.go` (currently 0.1.0)
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
  - [ ] Selection UI
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
**Current Focus**: v1.0 Release - Testing & Quality Assurance
