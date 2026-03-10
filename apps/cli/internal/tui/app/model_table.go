package app

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
)

func apiStatus(tokenSet bool) string {
	if tokenSet {
		return "CONNECTED"
	}
	return "OFFLINE"
}

func filterRows(rows []table.Row, rowIDs []string, query string) ([]table.Row, []string) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return rows, rowIDs
	}
	filtered := make([]table.Row, 0, len(rows))
	filteredIDs := make([]string, 0, len(rows))
	for i, row := range rows {
		joined := strings.ToLower(strings.Join(row, " "))
		if strings.Contains(joined, query) {
			filtered = append(filtered, row)
			if i < len(rowIDs) {
				filteredIDs = append(filteredIDs, rowIDs[i])
			}
		}
	}
	return filtered, filteredIDs
}

func defaultSortConfig(columns []table.Column) SortConfig {
	if len(columns) == 0 {
		return SortConfig{Column: -1, Desc: true}
	}
	preferred := []string{"Started", "Updated", "Received", "Created", "Last Used"}
	for _, target := range preferred {
		for i, col := range columns {
			if strings.EqualFold(strings.TrimSpace(col.Title), target) {
				return SortConfig{Column: i, Desc: true}
			}
		}
	}
	return SortConfig{Column: 0, Desc: true}
}

func columnsWithSortIndicators(columns []table.Column, sortColumn int, desc bool) []table.Column {
	decorated := make([]table.Column, len(columns))
	copy(decorated, columns)
	for i := range decorated {
		title := strings.TrimSpace(decorated[i].Title)
		if i == sortColumn && sortColumn >= 0 && sortColumn < len(decorated) {
			arrow := " ▲"
			if desc {
				arrow = " ▼"
			}
			title += arrow
		}
		if i > 0 {
			title = " " + title
		}
		decorated[i].Title = title
	}
	return decorated
}

func (m *Model) cycleSortColumn() {
	if len(m.columns) == 0 {
		return
	}
	indexes := serverSortableColumnIndexesForView(m.view, m.columns)
	if len(indexes) == 0 {
		indexes = make([]int, len(m.columns))
		for i := range m.columns {
			indexes[i] = i
		}
	}
	if len(indexes) == 0 {
		return
	}

	if m.sortColumn < 0 || !containsIndex(indexes, m.sortColumn) {
		m.sortColumn = indexes[0]
	} else {
		for i, idx := range indexes {
			if idx == m.sortColumn {
				m.sortColumn = indexes[(i+1)%len(indexes)]
				break
			}
		}
	}
	cfg := m.sortByView[m.view]
	cfg.Column = m.sortColumn
	cfg.Desc = m.sortDesc
	m.sortByView[m.view] = cfg
	m.refreshView()
}

func (m *Model) toggleSortDirection() {
	m.sortDesc = !m.sortDesc
	cfg := m.sortByView[m.view]
	cfg.Column = m.sortColumn
	cfg.Desc = m.sortDesc
	m.sortByView[m.view] = cfg
	m.refreshView()
}

func sortOrderFromDesc(desc bool) string {
	if desc {
		return "desc"
	}
	return "asc"
}

func setAPISortSpec(target *apiListSort, by string, order string) bool {
	by = strings.TrimSpace(by)
	order = strings.ToLower(strings.TrimSpace(order))
	if by == "" || (order != "asc" && order != "desc") {
		return false
	}
	if target.By == by && target.Order == order {
		return false
	}
	target.By = by
	target.Order = order
	return true
}

func (m *Model) syncServerSortForCurrentView() bool {
	if m.sortColumn < 0 || m.sortColumn >= len(m.columns) {
		return false
	}

	title := strings.TrimSpace(m.columns[m.sortColumn].Title)
	field, ok := serverSortFieldForViewColumn(m.view, title)
	if !ok {
		return false
	}
	order := sortOrderFromDesc(m.sortDesc)

	switch m.view {
	case ViewWorkflows:
		return setAPISortSpec(&m.snapshotSort.Workflows, field, order)
	case ViewDashboard, ViewRuns:
		return setAPISortSpec(&m.snapshotSort.Runs, field, order)
	case ViewTriggers:
		return setAPISortSpec(&m.snapshotSort.Triggers, field, order)
	case ViewEvents:
		return setAPISortSpec(&m.snapshotSort.Events, field, order)
	case ViewSecrets:
		return setAPISortSpec(&m.snapshotSort.Secrets, field, order)
	}

	return false
}

func serverSortFieldForViewColumn(view ViewID, title string) (string, bool) {
	title = strings.TrimSpace(title)
	switch view {
	case ViewWorkflows:
		if title == "Updated" {
			return "updatedAt", true
		}
	case ViewDashboard, ViewRuns:
		if title == "Started" {
			return "createdAt", true
		}
	case ViewTriggers:
		if title == "Created" {
			return "createdAt", true
		}
	case ViewEvents:
		if title == "Received" {
			return "receivedAt", true
		}
	case ViewSecrets:
		if title == "Created" {
			return "createdAt", true
		}
	}
	return "", false
}

func serverSortableColumnIndexesForView(view ViewID, columns []table.Column) []int {
	indexes := make([]int, 0, len(columns))
	for i, col := range columns {
		if _, ok := serverSortFieldForViewColumn(view, col.Title); ok {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func shouldUseClientSideSort(view ViewID, columns []table.Column) bool {
	return len(serverSortableColumnIndexesForView(view, columns)) == 0
}

func containsIndex(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func truncateRows(rows []table.Row, columns []table.Column) []table.Row {
	if len(rows) == 0 {
		return rows
	}
	widths := make([]int, 0, len(columns))
	for _, col := range columns {
		widths = append(widths, max(col.Width, 1))
	}
	for i := range rows {
		for c := 0; c < len(rows[i]) && c < len(widths); c++ {
			rows[i][c] = ansi.Truncate(rows[i][c], widths[c], "")
		}
	}
	return rows
}

func initialSize() (int, int) {
	width, height := 80, 24
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && w > 0 && h > 0 {
		width = w
		height = h
	}
	return width, height
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
