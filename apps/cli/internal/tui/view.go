package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	header := renderHeader("Taskforge TUI", m.serverURL, m.tokenSet, viewTitle(m.currentView))
	body := renderBody(m)

	viewContent := m.viewport
	viewContent.SetContent(body)
	content := lipgloss.JoinVertical(lipgloss.Left, header, "", viewContent.View())

	if m.errorMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", renderError(m.errorMsg))
	}

	footer := renderFooter(m.lastUpdated)
	content = lipgloss.JoinVertical(lipgloss.Left, content, "", footer)
	content = lipgloss.NewStyle().Padding(1, 2).Render(content)

	if m.showHelp {
		return renderModal(content, helpModal(), m.width, m.height)
	}

	return content
}

func renderHeader(title string, serverURL string, tokenSet bool, viewName string) string {
	left := lipgloss.NewStyle().Bold(true).Render(title)
	right := lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(viewName)
	line := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	tokenLine := "Token: missing"
	if tokenSet {
		tokenLine = "Token: set"
	}
	info := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
		fmt.Sprintf("Server: %s\n%s", serverURL, tokenLine),
	)

	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(
		strings.Repeat("─", max(0, lipgloss.Width(info))),
	)

	return lipgloss.JoinVertical(lipgloss.Left, line, info, divider)
}

func renderFooter(lastUpdated time.Time) string {
	text := "q quit • r refresh • w workflows • b back • ? help"
	if !lastUpdated.IsZero() {
		text = fmt.Sprintf("%s • updated %s", text, lastUpdated.Format("15:04:05"))
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(text)
}

func renderError(message string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	return style.Render("Error: " + message)
}

type modal struct {
	Title string
	Body  string
}

func helpModal() modal {
	return modal{
		Title: "Help",
		Body:  "q quit\nr refresh\nw workflows\nb back\n? help\nesc close",
	}
}

func renderModal(content string, m modal, width int, height int) string {
	modalStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Background(lipgloss.Color("235"))

	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.Title),
		"",
		m.Body,
	)

	box := modalStyle.Render(modalContent)
	canvasWidth := width
	canvasHeight := height
	if canvasWidth == 0 {
		canvasWidth = lipgloss.Width(content)
	}
	if canvasHeight == 0 {
		canvasHeight = lipgloss.Height(content)
	}

	return lipgloss.Place(canvasWidth, canvasHeight, lipgloss.Center, lipgloss.Center, box)
}

func viewTitle(view ViewID) string {
	switch view {
	case ViewWorkflows:
		return "Workflows"
	default:
		return "Home"
	}
}

func renderBody(m model) string {
	breadcrumb := renderBreadcrumb(m.currentView)
	switch m.currentView {
	case ViewWorkflows:
		return lipgloss.JoinVertical(lipgloss.Left, breadcrumb, "", renderWorkflowsView())
	default:
		return lipgloss.JoinVertical(lipgloss.Left, breadcrumb, "", renderHomeView(m))
	}
}

func renderHomeView(m model) string {
	if m.loading {
		return fmt.Sprintf("%s Loading health...", m.spinner.View())
	}
	if m.health == nil {
		return "No health data yet."
	}
	dbStatus := "down"
	if m.health.DB.Ok {
		dbStatus = "ok"
	}
	return fmt.Sprintf(
		"Status: %s\nVersion: %s\nDB: %s",
		m.health.Status,
		m.health.Version,
		dbStatus,
	)
}

func renderWorkflowsView() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(
		"Workflows view (coming soon)",
	)
}

func renderBreadcrumb(view ViewID) string {
	trail := "Home"
	if view == ViewWorkflows {
		trail = "Home / Workflows"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(trail)
}

func measureHeaderHeight(m model) int {
	return lipgloss.Height(renderHeader("Taskforge TUI", m.serverURL, m.tokenSet, viewTitle(m.currentView)))
}

func measureFooterHeight(m model) int {
	return lipgloss.Height(renderFooter(m.lastUpdated))
}
