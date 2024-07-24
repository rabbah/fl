package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type uInputModel struct {
	textarea textarea.Model
}

// message event for reading later (when contacting the API)
type userPromptMsg struct {
	prompt string
}

func newUserInputModel() uInputModel {
	m := uInputModel{}
	m.textarea = textarea.New()
	m.textarea.Placeholder = uInputPlaceholderTxt

	m.textarea.SetWidth(uInputWidth)
	m.textarea.SetHeight(uInputHeight)

	m.textarea.Prompt = uInputPrompt
	m.textarea.CharLimit = uInputCharLimit

	// the below are not expected to change so they are not moved to "styling"
	m.textarea.Focus()

	m.textarea.ShowLineNumbers = false
	m.textarea.KeyMap.InsertNewline.SetEnabled(false)
	return m
}

func (m uInputModel) Init() tea.Cmd {
	return textarea.Blink
}

func signalUserInput(prompt string) tea.Cmd {
	return func() tea.Msg {
		return userPromptMsg{prompt: prompt}
	}
}

func (m uInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg.(type) {
	case tea.KeyMsg:
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m uInputModel) UpdateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// emit a signal containing prompt value
			cmds = append(cmds, signalUserInput(m.textarea.Value()))
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m uInputModel) View() string {
	return m.textarea.View()
}
