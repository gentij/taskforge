package app

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gentij/taskforge/apps/cli/internal/tui/components"
)

func Render(m Model) string {
	if m.inspector.Active {
		return m.inspector.Render(m.width, m.height)
	}

	sidebar := renderSidebar(m)
	mainPanel := renderMainPanel(m)
	footer := renderFooter(m)

	row := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainPanel)
	row = clampSection(row, m.width, m.layout.MainHeight)
	base := lipgloss.JoinVertical(lipgloss.Left, row, footer)
	base = clampToViewport(base, m.width, m.height)

	output := base
	if m.showPalette {
		output = renderPaletteScreen(m)
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

func buildMainContent(m Model) string {
	width := max(m.layout.MainWidth-2, 1)
	header := renderMainHeader(m, width)
	body := renderMainBody(m, width)
	divider := m.styles.Divider.Render(strings.Repeat("─", width))
	sections := []string{header, divider, body}
	return strings.Join(sections, "\n")
}

func renderSidebar(m Model) string {
	width := m.layout.SidebarWidth
	height := m.layout.SidebarHeight
	innerWidth := max(width-2, 1)
	innerHeight := max(height-2, 1)
	border := m.styles.PanelBorder
	if m.focus == FocusSidebar {
		border = m.styles.PanelBorderFocus
	}

	brand := m.styles.SidebarTitle.Render("TASKFORGE")
	workspace := m.styles.SidebarMuted.Render("WS: " + m.workspaceName)
	section := m.styles.SidebarSection.Render("Navigation")
	listView := strings.TrimRight(m.sidebar.View(), "\n")

	status := []string{
		m.styles.SidebarSection.Render("Status"),
		chip(m, "API "+m.apiStatus, m.apiStatus == "CONNECTED"),
		chip(m, "Workers "+itoa(m.workerCount), false),
		chip(m, refreshChip(m), false),
		chip(m, "Theme "+strings.Title(strings.ReplaceAll(m.themeName, "-", " ")), false),
	}

	content := joinSidebarContent(
		[]string{brand, workspace, "", section},
		strings.Split(listView, "\n"),
		status,
		innerHeight,
	)
	filled := applyBackgroundLayer(content, innerWidth, innerHeight, m.styles.SidebarFill)
	body := lipgloss.Place(innerWidth, innerHeight, lipgloss.Left, lipgloss.Top, filled)
	return border.Width(innerWidth).Height(innerHeight).Render(body)
}

func joinSidebarContent(top []string, middle []string, bottom []string, height int) string {
	lines := append([]string{}, top...)
	lines = append(lines, middle...)
	remaining := height - len(lines) - len(bottom)
	if remaining < 0 {
		remaining = 0
	}
	for i := 0; i < remaining; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, bottom...)
	return strings.Join(lines, "\n")
}

func renderMainPanel(m Model) string {
	width := m.layout.MainWidth
	height := m.layout.MainHeight
	innerWidth := max(width-2, 1)
	innerHeight := max(height-2, 1)
	border := m.styles.PanelBorder
	if m.focus == FocusMain {
		border = m.styles.PanelBorderFocus
	}
	panelContent := strings.TrimRight(m.mainPanel.View(), "\n")
	topHeight := innerHeight
	if !m.contextCollapsed && m.layout.ContextHeight > 0 {
		topHeight -= m.layout.ContextHeight
	}
	if topHeight < 1 {
		topHeight = 1
	}
	panelContent = clampSection(panelContent, innerWidth, topHeight)
	top := applyBackgroundLayer(panelContent, innerWidth, topHeight, m.styles.PanelFill)
	contentParts := []string{top}
	if !m.contextCollapsed && m.layout.ContextHeight > 0 {
		contentParts = append(contentParts, renderContextDrawer(m, innerWidth))
	}
	composed := strings.Join(contentParts, "\n")
	composed = clampSection(composed, innerWidth, innerHeight)
	content := lipgloss.Place(innerWidth, innerHeight, lipgloss.Left, lipgloss.Top, composed)
	return border.Width(innerWidth).Height(innerHeight).Render(content)
}

