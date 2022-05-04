package help

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tobischo/kp2/ui/theme"
)

var (
	FooterHeight = 3

	helpTextStyle = lipgloss.NewStyle().Foreground(theme.DefaultTheme.SecondaryText)
	helpStyle     = lipgloss.NewStyle().
			Height(FooterHeight - 1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.DefaultTheme.PrimaryBorder)
)
