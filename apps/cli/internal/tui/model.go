package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentij/taskforge/apps/cli/internal/api"
)

type ViewID string

const (
	ViewHome      ViewID = "home"
	ViewWorkflows ViewID = "workflows"
)

type model struct {
	client      *api.Client
	serverURL   string
	tokenSet    bool
	width       int
	height      int
	spinner     spinner.Model
	loading     bool
	lastUpdated time.Time
	health      *api.Health
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
		headerHeight := measureHeaderHeight(m)
		footerHeight := measureFooterHeight(m)
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
