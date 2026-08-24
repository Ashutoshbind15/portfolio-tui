package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type blogItem struct {
	blog Blog
}

func (i blogItem) FilterValue() string { return i.blog.Title }
func (i blogItem) Title() string       { return i.blog.Title }
func (i blogItem) Description() string { return i.blog.Date + "  ·  " + i.blog.Summary }

type blogsModel struct {
	ctx      *Context
	list     list.Model
	viewport viewport.Model
	openID   string
}

func newBlogsModel(ctx *Context) blogsModel {
	items := make([]list.Item, 0, len(blogs()))
	for _, b := range blogs() {
		items = append(items, blogItem{blog: b})
	}
	return blogsModel{
		ctx:      ctx,
		list:     newItemList("Blogs", items, 0, 0),
		viewport: viewport.New(),
	}
}

func (m blogsModel) Init() tea.Cmd { return nil }

func (m blogsModel) Activate() (blogsModel, tea.Cmd) { return m, nil }

func (m *blogsModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	if m.openID != "" {
		m.refreshDetail()
	}
}

func (m *blogsModel) refreshDetail() {
	b, ok := blogByID(m.openID)
	if !ok {
		return
	}
	m.viewport.SetContent(renderBlogDetail(b, m.viewport.Width()))
}

func (m blogsModel) Update(msg tea.Msg) (blogsModel, tea.Cmd) {
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
		if it, ok := m.list.SelectedItem().(blogItem); ok {
			m.openID = it.blog.ID
			m.refreshDetail()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m blogsModel) View() string {
	if m.openID != "" {
		return m.viewport.View()
	}
	return m.list.View()
}

func renderBlogDetail(b Blog, width int) string {
	w := max(20, width-2)
	var sb strings.Builder
	sb.WriteString(styleHeading().Render(b.Title))
	sb.WriteString("\n")
	sb.WriteString(styleMuted().Render(b.Date))
	sb.WriteString("\n\n")
	sb.WriteString(styleBody().Width(w).Render(b.Body))
	sb.WriteString("\n\n")
	sb.WriteString(styleFaint().Render("esc back to the list"))
	return lipgloss.NewStyle().Padding(1, 1, 0, 1).Render(sb.String())
}
