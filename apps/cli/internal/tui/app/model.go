package app

import (
	"fmt"
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
	"github.com/charmbracelet/lipgloss"
	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/config"
	"github.com/gentij/taskforge/apps/cli/internal/tui/components"
	"github.com/gentij/taskforge/apps/cli/internal/tui/data"
	"github.com/gentij/taskforge/apps/cli/internal/tui/layout"
	"github.com/gentij/taskforge/apps/cli/internal/tui/styles"
	"github.com/gentij/taskforge/apps/cli/internal/tui/utils"
	"golang.org/x/term"
)

type ViewID string

type FocusPane int

const (
	ViewDashboard ViewID = "dashboard"
	ViewWorkflows ViewID = "workflows"
	ViewRuns      ViewID = "runs"
	ViewTriggers  ViewID = "triggers"
	ViewEvents    ViewID = "events"
	ViewSecrets   ViewID = "secrets"
	ViewTokens    ViewID = "tokens"
)

const (
	FocusSidebar FocusPane = iota
	FocusMain
	FocusContext
)

type paletteActionType int

const (
	paletteNoop paletteActionType = iota
	paletteGoToView
	paletteToggleRefresh
	paletteClearFilters
	paletteRunWorkflow
	paletteSetTheme
)

type pulseMsg struct{}

type uiLoadDoneMsg struct{}

type mockRefreshDoneMsg struct {
	failed bool
}

type toastClearMsg struct {
	id int
}

type SurfaceState int

const (
	SurfaceIdle SurfaceState = iota
	SurfaceLoading
	SurfaceRefreshing
	SurfaceSuccess
	SurfaceError
	SurfaceStale
	SurfaceEmpty
)

type ToastLevel int

const (
	ToastInfo ToastLevel = iota
	ToastSuccess
	ToastWarn
	ToastError
)

type ToastState struct {
	ID      int
	Active  bool
	Level   ToastLevel
	Message string
}

type paletteAction struct {
	Kind paletteActionType
	View ViewID
}

type ContextTab int

const (
	ContextTabOverview ContextTab = iota
	ContextTabJSON
	ContextTabSteps
	ContextTabLogs
)

type paletteItem struct {
	Label   string
	Detail  string
	Action  paletteAction
	Section bool
}

type SortConfig struct {
	Column int
	Desc   bool
}

func (p paletteItem) FilterValue() string {
	if p.Section {
		return ""
	}
	return p.Label
}
func (p paletteItem) Title() string       { return p.Label }
func (p paletteItem) Description() string { return p.Detail }

type navItem struct {
	ID    ViewID
	Label string
}

func (n navItem) FilterValue() string { return n.Label }
func (n navItem) Title() string       { return n.Label }
func (n navItem) Description() string { return "" }

type Model struct {
	client     *api.Client
	serverURL  string
	tokenSet   bool
	config     config.Config
	configPath string

	focus FocusPane

	width  int
	height int
	layout layout.Layout

	view  ViewID
	views []ViewID

	theme     styles.Theme
	themeName string
	styles    styles.StyleSet

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
	contextTab         ContextTab
	contextOffsets     map[ContextTab]int
	contextSelectedID  string
	contextSearchInput textinput.Model
	contextSearching   bool
	contextQuery       string

	palette     list.Model
	showPalette bool
	sidebar     list.Model
	mainPanel   viewport.Model

	inspector RunInspector

	showHelp bool
	help     help.Model
	keys     KeyMap

	autoRefresh    bool
	refreshEvery   time.Duration
	lastRefresh    time.Time
	pulseOn        bool
	refreshPending bool
	refreshCount   int
	workspaceName  string
	workerCount    int
	apiStatus      string
	paletteRecent  []paletteAction

	mainState    SurfaceState
	contextState SurfaceState
	toast        ToastState
	uiReady      bool

	sortByView map[ViewID]SortConfig
	sortColumn int
	sortDesc   bool

	paginator paginator.Model
}

