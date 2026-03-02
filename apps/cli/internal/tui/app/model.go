package app

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/tui/components"
	"github.com/gentij/taskforge/apps/cli/internal/tui/data"
	"github.com/gentij/taskforge/apps/cli/internal/tui/layout"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
	"github.com/gentij/taskforge/apps/cli/internal/tui/utils"
	"golang.org/x/term"
)

type ViewID string

const (
	ViewDashboard ViewID = "dashboard"
	ViewWorkflows ViewID = "workflows"
	ViewRuns      ViewID = "runs"
	ViewTriggers  ViewID = "triggers"
	ViewEvents    ViewID = "events"
	ViewSecrets   ViewID = "secrets"
	ViewTokens    ViewID = "tokens"
)

type paletteActionType int

const (
	paletteGoToView paletteActionType = iota
	paletteToggleRefresh
	paletteClearFilters
	paletteRunWorkflow
)

type paletteAction struct {
	Kind paletteActionType
	View ViewID
}

type paletteItem struct {
	Title  string
	Desc   string
	Action paletteAction
}

func (p paletteItem) FilterValue() string { return p.Title }

type Model struct {
	client    *api.Client
	serverURL string
	tokenSet  bool

	width  int
	height int
	layout layout.Layout

	view  ViewID
	views []ViewID

	theme  styles.Theme
	styles styles.StyleSet

	store data.Store

	columns        []table.Column
	baseRows       []table.Row
	baseRowIDs     []string
	filteredRows   []table.Row
	filteredRowIDs []string

	table table.Model

	searchInput textinput.Model
	searching   bool
	searchQuery string

	contextViewport    viewport.Model
	contextCollapsed   bool
	contextSearchInput textinput.Model
	contextSearching   bool
	contextQuery       string

	palette     list.Model
	showPalette bool

	inspector RunInspector

	showHelp bool
	help     help.Model
	keys     KeyMap

	autoRefresh   bool
	refreshEvery  time.Duration
	lastRefresh   time.Time
	workspaceName string
	workerCount   int
	apiStatus     string

	paginator paginator.Model
}

