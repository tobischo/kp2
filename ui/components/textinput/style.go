package textinput

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tobischo/kp2/ui/theme"
)

var (
	focusedStyle = lipgloss.NewStyle().
		Foreground(theme.DefaultTheme.PrimaryText).
		Background(theme.DefaultTheme.SecondaryText)
)
