package layout

type Layout struct {
	Width  int
	Height int

	HeaderHeight int
	FooterHeight int

	ContextHeight int
	PrimaryHeight int

	DashboardCardsHeight int
	DashboardLabelHeight int
	DashboardGap         int
	PrimaryTableHeight   int
}

func Compute(width int, height int, contextCollapsed bool) Layout {
	headerHeight := 2
	footerHeight := 1
	if width < 1 {
		width = 1
	}
	if height < headerHeight+footerHeight+1 {
		height = headerHeight + footerHeight + 1
	}

	available := height - headerHeight - footerHeight
	contextHeight := 0
	if !contextCollapsed {
		target := int(float64(available) * 0.28)
		if target < 6 {
			target = 6
		}
		maxContext := int(float64(available) * 0.5)
		if target > maxContext {
			target = maxContext
		}
		if target > available-4 {
			target = available - 4
		}
		if target < 0 {
			target = 0
		}
		contextHeight = target
	}
	primaryHeight := available - contextHeight
	if primaryHeight < 1 {
		primaryHeight = 1
	}

	cardsHeight := 5
	labelHeight := 1
	gap := 1
	primaryTableHeight := primaryHeight - 1
	if primaryHeight > cardsHeight+labelHeight+gap+2 {
		primaryTableHeight = primaryHeight - cardsHeight - labelHeight - gap - 1
	} else {
		cardsHeight = 0
		labelHeight = 0
		gap = 0
		primaryTableHeight = primaryHeight - 1
	}

	if primaryTableHeight < 1 {
		primaryTableHeight = 1
	}

	return Layout{
		Width:                width,
		Height:               height,
		HeaderHeight:         headerHeight,
		FooterHeight:         footerHeight,
		ContextHeight:        contextHeight,
		PrimaryHeight:        primaryHeight,
		DashboardCardsHeight: cardsHeight,
		DashboardLabelHeight: labelHeight,
		DashboardGap:         gap,
		PrimaryTableHeight:   primaryTableHeight,
	}
}
