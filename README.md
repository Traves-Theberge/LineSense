# LineSense

AI-powered shell command autocomplete and explanation tool.

## Overview

LineSense is an intelligent shell assistant that provides context-aware command suggestions and explanations. It integrates seamlessly with bash and zsh, learning from your usage patterns.

## Features

- **Context-Aware Suggestions**: Uses git info, shell history, and environment
- **Safety First**: Risk classification and configurable denylists
- **Multi-Shell Support**: Works with bash and zsh
- **Adaptive Learning**: Learns from your command usage patterns
- **OpenRouter Integration**: Powered by state-of-the-art LLMs via OpenRouter

## Project Structure

```
.
├── cmd/
│   └── linesense/          # Main CLI binary
│       └── main.go
├── internal/
│   ├── config/             # Configuration loading
│   │   ├── config.go       # Global config
│   │   └── providers.go    # Provider/model config
│   ├── core/               # Core engine
│   │   ├── context.go      # Context gathering
│   │   ├── engine.go       # Main suggest/explain engine
│   │   ├── git.go          # Git integration
│   │   ├── history.go      # Shell history
│   │   ├── safety.go       # Safety filters
│   │   └── usage.go        # Usage logging
│   └── ai/                 # AI provider implementations
│       ├── provider.go     # Provider factory
│       └── openrouter.go   # OpenRouter implementation
├── scripts/
│   ├── linesense.bash      # Bash integration
│   └── linesense.zsh       # Zsh integration
├── examples/
│   ├── config.toml         # Example global config
│   └── providers.toml      # Example providers config
└── PRD.md                  # Product requirements document
```

## Quick Start

```bash
# 1. Clone and build
git clone <repo-url>
cd LineSense
go install ./cmd/linesense

# 2. Set up configuration
mkdir -p ~/.config/linesense
cp examples/*.toml ~/.config/linesense/

# 3. Create .env file with your API key (for development)
echo "OPENROUTER_API_KEY=your-key-here" > .env

# 4. Try it out!
linesense suggest --line "list files"
linesense explain --line "rm -rf /"
```

## Installation

### Prerequisites

- Go 1.21 or later
- An OpenRouter API key (get one at https://openrouter.ai)

### Build from Source

```bash
# Clone the repository
git clone <repo-url>
cd LineSense

# Build the binary
go build -o linesense ./cmd/linesense

# Or install to $GOPATH/bin
go install ./cmd/linesense
```

### Configuration

1. Create configuration directory:
```bash
mkdir -p ~/.config/linesense
```

2. Copy example configurations:
```bash
cp examples/config.toml ~/.config/linesense/
cp examples/providers.toml ~/.config/linesense/
```

3. Set your OpenRouter API key:
```bash
# For development: create a .env file
echo "OPENROUTER_API_KEY=your-key-here" > .env

# For production: set environment variable
export OPENROUTER_API_KEY="your-api-key-here"
```

### Shell Integration

For **bash**, add to your `~/.bashrc`:
```bash
source /path/to/linesense/scripts/linesense.bash
```

For **zsh**, add to your `~/.zshrc`:
```bash
source /path/to/linesense/scripts/linesense.zsh
```

## Usage

### CLI Commands

LineSense provides two main commands:

#### Suggest Command
Generate command suggestions based on natural language input:
```bash
linesense suggest --line "list files"
linesense suggest --line "find large files" --cwd /var/log
linesense suggest --line "git com" --model openai/gpt-4o
```

Output:
```json
{
  "suggestions": [
    {
      "command": "ls -al",
      "risk": "low",
      "explanation": "Suggested based on: list files",
      "source": "llm"
    }
  ]
}
```

#### Explain Command
Get detailed explanations of commands:
```bash
linesense explain --line "rm -rf /"
linesense explain --line "docker ps -a"
```

Output:
```json
{
  "summary": "The `rm -rf /` command recursively and forcefully removes all files...",
  "risk": "high",
  "notes": [
    "This command is extremely dangerous...",
    "Important Flags: -r: Recursively delete..."
  ]
}
```

### Shell Integration (Coming Soon)

Once shell integration is implemented:
- Press `Ctrl+Space` (default) while typing to get AI-powered suggestions
- Press `Ctrl+E` to get an explanation of the current command

## Configuration

### Global Config (`~/.config/linesense/config.toml`)

Controls shell integration, keybindings, context gathering, and safety rules.

### Providers Config (`~/.config/linesense/providers.toml`)

Configures AI providers and models. Supports multiple profiles (default, fast, smart).

## Development Status

**Phase 1: Core Infrastructure & CLI** ✅ **COMPLETE**

LineSense is now fully functional for CLI usage! The following features are implemented and tested:

### ✅ Completed Features

1. **Configuration Loading** - Full TOML config support
   - XDG config directory resolution
   - Global config (`~/.config/linesense/config.toml`)
   - Provider profiles (`~/.config/linesense/providers.toml`)
   - .env file support for development

2. **Context Gathering** - Rich contextual awareness
   - Git repository detection (branch, status, remotes)
   - Shell history parsing (bash and zsh)
   - Environment variable filtering (security-aware)
   - Current working directory tracking

3. **AI Integration** - OpenRouter API
   - HTTP client with authentication and timeouts
   - Context-aware prompt construction
   - Response parsing and structuring
   - Risk assessment (low/medium/high)

4. **CLI** - Full-featured command-line interface
   - `suggest` command for generating suggestions
   - `explain` command for explaining commands
   - Flag parsing (--shell, --line, --cwd, --model)
   - Structured JSON output
   - Auto-detection of shell type

5. **Testing** - Comprehensive test coverage
   - End-to-end test suite
   - All operations verified working
   - Risk classification tested

### 🚧 Phase 2: Shell Integration & Safety (Next Steps)

1. **Shell Integration Scripts**
   - Bash integration with keybindings
   - Zsh integration with ZLE widgets
   - Config-driven keybindings
   - Better JSON parsing (use jq)

2. **Safety Filters**
   - Configurable command denylists
   - Enhanced risk classification
   - Command blocking for dangerous patterns

3. **Usage Logging**
   - Track accepted suggestions
   - Build usage summaries
   - Learn from patterns

See [PROGRESS.md](PROGRESS.md) for detailed implementation status.

## Testing

The PRD includes comprehensive Gherkin scenarios that should be used to drive test implementation.

## License

TODO

## Contributing

TODO
