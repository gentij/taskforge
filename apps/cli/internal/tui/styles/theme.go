package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Background lipgloss.Color
	Surface    lipgloss.Color
	SurfaceAlt lipgloss.Color
	Border     lipgloss.Color
	Text       lipgloss.Color
	Muted      lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Info       lipgloss.Color
	SuccessBg  lipgloss.Color
	WarningBg  lipgloss.Color
	ErrorBg    lipgloss.Color
	InfoBg     lipgloss.Color
	MutedBg    lipgloss.Color
	CRT        bool
	Scanline   lipgloss.Color
	Glow       lipgloss.Color
}

type StyleSet struct {
	Header           lipgloss.Style
	Footer           lipgloss.Style
	PanelBorder      lipgloss.Style
	PanelBorderFocus lipgloss.Style
	PanelTitle       lipgloss.Style
	BorderColor      lipgloss.Color
	Accent           lipgloss.Style
	SidebarTitle     lipgloss.Style
	SidebarMuted     lipgloss.Style
	SidebarSection   lipgloss.Style
	SidebarFill      lipgloss.Style
	Chip             lipgloss.Style
	ChipActive       lipgloss.Style
	TabActive        lipgloss.Style
	TabInactive      lipgloss.Style
	Divider          lipgloss.Style
	PanelFill        lipgloss.Style
	ContextFill      lipgloss.Style
	TableHeader      lipgloss.Style
	TableCell        lipgloss.Style
	TableSelected    lipgloss.Style
	RowAlt           lipgloss.Style
	CardTitle        lipgloss.Style
	CardValue        lipgloss.Style
	BadgeSuccess     lipgloss.Style
	BadgeFailed      lipgloss.Style
	BadgeRunning     lipgloss.Style
	BadgeQueued      lipgloss.Style
	BadgeMuted       lipgloss.Style
	Dim              lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		Name:       "Taskforge",
		Background: lipgloss.Color("#0B0F14"),
		Surface:    lipgloss.Color("#111827"),
		SurfaceAlt: lipgloss.Color("#0F172A"),
		Border:     lipgloss.Color("#1F2937"),
		Text:       lipgloss.Color("#E5E7EB"),
		Muted:      lipgloss.Color("#94A3B8"),
		Accent:     lipgloss.Color("#38BDF8"),
		Success:    lipgloss.Color("#22C55E"),
		Warning:    lipgloss.Color("#F59E0B"),
		Error:      lipgloss.Color("#EF4444"),
		Info:       lipgloss.Color("#60A5FA"),
		SuccessBg:  lipgloss.Color("#0B2E1A"),
		WarningBg:  lipgloss.Color("#2B1A05"),
		ErrorBg:    lipgloss.Color("#2A0F14"),
		InfoBg:     lipgloss.Color("#0B1F33"),
		MutedBg:    lipgloss.Color("#111827"),
		CRT:        false,
		Scanline:   lipgloss.Color("#0B0F14"),
		Glow:       lipgloss.Color("#38BDF8"),
	}
}

