# Testing Guide

LineSense has comprehensive test coverage across all core modules, ensuring reliability and maintainability.

## Overview

- **107 comprehensive tests**
- **90.7% coverage** on core business logic
- All tests passing ✅
- Enterprise-grade quality standards

## Quick Start

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # View in browser
```

## Test Coverage by Module

### Core Module (90.7% coverage)

**What's Covered:**
- Context building and gathering
- Git repository information collection
- Shell history parsing (bash and zsh)
- Environment variable filtering
- Safety filters and risk classification
- Command validation

**Test Files:**
- `internal/core/context_test.go` - Context gathering tests
- `internal/core/safety_test.go` - Safety filter tests

**Key Test Cases:**
```bash
✅ BuildContext with all features enabled
✅ BuildContext with features disabled
✅ Environment variable filtering (sensitive data protection)
✅ Git info collection (clean, dirty, untracked, deleted, added files)
✅ Shell history parsing (bash, zsh, malformed, blank lines)
✅ Risk classification (low/medium/high)
✅ Command blocking via denylist
✅ Integration tests with real git repos
```

### Config Module (84.8% coverage)

**What's Covered:**
- Configuration file loading and parsing
- XDG Base Directory specification compliance
- Provider profile management
- Default value handling
- Error handling for invalid configs

**Test Files:**
- `internal/config/config_test.go` - Main config tests
- `internal/config/providers_test.go` - Provider config tests

**Key Test Cases:**
```bash
✅ Config loading from TOML files
✅ XDG_CONFIG_HOME and HOME fallback
✅ Invalid TOML syntax handling
✅ Provider profile selection
✅ Default profile fallback
✅ Missing profile error handling
```

### AI Module (66.1% coverage)

**What's Covered:**
- Prompt building for suggestions and explanations
- Response parsing and cleanup
- Risk assessment algorithm
- Provider factory creation
- API key validation

**Test Files:**
- `internal/ai/prompts_test.go` - Prompt and parsing tests
- `internal/ai/provider_test.go` - Provider factory tests

**Key Test Cases:**
```bash
✅ System and user prompt generation
✅ Context inclusion in prompts
✅ Markdown cleanup from responses
✅ Structured vs unstructured explanation parsing
✅ Risk level extraction from responses
✅ Provider creation with valid API keys
✅ Error handling for missing API keys
```

## Test Organization

### Unit Tests

Each module has focused unit tests that verify:
- **Happy path** - Normal operation with valid inputs
- **Edge cases** - Boundary conditions, empty inputs, malformed data
- **Error handling** - Missing files, invalid configs, network failures
- **Integration** - Multiple components working together

### Test Fixtures

Tests use temporary directories and isolated environments:
```go
tmpDir := t.TempDir()  // Automatically cleaned up
```

Environment variables are saved and restored:
```go
originalXDG := os.Getenv("XDG_CONFIG_HOME")
defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
```

### Table-Driven Tests

Many tests use table-driven approaches for clarity:
```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"valid case", "input", "expected"},
    {"edge case", "", ""},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

## Running Specific Tests

```bash
# Run tests for a specific module
go test ./internal/core/...

# Run a specific test
go test -run TestBuildContext ./internal/core/

# Run tests matching a pattern
go test -run "TestCollectGit.*" ./internal/core/

# Run tests with race detection
go test -race ./...

# Run tests with timeout
go test -timeout 30s ./...
```

## Coverage Goals

| Module | Current | Target | Status |
|--------|---------|--------|--------|
| Core | 90.7% | >80% | ✅ Exceeded |
| Config | 84.8% | >80% | ✅ Exceeded |
| AI | 66.1% | >60% | ✅ Exceeded |

## What's Not Covered

The following code is intentionally not covered by unit tests:

### Engine (0%)
The engine requires integration with live AI providers and is better tested via E2E tests with actual API calls.

### Usage Logging (0%)
Not yet implemented - marked with `panic("not implemented")`.

