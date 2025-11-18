Got you. Here’s a **clean, self-contained spec** with:

1. **Key Entities** (for data models / modules)
2. **Full Gherkin feature set** (for behavior & tests)

You can paste this directly into a system prompt / task file for an AI dev agent.

---

## 1. Key Entities

### 1.1 Global Config (`Config`)

Represents: `~/.config/linesense/config.toml`

```text
Config
  shell:
    enable_bash: bool
    enable_zsh: bool
  keybindings:
    suggest: string        # e.g. "ctrl+space"
    explain: string        # e.g. "ctrl+e"
    alternatives: string   # e.g. "alt+a"
  context:
    history_length: int    # how many recent commands to use
    include_git: bool
    include_files: bool
    include_env: bool
  safety:
    require_confirm_patterns: []string
    denylist: []string
    default_execution: string  # "paste_only" (v0.1)
  ai:
    provider_profile: string  # "default" | "fast" | "smart" | etc.
```

---

### 1.2 Provider Config (`ProvidersConfig`)

Represents: `~/.config/linesense/providers.toml`
Primary provider: **OpenRouter**.

```text
ProvidersConfig
  default: ProfileConfig
  profiles: map[string]ProfileConfig
  openrouter: OpenRouterConfig

ProfileConfig
  provider: string          # "openrouter"
  model: string             # e.g. "openrouter/openai/gpt-4.1-mini"
  temperature: float
  max_tokens: int

OpenRouterConfig
  type: string              # "openrouter"
  api_key_env: string       # e.g. "OPENROUTER_API_KEY"
  base_url: string          # e.g. "https://openrouter.ai/api/v1"
  timeout_ms: int
```

---

### 1.3 Usage Log (`UsageEvent`)

Optional local log for adaptive behavior (stored e.g. at `~/.config/linesense/usage.log`).

```text
UsageEvent
  timestamp: string         # ISO 8601
  cwd: string
  command: string
  accepted: bool            # whether the user executed it as suggested
  source: string            # "llm"
```

---

### 1.4 Runtime Context (`ContextEnvelope`)

Collected before each suggestion / explanation.

```text
ContextEnvelope
  shell: string             # "bash" | "zsh"
  line: string              # current input line
  cwd: string
  git: GitInfo | null
  env: map[string]string    # filtered env (if enabled)
  history: []HistoryEntry   # last N commands
  usage_summary: UsageSummary | null

GitInfo
  is_repo: bool
  branch: string
  status_summary: string
  remotes: []string

HistoryEntry
  command: string
  timestamp: string | null
  exit_code: int | null

UsageSummary
  frequently_used_commands: []string    # top N commands in this cwd
```

---

### 1.5 Suggestions & Explanations

```text
SuggestionInput
  model_id: string
  prompt: string
  context: ContextEnvelope

Suggestion
  command: string
  risk: string              # "low" | "medium" | "high"
  explanation: string
  source: string            # "llm"

ExplanationInput
  model_id: string
  prompt: string
  context: ContextEnvelope

Explanation
  summary: string
  risk: string
  notes: []string
```

---

### 1.6 Provider Interface (`Provider`)

Backed by **OpenRouter** (via Charmbracelet `fantasy` or custom client).

```text
Provider
  Name() string
  Suggest(ctx, SuggestInput) -> ([]Suggestion, error)
  Explain(ctx, ExplainInput) -> (Explanation, error)
```

Concrete implementation: `OpenRouterProvider`.

---

### 1.7 Repositories / Modules

You may implement as multiple repos or packages in a monorepo; responsibilities must remain separated:

1. `linesense`

   * CLI binary (`cmd/linesense`)
   * Shell integration files (`linesense.bash`, `linesense.zsh`)

2. `linesense-config`

   * Loading & validating global `Config` and `ProvidersConfig`

3. `linesense-core`

   * Context collection (filesystem, git, history, env)
   * Usage log reading/writing + usage summary
   * Safety filters
   * Prompt building
   * Top-level engine: `Suggest` and `Explain` functions

4. `linesense-ai`

   * Provider implementations
   * OpenRouter integration (via fantasy or HTTP client)
   * Single `Provider` instance selected based on config/CLI flags

---

## 2. Gherkin Features & Scenarios

### Feature: Global configuration loading

```gherkin
Feature: Load global configuration
  As the LineSense engine
  I want to load configuration from standard locations
  So that user preferences control behavior

  Scenario: Load global config from XDG path
    Given the environment variable "XDG_CONFIG_HOME" is set to "/home/user/.config"
    And the file "/home/user/.config/linesense/config.toml" exists
    When the engine loads the global config
    Then it should parse "config.toml" into a Config struct
    And it should not produce an error

  Scenario: Fallback when XDG_CONFIG_HOME is not set
    Given the environment variable "XDG_CONFIG_HOME" is not set
    And the file "/home/user/.config/linesense/config.toml" exists
    When the engine loads the global config
    Then it should use "/home/user/.config/linesense/config.toml"
    And it should parse the config successfully

  Scenario: Load provider configuration
    Given the file "~/.config/linesense/providers.toml" exists
    And the [default] table contains provider "openrouter" and a "model" value
    When the engine loads the providers config
    Then it should create a ProvidersConfig struct
    And ProvidersConfig.default.model should equal the configured "model"
```

