package core

import (
	"os/exec"
	"strings"
)

// CollectGitInfo gathers git repository information from the current directory
func CollectGitInfo(cwd string) (*GitInfo, error) {
	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		// Not a git repository, return nil (not an error)
		return nil, nil
	}

	info := &GitInfo{
		IsRepo: true,
	}

	// Get current branch name
	if branch, err := gitCommand(cwd, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		info.Branch = strings.TrimSpace(branch)
	}

	// Get status summary
	if status, err := gitCommand(cwd, "status", "--porcelain"); err == nil {
		info.StatusSummary = summarizeGitStatus(status)
	}

	// Get remote URLs
	if remotes, err := gitCommand(cwd, "remote", "-v"); err == nil {
		info.Remotes = parseGitRemotes(remotes)
	}

	return info, nil
}

// gitCommand executes a git command and returns its output
func gitCommand(cwd string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// summarizeGitStatus creates a human-readable summary of git status
func summarizeGitStatus(porcelain string) string {
	if strings.TrimSpace(porcelain) == "" {
		return "clean"
	}

	lines := strings.Split(strings.TrimSpace(porcelain), "\n")
	modified := 0
	added := 0
	deleted := 0
	untracked := 0

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		status := line[0:2]
		switch {
		case strings.HasPrefix(status, "??"):
			untracked++
		case strings.HasPrefix(status, "A") || strings.HasPrefix(status, " A"):
			added++
		case strings.HasPrefix(status, "D") || strings.HasPrefix(status, " D"):
			deleted++
		case strings.HasPrefix(status, "M") || strings.HasPrefix(status, " M"):
			modified++
		}
	}

	var parts []string
	if modified > 0 {
		parts = append(parts, "modified")
	}
	if added > 0 {
		parts = append(parts, "added")
	}
	if deleted > 0 {
		parts = append(parts, "deleted")
	}
	if untracked > 0 {
		parts = append(parts, "untracked")
	}

	if len(parts) == 0 {
		return "uncommitted changes"
	}
	return strings.Join(parts, ", ")
}

// parseGitRemotes extracts unique remote URLs from git remote -v output
func parseGitRemotes(remotesOutput string) []string {
	lines := strings.Split(strings.TrimSpace(remotesOutput), "\n")
	seen := make(map[string]bool)
	var remotes []string

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			url := fields[1]
			if !seen[url] {
				seen[url] = true
				remotes = append(remotes, url)
			}
		}
	}

	return remotes
}
