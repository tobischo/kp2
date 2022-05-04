package theme

import "github.com/charmbracelet/lipgloss"

var (
	indigo       = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#383B5B"}
	subtleIndigo = lipgloss.AdaptiveColor{Light: "#5A57B5", Dark: "#242347"}
)

type Theme struct {
	PrimaryText           lipgloss.AdaptiveColor
	SecondaryText         lipgloss.AdaptiveColor
	PrimaryBorder         lipgloss.AdaptiveColor
	SecondaryBorder       lipgloss.AdaptiveColor
	ButtonBackgroundFocus lipgloss.AdaptiveColor
	ButtonBackground      lipgloss.AdaptiveColor
}

var DefaultTheme = Theme{
	PrimaryText:           lipgloss.AdaptiveColor{Light: "#242347", Dark: "#E2E1ED"},
	SecondaryText:         lipgloss.AdaptiveColor{Light: subtleIndigo.Dark, Dark: subtleIndigo.Light},
	PrimaryBorder:         indigo,
	SecondaryBorder:       lipgloss.AdaptiveColor{Light: indigo.Light, Dark: "#39386B"},
	ButtonBackgroundFocus: lipgloss.AdaptiveColor{Light: indigo.Dark, Dark: indigo.Light},
	ButtonBackground:      lipgloss.AdaptiveColor{Light: indigo.Light, Dark: indigo.Dark},
}
