package ui

import (
	"fl/web"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type explainViewModel struct {
	viewport viewport.Model
}

// message event for reading later (when contacting the API)
type webExplainMsg struct {
	res string
	err error
}

func newExplainViewModel() explainViewModel {
	m := explainViewModel{}
	m.viewport = viewport.New(explainViewWidth, explainViewHeight)
	m.viewport.SetContent("This box will explain commands as they are generated as long as the explain option is set.")
	return m
}

func (m explainViewModel) Init() tea.Cmd {
	return nil
}

func sendExplain(cmd string, language string) tea.Cmd {
	return func() tea.Msg {
		res, err := web.ExplainCommand(cmd, language)
		return webExplainMsg{res, err}
	}
}

func (m explainViewModel) updateExplain(msg webExplainMsg) (explainViewModel, tea.Cmd) {
	m.viewport.SetContent(msg.res)

	return m, nil
}

func (m explainViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case webExplainMsg:
		m, cmd = m.updateExplain(msg)
	}

	return m, cmd
}

func (m explainViewModel) View() string {
	return m.viewport.View()
}
