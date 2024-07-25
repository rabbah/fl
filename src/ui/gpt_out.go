package ui

import (
	"fl/exec"
	"fl/web"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// viewstate keeps track of the current state
// helps blocking new prompts until current execution loop is finished
type gptViewState uint

const (
	waitForPrompt gptViewState = iota
	waitForCommand
	waitForUserCommandExec
	waitForCommandExec
)

type gptViewModel struct {
	viewport viewport.Model
	state    gptViewState
	content  string
	prompt   string
	command  string
	output   string
}

// message event for reading later (when contacting the API)
type webCmdGenMsg struct {
	res string
	err error
}

// message event for command execution
type cmdExecMsg struct {
	res string
	err error
}

func newGPTViewModel() gptViewModel {
	m := gptViewModel{}
	m.viewport = viewport.New(gptViewWidth, gptViewHeight)
	m.viewport.SetContent(gptPlaceholderTxt)
	m.state = waitForPrompt
	m.content = ""
	m.prompt = ""
	m.command = ""
	m.output = ""
	return m
}

func (m gptViewModel) Init() tea.Cmd {
	return nil
}

func sendPrompt(prompt string) tea.Cmd {
	return func() tea.Msg {
		res, err := web.GenerateCommand(prompt)
		return webCmdGenMsg{res, err}
	}
}

func execCmd(prompt string) tea.Cmd {
	return func() tea.Msg {
		res, err := exec.Exec(prompt)
		return cmdExecMsg{res, err}
	}
}

func (m gptViewModel) updatePrompt(msg userPromptMsg) (gptViewModel, tea.Cmd) {
	m.prompt = msg.prompt
	m.content = "Prompt: " + m.prompt
	m.viewport.SetContent(m.content)
	m.state = waitForCommand

	cmd := sendPrompt(m.prompt)
	return m, cmd
}

func (m gptViewModel) updateCommand(msg webCmdGenMsg) (gptViewModel, tea.Cmd) {
	if msg.err != nil {
		m.command = "Something went wrong when generating command!"
		// should be logged!
	} else {
		m.command = msg.res
	}
	m.content = m.content + "\n\nCommand: " + m.command
	m.content = m.content + "\n\nDo you wish to execute the above? (enter = yes, anything else = no)"
	m.viewport.SetContent(m.content)
	m.state = waitForUserCommandExec
	return m, nil
}

func (m gptViewModel) updateExecPrompt(msg tea.KeyMsg) (gptViewModel, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg.String() {
	case "enter":
		m.state = waitForCommandExec
		m.content = ""
		cmd = execCmd(m.command)
	default:
		m.state = waitForPrompt
		m.content = "Waiting for next prompt..."
	}
	m.viewport.SetContent(m.content)
	return m, cmd
}

func (m gptViewModel) updateExec(msg cmdExecMsg) (gptViewModel, tea.Cmd) {
	if msg.err != nil {
		m.output = "Something went wrong when executing command!"
		// should be logged!
	} else {
		m.output = msg.res
	}
	m.content = m.output
	m.viewport.SetContent(m.content)
	m.state = waitForPrompt
	return m, nil
}

func (m gptViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case userPromptMsg:
		if m.state == waitForPrompt {
			m, cmd = m.updatePrompt(msg)
			cmds = append(cmds, cmd)
		}
	case webCmdGenMsg:
		if m.state == waitForCommand {
			m, cmd = m.updateCommand(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		if m.state == waitForUserCommandExec {
			m, cmd = m.updateExecPrompt(msg)
			cmds = append(cmds, cmd)
		}
	case cmdExecMsg:
		if m.state == waitForCommandExec {
			m, cmd = m.updateExec(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m gptViewModel) View() string {
	return m.viewport.View()
}
