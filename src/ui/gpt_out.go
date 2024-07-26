package ui

import (
	"fl/exec"
	"fl/helpers"
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
	Flags    *helpers.FlagStruct
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

func newGPTViewModel(Flags *helpers.FlagStruct) gptViewModel {
	m := gptViewModel{}
	m.viewport = viewport.New(gptViewWidth, gptViewHeight)
	// Flags is a shared global - do NOT rely on for setting the on-screen prompt or m.prompt
	prompt := Flags.Prompt
	// use placeholder if prompt not passed through CLI
	if prompt != "" {
		m.content = "Prompt: " + prompt
		// wait for command to be generated (signal is sent with Init)
		m.state = waitForCommand

	} else {
		m.content = gptPlaceholderTxt
		// wait for user executing prompt
		m.state = waitForPrompt
	}
	m.Flags = Flags
	m.prompt = prompt
	m.viewport.SetContent(m.content)
	m.command = ""
	m.output = ""
	return m
}

func (m gptViewModel) Init() tea.Cmd {
	if m.prompt != "" {
		return sendPrompt(m.prompt)
	}
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
		// err. restart state machine
		m.content = "Something went wrong when generating command! " + gptPlaceholderTxt
		m.command = ""
		m.state = waitForPrompt
		m.viewport.SetContent(m.content)
		return m, changeModelFocus(uInputView)
		// should be logged!
	} else {
		m.command = msg.res
	}

	if m.Flags.Autoexecute {
		// skip asking for user input (and execute command + change focus)
		m.state = waitForCommandExec
		return m, tea.Batch(execCmd(m.command), changeModelFocus(uInputView))
	} else if m.Flags.Noexec {
		// skip asking for user input (dont execute command but change focus)
		m.content = m.command
		m.viewport.SetContent(m.content)
		m.state = waitForPrompt
		return m, changeModelFocus(uInputView)
	} else {
		// change to prompt user state
		m.content = m.content + "\n\nDo you wish to execute the below? (y/n)"
		m.content = m.content + "\n\n" + m.command
		m.viewport.SetContent(m.content)
		m.state = waitForUserCommandExec
		return m, nil
	}
}

func (m gptViewModel) updateExecPrompt(msg tea.KeyMsg) (gptViewModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg.String() {
	case "y":
		m.state = waitForCommandExec
		m.content = ""
		// request command execution & change state
		cmds = append(cmds,
			execCmd(m.command),
			changeModelFocus(uInputView),
		)
	case "n":
		m.state = waitForPrompt
		m.content = "Waiting for next prompt..."
		// only restart state and change focus
		cmds = append(cmds,
			changeModelFocus(uInputView),
		)
	default:
		// dont issue command or change state
	}

	m.viewport.SetContent(m.content)
	return m, tea.Batch(cmds...)
}

func (m gptViewModel) updateExec(msg cmdExecMsg) (gptViewModel, tea.Cmd) {
	if msg.err != nil {
		// err. restart state machine
		m.content = "Something went wrong when executing command! " + gptPlaceholderTxt
		m.output = ""
		m.state = waitForPrompt
		m.viewport.SetContent(m.content)
		return m, changeModelFocus(uInputView)
		// should be logged!
	} else if helpers.IsEmpty(msg.res) {
		m.output = "Command response was empty"
	} else {
		m.output = msg.res
	}
	m.content = m.output
	m.viewport.SetContent(m.content)
	m.state = waitForPrompt
	return m, nil
}

func (m gptViewModel) UpdateFocused(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == waitForUserCommandExec {
			m, cmd = m.updateExecPrompt(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m gptViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case userPromptMsg:
		if m.state == waitForPrompt {
			m, cmd = m.updatePrompt(msg)
		}
	case webCmdGenMsg:
		if m.state == waitForCommand {
			m, cmd = m.updateCommand(msg)
		}
	case cmdExecMsg:
		if m.state == waitForCommandExec {
			m, cmd = m.updateExec(msg)
		}
	}

	return m, cmd
}

func (m gptViewModel) View() string {
	return m.viewport.View()
}
