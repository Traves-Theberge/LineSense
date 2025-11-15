package core

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DetectOS returns the operating system type
func DetectOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return runtime.GOOS
	}
}

// DetectDistribution detects the Linux distribution
// Returns empty string for non-Linux systems or if detection fails
func DetectDistribution() string {
	if runtime.GOOS != "linux" {
		return ""
	}

	// Try to read /etc/os-release (standard location)
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	// Parse os-release file
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			// Extract distro ID (e.g., ID=ubuntu)
			distro := strings.TrimPrefix(line, "ID=")
			distro = strings.Trim(distro, "\"")
			return strings.ToLower(distro)
		}
	}

	return ""
}

// DetectPackageManager detects the available package manager
// Returns the most common package manager for the system
func DetectPackageManager() string {
	osType := runtime.GOOS

	switch osType {
	case "darwin":
		// macOS - check for brew
		if commandExists("brew") {
			return "brew"
		}
		return ""

	case "linux":
		// Check for common Linux package managers in order of preference
		managers := []string{
			"apt",     // Debian/Ubuntu
			"dnf",     // Fedora/RHEL 8+
			"yum",     // RHEL/CentOS 7
			"pacman",  // Arch
			"zypper",  // openSUSE
			"apk",     // Alpine
		}

		for _, manager := range managers {
			if commandExists(manager) {
				return manager
			}
		}
		return ""

	case "windows":
		// Windows - check for common package managers
		if commandExists("choco") {
			return "choco"
		}
		if commandExists("winget") {
			return "winget"
		}
		if commandExists("scoop") {
			return "scoop"
		}
		return ""

	default:
		return ""
	}
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
