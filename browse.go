package main

import (
	"fmt"

	"github.com/tobischo/gokeepasslib/v3"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type viewMode string

var (
	viewModeList  viewMode = "list"
	viewModeEntry viewMode = "entry"
)

func browseCmd(_ *cobra.Command, _ []string) error {
	p := tea.NewProgram(initiateModel(&db.Content.Root.Groups[0]))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Encountered error: %w", err)
	}

	return nil
}

type groupCursor struct {
	group    *gokeepasslib.Group
	previous *groupCursor
}

type model struct {
	mode viewMode

	groupCursor *groupCursor

	choices []interface{} // items on the to-do list
	cursor  int           // which to-do list item our cursor is pointing at
}

func initiateModel(group *gokeepasslib.Group) model {
	return model{
		mode: viewModeList,

		groupCursor: &groupCursor{
			group: group,
		},

		// Our shopping list is a grocery list
		choices: prepareChoices(group),
	}
}

func prepareChoices(group *gokeepasslib.Group) []interface{} {
	choices := []interface{}{}
	for _, group := range group.Groups {
		choices = append(choices, group)
	}

	for _, entry := range group.Entries {
		choices = append(choices, entry)
	}

	return choices
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.mode == viewModeList {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.mode == viewModeList {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}

		case "left", "h":
			if m.mode == viewModeList {
				if m.groupCursor.previous != nil {
					m.groupCursor = m.groupCursor.previous
					m.choices = prepareChoices(m.groupCursor.group)
					m.cursor = 0
				}
			} else {
				m.mode = viewModeList
			}
		case "right", "l":
			choice := m.choices[m.cursor]
			switch t := choice.(type) {
			case gokeepasslib.Group:
				// go further to the right
				m.groupCursor = &groupCursor{
					group:    &t,
					previous: m.groupCursor,
				}
				m.choices = prepareChoices(&t)
				m.cursor = 0
			case gokeepasslib.Entry:
				m.mode = viewModeEntry
			}
		case "enter", " ":
			if m.mode == viewModeEntry {
				entry, ok := m.choices[m.cursor].(gokeepasslib.Entry)
				if !ok {
					// Do nothing for now
					break
				}

				// Ignore error here for now
				_ = clipboard.WriteAll(entry.GetPassword())
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := ""

	switch m.mode {
	case viewModeList:
		s += m.viewList()
	case viewModeEntry:
		s += m.viewEntry()
	}

	// The footer
	s += "\nPress q to quit.\n"

	return s
}

func (m model) viewList() string {
	// The header
	s := "Pick an entry\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		switch t := choice.(type) {
		case gokeepasslib.Group:
			s += fmt.Sprintf("%s %s [>]\n", cursor, t.Name)
		case gokeepasslib.Entry:
			s += fmt.Sprintf("%s %s\n", cursor, t.GetTitle())
		default:
			s += fmt.Sprintf("unsupported type %T\n", t)
		}
	}

	// Send the UI for rendering
	return s
}

func (m model) viewEntry() string {
	entry, ok := m.choices[m.cursor].(gokeepasslib.Entry)
	if !ok {
		return "fatal: choice is not entry"
	}

	s := "Entry\n\n"
	s += fmt.Sprintf("Title:             %s\n", entry.GetTitle())
	s += fmt.Sprintf("Creation:          %s\n", entry.Times.CreationTime.Time.Format(timeFormat))
	s += fmt.Sprintf(
		"Last Modification: %s\n",
		entry.Times.LastModificationTime.Time.Format(timeFormat),
	)
	s += fmt.Sprintf("Last Access:       %s\n", entry.Times.LastAccessTime.Time.Format(timeFormat))
	s += fmt.Sprintf("UserName:          %s\n", entry.GetContent("UserName"))
	s += fmt.Sprintf("URL:               %s\n", entry.GetContent("URL"))

	s += fmt.Sprintf("Notes:\n%s\n", entry.GetContent("Notes"))

	return s
}
