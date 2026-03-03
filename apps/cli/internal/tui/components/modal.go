package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
)

func RenderModal(title string, body string, width int, height int, styleSet styles.StyleSet) string {
	contentWidth := min(max(width-16, 36), 72)
	box := styleSet.PanelBorder.Copy().
		Width(contentWidth).
		Padding(1, 2).
		BorderForeground(styleSet.BorderColor)
	hint := styleSet.Dim.Render("Type to filter  |  / manual filter  |  enter select  |  esc close")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSet.PanelTitle.Render(title),
		hint,
		"",
		body,
	)
	panel := box.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, panel)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