func NewModel(client *api.Client, serverURL string, tokenSet bool, cfg config.Config, configPath string) Model {
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

	defaultTheme := styles.DefaultTheme()
	palette := buildPalette(defaultTheme, nil)
	sidebar := buildSidebar(defaultTheme, ViewDashboard)
	styleSet := styles.NewStyles(defaultTheme)

	tableModel := components.NewTable(nil, nil, 0, 0, styleSet)
	contextViewport := viewport.New(0, 0)
	mainPanel := viewport.New(0, 0)

	pager := paginator.New()
	pager.Type = paginator.Arabic
	pager.PerPage = 1
	pager.SetTotalPages(1)

	model := Model{
		client:             client,
		serverURL:          serverURL,
		tokenSet:           tokenSet,
		config:             cfg,
		configPath:         configPath,
		focus:              FocusSidebar,
		view:               ViewDashboard,
		views:              []ViewID{ViewDashboard, ViewWorkflows, ViewRuns, ViewTriggers, ViewEvents, ViewSecrets, ViewTokens},
		theme:              defaultTheme,
		themeName:          "taskforge",
		styles:             styleSet,
		store:              store,
		table:              tableModel,
		searchInput:        search,
		contextViewport:    contextViewport,
		contextTab:         ContextTabOverview,
		contextOffsets:     map[ContextTab]int{},
		contextSearchInput: contextSearch,
		palette:            palette,
		sidebar:            sidebar,
		mainPanel:          mainPanel,
		help:               helper,
		keys:               keys,
		refreshEvery:       2 * time.Second,
		lastRefresh:        now,
		pulseOn:            false,
		refreshPending:     false,
		refreshCount:       0,
		workspaceName:      "personal",
		workerCount:        2,
		apiStatus:          apiStatus(tokenSet),
		paginator:          pager,
		inspector:          NewInspector(styleSet, keys),
		mainState:          SurfaceLoading,
		contextState:       SurfaceLoading,
		uiReady:            false,
		sortByView:         map[ViewID]SortConfig{},
		sortColumn:         -1,
		sortDesc:           true,
	}

	model.applyTheme(cfg.Theme, false)

	width, height := initialSize()
	model.resize(width, height)
	model.refreshView()

	return model
}

func (m Model) Init() tea.Cmd {
	windowSizeCmd := func() tea.Msg {
		width, height := initialSize()
		return tea.WindowSizeMsg{Width: width, Height: height}
	}
	return tea.Batch(windowSizeCmd, pulseTick(), initialLoadTick())
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
	case pulseMsg:
		m.pulseOn = !m.pulseOn
		cmds := []tea.Cmd{pulseTick()}
		if m.autoRefresh && !m.refreshPending && time.Since(m.lastRefresh) >= m.refreshEvery {
			m.startMockRefresh(false)
			cmds = append(cmds, mockRefreshTick(m.refreshCount%5 == 4))
		}
		return m, tea.Batch(cmds...)
	case uiLoadDoneMsg:
		m.uiReady = true
		m.refreshView()
		m.syncSurfaceStates()
		return m, nil
	case mockRefreshDoneMsg:
		m.refreshPending = false
		m.lastRefresh = time.Now()
		if msg.failed {
			m.mainState = SurfaceStale
			m.contextState = SurfaceStale
			return m, m.pushToast(ToastWarn, "Refresh failed; showing cached data (ctrl+r retry)")
		}
		m.mainState = SurfaceIdle
		m.contextState = SurfaceIdle
		m.refreshView()
		m.syncSurfaceStates()
		return m, nil
	case toastClearMsg:
		if m.toast.Active && m.toast.ID == msg.id {
			m.toast = ToastState{}
		}
		return m, nil
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

	return m.updateMain(msg)
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

	innerWidth := max(m.layout.MainWidth-2, 1)
	innerHeight := max(m.layout.MainHeight-2, 1)
	m.table.SetWidth(innerWidth)
	m.table.SetHeight(m.layout.PrimaryTableHeight)

	m.paginator.PerPage = m.layout.PrimaryTableHeight

	contextBodyHeight := max(m.layout.ContextHeight-4, 1)
	contextWidth := max(innerWidth-2, 1)
	if m.layout.ContextHeight == 0 {
		contextBodyHeight = 0
	}
	m.contextViewport.Width = contextWidth
	m.contextViewport.Height = contextBodyHeight

	m.sidebar.SetSize(max(m.layout.SidebarWidth-2, 1), max(m.layout.SidebarHeight-10, 1))
	mainBodyHeight := innerHeight
	if !m.contextCollapsed && m.layout.ContextHeight > 0 {
		mainBodyHeight -= m.layout.ContextHeight
	}
	if mainBodyHeight < 1 {
		mainBodyHeight = 1
	}
	m.mainPanel.Width = innerWidth
	m.mainPanel.Height = mainBodyHeight
	if (m.contextCollapsed || m.layout.ContextHeight == 0) && m.focus == FocusContext {
		m.focus = FocusMain
	}

	m.inspector.Resize(width, height)
	m.resizePalette()
	m.updateMainPanel()
}

