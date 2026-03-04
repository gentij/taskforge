package app

import (
	"encoding/json"
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
	"github.com/charmbracelet/x/ansi"
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
	paletteRenameWorkflow
	paletteCreateTrigger
	paletteRenameTrigger
	paletteToggleTrigger
	paletteDeleteWorkflow
	paletteDeleteTrigger
	paletteCreateSecret
	paletteUpdateSecret
	paletteDeleteSecret
	paletteSetStatusScope
	paletteShowCLIHandoff
	paletteClearRecent
	paletteSetTheme
	paletteSetNetworkProfile
)

type statusScope int

const (
	statusScopeAll statusScope = iota
	statusScopeActive
	statusScopeInactive
)

type pulseMsg struct{}

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

type NetworkProfile int

const (
	NetworkFast NetworkProfile = iota
	NetworkNormal
	NetworkSlow
	NetworkFlaky
)

type paletteAction struct {
	Kind    paletteActionType
	View    ViewID
	Profile NetworkProfile
	Value   string
}

type ContextTab int

const (
	ContextTabOverview ContextTab = iota
	ContextTabJSON
	ContextTabSteps
	ContextTabLogs
)

type actionModalMode int

const (
	actionModalNone actionModalMode = iota
	actionModalRenameWorkflow
	actionModalRenameTrigger
	actionModalCreateTrigger
	actionModalUpdateTrigger
	actionModalCreateSecret
	actionModalUpdateSecret
	actionModalConfirmDelete
	actionModalCLIHandoff
)

type actionModalState struct {
	Active         bool
	Mode           actionModalMode
	Title          string
	Description    string
	Validation     string
	ShowValidation bool
	Primary        textinput.Model
	Secondary      textinput.Model
	Tertiary       textinput.Model
	Confirm        textinput.Model
	Focus          int
	WorkflowID     string
	TriggerID      string
	SecretID       string
	DeleteKind     string
	TriggerType    string
	TriggerActive  bool
	ConfirmPhrase  string
	CLICommand     string
}

type paletteItem struct {
	Label          string
	Detail         string
	Action         paletteAction
	Section        bool
	Enabled        bool
	DisabledReason string
	Keywords       []string
}

type paletteBuildState struct {
	View         ViewID
	HasSelection bool
	HasFilter    bool
	HasScope     bool
	Scope        statusScope
	AutoRefresh  bool
	Profile      NetworkProfile
	HasRecent    bool
}

type SortConfig struct {
	Column int
	Desc   bool
}

func (p paletteItem) FilterValue() string {
	if p.Section {
		return ""
	}
	parts := []string{p.Label, p.Detail, strings.Join(p.Keywords, " ")}
	if !p.Enabled && p.DisabledReason != "" {
		parts = append(parts, p.DisabledReason)
	}
	return strings.Join(parts, " ")
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
	action   actionModalState

	autoRefresh     bool
	refreshEvery    time.Duration
	lastRefresh     time.Time
	pulseOn         bool
	refreshPending  bool
	refreshCount    int
	mutationPending bool
	networkProfile  NetworkProfile
	workspaceName   string
	workerCount     int
	apiStatus       string
	paletteRecent   []paletteAction

	mainState    SurfaceState
	contextState SurfaceState
	toast        ToastState
	uiReady      bool

	sortByView map[ViewID]SortConfig
	sortColumn int
	sortDesc   bool

	statusScopeByView map[ViewID]statusScope

	paginator paginator.Model
}

func NewModel(client *api.Client, serverURL string, tokenSet bool, cfg config.Config, configPath string) Model {
	now := time.Now()
	store := data.Store{}
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
	palette := buildPalette(defaultTheme, nil, paletteBuildState{View: ViewDashboard, Profile: NetworkNormal})
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
		action:             actionModalState{},
		refreshEvery:       2 * time.Second,
		lastRefresh:        now,
		pulseOn:            false,
		refreshPending:     false,
		refreshCount:       0,
		mutationPending:    false,
		networkProfile:     NetworkNormal,
		workspaceName:      "personal",
		workerCount:        0,
		apiStatus:          apiStatus(tokenSet),
		paginator:          pager,
		inspector:          NewInspector(styleSet, keys),
		mainState:          SurfaceLoading,
		contextState:       SurfaceLoading,
		uiReady:            false,
		sortByView:         map[ViewID]SortConfig{},
		sortColumn:         -1,
		sortDesc:           true,
		statusScopeByView:  map[ViewID]statusScope{},
	}
	model.setNetworkProfile(NetworkNormal)

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
	return tea.Batch(windowSizeCmd, pulseTick(), fetchSnapshotCmd(m.client, m.profileDelay(), m.profileShouldFail(false)))
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
			cmds = append(cmds, fetchSnapshotCmd(m.client, m.profileDelay(), m.profileShouldFail(false)))
		}
		return m, tea.Batch(cmds...)
	case snapshotLoadedMsg:
		m.refreshPending = false
		m.lastRefresh = time.Now()
		if msg.err != nil {
			m.apiStatus = "OFFLINE"
			if !m.uiReady {
				m.uiReady = true
				m.mainState = SurfaceError
				m.contextState = SurfaceError
				m.refreshView()
			}
			if m.mainState != SurfaceError {
				m.mainState = SurfaceStale
			}
			if m.contextState != SurfaceError {
				m.contextState = SurfaceStale
			}
			return m, m.pushToast(ToastWarn, "Failed to sync data (ctrl+r retry)")
		}
		m.store = msg.store
		m.apiStatus = msg.apiStatus
		m.uiReady = true
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
	case mutationResultMsg:
		m.mutationPending = false
		if msg.err != nil {
			return m, m.pushToast(ToastError, mutationErrorMessage(msg.err))
		}
		cmds := []tea.Cmd{}
		if strings.TrimSpace(msg.successMessage) != "" {
			cmds = append(cmds, m.pushToast(ToastSuccess, msg.successMessage))
		}
		if msg.refresh {
			m.startMockRefresh(false)
			cmds = append(cmds, fetchSnapshotCmd(m.client, m.profileDelay(), m.profileShouldFail(false)))
		}
		return m, tea.Batch(cmds...)
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
	selectedID := m.selectedRowID()
	cursor := m.table.Cursor()
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
	m.applyFilterWithSelection(selectedID, cursor)
	m.syncSidebarSelection()
}

func (m *Model) applyFilter() {
	selectedID := m.selectedRowID()
	cursor := m.table.Cursor()
	m.applyFilterWithSelection(selectedID, cursor)
}

func (m *Model) applyFilterWithSelection(selectedID string, cursor int) {
	rows, rowIDs := m.scopeRowsForCurrentView(m.baseRows, m.baseRowIDs)
	rows, rowIDs = filterRows(rows, rowIDs, m.searchQuery)
	m.filteredRows = rows
	m.filteredRowIDs = rowIDs
	m.restoreSelection(selectedID, cursor)
	m.syncSurfaceStates()
	m.applyTableRows()
}

func (m *Model) scopeRowsForCurrentView(rows []table.Row, rowIDs []string) ([]table.Row, []string) {
	scope := m.currentStatusScope(m.view)
	if scope == statusScopeAll || !supportsStatusScope(m.view) {
		return rows, rowIDs
	}
	filteredRows := make([]table.Row, 0, len(rows))
	filteredIDs := make([]string, 0, len(rowIDs))
	for i, id := range rowIDs {
		if i >= len(rows) {
			break
		}
		active, ok := m.isRowActiveForScope(m.view, id)
		if !ok {
			continue
		}
		if (scope == statusScopeActive && active) || (scope == statusScopeInactive && !active) {
			filteredRows = append(filteredRows, rows[i])
			filteredIDs = append(filteredIDs, id)
		}
	}
	return filteredRows, filteredIDs
}

func (m *Model) isRowActiveForScope(view ViewID, rowID string) (bool, bool) {
	switch view {
	case ViewWorkflows:
		wf, ok := workflowByID(&m.store, rowID)
		if !ok {
			return false, false
		}
		return wf.Active, true
	case ViewTriggers:
		trg, ok := triggerByID(&m.store, rowID)
		if !ok {
			return false, false
		}
		return trg.Active, true
	default:
		return false, false
	}
}

