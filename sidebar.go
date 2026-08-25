package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	navIconWidth          = 2
	navGapWidth           = 1
	sidebarExpandedInner  = 12
	sidebarCollapsedInner = 2
	sidebarPadX           = 1
	sidebarBorderW        = 1
)

func sidebarInnerWidth(collapsed bool) int {
	if collapsed {
		return sidebarCollapsedInner
	}
	return sidebarExpandedInner
}

func sidebarFrameWidth(collapsed bool) int {
	return sidebarInnerWidth(collapsed) + sidebarPadX*2 + sidebarBorderW
}

func styleSidebar(collapsed bool) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(colorBorder)).
		Background(lipgloss.Color(colorBg)).
		Padding(0, sidebarPadX).
		Width(sidebarFrameWidth(collapsed))
}

func (m appModel) sidebarView(height int) string {
	st := styleSidebar(m.navCollapsed)
	innerH := max(0, height-st.GetVerticalFrameSize())
	innerW := max(0, sidebarInnerWidth(m.navCollapsed))

	var rows []string
	for _, page := range navPages() {
		rows = append(rows, m.navRow(page, innerW))
	}
	nav := lipgloss.JoinVertical(lipgloss.Left, rows...)

	dock := m.sidebarDock(innerW)
	navSlot := max(0, innerH-lipgloss.Height(dock))
	nav = lipgloss.NewStyle().
		Width(innerW).
		Height(navSlot).
		Render(nav)

	inner := lipgloss.JoinVertical(lipgloss.Left, nav, dock)
	return st.Height(max(0, height)).Render(inner)
}

func (m appModel) navRow(page Page, width int) string {
	active := page == m.page
	hover := m.hovering("nav-" + string(page))

	var row string
	if m.navCollapsed {
		row = styleNavIcon(active, hover).
			Width(width).
			Render(page.Icon())
	} else {
		labelW := max(0, width-navIconWidth-navGapWidth)
		icon := styleNavIcon(active, hover).Render(page.Icon())
		gap := styleNavGap(active, hover).Render(" ")
		label := styleNavLabel(active, hover).
			Width(labelW).
			Height(1).
			Align(lipgloss.Left, lipgloss.Center).
			Render(page.NavLabel())
		row = lipgloss.JoinHorizontal(lipgloss.Center, icon, gap, label)
	}

	row = lipgloss.NewStyle().Width(width).MarginBottom(1).Render(row)
	if m.ctx.zone != nil {
		row = m.ctx.zone.Mark("nav-"+string(page), row)
	}
	return row
}

func (m appModel) sidebarDock(width int) string {
	if width <= 0 {
		return ""
	}
	rule := styleFaint().Width(width).Render(strings.Repeat("─", width))
	picker := m.themePicker(width)
	toggle := m.navToggle(width)
	gap := lipgloss.NewStyle().Width(width).Height(1).Render(" ")
	return lipgloss.JoinVertical(lipgloss.Left, rule, picker, gap, toggle)
}

func (m appModel) navToggle(width int) string {
	glyph := "<"
	if m.navCollapsed {
		glyph = ">"
	}
	st := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color(colorFaint))
	if m.hovering("nav-toggle") {
		st = st.
			Foreground(lipgloss.Color(colorBright)).
			Background(lipgloss.Color(colorSurface))
	}
	row := st.Render(glyph)
	if m.ctx.zone != nil {
		row = m.ctx.zone.Mark("nav-toggle", row)
	}
	return row
}

func (m appModel) themePicker(width int) string {
	if width <= 0 {
		return ""
	}
	t := currentTheme()

	if m.navCollapsed {
		st := lipgloss.NewStyle().
			Width(width).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color(t.Accent))
		if m.hovering("theme-cycle") {
			st = st.Background(lipgloss.Color(colorSurface))
		}
		swatch := st.Render("◆")
		if m.ctx.zone != nil {
			swatch = m.ctx.zone.Mark("theme-cycle", swatch)
		}
		return swatch
	}

	nameSt := lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
	if m.hovering("theme-cycle") {
		nameSt = nameSt.
			Foreground(lipgloss.Color(colorBg)).
			Background(lipgloss.Color(colorAccent))
	}
	name := nameSt.Render(strings.ToLower(t.Name))
	if m.ctx.zone != nil {
		name = m.ctx.zone.Mark("theme-cycle", name)
	}

	label := styleFaint().Render("theme")
	swatches := m.themeSwatches(width)
	return lipgloss.NewStyle().Width(width).Render(
		lipgloss.JoinVertical(lipgloss.Left, label, name, swatches),
	)
}

func (m appModel) themeSwatches(width int) string {
	themes := allThemes()
	const cellW = 2
	perRow := max(1, width/cellW)
	t := currentTheme()
	var rows []string
	for i := 0; i < len(themes); i += perRow {
		end := min(i+perRow, len(themes))
		var cells []string
		for _, th := range themes[i:end] {
			id := "theme-" + th.ID
			hover := m.hovering(id)
			glyph := "●"
			st := lipgloss.NewStyle().Foreground(lipgloss.Color(th.Accent))
			if th.ID == t.ID {
				glyph = "◆"
				st = st.
					Foreground(lipgloss.Color(th.Bg)).
					Background(lipgloss.Color(th.Accent))
			} else if hover {
				st = st.Background(lipgloss.Color(colorSurface))
			}
			cell := lipgloss.JoinHorizontal(lipgloss.Center, st.Render(glyph), " ")
			if m.ctx.zone != nil {
				cell = m.ctx.zone.Mark(id, cell)
			}
			cells = append(cells, cell)
		}
		rows = append(rows, lipgloss.NewStyle().Width(width).Render(
			lipgloss.JoinHorizontal(lipgloss.Center, cells...),
		))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m appModel) toggleNav() appModel {
	m.navCollapsed = !m.navCollapsed
	return m
}
