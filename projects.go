package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type projectItem struct {
	project Project
}

func (i projectItem) FilterValue() string { return i.project.Name }
func (i projectItem) Title() string       { return i.project.Name }
func (i projectItem) Description() string { return i.project.Summary }

type projectsModel struct {
	ctx      *Context
	list     list.Model
	viewport viewport.Model
	openID   string
}

func newProjectsModel(ctx *Context) projectsModel {
	items := make([]list.Item, 0, len(projects()))
	for _, p := range projects() {
		items = append(items, projectItem{project: p})
	}
	return projectsModel{
		ctx:      ctx,
		list:     newItemList("Projects", items, 0, 0),
		viewport: viewport.New(),
	}
}

func (m projectsModel) Init() tea.Cmd { return nil }

func (m projectsModel) Activate() (projectsModel, tea.Cmd) {
	return m, nil
}

func (m *projectsModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
	sizeViewport(&m.viewport, width, height)
	if m.openID != "" {
		m.refreshDetail()
	}
}

func (m *projectsModel) refreshDetail() {
	p, ok := projectByID(m.openID)
	if !ok {
		return
	}
	m.viewport.SetContent(renderProjectDetail(p, m.viewport.Width()))
}

func (m projectsModel) Update(msg tea.Msg) (projectsModel, tea.Cmd) {
	if m.openID != "" {
		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
			m.openID = ""
			return m, nil
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
		if it, ok := m.list.SelectedItem().(projectItem); ok {
			m.openID = it.project.ID
			m.refreshDetail()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m projectsModel) View() string {
	if m.openID != "" {
		return viewWithScrollbar(m.viewport)
	}
	return m.list.View()
}

func (m *projectsModel) applyTheme() {
	restyleList(&m.list)
	if m.openID != "" {
		m.refreshDetail()
	}
}

func renderProjectDetail(p Project, width int) string {
	w := max(20, width-2)
	var b strings.Builder
	b.WriteString(styleHeading().Render(p.Name))
	if p.Kind != "" {
		b.WriteString(styleMuted().Render("  " + p.Kind))
	}
	b.WriteString("\n")
	if p.Tagline != "" {
		b.WriteString("\n")
		b.WriteString(styleMuted().Width(w).Render(p.Tagline))
	}

	writeDetailSection(&b, "what", p.What, w)
	writeDetailSection(&b, "why", p.Why, w)

	if projectHasLinks(p) {
		b.WriteString("\n\n")
		b.WriteString(styleFaint().Render("where"))
		b.WriteString("\n")
		for _, link := range projectLinks(p) {
			b.WriteString(styleFaint().Render(padLinkLabel(link.Label)))
			b.WriteString(styleMuted().Render(link.Value))
			b.WriteString("\n")
		}
	}

	writeDetailSection(&b, "technicals", p.Technicals, w)

	if len(p.Upcoming) > 0 {
		b.WriteString("\n\n")
		b.WriteString(styleFaint().Render("upcoming"))
		b.WriteString("\n")
		for _, item := range p.Upcoming {
			b.WriteString(styleMuted().Width(w).Render("•  " + item))
			b.WriteString("\n")
		}
	}

	if len(p.Tech) > 0 {
		b.WriteString("\n")
		b.WriteString(joinChips(p.Tech))
	}
	b.WriteString("\n\n")
	b.WriteString(styleFaint().Render("esc back to the list"))
	return lipgloss.NewStyle().Padding(1, 1, 0, 1).Render(b.String())
}

func writeDetailSection(b *strings.Builder, label, body string, width int) {
	if body == "" {
		return
	}
	b.WriteString("\n\n")
	b.WriteString(styleFaint().Render(label))
	b.WriteString("\n")
	b.WriteString(styleBody().Width(width).Render(body))
}
