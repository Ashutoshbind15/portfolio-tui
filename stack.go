package main

import (
	"charm.land/bubbles/v2/tree"
	tea "charm.land/bubbletea/v2"
)

type stackModel struct {
	ctx  *Context
	tree tree.Model
}

func newStackTree() *tree.Node {
	root := tree.Root("Stack")
	for _, cat := range stackCategories() {
		node := tree.Root(cat.Name)
		for _, item := range cat.Items {
			node.Child(item)
		}
		root.Child(node)
	}
	return root
}

func newStackModel(ctx *Context) stackModel {
	t := tree.New(newStackTree(), 0, 0)
	t.SetShowHelp(false)
	return stackModel{
		ctx:  ctx,
		tree: t,
	}
}

func (m stackModel) Init() tea.Cmd { return nil }

func (m stackModel) Activate() (stackModel, tea.Cmd) { return m, nil }

func (m *stackModel) SetSize(width, height int) {
	m.tree.SetSize(width, height)
}

func (m stackModel) Update(msg tea.Msg) (stackModel, tea.Cmd) {
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m stackModel) View() string {
	return m.tree.View()
}
