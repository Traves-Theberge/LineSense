# LineSense Installation Guide

This guide will walk you through installing and configuring LineSense on your system.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [From Source](#from-source)
  - [Using go install](#using-go-install)
- [Configuration](#configuration)
  - [Quick Setup](#quick-setup)
  - [Manual Setup](#manual-setup)
- [Shell Integration](#shell-integration)
  - [Bash Setup](#bash-setup)
  - [Zsh Setup](#zsh-setup)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before installing LineSense, ensure you have:

1. **Go 1.21 or later**
   ```bash
   go version
   # Should output: go version go1.21.x or higher
   ```

2. **An OpenRouter API key**
   - Sign up at [https://openrouter.ai](https://openrouter.ai)
   - Get your API key from the dashboard
   - Keep it safe - you'll need it during configuration

3. **Git** (for installation from source)
   ```bash
   git --version
   ```

4. **jq** (optional but recommended for better shell integration)
   ```bash
   # Ubuntu/Debian
   sudo apt-get install jq

   # macOS
   brew install jq

   # Fedora/RHEL
   sudo dnf install jq

   # Arch Linux
   sudo pacman -S jq
   ```

## Installation Methods

### From Source

This method gives you the latest development version and is recommended for development:

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/linesense.git
cd linesense

# 2. Build the binary
go build -o linesense ./cmd/linesense

# 3. Move to a directory in your PATH
sudo mv linesense /usr/local/bin/

# Verify installation
linesense version
```

### Using go install

This is the quickest method for users:

```bash
# Install directly to $GOPATH/bin
go install github.com/yourusername/linesense/cmd/linesense@latest

# Make sure $GOPATH/bin is in your PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify installation
linesense version
```

## Configuration

### Quick Setup

The easiest way to set up LineSense is using the interactive configuration wizard:

```bash
# 1. Initialize configuration
linesense config init

# This will:
# - Create ~/.config/linesense/ directory
# - Copy default config files
# - Display next steps

# 2. Set your API key (interactive)
linesense config set-key
# Enter your OpenRouter API key when prompted

# 3. Reload your shell
source ~/.bashrc  # for bash
# or
source ~/.zshrc   # for zsh

# 4. Verify configuration
linesense config show
```

### Manual Setup

If you prefer manual configuration:

```bash
# 1. Create config directory
mkdir -p ~/.config/linesense

# 2. Copy example configurations
# (Assuming you cloned the repository)
cp examples/config.toml ~/.config/linesense/
cp examples/providers.toml ~/.config/linesense/

# 3. Set API key manually
# Add this to your ~/.bashrc or ~/.zshrc
echo 'export OPENROUTER_API_KEY="your-api-key-here"' >> ~/.bashrc

# 4. Reload shell
source ~/.bashrc
```

## Shell Integration

Shell integration enables interactive keybindings for AI-powered suggestions.

### Bash Setup

Add this to your `~/.bashrc`:

```bash
# LineSense shell integration
if [ -f /path/to/linesense/scripts/linesense.bash ]; then
    source /path/to/linesense/scripts/linesense.bash
fi
```

**Custom keybindings** (optional):

```bash
# Set before sourcing the script
export LINESENSE_SUGGEST_KEY="\C-t"      # Ctrl+T for suggestions
export LINESENSE_EXPLAIN_KEY="\C-x\C-h"  # Ctrl+X Ctrl+H for explanations

source /path/to/linesense/scripts/linesense.bash
```

**Default keybindings:**
- `Ctrl+Space` - Get AI suggestions for current line
- `Ctrl+X Ctrl+E` - Explain current command

### Zsh Setup

Add this to your `~/.zshrc`:

```zsh
# LineSense shell integration
if [ -f /path/to/linesense/scripts/linesense.zsh ]; then
    source /path/to/linesense/scripts/linesense.zsh
fi
```

**Custom keybindings** (optional):

```zsh
# Set before sourcing the script
export LINESENSE_SUGGEST_KEY="^T"      # Ctrl+T for suggestions
export LINESENSE_EXPLAIN_KEY="^X^H"   # Ctrl+X Ctrl+H for explanations

source /path/to/linesense/scripts/linesense.zsh
```

**Default keybindings:**
- `Ctrl+Space` - Get AI suggestions for current line
- `Ctrl+X Ctrl+E` - Explain current command

### Reload Shell

After adding shell integration, reload your configuration:

```bash
# Bash
source ~/.bashrc

# Zsh
source ~/.zshrc
```

## Verification

Test your installation to ensure everything works:

### 1. Check Version

```bash
linesense version
# Should output: linesense version 0.4.0
```

### 2. View Configuration

```bash
linesense config show
# Should display your configuration with masked API key
```

### 3. Test Suggest Command

```bash
linesense suggest --line "list files"
# Should return JSON with command suggestions
```

### 4. Test Explain Command

```bash
linesense explain --line "ls -la"
# Should return JSON with command explanation
```

### 5. Test Shell Integration

Open a new terminal or reload your shell, then:

1. Type a partial command: `list fil`
2. Press `Ctrl+Space`
3. You should see AI-suggested completion

## Troubleshooting

### API Key Not Found

**Problem:** Error message "OPENROUTER_API_KEY environment variable not set"

**Solutions:**
```bash
# 1. Check if API key is set
echo $OPENROUTER_API_KEY

# 2. If empty, set it using config command
linesense config set-key

# 3. Reload shell
source ~/.bashrc  # or ~/.zshrc

# 4. Verify it's set
linesense config show
```

### Configuration Files Not Found

**Problem:** Error loading config files

**Solutions:**
```bash
# 1. Run config init
linesense config init

# 2. Check if files exist
ls -la ~/.config/linesense/

# 3. If missing, copy from examples
cp examples/*.toml ~/.config/linesense/
```

### Shell Integration Not Working

**Problem:** Keybindings don't trigger suggestions

**Solutions:**

For Bash:
```bash
# 1. Check if script is sourced
type _linesense_request
# Should show function definition

# 2. Check if linesense is in PATH
which linesense

# 3. Re-source the script
source /path/to/linesense/scripts/linesense.bash

# 4. Test bindings manually
bind -P | grep linesense
```

For Zsh:
```bash
# 1. Check if widget is loaded
zle -l | grep linesense

# 2. Check if linesense is in PATH
which linesense

# 3. Re-source the script
source /path/to/linesense/scripts/linesense.zsh

# 4. Test bindings manually
bindkey | grep linesense
```

### jq Not Installed Warning

**Problem:** Warning message about jq not being installed

**Impact:** Shell integration will work but use slower grep/sed fallback

**Solution:**
```bash
# Install jq for better performance
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq

# Fedora/RHEL
sudo dnf install jq
```

### Permission Denied Errors

**Problem:** Cannot create config files or write to directories

**Solutions:**
```bash
# 1. Check directory permissions
ls -la ~/.config/

# 2. Create config directory with proper permissions
mkdir -p ~/.config/linesense
chmod 700 ~/.config/linesense

# 3. For shell RC files
chmod 600 ~/.bashrc ~/.zshrc
```

### API Rate Limits

**Problem:** Too many API requests or slow responses

**Solutions:**
1. Use a faster model in your config:
   ```bash
   linesense config set-model meta-llama/llama-3.1-8b-instruct:free
   ```

2. Reduce context size in `~/.config/linesense/config.toml`:
   ```toml
   [context]
   history_length = 20  # Reduce from default 50
   ```

### Network Errors

**Problem:** Cannot connect to OpenRouter API

**Solutions:**
```bash
# 1. Check internet connectivity
curl https://openrouter.ai

# 2. Check if API endpoint is reachable
curl -I https://openrouter.ai/api/v1/chat/completions

# 3. Test with direct API call
linesense suggest --line "test" --model openai/gpt-4o-mini
```

## Next Steps

Once installed and configured, check out:

- [CONFIGURATION.md](CONFIGURATION.md) - Detailed configuration options
- [API.md](API.md) - CLI command reference
- [SECURITY.md](SECURITY.md) - Security features and best practices
- [README.md](../README.md) - General usage and examples

## Getting Help

If you encounter issues not covered here:

1. Check the [GitHub Issues](https://github.com/yourusername/linesense/issues)
2. Read the [FAQ](https://github.com/yourusername/linesense/wiki/FAQ)
3. Open a new issue with:
   - Your OS and shell version
   - Output of `linesense config show`
   - Error messages or unexpected behavior
   - Steps to reproduce the issue