func (m *Model) restoreSelection(selectedID string, cursor int) {
	if len(m.filteredRowIDs) == 0 {
		m.table.SetCursor(0)
		return
	}
	if selectedID != "" {
		for i, id := range m.filteredRowIDs {
			if id == selectedID {
				m.table.SetCursor(i)
				return
			}
		}
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(m.filteredRowIDs) {
		cursor = len(m.filteredRowIDs) - 1
	}
	m.table.SetCursor(cursor)
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
	m.palette = buildPalette(m.theme, m.paletteRecent, m.paletteState())
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
	if m.action.Active {
		return m.updateActionModal(msg)
	}
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
		return m, tea.Batch(fetchSnapshotCmd(m.client, m.profileDelay(), m.profileShouldFail(true)), clearToastCmd)
	}
	if key.Matches(msg, m.keys.Palette) {
		m.palette = buildPalette(m.theme, m.paletteRecent, m.paletteState())
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
		return m, m.queueRunForSelectedWorkflowCmd()
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.ToggleActive) {
		return m, m.toggleWorkflowActiveCmd()
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.Rename) {
		return m, m.openRenameWorkflowModalCmd()
	}
	if (m.view == ViewWorkflows || m.view == ViewTriggers) && key.Matches(msg, m.keys.CreateTrigger) {
		return m, m.openCreateTriggerModalCmd()
	}
	if m.view == ViewSecrets && key.Matches(msg, m.keys.CreateTrigger) {
		return m, m.openCreateSecretModalCmd()
	}
	if m.view == ViewTriggers && key.Matches(msg, m.keys.Rename) {
		return m, m.openRenameTriggerModalCmd()
	}
	if m.view == ViewSecrets && key.Matches(msg, m.keys.Rename) {
		return m, m.openUpdateSecretModalCmd()
	}
	if m.view == ViewTriggers && key.Matches(msg, m.keys.ToggleActive) {
		return m, m.toggleTriggerActiveCmd()
	}
	if m.view == ViewWorkflows && key.Matches(msg, m.keys.RevokeToken) {
		return m, m.openDeleteWorkflowModalCmd()
	}
	if m.view == ViewTriggers && key.Matches(msg, m.keys.RevokeToken) {
		return m, m.openDeleteTriggerModalCmd()
	}
	if m.view == ViewSecrets && key.Matches(msg, m.keys.RevokeToken) {
		return m, m.openDeleteSecretModalCmd()
	}
	if m.view == ViewTokens && key.Matches(msg, m.keys.RevokeToken) {
		return m, m.pushToast(ToastWarn, "API tokens are not available yet")
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
		if key.Matches(keyMsg, m.keys.CycleStatus) {
			if supportsStatusScope(m.view) {
				return m, m.cycleStatusScopeForCurrentViewCmd()
			}
		}
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
				if !item.Enabled {
					msg := "Command unavailable"
					if item.DisabledReason != "" {
						msg = item.DisabledReason
					}
					return m, m.pushToast(ToastWarn, msg)
				}
				cmd := m.runPaletteAction(item.Action)
				m.showPalette = false
				return m, cmd
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

func (m *Model) runPaletteAction(action paletteAction) tea.Cmd {
	switch action.Kind {
	case paletteNoop:
		return nil
	case paletteGoToView:
		m.rememberPaletteAction(action)
		m.view = action.View
		m.refreshView()
		return nil
	case paletteToggleRefresh:
		m.rememberPaletteAction(action)
		m.autoRefresh = !m.autoRefresh
		if m.autoRefresh {
			m.lastRefresh = time.Now()
		}
		return m.pushToast(ToastInfo, "Auto refresh toggled")
	case paletteClearFilters:
		m.rememberPaletteAction(action)
		m.clearCurrentViewFilters()
		return m.pushToast(ToastInfo, "Filters cleared")
	case paletteRunWorkflow:
		m.rememberPaletteAction(action)
		return m.queueRunForSelectedWorkflowCmd()
	case paletteRenameWorkflow:
		m.rememberPaletteAction(action)
		return m.openRenameWorkflowModalCmd()
	case paletteCreateTrigger:
		m.rememberPaletteAction(action)
		return m.openCreateTriggerModalCmd()
	case paletteRenameTrigger:
		m.rememberPaletteAction(action)
		return m.openRenameTriggerModalCmd()
	case paletteToggleTrigger:
		m.rememberPaletteAction(action)
		return m.toggleTriggerActiveCmd()
	case paletteDeleteWorkflow:
		m.rememberPaletteAction(action)
		return m.openDeleteWorkflowModalCmd()
	case paletteDeleteTrigger:
		m.rememberPaletteAction(action)
		return m.openDeleteTriggerModalCmd()
	case paletteCreateSecret:
		m.rememberPaletteAction(action)
		return m.openCreateSecretModalCmd()
	case paletteUpdateSecret:
		m.rememberPaletteAction(action)
		return m.openUpdateSecretModalCmd()
	case paletteDeleteSecret:
		m.rememberPaletteAction(action)
		return m.openDeleteSecretModalCmd()
	case paletteSetStatusScope:
		m.rememberPaletteAction(action)
		if !supportsStatusScope(m.view) {
			return m.pushToast(ToastWarn, "Status scope is available in Workflows and Triggers")
		}
		next := statusScopeFromValue(action.Value)
		m.setStatusScopeForView(m.view, next)
		return m.pushToast(ToastInfo, "Status filter: "+strings.ToLower(statusScopeLabel(next)))
	case paletteShowCLIHandoff:
		m.rememberPaletteAction(action)
		return m.openCLIHandoffModalCmd(action.Value)
	case paletteClearRecent:
		m.paletteRecent = nil
		return m.pushToast(ToastInfo, "Recent commands cleared")
	case paletteSetTheme:
		m.rememberPaletteAction(action)
		m.applyTheme(string(action.View), true)
		return m.pushToast(ToastInfo, "Theme switched")
	case paletteSetNetworkProfile:
		m.rememberPaletteAction(action)
		m.setNetworkProfile(action.Profile)
		return m.pushToast(ToastInfo, "Network profile set to "+strings.ToLower(networkProfileLabel(action.Profile)))
	}
	return nil
}

func pulseTick() tea.Cmd {
	return tea.Tick(650*time.Millisecond, func(time.Time) tea.Msg {
		return pulseMsg{}
	})
}

func (m Model) profileDelay() time.Duration {
	switch m.networkProfile {
	case NetworkFast:
		return 90 * time.Millisecond
	case NetworkSlow:
		return 900 * time.Millisecond
	case NetworkFlaky:
		return 550 * time.Millisecond
	default:
		return 240 * time.Millisecond
	}
}

func (m Model) profileShouldFail(manual bool) bool {
	if m.networkProfile != NetworkFlaky {
		return false
	}
	if manual {
		return m.refreshCount%4 == 0
	}
	return m.refreshCount%3 == 0
}

func (m Model) profileRefreshEvery() time.Duration {
	switch m.networkProfile {
	case NetworkFast:
		return 1200 * time.Millisecond
	case NetworkSlow:
		return 4 * time.Second
	case NetworkFlaky:
		return 2500 * time.Millisecond
	default:
		return 2 * time.Second
	}
}

func (m *Model) setNetworkProfile(profile NetworkProfile) {
	m.networkProfile = profile
	m.refreshEvery = m.profileRefreshEvery()
	m.lastRefresh = time.Now()
}

func networkProfileLabel(profile NetworkProfile) string {
	switch profile {
	case NetworkFast:
		return "FAST"
	case NetworkSlow:
		return "SLOW"
	case NetworkFlaky:
		return "FLAKY"
	default:
		return "NORMAL"
	}
}

func (m Model) paletteState() paletteBuildState {
	scope := m.currentStatusScope(m.view)
	hasScope := supportsStatusScope(m.view) && scope != statusScopeAll
	return paletteBuildState{
		View:         m.view,
		HasSelection: m.selectedRowID() != "",
		HasFilter:    strings.TrimSpace(m.searchQuery) != "" || hasScope,
		HasScope:     supportsStatusScope(m.view),
		Scope:        scope,
		AutoRefresh:  m.autoRefresh,
		Profile:      m.networkProfile,
		HasRecent:    len(m.paletteRecent) > 0,
	}
}

func supportsStatusScope(view ViewID) bool {
	return view == ViewWorkflows || view == ViewTriggers
}

func (m Model) currentStatusScope(view ViewID) statusScope {
	if !supportsStatusScope(view) {
		return statusScopeAll
	}
	scope, ok := m.statusScopeByView[view]
	if !ok {
		return statusScopeAll
	}
	return scope
}

func statusScopeLabel(scope statusScope) string {
	switch scope {
	case statusScopeActive:
		return "ACTIVE"
	case statusScopeInactive:
		return "INACTIVE"
	default:
		return "ALL"
	}
}

func statusScopeFromValue(value string) statusScope {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "active":
		return statusScopeActive
	case "inactive":
		return statusScopeInactive
	default:
		return statusScopeAll
	}
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

func mutationErrorMessage(err error) string {
	if err == nil {
		return "Action failed"
	}
	if apiErr := api.AsAPIError(err); apiErr != nil {
		code := strings.TrimSpace(apiErr.Code)
		msg := strings.TrimSpace(apiErr.Message)
		if code != "" && msg != "" {
			return code + ": " + msg
		}
		if msg != "" {
			return msg
		}
	}
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "Action failed"
	}
	return msg
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
	return fmt.Sprintf("%d:%s:%d:%s", action.Kind, action.View, action.Profile, action.Value)
}

func paletteItemFromAction(action paletteAction, state paletteBuildState) paletteItem {
	item := paletteItem{Label: "Action", Detail: "Recent", Action: action, Enabled: true}
	switch action.Kind {
	case paletteGoToView:
		switch action.View {
		case ViewDashboard:
			item.Label = "Go: Dashboard"
		case ViewWorkflows:
			item.Label = "Go: Workflows"
		case ViewRuns:
			item.Label = "Go: Runs"
		case ViewTriggers:
			item.Label = "Go: Triggers"
		case ViewEvents:
			item.Label = "Go: Events"
		case ViewSecrets:
			item.Label = "Go: Secrets"
		case ViewTokens:
			item.Label = "Go: API Tokens"
		}
	case paletteRunWorkflow:
		item.Label = "Action: Run selected workflow"
		if !(state.View == ViewWorkflows && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select workflow row"
			item.DisabledReason = "Select a workflow row in Workflows first"
		}
	case paletteRenameWorkflow:
		item.Label = "Action: Rename selected workflow"
		if !(state.View == ViewWorkflows && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select workflow row"
			item.DisabledReason = "Select a workflow row in Workflows first"
		}
	case paletteCreateTrigger:
		item.Label = "Action: Create trigger"
		if !((state.View == ViewWorkflows || state.View == ViewTriggers) && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select workflow or trigger row"
			item.DisabledReason = "Select a row in Workflows or Triggers first"
		}
	case paletteRenameTrigger:
		item.Label = "Action: Update selected trigger"
		if !(state.View == ViewTriggers && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select trigger row"
			item.DisabledReason = "Select a trigger row in Triggers first"
		}
	case paletteToggleTrigger:
		item.Label = "Action: Toggle selected trigger"
		if !(state.View == ViewTriggers && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select trigger row"
			item.DisabledReason = "Select a trigger row in Triggers first"
		}
	case paletteDeleteWorkflow:
		item.Label = "Action: Archive selected workflow"
		if !(state.View == ViewWorkflows && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select workflow row"
			item.DisabledReason = "Select a workflow row in Workflows first"
		}
	case paletteDeleteTrigger:
		item.Label = "Action: Archive selected trigger"
		if !(state.View == ViewTriggers && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select trigger row"
			item.DisabledReason = "Select a trigger row in Triggers first"
		}
	case paletteCreateSecret:
		item.Label = "Action: Create secret"
		if state.View != ViewSecrets {
			item.Enabled = false
			item.Detail = "Unavailable in this view"
			item.DisabledReason = "Open Secrets first"
		}
	case paletteUpdateSecret:
		item.Label = "Action: Update selected secret"
		if !(state.View == ViewSecrets && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select secret row"
			item.DisabledReason = "Select a secret row in Secrets first"
		}
	case paletteDeleteSecret:
		item.Label = "Action: Delete selected secret"
		if !(state.View == ViewSecrets && state.HasSelection) {
			item.Enabled = false
			item.Detail = "Unavailable: select secret row"
			item.DisabledReason = "Select a secret row in Secrets first"
		}
	case paletteSetStatusScope:
		if !state.HasScope {
			item.Enabled = false
			item.Detail = "Unavailable in this view"
			item.DisabledReason = "Status scope is available in Workflows and Triggers"
		}
		scope := statusScopeFromValue(action.Value)
		suffix := ""
		if state.Scope == scope {
			suffix = " (active)"
		}
		switch scope {
		case statusScopeActive:
			item.Label = "Filter: Active only" + suffix
		default:
			if scope == statusScopeInactive {
				item.Label = "Filter: Inactive only" + suffix
			} else {
				item.Label = "Filter: Show all" + suffix
			}
		}
	case paletteShowCLIHandoff:
		switch action.Value {
		case "workflow-create":
			item.Label = "CLI: Create workflow"
		case "workflow-version-create":
			item.Label = "CLI: Create workflow version"
		default:
			item.Label = "CLI: Open authoring command"
		}
	case paletteClearFilters:
		item.Label = "Action: Clear filters"
		if !state.HasFilter {
			item.Enabled = false
			item.Detail = "Unavailable: no active filters"
			item.DisabledReason = "No search filter is active"
		}
	case paletteToggleRefresh:
		item.Label = "Toggle: Auto refresh"
	case paletteSetTheme:
		item.Label = "Theme"
		switch action.View {
		case ViewID("catppuccin"):
			item.Label = "Theme: Catppuccin"
		case ViewID("tokyo-night"):
			item.Label = "Theme: Tokyo Night"
		case ViewID("fallout"):
			item.Label = "Theme: Fallout (CRT)"
		case ViewID("retro-amber"):
			item.Label = "Theme: Retro Amber"
		}
	case paletteSetNetworkProfile:
		switch action.Profile {
		case NetworkFast:
			item.Label = "Network: Fast"
		case NetworkNormal:
			item.Label = "Network: Normal"
		case NetworkSlow:
			item.Label = "Network: Slow"
		case NetworkFlaky:
			item.Label = "Network: Flaky"
		}
	case paletteClearRecent:
		item.Label = "Action: Clear recent commands"
		if !state.HasRecent {
			item.Enabled = false
			item.Detail = "Unavailable: no recent commands"
			item.DisabledReason = "Run commands from the palette first"
		}
	}
	return item
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

func (m *Model) setStatusScopeForView(view ViewID, scope statusScope) {
	if !supportsStatusScope(view) {
		return
	}
	m.statusScopeByView[view] = scope
	m.applyFilter()
}

func (m *Model) cycleStatusScopeForCurrentViewCmd() tea.Cmd {
	if !supportsStatusScope(m.view) {
		return m.pushToast(ToastWarn, "Status scope is available in Workflows and Triggers")
	}
	current := m.currentStatusScope(m.view)
	next := statusScopeAll
	switch current {
	case statusScopeAll:
		next = statusScopeActive
	case statusScopeActive:
		next = statusScopeInactive
	default:
		next = statusScopeAll
	}
	m.setStatusScopeForView(m.view, next)
	return m.pushToast(ToastInfo, "Status filter: "+strings.ToLower(statusScopeLabel(next)))
}

func (m *Model) clearCurrentViewFilters() {
	m.searchQuery = ""
	m.searchInput.SetValue("")
	if supportsStatusScope(m.view) {
		m.statusScopeByView[m.view] = statusScopeAll
	}
	m.applyFilter()
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

func (m *Model) toggleWorkflowActiveCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	selected := m.selectedRowID()
	if selected == "" {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	wf, ok := workflowByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	client := m.client
	next := !wf.Active
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.UpdateWorkflow(selected, map[string]any{"isActive": next})
		if err != nil {
			return mutationResultMsg{err: err}
		}
		message := "Workflow archived"
		if next {
			message = "Workflow restored"
		}
		return mutationResultMsg{successMessage: message, refresh: true}
	}
}

func (m *Model) queueRunForSelectedWorkflowCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewWorkflows {
		return m.pushToast(ToastWarn, "Open Workflows to run")
	}
	selected := m.selectedRowID()
	if selected == "" {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		result, err := client.RunWorkflow(selected, map[string]any{}, map[string]any{})
		if err != nil {
			return mutationResultMsg{err: err}
		}
		message := "Workflow run queued"
		if id, ok := result["workflowRunId"]; ok && strings.TrimSpace(id) != "" {
			message = "Workflow run queued: " + id
		}
		return mutationResultMsg{successMessage: message, refresh: true}
	}
}

func (m *Model) toggleTriggerActiveCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewTriggers {
		return m.pushToast(ToastWarn, "Open Triggers to toggle")
	}
	selected := m.selectedRowID()
	if selected == "" {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	trg, ok := triggerByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	client := m.client
	next := !trg.Active
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.UpdateTrigger(trg.WorkflowID, trg.ID, map[string]any{"isActive": next})
		if err != nil {
			return mutationResultMsg{err: err}
		}
		message := "Trigger archived"
		if next {
			message = "Trigger restored"
		}
		return mutationResultMsg{successMessage: message, refresh: true}
	}
}

