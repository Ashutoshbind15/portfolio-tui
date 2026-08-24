package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type experienceItem struct {
	exp Experience
}

func (i experienceItem) FilterValue() string { return i.exp.Role + " " + i.exp.Org }
func (i experienceItem) Title() string {
	title := i.exp.Role
	if i.exp.Org != "" {
		title += " · " + i.exp.Org
	}
	if i.exp.Current {
		title += "  now"
	}
	return title
}
func (i experienceItem) Description() string { return i.exp.Period }

type experienceModel struct {
	ctx      *Context
	list     list.Model
	viewport viewport.Model
	openID   string
}

func newExperienceModel(ctx *Context) experienceModel {
	items := make([]list.Item, 0, len(experiences()))
	for _, e := range experiences() {
		items = append(items, experienceItem{exp: e})
	}
	return experienceModel{
		ctx:      ctx,
		list:     newItemList("Experience", items, 0, 0),
		viewport: viewport.New(),
	}
}

func (m experienceModel) Init() tea.Cmd { return nil }

func (m experienceModel) Activate() (experienceModel, tea.Cmd) { return m, nil }

func (m *experienceModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	if m.openID != "" {
		m.refreshDetail()
	}
}

func (m *experienceModel) refreshDetail() {
	e, ok := experienceByID(m.openID)
	if !ok {
		return
	}
	m.viewport.SetContent(renderExperienceDetail(e, m.viewport.Width()))
}

func (m experienceModel) Update(msg tea.Msg) (experienceModel, tea.Cmd) {
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
		if it, ok := m.list.SelectedItem().(experienceItem); ok {
			m.openID = it.exp.ID
			m.refreshDetail()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m experienceModel) View() string {
	if m.openID != "" {
		return m.viewport.View()
	}
	return m.list.View()
}

func renderExperienceDetail(e Experience, width int) string {
	w := max(20, width-2)
	title := e.Role
	if e.Org != "" {
		title += " at " + e.Org
	}

	var b strings.Builder
	b.WriteString(styleHeading().Render(title))
	b.WriteString("\n")
	b.WriteString(styleMuted().Render(e.Period))
	if e.Current {
		b.WriteString(styleAccent().Render("  · current"))
	}
	b.WriteString("\n\n")
	b.WriteString(styleBody().Width(w).Render(e.Summary))
	b.WriteString("\n\n")
	for _, h := range e.Highlights {
		b.WriteString(styleMuted().Width(w).Render("•  " + h))
		b.WriteString("\n")
	}
	if len(e.Tech) > 0 {
		b.WriteString("\n")
		b.WriteString(joinChips(e.Tech))
	}
	b.WriteString("\n\n")
	b.WriteString(styleFaint().Render("esc back to the list"))
	return lipgloss.NewStyle().Padding(1, 1, 0, 1).Render(b.String())
}
