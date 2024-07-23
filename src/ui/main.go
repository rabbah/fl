package ui

import (
	"fl/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

/*
Based on examples found from the following links

Find them at:
https://github.com/charmbracelet/bubbletea/tree/master/tutorials
*/

// sessionState is used to track which model is focused
type sessionState uint

// expected views
/* @IMPORTANT
 * make sure that the sessionstate declarations are the SAME ORDER as views being added to map!!!
 * this ensures the map is sorted SEQUENTIALLY
 */
const (
	gptView sessionState = iota
	flagsView
)

type mainModel struct {
	state     sessionState
	altscreen bool
	quitting  bool
	models    map[sessionState]tea.Model
}

func newModel(Flags *helpers.FlagStruct) mainModel {
	m := mainModel{state: flagsView}
	m.models = make(map[sessionState]tea.Model)
	/* @IMPORTANT
	 * make sure that the views added to the map are the SAME ORDER as the sessionstate declarations!!!
	 * this ensures the map is sorted SEQUENTIALLY
	 */
	m.models[gptView] = newGPTViewModel()
	m.models[flagsView] = newFlagsModel(Flags)
	return m
}

func (m mainModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, model := range m.models {
		cmds = append(cmds, model.Init())
	}

	return tea.Batch(cmds...)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
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
			// session state is a uint
			m.state = m.state + 1
			if m.models[m.state] == nil {
				m.state = 0
			}
		}
	}

	// global updates for subviews (for spinners etc)

	for sessionState, model := range m.models {
		m.models[sessionState], cmd = model.Update(msg)
		cmds = append(cmds, cmd)
	}

	// focused updates for subviews (for items allowed only in focus)
	// must explicitly define as UpdateFocused is not part of the tea.Model interface
	switch m.state {
	case flagsView:
		m.models[flagsView], cmd = m.models[flagsView].(flagsModel).UpdateFocused(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string

	// clear screen on quit
	if m.quitting {
		return ""
	}

	if m.state == gptView {
		help := ""
		s += viewBuilder(m, setFocus(gptStyle), flagsStyle, help)
	} else if m.state == flagsView {
		help := "\nenter: toggle flag • j/up: scroll up • k/down: scroll down"
		s += viewBuilder(m, gptStyle, setFocus(flagsStyle), help)
	}
	s += helpStyle.Render("\ntab: focus next • q: exit • space: swap alt view\n")
	return s
}

func RunProgram(Flags *helpers.FlagStruct) (err error) {
	initialModel := newModel(Flags)
	p := tea.NewProgram(initialModel)
	_, err = p.Run()
	return err
}
