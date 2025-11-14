# LineSense Security Guide

This document describes LineSense's security features, best practices, and threat model.

## Table of Contents

- [Security Philosophy](#security-philosophy)
- [Risk Classification System](#risk-classification-system)
- [Built-in Protections](#built-in-protections)
- [API Key Security](#api-key-security)
- [Data Privacy](#data-privacy)
- [Safety Configuration](#safety-configuration)
- [Security Best Practices](#security-best-practices)
- [Threat Model](#threat-model)
- [Reporting Security Issues](#reporting-security-issues)

## Security Philosophy

LineSense is designed with security as a top priority:

1. **Defense in Depth** - Multiple layers of protection against dangerous commands
2. **Least Privilege** - Minimal permissions and data sharing by default
3. **Transparency** - Clear risk indicators and explanations for all suggestions
4. **User Control** - Extensive configuration options for security policies
5. **Secure by Default** - Safe default settings that can be relaxed if needed

**Core Principle:** LineSense assists users but never executes commands automatically. The user always has final control.

## Risk Classification System

Every command suggested by LineSense is classified into one of three risk levels:

### 🟢 Low Risk

**Definition:** Read-only commands that cannot modify system state or data.

**Examples:**
- `ls -la` - List files
- `cat file.txt` - Read file contents
- `grep pattern file` - Search file
- `pwd` - Print working directory
- `echo "hello"` - Print text
- `git status` - View git status
- `ps aux` - List processes

**Behavior:**
- No warnings shown
- Safe for all users
- Can be executed without concern

### 🟡 Medium Risk

**Definition:** Commands that modify system state but are commonly used and reversible.

**Examples:**
- `rm file.txt` - Delete a file
- `mv old.txt new.txt` - Move/rename file
- `chmod 644 file.txt` - Change permissions
- `sudo apt-get update` - Update package lists
- `git commit -m "msg"` - Create git commit
- `kill 1234` - Terminate process
- `systemctl restart nginx` - Restart service

**Behavior:**
- Yellow indicator shown in shell
- Execution allowed but user should review
- Generally safe for experienced users

### 🔴 High Risk

**Definition:** Dangerous commands that could cause data loss, system damage, or security compromise.

**Examples:**
- `rm -rf /` - Delete root filesystem
- `dd if=/dev/zero of=/dev/sda` - Wipe disk
- `mkfs.ext4 /dev/sda1` - Format partition
- `chmod 777 sensitive_file` - Overly permissive permissions
- `curl http://site.com | bash` - Execute remote script
- `:(){ :|:& };:` - Fork bomb
- `iptables -F` - Flush firewall rules

**Behavior:**
- ⚠️ Warning shown prominently
- Red indicator in shell output
- User strongly advised to review
- Can be blocked entirely via configuration

## Built-in Protections

LineSense includes comprehensive built-in protection patterns that are **always active**, regardless of configuration.

### High-Risk Patterns

These patterns trigger high-risk warnings:

```regex
rm\s+-rf\s+/              # Root filesystem deletion
dd\s+if=                  # Direct disk operations
mkfs                      # Filesystem formatting
>\s*/dev/                 # Writing to device files
chmod\s+777               # Overly permissive (777)
chmod\s+-R\s+777          # Recursive 777
curl.*\|\s*bash           # Curl to bash
wget.*\|\s*sh             # Wget to shell
:\(\)\{.*\};:             # Fork bomb
killall\s+-9              # Force kill all processes
```

### Medium-Risk Patterns

These patterns trigger medium-risk indicators:

```regex
sudo                      # Elevated privileges
rm\s+                     # File removal
mv\s+                     # File move/rename
chmod                     # Permission changes
chown                     # Ownership changes
kill                      # Process termination
pkill                     # Process killing by name
systemctl                 # System service management
reboot                    # System reboot
shutdown                  # System shutdown
iptables                  # Firewall changes
apt-get\s+remove          # Package removal (Debian/Ubuntu)
yum\s+remove              # Package removal (RedHat/Fedora)
```

### Pattern Matching

- All pattern matching is **case-insensitive**
- Patterns use **standard regex syntax**
- Patterns check the **full command line**, not just the command name
- User-defined patterns are checked **in addition to** built-in patterns

## API Key Security

### Storage

LineSense uses a secure multi-layered approach to API key storage:

**Best Practice (Recommended):**
```bash
# Use the config command - stores in shell RC file
linesense config set-key
```

**Storage locations (in order of preference):**

1. **Shell RC file** (`~/.bashrc` or `~/.zshrc`)
   - ✅ Persistent across sessions
   - ✅ Not in version control by default
   - ✅ Proper file permissions (0600)
   - ✅ Automatically loaded by shell

2. **Environment variable** (session-only)
   ```bash
   export OPENROUTER_API_KEY="sk-or-v1-..."
   ```
   - ⚠️ Lost when terminal closes
   - ✅ Not in shell history if exported directly

3. **.env file** (development only)
   ```bash
   echo 'OPENROUTER_API_KEY="sk-or-v1-..."' > .env
   ```
   - ⚠️ Must be in `.gitignore`
   - ⚠️ Only loaded in project directory
   - ✅ Good for development/testing

**Never:**
- ❌ Hardcode in source code
- ❌ Commit to version control
- ❌ Store in config files (config.toml)
- ❌ Share in public forums/issues
- ❌ Include in screenshots or logs

### File Permissions

LineSense automatically sets secure file permissions:

```bash
# Shell RC files
~/.bashrc:  0600 (rw-------)  # Owner read/write only
~/.zshrc:   0600 (rw-------)

# Config files
~/.config/linesense/:        0700 (rwx------)  # Owner only
~/.config/linesense/*.toml:  0600 (rw-------)
```

**What this means:**
- Only your user account can read/write these files
- Other users on the system cannot access your API key
- Web servers and services cannot read the files

### API Key Masking

When displaying configuration, API keys are automatically masked:

```bash
$ linesense config show
API Key: sk-or-v1...cf28 ✓
```

**Masking format:**
- Shows first 8 characters
- Shows last 4 characters
- Hides middle portion with `...`

### Key Rotation

To rotate your API key:

```bash
# 1. Generate new key at https://openrouter.ai
# 2. Update LineSense
linesense config set-key NEW_KEY_HERE
# 3. Reload shell
source ~/.bashrc
# 4. Verify
linesense config show
# 5. Revoke old key at OpenRouter dashboard
```

## Data Privacy

### What Data is Sent to OpenRouter

LineSense sends the following data to OpenRouter's API:

**Always Sent:**
1. **Partial command input** - The text you've typed
2. **System prompt** - Instructions for the AI model

**Optionally Sent (configurable):**
3. **Shell history** - Recent commands (default: last 50)
4. **Git information** - Current branch, status, remotes
5. **Environment variables** - Filtered through allowlist
6. **Current directory** - Working directory path

**Never Sent:**
- File contents (unless explicitly in command)
- SSH keys or credentials
- Browser history or cookies
- System logs or private data
- Other users' data

### Controlling Data Sharing

**Minimal data sharing:**
```toml
[context]
history_length = 0        # No history
include_git = false       # No git info
include_env = false       # No env vars
```

**Maximum privacy:**
```bash
# Use only the CLI without context
linesense suggest --line "list files" --cwd /tmp
```

### Data Retention

**By OpenRouter:**
- See [OpenRouter Privacy Policy](https://openrouter.ai/privacy)
- Requests may be logged for model providers
- Consider using models with strict data policies

**By LineSense:**
- No data stored locally except configuration
- No telemetry or analytics
- No usage tracking (yet - see Phase 3 roadmap)

### Environment Variable Filtering

Only variables in the allowlist are sent:

**Default allowlist:**
```toml
env_allowlist = [
    "PATH",      # System paths
    "USER",      # Username
    "HOME",      # Home directory
    "SHELL",     # Shell type
    "LANG",      # Language
    "LC_ALL",    # Locale
    "EDITOR",    # Default editor
    "VISUAL"     # Visual editor
]
```

**Never included (even if in allowlist):**
- Variables containing "PASSWORD", "SECRET", "KEY", "TOKEN"
- Variables with credential-like patterns
- SSH-related private keys

## Safety Configuration

### Configuring Risk Patterns

Add custom high-risk patterns to your `config.toml`:

```toml
[safety]
enable_filters = true

# These will show ⚠️ warnings
require_confirm_patterns = [
    # Database operations
    "DROP\\s+DATABASE",
    "TRUNCATE\\s+TABLE",
    "DELETE\\s+FROM.*WHERE",

    # Disk operations
    "format",
    "fdisk",
    "parted",

    # Security changes
    "setenforce\\s+0",      # Disable SELinux
    "ufw\\s+disable",       # Disable firewall
    "iptables.*-F",         # Flush firewall
]
```

### Blocking Commands

Completely prevent certain commands from being suggested:

```toml
[safety]
# These commands will NEVER be suggested
denylist = [
    "rm\\s+-rf\\s+/",           # Delete root
    "dd\\s+if=.*of=/dev/sd",    # Disk writes
    "mkfs",                     # Format disk
    ":\\(\\)\\{.*\\};:",        # Fork bomb
    "shutdown",                 # No shutdowns
    "reboot",                   # No reboots
]
```

### Disabling Safety (Not Recommended)

For testing or development only:

```toml
[safety]
enable_filters = false
```

⚠️ **Warning:** Only disable safety filters in controlled environments. Never disable in production or on systems with important data.

## Security Best Practices

### For Users

1. **Always review commands before execution**
   - Read what LineSense suggests
   - Understand what it will do
   - Don't blindly execute, especially for high-risk commands

2. **Use appropriate models**
   - More capable models = better safety analysis
   - Consider using `openai/gpt-4o` for critical operations

3. **Keep safety filters enabled**
   - Don't disable `enable_filters` unless necessary
   - Add custom patterns for your environment

4. **Secure your API key**
   - Use `linesense config set-key`
   - Don't share or commit keys
   - Rotate keys periodically

5. **Limit context sharing**
   - Only include necessary environment variables
   - Consider reducing history length
   - Disable git info if working with sensitive repos

6. **Monitor usage**
   - Check OpenRouter dashboard for unusual activity
   - Review API key usage regularly

### For Administrators

1. **Deploy with safe defaults**
   ```bash
   # Create organization config
   mkdir -p /etc/linesense/
   cp config.toml /etc/linesense/config.toml

   # Set restrictive defaults
   # Edit /etc/linesense/config.toml
   ```

2. **Use custom denylists**
   - Block commands specific to your environment
   - Add compliance-required patterns

3. **Audit logging** (future feature)
   - Enable usage logging when available
   - Monitor for policy violations

4. **Training**
   - Educate users about risk indicators
   - Provide clear guidelines for high-risk commands

### For Developers

1. **Never commit API keys**
   ```bash
   # Add to .gitignore
   echo ".env" >> .gitignore
   echo "*.toml" >> .gitignore  # If contains secrets
   ```

2. **Use .env for development**
   ```bash
   echo 'OPENROUTER_API_KEY="sk-or-..."' > .env
   # Already in .gitignore
   ```

3. **Validate inputs**
   - All user input is validated
   - Commands are sanitized before AI processing

4. **Review dependencies**
   - Keep dependencies updated
   - Monitor for security advisories

## Threat Model

### In Scope

**Threats LineSense protects against:**

1. **Accidental Execution**
   - User accidentally running destructive commands
   - Typos leading to dangerous operations
   - Misunderstanding command behavior

2. **AI Hallucinations**
   - AI suggesting incorrect or dangerous commands
   - Model generating unsafe command combinations
   - Context misinterpretation leading to bad suggestions

3. **Credential Exposure**
   - API keys leaked in logs or screenshots
   - Credentials sent to AI models
   - Keys committed to version control

4. **Data Leakage**
   - Sensitive environment variables sent to API
   - Private file paths revealed
   - User data in context

### Out of Scope

**Threats LineSense does NOT protect against:**

1. **Intentional Misuse**
   - User deliberately disabling safety features
   - Intentional execution of known-dangerous commands
   - Social engineering attacks

2. **Compromised System**
   - Malware on user's machine
   - Rootkits or kernel-level attacks
   - Compromised OpenRouter account

3. **Physical Access**
   - Attacker with physical access to machine
   - Stolen credentials from shoulder surfing
   - Unauthorized access to user account

4. **Supply Chain Attacks**
   - Compromised dependencies
   - Malicious OpenRouter API responses
   - Model poisoning

### Assumptions

LineSense security model assumes:

1. **User has control** - User can read and understand commands
2. **Shell is trusted** - The shell environment is not compromised
3. **OpenRouter is trusted** - API provider is secure and reliable
4. **Network is secure** - TLS/HTTPS protects API communications

## Reporting Security Issues

### Security Vulnerabilities

If you discover a security vulnerability in LineSense:

**DO:**
- Email security details to: security@yourproject.com (replace with actual email)
- Provide detailed reproduction steps
- Allow 90 days for fix before public disclosure
- Use PGP encryption for sensitive details

**DON'T:**
- Post vulnerabilities in public issues
- Share exploits publicly before fix
- Test vulnerabilities on production systems

### Reporting Format

```
Subject: [SECURITY] Brief description

Description:
[Detailed description of the vulnerability]

Impact:
[What an attacker could achieve]

Reproduction:
1. Step 1
2. Step 2
3. ...

Environment:
- LineSense version:
- OS:
- Shell:

Suggested Fix:
[If you have suggestions]
```

### Security Advisories

Security advisories will be published:
- GitHub Security Advisories
- Project README
- Release notes

## Security Checklist

Use this checklist to ensure secure LineSense usage:

- [ ] API key set using `linesense config set-key` (not manual editing)
- [ ] Shell RC file permissions are 0600
- [ ] Config directory permissions are 0700
- [ ] Safety filters enabled (`enable_filters = true`)
- [ ] Custom denylist configured for your environment
- [ ] Environment variable allowlist reviewed and minimal
- [ ] Using latest LineSense version
- [ ] OpenRouter API key rotated periodically
- [ ] .env file in .gitignore (if using)
- [ ] No API keys in version control
- [ ] Users trained on risk indicators
- [ ] Monitoring OpenRouter usage dashboard

## Additional Resources

- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [OpenRouter Security](https://openrouter.ai/security)
- [Shell Security Best Practices](https://www.gnu.org/software/bash/manual/html_node/Shell-Security.html)

## See Also

- [CONFIGURATION.md](CONFIGURATION.md) - Configuration options including safety settings
- [API.md](API.md) - CLI command reference
- [INSTALLATION.md](INSTALLATION.md) - Secure installation guide
- [README.md](../README.md) - General usage guide