func (m *Model) refreshView() {
	columns, rows, rowIDs := BuildRowsForView(m.view, &m.store, m.styles, max(m.layout.MainWidth-2, 1))
	cfg, ok := m.sortByView[m.view]
	if !ok || cfg.Column < 0 || cfg.Column >= len(columns) {
		cfg = defaultSortConfig(columns)
		m.sortByView[m.view] = cfg
	}
	m.sortColumn = cfg.Column
	m.sortDesc = cfg.Desc
	rows, rowIDs = SortRowsForView(m.view, &m.store, columns, rows, rowIDs, cfg.Column, cfg.Desc)
	m.columns = columns
	m.baseRows = rows
	m.baseRowIDs = rowIDs
	m.applyFilter()
	m.syncSidebarSelection()
}

func (m *Model) applyFilter() {
	rows, rowIDs := filterRows(m.baseRows, m.baseRowIDs, m.searchQuery)
	m.filteredRows = rows
	m.filteredRowIDs = rowIDs
	m.syncSurfaceStates()
	m.applyTableRows()
}

func (m *Model) applyTableRows() {
	m.table.SetRows(nil)
	m.table.SetColumns(columnsWithSortIndicators(m.columns, m.sortColumn, m.sortDesc))
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
	styled := components.StyleRows(truncated, m.columns, cursor, m.styles)
	m.table.SetRows(styled)
	m.updatePaginator()
	m.updateContext()
	m.updateMainPanel()
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
	if selectedID != m.contextSelectedID {
		m.contextSelectedID = selectedID
		m.contextOffsets = map[ContextTab]int{}
	}
	content := BuildContextTabContent(m.view, &m.store, selectedID, m.contextTab)
	content = utils.FilterLines(content, m.contextQuery)
	if m.contextViewport.Width > 0 {
		content = utils.WrapText(content, m.contextViewport.Width)
	}
	m.contextViewport.SetContent(content)
	offset := m.contextOffsets[m.contextTab]
	m.contextViewport.SetYOffset(offset)
	m.syncSurfaceStates()
	m.updateMainPanel()
}

func (m *Model) updateMainPanel() {
	if m.mainPanel.Width == 0 || m.mainPanel.Height == 0 {
		return
	}
	content := buildMainContent(*m)
	content = sanitizeRenderable(content)
	content = truncateLines(content, max(m.mainPanel.Width, 1))
	content = clampSection(content, max(m.mainPanel.Width, 1), max(m.mainPanel.Height, 1))
	m.mainPanel.SetContent(content)
	m.mainPanel.GotoTop()
}

