package main

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("62")  // Purple
	secondaryColor = lipgloss.Color("39")  // Blue
	groupColor     = lipgloss.Color("214") // Orange
	entryColor     = lipgloss.Color("252") // Light gray
	mutedColor     = lipgloss.Color("241") // Dark gray
	successColor   = lipgloss.Color("82")  // Green
	warningColor   = lipgloss.Color("214") // Yellow/orange
	errorColor     = lipgloss.Color("196") // Red

	// App title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Header style for current path/location
	headerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Footer/status bar style
	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Status message style (for copy feedback)
	statusStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Error status style
	errorStatusStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	// Selected item in list
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Unselected item in list
	unselectedStyle = lipgloss.NewStyle().
			Foreground(entryColor)

	// Group item style (folders)
	groupStyle = lipgloss.NewStyle().
			Foreground(groupColor)

	// Entry detail labels
	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(20)

	// Entry detail values
	valueStyle = lipgloss.NewStyle().
			Foreground(entryColor)

	// Help key style
	helpKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Help description style
	helpDescStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Search prompt style
	searchPromptStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	// Form label style
	formLabelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Width(12)

	// Active form input style
	formActiveInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1)

	// Inactive form input style
	formInactiveInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Padding(0, 1)

	// Confirmation dialog style
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			Align(lipgloss.Center)

	// Muted text style
	mutedStyle = lipgloss.NewStyle().Foreground(mutedColor)

	// Cursor styles for list
	cursorStyle     = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	cursorTextStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
)

// List delegate styles - used to customize list.Model rendering
type listStyles struct {
	NormalTitle   lipgloss.Style
	NormalDesc    lipgloss.Style
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
}

func defaultListStyles() listStyles {
	return listStyles{
		NormalTitle:   unselectedStyle,
		NormalDesc:    lipgloss.NewStyle().Foreground(mutedColor),
		SelectedTitle: selectedStyle,
		SelectedDesc:  lipgloss.NewStyle().Foreground(secondaryColor),
	}
}
