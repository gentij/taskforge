package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentij/taskforge/apps/cli/internal/api"
)

type App struct {
	client    *api.Client
	serverURL string
	tokenSet  bool
}

type ViewID string

const (
	ViewHome      ViewID = "home"
	ViewWorkflows ViewID = "workflows"
)

func NewApp(client *api.Client, serverURL string, tokenSet bool) *App {
	return &App{client: client, serverURL: serverURL, tokenSet: tokenSet}
}

func (a *App) Start() error {
	model := newModel(a.client, a.serverURL, a.tokenSet)
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

type healthResponse struct {
	Status  string  `json:"status"`
	Version string  `json:"version"`
	Uptime  float64 `json:"uptime"`
	DB      struct {
		Ok bool `json:"ok"`
	} `json:"db"`
}

type healthMsg struct {
	data *healthResponse
	err  error
}

type model struct {
	client      *api.Client
	serverURL   string
	tokenSet    bool
	width       int
	height      int
	spinner     spinner.Model
	loading     bool
	lastUpdated time.Time
	health      *healthResponse
	err         error
	errorMsg    string
	errorUntil  time.Time
	showHelp    bool
	currentView ViewID
	viewStack   []ViewID
	viewport    viewport.Model
}

func newModel(client *api.Client, serverURL string, tokenSet bool) model {
	spin := spinner.New(spinner.WithSpinner(spinner.Dot))
	vp := viewport.New(0, 0)
	return model{
		client:      client,
		serverURL:   serverURL,
		tokenSet:    tokenSet,
		spinner:     spin,
		loading:     true,
		currentView: ViewHome,
		viewStack:   []ViewID{},
		viewport:    vp,
	}
}

func (m model) pushView(next ViewID) model {
	if m.currentView != "" {
		m.viewStack = append(m.viewStack, m.currentView)
	}
	m.currentView = next
	return m
}

func (m model) popView() model {
	if len(m.viewStack) == 0 {
		return m
	}
	last := m.viewStack[len(m.viewStack)-1]
	m.viewStack = m.viewStack[:len(m.viewStack)-1]
	m.currentView = last
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHealthCmd(m.client))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = max(0, msg.Width-4)
		headerHeight := lipgloss.Height(renderHeader("Taskforge TUI", m.serverURL, m.tokenSet, viewTitle(m.currentView)))
		footerHeight := lipgloss.Height(renderFooter(m.lastUpdated))
		m.viewport.Height = max(0, msg.Height-headerHeight-footerHeight-6)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = true
			return m, nil
		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
		case "b":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			m = m.popView()
			return m, nil
		case "w":
			if !m.loading {
				m = m.pushView(ViewWorkflows)
			}
			return m, nil
		case "r":
			m.loading = true
			m.err = nil
			m.errorMsg = ""
			m.errorUntil = time.Time{}
			return m, tea.Batch(m.spinner.Tick, fetchHealthCmd(m.client))
		}

		if !m.showHelp {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case healthMsg:
		m.loading = false
		m.lastUpdated = time.Now()
		m.health = msg.data
		m.err = msg.err
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.errorUntil = time.Now().Add(5 * time.Second)
			return m, clearErrorCmd(m.errorUntil)
		}
		return m, nil
	case clearErrorMsg:
		if !m.errorUntil.IsZero() && time.Now().After(m.errorUntil) {
			m.errorMsg = ""
			m.errorUntil = time.Time{}
		}
		return m, nil
	}

	return m, nil
}

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

func fetchHealthCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return healthMsg{err: fmt.Errorf("missing API client")}
		}

		var health healthResponse
		err := client.GetJSON("/health", &health)
		if err != nil {
			return healthMsg{err: err}
		}
		return healthMsg{data: &health}
	}
}

type clearErrorMsg struct{}

func clearErrorCmd(until time.Time) tea.Cmd {
	return tea.Tick(time.Until(until), func(time.Time) tea.Msg {
		return clearErrorMsg{}
	})
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

	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("─", max(0, lipgloss.Width(info))))

	return lipgloss.JoinVertical(lipgloss.Left, line, info, divider)
}

func viewTitle(view ViewID) string {
	switch view {
	case ViewWorkflows:
		return "Workflows"
	default:
		return "Home"
	}
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
		Body:  "q quit\nr refresh\n? help\nesc close",
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
