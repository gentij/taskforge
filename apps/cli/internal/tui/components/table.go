package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
)

func NewTable(columns []table.Column, rows []table.Row, width int, height int, styleSet styles.StyleSet) table.Model {
	model := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	model.SetWidth(width)
	model.SetStyles(TableStyles(styleSet))

	return model
}

func TableStyles(styleSet styles.StyleSet) table.Styles {
	style := table.DefaultStyles()
	style.Header = styleSet.TableHeader
	style.Cell = styleSet.TableCell
	style.Selected = styleSet.TableSelected
	return style
}

func StyleRows(rows []table.Row, columns []table.Column, selected int, styleSet styles.StyleSet) []table.Row {
	styled := make([]table.Row, len(rows))
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = max(col.Width, 1)
	}
	for i, row := range rows {
		styledRow := make(table.Row, len(row))
		applyAlt := false
		var rowStyle lipgloss.Style
		for j, cell := range row {
			width := 0
			if j < len(widths) {
				width = widths[j]
			}
			cell = padCell(cell, width)
			if applyAlt {
				styledRow[j] = rowStyle.Render(cell)
			} else {
				styledRow[j] = cell
			}
		}
		styled[i] = styledRow
	}
	return styled
}

func padCell(value string, width int) string {
	if width <= 0 {
		return value
	}
	value = ansi.Truncate(value, width, "")
	pad := width - ansi.StringWidth(value)
	if pad < 0 {
		pad = 0
	}
	return value + strings.Repeat(" ", pad)
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
