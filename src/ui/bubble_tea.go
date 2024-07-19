package ui

import (
	"fl/helpers"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
Based on examples found from the following links

Find them at:
https://github.com/charmbracelet/bubbletea/tree/master/tutorials
*/

// sessionState is used to track which model is focused
type sessionState uint

// expected views
const (
	flagsView sessionState = iota
	flagsView2
)

// styling
var (
	// focus style
	modelStyle = lipgloss.NewStyle().
			Width(30).
			Height(5).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(30).
				Height(5).
				Align(lipgloss.Left, lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type mainModel struct {
	state       sessionState
	altscreen   bool
	flagsModel  flagsModel
	flagsModel2 flagsModel
}

func newModel() mainModel {
	m := mainModel{state: flagsView}
	m.flagsModel = newFlagsModel()
	m.flagsModel2 = newFlagsModel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		m.flagsModel.Init(),
		m.flagsModel2.Init(),
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd
		case "tab":
			if m.state == flagsView {
				m.state = flagsView2
			} else {
				m.state = flagsView
			}
		}
	}

	// global updates for subviews
	m.flagsModel, cmd = m.flagsModel.Update(msg)
	cmds = append(cmds, cmd)
	m.flagsModel2, cmd = m.flagsModel2.Update(msg)
	cmds = append(cmds, cmd)

	// focused updates for subviews
	switch m.state {
	case flagsView:
		m.flagsModel, cmd = m.flagsModel.UpdateFocused(msg)
		cmds = append(cmds, cmd)
	case flagsView2:
		m.flagsModel2, cmd = m.flagsModel2.UpdateFocused(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	if m.state == flagsView {
		s += lipgloss.JoinVertical(
			lipgloss.Top,
			focusedModelStyle.Render(fmt.Sprintf("%4s", m.flagsModel.View())),
			modelStyle.Render(m.flagsModel2.View()),
		)
		s += helpStyle.Render("\nenter: toggle flag")
	} else if m.state == flagsView2 {
		s += lipgloss.JoinVertical(
			lipgloss.Top,
			modelStyle.Render(fmt.Sprintf("%4s", m.flagsModel.View())),
			focusedModelStyle.Render(m.flagsModel2.View()),
		)
	}
	s += helpStyle.Render("\ntab: focus next • q: exit • space: swap alt view\n")
	return s
}

func RunProgram(Flags helpers.FlagStruct) (err error) {
	initialModel := newModel()
	p := tea.NewProgram(initialModel)
	_, err = p.Run()
	return err
}
