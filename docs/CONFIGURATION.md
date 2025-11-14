# LineSense Configuration Reference

This document provides a complete reference for all LineSense configuration options.

## Table of Contents

- [Configuration Overview](#configuration-overview)
- [Configuration Files](#configuration-files)
  - [Global Config](#global-config-configtoml)
  - [Providers Config](#providers-config-providerstoml)
- [Environment Variables](#environment-variables)
- [CLI Configuration Commands](#cli-configuration-commands)
- [Configuration Examples](#configuration-examples)
- [Advanced Configuration](#advanced-configuration)

## Configuration Overview

LineSense uses a multi-layered configuration system:

1. **Configuration files** (TOML format) in `~/.config/linesense/`
2. **Environment variables** for API keys and runtime overrides
3. **CLI flags** for per-command overrides

**Configuration precedence** (highest to lowest):
1. CLI flags (e.g., `--model`, `--cwd`)
2. Environment variables (e.g., `OPENROUTER_API_KEY`)
3. Configuration files (`config.toml`, `providers.toml`)
4. Built-in defaults

## Configuration Files

Configuration files are located in the XDG config directory:
- Linux/macOS: `~/.config/linesense/`
- Custom: Set `$XDG_CONFIG_HOME` to override

### Global Config (`config.toml`)

The main configuration file controlling LineSense behavior.

**Location:** `~/.config/linesense/config.toml`

**Full Example:**
```toml
# AI Provider Configuration
[ai]
# Which provider profile to use from providers.toml
# Options: "default", "fast", "smart", or any custom profile
provider_profile = "default"

# Context Gathering Configuration
[context]
# Number of shell history entries to include in context
# Higher values provide more context but increase API costs
# Range: 0-1000, Recommended: 20-100
history_length = 50

# Include git repository information (branch, status, remotes)
# Useful for git-related suggestions
include_git = true

# Include environment variables in context
# Variables are filtered through env_allowlist for security
include_env = true

# Which environment variables to include
# Only these variables will be sent to the AI
env_allowlist = [
    "PATH",
    "HOME",
    "USER",
    "SHELL",
    "LANG",
    "LC_ALL",
    "EDITOR",
    "VISUAL"
]

# Safety and Security Configuration
[safety]
# Enable safety filtering for commands
# When true, commands are classified by risk level
enable_filters = true

# Additional patterns to mark as high-risk
# Uses regex syntax - these commands will show ⚠️ warnings
require_confirm_patterns = [
    "format",           # Disk formatting operations
    "encrypt",          # Encryption operations
    "decrypt",          # Decryption operations
    "iptables.*-F",     # Firewall flush
    "ufw.*disable",     # Firewall disable
    "setenforce.*0",    # SELinux permissive mode
]

# Commands to completely block (will not be suggested)
# Uses regex syntax - matched commands are filtered out
denylist = [
    "rm\\s+-rf\\s+/",           # Delete root filesystem
    "dd\\s+if=.*of=/dev/sd",    # Disk operations to physical devices
    "mkfs.*",                   # Filesystem formatting
    ":\\(\\)\\{.*\\};:",        # Fork bombs
]

# Shell Integration Configuration
[shell]
# Enable shell integration features
enable_integration = true

# Show colored output in shell (risk indicators, formatting)
colored_output = true

# Maximum suggestions to show in shell
max_suggestions = 3

# Timeout for AI requests in shell (seconds)
timeout = 10
```

#### Configuration Sections Explained

##### `[ai]` Section

Controls which AI provider and model to use.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `provider_profile` | string | `"default"` | Profile name from `providers.toml` |

##### `[context]` Section

Controls what contextual information is gathered and sent to the AI.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `history_length` | int | `50` | Number of recent shell commands to include |
| `include_git` | bool | `true` | Include git repo info (branch, status, remotes) |
| `include_env` | bool | `true` | Include filtered environment variables |
| `env_allowlist` | array | See example | Which env vars to include |

##### `[safety]` Section

Controls command safety filtering and risk assessment.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enable_filters` | bool | `true` | Enable safety filtering |
| `require_confirm_patterns` | array | `[]` | Additional high-risk patterns (regex) |
| `denylist` | array | `[]` | Commands to completely block (regex) |

##### `[shell]` Section

Controls shell integration behavior.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enable_integration` | bool | `true` | Enable shell integration features |
| `colored_output` | bool | `true` | Use colored output with risk indicators |
| `max_suggestions` | int | `3` | Maximum suggestions to display |
| `timeout` | int | `10` | Request timeout in seconds |

### Providers Config (`providers.toml`)

Defines AI provider profiles with different models and parameters.

**Location:** `~/.config/linesense/providers.toml`

**Full Example:**
```toml
# Default Profile - Balanced performance and quality
[profiles.default]
provider = "openrouter"
model = "openai/gpt-4o-mini"
temperature = 0.3
max_tokens = 500
timeout = 10

# Fast Profile - Quick responses, lower quality
[profiles.fast]
provider = "openrouter"
model = "meta-llama/llama-3.1-8b-instruct:free"
temperature = 0.2
max_tokens = 300
timeout = 5

# Smart Profile - Highest quality, slower
[profiles.smart]
provider = "openrouter"
model = "openai/gpt-4o"
temperature = 0.4
max_tokens = 800
timeout = 15

# Custom Profile - Your own configuration
[profiles.custom]
provider = "openrouter"
model = "anthropic/claude-3.5-sonnet"
temperature = 0.35
max_tokens = 600
timeout = 12
```

#### Profile Options Explained

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `provider` | string | Yes | Always `"openrouter"` (only supported provider) |
| `model` | string | Yes | Model ID from OpenRouter (see [Available Models](#available-models)) |
| `temperature` | float | Yes | Creativity level (0.0-1.0, lower = more focused) |
| `max_tokens` | int | Yes | Maximum response length (100-2000) |
| `timeout` | int | No | Request timeout in seconds (default: 10) |

#### Available Models

Popular models available through OpenRouter:

**Fast & Free:**
- `meta-llama/llama-3.1-8b-instruct:free` - Fast, free, good quality
- `google/gemini-flash-1.5` - Very fast, free tier available

**Balanced:**
- `openai/gpt-4o-mini` - Best price/performance ratio (recommended default)
- `anthropic/claude-3-haiku` - Fast and accurate

**High Quality:**
- `openai/gpt-4o` - Very capable, higher cost
- `anthropic/claude-3.5-sonnet` - Excellent reasoning
- `anthropic/claude-3-opus` - Highest quality

See [OpenRouter Models](https://openrouter.ai/models) for full list and pricing.

## Environment Variables

### Required Variables

#### `OPENROUTER_API_KEY`

Your OpenRouter API key.

**Set using:**
```bash
# Preferred: Use config command (stores in shell RC file)
linesense config set-key

# Manual: Add to ~/.bashrc or ~/.zshrc
export OPENROUTER_API_KEY="sk-or-v1-..."

# Development: Use .env file
echo 'OPENROUTER_API_KEY="sk-or-v1-..."' > .env
```

### Optional Variables

#### `XDG_CONFIG_HOME`

Override config directory location.

**Default:** `~/.config`

**Example:**
```bash
export XDG_CONFIG_HOME="$HOME/.local/config"
# Config will be in ~/.local/config/linesense/
```

#### `LINESENSE_SUGGEST_KEY`

Customize shell keybinding for suggestions.

**Default:** `\C- ` (Ctrl+Space for bash), `^ ` (Ctrl+Space for zsh)

**Example:**
```bash
# Use Ctrl+T instead
export LINESENSE_SUGGEST_KEY="\C-t"  # bash
export LINESENSE_SUGGEST_KEY="^T"    # zsh
```

#### `LINESENSE_EXPLAIN_KEY`

Customize shell keybinding for explanations.

**Default:** `\C-x\C-e` (Ctrl+X Ctrl+E for bash), `^X^E` (Ctrl+X Ctrl+E for zsh)

**Example:**
```bash
# Use Ctrl+H instead
export LINESENSE_EXPLAIN_KEY="\C-h"  # bash
export LINESENSE_EXPLAIN_KEY="^H"    # zsh
```

## CLI Configuration Commands

### `linesense config init`

Initialize configuration with interactive setup wizard.

**Usage:**
```bash
linesense config init
```

**Actions:**
1. Creates `~/.config/linesense/` directory
2. Copies default `config.toml` and `providers.toml`
3. Prompts to overwrite if files exist
4. Displays next steps

### `linesense config set-key`

Set OpenRouter API key securely.

**Usage:**
```bash
# Interactive (more secure - key not in shell history)
linesense config set-key

# Direct (convenience - key appears in shell history)
linesense config set-key sk-or-v1-...
```

**Actions:**
1. Auto-detects shell (bash or zsh)
2. Adds/updates `OPENROUTER_API_KEY` in shell RC file
3. Sets file permissions to 0600
4. Prompts before overwriting existing key

**Security:**
- API key stored in `~/.bashrc` or `~/.zshrc`
- NOT stored in config files
- File permissions set to 0600 (owner read/write only)

### `linesense config set-model`

Change the default model in your config.

**Usage:**
```bash
linesense config set-model MODEL_ID
```

**Example:**
```bash
# Switch to GPT-4o
linesense config set-model openai/gpt-4o

# Switch to free Llama model
linesense config set-model meta-llama/llama-3.1-8b-instruct:free
```

**Actions:**
1. Updates `model` field in `providers.toml` default profile
2. Displays popular model options

### `linesense config show`

Display current configuration.

**Usage:**
```bash
linesense config show
```

**Output:**
```
📋 LineSense Configuration
==========================

Config directory: /home/user/.config/linesense

API Key: sk-or-v1...cf28 ✓

Configuration:
  Provider profile: default
  History length: 50
  Include git: true
  Include env: true

Provider settings:
  Model: openai/gpt-4o-mini
  Temperature: 0.3
  Max tokens: 500
```

## Configuration Examples

### Minimal Configuration

For basic usage with minimal context:

**config.toml:**
```toml
[ai]
provider_profile = "fast"

[context]
history_length = 10
include_git = false
include_env = false

[safety]
enable_filters = true
```

### Maximum Context Configuration

For detailed, context-aware suggestions:

**config.toml:**
```toml
[ai]
provider_profile = "smart"

[context]
history_length = 100
include_git = true
include_env = true
env_allowlist = [
    "PATH", "HOME", "USER", "SHELL",
    "EDITOR", "VISUAL", "LANG", "LC_ALL",
    "PWD", "OLDPWD", "TERM",
    "GIT_*", "NODE_*", "PYTHON_*"
]

[safety]
enable_filters = true
```

### High Security Configuration

For maximum safety and command filtering:

**config.toml:**
```toml
[ai]
provider_profile = "default"

[context]
history_length = 20
include_git = true
include_env = false  # Don't send env vars

[safety]
enable_filters = true

require_confirm_patterns = [
    "rm",
    "mv",
    "dd",
    "format",
    "chmod",
    "chown",
    "sudo",
    "systemctl",
    "service",
    "iptables",
    "ufw"
]

denylist = [
    "rm\\s+-rf\\s+/",
    "dd\\s+if=",
    "mkfs",
    ":\\(\\)\\{.*\\};:",
    "chmod\\s+777",
    "curl.*\\|.*bash",
    "wget.*\\|.*sh"
]
```

### Development Configuration

For development and testing:

**config.toml:**
```toml
[ai]
provider_profile = "fast"  # Use free/cheap model

[context]
history_length = 20
include_git = true
include_env = true

[safety]
enable_filters = false  # Disable for testing

[shell]
timeout = 30  # Longer timeout for debugging
```

## Advanced Configuration

### Multiple Profiles

Create multiple provider profiles for different use cases:

**providers.toml:**
```toml
# Quick suggestions for common tasks
[profiles.quick]
provider = "openrouter"
model = "meta-llama/llama-3.1-8b-instruct:free"
temperature = 0.1
max_tokens = 200

# Detailed explanations
[profiles.explain]
provider = "openrouter"
model = "openai/gpt-4o"
temperature = 0.5
max_tokens = 1000

# Git-specific suggestions
[profiles.git]
provider = "openrouter"
model = "openai/gpt-4o-mini"
temperature = 0.3
max_tokens = 400
```

**Usage:**
```bash
# Change default profile
linesense config set-model anthropic/claude-3.5-sonnet

# Or override per-command
linesense suggest --line "git rebase" --model "openai/gpt-4o"
```

### Custom Environment Variable Filtering

Include project-specific environment variables:

**config.toml:**
```toml
[context]
include_env = true
env_allowlist = [
    # Standard
    "PATH", "HOME", "USER",

    # Development
    "NODE_ENV", "PYTHON_PATH", "GOPATH",

    # Project-specific (use patterns)
    "MYAPP_*",
    "DATABASE_*",

    # Cloud providers
    "AWS_REGION", "GCP_PROJECT",
]
```

### Risk Pattern Customization

Fine-tune risk classification for your workflow:

**config.toml:**
```toml
[safety]
enable_filters = true

# Mark these as high-risk (show warnings)
require_confirm_patterns = [
    # Database operations
    "DROP\\s+DATABASE",
    "TRUNCATE",
    "DELETE\\s+FROM.*WHERE.*",

    # File operations
    "rm.*-rf",
    "shred",

    # Network
    "nc.*-l",  # Netcat listening
    "nmap",    # Port scanning
]

# Completely block these
denylist = [
    # Destructive
    "rm\\s+-rf\\s+/",
    "dd\\s+if=.*of=/dev/sd",

    # Malicious
    ":\\(\\)\\{.*\\};:",  # Fork bomb
    "wget.*\\|.*sh",      # Remote script execution
]
```

### Shell-Specific Configuration

Configure different settings per shell:

**~/.bashrc:**
```bash
# Bash-specific settings
export LINESENSE_SUGGEST_KEY="\C-t"
export LINESENSE_EXPLAIN_KEY="\C-h"

# Use fast model in bash
export LINESENSE_PROFILE="fast"

source /path/to/linesense/scripts/linesense.bash
```

**~/.zshrc:**
```zsh
# Zsh-specific settings
export LINESENSE_SUGGEST_KEY="^O"
export LINESENSE_EXPLAIN_KEY="^P"

# Use smart model in zsh
export LINESENSE_PROFILE="smart"

source /path/to/linesense/scripts/linesense.zsh
```

## Configuration Best Practices

1. **Start with defaults** - Use `linesense config init` and adjust from there

2. **Use appropriate profiles** - Fast for simple tasks, smart for complex operations

3. **Limit history length** - Balance context quality with API costs (20-50 recommended)

4. **Enable safety filters** - Always keep `enable_filters = true` in production

5. **Secure API keys** - Use `linesense config set-key` instead of manual editing

6. **Test changes** - After modifying config, run `linesense config show` to verify

7. **Version control** - Consider committing `config.toml` but NOT shell RC files with API keys

8. **Regular updates** - Check for new models and features periodically

## Troubleshooting Configuration

### Configuration Not Loading

```bash
# Check if files exist
ls -la ~/.config/linesense/

# Verify file permissions
chmod 600 ~/.config/linesense/*.toml

# Test with explicit path
LINESENSE_CONFIG=~/.config/linesense/config.toml linesense config show
```

### Invalid TOML Syntax

```bash
# Validate TOML syntax online
# Copy your config to https://www.toml-lint.com/

# Or reinstall from examples
mv ~/.config/linesense/config.toml ~/.config/linesense/config.toml.backup
linesense config init
```

### API Key Issues

```bash
# Verify API key is set
echo $OPENROUTER_API_KEY

# Re-set API key
linesense config set-key

# Check shell RC file
grep OPENROUTER_API_KEY ~/.bashrc ~/.zshrc
```

## See Also

- [INSTALLATION.md](INSTALLATION.md) - Installation and setup
- [API.md](API.md) - CLI command reference
- [SECURITY.md](SECURITY.md) - Security features and best practices
- [README.md](../README.md) - General usage guide
