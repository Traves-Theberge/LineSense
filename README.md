# LineSense

```text
   __    _            ____                       
  / /   (_)_ __   ___/ ___|  ___ _ __  ___  ___  
 / /    | | '_ \ / _ \___ \ / _ \ '_ \/ __|/ _ \ 
/ /___  | | | | |  __/___) |  __/ | | \__ \  __/ 
\____/  |_|_| |_|\___|____/ \___|_| |_|___/\___| 
```

**AI-Powered Shell Assistant** (v0.6.1)

[![CI](https://github.com/Traves-Theberge/LineSense/workflows/CI/badge.svg)](https://github.com/Traves-Theberge/LineSense/actions)
[![Tests](https://img.shields.io/badge/tests-107%20passing-success)](.)
[![Coverage](https://img.shields.io/badge/coverage-90.7%25-brightgreen)](.)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](.)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Traves-Theberge/LineSense)](https://goreportcard.com/report/github.com/Traves-Theberge/LineSense)

## Overview

LineSense is an intelligent shell assistant that provides context-aware command suggestions and explanations. It integrates seamlessly with bash and zsh, learning from your usage patterns.

## Features

- **🎯 OS-Aware Suggestions** (NEW in v0.5.2): Automatically detects your OS, distribution, and package manager
  - Get `brew install` on macOS, `apt install` on Ubuntu, `pacman -S` on Arch
  - Smart command suggestions tailored to your system
  - Zero configuration required
- **🎨 Beautiful Terminal UI**: Styled output with colors, borders, and dynamic width adjustment
- **🔄 Dual Output Modes**: Pretty format for humans, JSON for scripting (`--format` flag)
- **🌍 Global Instructions**: Define personal rules that apply everywhere (e.g., "Always use bat")
- **📁 Project Context**: Add `.linesense_context` files for directory-specific AI knowledge
- **⚡ Loading Indicators**: Animated spinner while AI processes your request
- **🧠 Context-Aware Suggestions**: Uses git info, shell history, environment, and OS context
- **🛡️ Safety First**: Risk classification and configurable denylists
- **🐚 Multi-Shell Support**: Works with bash and zsh
- **🚀 OpenRouter Integration**: Powered by state-of-the-art LLMs via OpenRouter
- **📏 Responsive Design**: Output automatically adapts to terminal width
- **💡 Smart Explanations**: Each suggestion includes a brief 5-10 word explanation
- **🔢 Multiple Suggestions**: Get 3-5 alternative command options for every request

## Project Structure

```
.
├── cmd/
│   └── linesense/          # Main CLI binary
│       ├── main.go         # CLI entry point
│       └── ui.go           # Terminal UI (Lipgloss/Bubbletea)
├── internal/
│   ├── config/             # Configuration loading
│   │   ├── config.go       # Global config
│   │   └── providers.go    # Provider/model config
│   ├── core/               # Core engine
│   │   ├── context.go      # Context gathering
│   │   ├── engine.go       # Main suggest/explain engine
│   │   ├── git.go          # Git integration
│   │   ├── history.go      # Shell history
│   │   ├── osdetect.go     # OS & package manager detection
│   │   ├── safety.go       # Safety filters
│   │   └── usage.go        # Usage logging
│   └── ai/                 # AI provider implementations
│       ├── provider.go     # Provider factory
│       ├── prompts.go      # AI prompts & parsing
│       └── openrouter.go   # OpenRouter implementation
├── scripts/
│   ├── linesense.bash      # Bash integration
│   └── linesense.zsh       # Zsh integration
├── examples/
│   ├── config.toml         # Example global config
│   └── providers.toml      # Example providers config
├── docs/                   # Comprehensive documentation
│   ├── INSTALLATION.md     # Installation guide
│   ├── CONFIGURATION.md    # Configuration reference
│   ├── SECURITY.md         # Security features
│   ├── API.md              # CLI reference
│   ├── TESTING.md          # Testing guide
│   └── CI_CD.md            # CI/CD and release process
└── docs/               # Documentation (including PRD)
```

## Quick Start

### Automated Installation (Recommended)

The easiest way to install LineSense is using the automated installation script:

```bash
# Download and run the installer
curl -fsSL https://raw.githubusercontent.com/traves-theberge/LineSense/main/install.sh | bash
```

This will:
- ✅ Build and install LineSense
- ✅ Set up shell integration (bash/zsh)
- ✅ Initialize configuration
- ✅ Guide you through API key setup

**Then restart your shell and set your API key:**
```bash
# Restart your shell or reload config
source ~/.bashrc  # or ~/.zshrc

# Set your OpenRouter API key (interactive)
linesense config set-key
```

### Manual Installation

If you prefer to install manually or want more control:

#### Prerequisites

- Go 1.21 or later
- Git
- An OpenRouter API key (get one at https://openrouter.ai)

#### Install from Source

```bash
# 1. Clone the repository
git clone https://github.com/Traves-Theberge/LineSense.git
cd LineSense

# 2. Build and install
go install ./cmd/linesense

# 3. Initialize configuration
linesense config init

# 4. Set your OpenRouter API key
linesense config set-key

# 5. Set up shell integration
# For bash, add to ~/.bashrc:
echo '[ -f "$HOME/.config/linesense/shell/linesense.bash" ] && source "$HOME/.config/linesense/shell/linesense.bash"' >> ~/.bashrc

# For zsh, add to ~/.zshrc:
echo '[ -f "$HOME/.config/linesense/shell/linesense.zsh" ] && source "$HOME/.config/linesense/shell/linesense.zsh"' >> ~/.zshrc

# 6. Reload your shell
source ~/.bashrc  # or ~/.zshrc
```

#### Alternative: Go Install

Install directly from the repository:

```bash
go install github.com/Traves-Theberge/LineSense/cmd/linesense@latest
linesense config init
linesense config set-key
```

### Verify Installation

```bash
# Check version
linesense --version

# View configuration
linesense config show

# Try it out!
linesense suggest --line "list files sorted by size"
linesense explain --line "docker ps -a"
```

## Usage

### CLI Commands

LineSense provides two main commands:

#### Suggest Command
Generate command suggestions based on natural language input:
```bash
# Pretty output (default) - beautiful terminal UI
linesense suggest --line "list files"

# JSON output for scripting
linesense suggest --line "list files" --format json

# Advanced options
linesense suggest --line "find large files" --cwd /var/log
linesense suggest --line "git com" --model openai/gpt-4o
```

**Pretty Output** (default):
```
💡 Command Suggestions
────────────────────────────────────────────────────────────────────────────────

1. ls -lhS
   ✓ Risk: low
   List files sorted by size in human-readable format

2. find . -type f -exec du -h {} + | sort -rh | head -20
   ✓ Risk: low
   Find and display 20 largest files

3. du -ah . | sort -rh | head -20
   ✓ Risk: low
   Show disk usage sorted by size
```

**OS-Aware Examples:**

Ubuntu user types "install nginx":
```
1. sudo apt install nginx
   ⚠ Risk: medium
   Install nginx web server using apt

2. sudo apt install nginx-full
   ⚠ Risk: medium
   Install nginx with all available modules
```

macOS user types "install nginx":
```
1. brew install nginx
   ✓ Risk: low
   Install nginx web server using Homebrew

2. brew install nginx --with-pcre
   ✓ Risk: low
   Install nginx with PCRE support
```

**JSON Output** (`--format json`):
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
# Pretty output (default) - beautiful terminal UI
linesense explain --line "docker ps -a"

# JSON output for scripting
linesense explain --line "docker ps -a" --format json
```

**Pretty Output** (default):
```
📖 Command Explanation
────────────────────────────────────────────────────────────────────────────────

╭──────────────────────────────────────────────────────────────────────────────╮
│ Summary                                                                      │
│                                                                              │
│ Lists all Docker containers (running and stopped) with their details        │
╰──────────────────────────────────────────────────────────────────────────────╯

╭──────────────────────────────────────────────────────────────────────────────╮
│ ✓ Risk Level: low                                                           │
╰──────────────────────────────────────────────────────────────────────────────╯

Details

╭──────────────────────────────────────────────────────────────────────────────╮
│ What it does                                                                 │
│ - Shows container ID, image, command, status, ports, and names              │
│ - The -a flag includes stopped containers                                   │
╰──────────────────────────────────────────────────────────────────────────────╯
```

**JSON Output** (`--format json`):
```json
{
  "summary": "Lists all Docker containers (running and stopped)...",
  "risk": "low",
  "notes": [
    "What it does",
    "- Shows container ID, image, command..."
  ]
}
```

### Shell Integration

LineSense provides interactive shell integration for both bash and zsh. The integration loads silently in the background - no startup messages or notifications.

**Default Keybindings:**
- Press `Ctrl+Space` to get AI-powered command suggestions
- Press `Ctrl+X` to get an explanation of the current command

**Customization:**
You can customize keybindings by setting environment variables before sourcing the script:

```bash
# In your ~/.bashrc or ~/.zshrc
export LINESENSE_SUGGEST_KEY="\C-t"      # Change suggest to Ctrl+T
export LINESENSE_EXPLAIN_KEY="\C-x\C-h"  # Change explain to Ctrl+X Ctrl+H
source ~/.config/linesense/shell/linesense.bash
```

**Features:**
- 💡 Smart suggestions - handles typos and provides intent-based alternatives
- 📖 Detailed explanations - comprehensive command breakdowns with risk assessment
- 🧠 Context-aware - uses current directory, git status, and shell history
- 🖥️ OS-aware - detects your operating system, distribution, and package manager for tailored suggestions

## Configuration

LineSense uses a TOML configuration file located at `~/.config/linesense/config.toml`.

### Global Instructions

You can define global rules that the AI should always follow. This is useful for enforcing personal preferences or tools.

To edit your configuration, run:
```bash
linesense config edit
```

Then add your instructions to the `[context]` section:

```toml
[context]
global_instructions = """
- Always prefer 'bat' over 'cat'
- Use 'podman' instead of 'docker'
- When suggesting git commands, always include '--verbose'
"""
```

### Project-Specific Context

For project-specific rules, create a `.linesense_context` file in your project's root directory. LineSense will automatically read this file when you are working in that directory.

You can easily create this file using the CLI:

```bash
linesense config init-project
```

**Example `.linesense_context`:**
```text
This project uses a custom CLI tool called 'ops-cli'.
- To build: ops-cli build --env=prod
- To deploy: ops-cli deploy --region=us-east-1
- Never use 'kubectl' directly, always use 'ops-cli k8s' wrapper.
```
