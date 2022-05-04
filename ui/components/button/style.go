package button

import (
  "github.com/charmbracelet/lipgloss"
  "github.com/tobischo/kp2/ui/theme"
)

var (
  paddingLeftRight = 2

  buttonStyle = lipgloss.NewStyle().
      PaddingLeft(paddingLeftRight).
      PaddingRight(paddingLeftRight).
      BorderStyle(lipgloss.NormalBorder()).
      BorderForeground(theme.DefaultTheme.PrimaryBorder)
)
