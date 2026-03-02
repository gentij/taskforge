package app

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gentij/taskforge/apps/cli/internal/tui/components"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
)

func Render(m Model) string {
	if m.inspector.Active {
		return m.inspector.Render(m.width, m.height)
	}

	header := renderHeader(m)
	primary := renderPrimary(m)
	footer := renderFooter(m)
	context := ""
	if m.layout.ContextHeight > 0 {
		context = renderContext(m)
	}

	sections := []string{header, primary}
	if context != "" {
		sections = append(sections, context)
	}
	sections = append(sections, footer)
	base := lipgloss.JoinVertical(lipgloss.Left, sections...)
	base = clampToViewport(base, m.width, m.height)

	output := base
	if m.showPalette {
		modal := components.RenderModal("Command Palette", m.palette.View(), m.width, m.height, m.styles)
		output = renderOverlay(base, modal, m)
	}
	if m.showHelp {
		modal := components.RenderModal("Help", m.help.View(m.keys), m.width, m.height, m.styles)
		output = renderOverlay(base, modal, m)
	}
	if m.theme.CRT {
		output = applyScanlines(output, m.width, m.theme.Scanline)
	}
	return output
}

func renderHeader(m Model) string {
	section := viewTitle(m.view)
	refresh := "OFF"
	if m.autoRefresh {
		refresh = "2s"
	}
	apiBadge := styles.Badge(m.styles.BadgeMuted, m.apiStatus)
	if m.apiStatus == "CONNECTED" {
		apiBadge = styles.Badge(m.styles.BadgeSuccess, m.apiStatus)
	}
	accent := m.styles.Accent.Render("▌")
	appName := m.styles.Accent.Render("TASKFORGE")
	sectionLabel := m.styles.PanelTitle.Render(section)
	left := accent + " " + appName + "  " + m.styles.Dim.Render("•") + " " + sectionLabel
	right := "WS: " + m.workspaceName + "  Workers: " + itoa(m.workerCount) + "  ⟳ " + refresh + "  API: " + apiBadge
	line1 := joinLeftRight(left, right, m.width)

	summaryLeft := viewSummary(m)
	summaryRight := ""
	if m.searchQuery != "" {
		summaryRight = "Filter: " + m.searchQuery
	}
	if m.contextCollapsed {
		if summaryRight != "" {
			summaryRight += "  "
		}
		summaryRight += "Context: hidden"
	}
	line2 := joinLeftRight(summaryLeft, summaryRight, m.width)

	line1 = " " + ansi.Truncate(line1, m.width-1, "")
	line2 = " " + ansi.Truncate(line2, m.width-1, "")
	line1 = m.styles.Header.Width(m.width).Render(line1)
	if m.theme.CRT {
		glowStyle := lipgloss.NewStyle().Foreground(m.theme.Glow).Faint(true)
		line2 = glowStyle.Render(line2)
	}
	line2 = m.styles.Header.Width(m.width).Render(line2)
	content := line1 + "\n" + line2
	return clampSection(content, m.width, m.layout.HeaderHeight)
}

func renderPrimary(m Model) string {
	content := ""
	if m.view == ViewDashboard {
		content = renderDashboard(m)
	} else {
		content = strings.TrimRight(m.table.View(), "\n")
	}
	return clampSection(content, m.width, m.layout.PrimaryHeight)
}

func renderDashboard(m Model) string {
	if m.layout.DashboardCardsHeight == 0 {
		return strings.TrimRight(m.table.View(), "\n")
	}
	active := 0
	for _, wf := range m.store.Workflows {
		if wf.Active {
			active++
		}
	}
	failed := 0
	for _, run := range m.store.Runs {
		if run.Status == "FAILED" {
			failed++
		}
	}
	cards := components.RenderCards([]components.StatCard{
		{Title: "Workflows", Value: itoa(active) + " active", Subtitle: itoa(len(m.store.Workflows)) + " total"},
		{Title: "Runs 24h", Value: itoa(len(m.store.Runs)) + " total", Subtitle: "recent activity"},
		{Title: "Failures", Value: itoa(failed), Subtitle: "needs attention"},
	}, m.width, m.styles)
	label := m.styles.Dim.Render("Recent Runs")
	tableView := strings.TrimRight(m.table.View(), "\n")
	rows := []string{cards}
	for i := 0; i < m.layout.DashboardGap; i++ {
		rows = append(rows, "")
	}
	rows = append(rows, label, tableView)
	return strings.Join(rows, "\n")
}

func renderContext(m Model) string {
	width := m.width
	height := m.layout.ContextHeight
	if height <= 0 {
		return ""
	}
	innerWidth := max(width-2, 1)
	innerHeight := max(height-2, 1)
	bodyHeight := max(innerHeight-1, 1)
	content := strings.TrimRight(m.contextViewport.View(), "\n")
	body := lipgloss.Place(innerWidth, bodyHeight, lipgloss.Left, lipgloss.Top, content)
	accent := m.styles.Accent.Render("▌")
	title := accent + " " + m.styles.PanelTitle.Render(contextTitle(m))
	inner := lipgloss.JoinVertical(lipgloss.Left, title, body)
	box := m.styles.PanelBorder.Width(width).Height(height).MaxHeight(height).MaxWidth(width)
	return clampSection(box.Render(inner), width, height)
}

func renderFooter(m Model) string {
	left := m.help.View(m.keys)
	rightParts := []string{}
	if m.searching {
		rightParts = append(rightParts, m.searchInput.View())
	} else if m.contextSearching {
		rightParts = append(rightParts, m.contextSearchInput.View())
	} else if m.searchQuery != "" {
		rightParts = append(rightParts, "Filter: "+m.searchQuery)
	}
	rightParts = append(rightParts, m.paginator.View())
	right := strings.TrimSpace(strings.Join(rightParts, " │ "))
	line := joinLeftRight(left, right, m.width)
	content := m.styles.Footer.Width(m.width).Render(line)
	return clampSection(content, m.width, m.layout.FooterHeight)
}

