package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

// View modes for the TUI
type viewMode int

const (
	viewModeList viewMode = iota
	viewModeEntry
	viewModeSearch
	viewModeAddForm
	viewModeEditForm
	viewModeConfirm
	viewModeCreateGroup
)

// Form field indices
const (
	fieldTitle = iota
	fieldUsername
	fieldPassword
	fieldURL
	fieldNotes
)

// groupCursor tracks navigation through the group hierarchy
type groupCursor struct {
	group    *gokeepasslib.Group
	previous *groupCursor
}

// listItem represents an item in the list (group or entry)
type listItem struct {
	isGroup bool
	group   gokeepasslib.Group
	entry   gokeepasslib.Entry
}

func (i listItem) Title() string {
	if i.isGroup {
		return i.group.Name
	}
	return i.entry.GetTitle()
}

func (i listItem) Description() string {
	if i.isGroup {
		return fmt.Sprintf("%d entries, %d groups", len(i.group.Entries), len(i.group.Groups))
	}
	return i.entry.GetContent("UserName")
}

// model is the main TUI model
type model struct {
	mode        viewMode
	groupCursor *groupCursor

	// List state
	items         []listItem
	filteredItems []listItem
	cursor        int

	// Search state
	searchInput textinput.Model
	searching   bool

	// Entry detail view
	viewport      viewport.Model
	selectedEntry *gokeepasslib.Entry

	// Form state (for add/edit)
	formInputs        []textinput.Model
	formFocus         int
	editingEntry      *gokeepasslib.Entry
	editingEntryIndex int

	// Group creation form
	groupInput textinput.Model

	// Confirm dialog state
	confirmAction string
	confirmTarget string

	// UI components
	help   help.Model
	keys   keyMap
	width  int
	height int

	// Status message
	statusMsg     string
	statusIsError bool
}

// Status message command
type statusMsg struct {
	message string
	isError bool
}

func clearStatusCmd() tea.Msg {
	return statusMsg{"", false}
}

func browseCmd(_ *cobra.Command, _ []string) error {
	p := tea.NewProgram(initialModel(&db.Content.Root.Groups[0]), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("encountered error: %w", err)
	}
	return nil
}

func initialModel(group *gokeepasslib.Group) model {
	// Initialize search input
	si := textinput.New()
	si.Placeholder = "Search entries..."
	si.CharLimit = 100
	si.Width = 40
	si.PromptStyle = searchPromptStyle
	si.Prompt = "/ "

	// Initialize group input
	gi := textinput.New()
	gi.Placeholder = "Group name"
	gi.CharLimit = 100
	gi.Width = 40
	gi.PromptStyle = formLabelStyle
	gi.Prompt = "Name: "

	// Initialize help
	h := help.New()
	h.Styles.ShortKey = helpKeyStyle
	h.Styles.ShortDesc = helpDescStyle
	h.Styles.FullKey = helpKeyStyle
	h.Styles.FullDesc = helpDescStyle

	items := prepareItems(group)

	return model{
		mode: viewModeList,
		groupCursor: &groupCursor{
			group: group,
		},
		items:         items,
		filteredItems: items,
		searchInput:   si,
		groupInput:    gi,
		help:          h,
		keys:          DefaultKeyMap(),
		width:         80,
		height:        24,
	}
}

