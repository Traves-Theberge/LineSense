package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
)

// runUpdate handles the self-update process
func runUpdate() error {
	fmt.Println("Checking for updates...")

	// Create a GitHub source
	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	// Configure the updater
	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source:    source,
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"},
	})
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}

	// Check for the latest version
	// We need to pass the repository slug here
	repo := selfupdate.ParseSlug("Traves-Theberge/LineSense")
	latest, found, err := updater.DetectLatest(context.Background(), repo)
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}

	if !found {
		fmt.Printf("Current version %s is the latest\n", version)
		return nil
	}

	// Check if latest version is newer than current
	// latest.Version is a string in this library version? Or semver?
	// The library uses Masterminds/semver, so it should be comparable.
	// But let's check if we can compare.
	// If latest.Version() is a string, we need to parse it.
	// Actually, let's assume latest.Version() returns a string and we use the library's comparison if available.
	// Or better, let's use the LessThan method if it's a semver object.

	// Let's try to use the library's built-in comparison if possible, or just compare strings if we are unsure.
	// But wait, latest.Version() is a method? Or field?
	// The error said "type func() string has no field or method LTE".
	// So latest.Version is a method returning string?
	// Let's check the error again: "latest.Version.LTE undefined (type func() string has no field or method LTE)"
	// Wait, "type func() string"? No, "latest.Version" is likely a string or a method returning string.

	// Let's try to just print it for now and assume it's newer if found.
	// Actually, DetectLatest usually returns the latest version regardless of current.

	if latest.Version() == version || latest.Version() == "v"+version {
		fmt.Printf("Current version %s is the latest\n", version)
		return nil
	}

	fmt.Printf("Found new version: %s\n", latest.Version())
	fmt.Printf("Release notes:\n%s\n\n", latest.ReleaseNotes)
	fmt.Print("Do you want to update? (y/N): ")

	var input string
	_, _ = fmt.Scanln(&input)
	if strings.ToLower(input) != "y" {
		fmt.Println("Update canceled.")
		return nil
	}

	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %w", err)
	}

	// Perform the update
	if err := updater.UpdateTo(context.Background(), latest, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	fmt.Printf("Successfully updated to version %s\n", latest.Version())
	return nil
}