func (m *Model) applyTheme(themeKey string, persist bool) {
	registry := styles.ThemeRegistry()
	key := strings.ToLower(strings.TrimSpace(themeKey))
	if key == "" {
		key = "taskforge"
	}
	selected, ok := registry[key]
	if !ok {
		selected = styles.DefaultTheme()
		key = "taskforge"
	}

	m.theme = selected
	m.themeName = key
	m.styles = styles.NewStyles(selected)
	m.table.SetStyles(components.TableStyles(m.styles))
	m.inspector.ApplyStyles(m.styles)
	m.palette = buildPalette(m.theme, m.paletteRecent)
	m.sidebar = buildSidebar(m.theme, m.view)
	m.resizePalette()

	if persist {
		m.config.Theme = key
		_ = config.Save(m.configPath, m.config)
	}

	m.refreshView()
}

func (m *Model) resizePalette() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	width := max(m.width-2, 1)
	height := max(m.height-5, 1)
	m.palette.SetSize(width, height)
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
	if key.Matches(msg, m.keys.Retry) && m.canRetry() {
		m.startMockRefresh(true)
		clearToastCmd := m.pushToast(ToastInfo, "Retrying refresh...")
		return m, tea.Batch(mockRefreshTick(false), clearToastCmd)
	}
	if key.Matches(msg, m.keys.Palette) {
		m.palette = buildPalette(m.theme, m.paletteRecent)
		m.resizePalette()
		m.palette.ResetFilter()
		m.ensurePaletteSelection(true)
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
		if m.contextCollapsed && m.focus == FocusContext {
			m.focus = FocusMain
		}
		m.resize(m.width, m.height)
		m.updateContext()
		return m, nil
	}
	if key.Matches(msg, m.keys.NextScreen) {
		m.focusNext()
		return m, nil
	}
	if key.Matches(msg, m.keys.PrevScreen) {
		m.focusPrev()
		return m, nil
	}
	if msg.String() == "left" {
		if m.focus == FocusContext {
			m.focus = FocusMain
		} else {
			m.focus = FocusSidebar
		}
		return m, nil
	}
	if msg.String() == "right" {
		if m.focus == FocusSidebar {
			m.focus = FocusMain
			return m, nil
		}
		if m.focus == FocusMain && !m.contextCollapsed && m.layout.ContextHeight > 0 {
			m.focus = FocusContext
			return m, nil
		}
		return m, nil
	}
	if m.focus == FocusMain && m.view == ViewRuns && key.Matches(msg, m.keys.Enter) {
		m.openInspector()
		return m, nil
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.RunWorkflow) {
		m.queueRunForSelectedWorkflow()
		return m, m.pushToast(ToastSuccess, "Workflow run queued")
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.ToggleActive) {
		m.toggleWorkflowActive()
		return m, m.pushToast(ToastInfo, "Workflow status updated")
	}
	if m.view == ViewTokens && key.Matches(msg, m.keys.RevokeToken) {
		m.toggleTokenRevoked()
		return m, m.pushToast(ToastInfo, "Token status updated")
	}

	if m.focus == FocusSidebar {
		return m.updateSidebar(msg)
	}
	if m.focus == FocusContext {
		return m.updateContextPane(msg)
	}
	return m.updateMain(msg)
}

func (m *Model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.SortColumn) {
			m.cycleSortColumn()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.SortDirection) {
			m.toggleSortDirection()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.JumpTop) {
			m.table.SetCursor(0)
			m.applyTableRows()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.JumpBottom) {
			if n := len(m.filteredRows); n > 0 {
				m.table.SetCursor(n - 1)
				m.applyTableRows()
			}
			return m, nil
		}
		if m.handleMainScroll(keyMsg) {
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
	case "up", "k", "shift+up", "shift+k":
		m.contextViewport.LineUp(1)
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	case "down", "j", "shift+down", "shift+j":
		m.contextViewport.LineDown(1)
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	case "pgup", "pageup", "shift+pgup", "shift+pageup", "ctrl+u":
		m.contextViewport.LineUp(m.contextViewport.Height)
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	case "pgdown", "pagedown", "pgdn", "shift+pgdown", "shift+pagedown", "shift+pgdn", "ctrl+d":
		m.contextViewport.LineDown(m.contextViewport.Height)
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	case "home":
		m.contextViewport.GotoTop()
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	case "end":
		m.contextViewport.GotoBottom()
		m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
		m.updateMainPanel()
		return true
	}
	return false
}