---

### Feature: Usage logging and preference awareness

```gherkin
Feature: Usage logging and preference awareness
  As the LineSense engine
  I want to record accepted suggestions locally
  So that future suggestions reflect the user's habits

  Scenario: Record accepted suggestion
    Given the user accepts a suggested command and executes it
    And usage logging is enabled by default
    When the engine receives confirmation that the command was executed
    Then it should append a UsageEvent entry to "usage.log"
    And the entry should contain cwd, command, timestamp, accepted=true and source

  Scenario: Build usage summary for a directory
    Given "usage.log" contains prior events for cwd "/home/user/app"
    And some of these commands were accepted multiple times
    When the ContextEnvelope is built for "/home/user/app"
    Then ContextEnvelope.usage_summary.frequently_used_commands should include the most frequently accepted commands for that cwd

  Scenario: Do not fail when usage log is missing
    Given "usage.log" does not exist
    When the ContextEnvelope is built
    Then ContextEnvelope.usage_summary may be null
    And the engine should not error
```

---

### Feature: Shell integration for bash

```gherkin
Feature: Bash shell integration
  As a bash user
  I want LineSense bound to a key
  So that I can trigger suggestions without changing my workflow

  Scenario: Source bash integration script
    Given the "linesense" binary is installed on PATH
    And the user has added a source line for "linesense.bash" in ".bashrc"
    When a new bash session starts
    Then the function "_linesense_request" should be defined
    And a readline keybinding for the suggest action should be registered according to Config.keybindings.suggest

  Scenario: Trigger suggestion from bash
    Given a bash session with the integration script loaded
    And the current command line is "deploy backend to staging"
    When the user presses the configured suggest keybinding
    Then bash should execute "linesense suggest" with the current line and cwd
    And if the JSON response contains at least one suggestion
    Then READLINE_LINE should be replaced with the first suggestion.command
    And READLINE_POINT should be set to the end of the new line

  Scenario: Suggestion failure is graceful in bash
    Given "linesense suggest" terminates with a non-zero exit code
    When the user presses the suggest keybinding
    Then the current bash prompt should remain usable
    And the current line should not be corrupted with partial JSON or errors
```

---

### Feature: Shell integration for zsh

```gherkin
Feature: Zsh shell integration
  As a zsh user
  I want LineSense wired to a ZLE widget
  So that I can trigger suggestions with a keybinding

  Scenario: Source zsh integration script
    Given the "linesense" binary is installed on PATH
    And the user has added a source line for "linesense.zsh" in ".zshrc"
    When a new zsh session starts
    Then a ZLE widget "linesense-widget" should be defined
    And a keybinding should be registered for that widget according to Config.keybindings.suggest

  Scenario: Trigger suggestion from zsh
    Given a zsh session with the integration script loaded
    And the current BUFFER is "docker run postgres"
    When the user presses the suggest keybinding
    Then zsh should execute "linesense suggest" with the current BUFFER and cwd
    And if the JSON response contains at least one suggestion
    Then BUFFER should be replaced with the first suggestion.command
    And CURSOR should be set to the end of BUFFER
```

---

### Feature: Suggest API behavior

```gherkin
Feature: Suggest API behavior
  As the linesense CLI
  I want to generate suggestions based on the current line and context
  So that users can complete shell commands quickly

  Scenario: Suggest with context
    Given the user runs
      """
      linesense suggest --shell bash --line "deploy backend to staging" --cwd "/home/user/app"
      """
    And "/home/user/app" is a git repository
    And history_length in Config.context is greater than 0
    When the Suggest function is invoked
    Then it should build a ContextEnvelope including git info and recent history
    And it should call the Provider.Suggest method with a SuggestInput containing the resolved model_id and prompt
    And it should receive one or more Suggestion structs

  Scenario: Suggest returns JSON
    Given Suggestion structs were produced
    When the CLI prints the result
    Then the output should be valid JSON
    And the top-level value should be an object containing "suggestions"
    And "suggestions" should be an array of objects with "command", "risk", "explanation" and "source" fields
```

---

### Feature: OpenRouter provider and model selection

