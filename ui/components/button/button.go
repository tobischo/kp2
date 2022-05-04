package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tobischo/kp2/ui/theme"
)

type Model struct {
	label    string
	focus    bool
	command  tea.Cmd
	minWidth int
}

func New(label string, minWidth int, command tea.Cmd) Model {
	return Model{
		label:    label,
		command:  command,
		minWidth: minWidth,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch()
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Focused() bool {
	return m.focus
}

func (m *Model) Blur() {
	m.focus = false
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focus {
				return m, m.command
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	foregroundColor := theme.DefaultTheme.PrimaryText
	backgroundColor := theme.DefaultTheme.ButtonBackgroundFocus

	if !m.focus {
		foregroundColor = theme.DefaultTheme.SecondaryText
		backgroundColor = theme.DefaultTheme.ButtonBackground
	}

	return buttonStyle.Copy().
		Foreground(foregroundColor).
		Background(backgroundColor).
		Width(m.minWidth).
		Align(lipgloss.Center).
		Render(m.label)
}