func (m *Model) updateContextPane(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "]", "l":
			m.nextContextTab()
			return m, nil
		case "[", "h":
			m.prevContextTab()
			return m, nil
		case "1":
			m.setContextTab(ContextTabOverview)
			return m, nil
		case "2":
			m.setContextTab(ContextTabJSON)
			return m, nil
		case "3":
			m.setContextTab(ContextTabSteps)
			return m, nil
		case "4":
			m.setContextTab(ContextTabLogs)
			return m, nil
		}
		if m.handleContextScroll(keyMsg) {
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) nextContextTab() {
	m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
	m.contextTab = (m.contextTab + 1) % 4
	m.updateContext()
}

func (m *Model) prevContextTab() {
	m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
	m.contextTab = (m.contextTab + 3) % 4
	m.updateContext()
}

func (m *Model) setContextTab(tab ContextTab) {
	m.contextOffsets[m.contextTab] = m.contextViewport.YOffset
	m.contextTab = tab
	m.updateContext()
}

func (m *Model) handleMainScroll(msg tea.KeyMsg) bool {
	if m.mainPanel.Height <= 0 {
		return false
	}
	switch msg.String() {
	case "alt+up", "alt+k":
		m.mainPanel.LineUp(1)
		return true
	case "alt+down", "alt+j":
		m.mainPanel.LineDown(1)
		return true
	case "pgup", "pageup":
		m.mainPanel.LineUp(m.mainPanel.Height)
		return true
	case "pgdown", "pagedown", "pgdn":
		m.mainPanel.LineDown(m.mainPanel.Height)
		return true
	case "home":
		m.mainPanel.GotoTop()
		return true
	case "end":
		m.mainPanel.GotoBottom()
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
			if ok && !item.Section {
				m.runPaletteAction(item.Action)
				m.showPalette = false
			}
			return m, nil
		}
		if shouldAutoStartPaletteFilter(keyMsg, m.palette.SettingFilter()) {
			startMsg := tea.KeyMsg(tea.Key{Type: tea.KeyRunes, Runes: []rune{'/'}})
			var startCmd tea.Cmd
			m.palette, startCmd = m.palette.Update(startMsg)
			var typeCmd tea.Cmd
			m.palette, typeCmd = m.palette.Update(keyMsg)
			return m, tea.Batch(startCmd, typeCmd)
		}
	}
	var cmd tea.Cmd
	m.palette, cmd = m.palette.Update(msg)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			m.ensurePaletteSelection(false)
		case "down", "j", "pgup", "pageup", "pgdown", "pagedown", "pgdn", "home", "end":
			m.ensurePaletteSelection(true)
		default:
			if m.palette.SettingFilter() || m.palette.IsFiltered() {
				m.ensurePaletteSelection(true)
			}
		}
	}
	return m, cmd
}

func (m *Model) ensurePaletteSelection(forward bool) {
	items := m.palette.VisibleItems()
	if len(items) == 0 {
		return
	}
	for attempts := 0; attempts < len(items); attempts++ {
		item, ok := m.palette.SelectedItem().(paletteItem)
		if !ok || !item.Section {
			return
		}
		if forward {
			m.palette.CursorDown()
		} else {
			m.palette.CursorUp()
		}
	}
}

func shouldAutoStartPaletteFilter(msg tea.KeyMsg, alreadyFiltering bool) bool {
	if alreadyFiltering {
		return false
	}
	if msg.Type != tea.KeyRunes || len(msg.Runes) == 0 || msg.Alt {
		return false
	}
	for _, r := range msg.Runes {
		if r >= 32 && r != '/' {
			return true
		}
	}
	return false
}