func (m *Model) deleteWorkflowCmd(workflowID string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	workflowID = strings.TrimSpace(workflowID)
	if workflowID == "" {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.DeleteWorkflow(workflowID)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Workflow archived", refresh: true}
	}
}

func (m *Model) deleteTriggerCmd(workflowID string, triggerID string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	workflowID = strings.TrimSpace(workflowID)
	triggerID = strings.TrimSpace(triggerID)
	if workflowID == "" || triggerID == "" {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.DeleteTrigger(workflowID, triggerID)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Trigger archived", refresh: true}
	}
}

func (m *Model) createSecretCmd(name string, value string, description string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return m.pushToast(ToastWarn, "Secret name cannot be empty")
	}
	if strings.TrimSpace(value) == "" {
		return m.pushToast(ToastWarn, "Secret value cannot be empty")
	}
	payload := map[string]any{"name": name, "value": value}
	description = strings.TrimSpace(description)
	if description != "" {
		payload["description"] = description
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.CreateSecret(payload)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Secret created", refresh: true}
	}
}

func (m *Model) updateSecretCmd(secretID string, name string, value string, description string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	secretID = strings.TrimSpace(secretID)
	if secretID == "" {
		return m.pushToast(ToastWarn, "Select a secret first")
	}
	secret, ok := secretByID(&m.store, secretID)
	if !ok {
		return m.pushToast(ToastWarn, "Select a secret first")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return m.pushToast(ToastWarn, "Secret name cannot be empty")
	}
	patch := map[string]any{}
	if name != secret.Name {
		patch["name"] = name
	}
	if description != "" && description != strings.TrimSpace(secret.Description) {
		patch["description"] = description
	}
	if strings.TrimSpace(value) != "" {
		patch["value"] = value
	}
	if len(patch) == 0 {
		return m.pushToast(ToastWarn, "No changes to update")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.UpdateSecret(secretID, patch)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Secret updated", refresh: true}
	}
}

