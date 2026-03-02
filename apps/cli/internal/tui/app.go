package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/tui/app"
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
	model := app.NewModel(a.client, a.serverURL, a.tokenSet)
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}
