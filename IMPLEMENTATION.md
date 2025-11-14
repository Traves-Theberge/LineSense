# LineSense Implementation Guide

This document provides guidance for implementing the LineSense features based on the PRD.

## Current Status

**✅ Phase 1 (Configuration & Context) - COMPLETED**

All core infrastructure for configuration loading and context gathering is implemented and tested.

## Implementation Priority

### Phase 1: Core Infrastructure (v0.1-alpha) ✅ COMPLETED

1. **Configuration Loading** (`internal/config`) ✅
   - [x] Implement XDG config directory resolution
   - [x] Add TOML parsing (`github.com/BurntSushi/toml` v1.5.0)
   - [x] Implement `LoadConfig()` in [config.go](internal/config/config.go:46-66)
   - [x] Implement `LoadProvidersConfig()` in [providers.go](internal/config/providers.go:27-53)
   - [x] Implement `GetProfile()` for profile lookup in [providers.go](internal/config/providers.go:56-67)
   - [x] ~~Project config removed (global config only)~~

2. **Basic Context Gathering** (`internal/core`) ✅
   - [x] Implement `CollectGitInfo()` in [git.go](internal/core/git.go:9-38)
   - [x] Implement `BuildContext()` in [context.go](internal/core/context.go:37-75)
   - [x] Implement shell history collection in [history.go](internal/core/history.go:11-57)
   - [x] Implement environment variable filtering in [context.go](internal/core/context.go:77-112)
   - [x] Support for both bash and zsh history formats

3. **OpenRouter Integration** (`internal/ai`) ✅
   - [x] Implement OpenRouter HTTP client in [openrouter.go](internal/ai/openrouter.go:106-173)
   - [x] Implement prompt construction in [prompts.go](internal/ai/prompts.go)
   - [x] Implement response parsing in [prompts.go](internal/ai/prompts.go:101-176)
   - [x] Add proper error handling and timeouts
   - [x] Implement provider factory in [provider.go](internal/ai/provider.go:11-26)
   - [x] Support for Suggest and Explain operations

4. **CLI Implementation** (`cmd/linesense`)
   - [ ] Add flag parsing (consider using `github.com/spf13/cobra`)
   - [ ] Implement `runSuggest()` properly in [main.go](cmd/linesense/main.go:35)
   - [ ] Implement `runExplain()` properly in [main.go](cmd/linesense/main.go:75)
   - [ ] Add proper error messages

### Phase 2: Safety & Usability (v0.1-beta)

5. **Safety Filters** (`internal/core`)
   - [ ] Implement `ApplySafetyFilters()` in [safety.go](internal/core/safety.go:18)
   - [ ] Implement `ClassifyRisk()` in [safety.go](internal/core/safety.go:25)
   - [ ] Implement `IsBlocked()` in [safety.go](internal/core/safety.go:32)
   - [ ] Add pattern matching for safety rules

6. **Shell Integration Improvements** (`scripts`)
   - [ ] Improve JSON parsing in bash/zsh scripts
   - [ ] Add config-driven keybindings
   - [ ] Better error handling in shell integration
   - [ ] Add visual feedback for suggestions

### Phase 3: Advanced Features (v0.2)

7. **Usage Logging** (`internal/core`)
   - [ ] Implement `LogUsage()` in [usage.go](internal/core/usage.go:16)
   - [ ] Implement `BuildUsageSummary()` in [usage.go](internal/core/usage.go:23)
   - [ ] Add usage pattern analysis

8. **Command Explanation**
   - [ ] Implement explanation prompt construction
   - [ ] Add formatted output for explanations
   - [ ] Add explanation UI in shell integration

9. **Environment Collection**
   - [ ] Implement environment variable filtering
   - [ ] Add configurable env variable include/exclude lists

## Dependencies to Add

Add these to `go.mod`:

```bash
go get github.com/BurntSushi/toml          # TOML parsing
go get github.com/spf13/cobra              # CLI framework (optional)
go get github.com/spf13/viper              # Config management (optional)
```