func (m *Model) deleteSecretCmd(secretID string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	secretID = strings.TrimSpace(secretID)
	if secretID == "" {
		return m.pushToast(ToastWarn, "Select a secret first")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.DeleteSecret(secretID)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Secret deleted", refresh: true}
	}
}

func (m *Model) renameWorkflowCmd(workflowID string, name string) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	name = strings.TrimSpace(name)
	if workflowID == "" {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	if name == "" {
		return m.pushToast(ToastWarn, "Workflow name cannot be empty")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.UpdateWorkflow(workflowID, map[string]any{"name": name})
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Workflow renamed", refresh: true}
	}
}

func (m *Model) updateTriggerCmd(workflowID string, triggerID string, name string, active bool, configValue map[string]any) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	name = strings.TrimSpace(name)
	if workflowID == "" || triggerID == "" {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	if name == "" {
		return m.pushToast(ToastWarn, "Trigger name cannot be empty")
	}
	trigger, ok := triggerByID(&m.store, triggerID)
	if !ok {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	if trigger.WorkflowID != workflowID {
		return m.pushToast(ToastWarn, "Trigger does not match selected workflow")
	}
	patch := map[string]any{}
	if name != strings.TrimSpace(trigger.Name) {
		patch["name"] = name
	}
	if active != trigger.Active {
		patch["isActive"] = active
	}
	if configValue == nil {
		configValue = map[string]any{}
	}
	existingConfig, err := parseJSONObject(trigger.ConfigJSON)
	if err != nil {
		existingConfig = map[string]any{}
	}
	if isCronTriggerType(trigger.Type) {
		existingConfig = normalizeCronConfigMap(existingConfig)
		configValue = normalizeCronConfigMap(configValue)
	}
	if !jsonMapsEqual(existingConfig, configValue) {
		patch["config"] = configValue
	}
	if len(patch) == 0 {
		return m.pushToast(ToastWarn, "No changes to update")
	}
	client := m.client
	m.mutationPending = true
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		_, err := client.UpdateTrigger(workflowID, triggerID, patch)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Trigger updated", refresh: true}
	}
}

func (m *Model) createTriggerCmd(workflowID string, triggerType string, name string, active bool, configValue map[string]any) tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	workflowID = strings.TrimSpace(workflowID)
	name = strings.TrimSpace(name)
	if workflowID == "" {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	if name == "" {
		return m.pushToast(ToastWarn, "Trigger name cannot be empty")
	}
	if configValue == nil {
		configValue = map[string]any{}
	}
	client := m.client
	m.mutationPending = true
	triggerType = strings.ToUpper(strings.TrimSpace(triggerType))
	return func() tea.Msg {
		if client == nil {
			return mutationResultMsg{err: fmt.Errorf("api client unavailable")}
		}
		payload := map[string]any{
			"type":     triggerType,
			"name":     name,
			"isActive": active,
			"config":   configValue,
		}
		_, err := client.CreateTrigger(workflowID, payload)
		if err != nil {
			return mutationResultMsg{err: err}
		}
		return mutationResultMsg{successMessage: "Trigger created", refresh: true}
	}
}

func (m *Model) openRenameWorkflowModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewWorkflows {
		return m.pushToast(ToastWarn, "Open Workflows to rename")
	}
	selected := m.selectedRowID()
	wf, ok := workflowByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	input := newActionInput("name> ", "Workflow name", wf.Name, 120)
	m.action = actionModalState{
		Active:      true,
		Mode:        actionModalRenameWorkflow,
		Title:       "Rename Workflow",
		Description: "Update selected workflow name",
		Primary:     input,
		Focus:       0,
		WorkflowID:  wf.ID,
	}
	m.syncActionModalFocus()
	return nil
}

func (m *Model) openRenameTriggerModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewTriggers {
		return m.pushToast(ToastWarn, "Open Triggers to update")
	}
	selected := m.selectedRowID()
	trg, ok := triggerByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	nameInput := newActionInput("name> ", "Trigger name", trg.Name, 120)
	configInput := newActionInput("config> ", "JSON object", "{}", 2000)
	timezoneInput := newActionInput("tz> ", "UTC", "", 120)
	triggerType := strings.ToUpper(strings.TrimSpace(trg.Type))
	if triggerType == "CRON" {
		cronExpr, timezone := triggerCronFieldsFromJSON(trg.ConfigJSON)
		configInput = newActionInput("cron> ", "*/5 * * * *", cronExpr, 256)
		timezoneInput = newActionInput("tz> ", "UTC", timezone, 120)
	} else {
		configValue := strings.TrimSpace(trg.ConfigJSON)
		if configValue == "" {
			configValue = "{}"
		}
		configInput.SetValue(configValue)
		configInput.CursorEnd()
	}
	m.action = actionModalState{
		Active:        true,
		Mode:          actionModalUpdateTrigger,
		Title:         "Update Trigger",
		Description:   "Edit trigger fields and config",
		Primary:       nameInput,
		Secondary:     configInput,
		Tertiary:      timezoneInput,
		Focus:         0,
		WorkflowID:    trg.WorkflowID,
		TriggerID:     trg.ID,
		TriggerType:   triggerType,
		TriggerActive: trg.Active,
	}
	m.syncActionModalFocus()
	return nil
}

func (m *Model) openCreateTriggerModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	workflowID := m.selectedWorkflowIDForTriggerMutation()
	if workflowID == "" {
		return m.pushToast(ToastWarn, "Select a workflow in Workflows, or a trigger row in Triggers")
	}
	nameInput := newActionInput("name> ", "Trigger name", "", 120)
	configInput := newActionInput("config> ", "JSON object", "{}", 2000)
	timezoneInput := newActionInput("tz> ", "UTC", "UTC", 120)
	m.action = actionModalState{
		Active:        true,
		Mode:          actionModalCreateTrigger,
		Title:         "Create Trigger",
		Description:   "Create a trigger for selected workflow",
		Primary:       nameInput,
		Secondary:     configInput,
		Tertiary:      timezoneInput,
		Focus:         1,
		WorkflowID:    workflowID,
		TriggerType:   "MANUAL",
		TriggerActive: true,
	}
	m.setActionTriggerType("MANUAL")
	m.syncActionModalFocus()
	return nil
}

func (m *Model) openCreateSecretModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewSecrets {
		return m.pushToast(ToastWarn, "Open Secrets to create")
	}
	nameInput := newActionInput("name> ", "Secret name", "", 120)
	valueInput := newMaskedActionInput("value> ", 5000)
	descriptionInput := newActionInput("desc> ", "Optional description", "", 500)
	m.action = actionModalState{
		Active:      true,
		Mode:        actionModalCreateSecret,
		Title:       "Create Secret",
		Description: "Secret value is masked and never shown in context.",
		Primary:     nameInput,
		Secondary:   valueInput,
		Tertiary:    descriptionInput,
		Focus:       0,
	}
	m.syncActionModalFocus()
	return nil
}