```gherkin
Feature: OpenRouter provider and model selection
  As the AI engine
  I want to use OpenRouter with a configurable model
  So that users can choose the balance of cost and quality

  Scenario: Create OpenRouter provider from config
    Given providers.toml contains a [openrouter] section
    And that section defines "api_key_env" and "base_url"
    And the environment variable named by "api_key_env" is set to a valid API key
    When the engine creates a Provider instance
    Then it should create an OpenRouterProvider with the configured base URL and API key
    And Provider.Name() should return "openrouter"

  Scenario: Select model from default profile
    Given providers.toml defines a [default] profile with a "model" value
    And config.toml does not override "ai.provider_profile"
    When the engine prepares a SuggestInput
    Then SuggestInput.model_id should equal the "default.model" value

  Scenario: Override model using provider_profile
    Given providers.toml defines a profile "[profile.smart]" with a "model" value
    And config.toml contains "ai.provider_profile = 'smart'"
    When the engine prepares a SuggestInput
    Then SuggestInput.model_id should equal the "profile.smart.model" value

  Scenario: Override model using CLI flag
    Given the user executes
      """
      linesense suggest --model "openrouter/meta-llama/llama-3.1-8b-instruct"
      """
    When the engine prepares a SuggestInput
    Then SuggestInput.model_id should equal "openrouter/meta-llama/llama-3.1-8b-instruct"
    And this should override any profile or default configuration
```

---

### Feature: Safety and risk classification

```gherkin
Feature: Safety and risk classification
  As the AI engine
  I want to apply safety rules to suggested commands
  So that destructive commands are flagged or blocked

  Scenario: Denylist removes dangerous suggestions
    Given Config.safety.denylist contains "rm -rf /"
    And a Suggestion.command equals "sudo rm -rf /"
    When safety filters are applied
    Then this Suggestion should be removed from the suggestions list

  Scenario: Require-confirm marks suggestions as high risk
    Given Config.safety.require_confirm_patterns contains "rm "
    And a Suggestion.command equals "rm -rf ./build"
    When safety filters are applied
    Then Suggestion.risk should be set to "high"

  Scenario: Benign command defaults to low risk
    Given a Suggestion.command equals "npm test"
    And no denylist or require-confirm patterns match
    When safety filters are applied
    Then Suggestion.risk should be "low"
```

---

### Feature: Context gathering (git, history, env, usage)

```gherkin
Feature: Context gathering
  As the AI engine
  I want to gather contextual information
  So that suggestions are tailored to the user's environment and habits

  Scenario: Collect git context inside repository
    Given the current working directory is inside a git repository
    When the ContextEnvelope is built
    Then ContextEnvelope.git.is_repo should be true
    And ContextEnvelope.git.branch should be the current branch name
    And ContextEnvelope.git.status_summary should describe the working tree state

  Scenario: No git context outside repository
    Given the current working directory is not inside a git repository
    When the ContextEnvelope is built
    Then ContextEnvelope.git.is_repo should be false

  Scenario: Include recent history up to configured length
    Given Config.context.history_length is 100
    And the underlying shell history file contains at least 150 entries
    When the ContextEnvelope is built
    Then ContextEnvelope.history should contain the last 100 entries

  Scenario: Include environment variables when enabled
    Given Config.context.include_env is true
    When the ContextEnvelope is built
    Then ContextEnvelope.env should contain a filtered subset of environment variables

  Scenario: Exclude environment variables when disabled
    Given Config.context.include_env is false
    When the ContextEnvelope is built
    Then ContextEnvelope.env should be empty or null
```

---

### Feature: Explain API (optional v0.2+)

```gherkin
Feature: Explain API
  As a user
  I want to understand commands before running them
  So that I can avoid unintended side effects

  Scenario: Explain a docker prune command
    Given the user invokes
      """
      linesense explain --shell bash --line "docker system prune -a --volumes" --cwd "/home/user/app"
      """
    And provider configuration is valid
    When the Explain function is executed
    Then the CLI should output JSON containing "summary" and "risk"
    And the summary should describe that the command removes unused Docker data including volumes
    And the risk should be "high" or "medium"

  Scenario: Explain keybinding in shell
    Given the zsh integration binds Config.keybindings.explain to a ZLE widget
    And the current BUFFER contains a shell command
    When the user presses the explain keybinding
    Then "linesense explain" should be invoked with the BUFFER and cwd
    And the explanation summary should be displayed above the prompt or in a small panel
```

---

### Feature: Failure modes and offline behavior

```gherkin
Feature: Failure modes and offline behavior
  As a user
  I want my shell to remain stable even if AI calls fail
  So that I am never blocked by network or provider issues

  Scenario: Network error from OpenRouter
    Given the OpenRouter API is unreachable
    And the user presses the suggest keybinding
    When "linesense suggest" attempts to call the provider
    Then the program should exit with a non-zero status
    And it should print a concise error message to stderr
    And the shell integration should not modify the current line buffer

  Scenario: Missing API key
    Given the environment variable specified by OpenRouterConfig.api_key_env is not set
    When the user runs "linesense suggest"
    Then the program should exit with an error
    And the error should state that the API key is missing

  Scenario: No suggestions returned
    Given Provider.Suggest returns an empty list
    When the CLI prints the result
    Then the JSON "suggestions" array should be empty
    And the shell integration should not alter the current command line
```

---

This is everything an AI dev agent needs:

* Clear data structures (Key Entities)
* Explicit behavior (Gherkin features)
* Concrete integration points (bash, zsh, OpenRouter, config, history, usage)

If you want, next step I can generate **stub Go package structures** matching this spec (file/folder names + empty types) so an agent can just fill in the logic.
