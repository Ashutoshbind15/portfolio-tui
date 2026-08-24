package main

import (
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

	toggle := m.navToggle(innerW)
	toggleH := lipgloss.Height(toggle)
	navSlot := max(0, innerH-toggleH)
	nav = lipgloss.NewStyle().
		Width(innerW).
		Height(navSlot).
		Render(nav)

	inner := lipgloss.JoinVertical(lipgloss.Left, nav, toggle)
	return st.Height(max(0, height)).Render(inner)
}

func (m appModel) navRow(page Page, width int) string {
	active := page == m.page

	var row string
	if m.navCollapsed {
		row = styleNavIcon(active).
			Width(width).
			Render(page.Icon())
	} else {
		labelW := max(0, width-navIconWidth-navGapWidth)
		icon := styleNavIcon(active).Render(page.Icon())
		gap := styleNavGap(active).Render(" ")
		label := styleNavLabel(active).
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

func (m appModel) navToggle(width int) string {
	glyph := "«"
	if m.navCollapsed {
		glyph = "»"
	}
	row := styleFaint().
		Width(width).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(glyph)
	if m.ctx.zone != nil {
		row = m.ctx.zone.Mark("nav-toggle", row)
	}
	return row
}

func (m appModel) toggleNav() appModel {
	m.navCollapsed = !m.navCollapsed
	return m
}
