package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styling
var (
	flags     = []string{"output", "autoexecute"}
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type flagsModel struct {
	flags          []string
	flags_cursor   int
	flags_selected map[int]struct{}
}

func newFlagsModel() flagsModel {
	m := flagsModel{}
	m.flags = flags
	m.flags_selected = make(map[int]struct{})
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
				m.flags_cursor = len(flags) - 1
			}
		case "j", "down":
			if m.flags_cursor < len(flags)-1 {
				m.flags_cursor++
			} else {
				m.flags_cursor = 0
			}
		case "enter":
			_, ok := m.flags_selected[m.flags_cursor]
			if ok {
				delete(m.flags_selected, m.flags_cursor)
			} else {
				m.flags_selected[m.flags_cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m flagsModel) View() string {
	var s string

	// Iterate over our choices
	for i, choice := range m.flags {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.flags_cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.flags_selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	return s
}