func (m *Model) openUpdateSecretModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewSecrets {
		return m.pushToast(ToastWarn, "Open Secrets to update")
	}
	selected := m.selectedRowID()
	secret, ok := secretByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a secret first")
	}
	nameInput := newActionInput("name> ", "Secret name", secret.Name, 120)
	valueInput := newMaskedActionInput("value> ", 5000)
	descriptionInput := newActionInput("desc> ", "Optional description", secret.Description, 500)
	m.action = actionModalState{
		Active:      true,
		Mode:        actionModalUpdateSecret,
		Title:       "Update Secret",
		Description: "Leave value empty to keep existing secret value.",
		Primary:     nameInput,
		Secondary:   valueInput,
		Tertiary:    descriptionInput,
		Focus:       0,
		SecretID:    secret.ID,
	}
	m.syncActionModalFocus()
	return nil
}

func (m *Model) openDeleteSecretModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewSecrets {
		return m.pushToast(ToastWarn, "Open Secrets to delete")
	}
	selected := m.selectedRowID()
	secret, ok := secretByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a secret first")
	}
	phrase := "DELETE " + secret.ID
	description := "Permanently delete secret \"" + secret.Name + "\" (" + secret.ID + ")"
	m.openDeleteConfirmModal("Delete Secret", description, phrase, "secret", "", "", secret.ID)
	return nil
}

func (m *Model) openDeleteWorkflowModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewWorkflows {
		return m.pushToast(ToastWarn, "Open Workflows to archive")
	}
	selected := m.selectedRowID()
	wf, ok := workflowByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a workflow first")
	}
	if !wf.Active {
		return m.pushToast(ToastInfo, "Workflow already archived; press e to restore")
	}
	phrase := "ARCHIVE " + wf.ID
	description := "Archive workflow \"" + wf.Name + "\" (" + wf.ID + ") and set it inactive"
	m.openDeleteConfirmModal("Archive Workflow", description, phrase, "workflow", wf.ID, "", "")
	return nil
}

func (m *Model) openDeleteTriggerModalCmd() tea.Cmd {
	if m.mutationPending {
		return m.pushToast(ToastWarn, "Another action is still in progress")
	}
	if m.view != ViewTriggers {
		return m.pushToast(ToastWarn, "Open Triggers to archive")
	}
	selected := m.selectedRowID()
	trg, ok := triggerByID(&m.store, selected)
	if !ok {
		return m.pushToast(ToastWarn, "Select a trigger first")
	}
	if !trg.Active {
		return m.pushToast(ToastInfo, "Trigger already archived; press e to restore")
	}
	phrase := "ARCHIVE " + trg.ID
	description := "Archive trigger \"" + trg.Name + "\" (" + trg.ID + ") and set it inactive"
	m.openDeleteConfirmModal("Archive Trigger", description, phrase, "trigger", trg.WorkflowID, trg.ID, "")
	return nil
}

func (m *Model) openCLIHandoffModalCmd(topic string) tea.Cmd {
	description := "Use the CLI for definition authoring"
	command := ""
	title := "CLI Handoff"
	selectedWorkflow := m.selectedWorkflowIDForTriggerMutation()
	if selectedWorkflow == "" && m.view == ViewWorkflows {
		selectedWorkflow = m.selectedRowID()
	}
	workflowTarget := "<workflow-id>"
	if strings.TrimSpace(selectedWorkflow) != "" {
		workflowTarget = selectedWorkflow
	}
	switch strings.TrimSpace(topic) {
	case "workflow-version-create":
		title = "Create Workflow Version"
		description = "Workflow versions stay CLI-first for JSON authoring and validation"
		command = "taskforge workflow version create " + workflowTarget + " --definition ./workflow-definition.json"
	default:
		title = "Create Workflow"
		description = "Workflow creation stays CLI-first for JSON authoring and validation"
		command = "taskforge workflow create --name \"my-workflow\" --definition ./workflow-definition.json"
	}
	m.action = actionModalState{
		Active:      true,
		Mode:        actionModalCLIHandoff,
		Title:       title,
		Description: description,
		CLICommand:  command,
	}
	return nil
}

func newActionInput(prompt string, placeholder string, value string, limit int) textinput.Model {
	input := textinput.New()
	input.Prompt = prompt
	input.Placeholder = placeholder
	input.CharLimit = limit
	input.SetValue(value)
	input.CursorEnd()
	return input
}

func newMaskedActionInput(prompt string, limit int) textinput.Model {
	input := newActionInput(prompt, "", "", limit)
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '*'
	return input
}

func parseJSONObject(raw string) (map[string]any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		trimmed = "{}"
	}
	var parsed any
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, err
	}
	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("not an object")
	}
	return obj, nil
}

func jsonMapsEqual(a map[string]any, b map[string]any) bool {
	left, errLeft := json.Marshal(a)
	right, errRight := json.Marshal(b)
	if errLeft != nil || errRight != nil {
		return false
	}
	return string(left) == string(right)
}

func triggerCronFieldsFromJSON(configJSON string) (string, string) {
	config, err := parseJSONObject(configJSON)
	if err != nil {
		return "*/5 * * * *", "UTC"
	}
	cron := "*/5 * * * *"
	if value, ok := config["cron"].(string); ok && strings.TrimSpace(value) != "" {
		cron = strings.TrimSpace(value)
	}
	timezone := "UTC"
	if value, ok := config["timezone"].(string); ok && strings.TrimSpace(value) != "" {
		timezone = strings.TrimSpace(value)
	}
	return cron, timezone
}

func normalizeCronConfigMap(config map[string]any) map[string]any {
	normalized := map[string]any{}
	for key, value := range config {
		normalized[key] = value
	}
	cron := strings.TrimSpace(stringFromAny(normalized["cron"]))
	if cron != "" {
		normalized["cron"] = cron
	}
	timezone := strings.TrimSpace(stringFromAny(normalized["timezone"]))
	if timezone == "" {
		timezone = "UTC"
	}
	normalized["timezone"] = timezone
	if _, ok := normalized["input"]; !ok {
		normalized["input"] = map[string]any{}
	}
	return normalized
}

func stringFromAny(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return fmt.Sprintf("%v", value)
}

func isCronTriggerType(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "CRON")
}

func (m *Model) setActionTriggerType(triggerType string) {
	next := strings.ToUpper(strings.TrimSpace(triggerType))
	if !isAllowedTriggerType(next) {
		next = "MANUAL"
	}
	previous := strings.ToUpper(strings.TrimSpace(m.action.TriggerType))
	m.action.TriggerType = next
	if isCronTriggerType(next) {
		m.action.Secondary.Prompt = "cron> "
		m.action.Secondary.Placeholder = "*/5 * * * *"
		if strings.TrimSpace(m.action.Secondary.Value()) == "" || strings.TrimSpace(m.action.Secondary.Value()) == "{}" || !isCronTriggerType(previous) {
			m.action.Secondary.SetValue("*/5 * * * *")
			m.action.Secondary.CursorEnd()
		}
		m.action.Tertiary.Prompt = "tz> "
		m.action.Tertiary.Placeholder = "UTC"
		if strings.TrimSpace(m.action.Tertiary.Value()) == "" {
			m.action.Tertiary.SetValue("UTC")
			m.action.Tertiary.CursorEnd()
		}
		return
	}
	m.action.Secondary.Prompt = "config> "
	m.action.Secondary.Placeholder = "JSON object"
	if isCronTriggerType(previous) {
		m.action.Secondary.SetValue("{}")
		m.action.Secondary.CursorEnd()
	}
	m.action.Tertiary.Prompt = "tz> "
	m.action.Tertiary.Placeholder = "UTC"
}

func (m *Model) triggerConfigFromAction() (map[string]any, error) {
	if isCronTriggerType(m.action.TriggerType) {
		cronExpr := strings.TrimSpace(m.action.Secondary.Value())
		timezone := strings.TrimSpace(m.action.Tertiary.Value())
		if timezone == "" {
			timezone = "UTC"
		}
		return map[string]any{
			"cron":     cronExpr,
			"timezone": timezone,
			"input":    map[string]any{},
		}, nil
	}
	return parseJSONObject(m.action.Secondary.Value())
}

func (m *Model) refreshActionValidation() {
	if !m.action.ShowValidation {
		return
	}
	errMessage := m.actionModalValidationError()
	m.action.Validation = errMessage
	if errMessage == "" {
		m.action.ShowValidation = false
	}
}