func NewModel(client *api.Client, serverURL string, tokenSet bool) Model {
	now := time.Now()
	store := data.MockStore(now)
	keys := DefaultKeyMap()
	helper := help.New()
	helper.ShowAll = false

	search := textinput.New()
	search.Prompt = "/ "
	search.Placeholder = "Filter"
	search.CharLimit = 64

	contextSearch := textinput.New()
	contextSearch.Prompt = "panel/ "
	contextSearch.Placeholder = "Search context"
	contextSearch.CharLimit = 64

	palette := buildPalette()

	theme := styles.DefaultTheme()
	styleSet := styles.NewStyles(theme)

	tableModel := components.NewTable(nil, nil, 0, 0, styleSet)
	contextViewport := viewport.New(0, 0)

	pager := paginator.New()
	pager.Type = paginator.Arabic
	pager.PerPage = 1
	pager.SetTotalPages(1)

	model := Model{
		client:             client,
		serverURL:          serverURL,
		tokenSet:           tokenSet,
		view:               ViewDashboard,
		views:              []ViewID{ViewDashboard, ViewWorkflows, ViewRuns, ViewTriggers, ViewEvents, ViewSecrets, ViewTokens},
		theme:              theme,
		styles:             styleSet,
		store:              store,
		table:              tableModel,
		searchInput:        search,
		contextViewport:    contextViewport,
		contextSearchInput: contextSearch,
		palette:            palette,
		help:               helper,
		keys:               keys,
		refreshEvery:       2 * time.Second,
		lastRefresh:        now,
		workspaceName:      "personal",
		workerCount:        2,
		apiStatus:          apiStatus(tokenSet),
		paginator:          pager,
		inspector:          NewInspector(styleSet, keys),
	}

	width, height := initialSize()
	model.resize(width, height)
	model.refreshView()

	return model
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		width, height := initialSize()
		return tea.WindowSizeMsg{Width: width, Height: height}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inspector.Active {
		return m.updateInspector(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		m.refreshView()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.searching {
		return m.updateSearch(msg)
	}
	if m.contextSearching {
		return m.updateContextSearch(msg)
	}
	if m.showPalette {
		return m.updatePalette(msg)
	}

	return m.updateTable(msg)
}

func (m Model) View() string {
	return Render(m)
}

func (m *Model) resize(width int, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	m.width = width
	m.height = height
	m.layout = layout.Compute(width, height, m.contextCollapsed)

	m.table.SetWidth(width)
	m.table.SetHeight(m.layout.PrimaryTableHeight)

	m.paginator.PerPage = m.layout.PrimaryTableHeight

	contextBodyHeight := max(m.layout.ContextHeight-3, 1)
	contextWidth := max(width-2, 1)
	if m.layout.ContextHeight == 0 {
		contextBodyHeight = 0
	}
	m.contextViewport.Width = contextWidth
	m.contextViewport.Height = contextBodyHeight

	m.inspector.Resize(width, height)
}

func (m *Model) refreshView() {
	columns, rows, rowIDs := BuildRowsForView(m.view, &m.store, m.styles, m.width)
	m.columns = columns
	m.baseRows = rows
	m.baseRowIDs = rowIDs
	m.applyFilter()
}

func (m *Model) applyFilter() {
	rows, rowIDs := filterRows(m.baseRows, m.baseRowIDs, m.searchQuery)
	m.filteredRows = rows
	m.filteredRowIDs = rowIDs
	m.applyTableRows()
}

func (m *Model) applyTableRows() {
	m.table.SetRows(nil)
	m.table.SetColumns(m.columns)
	if len(m.filteredRows) == 0 {
		m.updatePaginator()
		m.updateContext()
		return
	}

	truncated := truncateRows(m.filteredRows, m.columns)
	m.filteredRows = truncated

	cursor := m.table.Cursor()
	if cursor >= len(truncated) {
		cursor = len(truncated) - 1
		m.table.SetCursor(cursor)
	}
	if cursor < 0 {
		cursor = 0
		m.table.SetCursor(0)
	}
	styled := components.StyleRows(truncated, cursor, m.styles)
	m.table.SetRows(styled)
	m.updatePaginator()
	m.updateContext()
}

func (m *Model) updatePaginator() {
	perPage := max(m.layout.PrimaryTableHeight, 1)
	pages := (len(m.filteredRows) + perPage - 1) / perPage
	if pages < 1 {
		pages = 1
	}
	m.paginator.SetTotalPages(pages)
	page := 0
	if len(m.filteredRows) > 0 {
		page = m.table.Cursor() / perPage
	}
	if page >= pages {
		page = pages - 1
	}
	m.paginator.Page = page
}

func (m *Model) updateContext() {
	selectedID := m.selectedRowID()
	content := BuildContextContent(m.view, &m.store, selectedID)
	content = utils.FilterLines(content, m.contextQuery)
	if m.contextViewport.Width > 0 {
		content = utils.WrapText(content, m.contextViewport.Width)
	}
	m.contextViewport.SetContent(content)
	m.contextViewport.GotoTop()
}

func (m *Model) selectedRowID() string {
	if len(m.filteredRowIDs) == 0 {
		return ""
	}
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.filteredRowIDs) {
		return ""
	}
	return m.filteredRowIDs[idx]
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showHelp {
		if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Back) {
			m.showHelp = false
			m.help.ShowAll = false
			return m, nil
		}
		return m, nil
	}
	if m.showPalette {
		return m.updatePalette(msg)
	}
	if m.searching {
		return m.updateSearch(msg)
	}
	if m.contextSearching {
		return m.updateContextSearch(msg)
	}

	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit
	}
	if key.Matches(msg, m.keys.Help) {
		m.showHelp = true
		m.help.ShowAll = true
		return m, nil
	}
	if key.Matches(msg, m.keys.Palette) {
		m.showPalette = true
		return m, nil
	}
	if key.Matches(msg, m.keys.Search) {
		m.searching = true
		m.searchInput.SetValue(m.searchQuery)
		m.searchInput.CursorEnd()
		m.searchInput.Focus()
		return m, nil
	}
	if key.Matches(msg, m.keys.ContextSearch) {
		m.contextSearching = true
		m.contextSearchInput.SetValue(m.contextQuery)
		m.contextSearchInput.CursorEnd()
		m.contextSearchInput.Focus()
		return m, nil
	}
	if key.Matches(msg, m.keys.ToggleContext) {
		m.contextCollapsed = !m.contextCollapsed
		m.resize(m.width, m.height)
		m.updateContext()
		return m, nil
	}
	if key.Matches(msg, m.keys.NextScreen) {
		m.nextView()
		return m, nil
	}
	if key.Matches(msg, m.keys.PrevScreen) {
		m.prevView()
		return m, nil
	}
	if m.view == ViewRuns && key.Matches(msg, m.keys.Enter) {
		m.openInspector()
		return m, nil
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.RunWorkflow) {
		m.queueRunForSelectedWorkflow()
		return m, nil
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.ToggleActive) {
		m.toggleWorkflowActive()
		return m, nil
	}
	if m.view == ViewTokens && key.Matches(msg, m.keys.RevokeToken) {
		m.toggleTokenRevoked()
		return m, nil
	}

	return m.updateTable(msg)
}