func (m *Model) runPaletteAction(action paletteAction) {
	m.rememberPaletteAction(action)
	switch action.Kind {
	case paletteNoop:
		return
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
	case paletteSetTheme:
		m.applyTheme(string(action.View), true)
	}
}

func pulseTick() tea.Cmd {
	return tea.Tick(650*time.Millisecond, func(time.Time) tea.Msg {
		return pulseMsg{}
	})
}

func initialLoadTick() tea.Cmd {
	return tea.Tick(220*time.Millisecond, func(time.Time) tea.Msg {
		return uiLoadDoneMsg{}
	})
}

func mockRefreshTick(failed bool) tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg {
		return mockRefreshDoneMsg{failed: failed}
	})
}

func (m *Model) startMockRefresh(manual bool) {
	m.refreshPending = true
	m.refreshCount++
	m.mainState = SurfaceRefreshing
	m.contextState = SurfaceRefreshing
	if manual {
		m.toast = ToastState{}
	}
}

func (m *Model) syncSurfaceStates() {
	if !m.uiReady {
		return
	}
	if m.mainState != SurfaceRefreshing && m.mainState != SurfaceError && m.mainState != SurfaceStale {
		if len(m.filteredRows) == 0 {
			m.mainState = SurfaceEmpty
		} else {
			m.mainState = SurfaceSuccess
		}
	}
	if m.contextState != SurfaceRefreshing && m.contextState != SurfaceError && m.contextState != SurfaceStale {
		content := strings.TrimSpace(m.contextViewport.View())
		if content == "" {
			m.contextState = SurfaceEmpty
		} else {
			m.contextState = SurfaceSuccess
		}
	}
}

func (m *Model) pushToast(level ToastLevel, message string) tea.Cmd {
	m.toast.ID++
	m.toast.Active = true
	m.toast.Level = level
	m.toast.Message = message
	id := m.toast.ID
	return tea.Tick(2400*time.Millisecond, func(time.Time) tea.Msg {
		return toastClearMsg{id: id}
	})
}

func (m *Model) canRetry() bool {
	return m.mainState == SurfaceError || m.mainState == SurfaceStale || m.contextState == SurfaceError || m.contextState == SurfaceStale
}

func (m *Model) rememberPaletteAction(action paletteAction) {
	if action.Kind == paletteNoop {
		return
	}
	maxRecent := 5
	recent := []paletteAction{action}
	for _, existing := range m.paletteRecent {
		if paletteActionKey(existing) == paletteActionKey(action) {
			continue
		}
		recent = append(recent, existing)
		if len(recent) >= maxRecent {
			break
		}
	}
	m.paletteRecent = recent
}

func paletteActionKey(action paletteAction) string {
	return fmt.Sprintf("%d:%s", action.Kind, action.View)
}

