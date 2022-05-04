package password

import (
  "fmt"

  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "github.com/tobischo/kp2/ui/commands"
  "github.com/tobischo/kp2/ui/components/button"
  "github.com/tobischo/kp2/ui/components/textinput"
  "github.com/tobischo/kp2/ui/theme"
)

type tickMsg struct{}

type Model struct {
  passwordField textinput.Model

  openButton   button.Model
  cancelButton button.Model

  err error
}

func NewModel() *Model {
  ti := textinput.New(
    textinput.AsPasswordField,
    textinput.WithPrompt(""),
    textinput.WithEnterCommand(commands.OpenDatabase),
  )
  ti.Focus()
  ti.CharLimit = 156
  ti.Width = 30

  return &Model{
    passwordField: ti,
    openButton: button.New(
      "Open",
      10,
      commands.OpenDatabase,
    ),
    cancelButton: button.New(
      "Cancel",
      10,
      tea.Quit,
    ),
  }
}

func (m *Model) Init() tea.Cmd {
  return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmds []tea.Cmd
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyCtrlC, tea.KeyEsc:
      cmds = append(cmds, tea.Quit)
    case tea.KeyTab:
      switch {
      case m.passwordField.Focused():
        m.passwordField.Blur()
        m.openButton.Focus()
      case m.openButton.Focused():
        m.openButton.Blur()
        m.cancelButton.Focus()
      case m.cancelButton.Focused():
        m.cancelButton.Blur()
        m.passwordField.Focus()
      }
    }
  case commands.OpenDatabaseMsg:
    cmds = append(cmds, commands.PassDatabasePassword(m.passwordField.Value()))
  // We handle errors just like any other message
  case commands.ErrorMsg:
    m.err = msg
  }

  m.passwordField, cmd = m.passwordField.Update(msg)
  cmds = append(cmds, cmd)
  _, cmd = m.openButton.Update(msg)
  cmds = append(cmds, cmd)
  _, cmd = m.cancelButton.Update(msg)
  cmds = append(cmds, cmd)
  return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
  viewContent := fmt.Sprintf(
    "Password:\n\n%s\n\n%s",
    m.passwordField.View(),
    lipgloss.JoinHorizontal(
      lipgloss.Center,
      m.openButton.View(),
      m.cancelButton.View(),
    ),
  ) + "\n"

  // TODO find a better place outside of the border
  if m.err != nil {
    viewContent = fmt.Sprintf("%s\n%s\n", viewContent, m.err.Error())
  }

  return lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(theme.DefaultTheme.PrimaryBorder).
    Render(
      lipgloss.Place(
        40,
        10,
        lipgloss.Center,
        lipgloss.Center,
        viewContent,
      ),
    )
}