## Key Implementation Notes

### Configuration Loading

Use standard XDG directories:
```go
configDir := os.Getenv("XDG_CONFIG_HOME")
if configDir == "" {
    configDir = filepath.Join(os.Getenv("HOME"), ".config")
}
configPath := filepath.Join(configDir, "linesense", "config.toml")
```

### Git Information

Use `git` command directly rather than libgit2:
```go
cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
// Check output for git repo detection
```

### Shell History

Locations:
- Bash: `~/.bash_history` or `$HISTFILE`
- Zsh: `~/.zsh_history` or `$HISTFILE`

### OpenRouter API

Endpoint: `https://openrouter.ai/api/v1/chat/completions`

Request format:
```json
{
  "model": "openrouter/openai/gpt-4o-mini",
  "messages": [
    {"role": "system", "content": "You are a shell command assistant..."},
    {"role": "user", "content": "deploy backend to staging"}
  ]
}
```

Response parsing: Extract suggestions from model output and structure as JSON.

### Prompt Engineering

The prompt should include:
1. Shell type and current working directory
2. Current input line
3. Recent command history (if enabled)
4. Git context (branch, status, remotes)
5. Project presets (if in a project with config)
6. Usage patterns (frequently used commands)

Example prompt structure:
```
You are a shell command assistant. Suggest a complete command based on this context:

Shell: bash
CWD: /home/user/myproject
Current input: "deploy backend to staging"

Git:
- Branch: main
- Status: clean
- Remotes: origin (github.com/user/myproject)

Recent commands:
- git pull
- npm run build
- docker-compose up -d

Project presets:
- deploy_staging: kubectl apply -f k8s/staging && kubectl rollout status...

Provide a JSON response with suggestions.
```

## Testing Strategy

1. **Unit Tests**: Test each package independently
   - Config loading with various TOML inputs
   - Context gathering with mocked git/history
   - Safety filters with test patterns
   - Prompt construction

2. **Integration Tests**: Test end-to-end flows
   - Full suggest workflow
   - Full explain workflow
   - Error handling paths

3. **Gherkin Scenarios**: Implement BDD tests based on PRD scenarios
   - Use `github.com/cucumber/godog` for Gherkin tests
   - Implement all scenarios from [PRD.md](PRD.md)

## Security Considerations

1. **API Key Handling**
   - Never log API keys
   - Read from environment only
   - Consider using keyring for storage

2. **Command Filtering**
   - Always apply safety filters before displaying suggestions
   - Never auto-execute high-risk commands
   - Maintain comprehensive denylist

3. **Environment Variables**
   - Filter sensitive env vars (API keys, passwords)
   - Use allowlist approach rather than blocklist

4. **Usage Logging**
   - Keep logs local only
   - Don't send usage data to external services
   - Make logging opt-in (or clearly documented)

## Performance Optimization

1. **Caching**
   - Cache git info for short periods (10s)
   - Cache project config until file changes
   - Cache provider configuration

2. **Async Operations**
   - Make API calls with reasonable timeouts (5-10s)
   - Show loading indicators in shell

3. **History Parsing**
   - Only read last N lines of history file
   - Use efficient file reading (tail)

## Documentation Needed

1. Installation guide with shell-specific instructions
2. Configuration reference for all TOML options
3. Project preset examples for common stacks
4. Troubleshooting guide
5. Contributing guide
6. API documentation for provider interface

## Release Checklist

Before v0.1 release:
- [ ] All Phase 1 features implemented
- [ ] Basic tests passing
- [ ] Example configs working
- [ ] Shell integration tested on bash and zsh
- [ ] README updated with installation instructions
- [ ] License chosen and added
- [ ] GitHub releases set up

## Next Steps

Start with Phase 1, item 1 (Configuration Loading). This is the foundation for everything else.

Suggested first task:
```bash
# Implement LoadConfig in internal/config/config.go
go get github.com/BurntSushi/toml
# Then implement the XDG directory resolution and TOML parsing
```