func (m *Model) updateTable(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if m.handleContextScroll(keyMsg) {
			return m, nil
		}
	}

	prevCursor := m.table.Cursor()
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	if prevCursor != m.table.Cursor() {
		m.applyTableRows()
	}
	return m, cmd
}

func (m *Model) handleContextScroll(msg tea.KeyMsg) bool {
	if m.contextViewport.Height <= 0 {
		return false
	}
	switch msg.String() {
	case "alt+up", "alt+k", "shift+up", "shift+k":
		m.contextViewport.LineUp(1)
		return true
	case "alt+down", "alt+j", "shift+down", "shift+j":
		m.contextViewport.LineDown(1)
		return true
	case "pgup", "pageup":
		m.contextViewport.LineUp(m.contextViewport.Height)
		return true
	case "pgdown", "pagedown", "pgdn":
		m.contextViewport.LineDown(m.contextViewport.Height)
		return true
	case "home":
		m.contextViewport.GotoTop()
		return true
	case "end":
		m.contextViewport.GotoBottom()
		return true
	}
	return false
}

func (m *Model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) {
			m.searching = false
			m.searchInput.Blur()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) {
			m.searching = false
			m.searchInput.Blur()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Clear) {
			m.searchInput.SetValue("")
			m.searchQuery = ""
			m.applyFilter()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.searchQuery = strings.TrimSpace(m.searchInput.Value())
	m.applyFilter()
	return m, cmd
}

func (m *Model) updateContextSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) {
			m.contextSearching = false
			m.contextSearchInput.Blur()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) {
			m.contextSearching = false
			m.contextSearchInput.Blur()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Clear) {
			m.contextSearchInput.SetValue("")
			m.contextQuery = ""
			m.updateContext()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.contextSearchInput, cmd = m.contextSearchInput.Update(msg)
	m.contextQuery = strings.TrimSpace(m.contextSearchInput.Value())
	m.updateContext()
	return m, cmd
}

func (m *Model) updatePalette(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) {
			m.showPalette = false
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) {
			item, ok := m.palette.SelectedItem().(paletteItem)
			if ok {
				m.runPaletteAction(item.Action)
			}
			m.showPalette = false
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.palette, cmd = m.palette.Update(msg)
	return m, cmd
}