func prepareItems(group *gokeepasslib.Group) []listItem {
	items := make([]listItem, 0, len(group.Groups)+len(group.Entries))
	for _, g := range group.Groups {
		items = append(items, listItem{isGroup: true, group: g})
	}
	for _, e := range group.Entries {
		items = append(items, listItem{isGroup: false, entry: e})
	}
	return items
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 6

	case statusMsg:
		m.statusMsg = msg.message
		m.statusIsError = msg.isError

	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, m.keys.Quit) && m.mode == viewModeList && !m.searching {
			return m, tea.Quit
		}

		switch m.mode {
		case viewModeList:
			return m.updateList(msg)
		case viewModeEntry:
			return m.updateEntry(msg)
		case viewModeSearch:
			return m.updateSearch(msg)
		case viewModeAddForm:
			return m.updateAddForm(msg)
		case viewModeEditForm:
			return m.updateEditForm(msg)
		case viewModeConfirm:
			return m.updateConfirm(msg)
		case viewModeCreateGroup:
			return m.updateCreateGroup(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.filteredItems)-1 {
			m.cursor++
		}

	case key.Matches(msg, m.keys.Left):
		if m.groupCursor.previous != nil {
			m.groupCursor = m.groupCursor.previous
			m.items = prepareItems(m.groupCursor.group)
			m.filteredItems = m.items
			m.cursor = 0
		}

	case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Enter):
		if len(m.filteredItems) > 0 {
			item := m.filteredItems[m.cursor]
			if item.isGroup {
				m.groupCursor = &groupCursor{
					group:    &item.group,
					previous: m.groupCursor,
				}
				m.items = prepareItems(&item.group)
				m.filteredItems = m.items
				m.cursor = 0
			} else {
				m.selectedEntry = &item.entry
				m.mode = viewModeEntry
				m.viewport = viewport.New(m.width, m.height-6)
				m.viewport.SetContent(m.renderEntryContent())
			}
		}

	case key.Matches(msg, m.keys.Filter):
		m.mode = viewModeSearch
		m.searchInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.AddEntry):
		m.mode = viewModeAddForm
		m.initFormInputs()
		m.formFocus = 0
		m.formInputs[0].Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.CreateGroup):
		m.mode = viewModeCreateGroup
		m.groupInput.SetValue("")
		m.groupInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.DeleteEntry):
		if len(m.filteredItems) > 0 && !m.filteredItems[m.cursor].isGroup {
			item := m.filteredItems[m.cursor]
			m.confirmAction = "delete"
			m.confirmTarget = item.entry.GetTitle()
			m.mode = viewModeConfirm
		}

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll

	case key.Matches(msg, m.keys.CopyPassword):
		if len(m.filteredItems) > 0 && !m.filteredItems[m.cursor].isGroup {
			item := m.filteredItems[m.cursor]
			if err := clipboard.WriteAll(item.entry.GetPassword()); err != nil {
				m.statusMsg = "Failed to copy password"
				m.statusIsError = true
			} else {
				m.statusMsg = "Password copied!"
				m.statusIsError = false
			}
		}

	case key.Matches(msg, m.keys.CopyUsername):
		if len(m.filteredItems) > 0 && !m.filteredItems[m.cursor].isGroup {
			item := m.filteredItems[m.cursor]
			if err := clipboard.WriteAll(item.entry.GetContent("UserName")); err != nil {
				m.statusMsg = "Failed to copy username"
				m.statusIsError = true
			} else {
				m.statusMsg = "Username copied!"
				m.statusIsError = false
			}
		}

	case key.Matches(msg, m.keys.CopyURL):
		if len(m.filteredItems) > 0 && !m.filteredItems[m.cursor].isGroup {
			item := m.filteredItems[m.cursor]
			if err := clipboard.WriteAll(item.entry.GetContent("URL")); err != nil {
				m.statusMsg = "Failed to copy URL"
				m.statusIsError = true
			} else {
				m.statusMsg = "URL copied!"
				m.statusIsError = false
			}
		}
	}

	return m, nil
}

