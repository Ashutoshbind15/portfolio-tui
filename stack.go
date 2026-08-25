package main

import (
	"strings"

	"charm.land/bubbles/v2/tree"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func treeStyles() tree.Styles {
	s := tree.DefaultDarkStyles()
	s.NodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
	s.SelectedNodeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBright)).
		Bold(true)
	s.RootNodeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		Bold(true)
	s.ParentNodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
	s.CursorStyle = s.CursorStyle.Foreground(lipgloss.Color(colorAccent))
	s.EnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorFaint))
	s.SelectedEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
	s.IndenterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorFaint))
	s.OpenIndicatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
	return s
}

func newStackModel(ctx *Context) stackModel {
	t := tree.New(newStackTree(), 0, 0)
	t.SetShowHelp(false)
	t.SetStyles(treeStyles())
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
	switch msg := msg.(type) {
	case tea.MouseMotionMsg:
		if line := m.lineAt(msg); line >= 0 {
			if m.tree.YOffset() != line {
				m.tree.SetYOffset(line)
			}
			return m, nil
		}
	case tea.MouseClickMsg:
		if line := m.lineAt(msg); line >= 0 {
			if m.tree.YOffset() != line {
				m.tree.SetYOffset(line)
			}
			m.tree.ToggleCurrentNode()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m stackModel) lineAt(msg tea.MouseMsg) int {
	if m.ctx.zone == nil {
		return -1
	}
	for i := 0; i < m.tree.Height(); i++ {
		if m.ctx.zone.Get(stackZoneID(i)).InBounds(msg) {
			return m.tree.ViewportYOffset() + i
		}
	}
	return -1
}

func (m stackModel) View() string {
	view := m.tree.View()
	if m.ctx.zone == nil {
		return view
	}
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		lines[i] = m.ctx.zone.Mark(stackZoneID(i), line)
	}
	return strings.Join(lines, "\n")
}

func (m *stackModel) applyTheme() {
	m.tree.SetStyles(treeStyles())
}
