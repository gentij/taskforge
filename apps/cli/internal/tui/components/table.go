package components

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
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

func StyleRows(rows []table.Row, selected int, styleSet styles.StyleSet) []table.Row {
	styled := make([]table.Row, len(rows))
	for i, row := range rows {
		styledRow := make(table.Row, len(row))
		applyAlt := i%2 == 1 && i != selected
		var rowStyle lipgloss.Style
		if applyAlt {
			rowStyle = styleSet.RowAlt
		}
		for j, cell := range row {
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
