# LineSense CLI API Reference

Complete reference for all LineSense command-line interface commands.

## Table of Contents

- [Command Overview](#command-overview)
- [Global Options](#global-options)
- [Commands](#commands)
  - [suggest](#suggest)
  - [explain](#explain)
  - [config](#config)
  - [version](#version)
  - [help](#help)
- [Exit Codes](#exit-codes)
- [Output Formats](#output-formats)
- [Examples](#examples)

## Command Overview

LineSense provides a simple CLI with the following commands:

```
linesense <command> [options]

Commands:
  suggest     Generate command suggestions from natural language
  explain     Explain what a command does
  config      Manage LineSense configuration
  version     Show version information
  help        Show help message
```

## Global Options

These options are available for all commands:

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help for the command |
| `--version`, `-v` | Show version (when used as command) |

**Environment Variables:**

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENROUTER_API_KEY` | OpenRouter API key (required) | - |
| `XDG_CONFIG_HOME` | Config directory location | `~/.config` |

## Commands

### suggest

Generate command suggestions based on natural language input.

**Syntax:**
```bash
linesense suggest [options]
```

**Required Options:**

| Option | Type | Description |
|--------|------|-------------|
| `--line <text>` | string | Natural language description or partial command |

**Optional Options:**

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--shell <type>` | string | auto-detect | Shell type: `bash` or `zsh` |
| `--cwd <path>` | string | current dir | Current working directory |
| `--model <id>` | string | from config | Override model ID from config |

**Examples:**

```bash
# Basic usage
linesense suggest --line "list files"

# With specific shell
linesense suggest --line "find large files" --shell bash

# With custom directory
linesense suggest --line "git status" --cwd ~/projects/myapp

# Override model
linesense suggest --line "complex docker command" --model openai/gpt-4o
```

**Output:**

JSON object with array of suggestions:

```json
{
  "suggestions": [
    {
      "command": "ls -la",
      "risk": "low",
      "explanation": "Lists all files in long format including hidden files",
      "source": "llm"
    }
  ]
}
```

**Output Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `command` | string | The suggested shell command |
| `risk` | string | Risk level: `low`, `medium`, or `high` |
| `explanation` | string | Why this command was suggested |
| `source` | string | Source of suggestion: `llm`, `history`, or `builtin` |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - suggestions generated |
| `1` | Error - missing required flags, API error, or config error |

**Common Errors:**

```bash
# Missing --line flag
$ linesense suggest
Error: --line flag is required

# API key not set
$ linesense suggest --line "test"
Error: failed to create provider: OPENROUTER_API_KEY not set

# Invalid shell type
$ linesense suggest --line "test" --shell invalid
Error: invalid shell type: invalid
```

---

### explain

Get a detailed explanation of what a command does.

**Syntax:**
```bash
linesense explain [options]
```

**Required Options:**

| Option | Type | Description |
|--------|------|-------------|
| `--line <command>` | string | The command to explain |

**Optional Options:**

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--shell <type>` | string | auto-detect | Shell type: `bash` or `zsh` |
| `--cwd <path>` | string | current dir | Current working directory |
| `--model <id>` | string | from config | Override model ID from config |

**Examples:**

```bash
# Basic usage
linesense explain --line "ls -la"

# Explain complex command
linesense explain --line "find . -type f -name '*.go' -exec grep -l 'TODO' {} \;"

# Explain dangerous command
linesense explain --line "rm -rf /"

# With specific context
linesense explain --line "git rebase -i HEAD~3" --cwd ~/project
```

**Output:**

JSON object with explanation details:

```json
{
  "summary": "The `ls -la` command lists all files and directories in the current directory in long format, including hidden files.",
  "risk": "low",
  "notes": [
    "The -l flag shows detailed information (permissions, owner, size, date)",
    "The -a flag includes hidden files (those starting with .)",
    "This is a read-only command and cannot modify any files",
    "Commonly used to view directory contents with full details"
  ]
}
```

**Output Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `summary` | string | High-level explanation of the command |
| `risk` | string | Risk level: `low`, `medium`, or `high` |
| `notes` | array | Detailed notes about flags, behavior, and warnings |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - explanation generated |
| `1` | Error - missing required flags, API error, or config error |

**Risk Classification Examples:**

**Low Risk:**
```bash
$ linesense explain --line "cat README.md"
{
  "summary": "Read and display the contents of README.md",
  "risk": "low",
  "notes": ["Read-only operation", "Cannot modify files"]
}
```

**Medium Risk:**
```bash
$ linesense explain --line "sudo apt-get update"
{
  "summary": "Update package lists with elevated privileges",
  "risk": "medium",
  "notes": ["Requires sudo privileges", "Modifies system state", ...]
}
```

**High Risk:**
```bash
$ linesense explain --line "rm -rf /"
{
  "summary": "DANGER: Recursively delete all files starting from root",
  "risk": "high",
  "notes": ["Extremely destructive", "Will delete entire filesystem", ...]
}
```

---

### config

Manage LineSense configuration.

**Syntax:**
```bash
linesense config <subcommand> [options]
```

**Subcommands:**

| Subcommand | Description |
|------------|-------------|
| `init` | Initialize configuration with interactive setup |
| `set-key` | Set OpenRouter API key securely |
| `set-model` | Change the default model |
| `show` | Display current configuration |

#### config init

Initialize configuration files.

**Syntax:**
```bash
linesense config init
```

**Behavior:**
1. Creates `~/.config/linesense/` directory
2. Copies default `config.toml` and `providers.toml`
3. Prompts before overwriting existing files
4. Sets proper file permissions (0600)
5. Displays next steps

**Example:**

```bash
$ linesense config init
🚀 LineSense Configuration Setup
================================

Configuration directory: /home/user/.config/linesense

✓ Created /home/user/.config/linesense/config.toml
✓ Created /home/user/.config/linesense/providers.toml

📝 Next steps:
1. Set your OpenRouter API key:
   linesense config set-key YOUR_API_KEY

2. (Optional) Change the default model:
   linesense config set-model openai/gpt-4o

3. Test it:
   linesense suggest --line "list files"
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - config initialized |
| `1` | Error - cannot create directory or files |

#### config set-key

Set or update OpenRouter API key.

**Syntax:**
```bash
# Interactive (recommended - key not in shell history)
linesense config set-key

# Direct (key will appear in shell history)
linesense config set-key <api-key>
```

**Arguments:**

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `<api-key>` | string | No | OpenRouter API key (if omitted, prompts interactively) |

**Behavior:**
1. Auto-detects shell (bash or zsh)
2. Adds/updates `OPENROUTER_API_KEY` in shell RC file
3. Prompts before overwriting existing key
4. Sets file permissions to 0600
5. Displays reload instructions

**Example:**

```bash
$ linesense config set-key
Enter your OpenRouter API key: sk-or-v1-...
✓ API key saved to /home/user/.bashrc

⚠️  IMPORTANT: Reload your shell or run:
   source /home/user/.bashrc
```

**With Existing Key:**

```bash
$ linesense config set-key
Enter your OpenRouter API key: sk-or-v1-...
⚠️  OPENROUTER_API_KEY already exists in /home/user/.bashrc
Do you want to update it? (y/N): y
✓ API key saved to /home/user/.bashrc
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - API key set |
| `1` | Error - cannot write to RC file or user cancelled |

#### config set-model

Change the default model.

**Syntax:**
```bash
linesense config set-model <model-id>
```

**Arguments:**

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `<model-id>` | string | Yes | Model ID from OpenRouter |

**Behavior:**
1. Updates `model` field in default profile
2. Modifies `~/.config/linesense/providers.toml`
3. Displays popular model options

**Example:**

```bash
$ linesense config set-model openai/gpt-4o
✓ Default model updated to: openai/gpt-4o

Popular models:
  openai/gpt-4o-mini       - Fast and cheap
  openai/gpt-4o            - Most capable
  meta-llama/llama-3.1-8b  - Open source, fast
```

**Popular Models:**

| Model ID | Description | Speed | Cost |
|----------|-------------|-------|------|
| `openai/gpt-4o-mini` | Best balance (default) | Fast | Low |
| `openai/gpt-4o` | Highest quality | Medium | High |
| `meta-llama/llama-3.1-8b-instruct:free` | Free tier | Very Fast | Free |
| `anthropic/claude-3.5-sonnet` | Excellent reasoning | Medium | Medium |

See [OpenRouter Models](https://openrouter.ai/models) for complete list.

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - model updated |
| `1` | Error - cannot write to config or invalid model |

#### config show

Display current configuration.

**Syntax:**
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

**API Key Display:**
- First 8 characters shown
- Last 4 characters shown
- Middle portion masked with `...`
- Example: `sk-or-v1...cf28`

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Success - config displayed |
| `1` | Error - cannot load config files |

---

### version

Display version information.

**Syntax:**
```bash
linesense version
# or
linesense --version
# or
linesense -v
```

**Output:**
```
linesense version 0.1.0
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Always successful |

---

### help

Display help information.

**Syntax:**
```bash
linesense help
# or
linesense --help
# or
linesense -h
```

**Output:**
```
linesense - AI-powered shell command assistant

Usage:
  linesense config [subcommand]  Configure LineSense
  linesense suggest [flags]      Generate command suggestions
  linesense explain [flags]      Explain a command
  linesense version              Show version information
  linesense help                 Show this help message

Config Subcommands:
  init            Initialize configuration with interactive setup
  set-key         Set OpenRouter API key securely
  set-model       Change the default model
  show            Display current configuration

[... full help text ...]
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| `0` | Always successful |

---

## Exit Codes

All LineSense commands use these exit codes:

| Code | Meaning | When |
|------|---------|------|
| `0` | Success | Command completed successfully |
| `1` | Error | Any error occurred |

**Common Error Scenarios:**

- Missing required flags
- API key not set
- Configuration file errors
- Network errors
- API errors
- Invalid input
- Permission errors

---

## Output Formats

### JSON Output

`suggest` and `explain` commands return JSON for easy parsing.

**Parsing with jq:**
```bash
# Get first suggestion
linesense suggest --line "list files" | jq -r '.suggestions[0].command'

# Get risk level
linesense explain --line "rm -rf /" | jq -r '.risk'

# Get all suggestions
linesense suggest --line "git commit" | jq -r '.suggestions[].command'
```

**Parsing with grep/sed:**
```bash
# Get first command (without jq)
linesense suggest --line "list" | grep -o '"command":"[^"]*"' | head -1 | sed 's/"command":"\(.*\)"/\1/'
```

**Parsing in scripts:**
```bash
#!/bin/bash
result=$(linesense suggest --line "list files")
command=$(echo "$result" | jq -r '.suggestions[0].command')
risk=$(echo "$result" | jq -r '.suggestions[0].risk')

if [ "$risk" = "high" ]; then
    echo "⚠️  High risk command: $command"
    exit 1
fi

echo "Suggested: $command"
```

### Human-Readable Output

`config` commands use formatted text output for readability.

**Features:**
- Clear section headers
- Emoji indicators (✓, ❌, ⚠️, 📋, 🚀)
- Masked sensitive data
- Structured information
- Color support (in shell integration)

---

## Examples

### Basic Usage

```bash
# Get command suggestion
linesense suggest --line "show running processes"

# Explain a command
linesense explain --line "ps aux | grep nginx"

# Initialize config
linesense config init

# Set API key
linesense config set-key
```

### Advanced Usage

```bash
# Use specific model for complex task
linesense suggest --line "create docker compose for nginx and postgres" \
  --model openai/gpt-4o

# Context-aware git suggestion
cd ~/myproject
linesense suggest --line "commit these changes" --cwd ~/myproject

# Explain with specific shell syntax
linesense explain --line "source ~/.bashrc" --shell bash
```

### Scripting

**Automated suggestions:**
```bash
#!/bin/bash
# suggest.sh - Get AI suggestions

prompt="$*"
if [ -z "$prompt" ]; then
    echo "Usage: suggest.sh <description>"
    exit 1
fi

result=$(linesense suggest --line "$prompt" 2>/dev/null)
if [ $? -ne 0 ]; then
    echo "Error getting suggestion"
    exit 1
fi

command=$(echo "$result" | jq -r '.suggestions[0].command')
risk=$(echo "$result" | jq -r '.suggestions[0].risk')

case "$risk" in
    high)
        echo "⚠️  HIGH RISK: $command"
        ;;
    medium)
        echo "⚡ MEDIUM RISK: $command"
        ;;
    *)
        echo "✓ $command"
        ;;
esac
```

**Usage:**
```bash
$ ./suggest.sh find large files
✓ find . -type f -size +100M

$ ./suggest.sh delete everything
⚠️  HIGH RISK: rm -rf /
```

**Safe execution wrapper:**
```bash
#!/bin/bash
# safe-suggest.sh - Only execute low-risk commands

prompt="$*"
result=$(linesense suggest --line "$prompt")
command=$(echo "$result" | jq -r '.suggestions[0].command')
risk=$(echo "$result" | jq -r '.suggestions[0].risk')

if [ "$risk" != "low" ]; then
    echo "⚠️  Command is $risk risk. Manual execution required:"
    echo "  $command"
    exit 1
fi

echo "Executing: $command"
eval "$command"
```

### Pipeline Usage

```bash
# Suggest and execute (dangerous - review first!)
linesense suggest --line "list files" | jq -r '.suggestions[0].command' | bash

# Explain multiple commands
cat commands.txt | while read cmd; do
    echo "=== $cmd ==="
    linesense explain --line "$cmd" | jq -r '.summary'
    echo
done

# Batch suggestions
echo "list files
find large files
show processes" | while read prompt; do
    linesense suggest --line "$prompt" | jq -r '.suggestions[0].command'
done
```

### Error Handling

```bash
#!/bin/bash
# robust-suggest.sh - Proper error handling

if [ -z "$OPENROUTER_API_KEY" ]; then
    echo "Error: OPENROUTER_API_KEY not set. Run:"
    echo "  linesense config set-key"
    exit 1
fi

result=$(linesense suggest --line "$1" 2>&1)
exit_code=$?

if [ $exit_code -ne 0 ]; then
    echo "Error: $result" >&2
    exit $exit_code
fi

if ! echo "$result" | jq empty 2>/dev/null; then
    echo "Error: Invalid JSON response" >&2
    exit 1
fi

command=$(echo "$result" | jq -r '.suggestions[0].command')
if [ "$command" = "null" ] || [ -z "$command" ]; then
    echo "Error: No suggestion generated" >&2
    exit 1
fi

echo "$command"
```

## See Also

- [README.md](../README.md) - General usage guide
- [INSTALLATION.md](INSTALLATION.md) - Installation instructions
- [CONFIGURATION.md](CONFIGURATION.md) - Configuration reference
- [SECURITY.md](SECURITY.md) - Security features and best practices
