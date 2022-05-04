package help

import (
  bubblesHelp "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/tobischo/kp2/ui/keys"
  "github.com/tobischo/kp2/ui/viewctx"
)

type Model struct {
  help bubblesHelp.Model
}

func NewModel() Model {
  help := bubblesHelp.NewModel()
  help.Styles = bubblesHelp.Styles{
    ShortDesc:      helpTextStyle.Copy(),
    FullDesc:       helpTextStyle.Copy(),
    ShortSeparator: helpTextStyle.Copy(),
    FullSeparator:  helpTextStyle.Copy(),
    FullKey:        helpTextStyle.Copy(),
    ShortKey:       helpTextStyle.Copy(),
    Ellipsis:       helpTextStyle.Copy(),
  }

  return Model{
    help: help,
  }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, keys.Keys.Help):
      m.help.ShowAll = !m.help.ShowAll
    }
  }

  return m, nil
}

func (m Model) View(ctx viewctx.ViewContext) string {
  return helpStyle.Copy().
    Width(ctx.ScreenWidth).
    Render(m.help.View(keys.Keys))
}

func (m *Model) SetWidth(width int) {
  m.help.Width = width
}
