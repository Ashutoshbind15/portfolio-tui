package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type homeModel struct {
	ctx *Context
}

func newHomeModel(ctx *Context) homeModel {
	return homeModel{ctx: ctx}
}

func (m homeModel) Init() tea.Cmd                  { return nil }
func (m homeModel) Activate() (homeModel, tea.Cmd) { return m, nil }
func (m homeModel) Update(msg tea.Msg) (homeModel, tea.Cmd) {
	return m, nil
}

func (m homeModel) View() string {
	p := profile()
	width := max(20, m.ctx.innerW)

	name := styleHeading().Render(p.Name)
	role := styleMuted().Render(p.Role)

	var bullets []string
	for _, line := range p.Bullets {
		bullets = append(bullets, styleBody().Width(width).Render("•  "+line))
	}

	contact := lipgloss.JoinVertical(lipgloss.Left,
		styleFaint().Render("@  ")+styleBody().Render(p.Email),
		styleFaint().Render("$  ")+styleBody().Render(p.SSH),
		styleFaint().Render("   ")+styleMuted().Render(p.GitHub),
		styleFaint().Render("   ")+styleMuted().Render(p.Website),
	)

	hint := styleFaint().Render("tab cycles pages · [ toggles nav · enter opens a list item")

	return lipgloss.JoinVertical(lipgloss.Left,
		name,
		role,
		"",
		lipgloss.JoinVertical(lipgloss.Left, bullets...),
		"",
		contact,
		"",
		hint,
	)
}
