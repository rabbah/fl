package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type gptViewModel struct {
	viewport viewport.Model
	messages []string
}

func newGPTViewModel() gptViewModel {
	m := gptViewModel{}
	m.viewport = viewport.New(gptViewWidth, gptViewHeight)
	m.viewport.SetContent(gptPlaceholderTxt)
	m.messages = []string{}
	return m
}

func (m gptViewModel) Init() tea.Cmd {
	return nil
}

func (m gptViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg.(type) {
	case userPromptMsg:
		m.viewport.SetContent(
			"prompt: " +
				msg.(userPromptMsg).prompt,
		)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m gptViewModel) View() string {
	return m.viewport.View()
}
