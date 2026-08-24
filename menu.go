package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type pageMenuItem struct {
	page        Page
	title       string
	description string
}

func (i pageMenuItem) FilterValue() string { return i.title }
func (i pageMenuItem) Title() string       { return i.title }
func (i pageMenuItem) Description() string { return i.description }

func newPageList(width, height int) list.Model {
	items := make([]list.Item, 0, len(navPages()))
	for _, page := range navPages() {
		items = append(items, pageMenuItem{
			page:        page,
			title:       page.Title(),
			description: page.Description(),
		})
	}
	return newItemList("Pages", items, width, height)
}

type menuModel struct {
	ctx      *Context
	pageList list.Model
}

func newMenuModel(ctx *Context) menuModel {
	return menuModel{
		ctx:      ctx,
		pageList: newPageList(0, 0),
	}
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Activate() (menuModel, tea.Cmd) { return m, nil }

func (m *menuModel) SetSize(width, height int) {
	m.pageList.SetSize(width, height)
}

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "enter":
			if it, ok := m.pageList.SelectedItem().(pageMenuItem); ok {
				return m, navigateTo(it.page)
			}
		case "esc":
			return m, func() tea.Msg { return closeMenuMsg{} }
		}
	}
	var cmd tea.Cmd
	m.pageList, cmd = m.pageList.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	return m.pageList.View()
}