func (m model) updateEntry(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Left):
		m.mode = viewModeList
		m.selectedEntry = nil

	case key.Matches(msg, m.keys.CopyPassword):
		if m.selectedEntry != nil {
			if err := clipboard.WriteAll(m.selectedEntry.GetPassword()); err != nil {
				m.statusMsg = "Failed to copy password"
				m.statusIsError = true
			} else {
				m.statusMsg = "Password copied!"
				m.statusIsError = false
			}
		}

	case key.Matches(msg, m.keys.CopyUsername):
		if m.selectedEntry != nil {
			if err := clipboard.WriteAll(m.selectedEntry.GetContent("UserName")); err != nil {
				m.statusMsg = "Failed to copy username"
				m.statusIsError = true
			} else {
				m.statusMsg = "Username copied!"
				m.statusIsError = false
			}
		}

	case key.Matches(msg, m.keys.CopyURL):
		if m.selectedEntry != nil {
			if err := clipboard.WriteAll(m.selectedEntry.GetContent("URL")); err != nil {
				m.statusMsg = "Failed to copy URL"
				m.statusIsError = true
			} else {
				m.statusMsg = "URL copied!"
				m.statusIsError = false
			}
		}

	case key.Matches(msg, m.keys.EditEntry):
		if m.selectedEntry != nil {
			m.mode = viewModeEditForm
			m.editingEntry = m.selectedEntry
			// Find the index in the current group's entries
			for i, e := range m.groupCursor.group.Entries {
				if e.GetTitle() == m.selectedEntry.GetTitle() {
					m.editingEntryIndex = i
					break
				}
			}
			m.initFormInputsWithEntry(m.selectedEntry)
			m.formFocus = 0
			m.formInputs[0].Focus()
			return m, textinput.Blink
		}

	case key.Matches(msg, m.keys.DeleteEntry):
		if m.selectedEntry != nil {
			m.confirmAction = "delete"
			m.confirmTarget = m.selectedEntry.GetTitle()
			m.mode = viewModeConfirm
		}

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll

	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = viewModeList
		m.searchInput.SetValue("")
		m.filteredItems = m.items
		m.cursor = 0
		m.searchInput.Blur()

	case key.Matches(msg, m.keys.Enter):
		if len(m.filteredItems) > 0 {
			item := m.filteredItems[m.cursor]
			if item.isGroup {
				m.groupCursor = &groupCursor{
					group:    &item.group,
					previous: m.groupCursor,
				}
				m.items = prepareItems(&item.group)
				m.filteredItems = m.items
				m.cursor = 0
				m.mode = viewModeList
				m.searchInput.SetValue("")
				m.searchInput.Blur()
			} else {
				m.selectedEntry = &item.entry
				m.mode = viewModeEntry
				m.viewport = viewport.New(m.width, m.height-6)
				m.viewport.SetContent(m.renderEntryContent())
				m.searchInput.SetValue("")
				m.searchInput.Blur()
			}
		}

	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.filteredItems)-1 {
			m.cursor++
		}

	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		// Filter items based on search
		m.filterItems()
		return m, cmd
	}

	return m, nil
}

func (m *model) filterItems() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filteredItems = m.items
		m.cursor = 0
		return
	}

	filtered := make([]listItem, 0)
	for _, item := range m.items {
		title := strings.ToLower(item.Title())
		desc := strings.ToLower(item.Description())
		if strings.Contains(title, query) || strings.Contains(desc, query) {
			filtered = append(filtered, item)
		}
	}
	m.filteredItems = filtered
	if m.cursor >= len(m.filteredItems) {
		m.cursor = max(0, len(m.filteredItems)-1)
	}
}

func (m model) updateAddForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = viewModeList
		m.blurAllFormInputs()

	case key.Matches(msg, m.keys.NextField):
		m.formInputs[m.formFocus].Blur()
		m.formFocus = (m.formFocus + 1) % len(m.formInputs)
		m.formInputs[m.formFocus].Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.PrevField):
		m.formInputs[m.formFocus].Blur()
		m.formFocus = (m.formFocus - 1 + len(m.formInputs)) % len(m.formInputs)
		m.formInputs[m.formFocus].Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.Enter):
		// Submit form - create new entry
		entry := m.createEntryFromForm()
		m.groupCursor.group.Entries = append(m.groupCursor.group.Entries, entry)
		changed = true
		m.items = prepareItems(m.groupCursor.group)
		m.filteredItems = m.items
		m.mode = viewModeList
		m.blurAllFormInputs()
		m.statusMsg = "Entry created!"
		m.statusIsError = false

	default:
		var cmd tea.Cmd
		m.formInputs[m.formFocus], cmd = m.formInputs[m.formFocus].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) updateEditForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = viewModeEntry
		m.blurAllFormInputs()

	case key.Matches(msg, m.keys.NextField):
		m.formInputs[m.formFocus].Blur()
		m.formFocus = (m.formFocus + 1) % len(m.formInputs)
		m.formInputs[m.formFocus].Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.PrevField):
		m.formInputs[m.formFocus].Blur()
		m.formFocus = (m.formFocus - 1 + len(m.formInputs)) % len(m.formInputs)
		m.formInputs[m.formFocus].Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.Enter):
		// Submit form - update entry
		m.updateEntryFromForm()
		changed = true
		m.items = prepareItems(m.groupCursor.group)
		m.filteredItems = m.items
		// Update viewport content
		m.viewport.SetContent(m.renderEntryContent())
		m.mode = viewModeEntry
		m.blurAllFormInputs()
		m.statusMsg = "Entry updated!"
		m.statusIsError = false

	default:
		var cmd tea.Cmd
		m.formInputs[m.formFocus], cmd = m.formInputs[m.formFocus].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if m.confirmAction == "delete" {
			// Delete the entry
			m.deleteCurrentEntry()
			m.items = prepareItems(m.groupCursor.group)
			m.filteredItems = m.items
			if m.cursor >= len(m.filteredItems) {
				m.cursor = max(0, len(m.filteredItems)-1)
			}
			m.mode = viewModeList
			m.selectedEntry = nil
			m.statusMsg = "Entry deleted!"
			m.statusIsError = false
			changed = true
		}

	case "n", "N", "esc":
		if m.selectedEntry != nil {
			m.mode = viewModeEntry
		} else {
			m.mode = viewModeList
		}
	}

	return m, nil
}

