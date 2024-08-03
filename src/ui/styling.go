package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// avoid magic numbers in other files using these
var (
	flagsHeight   = 10
	flagsWidth    = 30
	flagsCursor   = themeStyle.Render(">")
	flagsSelected = themeStyle.Render("x")

	gptViewHeight     = 10
	gptViewWidth      = 78
	gptPlaceholderTxt = "Waiting for prompt..."

	explainViewHeight = 14
	explainViewWidth  = 50

	uInputHeight         = 2
	uInputWidth          = 110
	uInputPlaceholderTxt = "Describe the command you would like to generate..."
	uInputPrompt         = themeStyle.Render("┃ ")
	uInputCharLimit      = 1000

	mainHelp    = "tab: focus next • esc/ctrl+c: exit • ctrl+f: swap alt view\nctrl+y: copy generated command to clipboard"
	helpColor   = postmanColor
	borderColor = postmanColor
	themeColor  = postmanColor
)

// styling
var (
	// specific models
	flagsStyle = lipgloss.NewStyle().
			Width(flagsWidth).
			Height(flagsHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())
	gptStyle = lipgloss.NewStyle().
			Width(gptViewWidth).
			Height(gptViewHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())
	uInputStyle = lipgloss.NewStyle().
			Width(uInputWidth).
			Height(uInputHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())
	explainStyle = lipgloss.NewStyle().
			Width(explainViewWidth).
			Height(explainViewHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())
	// help
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColor)).
			Width(uInputWidth)
	// themed styling
	themeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor))
	// extra effects
	focusedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color(borderColor))
)

func setFocus(baseView lipgloss.Style) (style lipgloss.Style) {
	return lipgloss.NewStyle().Inherit(focusedModelStyle).Inherit(baseView)
}

func viewBuilder(m mainModel,
	uInputStyle lipgloss.Style,
	gptStyle lipgloss.Style,
	flagStyle lipgloss.Style,
	explainStyle lipgloss.Style,
	help string,
) (render string) {

	// main
	render = lipgloss.JoinHorizontal(
		lipgloss.Left,
		gptStyle.Render(m.models[gptView].View()),
		flagStyle.Render(fmt.Sprintf("%4s", m.models[flagsView].View())),
	)
	render = lipgloss.JoinVertical(
		lipgloss.Left,
		render,
		uInputStyle.Render(m.models[uInputView].View()),
	)
	render = lipgloss.JoinHorizontal(
		lipgloss.Left,
		render,
		explainStyle.Render(m.models[explainView].View()),
	)

	// help
	render = lipgloss.JoinVertical(lipgloss.Top, render, helpStyle.Render(help))
	render = lipgloss.JoinVertical(lipgloss.Top, render, helpStyle.Render(mainHelp))

	return render
}
