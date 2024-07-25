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
	uInputView sessionState = iota
	gptView
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
	m.models[uInputView] = newUserInputModel()
	m.models[gptView] = newGPTViewModel()
	m.models[flagsView] = newFlagsModel(Flags)

	// start with focus on the input (see constructor...)
	m.state = uInputView
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
		case "ctrl+c", "q", "esc":
			m.quitting = true
			if m.altscreen {
				cmds = append(cmds, tea.ExitAltScreen)
			}
			cmds = append(cmds, tea.ClearScreen, tea.Quit)
			// clear screen and quit!
			return m, tea.Batch(cmds...)
		case "ctrl+f":
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
	case uInputView:
		m.models[uInputView], cmd = m.models[uInputView].(uInputModel).UpdateFocused(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string

	if m.quitting {
		return ""
	}

	help := ""
	switch m.state {
	case uInputView:
		help := "\nenter: submit prompt"
		s += viewBuilder(m, setFocus(uInputStyle), gptStyle, flagsStyle, help)
	case gptView:
		s += viewBuilder(m, uInputStyle, setFocus(gptStyle), flagsStyle, help)
	case flagsView:
		help := "\nenter: toggle flag • j/up: scroll up • k/down: scroll down"
		s += viewBuilder(m, uInputStyle, gptStyle, setFocus(flagsStyle), help)
	}

	return s
}

func RunProgram(Flags *helpers.FlagStruct) (err error) {
	initialModel := newModel(Flags)
	p := tea.NewProgram(initialModel)
	_, err = p.Run()
	return err
}
