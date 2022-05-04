package ui

import (
  "fmt"
  "os"
  "strings"

  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "github.com/tobischo/gokeepasslib/v3"
  "github.com/tobischo/kp2/config"
  "github.com/tobischo/kp2/ui/commands"
  "github.com/tobischo/kp2/ui/components/help"
  "github.com/tobischo/kp2/ui/keys"
  "github.com/tobischo/kp2/ui/viewctx"
  "github.com/tobischo/kp2/ui/views/password"
)

type Model struct {
  viewContext viewctx.ViewContext

  keys keys.KeyMap
  err  error
  // sidebar       sidebar.Model
  viewMode int

  views map[int]tea.Model

  help help.Model

  db *gokeepasslib.Database
  // prs           []section.Section
  // issues        []section.Section
  // ready         bool
  // isSidebarOpen bool
}

func NewModel() *Model {
  // tabsModel := tabs.NewModel()
  return &Model{
    viewContext: viewctx.NewViewContext(),
    keys:        keys.Keys,
    help:        help.NewModel(),
    views: map[int]tea.Model{
      0: password.NewModel(),
    },
    viewMode: 0,
    // sidebar:       sidebar.NewModel(),
  }
}

func initScreen() tea.Msg {
  config, err := config.LoadConfig()
  if err != nil {
    return commands.SetError(err)
  }

  return configMsg{Config: config}
}

func (m *Model) Init() tea.Cmd {
  return tea.Batch(initScreen, tea.EnterAltScreen)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmds []tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.Quit):
      cmds = append(cmds, tea.Quit)
    }
  case tea.WindowSizeMsg:
    m.updateWindowSize(msg)
  case commands.DatabasePasswordMsg:
    err := m.openDatabase(msg)

    cmds = append(cmds, commands.SetError(err))
  }

  _, cmd := m.views[m.viewMode].Update(msg)

  cmds = append(cmds, cmd)

  return m, tea.Batch(cmds...)
}

func (m *Model) updateWindowSize(msg tea.WindowSizeMsg) {
  m.viewContext.ScreenHeight = msg.Height
  m.viewContext.ScreenWidth = msg.Width

  m.help.SetWidth(m.viewContext.ScreenWidth)
}

func (m *Model) openDatabase(msg commands.DatabasePasswordMsg) error {
  filePath := os.Getenv("KP2FILE")
  credentials := gokeepasslib.NewPasswordCredentials(
    msg.Password(),
  )

  m.db = &gokeepasslib.Database{}
  m.db.Credentials = credentials

  file, err := os.Open(filePath)
  if err != nil {
    return fmt.Errorf("Failed to open Keepass2 file %s: '%s'", filePath, err)
  }

  err = gokeepasslib.NewDecoder(file).Decode(m.db)
  if err != nil {
    return fmt.Errorf("Failed to decode Keepass2 file: %s", err)
  }

  if err := m.db.UnlockProtectedEntries(); err != nil {
    return err
  }

  return nil
}

func (m *Model) View() string {
  s := strings.Builder{}

  centered := lipgloss.Place(
    m.viewContext.ScreenWidth,
    m.viewContext.ScreenHeight-help.FooterHeight,
    lipgloss.Center,
    lipgloss.Center,
    m.views[m.viewMode].View(),
  )

  s.WriteString(centered)

  s.WriteString("\n")

  s.WriteString(m.help.View(m.viewContext))
  return s.String()
}

type configMsg struct {
  Config config.Config
}
