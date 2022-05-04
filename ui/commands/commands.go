package commands

import tea "github.com/charmbracelet/bubbletea"

type OpenDatabaseMsg struct{}

func OpenDatabase() tea.Msg {
	return OpenDatabaseMsg{}
}

type DatabasePasswordMsg struct {
	password string
}

func (msg DatabasePasswordMsg) Password() string {
	return msg.password
}

func PassDatabasePassword(password string) tea.Cmd {
	return func() tea.Msg {
		return DatabasePasswordMsg{
			password: password,
		}
	}
}

type ErrorMsg error

func SetError(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}
