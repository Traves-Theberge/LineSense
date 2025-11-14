package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/traves/linesense/internal/config"
)

func TestBuildContext(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create minimal config
	cfg := &config.Config{
		Context: config.ContextConfig{
			HistoryLength: 10,
			IncludeGit:    true,
			IncludeEnv:    true,
		},
	}

	ctx, err := BuildContext("bash", "test command", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	// Verify basic fields
	if ctx.Shell != "bash" {
		t.Errorf("Shell = %v, want bash", ctx.Shell)
	}
	if ctx.Line != "test command" {
		t.Errorf("Line = %v, want test command", ctx.Line)
	}
	if ctx.CWD != tmpDir {
		t.Errorf("CWD = %v, want %v", ctx.CWD, tmpDir)
	}
}

func TestBuildContext_DisabledFeatures(t *testing.T) {
	tmpDir := t.TempDir()

	// Config with everything disabled
	cfg := &config.Config{
		Context: config.ContextConfig{
			HistoryLength: 0,
			IncludeGit:    false,
			IncludeEnv:    false,
		},
	}

	ctx, err := BuildContext("zsh", "ls", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	// Should still have basic fields
	if ctx.Shell != "zsh" {
		t.Errorf("Shell = %v, want zsh", ctx.Shell)
	}

	// Git should be nil or empty when disabled
	if ctx.Git != nil && ctx.Git.Branch != "" {
		t.Error("Git info should be empty when IncludeGit is false")
	}

	// Environment should be empty when disabled
	if len(ctx.Env) > 0 {
		t.Error("Environment should be empty when IncludeEnv is false")
	}
}

func TestBuildContext_EnvironmentFiltering(t *testing.T) {
	// Set up test environment
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("TEST_PASSWORD", "secret")
	os.Setenv("TEST_API_KEY", "key123")
	os.Setenv("TEST_TOKEN", "token456")
	defer os.Unsetenv("TEST_VAR")
	defer os.Unsetenv("TEST_PASSWORD")
	defer os.Unsetenv("TEST_API_KEY")
	defer os.Unsetenv("TEST_TOKEN")

	tmpDir := t.TempDir()
	cfg := &config.Config{
		Context: config.ContextConfig{
			IncludeEnv: true,
		},
	}

	ctx, err := BuildContext("bash", "ls", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	// Should include non-sensitive vars
	found := false
	for k := range ctx.Env {
		if k == "TEST_VAR" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should include TEST_VAR")
	}

	// Should exclude sensitive vars
	sensitiveVars := []string{"TEST_PASSWORD", "TEST_API_KEY", "TEST_TOKEN"}
	for _, sensitive := range sensitiveVars {
		if _, exists := ctx.Env[sensitive]; exists {
			t.Errorf("Should not include sensitive var: %s", sensitive)
		}
	}
}

func TestCollectHistory_Bash(t *testing.T) {
	// Create temporary bash history
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	histContent := `ls -la
cd /tmp
echo "test"
git status
pwd
`
	if err := os.WriteFile(histFile, []byte(histContent), 0644); err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	// Set HISTFILE env var
	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("bash", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	// Should have commands
	if len(history) == 0 {
		t.Error("History should not be empty")
	}

	// Should contain our test commands
	found := false
	for _, entry := range history {
		if entry.Command == "git status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("History should contain 'git status'")
	}
}

func TestCollectHistory_Zsh(t *testing.T) {
	// Create temporary zsh history with extended format
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".zsh_history")

	// Zsh extended history format: : timestamp:0;command
	histContent := `: 1234567890:0;ls -la
: 1234567891:0;cd /tmp
: 1234567892:0;git status
`
	if err := os.WriteFile(histFile, []byte(histContent), 0644); err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("zsh", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	if len(history) == 0 {
		t.Error("History should not be empty")
	}

	// Check that timestamp is stripped
	for _, entry := range history {
		if entry.Command == "git status" {
			return // Found it
		}
	}
	t.Error("History should contain 'git status' without timestamp")
}

func TestCollectHistory_LimitEntries(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	// Create history with many entries
	histContent := ""
	for i := 0; i < 100; i++ {
		histContent += "command" + string(rune('0'+i%10)) + "\n"
	}

	if err := os.WriteFile(histFile, []byte(histContent), 0644); err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	// Request only 5 entries
	history, err := CollectHistory("bash", 5)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	if len(history) > 5 {
		t.Errorf("History length = %d, want <= 5", len(history))
	}
}

func TestCollectHistory_NonexistentFile(t *testing.T) {
	// Set HISTFILE to non-existent path
	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", "/nonexistent/path/.bash_history")

	history, err := CollectHistory("bash", 10)

	// Should return empty slice, not error
	if err != nil {
		t.Fatalf("CollectHistory() should not error for missing file, got: %v", err)
	}
	if history == nil {
		t.Error("History should not be nil")
	}
	if len(history) != 0 {
		t.Errorf("History should be empty for nonexistent file, got %d entries", len(history))
	}
}

func TestCollectGitInfo_NotARepo(t *testing.T) {
	// Use a temp directory that's not a git repo
	tmpDir := t.TempDir()

	gitInfo, err := CollectGitInfo(tmpDir)

	// Should return nil for non-repo (not an error)
	if err != nil {
		t.Errorf("CollectGitInfo() should not error for non-repo, got: %v", err)
	}
	if gitInfo != nil {
		t.Error("GitInfo should be nil for non-repo")
	}
}

func TestCollectGitInfo_InRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create an initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	// Collect git info
	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if gitInfo == nil {
		t.Fatal("GitInfo should not be nil for git repo")
	}

	// Verify basic fields
	if !gitInfo.IsRepo {
		t.Error("IsRepo should be true")
	}

	// Should have a branch (likely "master" or "main")
	if gitInfo.Branch == "" {
		t.Error("Branch should not be empty")
	}

	// Status should be "clean" since we just committed
	if gitInfo.StatusSummary != "clean" {
		t.Errorf("StatusSummary = %q, want clean", gitInfo.StatusSummary)
	}
}

func TestCollectGitInfo_DirtyRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create and commit initial file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	// Now modify the file (create dirty state)
	if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Collect git info
	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if gitInfo == nil {
		t.Fatal("GitInfo should not be nil")
	}

	// Status should indicate modifications
	if gitInfo.StatusSummary == "clean" {
		t.Error("StatusSummary should not be clean for dirty repo")
	}
	if !strings.Contains(gitInfo.StatusSummary, "modified") {
		t.Errorf("StatusSummary should contain 'modified', got: %q", gitInfo.StatusSummary)
	}
}

func TestCollectGitInfo_WithRemotes(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Add a remote
	exec.Command("git", "-C", tmpDir, "remote", "add", "origin", "https://github.com/test/repo.git").Run()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if len(gitInfo.Remotes) == 0 {
		t.Error("Should have at least one remote")
	}

	found := false
	for _, remote := range gitInfo.Remotes {
		if strings.Contains(remote, "github.com/test/repo.git") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should contain the added remote URL")
	}
}

func TestCollectGitInfo_UntrackedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	// Add untracked file
	untrackedFile := filepath.Join(tmpDir, "untracked.txt")
	os.WriteFile(untrackedFile, []byte("untracked"), 0644)

	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if !strings.Contains(gitInfo.StatusSummary, "untracked") {
		t.Errorf("StatusSummary should contain 'untracked', got: %q", gitInfo.StatusSummary)
	}
}

func TestCollectHistory_DefaultPath(t *testing.T) {
	// Test with default history path (no HISTFILE set)
	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Unsetenv("HISTFILE")

	// This should use default path and gracefully handle if it doesn't exist
	history, err := CollectHistory("bash", 5)

	// Should not error even if file doesn't exist
	if err != nil {
		t.Fatalf("CollectHistory() should not error with default path: %v", err)
	}

	// History may be empty if no history file exists
	if history == nil {
		t.Error("History should not be nil")
	}
}

func TestCollectHistory_ZshWithoutTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".zsh_history")

	// Zsh history without timestamps (simple format)
	histContent := `ls -la
cd /tmp
pwd
`
	if err := os.WriteFile(histFile, []byte(histContent), 0644); err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("zsh", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	if len(history) == 0 {
		t.Error("History should not be empty")
	}

	// Should parse simple format too
	found := false
	for _, entry := range history {
		if entry.Command == "pwd" {
			found = true
			break
		}
	}
	if !found {
		t.Error("History should contain 'pwd'")
	}
}

