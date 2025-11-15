package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/traves/linesense/internal/core"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Color scheme using Lipgloss
var (
	// Primary colors
	primaryColor   = lipgloss.Color("39")  // Bright blue
	secondaryColor = lipgloss.Color("205") // Pink
	successColor   = lipgloss.Color("42")  // Green
	warningColor   = lipgloss.Color("226") // Yellow
	dangerColor    = lipgloss.Color("196") // Red
	mutedColor     = lipgloss.Color("241") // Gray

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 1)

	riskLowStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	riskMediumStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	riskHighStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	headerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginTop(1)
)

// spinnerModel is the Bubble Tea model for the loading spinner
type spinnerModel struct {
	spinner  spinner.Model
	message  string
	quitting bool
	done     bool
}

func newSpinnerModel(message string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)
	return spinnerModel{
		spinner: s,
		message: message,
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		if !m.done {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case doneMsg:
		m.done = true
		m.quitting = true
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n  %s %s\n\n", m.spinner.View(), m.message)
}

type doneMsg struct{}

// showSpinner displays a loading spinner while executing a function
func showSpinner(message string, fn func() error) error {
	done := make(chan error, 1)

	p := tea.NewProgram(newSpinnerModel(message))

	// Run the function in a goroutine
	go func() {
		done <- fn()
	}()

	// Start the spinner in another goroutine
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running spinner: %v\n", err)
		}
	}()

	// Wait for function to complete
	err := <-done

	// Signal spinner to quit
	p.Send(doneMsg{})

	// Give spinner time to quit gracefully
	time.Sleep(50 * time.Millisecond)

	return err
}

// printSuggestionsStyled prints command suggestions with Lipgloss styling
func printSuggestionsStyled(suggestions []core.Suggestion) {
	if len(suggestions) == 0 {
		fmt.Println(mutedStyle.Render("\n  No suggestions found.\n"))
		return
	}

	// Header
	header := titleStyle.Render("💡 Command Suggestions")
	divider := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n%s\n", header, mutedStyle.Render(divider))

	for i, suggestion := range suggestions {
		// Build the suggestion box content
		var parts []string

		// Number and command
		number := lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Render(fmt.Sprintf("%d.", i+1))

		command := commandStyle.Render(suggestion.Command)
		parts = append(parts, fmt.Sprintf("%s %s", number, command))

		// Risk indicator
		var riskStyle lipgloss.Style
		var riskIcon string
		switch suggestion.Risk {
		case "low":
			riskStyle = riskLowStyle
			riskIcon = "✓"
		case "medium":
			riskStyle = riskMediumStyle
			riskIcon = "⚠"
		case "high":
			riskStyle = riskHighStyle
			riskIcon = "⚠"
		default:
			riskStyle = mutedStyle
			riskIcon = "•"
		}

		risk := fmt.Sprintf("   %s Risk: %s",
			riskStyle.Render(riskIcon),
			riskStyle.Render(string(suggestion.Risk)),
		)
		parts = append(parts, risk)

		// Explanation
		if suggestion.Explanation != "" {
			explanation := mutedStyle.Render(fmt.Sprintf("   %s", suggestion.Explanation))
			parts = append(parts, explanation)
		}

		fmt.Printf("\n%s\n", strings.Join(parts, "\n"))
	}

	fmt.Println()
}

// printExplanationStyled prints a command explanation with Lipgloss styling
func printExplanationStyled(explanation core.Explanation) {
	// Header
	header := titleStyle.Render("📖 Command Explanation")
	divider := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n%s\n\n", header, mutedStyle.Render(divider))

	// Summary box
	summaryTitle := headerStyle.Render("Summary")
	summaryText := lipgloss.NewStyle().
		Width(70).
		Render(explanation.Summary)
	summaryBox := boxStyle.Render(fmt.Sprintf("%s\n\n%s", summaryTitle, summaryText))
	fmt.Println(summaryBox)

	// Risk level
	var riskStyle lipgloss.Style
	var riskIcon string
	switch explanation.Risk {
	case "low":
		riskStyle = riskLowStyle
		riskIcon = "✓"
	case "medium":
		riskStyle = riskMediumStyle
		riskIcon = "⚠"
	case "high":
		riskStyle = riskHighStyle
		riskIcon = "⚠"
	default:
		riskStyle = mutedStyle
		riskIcon = "•"
	}

	riskBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(riskStyle.GetForeground()).
		Padding(0, 1).
		Render(fmt.Sprintf("%s Risk Level: %s",
			riskStyle.Render(riskIcon),
			riskStyle.Render(string(explanation.Risk)),
		))
	fmt.Println(riskBox)

	// Details
	if len(explanation.Notes) > 0 {
		detailsTitle := headerStyle.Render("\nDetails")
		fmt.Println(detailsTitle)

		var detailsParts []string
		for _, note := range explanation.Notes {
			// Check if this is a header (no leading spaces/punctuation)
			if len(note) > 0 && note[0] != ' ' && note[0] != '-' && !strings.HasPrefix(note, "  ") {
				// Section header
				detailsParts = append(detailsParts, "")
				detailsParts = append(detailsParts, lipgloss.NewStyle().
					Foreground(secondaryColor).
					Bold(true).
					Render(note))
			} else {
				// Regular note
				detailsParts = append(detailsParts, mutedStyle.Render(note))
			}
		}

		detailsBox := boxStyle.Render(strings.Join(detailsParts, "\n"))
		fmt.Println(detailsBox)
	}

	fmt.Println()
}

// withSpinner wraps an AI call with a nice loading spinner
func withSpinner(message string, fn func(context.Context) error) error {
	ctx := context.Background()
	return showSpinner(message, func() error {
		return fn(ctx)
	})
}