func joinLeftRight(left string, right string, width int) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if right == "" {
		return ansi.Truncate(left, width, "")
	}
	space := width - lipgloss.Width(left) - lipgloss.Width(right)
	if space < 1 {
		space = 1
	}
	line := left + strings.Repeat(" ", space) + right
	return ansi.Truncate(line, width, "")
}

func renderOverlay(base string, modal string, m Model) string {
	background := m.styles.Dim.Render(base)
	overlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	return clampToViewport(mergeOverlay(background, overlay, m.height), m.width, m.height)
}

func itoa(value int) string {
	return strconv.Itoa(value)
}

func mergeOverlay(base string, overlay string, height int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")
	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		var baseLine string
		if i < len(baseLines) {
			baseLine = baseLines[i]
		}
		overlayLine := ""
		if i < len(overlayLines) {
			overlayLine = overlayLines[i]
		}
		if strings.TrimSpace(overlayLine) == "" {
			lines = append(lines, baseLine)
			continue
		}
		lines = append(lines, overlayLine)
	}
	return strings.Join(lines, "\n")
}

func applyScanlines(content string, width int, color lipgloss.Color) string {
	lines := strings.Split(content, "\n")
	scanStyle := lipgloss.NewStyle().Foreground(color).Faint(true)
	scanChar := scanStyle.Render(".")
	for i, line := range lines {
		if i%2 == 1 {
			lines[i] = scanlineLine(line, width, scanChar)
		}
	}
	return strings.Join(lines, "\n")
}

func scanlineLine(line string, width int, scanChar string) string {
	if width < 1 {
		width = 1
	}
	var builder strings.Builder
	builder.Grow(len(line) + width)
	inEscape := false
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '\x1b' {
			inEscape = true
			builder.WriteByte(ch)
			continue
		}
		if inEscape {
			builder.WriteByte(ch)
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		if ch == ' ' {
			builder.WriteString(scanChar)
		} else {
			builder.WriteByte(ch)
		}
	}
	padding := width - ansi.StringWidth(line)
	for i := 0; i < padding; i++ {
		builder.WriteString(scanChar)
	}
	return builder.String()
}

func clampToViewport(content string, width int, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	lines := strings.Split(content, "\n")
	clamped := make([]string, 0, height)
	for i := 0; i < len(lines) && i < height; i++ {
		clamped = append(clamped, ansi.Truncate(lines[i], width, ""))
	}
	for len(clamped) < height {
		clamped = append(clamped, "")
	}
	return strings.Join(clamped, "\n")
}

func clampSection(content string, width int, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	lines := strings.Split(content, "\n")
	clamped := make([]string, 0, height)
	for i := 0; i < len(lines) && i < height; i++ {
		clamped = append(clamped, ansi.Truncate(lines[i], width, ""))
	}
	for len(clamped) < height {
		clamped = append(clamped, "")
	}
	return strings.Join(clamped, "\n")
}

func viewTitle(view ViewID) string {
	switch view {
	case ViewDashboard:
		return "Dashboard"
	case ViewWorkflows:
		return "Workflows"
	case ViewRuns:
		return "Runs"
	case ViewTriggers:
		return "Triggers"
	case ViewEvents:
		return "Events"
	case ViewSecrets:
		return "Secrets"
	case ViewTokens:
		return "API Tokens"
	default:
		return "Overview"
	}
}

func viewSummary(m Model) string {
	switch m.view {
	case ViewDashboard:
		return "Overview" + "  " + m.styles.Dim.Render("•") + "  Recent runs " + itoa(min(len(m.store.Runs), 6))
	case ViewWorkflows:
		active := 0
		for _, wf := range m.store.Workflows {
			if wf.Active {
				active++
			}
		}
		return "Workflows " + itoa(len(m.store.Workflows)) + " total" + "  " + m.styles.Dim.Render("•") + "  " + itoa(active) + " active"
	case ViewRuns:
		failed := 0
		running := 0
		queued := 0
		for _, run := range m.store.Runs {
			switch run.Status {
			case "FAILED":
				failed++
			case "RUNNING":
				running++
			case "QUEUED":
				queued++
			}
		}
		return "Runs " + itoa(len(m.store.Runs)) + " total" + "  " + m.styles.Dim.Render("•") + "  " + itoa(running) + " running" + "  " + m.styles.Dim.Render("•") + "  " + itoa(failed) + " failed" + "  " + m.styles.Dim.Render("•") + "  " + itoa(queued) + " queued"
	case ViewTriggers:
		active := 0
		for _, trg := range m.store.Triggers {
			if trg.Active {
				active++
			}
		}
		return "Triggers " + itoa(len(m.store.Triggers)) + " total" + "  " + m.styles.Dim.Render("•") + "  " + itoa(active) + " active"
	case ViewEvents:
		return "Events " + itoa(len(m.store.Events)) + " total"
	case ViewSecrets:
		return "Secrets " + itoa(len(m.store.Secrets)) + " total"
	case ViewTokens:
		revoked := 0
		for _, tok := range m.store.ApiTokens {
			if tok.Revoked {
				revoked++
			}
		}
		return "Tokens " + itoa(len(m.store.ApiTokens)) + " total" + "  " + m.styles.Dim.Render("•") + "  " + itoa(revoked) + " revoked"
	default:
		return ""
	}
}

func contextTitle(m Model) string {
	selected := m.selectedRowID()
	if selected == "" {
		return "Context Panel"
	}
	return "Context — " + selected
}
