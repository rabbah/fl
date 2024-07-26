package ui

import (
	"fl/helpers"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
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
const default_placeholder = "enter filename"

var (
	flags_allowed = []string{"autoexecute", "noexec", "output"}
)

// return string for future opts
func (m flagsModel) flagsSelected(cursor int) (opts string, ok bool) {
	switch cursor {
	case autoexecute:
		return "", m.Flags.Autoexecute
	case noexec:
		return "", m.Flags.Noexec
	case output:
		return "", m.Flags.Output
	default:
		return "", false
	}
}

func (m flagsModel) setFlag(cursor int, setValue bool) {
	switch cursor {
	case autoexecute:
		m.Flags.Autoexecute = setValue
	case noexec:
		m.Flags.Noexec = setValue
	case output:
		m.Flags.Output = setValue
	default:
		return
	}
}

func (m flagsModel) toggleFlag(cursor int) (newValue bool) {
	switch cursor {
	case autoexecute:
		m.Flags.Autoexecute = !m.Flags.Autoexecute
		return m.Flags.Autoexecute
	case noexec:
		m.Flags.Noexec = !m.Flags.Noexec
		return m.Flags.Noexec
	case output:
		m.Flags.Output = !m.Flags.Output
		return m.Flags.Output
	default:
		return false
	}
}

func (m flagsModel) validateFlags(cursor int) {
	// check the mutual exclusive flags, flip the one not recently toggled
	if m.Flags.Autoexecute && m.Flags.Noexec {
		if cursor == autoexecute {
			m.toggleFlag(noexec)
		} else {
			m.toggleFlag(autoexecute)
		}
	}

	// provide default outfile name if not already given
	if m.Flags.Output && m.Flags.Outfile == "" {
		m.Flags.Outfile = default_outfile_name
	}
}

type flagsModel struct {
	flags_allowed []string
	flags_cursor  int
	Flags         *helpers.FlagStruct
	outfileEntry  textinput.Model
}

func newOutputFilenameModel(outfile string) textinput.Model {
	m := textinput.New()
	m.Prompt = "output "
	m.Placeholder = default_placeholder
	if outfile != "" {
		m.SetValue(outfile)
	}

	return m
}

func newFlagsModel(Flags *helpers.FlagStruct) flagsModel {
	m := flagsModel{}
	m.Flags = Flags
	m.flags_allowed = flags_allowed
	// if initialized with output enabled, we know outfile name was parsed with helpers.argparse
	m.outfileEntry = newOutputFilenameModel(m.Flags.Outfile)
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

// blur the cursors when unfocused
func (m flagsModel) BlurUnfocused() flagsModel {
	m.outfileEntry.Blur()
	return m
}

// update code for only when focused
func (m flagsModel) UpdateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.flags_cursor > 0 {
				m.flags_cursor--
			} else {
				m.flags_cursor = len(flags_allowed) - 1
			}
		case "down":
			if m.flags_cursor < len(flags_allowed)-1 {
				m.flags_cursor++
			} else {
				m.flags_cursor = 0
			}
		case "enter":
			// if output flag, swap between placeholder and value when toggling
			m.toggleFlag(m.flags_cursor)
			m.validateFlags(m.flags_cursor)
		}
	}

	// allow user to modify output flag if cursor is over it
	// set false if output is empty, true otherwise!
	if m.flags_cursor == output {
		m.outfileEntry.Focus()

		m.outfileEntry, cmd = m.outfileEntry.Update(msg)
		outfile := m.outfileEntry.Value()
		if helpers.IsEmpty(outfile) {
			m.setFlag(output, false)
			m.Flags.Outfile = ""
		} else {
			m.setFlag(output, true)
			m.Flags.Outfile = outfile
		}

	} else {
		m.outfileEntry.Blur()
	}

	return m, cmd
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
		// output row is treated specially
		if i == output {
			choice = m.outfileEntry.View()
		}
		s += fmt.Sprintf("%s [%s] %s", cursor, checked, choice)

		s += "\n"
	}

	return s
}