func renderMainHeader(m Model, width int) string {
	section := viewTitle(m.view)
	summary := viewSummary(m)
	left := m.styles.PanelTitle.Render(section) + "  " + m.styles.Dim.Render("•") + "  " + summary
	filter := ""
	if m.searchQuery != "" {
		filter = "Filter: " + m.searchQuery
	}
	line1 := joinLeftRight(left, filter, width)
	chips := []string{
		chip(m, "API "+m.apiStatus, m.apiStatus == "CONNECTED"),
		chip(m, "Workers "+itoa(m.workerCount), false),
		chip(m, refreshChip(m), false),
	}
	if m.searchQuery != "" {
		chips = append(chips, chip(m, "Filter", true))
	}
	line2 := strings.Join(chips, " ")

	line1 = ansi.Truncate(line1, width, "")
	line2 = ansi.Truncate(line2, width, "")
	line1 = m.styles.Header.Width(width).Render(line1)
	line2 = m.styles.Header.Width(width).Render(line2)
	return strings.Join([]string{line1, line2}, "\n")
}

func renderMainBody(m Model, width int) string {
	parts := []string{}
	if m.view == ViewDashboard {
		parts = append(parts, renderDashboard(m, width))
	} else {
		label := m.styles.SidebarSection.Render("List")
		if len(m.filteredRows) == 0 {
			parts = append(parts, label, renderEmptyState(m, width))
		} else {
			parts = append(parts, label, strings.TrimRight(m.table.View(), "\n"))
		}
	}
	return strings.Join(parts, "\n")
}

func renderDashboard(m Model, width int) string {
	if m.layout.DashboardCardsHeight == 0 {
		if len(m.filteredRows) == 0 {
			return renderEmptyState(m, width)
		}
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
	}, width, m.styles)
	label := m.styles.SidebarSection.Render("Recent Runs")
	if len(m.filteredRows) == 0 {
		return strings.Join([]string{cards, label, renderEmptyState(m, width)}, "\n")
	}
	return strings.Join([]string{cards, label, strings.TrimRight(m.table.View(), "\n")}, "\n")
}

func renderContextDrawer(m Model, width int) string {
	innerWidth := max(width-2, 1)
	tabs := renderTabs(m, []string{"Overview", "JSON", "Steps", "Logs"}, contextTabLabel(m.contextTab))
	content := strings.TrimRight(m.contextViewport.View(), "\n")
	body := lipgloss.Place(innerWidth, max(m.layout.ContextHeight-4, 1), lipgloss.Left, lipgloss.Top, content)
	inner := lipgloss.JoinVertical(lipgloss.Left, tabs, body)
	inner = applyBackgroundLayer(inner, innerWidth, max(m.layout.ContextHeight-2, 1), m.styles.ContextFill)
	box := m.styles.PanelBorder.Width(innerWidth).Height(max(m.layout.ContextHeight-2, 1))
	if m.focus == FocusContext {
		box = m.styles.PanelBorderFocus.Width(innerWidth).Height(max(m.layout.ContextHeight-2, 1))
	}
	return box.Render(inner)
}

func renderTabs(m Model, tabs []string, active string) string {
	items := make([]string, 0, len(tabs))
	for _, tab := range tabs {
		if tab == active {
			items = append(items, m.styles.TabActive.Render(" "+tab+" "))
		} else {
			items = append(items, m.styles.TabInactive.Render(" "+tab+" "))
		}
	}
	return strings.Join(items, " ")
}

func contextTabLabel(tab ContextTab) string {
	switch tab {
	case ContextTabJSON:
		return "JSON"
	case ContextTabSteps:
		return "Steps"
	case ContextTabLogs:
		return "Logs"
	default:
		return "Overview"
	}
}

func renderFooter(m Model) string {
	left := m.help.View(m.keys)
	right := m.paginator.View()
	line := joinLeftRight(left, right, m.width)
	content := m.styles.Footer.Width(m.width).Render(line)
	return clampSection(content, m.width, m.layout.FooterHeight)
}

func chip(m Model, text string, active bool) string {
	if active {
		return m.styles.ChipActive.Render(" " + text + " ")
	}
	return m.styles.Chip.Render(" " + text + " ")
}

