package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/traves/linesense/internal/ai"
	"github.com/traves/linesense/internal/config"
	"github.com/traves/linesense/internal/core"
)

const version = "0.6.6"

func main() {
	// Load LineSense .env file from config directory (secure location)
	loadSecureEnv()

	// Also try to load .env file from current directory (for development)
	if cwd, err := os.Getwd(); err == nil {
		_ = godotenv.Load(filepath.Join(cwd, ".env"))
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// loadSecureEnv loads the API key from ~/.config/linesense/.env
// Falls back to environment variables if not present
func loadSecureEnv() {
	configDir := config.GetConfigDir()
	envPath := filepath.Join(configDir, ".env")

	// Silently load if exists
	_ = godotenv.Load(envPath)
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return fmt.Errorf("no command specified")
	}

	command := os.Args[1]

	switch command {
	case "suggest":
		return runSuggest(os.Args[2:])
	case "explain":
		return runExplain(os.Args[2:])
	case "config":
		return runConfig(os.Args[2:])
	case "update":
		return runUpdate()
	case "version", "--version", "-v":
		fmt.Printf("linesense version %s\n", version)
		return nil
	case "help", "--help", "-h":
		printUsage()
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `linesense - AI-powered shell command assistant

Usage:
  linesense config [subcommand]  Configure LineSense
  linesense suggest [flags]      Generate command suggestions
  linesense explain [flags]      Explain a command
  linesense update               Update LineSense to the latest version
  linesense version              Show version information
  linesense help                 Show this help message

Config Subcommands:
  init            Initialize configuration with interactive setup
  init-project    Initialize project-specific context in current directory
  set-key         Set OpenRouter API key securely
  set-model       Change the default model
  edit            Open configuration file in default editor
  show            Display current configuration

Suggest Flags:
  --shell string     Shell type (bash, zsh) (default: auto-detect)
  --line string      Partial command line to complete (required)
  --cwd string       Current working directory (default: current directory)
  --model string     Override model ID from config
  --format string    Output format: pretty or json (default: pretty)

Explain Flags:
  --shell string     Shell type (bash, zsh) (default: auto-detect)
  --line string      Command to explain (required)
  --cwd string       Current working directory (default: current directory)
  --model string     Override model ID from config
  --format string    Output format: pretty or json (default: pretty)

Examples:
  linesense suggest --line "list files"
  linesense explain --line "rm -rf /"
  linesense suggest --line "git com" --shell bash
  linesense explain --line "docker ps -a" --model gpt-4

Configuration:
  Config files: ~/.config/linesense/config.toml
                ~/.config/linesense/providers.toml
`)
}

func runSuggest(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("suggest", flag.ExitOnError)
	shell := fs.String("shell", "", "Shell type (bash, zsh)")
	line := fs.String("line", "", "Partial command line to complete")
	cwd := fs.String("cwd", "", "Current working directory")
	model := fs.String("model", "", "Override model ID from config")
	format := fs.String("format", "pretty", "Output format: json or pretty")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *line == "" {
		return fmt.Errorf("--line flag is required")
	}

	// Auto-detect shell if not provided
	if *shell == "" {
		*shell = detectShell()
	}

	// Use current directory if not provided
	if *cwd == "" {
		var err error
		*cwd, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	providersCfg, err := config.LoadProvidersConfig()
	if err != nil {
		return fmt.Errorf("failed to load providers config: %w", err)
	}

	// Create provider
	provider, err := ai.NewProvider(providersCfg, cfg.AI.ProviderProfile)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Build context
	contextEnv, err := core.BuildContext(*shell, *line, *cwd, cfg)
	if err != nil {
		return fmt.Errorf("failed to build context: %w", err)
	}

	// Create suggest input
	input := core.SuggestInput{
		ModelID: *model,
		Prompt:  *line,
		Context: contextEnv,
	}

	// Generate suggestions with spinner
	var suggestions []core.Suggestion
	err = withSpinner("Generating suggestions...", func(ctx context.Context) error {
		var err error
		suggestions, err = provider.Suggest(ctx, input)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to generate suggestions: %w", err)
	}

	// Output based on format
	if *format == "json" {
		output := map[string]interface{}{
			"suggestions": suggestions,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	// Pretty format (default) with styled output
	printSuggestionsStyled(suggestions)
	return nil
}

func runExplain(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("explain", flag.ExitOnError)
	shell := fs.String("shell", "", "Shell type (bash, zsh)")
	line := fs.String("line", "", "Command to explain")
	cwd := fs.String("cwd", "", "Current working directory")
	model := fs.String("model", "", "Override model ID from config")
	format := fs.String("format", "pretty", "Output format: json or pretty")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if *line == "" {
		return fmt.Errorf("--line flag is required")
	}

	// Auto-detect shell if not provided
	if *shell == "" {
		*shell = detectShell()
	}

	// Use current directory if not provided
	if *cwd == "" {
		var err error
		*cwd, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	providersCfg, err := config.LoadProvidersConfig()
	if err != nil {
		return fmt.Errorf("failed to load providers config: %w", err)
	}

	// Create provider
	provider, err := ai.NewProvider(providersCfg, cfg.AI.ProviderProfile)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Build context
	contextEnv, err := core.BuildContext(*shell, *line, *cwd, cfg)
	if err != nil {
		return fmt.Errorf("failed to build context: %w", err)
	}

	// Create explain input
	input := core.ExplainInput{
		ModelID: *model,
		Prompt:  *line,
		Context: contextEnv,
	}

	// Generate explanation with spinner
	var explanation core.Explanation
	err = withSpinner("Analyzing command...", func(ctx context.Context) error {
		var err error
		explanation, err = provider.Explain(ctx, input)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to generate explanation: %w", err)
	}

	// Output based on format
	if *format == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(explanation)
	}

	// Pretty format (default) with styled output
	printExplanationStyled(explanation)
	return nil
}

// detectShell attempts to auto-detect the current shell
func detectShell() string {
	// Try SHELL environment variable
	if shell := os.Getenv("SHELL"); shell != "" {
		if strings.Contains(shell, "zsh") {
			return "zsh"
		}
		if strings.Contains(shell, "bash") {
			return "bash"
		}
	}

	// Default to bash
	return "bash"
}

// runConfig handles configuration management
func runConfig(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("config subcommand required (init, set-key, set-model, show)")
	}

	subcommand := args[0]

	switch subcommand {
	case "init":
		return runConfigInit()
	case "init-project":
		return runConfigInitProject()
	case "set-key":
		return runConfigSetKey(args[1:])
	case "set-model":
		return runConfigSetModel(args[1:])
	case "edit":
		return runConfigEdit()
	case "show":
		return runConfigShow()
	default:
		return fmt.Errorf("unknown config subcommand: %s", subcommand)
	}
}

// runConfigEdit opens the config file in the default editor
func runConfigEdit() error {
	configDir := config.GetConfigDir()
	configPath := filepath.Join(configDir, "config.toml")

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // Fallback to nano
		if _, err := os.Stat("/usr/bin/vim"); err == nil {
			editor = "vim"
		}
		if _, err := os.Stat("/usr/bin/vi"); err == nil {
			editor = "vi"
		}
	}

	fmt.Printf("Opening global configuration %s in %s...\n", configPath, editor)

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// runConfigInitProject initializes a project-specific context file
func runConfigInitProject() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	contextPath := filepath.Join(cwd, ".linesense_context")

	if _, err := os.Stat(contextPath); err == nil {
		fmt.Printf("⚠️  Project context file already exists at %s\n", contextPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("operation canceled")
		}
	}

	content := `# LineSense Project Context
# Add project-specific instructions here.
# The AI will use this context when generating suggestions in this directory.

# Example:
# - Use 'npm run build' instead of 'make'
# - The main branch is 'develop'
# - Always run tests with '--race'
`

	if err := os.WriteFile(contextPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create project context file: %w", err)
	}

	fmt.Printf("✓ Created project-specific context file at %s\n", contextPath)
	fmt.Println("  Edit this file to add project-specific instructions.")
	return nil
}

// runConfigInit initializes configuration with interactive setup
func runConfigInit() error {
	fmt.Println("🚀 LineSense Configuration Setup")
	fmt.Println("================================")
	fmt.Println()

	// Get config directory
	configDir := config.GetConfigDir()
	fmt.Printf("Configuration directory: %s\n", configDir)
	fmt.Println()

	// Check if config already exists
	configPath := filepath.Join(configDir, "config.toml")
	providersPath := filepath.Join(configDir, "providers.toml")

	if _, err := os.Stat(configPath); err == nil {
		// Check for interactive terminal
		fileInfo, _ := os.Stdin.Stat()
		if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("⚠️  Configuration files already exist. Skipping overwrite (non-interactive mode).")
			return nil
		}

		fmt.Println("⚠️  Configuration files already exist.")
		fmt.Print("Do you want to overwrite them? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("configuration setup canceled")
		}
	}

	// Create config directory
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Copy example configs
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Try to find examples directory
	examplesDir := filepath.Join(exeDir, "..", "examples")
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		// Try current directory
		examplesDir = "examples"
	}

	// Copy config.toml
	if err := copyFile(filepath.Join(examplesDir, "config.toml"), configPath); err != nil {
		return fmt.Errorf("failed to copy config.toml: %w", err)
	}
	fmt.Printf("✓ Created %s\n", configPath)

	// Copy providers.toml
	if err := copyFile(filepath.Join(examplesDir, "providers.toml"), providersPath); err != nil {
		return fmt.Errorf("failed to copy providers.toml: %w", err)
	}
	fmt.Printf("✓ Created %s\n", providersPath)

	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Println("1. Set your OpenRouter API key:")
	fmt.Println("   linesense config set-key YOUR_API_KEY")
	fmt.Println()
	fmt.Println("2. (Optional) Change the default model:")
	fmt.Println("   linesense config set-model openai/gpt-4o")
	fmt.Println()
	fmt.Println("3. Test it:")
	fmt.Println("   linesense suggest --line \"list files\"")

	return nil
}

// runConfigSetKey sets the OpenRouter API key securely
func runConfigSetKey(args []string) error {
	var apiKey string

	if len(args) > 0 {
		apiKey = args[0]
	} else {
		// Read from stdin (more secure)
		fmt.Print("Enter your OpenRouter API key: ")
		_, _ = fmt.Scanln(&apiKey)
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Store in secure .env file
	configDir := config.GetConfigDir()
	envPath := filepath.Join(configDir, ".env")

	// Ensure config directory exists with secure permissions
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if .env file already exists
	var existingContent []byte
	var hasExistingKey bool

	if content, err := os.ReadFile(envPath); err == nil {
		existingContent = content
		hasExistingKey = strings.Contains(string(content), "OPENROUTER_API_KEY")
	}

	// Check for legacy shell rc file installations
	homeDir, _ := os.UserHomeDir()
	shell := detectShell()
	var rcFile string
	switch shell {
	case "zsh":
		rcFile = filepath.Join(homeDir, ".zshrc")
	case "bash":
		rcFile = filepath.Join(homeDir, ".bashrc")
	default:
		rcFile = filepath.Join(homeDir, ".bashrc")
	}

	var hasLegacyKey bool
	if rcContent, err := os.ReadFile(rcFile); err == nil {
		hasLegacyKey = strings.Contains(string(rcContent), "OPENROUTER_API_KEY")
	}

	// Warn about migration if legacy key exists
	if hasLegacyKey && !hasExistingKey {
		fmt.Println("⚠️  MIGRATION NOTICE:")
		fmt.Printf("   An API key was found in %s\n", rcFile)
		fmt.Println("   LineSense now stores API keys in a secure file with restricted permissions.")
		fmt.Printf("   After setup, you can remove the old export from %s\n", rcFile)
		fmt.Println()
	}

	if hasExistingKey {
		fmt.Printf("⚠️  OPENROUTER_API_KEY already exists in %s\n", envPath)
		fmt.Print("Do you want to update it? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("API key setup canceled")
		}

		// Remove old key from .env content
		lines := strings.Split(string(existingContent), "\n")
		var newLines []string
		for _, line := range lines {
			if !strings.Contains(line, "OPENROUTER_API_KEY") && !strings.Contains(line, "# LineSense OpenRouter API Key") {
				newLines = append(newLines, line)
			}
		}
		existingContent = []byte(strings.Join(newLines, "\n"))
	}

	// Create new .env content
	var newContent string
	if len(existingContent) > 0 {
		newContent = string(existingContent)
		if !strings.HasSuffix(newContent, "\n") {
			newContent += "\n"
		}
	}
	newContent += fmt.Sprintf("# LineSense OpenRouter API Key\nOPENROUTER_API_KEY=%s\n", apiKey)

	// Write with secure permissions (0600 = owner read/write only)
	if err := os.WriteFile(envPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", envPath, err)
	}

	fmt.Printf("✓ API key saved securely to %s\n", envPath)
	fmt.Println("  (File permissions: 0600 - owner read/write only)")
	fmt.Println()

	if hasLegacyKey {
		fmt.Println("📝 Optional cleanup:")
		fmt.Printf("   You can now remove the old export from %s\n", rcFile)
		fmt.Println("   Look for lines containing 'OPENROUTER_API_KEY'")
		fmt.Println()
	}

	fmt.Println("✓ Setup complete! API key is loaded automatically.")

	return nil
}

// runConfigSetModel changes the default model
func runConfigSetModel(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("model name required")
	}

	model := args[0]

	// Load current config
	configDir := config.GetConfigDir()
	providersPath := filepath.Join(configDir, "providers.toml")

	content, err := os.ReadFile(providersPath)
	if err != nil {
		return fmt.Errorf("failed to read providers config: %w", err)
	}

	// Update model line
	lines := strings.Split(string(content), "\n")
	var updated bool
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "model =") {
			lines[i] = fmt.Sprintf("model = \"%s\"", model)
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("could not find model setting in config")
	}

	// Write back
	if err := os.WriteFile(providersPath, []byte(strings.Join(lines, "\n")), 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("✓ Default model updated to: %s\n", model)
	fmt.Println()
	fmt.Println("Popular models:")
	fmt.Println("  openai/gpt-4o-mini       - Fast and cheap")
	fmt.Println("  openai/gpt-4o            - Most capable")
	fmt.Println("  meta-llama/llama-3.1-8b  - Open source, fast")

	return nil
}

// runConfigShow displays current configuration
func runConfigShow() error {
	fmt.Println("📋 LineSense Configuration")
	fmt.Println("==========================")
	fmt.Println()

	// Config directory
	configDir := config.GetConfigDir()
	fmt.Printf("Config directory: %s\n", configDir)
	fmt.Println()

	// Check API key
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	envPath := filepath.Join(configDir, ".env")

	if apiKey != "" {
		masked := apiKey[:8] + "..." + apiKey[len(apiKey)-4:]

		// Check source
		if _, err := os.Stat(envPath); err == nil {
			fmt.Printf("API Key: %s ✓\n", masked)
			fmt.Printf("  Location: %s (secure)\n", envPath)
		} else {
			fmt.Printf("API Key: %s ✓\n", masked)
			fmt.Println("  Location: Environment variable (consider migrating to secure storage)")
			fmt.Println("  Run: linesense config set-key")
		}
	} else {
		fmt.Println("API Key: Not set ❌")
		fmt.Println("  Set with: linesense config set-key")
	}
	fmt.Println()

	// Try to load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Config file: Error loading (%v)\n", err)
		return nil
	}

	fmt.Println("Configuration:")
	fmt.Printf("  Provider profile: %s\n", cfg.AI.ProviderProfile)
	fmt.Printf("  History length: %d\n", cfg.Context.HistoryLength)
	fmt.Printf("  Include git: %v\n", cfg.Context.IncludeGit)
	fmt.Printf("  Include env: %v\n", cfg.Context.IncludeEnv)
	fmt.Println()

	// Try to load providers config
	providersCfg, err := config.LoadProvidersConfig()
	if err != nil {
		fmt.Printf("Providers config: Error loading (%v)\n", err)
		return nil
	}

	profile, err := providersCfg.GetProfile(cfg.AI.ProviderProfile)
	if err != nil {
		fmt.Printf("Provider profile: Error (%v)\n", err)
		return nil
	}

	fmt.Println("Provider settings:")
	fmt.Printf("  Model: %s\n", profile.Model)
	fmt.Printf("  Temperature: %.1f\n", profile.Temperature)
	fmt.Printf("  Max tokens: %d\n", profile.MaxTokens)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0600)
}
