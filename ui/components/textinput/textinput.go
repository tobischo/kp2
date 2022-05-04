package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textinput.Model
	enterCommand tea.Cmd
}

type Option func(model *Model)

func AsPasswordField(model *Model) {
	model.Model.EchoMode = textinput.EchoPassword
	model.Model.EchoCharacter = '•'
}

func WithPrompt(prompt string) Option {
	return func(model *Model) {
		model.Model.Prompt = prompt
	}
}

func WithEnterCommand(cmd tea.Cmd) Option {
	return func(model *Model) {
		model.enterCommand = cmd
	}
}

func New(options ...Option) Model {
	ti := textinput.New()

	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle

	model := Model{
		Model: ti,
	}

	for _, option := range options {
		option(&model)
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Model.Focused() {
				return m, m.enterCommand
			}
		}
	}

	model, cmd := m.Model.Update(msg)
	m.Model = model

	return m, cmd
}