func (m model) updateCreateGroup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = viewModeList
		m.groupInput.Blur()

	case key.Matches(msg, m.keys.Enter):
		groupName := m.groupInput.Value()
		if groupName != "" {
			newGroup := gokeepasslib.NewGroup()
			newGroup.Name = groupName
			m.groupCursor.group.Groups = append(m.groupCursor.group.Groups, newGroup)
			changed = true
			m.items = prepareItems(m.groupCursor.group)
			m.filteredItems = m.items
			m.statusMsg = "Group created!"
			m.statusIsError = false
		}
		m.mode = viewModeList
		m.groupInput.Blur()

	default:
		var cmd tea.Cmd
		m.groupInput, cmd = m.groupInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) initFormInputs() {
	m.formInputs = make([]textinput.Model, 5)

	labels := []string{"Title", "Username", "Password", "URL", "Notes"}
	placeholders := []string{"Entry title", "username@example.com", "••••••••", "https://example.com", "Additional notes"}

	for i := range m.formInputs {
		ti := textinput.New()
		ti.Placeholder = placeholders[i]
		ti.CharLimit = 256
		ti.Width = 40
		ti.Prompt = labels[i] + ": "
		ti.PromptStyle = formLabelStyle

		if i == fieldPassword {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}

		m.formInputs[i] = ti
	}
}

func (m *model) initFormInputsWithEntry(entry *gokeepasslib.Entry) {
	m.initFormInputs()
	m.formInputs[fieldTitle].SetValue(entry.GetTitle())
	m.formInputs[fieldUsername].SetValue(entry.GetContent("UserName"))
	m.formInputs[fieldPassword].SetValue(entry.GetPassword())
	m.formInputs[fieldURL].SetValue(entry.GetContent("URL"))
	m.formInputs[fieldNotes].SetValue(entry.GetContent("Notes"))
}

func (m *model) blurAllFormInputs() {
	for i := range m.formInputs {
		m.formInputs[i].Blur()
	}
}

func (m model) createEntryFromForm() gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values,
		gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: m.formInputs[fieldTitle].Value()}},
		gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: m.formInputs[fieldUsername].Value()}},
		gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: m.formInputs[fieldPassword].Value(), Protected: wrappers.NewBoolWrapper(true)}},
		gokeepasslib.ValueData{Key: "URL", Value: gokeepasslib.V{Content: m.formInputs[fieldURL].Value()}},
		gokeepasslib.ValueData{Key: "Notes", Value: gokeepasslib.V{Content: m.formInputs[fieldNotes].Value()}},
	)
	return entry
}

func (m *model) updateEntryFromForm() {
	if m.editingEntry == nil || m.editingEntryIndex < 0 {
		return
	}

	entry := &m.groupCursor.group.Entries[m.editingEntryIndex]

	// Update values
	for i := range entry.Values {
		switch entry.Values[i].Key {
		case "Title":
			entry.Values[i].Value.Content = m.formInputs[fieldTitle].Value()
		case "UserName":
			entry.Values[i].Value.Content = m.formInputs[fieldUsername].Value()
		case "Password":
			entry.Values[i].Value.Content = m.formInputs[fieldPassword].Value()
		case "URL":
			entry.Values[i].Value.Content = m.formInputs[fieldURL].Value()
		case "Notes":
			entry.Values[i].Value.Content = m.formInputs[fieldNotes].Value()
		}
	}

	// Update selected entry pointer
	m.selectedEntry = entry
}

func (m *model) deleteCurrentEntry() {
	if m.selectedEntry == nil {
		// Try to delete from list view
		if len(m.filteredItems) > 0 && !m.filteredItems[m.cursor].isGroup {
			title := m.filteredItems[m.cursor].entry.GetTitle()
			for i, e := range m.groupCursor.group.Entries {
				if e.GetTitle() == title {
					m.groupCursor.group.Entries = append(
						m.groupCursor.group.Entries[:i],
						m.groupCursor.group.Entries[i+1:]...,
					)
					return
				}
			}
		}
		return
	}

	// Delete selected entry
	for i, e := range m.groupCursor.group.Entries {
		if e.GetTitle() == m.selectedEntry.GetTitle() {
			m.groupCursor.group.Entries = append(
				m.groupCursor.group.Entries[:i],
				m.groupCursor.group.Entries[i+1:]...,
			)
			return
		}
	}
}

