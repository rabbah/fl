package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// styling
var (
	// specific models
	flagsStyle = lipgloss.NewStyle().
			Width(30).
			Height(5).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	// help
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	// extra effects
	focusedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

func setFocus(baseView lipgloss.Style) (style lipgloss.Style) {
	return lipgloss.NewStyle().Inherit(focusedModelStyle).Inherit(baseView)
}

func viewBuilder(m mainModel,
	flagStyle1 lipgloss.Style,
	flagStyle2 lipgloss.Style,
	help string,
) (render string) {
	render = lipgloss.JoinHorizontal(
		lipgloss.Top,
		flagStyle1.Render(fmt.Sprintf("%4s", m.models[flagsView].View())),
		flagStyle2.Render(m.models[flagsView2].View()),
	)
	render += helpStyle.Render(help)
	return render
}
