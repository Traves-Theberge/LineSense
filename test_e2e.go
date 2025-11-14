package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/traves/linesense/internal/ai"
	"github.com/traves/linesense/internal/config"
	"github.com/traves/linesense/internal/core"
)

func main() {
	// Try to load .env file (silently ignore if not present)
	if cwd, err := os.Getwd(); err == nil {
		godotenv.Load(filepath.Join(cwd, ".env"))
	}

	fmt.Println("🧪 LineSense End-to-End Test")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Check for API key
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		log.Fatal("❌ OPENROUTER_API_KEY environment variable not set")
	}
	fmt.Println("✓ API key found")

	// Load configuration
	fmt.Println("\n1️⃣  Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	fmt.Printf("   ✓ Config loaded (history_length: %d)\n", cfg.Context.HistoryLength)

	providersCfg, err := config.LoadProvidersConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load providers config: %v", err)
	}
	fmt.Printf("   ✓ Providers config loaded (model: %s)\n", providersCfg.Default.Model)

	// Create provider
	fmt.Println("\n2️⃣  Creating AI provider...")
	provider, err := ai.NewProvider(providersCfg, cfg.AI.ProviderProfile)
	if err != nil {
		log.Fatalf("❌ Failed to create provider: %v", err)
	}
	fmt.Printf("   ✓ Provider created (%s)\n", provider.Name())

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Failed to get working directory: %v", err)
	}

	// Build context
	fmt.Println("\n3️⃣  Building context...")
	ctx := context.Background()
	contextEnv, err := core.BuildContext("bash", "list files", cwd, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to build context: %v", err)
	}
	fmt.Printf("   ✓ Context built\n")
	fmt.Printf("     - Shell: %s\n", contextEnv.Shell)
	fmt.Printf("     - CWD: %s\n", contextEnv.CWD)
	if contextEnv.Git != nil {
		fmt.Printf("     - Git: branch=%s, status=%s\n", contextEnv.Git.Branch, contextEnv.Git.StatusSummary)
	}
	if len(contextEnv.History) > 0 {
		fmt.Printf("     - History: %d commands\n", len(contextEnv.History))
	}

	// Test Suggest
	fmt.Println("\n4️⃣  Testing Suggest operation...")
	fmt.Println("   Input: 'list files'")
	suggestInput := core.SuggestInput{
		ModelID: "", // Use default from profile
		Prompt:  "list files",
		Context: contextEnv,
	}

	suggestions, err := provider.Suggest(ctx, suggestInput)
	if err != nil {
		log.Fatalf("❌ Suggest failed: %v", err)
	}

	if len(suggestions) == 0 {
		log.Fatal("❌ No suggestions returned")
	}

	fmt.Printf("   ✓ Received %d suggestion(s)\n", len(suggestions))
	for i, suggestion := range suggestions {
		fmt.Printf("\n   Suggestion %d:\n", i+1)
		fmt.Printf("     Command: %s\n", suggestion.Command)
		fmt.Printf("     Risk: %s\n", suggestion.Risk)
		fmt.Printf("     Source: %s\n", suggestion.Source)
		if suggestion.Explanation != "" {
			fmt.Printf("     Explanation: %s\n", suggestion.Explanation)
		}
	}

	// Test Explain
	fmt.Println("\n5️⃣  Testing Explain operation...")
	explainLine := suggestions[0].Command // Explain the suggested command
	fmt.Printf("   Input: '%s'\n", explainLine)

	contextEnv.Line = explainLine // Update context with new line
	explainInput := core.ExplainInput{
		ModelID: "",
		Prompt:  explainLine,
		Context: contextEnv,
	}

	explanation, err := provider.Explain(ctx, explainInput)
	if err != nil {
		log.Fatalf("❌ Explain failed: %v", err)
	}

	fmt.Printf("   ✓ Explanation received\n")
	fmt.Printf("\n   Summary: %s\n", explanation.Summary)
	fmt.Printf("   Risk: %s\n", explanation.Risk)
	if len(explanation.Notes) > 0 {
		fmt.Println("   Notes:")
		for _, note := range explanation.Notes {
			fmt.Printf("     - %s\n", note)
		}
	}

	// Success!
	fmt.Println("\n" + string(make([]byte, 50)))
	fmt.Println("✅ All tests passed!")
	fmt.Println("\nLineSense is working end-to-end:")
	fmt.Println("  ✓ Configuration loading")
	fmt.Println("  ✓ Context gathering (git, history, env)")
	fmt.Println("  ✓ AI provider integration")
	fmt.Println("  ✓ Command suggestions")
	fmt.Println("  ✓ Command explanations")
	fmt.Println("\nReady for CLI implementation! 🚀")
}