func (m model) View() string {
	var s strings.Builder

	switch m.mode {
	case viewModeList:
		s.WriteString(m.viewList())
	case viewModeEntry:
		s.WriteString(m.viewEntry())
	case viewModeSearch:
		s.WriteString(m.viewSearch())
	case viewModeAddForm:
		s.WriteString(m.viewForm("Add Entry"))
	case viewModeEditForm:
		s.WriteString(m.viewForm("Edit Entry"))
	case viewModeConfirm:
		s.WriteString(m.viewConfirm())
	case viewModeCreateGroup:
		s.WriteString(m.viewCreateGroup())
	}

	// Footer with help
	s.WriteString("\n")
	if m.statusMsg != "" {
		if m.statusIsError {
			s.WriteString(errorStatusStyle.Render(m.statusMsg))
		} else {
			s.WriteString(statusStyle.Render(m.statusMsg))
		}
		s.WriteString("\n")
	}
	s.WriteString(m.renderHelp())

	return s.String()
}

// visibleRange computes the start and end indices for a scrollable list
// window that keeps the cursor visible. headerLines is the number of lines
// used above the list (title, path, blank lines, etc.).
func (m model) visibleRange(itemCount, headerLines int) (start, end int) {
	// footer: 1 blank + optional status + help ≈ 3 lines
	footerLines := 3
	if m.statusMsg != "" {
		footerLines++
	}

	available := m.height - headerLines - footerLines
	if available < 1 {
		available = 1
	}
	if available >= itemCount {
		return 0, itemCount
	}

	// Keep cursor roughly centred, but clamp to list bounds.
	half := available / 2
	start = m.cursor - half
	if start < 0 {
		start = 0
	}
	end = start + available
	if end > itemCount {
		end = itemCount
		start = end - available
	}
	return start, end
}

func (m model) viewList() string {
	var s strings.Builder

	// Header with path
	path := m.getPath()
	s.WriteString(titleStyle.Render("KP2 Browser"))
	s.WriteString("\n")
	s.WriteString(headerStyle.Render("📁 " + path))
	s.WriteString("\n\n")

	if len(m.filteredItems) == 0 {
		s.WriteString(mutedStyle.Render("  (empty)"))
		s.WriteString("\n")
	} else {
		// 3 header lines: title, path, blank line
		start, end := m.visibleRange(len(m.filteredItems), 3)

		if start > 0 {
			s.WriteString(mutedStyle.Render(fmt.Sprintf("  ↑ %d more", start)))
			s.WriteString("\n")
		}

		for i := start; i < end; i++ {
			item := m.filteredItems[i]
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("> ")
			}

			var line string
			if item.isGroup {
				icon := "📁 "
				name := item.group.Name
				if m.cursor == i {
					line = cursorTextStyle.Render(icon + name)
				} else {
					line = groupStyle.Render(icon + name)
				}
			} else {
				icon := "🔑 "
				name := item.entry.GetTitle()
				user := item.entry.GetContent("UserName")
				if m.cursor == i {
					line = cursorTextStyle.Render(icon+name) + " " + lipgloss.NewStyle().Foreground(secondaryColor).Render(user)
				} else {
					line = unselectedStyle.Render(icon+name) + " " + lipgloss.NewStyle().Foreground(mutedColor).Render(user)
				}
			}

			s.WriteString(cursor + line + "\n")
		}

		if end < len(m.filteredItems) {
			s.WriteString(mutedStyle.Render(fmt.Sprintf("  ↓ %d more", len(m.filteredItems)-end)))
			s.WriteString("\n")
		}
	}

	return s.String()
}

func (m model) viewEntry() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Entry Details"))
	s.WriteString("\n\n")
	s.WriteString(m.viewport.View())

	return s.String()
}