func (m *Model) actionModalValidationError() string {
	switch m.action.Mode {
	case actionModalRenameWorkflow:
		if strings.TrimSpace(m.action.WorkflowID) == "" {
			return "Select a workflow first"
		}
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Workflow name cannot be empty"
		}
	case actionModalRenameTrigger:
		if strings.TrimSpace(m.action.WorkflowID) == "" || strings.TrimSpace(m.action.TriggerID) == "" {
			return "Select a trigger first"
		}
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Trigger name cannot be empty"
		}
	case actionModalUpdateTrigger:
		if strings.TrimSpace(m.action.WorkflowID) == "" || strings.TrimSpace(m.action.TriggerID) == "" {
			return "Select a trigger first"
		}
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Trigger name cannot be empty"
		}
		if isCronTriggerType(m.action.TriggerType) {
			if len(strings.Fields(strings.TrimSpace(m.action.Secondary.Value()))) != 5 {
				return "Cron expression must have 5 fields"
			}
			if strings.TrimSpace(m.action.Tertiary.Value()) == "" {
				return "Timezone cannot be empty"
			}
		} else {
			if _, err := parseJSONObject(m.action.Secondary.Value()); err != nil {
				return "Trigger config must be a valid JSON object"
			}
		}
	case actionModalCreateTrigger:
		if strings.TrimSpace(m.action.WorkflowID) == "" {
			return "Select a workflow first"
		}
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Trigger name cannot be empty"
		}
		if !isAllowedTriggerType(m.action.TriggerType) {
			return "Trigger type must be MANUAL, CRON, or WEBHOOK"
		}
		if isCronTriggerType(m.action.TriggerType) {
			if len(strings.Fields(strings.TrimSpace(m.action.Secondary.Value()))) != 5 {
				return "Cron expression must have 5 fields"
			}
			if strings.TrimSpace(m.action.Tertiary.Value()) == "" {
				return "Timezone cannot be empty"
			}
		} else {
			if _, err := parseJSONObject(m.action.Secondary.Value()); err != nil {
				return "Trigger config must be a valid JSON object"
			}
		}
	case actionModalCreateSecret:
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Secret name cannot be empty"
		}
		if strings.TrimSpace(m.action.Secondary.Value()) == "" {
			return "Secret value cannot be empty"
		}
	case actionModalUpdateSecret:
		if strings.TrimSpace(m.action.SecretID) == "" {
			return "Select a secret first"
		}
		secret, ok := secretByID(&m.store, m.action.SecretID)
		if !ok {
			return "Select a secret first"
		}
		if strings.TrimSpace(m.action.Primary.Value()) == "" {
			return "Secret name cannot be empty"
		}
		nameChanged := strings.TrimSpace(m.action.Primary.Value()) != secret.Name
		description := strings.TrimSpace(m.action.Tertiary.Value())
		descriptionChanged := description != "" && description != strings.TrimSpace(secret.Description)
		valueChanged := strings.TrimSpace(m.action.Secondary.Value()) != ""
		if !nameChanged && !descriptionChanged && !valueChanged {
			return "No changes to update"
		}
	case actionModalConfirmDelete:
		kind := strings.TrimSpace(strings.ToLower(m.action.DeleteKind))
		switch kind {
		case "workflow":
			if strings.TrimSpace(m.action.WorkflowID) == "" {
				return "Missing delete target"
			}
		case "trigger":
			if strings.TrimSpace(m.action.WorkflowID) == "" || strings.TrimSpace(m.action.TriggerID) == "" {
				return "Missing delete target"
			}
		case "secret":
			if strings.TrimSpace(m.action.SecretID) == "" {
				return "Missing delete target"
			}
		default:
			if strings.TrimSpace(m.action.SecretID) == "" && strings.TrimSpace(m.action.WorkflowID) == "" {
				return "Missing delete target"
			}
		}
		if strings.TrimSpace(m.action.ConfirmPhrase) == "" {
			return "Confirmation phrase is required"
		}
		if strings.TrimSpace(m.action.Confirm.Value()) != strings.TrimSpace(m.action.ConfirmPhrase) {
			return "Type the exact confirmation phrase"
		}
	}
	return ""
}

func isAllowedTriggerType(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "MANUAL", "CRON", "WEBHOOK":
		return true
	default:
		return false
	}
}

func (m *Model) openDeleteConfirmModal(title string, description string, phrase string, deleteKind string, workflowID string, triggerID string, secretID string) {
	confirmInput := newActionInput("confirm> ", "", "", 80)
	m.action = actionModalState{
		Active:        true,
		Mode:          actionModalConfirmDelete,
		Title:         title,
		Description:   description,
		Confirm:       confirmInput,
		Focus:         0,
		DeleteKind:    deleteKind,
		WorkflowID:    workflowID,
		TriggerID:     triggerID,
		SecretID:      secretID,
		ConfirmPhrase: phrase,
	}
	m.syncActionModalFocus()
}

func (m *Model) selectedWorkflowIDForTriggerMutation() string {
	selected := m.selectedRowID()
	if selected == "" {
		return ""
	}
	if m.view == ViewWorkflows {
		if _, ok := workflowByID(&m.store, selected); ok {
			return selected
		}
		return ""
	}
	if m.view == ViewTriggers {
		if trg, ok := triggerByID(&m.store, selected); ok {
			return trg.WorkflowID
		}
	}
	return ""
}

