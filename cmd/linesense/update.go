package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
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

	// Parse versions to compare
	vCurrent, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("failed to parse current version: %w", err)
	}

	vLatest, err := semver.NewVersion(latest.Version())
	if err != nil {
		return fmt.Errorf("failed to parse latest version: %w", err)
	}

	// Check if latest version is actually newer than current
	if !vLatest.GreaterThan(vCurrent) {
		fmt.Printf("Current version %s is the latest (found %s)\n", version, latest.Version())
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