### OpenRouter HTTP Methods (0%)
Require actual HTTP calls or complex mocking. Better covered by E2E tests.

## Writing New Tests

### Test Checklist

When adding new functionality, ensure tests cover:

- [ ] **Happy path** - Normal operation
- [ ] **Nil/empty inputs** - Graceful handling
- [ ] **Invalid inputs** - Proper error messages
- [ ] **Boundary conditions** - Min/max values
- [ ] **Integration** - Works with other components
- [ ] **Cleanup** - No side effects, temp files cleaned

### Example Test Structure

```go
func TestNewFeature(t *testing.T) {
    // Setup
    tmpDir := t.TempDir()

    // Create test data
    testData := createTestData()

    // Execute
    result, err := NewFeature(testData)

    // Verify
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Test Naming Conventions

- `TestFunctionName` - Basic functionality
- `TestFunctionName_SpecificCase` - Specific scenario
- `TestFunctionName_ErrorCondition` - Error handling

Examples:
- `TestBuildContext` - Basic context building
- `TestBuildContext_DisabledFeatures` - With features off
- `TestLoadConfig_MissingFile` - Error case

## Shell Integration Tests

Shell integration is tested via `test_shell_integration.sh`:

```bash
./test_shell_integration.sh
```

This verifies:
- ✅ Bash integration script loads correctly
- ✅ Zsh integration script loads correctly
- ✅ JSON parsing works (jq and fallback)
- ✅ Suggest command integration
- ✅ Explain command integration

## End-to-End Tests

E2E tests require a valid API key:

```bash
# Set your API key
export OPENROUTER_API_KEY="sk-or-v1-..."

# Run E2E tests
go test -v ./test_e2e.go
```

E2E tests verify:
- Complete suggest workflow
- Complete explain workflow
- Real API integration
- Full context gathering pipeline

## Continuous Integration

### GitHub Actions (Planned)

```yaml
- Run tests on push/PR
- Generate coverage reports
- Test multiple Go versions (1.21, 1.22, latest)
- Test on multiple platforms (Linux, macOS)
- Upload coverage to codecov.io
```

## Coverage Improvement

To improve coverage further:

1. **Identify gaps:**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out | grep -v "100.0%"
   ```

2. **Add targeted tests:**
   Focus on uncovered functions shown in the report

3. **Verify improvement:**
   ```bash
   go test -cover ./...
   ```

## Best Practices

### ✅ Do's

- ✅ Use `t.TempDir()` for temporary directories
- ✅ Save and restore environment variables
- ✅ Use table-driven tests for multiple cases
- ✅ Test both success and error paths
- ✅ Keep tests isolated and independent
- ✅ Use clear, descriptive test names
- ✅ Add comments for complex test logic

### ❌ Don'ts

- ❌ Don't rely on external state
- ❌ Don't use hardcoded paths
- ❌ Don't skip cleanup in defers
- ❌ Don't test implementation details
- ❌ Don't make tests depend on each other
- ❌ Don't use sleep() for synchronization
- ❌ Don't commit test files to git

## Debugging Tests

```bash
# Run with verbose output
go test -v ./internal/core/

# Print test output (including t.Log)
go test -v -test.v ./...

# Run with stack traces
go test -v -trace=trace.out ./...

# Use delve debugger
dlv test ./internal/core/ -- -test.run TestBuildContext
```

## Performance Testing

```bash
# Run benchmarks
go test -bench=. ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

## Contributing

When submitting PRs:

1. Ensure all tests pass: `go test ./...`
2. Maintain or improve coverage
3. Add tests for new features
4. Update this document if adding new test patterns

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Coverage Guide](https://go.dev/blog/cover)
- [Testing Best Practices](https://github.com/golang/go/wiki/CodeReviewComments#tests)

---

**Last Updated:** 2025-11-14
**Test Count:** 107
**Coverage:** Core 90.7% | Config 84.8% | AI 66.1%
