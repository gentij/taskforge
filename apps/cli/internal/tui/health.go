package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentij/taskforge/apps/cli/internal/api"
)

type healthMsg struct {
	data *api.Health
	err  error
}

func fetchHealthCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return healthMsg{err: fmt.Errorf("missing API client")}
		}

		var health api.Health
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
