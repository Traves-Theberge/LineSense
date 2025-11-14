package core

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// CollectHistory reads recent commands from shell history
func CollectHistory(shell string, limit int) ([]HistoryEntry, error) {
	historyPath, err := getHistoryPath(shell)
	if err != nil {
		return nil, err
	}

	// Check if history file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		// No history file, return empty slice (not an error)
		return []HistoryEntry{}, nil
	}

	// Read history file
	file, err := os.Open(historyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read all lines
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Get last N lines
	start := 0
	if len(lines) > limit {
		start = len(lines) - limit
	}
	lines = lines[start:]

	// Parse into HistoryEntry
	var entries []HistoryEntry
	for _, line := range lines {
		entry := parseHistoryLine(shell, line)
		if entry.Command != "" {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// getHistoryPath returns the path to the shell history file
func getHistoryPath(shell string) (string, error) {
	// Check HISTFILE environment variable first
	if histFile := os.Getenv("HISTFILE"); histFile != "" {
		return histFile, nil
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Default history file locations
	switch shell {
	case "bash":
		return filepath.Join(home, ".bash_history"), nil
	case "zsh":
		return filepath.Join(home, ".zsh_history"), nil
	default:
		// Default to bash
		return filepath.Join(home, ".bash_history"), nil
	}
}

// parseHistoryLine parses a single history line based on shell type
func parseHistoryLine(shell string, line string) HistoryEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return HistoryEntry{}
	}

	switch shell {
	case "zsh":
		// Zsh extended history format: ": timestamp:duration;command"
		if strings.HasPrefix(line, ":") {
			parts := strings.SplitN(line[1:], ";", 2)
			if len(parts) == 2 {
				command := strings.TrimSpace(parts[1])
				return HistoryEntry{Command: command}
			}
		}
		// Fallthrough to simple format
		fallthrough
	case "bash":
		fallthrough
	default:
		// Simple format: just the command
		return HistoryEntry{Command: line}
	}
}