func paletteItemFromAction(action paletteAction) paletteItem {
	switch action.Kind {
	case paletteGoToView:
		switch action.View {
		case ViewDashboard:
			return paletteItem{Label: "Go: Dashboard", Detail: "Recent", Action: action}
		case ViewWorkflows:
			return paletteItem{Label: "Go: Workflows", Detail: "Recent", Action: action}
		case ViewRuns:
			return paletteItem{Label: "Go: Runs", Detail: "Recent", Action: action}
		case ViewTriggers:
			return paletteItem{Label: "Go: Triggers", Detail: "Recent", Action: action}
		case ViewEvents:
			return paletteItem{Label: "Go: Events", Detail: "Recent", Action: action}
		case ViewSecrets:
			return paletteItem{Label: "Go: Secrets", Detail: "Recent", Action: action}
		case ViewTokens:
			return paletteItem{Label: "Go: API Tokens", Detail: "Recent", Action: action}
		}
	case paletteRunWorkflow:
		return paletteItem{Label: "Action: Run selected workflow", Detail: "Recent", Action: action}
	case paletteClearFilters:
		return paletteItem{Label: "Action: Clear filters", Detail: "Recent", Action: action}
	case paletteToggleRefresh:
		return paletteItem{Label: "Toggle: Auto refresh", Detail: "Recent", Action: action}
	case paletteSetTheme:
		switch action.View {
		case ViewID("catppuccin"):
			return paletteItem{Label: "Theme: Catppuccin", Detail: "Recent", Action: action}
		case ViewID("tokyo-night"):
			return paletteItem{Label: "Theme: Tokyo Night", Detail: "Recent", Action: action}
		case ViewID("fallout"):
			return paletteItem{Label: "Theme: Fallout (CRT)", Detail: "Recent", Action: action}
		case ViewID("retro-amber"):
			return paletteItem{Label: "Theme: Retro Amber", Detail: "Recent", Action: action}
		}
	}
	return paletteItem{Label: "Action", Detail: "Recent", Action: action}
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

func (m *Model) focusNext() {
	if m.focus == FocusSidebar {
		m.focus = FocusMain
		return
	}
	if m.focus == FocusMain {
		if !m.contextCollapsed && m.layout.ContextHeight > 0 {
			m.focus = FocusContext
		} else {
			m.focus = FocusSidebar
		}
		return
	}
	m.focus = FocusSidebar
}

func (m *Model) focusPrev() {
	if m.focus == FocusSidebar {
		if !m.contextCollapsed && m.layout.ContextHeight > 0 {
			m.focus = FocusContext
		} else {
			m.focus = FocusMain
		}
		return
	}
	if m.focus == FocusMain {
		m.focus = FocusSidebar
		return
	}
	m.focus = FocusMain
}

func (m *Model) updateSidebar(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "right", "enter":
			m.focus = FocusMain
			return m, nil
		}
	}
	prev := m.sidebar.Index()
	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	if m.sidebar.Index() != prev {
		m.syncViewFromSidebar()
	}
	return m, cmd
}

func (m *Model) syncViewFromSidebar() {
	item, ok := m.sidebar.SelectedItem().(navItem)
	if !ok {
		return
	}
	if item.ID == m.view {
		return
	}
	m.view = item.ID
	m.refreshView()
}

