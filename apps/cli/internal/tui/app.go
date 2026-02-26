package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentij/taskforge/apps/cli/internal/api"
)

type App struct {
	client    *api.Client
	serverURL string
	tokenSet  bool
}

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
}

func newModel(client *api.Client, serverURL string, tokenSet bool) model {
	spin := spinner.New(spinner.WithSpinner(spinner.Dot))
	return model{
		client:    client,
		serverURL: serverURL,
		tokenSet:  tokenSet,
		spinner:   spin,
		loading:   true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHealthCmd(m.client))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		case "r":
			m.loading = true
			m.err = nil
			m.errorMsg = ""
			m.errorUntil = time.Time{}
			return m, tea.Batch(m.spinner.Tick, fetchHealthCmd(m.client))
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
	header := renderHeader("Taskforge TUI")
	serverLine := fmt.Sprintf("Server: %s", m.serverURL)
	tokenLine := "Token: missing"
	if m.tokenSet {
		tokenLine = "Token: set"
	}

	statusBlock := renderStatus(serverLine, tokenLine)

	body := ""
	if m.loading {
		body = fmt.Sprintf("%s Loading health...", m.spinner.View())
	} else if m.health != nil {
		dbStatus := "down"
		if m.health.DB.Ok {
			dbStatus = "ok"
		}
		body = fmt.Sprintf(
			"Status: %s\nVersion: %s\nDB: %s",
			m.health.Status,
			m.health.Version,
			dbStatus,
		)
	}

	footer := renderFooter(m.lastUpdated)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		statusBlock,
		"",
		body,
	)

	if m.errorMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", renderError(m.errorMsg))
	}

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

func renderHeader(title string) string {
	return lipgloss.NewStyle().Bold(true).Render(title)
}

func renderStatus(serverLine string, tokenLine string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
		fmt.Sprintf("%s\n%s", serverLine, tokenLine),
	)
}

func renderFooter(lastUpdated time.Time) string {
	text := "q quit • r refresh • ? help"
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