func TestCollectHistory_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	// Empty history file
	if err := os.WriteFile(histFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("bash", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	if len(history) != 0 {
		t.Errorf("History should be empty for empty file, got %d entries", len(history))
	}
}

func TestBuildContext_HistoryIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	histContent := `git status
git add .
git commit -m "test"
`
	os.WriteFile(histFile, []byte(histContent), 0644)

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	cfg := &config.Config{
		Context: config.ContextConfig{
			HistoryLength: 5,
		},
	}

	ctx, err := BuildContext("bash", "git push", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	if len(ctx.History) == 0 {
		t.Error("Context should include history")
	}

	// Verify history is populated
	found := false
	for _, entry := range ctx.History {
		if strings.Contains(entry.Command, "git commit") {
			found = true
			break
		}
	}
	if !found {
		t.Error("History should contain git commit command")
	}
}

func TestCollectGitInfo_DeletedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	exec.Command("git", "-C", tmpDir, "init").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create and commit files
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	// Delete one file
	os.Remove(file1)

	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if !strings.Contains(gitInfo.StatusSummary, "deleted") {
		t.Errorf("StatusSummary should contain 'deleted', got: %q", gitInfo.StatusSummary)
	}
}

func TestCollectGitInfo_AddedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	exec.Command("git", "-C", tmpDir, "init").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create initial commit
	file1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(file1, []byte("content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Add a new file to staging
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file2, []byte("new content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "file2.txt").Run()

	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	if !strings.Contains(gitInfo.StatusSummary, "added") {
		t.Errorf("StatusSummary should contain 'added', got: %q", gitInfo.StatusSummary)
	}
}

func TestCollectGitInfo_MultipleChanges(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	exec.Command("git", "-C", tmpDir, "init").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create initial commit
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Modify file1
	os.WriteFile(file1, []byte("modified content"), 0644)

	// Delete file2
	os.Remove(file2)

	// Add file3 to staging
	file3 := filepath.Join(tmpDir, "file3.txt")
	os.WriteFile(file3, []byte("new"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "file3.txt").Run()

	// Add untracked file4
	file4 := filepath.Join(tmpDir, "file4.txt")
	os.WriteFile(file4, []byte("untracked"), 0644)

	gitInfo, err := CollectGitInfo(tmpDir)
	if err != nil {
		t.Fatalf("CollectGitInfo() error = %v", err)
	}

	summary := gitInfo.StatusSummary
	// Should contain all types of changes
	if !strings.Contains(summary, "modified") {
		t.Errorf("StatusSummary should contain 'modified', got: %q", summary)
	}
	if !strings.Contains(summary, "deleted") {
		t.Errorf("StatusSummary should contain 'deleted', got: %q", summary)
	}
	if !strings.Contains(summary, "added") {
		t.Errorf("StatusSummary should contain 'added', got: %q", summary)
	}
	if !strings.Contains(summary, "untracked") {
		t.Errorf("StatusSummary should contain 'untracked', got: %q", summary)
	}
}

func TestCollectHistory_UnknownShell(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	histContent := `ls
pwd
`
	os.WriteFile(histFile, []byte(histContent), 0644)

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	// Use unknown shell - should default to bash behavior
	history, err := CollectHistory("fish", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	if len(history) == 0 {
		t.Error("History should not be empty for unknown shell")
	}
}

func TestCollectHistory_ZshMalformedTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".zsh_history")

	// Malformed zsh history - timestamp without command
	histContent := `: 1234567890:0;
: malformed
regular command
`
	os.WriteFile(histFile, []byte(histContent), 0644)

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("zsh", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	// Should handle malformed lines gracefully
	if len(history) == 0 {
		t.Error("History should contain at least the regular command")
	}

	// Check that "regular command" was parsed
	found := false
	for _, entry := range history {
		if entry.Command == "regular command" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should have parsed the regular command line")
	}
}

func TestCollectHistory_WithBlankLines(t *testing.T) {
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")

	histContent := `ls

pwd

cd /tmp

`
	os.WriteFile(histFile, []byte(histContent), 0644)

	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	history, err := CollectHistory("bash", 10)
	if err != nil {
		t.Fatalf("CollectHistory() error = %v", err)
	}

	// Should skip blank lines
	if len(history) != 3 {
		t.Errorf("Expected 3 commands (skipping blanks), got %d", len(history))
	}
}

func TestBuildContext_GitIntegration(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tmpDir := t.TempDir()
	exec.Command("git", "-C", tmpDir, "init").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	cfg := &config.Config{
		Context: config.ContextConfig{
			IncludeGit: true,
		},
	}

	ctx, err := BuildContext("bash", "git status", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	if ctx.Git == nil {
		t.Fatal("Git info should be included")
	}

	if !ctx.Git.IsRepo {
		t.Error("IsRepo should be true")
	}

	if ctx.Git.Branch == "" {
		t.Error("Branch should not be empty")
	}
}

func TestBuildContext_AllFeatures(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	// Set up environment
	os.Setenv("TEST_VISIBLE_VAR", "visible")
	os.Setenv("TEST_SECRET", "hidden")
	defer os.Unsetenv("TEST_VISIBLE_VAR")
	defer os.Unsetenv("TEST_SECRET")

	// Set up history
	tmpDir := t.TempDir()
	histFile := filepath.Join(tmpDir, ".bash_history")
	os.WriteFile(histFile, []byte("ls\npwd\n"), 0644)
	originalHist := os.Getenv("HISTFILE")
	defer os.Setenv("HISTFILE", originalHist)
	os.Setenv("HISTFILE", histFile)

	// Set up git repo
	exec.Command("git", "-C", tmpDir, "init").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	cfg := &config.Config{
		Context: config.ContextConfig{
			HistoryLength: 10,
			IncludeGit:    true,
			IncludeEnv:    true,
		},
	}

	ctx, err := BuildContext("bash", "test", tmpDir, cfg)
	if err != nil {
		t.Fatalf("BuildContext() error = %v", err)
	}

	// Verify all features are populated
	if ctx.Git == nil {
		t.Error("Git should be populated")
	}

	if len(ctx.History) == 0 {
		t.Error("History should be populated")
	}

	if len(ctx.Env) == 0 {
		t.Error("Environment should be populated")
	}

	// Check env filtering works
	if _, exists := ctx.Env["TEST_SECRET"]; exists {
		t.Error("Sensitive env var should be filtered")
	}
}
