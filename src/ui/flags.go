package ui

import (
	"fl/helpers"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// flag structure
// when adding flags for TUI support, update at minimum:
// 1. this structure
// 2. the var structure
// 3. flagSelected
// 4. toggleFlag
const (
	autoexecute = iota
	noexec
	output
)

// avoid prompting for input by specifying default outfile name
const default_outfile_name = "fl.out"

var (
	flags_allowed = []string{"autoexecute", "noexec", "output"}
)

// return string for future opts
func (m flagsModel) flagsSelected(cursor int) (opts string, ok bool) {
	switch cursor {
	case autoexecute:
		return "", m.flags.Autoexecute
	case noexec:
		return "", m.flags.Noexec
	case output:
		return "", m.flags.Output
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
	case output:
		m.flags.Output = !m.flags.Output
		return m.flags.Output
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

	// provide default outfile name if not already given
	if m.flags.Output && m.flags.Outfile == "" {
		m.flags.Outfile = default_outfile_name
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
			cursor = flagsCursor // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.flagsSelected(i); ok {
			checked = flagsSelected // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s", cursor, checked, choice)

		// add filename if row contains output and is selected
		if i == output && m.flags.Output {
			s += " " + m.flags.Outfile
		}

		s += "\n"
	}

	return s
}
