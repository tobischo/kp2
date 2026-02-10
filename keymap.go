package main

import "github.com/charmbracelet/bubbles/key"

// keyMap defines all keybindings for the TUI
type keyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Filter key.Binding

	// Actions
	CopyPassword key.Binding
	CopyUsername key.Binding
	CopyURL      key.Binding
	AddEntry     key.Binding
	EditEntry    key.Binding
	DeleteEntry  key.Binding
	CreateGroup  key.Binding

	// General
	Help   key.Binding
	Quit   key.Binding
	Cancel key.Binding
	Submit key.Binding

	// Form navigation
	NextField key.Binding
	PrevField key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "back"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "enter"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		CopyPassword: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy password"),
		),
		CopyUsername: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "copy username"),
		),
		CopyURL: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "copy URL"),
		),
		AddEntry: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add entry"),
		),
		EditEntry: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit entry"),
		),
		DeleteEntry: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		CreateGroup: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "new group"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
		NextField: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		PrevField: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev field"),
		),
	}
}

// ShortHelp returns keybindings for the short help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Filter, k.Help, k.Quit}
}

// FullHelp returns keybindings for the full help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Back},
		{k.CopyPassword, k.CopyUsername, k.CopyURL},
		{k.AddEntry, k.EditEntry, k.DeleteEntry, k.CreateGroup},
		{k.Filter, k.Help, k.Quit},
	}
}

// ListKeyMap returns keybindings shown in list view
func (k keyMap) ListKeyMap() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Right, k.Left, k.Filter, k.AddEntry, k.CreateGroup, k.Help, k.Quit}
}

// EntryKeyMap returns keybindings shown in entry detail view
func (k keyMap) EntryKeyMap() []key.Binding {
	return []key.Binding{k.CopyPassword, k.CopyUsername, k.CopyURL, k.EditEntry, k.DeleteEntry, k.Back, k.Help, k.Quit}
}

// SearchKeyMap returns keybindings shown in search view
func (k keyMap) SearchKeyMap() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Cancel}
}

// FormKeyMap returns keybindings shown in form view
func (k keyMap) FormKeyMap() []key.Binding {
	return []key.Binding{k.NextField, k.PrevField, k.Submit, k.Cancel}
}

// ConfirmKeyMap returns keybindings shown in confirm view
func (k keyMap) ConfirmKeyMap() []key.Binding {
	return []key.Binding{k.Submit, k.Cancel}
}
