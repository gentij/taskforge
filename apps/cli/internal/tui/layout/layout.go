package layout

type Layout struct {
	Width  int
	Height int

	FooterHeight  int
	SidebarWidth  int
	SidebarHeight int
	MainWidth     int
	MainHeight    int
	MainHeader    int
	ContentHeight int
	ContextHeight int

	DashboardCardsHeight int
	DashboardLabelHeight int
	DashboardGap         int
	PrimaryTableHeight   int
}

func Compute(width int, height int, contextCollapsed bool) Layout {
	footerHeight := 1
	if width < 1 {
		width = 1
	}
	if height < footerHeight+1 {
		height = footerHeight + 1
	}

	sidebarWidth := 22
	if width < 80 {
		sidebarWidth = 18
	}
	if width < 60 {
		sidebarWidth = 14
	}
	mainWidth := width - sidebarWidth
	if mainWidth < 20 {
		mainWidth = 20
		sidebarWidth = max(width-mainWidth, 12)
	}

	mainHeight := height - footerHeight
	if mainHeight < 1 {
		mainHeight = 1
	}
	innerHeight := mainHeight - 2
	if innerHeight < 1 {
		innerHeight = 1
	}

	mainHeader := 3
	contextHeight := 0
	if !contextCollapsed {
		target := int(float64(innerHeight) * 0.28)
		if target < 6 {
			target = 6
		}
		maxContext := int(float64(innerHeight) * 0.5)
		if target > maxContext {
			target = maxContext
		}
		if target > innerHeight-3 {
			target = innerHeight - 3
		}
		if target < 0 {
			target = 0
		}
		contextHeight = target
	}
	contentHeight := innerHeight - mainHeader - contextHeight
	if contentHeight < 3 {
		contentHeight = 3
	}

	cardsHeight := 5
	labelHeight := 1
	gap := 1
	primaryTableHeight := contentHeight - 2
	if contentHeight > cardsHeight+labelHeight+gap+2 {
		primaryTableHeight = contentHeight - cardsHeight - labelHeight - gap - 2
	} else {
		cardsHeight = 0
		labelHeight = 0
		gap = 0
		primaryTableHeight = contentHeight - 2
	}

	if primaryTableHeight < 1 {
		primaryTableHeight = 1
	}

	return Layout{
		Width:                width,
		Height:               height,
		FooterHeight:         footerHeight,
		SidebarWidth:         sidebarWidth,
		SidebarHeight:        mainHeight,
		MainWidth:            mainWidth,
		MainHeight:           mainHeight,
		MainHeader:           mainHeader,
		ContentHeight:        contentHeight,
		ContextHeight:        contextHeight,
		DashboardCardsHeight: cardsHeight,
		DashboardLabelHeight: labelHeight,
		DashboardGap:         gap,
		PrimaryTableHeight:   primaryTableHeight,
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
