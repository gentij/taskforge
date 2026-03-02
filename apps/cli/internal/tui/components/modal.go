package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
)

func RenderModal(title string, body string, width int, height int, styleSet styles.StyleSet) string {
	box := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styleSet.BorderColor)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSet.PanelTitle.Render(title),
		"",
		body,
	)
	panel := box.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, panel)
}