func (m model) renderEntryContent() string {
	if m.selectedEntry == nil {
		return "No entry selected"
	}

	entry := m.selectedEntry
	var s strings.Builder

	s.WriteString(labelStyle.Render("Title:") + valueStyle.Render(entry.GetTitle()) + "\n")
	s.WriteString(labelStyle.Render("Username:") + valueStyle.Render(entry.GetContent("UserName")) + "\n")
	s.WriteString(labelStyle.Render("Password:") + valueStyle.Render("••••••••") + "\n")
	s.WriteString(labelStyle.Render("URL:") + valueStyle.Render(entry.GetContent("URL")) + "\n")
	s.WriteString("\n")
	s.WriteString(labelStyle.Render("Created:") + valueStyle.Render(entry.Times.CreationTime.Time.Format(timeFormat)) + "\n")
	s.WriteString(labelStyle.Render("Modified:") + valueStyle.Render(entry.Times.LastModificationTime.Time.Format(timeFormat)) + "\n")
	s.WriteString(labelStyle.Render("Accessed:") + valueStyle.Render(entry.Times.LastAccessTime.Time.Format(timeFormat)) + "\n")

	notes := entry.GetContent("Notes")
	if notes != "" {
		s.WriteString("\n")
		s.WriteString(labelStyle.Render("Notes:") + "\n")
		s.WriteString(valueStyle.Render(notes) + "\n")
	}

	return s.String()
}

func (m model) viewSearch() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Search"))
	s.WriteString("\n\n")
	s.WriteString(m.searchInput.View())
	s.WriteString("\n\n")

	if len(m.filteredItems) == 0 {
		s.WriteString(mutedStyle.Render("  No matches"))
		s.WriteString("\n")
	} else {
		// 4 header lines: title, blank, search input, blank
		start, end := m.visibleRange(len(m.filteredItems), 4)

		if start > 0 {
			s.WriteString(mutedStyle.Render(fmt.Sprintf("  ↑ %d more", start)))
			s.WriteString("\n")
		}

		for i := start; i < end; i++ {
			item := m.filteredItems[i]
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("> ")
			}

			var line string
			if item.isGroup {
				icon := "📁 "
				if m.cursor == i {
					line = cursorTextStyle.Render(icon + item.group.Name)
				} else {
					line = groupStyle.Render(icon + item.group.Name)
				}
			} else {
				icon := "🔑 "
				if m.cursor == i {
					line = cursorTextStyle.Render(icon + item.entry.GetTitle())
				} else {
					line = unselectedStyle.Render(icon + item.entry.GetTitle())
				}
			}

			s.WriteString(cursor + line + "\n")
		}

		if end < len(m.filteredItems) {
			s.WriteString(mutedStyle.Render(fmt.Sprintf("  ↓ %d more", len(m.filteredItems)-end)))
			s.WriteString("\n")
		}
	}

	return s.String()
}

func (m model) viewForm(title string) string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	for i, input := range m.formInputs {
		if i == m.formFocus {
			s.WriteString(formActiveInputStyle.Render(input.View()))
		} else {
			s.WriteString(formInactiveInputStyle.Render(input.View()))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Tab: next field • Enter: submit • Esc: cancel"))

	return s.String()
}

func (m model) viewConfirm() string {
	var s strings.Builder

	msg := fmt.Sprintf("Delete '%s'?\n\nPress y to confirm, n to cancel", m.confirmTarget)
	s.WriteString(confirmStyle.Render(msg))

	return s.String()
}

func (m model) viewCreateGroup() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Create Group"))
	s.WriteString("\n\n")
	s.WriteString(formActiveInputStyle.Render(m.groupInput.View()))
	s.WriteString("\n\n")
	s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Enter/Ctrl+S: create • Esc: cancel"))

	return s.String()
}

func (m model) renderHelp() string {
	var bindings []key.Binding

	switch m.mode {
	case viewModeList:
		bindings = m.keys.ListKeyMap()
	case viewModeEntry:
		bindings = m.keys.EntryKeyMap()
	case viewModeSearch:
		bindings = m.keys.SearchKeyMap()
	case viewModeAddForm, viewModeEditForm:
		bindings = m.keys.FormKeyMap()
	case viewModeConfirm:
		bindings = m.keys.ConfirmKeyMap()
	default:
		bindings = m.keys.ShortHelp()
	}

	if m.help.ShowAll {
		return m.help.FullHelpView(m.keys.FullHelp())
	}

	return m.help.ShortHelpView(bindings)
}

func (m model) getPath() string {
	var parts []string
	cursor := m.groupCursor
	for cursor != nil {
		parts = append([]string{cursor.group.Name}, parts...)
		cursor = cursor.previous
	}
	return strings.Join(parts, " / ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