func (m *Model) updateActionModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.action.Active {
		return m, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) || key.Matches(keyMsg, m.keys.Quit) {
			m.action = actionModalState{}
			return m, nil
		}
		if m.action.Mode == actionModalCLIHandoff {
			if key.Matches(keyMsg, m.keys.Enter) {
				m.action = actionModalState{}
			}
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) {
			if errMessage := m.actionModalValidationError(); errMessage != "" {
				m.action.Validation = errMessage
				m.action.ShowValidation = true
				return m, nil
			}
			m.action.Validation = ""
			m.action.ShowValidation = false
			return m, m.submitActionModal()
		}
		if key.Matches(keyMsg, m.keys.NextScreen) {
			m.cycleActionModalFocus(1)
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.PrevScreen) {
			m.cycleActionModalFocus(-1)
			return m, nil
		}
		if m.action.Mode == actionModalCreateTrigger || m.action.Mode == actionModalUpdateTrigger {
			switch keyMsg.String() {
			case "left", "h":
				if m.action.Mode == actionModalCreateTrigger && m.action.Focus == 0 {
					m.cycleActionTriggerType(-1)
					m.refreshActionValidation()
					return m, nil
				}
			case "right", "l":
				if m.action.Mode == actionModalCreateTrigger && m.action.Focus == 0 {
					m.cycleActionTriggerType(1)
					m.refreshActionValidation()
					return m, nil
				}
			case " ":
				toggleFocus := 2
				if m.action.Mode == actionModalUpdateTrigger {
					toggleFocus = 1
				}
				if m.action.Focus == toggleFocus {
					m.action.TriggerActive = !m.action.TriggerActive
					m.refreshActionValidation()
					return m, nil
				}
			}
		}
		if m.action.Mode == actionModalConfirmDelete {
			if key.Matches(keyMsg, m.keys.Clear) {
				m.action.Confirm.SetValue("")
				m.action.Confirm.CursorEnd()
				m.refreshActionValidation()
				return m, nil
			}
		}
		if key.Matches(keyMsg, m.keys.Clear) {
			switch m.action.Mode {
			case actionModalRenameWorkflow, actionModalRenameTrigger:
				m.action.Primary.SetValue("")
				m.action.Primary.CursorEnd()
				m.refreshActionValidation()
				return m, nil
			case actionModalCreateTrigger, actionModalUpdateTrigger:
				nameFocus := 1
				configFocus := 3
				timezoneFocus := 4
				if m.action.Mode == actionModalUpdateTrigger {
					nameFocus = 0
					configFocus = 2
					timezoneFocus = 3
				}
				if m.action.Focus == nameFocus {
					m.action.Primary.SetValue("")
					m.action.Primary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
				if m.action.Focus == configFocus {
					if isCronTriggerType(m.action.TriggerType) {
						m.action.Secondary.SetValue("*/5 * * * *")
					} else {
						m.action.Secondary.SetValue("{}")
					}
					m.action.Secondary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
				if isCronTriggerType(m.action.TriggerType) && m.action.Focus == timezoneFocus {
					m.action.Tertiary.SetValue("UTC")
					m.action.Tertiary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
			case actionModalCreateSecret, actionModalUpdateSecret:
				if m.action.Focus == 0 {
					m.action.Primary.SetValue("")
					m.action.Primary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
				if m.action.Focus == 1 {
					m.action.Secondary.SetValue("")
					m.action.Secondary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
				if m.action.Focus == 2 {
					m.action.Tertiary.SetValue("")
					m.action.Tertiary.CursorEnd()
					m.refreshActionValidation()
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.action.Mode {
	case actionModalRenameWorkflow, actionModalRenameTrigger:
		m.action.Primary, cmd = m.action.Primary.Update(msg)
		m.refreshActionValidation()
		return m, cmd
	case actionModalCreateTrigger, actionModalUpdateTrigger:
		nameFocus := 1
		configFocus := 3
		timezoneFocus := 4
		if m.action.Mode == actionModalUpdateTrigger {
			nameFocus = 0
			configFocus = 2
			timezoneFocus = 3
		}
		if m.action.Focus == nameFocus {
			m.action.Primary, cmd = m.action.Primary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
		if m.action.Focus == configFocus {
			m.action.Secondary, cmd = m.action.Secondary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
		if isCronTriggerType(m.action.TriggerType) && m.action.Focus == timezoneFocus {
			m.action.Tertiary, cmd = m.action.Tertiary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
	case actionModalCreateSecret, actionModalUpdateSecret:
		if m.action.Focus == 0 {
			m.action.Primary, cmd = m.action.Primary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
		if m.action.Focus == 1 {
			m.action.Secondary, cmd = m.action.Secondary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
		if m.action.Focus == 2 {
			m.action.Tertiary, cmd = m.action.Tertiary.Update(msg)
			m.refreshActionValidation()
			return m, cmd
		}
	case actionModalConfirmDelete:
		m.action.Confirm, cmd = m.action.Confirm.Update(msg)
		m.refreshActionValidation()
		return m, cmd
	}
	return m, nil
}

func (m *Model) submitActionModal() tea.Cmd {
	switch m.action.Mode {
	case actionModalRenameWorkflow:
		name := strings.TrimSpace(m.action.Primary.Value())
		if name == "" {
			return m.pushToast(ToastWarn, "Workflow name cannot be empty")
		}
		workflowID := m.action.WorkflowID
		m.action = actionModalState{}
		return m.renameWorkflowCmd(workflowID, name)
	case actionModalRenameTrigger, actionModalUpdateTrigger:
		name := strings.TrimSpace(m.action.Primary.Value())
		workflowID := m.action.WorkflowID
		triggerID := m.action.TriggerID
		active := m.action.TriggerActive
		configValue, err := m.triggerConfigFromAction()
		if err != nil {
			return m.pushToast(ToastWarn, "Trigger config must be a valid JSON object")
		}
		m.action = actionModalState{}
		return m.updateTriggerCmd(workflowID, triggerID, name, active, configValue)
	case actionModalCreateTrigger:
		name := strings.TrimSpace(m.action.Primary.Value())
		configValue, err := m.triggerConfigFromAction()
		if err != nil {
			return m.pushToast(ToastWarn, "Trigger config must be a valid JSON object")
		}
		workflowID := m.action.WorkflowID
		triggerType := m.action.TriggerType
		active := m.action.TriggerActive
		m.action = actionModalState{}
		return m.createTriggerCmd(workflowID, triggerType, name, active, configValue)
	case actionModalCreateSecret:
		name := strings.TrimSpace(m.action.Primary.Value())
		value := m.action.Secondary.Value()
		description := strings.TrimSpace(m.action.Tertiary.Value())
		m.action = actionModalState{}
		return m.createSecretCmd(name, value, description)
	case actionModalUpdateSecret:
		secretID := strings.TrimSpace(m.action.SecretID)
		name := strings.TrimSpace(m.action.Primary.Value())
		value := m.action.Secondary.Value()
		description := strings.TrimSpace(m.action.Tertiary.Value())
		m.action = actionModalState{}
		return m.updateSecretCmd(secretID, name, value, description)
	case actionModalConfirmDelete:
		if errMessage := m.actionModalValidationError(); errMessage != "" {
			m.action.Validation = errMessage
			m.action.ShowValidation = true
			return nil
		}
		kind := strings.TrimSpace(strings.ToLower(m.action.DeleteKind))
		workflowID := m.action.WorkflowID
		triggerID := m.action.TriggerID
		secretID := m.action.SecretID
		m.action = actionModalState{}
		if kind == "secret" {
			return m.deleteSecretCmd(secretID)
		}
		if strings.TrimSpace(triggerID) != "" {
			return m.deleteTriggerCmd(workflowID, triggerID)
		}
		return m.deleteWorkflowCmd(workflowID)
	case actionModalCLIHandoff:
		m.action = actionModalState{}
		return nil
	default:
		m.action = actionModalState{}
		return nil
	}
}

func (m *Model) cycleActionModalFocus(delta int) {
	total := 0
	switch m.action.Mode {
	case actionModalRenameWorkflow, actionModalRenameTrigger:
		total = 1
	case actionModalCreateTrigger:
		total = 4
		if isCronTriggerType(m.action.TriggerType) {
			total = 5
		}
	case actionModalUpdateTrigger:
		total = 3
		if isCronTriggerType(m.action.TriggerType) {
			total = 4
		}
	case actionModalCreateSecret, actionModalUpdateSecret:
		total = 3
	case actionModalConfirmDelete:
		total = 1
	default:
		total = 0
	}
	if total <= 1 {
		return
	}
	next := m.action.Focus + delta
	for next < 0 {
		next += total
	}
	m.action.Focus = next % total
	m.syncActionModalFocus()
}

func (m *Model) syncActionModalFocus() {
	m.action.Primary.Blur()
	m.action.Secondary.Blur()
	m.action.Tertiary.Blur()
	m.action.Confirm.Blur()
	switch m.action.Mode {
	case actionModalRenameWorkflow, actionModalRenameTrigger:
		m.action.Primary.Focus()
	case actionModalCreateTrigger:
		if m.action.Focus == 1 {
			m.action.Primary.Focus()
		}
		if m.action.Focus == 3 {
			m.action.Secondary.Focus()
		}
		if isCronTriggerType(m.action.TriggerType) && m.action.Focus == 4 {
			m.action.Tertiary.Focus()
		}
	case actionModalUpdateTrigger:
		if m.action.Focus == 0 {
			m.action.Primary.Focus()
		}
		if m.action.Focus == 2 {
			m.action.Secondary.Focus()
		}
		if isCronTriggerType(m.action.TriggerType) && m.action.Focus == 3 {
			m.action.Tertiary.Focus()
		}
	case actionModalCreateSecret, actionModalUpdateSecret:
		if m.action.Focus == 0 {
			m.action.Primary.Focus()
		}
		if m.action.Focus == 1 {
			m.action.Secondary.Focus()
		}
		if m.action.Focus == 2 {
			m.action.Tertiary.Focus()
		}
	case actionModalConfirmDelete:
		m.action.Confirm.Focus()
	}
}

func (m *Model) cycleActionTriggerType(delta int) {
	order := []string{"MANUAL", "CRON", "WEBHOOK"}
	index := 0
	for i, item := range order {
		if strings.EqualFold(item, m.action.TriggerType) {
			index = i
			break
		}
	}
	next := index + delta
	for next < 0 {
		next += len(order)
	}
	m.setActionTriggerType(order[next%len(order)])
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

func buildPalette(theme styles.Theme, recentActions []paletteAction, state paletteBuildState) list.Model {
	section := func(label string) paletteItem {
		return paletteItem{Label: label, Section: true, Enabled: false, Action: paletteAction{Kind: paletteNoop}}
	}
	command := func(label string, detail string, action paletteAction, keywords ...string) paletteItem {
		return paletteItem{Label: label, Detail: detail, Action: action, Enabled: true, Keywords: keywords}
	}

	autoStatus := "OFF"
	if state.AutoRefresh {
		autoStatus = "ON"
	}

	runSelected := command("Action: Run selected workflow", "Workflow", paletteAction{Kind: paletteRunWorkflow}, "run", "workflow", "queue")
	if !(state.View == ViewWorkflows && state.HasSelection) {
		runSelected.Enabled = false
		runSelected.Detail = "Unavailable: select workflow row"
		runSelected.DisabledReason = "Select a workflow row in Workflows first"
	}

	renameWorkflow := command("Action: Rename selected workflow", "Workflow", paletteAction{Kind: paletteRenameWorkflow}, "rename", "workflow", "name")
	if !(state.View == ViewWorkflows && state.HasSelection) {
		renameWorkflow.Enabled = false
		renameWorkflow.Detail = "Unavailable: select workflow row"
		renameWorkflow.DisabledReason = "Select a workflow row in Workflows first"
	}

	createTrigger := command("Action: Create trigger", "Trigger", paletteAction{Kind: paletteCreateTrigger}, "create", "trigger", "workflow")
	if !((state.View == ViewWorkflows || state.View == ViewTriggers) && state.HasSelection) {
		createTrigger.Enabled = false
		createTrigger.Detail = "Unavailable: select workflow or trigger row"
		createTrigger.DisabledReason = "Select a row in Workflows or Triggers first"
	}

	renameTrigger := command("Action: Update selected trigger", "Trigger", paletteAction{Kind: paletteRenameTrigger}, "update", "trigger", "name", "config")
	if !(state.View == ViewTriggers && state.HasSelection) {
		renameTrigger.Enabled = false
		renameTrigger.Detail = "Unavailable: select trigger row"
		renameTrigger.DisabledReason = "Select a trigger row in Triggers first"
	}

	toggleTrigger := command("Action: Toggle selected trigger", "Trigger", paletteAction{Kind: paletteToggleTrigger}, "toggle", "trigger", "active")
	if !(state.View == ViewTriggers && state.HasSelection) {
		toggleTrigger.Enabled = false
		toggleTrigger.Detail = "Unavailable: select trigger row"
		toggleTrigger.DisabledReason = "Select a trigger row in Triggers first"
	}

	deleteWorkflow := command("Action: Archive selected workflow", "Workflow", paletteAction{Kind: paletteDeleteWorkflow}, "archive", "workflow", "soft-delete")
	if !(state.View == ViewWorkflows && state.HasSelection) {
		deleteWorkflow.Enabled = false
		deleteWorkflow.Detail = "Unavailable: select workflow row"
		deleteWorkflow.DisabledReason = "Select a workflow row in Workflows first"
	}

	deleteTrigger := command("Action: Archive selected trigger", "Trigger", paletteAction{Kind: paletteDeleteTrigger}, "archive", "trigger", "soft-delete")
	if !(state.View == ViewTriggers && state.HasSelection) {
		deleteTrigger.Enabled = false
		deleteTrigger.Detail = "Unavailable: select trigger row"
		deleteTrigger.DisabledReason = "Select a trigger row in Triggers first"
	}

	createSecret := command("Action: Create secret", "Secret", paletteAction{Kind: paletteCreateSecret}, "create", "secret", "credential")
	if state.View != ViewSecrets {
		createSecret.Enabled = false
		createSecret.Detail = "Unavailable in this view"
		createSecret.DisabledReason = "Open Secrets first"
	}

	updateSecret := command("Action: Update selected secret", "Secret", paletteAction{Kind: paletteUpdateSecret}, "update", "secret", "edit")
	if !(state.View == ViewSecrets && state.HasSelection) {
		updateSecret.Enabled = false
		updateSecret.Detail = "Unavailable: select secret row"
		updateSecret.DisabledReason = "Select a secret row in Secrets first"
	}

	deleteSecret := command("Action: Delete selected secret", "Secret", paletteAction{Kind: paletteDeleteSecret}, "delete", "secret", "hard-delete")
	if !(state.View == ViewSecrets && state.HasSelection) {
		deleteSecret.Enabled = false
		deleteSecret.Detail = "Unavailable: select secret row"
		deleteSecret.DisabledReason = "Select a secret row in Secrets first"
	}

	showAllScope := command("Filter: Show all", "Status scope", paletteAction{Kind: paletteSetStatusScope, Value: "all"}, "filter", "status", "all")
	showActiveScope := command("Filter: Active only", "Status scope", paletteAction{Kind: paletteSetStatusScope, Value: "active"}, "filter", "status", "active")
	showInactiveScope := command("Filter: Inactive only", "Status scope", paletteAction{Kind: paletteSetStatusScope, Value: "inactive"}, "filter", "status", "inactive")
	if !state.HasScope {
		showAllScope.Enabled = false
		showAllScope.Detail = "Unavailable in this view"
		showAllScope.DisabledReason = "Status scope is available in Workflows and Triggers"
		showActiveScope.Enabled = false
		showActiveScope.Detail = "Unavailable in this view"
		showActiveScope.DisabledReason = "Status scope is available in Workflows and Triggers"
		showInactiveScope.Enabled = false
		showInactiveScope.Detail = "Unavailable in this view"
		showInactiveScope.DisabledReason = "Status scope is available in Workflows and Triggers"
	}
	if state.Scope == statusScopeAll {
		showAllScope.Label += " (active)"
	}
	if state.Scope == statusScopeActive {
		showActiveScope.Label += " (active)"
	}
	if state.Scope == statusScopeInactive {
		showInactiveScope.Label += " (active)"
	}

	clearFilters := command("Action: Clear filters", "Table", paletteAction{Kind: paletteClearFilters}, "clear", "filter", "reset")
	if !state.HasFilter {
		clearFilters.Enabled = false
		clearFilters.Detail = "Unavailable: no active filters"
		clearFilters.DisabledReason = "No search filter is active"
	}

	clearRecent := command("Action: Clear recent commands", "System", paletteAction{Kind: paletteClearRecent}, "recent", "history", "clear")
	if !state.HasRecent {
		clearRecent.Enabled = false
		clearRecent.Detail = "Unavailable: no recent commands"
		clearRecent.DisabledReason = "Run commands from the palette first"
	}

	profileItem := func(label string, profile NetworkProfile, active bool) paletteItem {
		detail := "Network"
		if active {
			detail = "Network (active)"
		}
		return command(label, detail, paletteAction{Kind: paletteSetNetworkProfile, Profile: profile}, "network", "latency", "profile")
	}

	items := []list.Item{
		section(":: Navigation"),
		command("Go: Dashboard", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewDashboard}, "dash", "home", "overview"),
		command("Go: Workflows", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewWorkflows}, "wf", "workflow"),
		command("Go: Runs", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewRuns}, "run", "jobs"),
		command("Go: Triggers", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewTriggers}, "trigger"),
		command("Go: Events", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewEvents}, "event", "webhook"),
		command("Go: Secrets", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewSecrets}, "secret", "vault"),
		command("Go: API Tokens", "Navigation", paletteAction{Kind: paletteGoToView, View: ViewTokens}, "token", "auth", "api"),
		section(":: Actions"),
		runSelected,
		renameWorkflow,
		createTrigger,
		renameTrigger,
		toggleTrigger,
		deleteWorkflow,
		deleteTrigger,
		createSecret,
		updateSecret,
		deleteSecret,
		showAllScope,
		showActiveScope,
		showInactiveScope,
		clearFilters,
		command("Toggle: Auto refresh", "System ("+autoStatus+")", paletteAction{Kind: paletteToggleRefresh}, "refresh", "polling", "live"),
		section(":: CLI Handoff"),
		command("CLI: Create workflow", "Workflow authoring", paletteAction{Kind: paletteShowCLIHandoff, Value: "workflow-create"}, "cli", "workflow", "create", "definition"),
		command("CLI: Create workflow version", "Version authoring", paletteAction{Kind: paletteShowCLIHandoff, Value: "workflow-version-create"}, "cli", "version", "definition"),
		clearRecent,
		section(":: Network"),
		profileItem("Network: Fast", NetworkFast, state.Profile == NetworkFast),
		profileItem("Network: Normal", NetworkNormal, state.Profile == NetworkNormal),
		profileItem("Network: Slow", NetworkSlow, state.Profile == NetworkSlow),
		profileItem("Network: Flaky", NetworkFlaky, state.Profile == NetworkFlaky),
		section(":: Themes"),
		command("Theme: Catppuccin", "Theme", paletteAction{Kind: paletteSetTheme, View: ViewID("catppuccin")}, "theme", "pastel"),
		command("Theme: Tokyo Night", "Theme", paletteAction{Kind: paletteSetTheme, View: ViewID("tokyo-night")}, "theme", "blue"),
		command("Theme: Fallout (CRT)", "Theme", paletteAction{Kind: paletteSetTheme, View: ViewID("fallout")}, "theme", "crt", "green"),
		command("Theme: Retro Amber", "Theme", paletteAction{Kind: paletteSetTheme, View: ViewID("retro-amber")}, "theme", "amber", "crt"),
	}
	if len(recentActions) > 0 {
		recent := make([]list.Item, 0, len(recentActions)+1)
		recent = append(recent, section(":: Recent"))
		for _, action := range recentActions {
			recent = append(recent, paletteItemFromAction(action, state))
		}
		items = append(recent, items...)
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
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
			rows[i][c] = ansi.Truncate(rows[i][c], widths[c], "")
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