func ThemeRegistry() map[string]Theme {
	return map[string]Theme{
		"taskforge": DefaultTheme(),
		"tokyo-night": {
			Name:       "Tokyo Night",
			Background: lipgloss.Color("#0B0F1A"),
			Surface:    lipgloss.Color("#151A2C"),
			SurfaceAlt: lipgloss.Color("#101522"),
			Border:     lipgloss.Color("#262B3F"),
			Text:       lipgloss.Color("#DDE3F0"),
			Muted:      lipgloss.Color("#8A93A6"),
			Accent:     lipgloss.Color("#7AA2F7"),
			Success:    lipgloss.Color("#9ECE6A"),
			Warning:    lipgloss.Color("#E0AF68"),
			Error:      lipgloss.Color("#F7768E"),
			Info:       lipgloss.Color("#7DCFFF"),
			SuccessBg:  lipgloss.Color("#1F2A1A"),
			WarningBg:  lipgloss.Color("#2A1F10"),
			ErrorBg:    lipgloss.Color("#2A151A"),
			InfoBg:     lipgloss.Color("#10233A"),
			MutedBg:    lipgloss.Color("#151A2C"),
		},
		"catppuccin": {
			Name:       "Catppuccin",
			Background: lipgloss.Color("#0F0E17"),
			Surface:    lipgloss.Color("#1C1B24"),
			SurfaceAlt: lipgloss.Color("#16151D"),
			Border:     lipgloss.Color("#2A2934"),
			Text:       lipgloss.Color("#E6E1F0"),
			Muted:      lipgloss.Color("#A59BBE"),
			Accent:     lipgloss.Color("#8AADF4"),
			Success:    lipgloss.Color("#A6DA95"),
			Warning:    lipgloss.Color("#EED49F"),
			Error:      lipgloss.Color("#ED8796"),
			Info:       lipgloss.Color("#91D7E3"),
			SuccessBg:  lipgloss.Color("#1A2418"),
			WarningBg:  lipgloss.Color("#2A2216"),
			ErrorBg:    lipgloss.Color("#2A151A"),
			InfoBg:     lipgloss.Color("#13252B"),
			MutedBg:    lipgloss.Color("#1C1B24"),
			CRT:        false,
			Scanline:   lipgloss.Color("#0F0E17"),
			Glow:       lipgloss.Color("#8AADF4"),
		},
		"fallout": {
			Name:       "Fallout",
			Background: lipgloss.Color("#040804"),
			Surface:    lipgloss.Color("#091109"),
			SurfaceAlt: lipgloss.Color("#0C170C"),
			Border:     lipgloss.Color("#3A8F4A"),
			Text:       lipgloss.Color("#C9F7BE"),
			Muted:      lipgloss.Color("#86C884"),
			Accent:     lipgloss.Color("#8CFF78"),
			Success:    lipgloss.Color("#8CFF78"),
			Warning:    lipgloss.Color("#F3DB5D"),
			Error:      lipgloss.Color("#FF7A7A"),
			Info:       lipgloss.Color("#9CEF95"),
			SuccessBg:  lipgloss.Color("#113015"),
			WarningBg:  lipgloss.Color("#31290B"),
			ErrorBg:    lipgloss.Color("#311515"),
			InfoBg:     lipgloss.Color("#123016"),
			MutedBg:    lipgloss.Color("#0D190D"),
			CRT:        true,
			Scanline:   lipgloss.Color("#0C1A0C"),
			Glow:       lipgloss.Color("#A8FF9A"),
		},
		"retro-amber": {
			Name:       "Retro Amber",
			Background: lipgloss.Color("#120A05"),
			Surface:    lipgloss.Color("#1A0F08"),
			SurfaceAlt: lipgloss.Color("#140B06"),
			Border:     lipgloss.Color("#8B5C24"),
			Text:       lipgloss.Color("#FCE9C1"),
			Muted:      lipgloss.Color("#C98D4A"),
			Accent:     lipgloss.Color("#F2B24C"),
			Success:    lipgloss.Color("#F2B24C"),
			Warning:    lipgloss.Color("#F59E0B"),
			Error:      lipgloss.Color("#FF9A6B"),
			Info:       lipgloss.Color("#F2B24C"),
			SuccessBg:  lipgloss.Color("#2B1A08"),
			WarningBg:  lipgloss.Color("#2B1A05"),
			ErrorBg:    lipgloss.Color("#2A1208"),
			InfoBg:     lipgloss.Color("#2B1A08"),
			MutedBg:    lipgloss.Color("#1A0F08"),
			CRT:        false,
			Scanline:   lipgloss.Color("#120A05"),
			Glow:       lipgloss.Color("#F2B24C"),
		},
	}
}