func refreshLabel(m Model) string {
	if m.autoRefresh {
		return "2s"
	}
	return "OFF"
}

func refreshChip(m Model) string {
	if !m.autoRefresh {
		return "⟳ " + refreshLabel(m)
	}
	if m.pulseOn {
		return "↻ " + refreshLabel(m)
	}
	return "⟳ " + refreshLabel(m)
}

func renderEmptyState(m Model, width int) string {
	title := "No items to display"
	if m.searchQuery != "" {
		title = "No matches for \"" + m.searchQuery + "\""
	}
	hint := "Try a different filter, or open the command palette with ctrl+k"
	line1 := m.styles.Dim.Bold(true).Render(title)
	line2 := m.styles.Dim.Render(hint)
	inner := strings.Join([]string{line1, "", line2}, "\n")
	return applyBackgroundLayer(inner, width, 6, m.styles.ContextFill)
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

func renderPaletteScreen(m Model) string {
	innerWidth := max(m.width-2, 1)
	innerHeight := max(m.height-2, 1)
	headerTitle := m.styles.PanelTitle.Render("Command Palette")
	headerHint := m.styles.Dim.Render("Type to filter  |  / manual filter  |  enter select  |  esc close")
	divider := m.styles.Divider.Render(strings.Repeat("─", innerWidth))

	listHeight := max(innerHeight-3, 1)
	paletteView := strings.TrimRight(m.palette.View(), "\n")
	paletteView = clampSection(paletteView, innerWidth, listHeight)

	body := strings.Join([]string{headerTitle, headerHint, divider, paletteView}, "\n")
	filled := applyBackgroundLayer(body, innerWidth, innerHeight, m.styles.PanelFill)
	content := lipgloss.Place(innerWidth, innerHeight, lipgloss.Left, lipgloss.Top, filled)
	panel := m.styles.PanelBorderFocus.Width(innerWidth).Height(innerHeight).Render(content)
	return clampToViewport(panel, m.width, m.height)
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

func truncateLines(content string, width int) string {
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], width, "")
	}
	return strings.Join(lines, "\n")
}

func sanitizeRenderable(content string) string {
	content = strings.ReplaceAll(content, "\r", "")
	return stripNonSGRANSI(content)
}

func stripNonSGRANSI(content string) string {
	var builder strings.Builder
	builder.Grow(len(content))
	for i := 0; i < len(content); i++ {
		ch := content[i]
		if ch != '\x1b' {
			builder.WriteByte(ch)
			continue
		}
		if i+1 >= len(content) {
			break
		}
		next := content[i+1]
		if next != '[' {
			continue
		}
		j := i + 2
		for j < len(content) {
			final := content[j]
			if final >= 0x40 && final <= 0x7E {
				if final == 'm' {
					builder.WriteString(content[i : j+1])
				}
				i = j
				break
			}
			j++
		}
		if j >= len(content) {
			break
		}
	}
	return builder.String()
}

func applyBackgroundLayer(content string, width int, height int, style lipgloss.Style) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	prefix, reset := backgroundCodes(style)
	lines := strings.Split(content, "\n")
	filled := make([]string, 0, height)
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = ansi.Truncate(lines[i], width, "")
		}
		if prefix != "" {
			line = strings.ReplaceAll(line, reset, reset+prefix)
			line = strings.ReplaceAll(line, "\x1b[m", "\x1b[m"+prefix)
			line = strings.ReplaceAll(line, "\x1b[49m", "\x1b[49m"+prefix)
			line = strings.ReplaceAll(line, "\x1b[39m", "\x1b[39m"+prefix)
		}
		pad := width - ansi.StringWidth(line)
		if pad < 0 {
			pad = 0
		}
		line = prefix + line + strings.Repeat(" ", pad) + reset
		filled = append(filled, line)
	}
	return strings.Join(filled, "\n")
}

func backgroundCodes(style lipgloss.Style) (string, string) {
	sample := style.Render("X")
	idx := strings.Index(sample, "X")
	if idx == -1 {
		return "", "\x1b[0m"
	}
	prefix := sample[:idx]
	suffix := sample[idx+1:]
	if suffix == "" {
		suffix = "\x1b[0m"
	}
	return prefix, suffix
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
