package ui

import (
	"fl/helpers"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// flag structure
const (
	autoexecute = iota
	noexec
)

var (
	flags_allowed = []string{"autoexecute", "noexec"}
)

func (m flagsModel) flagsSelected(cursor int) (opts string, ok bool) {
	switch cursor {
	case autoexecute:
		return "", m.flags.Autoexecute
	case noexec:
		return "", m.flags.Noexec
	default:
		return "", false
	}
}

func (m flagsModel) toggleFlag(cursor int) (newValue bool) {
	switch cursor {
	case autoexecute:
		m.flags.Autoexecute = !m.flags.Autoexecute
		return m.flags.Autoexecute
	case noexec:
		m.flags.Noexec = !m.flags.Noexec
		return m.flags.Noexec
	default:
		return false
	}
}

func (m flagsModel) validateFlags(cursor int) {
	// check the mutual exclusive flags, flip the one not recently toggled
	if m.flags.Autoexecute && m.flags.Noexec {
		if cursor == autoexecute {
			m.toggleFlag(noexec)
		} else {
			m.toggleFlag(autoexecute)
		}
	}
}

type flagsModel struct {
	flags_allowed []string
	flags_cursor  int
	flags         *helpers.FlagStruct
}

func newFlagsModel(Flags *helpers.FlagStruct) flagsModel {
	m := flagsModel{}
	m.flags = Flags
	m.flags_allowed = flags_allowed
	return m
}

func (m flagsModel) Init() tea.Cmd {
	return nil
}

// placeholder code, this update function updates globally
func (m flagsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// update code for only when focused
func (m flagsModel) UpdateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			if m.flags_cursor > 0 {
				m.flags_cursor--
			} else {
				m.flags_cursor = len(flags_allowed) - 1
			}
		case "j", "down":
			if m.flags_cursor < len(flags_allowed)-1 {
				m.flags_cursor++
			} else {
				m.flags_cursor = 0
			}
		case "enter":
			m.toggleFlag(m.flags_cursor)
			m.validateFlags(m.flags_cursor)
		}
	}

	return m, nil
}

func (m flagsModel) View() string {
	var s string

	// Iterate over our choices
	for i, choice := range m.flags_allowed {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.flags_cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.flagsSelected(i); ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	return s
}
