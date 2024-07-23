package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// avoid magic numbers in other files using these
var (
	flagsHeight = 3
	flagsWidth  = 30

	gptViewHeight = 4
	gptViewWidth  = 80

	helpColor   = "241"
	borderColor = "69"
)

// styling
var (
	// specific models
	flagsStyle = lipgloss.NewStyle().
			Width(flagsWidth).
			Height(flagsHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	gptStyle = lipgloss.NewStyle().
			Width(gptViewWidth).
			Height(gptViewHeight).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	// help
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(helpColor))
	// extra effects
	focusedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(borderColor))
)

func setFocus(baseView lipgloss.Style) (style lipgloss.Style) {
	return lipgloss.NewStyle().Inherit(focusedModelStyle).Inherit(baseView)
}

func viewBuilder(m mainModel,
	gptStyle lipgloss.Style,
	flagStyle lipgloss.Style,
	help string,
) (render string) {
	render = lipgloss.JoinHorizontal(
		lipgloss.Top,
		gptStyle.Render(m.models[gptView].View()),
		flagStyle.Render(fmt.Sprintf("%4s", m.models[flagsView].View())),
	)
	render += helpStyle.Render(help)
	return render
}