func (m *Model) syncSidebarSelection() {
	items := m.sidebar.Items()
	for i, item := range items {
		nav, ok := item.(navItem)
		if !ok {
			continue
		}
		if nav.ID == m.view {
			m.sidebar.Select(i)
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

func buildPalette(theme styles.Theme, recentActions []paletteAction) list.Model {
	items := []list.Item{
		paletteItem{Label: ":: Navigation", Detail: "", Section: true, Action: paletteAction{Kind: paletteNoop}},
		paletteItem{Label: "Go: Dashboard", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewDashboard}},
		paletteItem{Label: "Go: Workflows", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewWorkflows}},
		paletteItem{Label: "Go: Runs", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewRuns}},
		paletteItem{Label: "Go: Triggers", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewTriggers}},
		paletteItem{Label: "Go: Events", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewEvents}},
		paletteItem{Label: "Go: Secrets", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewSecrets}},
		paletteItem{Label: "Go: API Tokens", Detail: "Navigation", Action: paletteAction{Kind: paletteGoToView, View: ViewTokens}},
		paletteItem{Label: ":: Actions", Detail: "", Section: true, Action: paletteAction{Kind: paletteNoop}},
		paletteItem{Label: "Action: Run selected workflow", Detail: "Workflow", Action: paletteAction{Kind: paletteRunWorkflow}},
		paletteItem{Label: "Action: Clear filters", Detail: "Table", Action: paletteAction{Kind: paletteClearFilters}},
		paletteItem{Label: "Toggle: Auto refresh", Detail: "System", Action: paletteAction{Kind: paletteToggleRefresh}},
		paletteItem{Label: ":: Themes", Detail: "", Section: true, Action: paletteAction{Kind: paletteNoop}},
		paletteItem{Label: "Theme: Catppuccin", Detail: "Theme", Action: paletteAction{Kind: paletteSetTheme, View: ViewID("catppuccin")}},
		paletteItem{Label: "Theme: Tokyo Night", Detail: "Theme", Action: paletteAction{Kind: paletteSetTheme, View: ViewID("tokyo-night")}},
		paletteItem{Label: "Theme: Fallout (CRT)", Detail: "Theme", Action: paletteAction{Kind: paletteSetTheme, View: ViewID("fallout")}},
		paletteItem{Label: "Theme: Retro Amber", Detail: "Theme", Action: paletteAction{Kind: paletteSetTheme, View: ViewID("retro-amber")}},
	}
	if len(recentActions) > 0 {
		recent := make([]list.Item, 0, len(recentActions)+1)
		recent = append(recent, paletteItem{Label: ":: Recent", Detail: "", Section: true, Action: paletteAction{Kind: paletteNoop}})
		for _, action := range recentActions {
			recent = append(recent, paletteItemFromAction(action))
		}
		items = append(recent, items...)
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetHeight(1)
	selectedBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.Accent)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Background(theme.SurfaceAlt).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.Accent)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(theme.Text).
		Background(theme.SurfaceAlt).
		Inherit(selectedBorder)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(theme.Text)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(theme.Muted)
	model := list.New(items, delegate, 64, 16)
	listStyles := list.DefaultStyles()
	listStyles.Title = lipgloss.NewStyle().Foreground(theme.Text).Bold(true)
	listStyles.FilterPrompt = lipgloss.NewStyle().Foreground(theme.Accent)
	listStyles.FilterCursor = lipgloss.NewStyle().Foreground(theme.Accent)
	listStyles.PaginationStyle = lipgloss.NewStyle().Foreground(theme.Muted)
	listStyles.HelpStyle = lipgloss.NewStyle().Foreground(theme.Muted)
	listStyles.NoItems = lipgloss.NewStyle().Foreground(theme.Muted)
	model.Styles = listStyles
	model.SetFilteringEnabled(true)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)
	model.SetShowTitle(false)
	model.DisableQuitKeybindings()
	return model
}

func buildSidebar(theme styles.Theme, selected ViewID) list.Model {
	items := []list.Item{
		navItem{ID: ViewDashboard, Label: "Dashboard"},
		navItem{ID: ViewWorkflows, Label: "Workflows"},
		navItem{ID: ViewRuns, Label: "Runs"},
		navItem{ID: ViewTriggers, Label: "Triggers"},
		navItem{ID: ViewEvents, Label: "Events"},
		navItem{ID: ViewSecrets, Label: "Secrets"},
		navItem{ID: ViewTokens, Label: "API Tokens"},
	}
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(theme.Background).Background(theme.Accent).Bold(true)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(theme.Muted)
	model := list.New(items, delegate, 20, 10)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)
	model.SetFilteringEnabled(false)
	model.SetShowTitle(false)
	model.SetShowPagination(false)
	model.DisableQuitKeybindings()

	for i, item := range items {
		nav, ok := item.(navItem)
		if ok && nav.ID == selected {
			model.Select(i)
			break
		}
	}
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
	if sortColumn < 0 || sortColumn >= len(decorated) {
		return decorated
	}
	arrow := " ▲"
	if desc {
		arrow = " ▼"
	}
	decorated[sortColumn].Title = strings.TrimSpace(decorated[sortColumn].Title) + arrow
	return decorated
}

func (m *Model) cycleSortColumn() {
	if len(m.columns) == 0 {
		return
	}
	if m.sortColumn < 0 {
		m.sortColumn = 0
	} else {
		m.sortColumn = (m.sortColumn + 1) % len(m.columns)
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