func (m *Model) runPaletteAction(action paletteAction) {
	switch action.Kind {
	case paletteGoToView:
		m.view = action.View
		m.refreshView()
	case paletteToggleRefresh:
		m.autoRefresh = !m.autoRefresh
	case paletteClearFilters:
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.applyFilter()
	case paletteRunWorkflow:
		m.queueRunForSelectedWorkflow()
	}
}

func (m *Model) nextView() {
	for i, view := range m.views {
		if view == m.view {
			m.view = m.views[(i+1)%len(m.views)]
			m.refreshView()
			return
		}
	}
}

func (m *Model) prevView() {
	for i, view := range m.views {
		if view == m.view {
			idx := i - 1
			if idx < 0 {
				idx = len(m.views) - 1
			}
			m.view = m.views[idx]
			m.refreshView()
			return
		}
	}
}

func (m *Model) toggleWorkflowActive() {
	selected := m.selectedRowID()
	if selected == "" {
		return
	}
	for i, wf := range m.store.Workflows {
		if wf.ID == selected {
			m.store.Workflows[i].Active = !wf.Active
			m.refreshView()
			return
		}
	}
}

func (m *Model) toggleTokenRevoked() {
	selected := m.selectedRowID()
	if selected == "" {
		return
	}
	for i, tok := range m.store.ApiTokens {
		if tok.ID == selected {
			m.store.ApiTokens[i].Revoked = !tok.Revoked
			m.refreshView()
			return
		}
	}
}

func (m *Model) queueRunForSelectedWorkflow() {
	if m.view != ViewWorkflows {
		return
	}
	selected := m.selectedRowID()
	if selected == "" {
		return
	}
	newRun := data.WorkflowRun{
		ID:          "run_" + time.Now().Format("150405"),
		WorkflowID:  selected,
		Status:      "QUEUED",
		TriggerType: "manual",
		StartedAt:   time.Now(),
		Duration:    0,
		InputJSON:   `{"manual":true}`,
		OutputJSON:  `{}`,
	}
	m.store.Runs = append([]data.WorkflowRun{newRun}, m.store.Runs...)
	m.refreshView()
}

func buildPalette() list.Model {
	items := []list.Item{
		paletteItem{Title: "Go to Dashboard", Desc: "Overview", Action: paletteAction{Kind: paletteGoToView, View: ViewDashboard}},
		paletteItem{Title: "Go to Workflows", Desc: "Workflow list", Action: paletteAction{Kind: paletteGoToView, View: ViewWorkflows}},
		paletteItem{Title: "Go to Runs", Desc: "Workflow runs", Action: paletteAction{Kind: paletteGoToView, View: ViewRuns}},
		paletteItem{Title: "Go to Triggers", Desc: "Trigger list", Action: paletteAction{Kind: paletteGoToView, View: ViewTriggers}},
		paletteItem{Title: "Go to Events", Desc: "Event list", Action: paletteAction{Kind: paletteGoToView, View: ViewEvents}},
		paletteItem{Title: "Go to Secrets", Desc: "Secret registry", Action: paletteAction{Kind: paletteGoToView, View: ViewSecrets}},
		paletteItem{Title: "Go to API Tokens", Desc: "Token list", Action: paletteAction{Kind: paletteGoToView, View: ViewTokens}},
		paletteItem{Title: "Run workflow", Desc: "Queue selected workflow", Action: paletteAction{Kind: paletteRunWorkflow}},
		paletteItem{Title: "Toggle auto refresh", Desc: "Enable/disable polling", Action: paletteAction{Kind: paletteToggleRefresh}},
		paletteItem{Title: "Clear filters", Desc: "Reset table filters", Action: paletteAction{Kind: paletteClearFilters}},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetHeight(2)
	model := list.New(items, delegate, 60, 12)
	model.SetFilteringEnabled(true)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)
	model.SetShowTitle(false)
	model.DisableQuitKeybindings()
	return model
}

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
			rows[i][c] = utils.Truncate(rows[i][c], widths[c])
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
