# Contributing to LineSense

Thank you for your interest in contributing to LineSense! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Review Process](#review-process)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/LineSense.git
   cd LineSense
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/traves/LineSense.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- A shell (bash or zsh) for testing integrations
- Optional: OpenRouter API key for E2E tests

### Install Dependencies

```bash
# Install the project
go install ./cmd/linesense

# Verify installation
linesense --version
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run shell integration tests
./test_shell_integration.sh
```

See [docs/TESTING.md](docs/TESTING.md) for detailed testing information.

### Continuous Integration

All pull requests run through GitHub Actions CI:
- Tests on Linux and macOS
- Multiple Go versions (1.21, 1.22, 1.23)
- Code quality checks (golangci-lint, gofmt, go vet)
- Security scanning (Gosec)
- Shell integration tests

See [docs/CI_CD.md](docs/CI_CD.md) for details on the CI/CD pipeline.

## Making Changes

### Before You Start

1. **Check existing issues** - Your idea might already be discussed
2. **Open an issue** - Discuss major changes before implementing
3. **Create a branch** - Use a descriptive name:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

### Branch Naming Conventions

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `test/` - Test improvements
- `refactor/` - Code refactoring
- `chore/` - Maintenance tasks

## Testing

**All code changes must include tests.**

### Test Requirements

- [ ] Unit tests for new functions
- [ ] Integration tests for new features
- [ ] Update existing tests if behavior changes
- [ ] Maintain or improve code coverage
- [ ] All tests must pass

### Writing Tests

See [docs/TESTING.md](docs/TESTING.md) for comprehensive testing guidelines.

Quick checklist:
- Use `t.TempDir()` for temporary directories
- Save and restore environment variables
- Use table-driven tests for multiple cases
- Test both success and error paths
- Add clear test names and comments

### Coverage Standards

- Core logic: >80% coverage
- New features: >80% coverage
- Bug fixes: Must include regression test

## Submitting Changes

### Pull Request Process

1. **Update your branch** with latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests** locally:
   ```bash
   go test ./...
   go test -cover ./...
   ```

3. **Push your changes**:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request** on GitHub with:
   - Clear title describing the change
   - Description of what changed and why
   - Reference to related issues (e.g., "Fixes #123")
   - Screenshots/examples if applicable

### PR Checklist

Before submitting, ensure:

- [ ] Tests pass locally
- [ ] Code follows project style (run `gofmt`)
- [ ] Tests are included for new functionality
- [ ] Documentation is updated if needed
- [ ] CHANGELOG.md is updated (if applicable)
- [ ] Commit messages are clear and descriptive
- [ ] No merge conflicts with main branch

## Code Style

### Go Code Style

We follow standard Go conventions:

```bash
# Format your code
gofmt -w .

# Run linters
go vet ./...
```

### Style Guidelines

- **Naming**: Use clear, descriptive names
  - Exported: `BuildContext`, `CollectHistory`
  - Unexported: `parseHistoryLine`, `filterEnvironment`

- **Comments**: Add comments for:
  - Exported functions (godoc format)
  - Complex logic
  - Non-obvious behavior

- **Error Handling**: Always handle errors
  ```go
  // Good
  result, err := DoSomething()
  if err != nil {
      return fmt.Errorf("doing something: %w", err)
  }

  // Bad
  result, _ := DoSomething()
  ```

- **Keep functions small**: Aim for < 50 lines
- **Single responsibility**: One function, one job

### Documentation

- Update README.md for user-facing changes
- Update relevant docs/ files for feature changes
- Add godoc comments for exported functions
- Include examples in comments when helpful

## Commit Messages

Write clear, descriptive commit messages:

### Format

```
type(scope): short description

Longer explanation if needed. Wrap at 72 characters.
Include motivation for the change and contrast with previous behavior.

Fixes #123
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Examples

```
feat(ai): add support for custom temperature settings

Allows users to configure LLM temperature per profile.
This gives more control over suggestion creativity.

Fixes #45
```

```
fix(context): handle empty shell history gracefully

Previously crashed on empty history files. Now returns
empty slice and continues normally.

Fixes #78
```

## Review Process

### What to Expect

1. **Automated Checks**: CI/CD runs tests automatically
2. **Code Review**: Maintainers review your code
3. **Feedback**: You may receive change requests
4. **Approval**: Once approved, maintainers will merge

### Review Timeline

- Small fixes: 1-3 days
- New features: 3-7 days
- Major changes: 1-2 weeks

### Responding to Feedback

- Be open to suggestions
- Ask questions if unclear
- Make requested changes in new commits
- Don't force-push after review starts

## Areas for Contribution

### Good First Issues

Look for issues labeled `good-first-issue`:
- Documentation improvements
- Test coverage improvements
- Small bug fixes
- Example additions

### High Priority

- Shell integration improvements (fish, nushell)
- Additional AI provider support
- Performance optimizations
- Error message improvements

### Feature Requests

- Local model support (Ollama)
- Custom prompt templates
- Usage analytics dashboard
- Mobile support (Termux)

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open an issue with reproduction steps
- **Security**: See SECURITY.md for reporting vulnerabilities
- **Chat**: Join our community (link TBD)

## Recognition

Contributors are recognized in:
- CHANGELOG.md for their contributions
- GitHub contributors page
- Release notes for significant features

Thank you for contributing to LineSense! 🎉
