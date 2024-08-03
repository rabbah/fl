package ui

import (
	"fl/helpers"
	"fl/io"
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
	prompt
	explain
	output
	language
)

// avoid prompting for input by specifying default outfile name
const default_placeholder = "enter filename"

var (
	flags_allowed = []string{"autoexecute", "prompt", "explain", "output", "language"}
)

// return string for future opts
func (m flagsModel) flagsSelected(cursor int) (ok bool) {
	switch cursor {
	case autoexecute:
		return m.Config.Autoexec
	case prompt:
		return m.Flags.PromptExec
	case explain:
		return m.Flags.Explain
	case output:
		return m.Flags.Output
	case language:
		return m.language != ""
	default:
		return false
	}
}

func (m flagsModel) setFlag(cursor int, setValue bool) {
	switch cursor {
	case autoexecute:
		m.Config.Autoexec = setValue
	case prompt:
		m.Flags.PromptExec = setValue
	case explain:
		m.Flags.Explain = setValue
	case output:
		m.Flags.Output = setValue
	default:
		return
	}
}

func (m flagsModel) toggleFlag(cursor int) (newValue bool) {
	switch cursor {
	case autoexecute:
		m.Config.Autoexec = !m.Config.Autoexec
		return m.Config.Autoexec
	case prompt:
		m.Flags.PromptExec = !m.Flags.PromptExec
		return m.Flags.PromptExec
	case explain:
		m.Flags.Explain = !m.Flags.Explain
		return m.Flags.Explain
	case output:
		m.Flags.Output = !m.Flags.Output
		return m.Flags.Output
	default:
		return false
	}
}

type flagsModel struct {
	flags_allowed []string
	flags_cursor  int
	Flags         *helpers.FlagStruct
	Config        *io.Config
	outfileEntry  textinput.Model
	language      string
	languageEntry textinput.Model
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

func newLanguageModel(language string, placeholder_lang string) textinput.Model {
	m := textinput.New()
	m.Prompt = "language "
	m.Placeholder = placeholder_lang
	// language will always have a value
	// should be initialized such that if language = config_default then should look empty
	if language != "" {
		m.SetValue(language)
	}

	return m
}

func newFlagsModel(Flags *helpers.FlagStruct, Config *io.Config) flagsModel {
	m := flagsModel{}
	m.Flags = Flags
	m.Config = Config
	m.flags_allowed = flags_allowed
	// if initialized with output enabled, we know outfile name was parsed with helpers.argparse
	m.outfileEntry = newOutputFilenameModel(m.Flags.Outfile)
	// language will always have a value
	// should be initialized such that if language = config_default then should look empty
	m.language = ""
	if m.Flags.Language != m.Config.Language {
		m.language = m.Flags.Language
	}
	m.languageEntry = newLanguageModel(m.language, m.Config.Language)
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
	var cmds []tea.Cmd
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
		}
	}

	m, cmd = m.updateOutfileEntry(msg)
	cmds = append(cmds, cmd)
	m, cmd = m.updateLanguageEntry(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m flagsModel) updateOutfileEntry(msg tea.Msg) (flagsModel, tea.Cmd) {
	// allow user to modify output flag if cursor is over it
	// set false if output is empty, true otherwise!
	var cmd tea.Cmd
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

func (m flagsModel) updateLanguageEntry(msg tea.Msg) (flagsModel, tea.Cmd) {
	// allow user to modify lang flag if cursor is over it
	// set false if output is empty, true otherwise!
	// actual flags value should be m.Config.Language if empty, though!
	var cmd tea.Cmd
	if m.flags_cursor == language {
		m.languageEntry.Focus()

		m.languageEntry, cmd = m.languageEntry.Update(msg)
		language := m.languageEntry.Value()
		if helpers.IsEmpty(language) {
			m.language = ""
			m.Flags.Language = m.Config.Language
		} else {
			m.language = language
			m.Flags.Language = m.language
		}

	} else {
		m.languageEntry.Blur()
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
		if ok := m.flagsSelected(i); ok {
			checked = flagsSelected // selected!
		}

		// Render the row
		// output and language row are treated specially
		if i == output {
			choice = m.outfileEntry.View()
		} else if i == language {
			choice = m.languageEntry.View()
		}
		s += fmt.Sprintf("%s [%s] %s", cursor, checked, choice)

		s += "\n"
	}

	return s
}