func NewStyles(theme Theme) StyleSet {
	borderShape := lipgloss.RoundedBorder()
	focusBorderShape := lipgloss.ThickBorder()
	titleColor := theme.Text
	accentColor := theme.Accent
	chipActiveForeground := theme.Background
	chipActiveBackground := theme.Accent
	tableHeaderColor := theme.Muted
	tableSelectedForeground := theme.Text
	tableSelectedBackground := theme.SurfaceAlt
	rowAltBackground := theme.MutedBg
	panelFillBackground := theme.Surface
	contextFillBackground := theme.SurfaceAlt
	tableSelectedUnderline := false
	dimFaint := true

	if theme.CRT {
		borderShape = lipgloss.NormalBorder()
		focusBorderShape = lipgloss.NormalBorder()
		titleColor = theme.Glow
		accentColor = theme.Glow
		chipActiveForeground = theme.Surface
		chipActiveBackground = theme.Glow
		tableHeaderColor = theme.Glow
		tableSelectedForeground = theme.Glow
		tableSelectedBackground = theme.SuccessBg
		tableSelectedUnderline = true
		rowAltBackground = theme.SurfaceAlt
		panelFillBackground = theme.SurfaceAlt
		contextFillBackground = theme.Surface
		dimFaint = false
	}

	return StyleSet{
		Header: lipgloss.NewStyle().
			Background(theme.SurfaceAlt).
			Foreground(theme.Text),
		Footer: lipgloss.NewStyle().
			Background(theme.SurfaceAlt).
			Foreground(theme.Muted),
		PanelBorder: lipgloss.NewStyle().
			Border(borderShape).
			BorderForeground(theme.Muted).
			Background(theme.Surface),
		PanelBorderFocus: lipgloss.NewStyle().
			Border(focusBorderShape).
			BorderForeground(accentColor).
			Bold(theme.CRT).
			Background(theme.Surface),
		PanelTitle: lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true),
		BorderColor: theme.Border,
		Accent: lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true),
		SidebarTitle: lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true),
		SidebarMuted: lipgloss.NewStyle().
			Foreground(theme.Muted),
		SidebarSection: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Bold(true),
		SidebarFill: lipgloss.NewStyle().
			Background(theme.SurfaceAlt),
		Chip: lipgloss.NewStyle().
			Foreground(theme.Text).
			Background(theme.SurfaceAlt).
			Padding(0, 1),
		ChipActive: lipgloss.NewStyle().
			Foreground(chipActiveForeground).
			Background(chipActiveBackground).
			Padding(0, 1).
			Bold(true),
		TabActive: lipgloss.NewStyle().
			Foreground(chipActiveForeground).
			Background(chipActiveBackground).
			Bold(true),
		TabInactive: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.SurfaceAlt),
		Divider: lipgloss.NewStyle().
			Foreground(theme.Border),
		PanelFill: lipgloss.NewStyle().
			Background(panelFillBackground),
		ContextFill: lipgloss.NewStyle().
			Background(contextFillBackground),
		TableHeader: lipgloss.NewStyle().
			Foreground(tableHeaderColor).
			Bold(true),
		TableCell: lipgloss.NewStyle().
			Foreground(theme.Text),
		TableSelected: lipgloss.NewStyle().
			Foreground(tableSelectedForeground).
			Background(tableSelectedBackground).
			Underline(tableSelectedUnderline).
			Bold(true),
		RowAlt: lipgloss.NewStyle().
			Background(rowAltBackground),
		CardTitle: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Bold(true),
		CardValue: lipgloss.NewStyle().
			Foreground(theme.Text).
			Bold(true),
		BadgeSuccess: lipgloss.NewStyle().
			Foreground(theme.Success).
			Background(theme.SuccessBg).
			Bold(true),
		BadgeFailed: lipgloss.NewStyle().
			Foreground(theme.Error).
			Background(theme.ErrorBg).
			Bold(true),
		BadgeRunning: lipgloss.NewStyle().
			Foreground(theme.Info).
			Background(theme.InfoBg).
			Bold(true),
		BadgeQueued: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Background(theme.WarningBg).
			Bold(true),
		BadgeMuted: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.MutedBg).
			Bold(true),
		Dim: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Faint(dimFaint),
	}
}

func Badge(style lipgloss.Style, text string) string {
	return style.Render("[" + text + "]")
}
