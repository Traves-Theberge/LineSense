package core

import (
	"runtime"
	"testing"
)

func TestDetectOS(t *testing.T) {
	os := DetectOS()

	// Should return a valid OS type
	validOS := map[string]bool{
		"linux":   true,
		"darwin":  true,
		"windows": true,
	}

	if !validOS[os] && os != runtime.GOOS {
		t.Errorf("DetectOS() = %q, expected linux, darwin, windows, or %q", os, runtime.GOOS)
	}

	// Should match runtime.GOOS for standard systems
	expectedOS := runtime.GOOS
	if expectedOS == "darwin" || expectedOS == "linux" || expectedOS == "windows" {
		if os != expectedOS {
			t.Errorf("DetectOS() = %q, expected %q", os, expectedOS)
		}
	}
}

func TestDetectDistribution(t *testing.T) {
	distro := DetectDistribution()

	// On non-Linux systems, should return empty string
	if runtime.GOOS != "linux" {
		if distro != "" {
			t.Errorf("DetectDistribution() on %s = %q, expected empty string", runtime.GOOS, distro)
		}
		return
	}

	// On Linux, we can't predict the exact distro, but we can verify format
	if distro != "" {
		// Should be lowercase
		if distro != "" && distro[0] >= 'A' && distro[0] <= 'Z' {
			t.Errorf("DetectDistribution() = %q, expected lowercase distro name", distro)
		}
	}
	// Note: distro might be empty if /etc/os-release doesn't exist or can't be read
}

func TestDetectPackageManager(t *testing.T) {
	pm := DetectPackageManager()

	// Valid package managers by OS
	validLinux := map[string]bool{
		"apt":    true,
		"dnf":    true,
		"yum":    true,
		"pacman": true,
		"zypper": true,
		"apk":    true,
		"":       true, // Empty is valid if none found
	}

	validDarwin := map[string]bool{
		"brew": true,
		"":     true, // Empty is valid if brew not installed
	}

	validWindows := map[string]bool{
		"choco":  true,
		"winget": true,
		"scoop":  true,
		"":       true, // Empty is valid if none found
	}

	switch runtime.GOOS {
	case "linux":
		if !validLinux[pm] {
			t.Errorf("DetectPackageManager() on Linux = %q, expected one of: apt, dnf, yum, pacman, zypper, apk, or empty", pm)
		}
	case "darwin":
		if !validDarwin[pm] {
			t.Errorf("DetectPackageManager() on macOS = %q, expected brew or empty", pm)
		}
	case "windows":
		if !validWindows[pm] {
			t.Errorf("DetectPackageManager() on Windows = %q, expected choco, winget, scoop, or empty", pm)
		}
	}
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name    string
		command string
		// We can't predict if specific commands exist, so we just test the function works
	}{
		{"check ls", "ls"},
		{"check nonexistent", "this-command-definitely-does-not-exist-12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			_ = commandExists(tt.command)
		})
	}

	// Test that a command that definitely doesn't exist returns false
	if commandExists("this-command-definitely-does-not-exist-12345-xyz") {
		t.Error("commandExists should return false for non-existent command")
	}
}
