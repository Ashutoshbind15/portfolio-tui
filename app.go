package main

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Context struct {
	zone   *zone.Manager
	width  int
	height int
}

type keyMap struct {
	Menu key.Binding
	Open key.Binding
	Back key.Binding
	Quit key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Menu: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "pages")),
		Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Menu, k.Open, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

type appModel struct {
	ctx          *Context
	page         Page
	previousPage *Page
	keys         keyMap
	help         help.Model

	home       homeModel
	menu       menuModel
	projects   projectsModel
	experience experienceModel
	stack      stackModel
	blogs      blogsModel
}

func newAppModel() appModel {
	ctx := &Context{zone: zone.New()}
	return appModel{
		ctx:        ctx,
		page:       PageHome,
		keys:       newKeyMap(),
		help:       help.New(),
		home:       newHomeModel(ctx),
		menu:       newMenuModel(ctx),
		projects:   newProjectsModel(ctx),
		experience: newExperienceModel(ctx),
		stack:      newStackModel(ctx),
		blogs:      newBlogsModel(ctx),
	}
}

func (m appModel) Init() tea.Cmd {
	return m.home.Init()
}

func (m appModel) openPageSelect() appModel {
	if m.page == PageSelect {
		return m
	}
	prev := m.page
	m.previousPage = &prev
	m.page = PageSelect
	return m
}

func (m appModel) closePageSelect() appModel {
	if m.previousPage == nil {
		return m
	}
	m.page = *m.previousPage
	m.previousPage = nil
	return m
}

func (m appModel) effectivePage() Page {
	if m.page == PageSelect && m.previousPage != nil {
		return *m.previousPage
	}
	return m.page
}

func (m appModel) navigateTo(page Page) (appModel, tea.Cmd) {
	m.page = page
	m.previousPage = nil
	return m.activateCurrentPage()
}

func (m appModel) activateCurrentPage() (appModel, tea.Cmd) {
	switch m.page {
	case PageHome:
		var cmd tea.Cmd
		m.home, cmd = m.home.Activate()
		return m, cmd
	case PageProjects:
		var cmd tea.Cmd
		m.projects, cmd = m.projects.Activate()
		return m, cmd
	case PageExperience:
		var cmd tea.Cmd
		m.experience, cmd = m.experience.Activate()
		return m, cmd
	case PageStack:
		var cmd tea.Cmd
		m.stack, cmd = m.stack.Activate()
		return m, cmd
	case PageBlogs:
		var cmd tea.Cmd
		m.blogs, cmd = m.blogs.Activate()
		return m, cmd
	case PageSelect:
		var cmd tea.Cmd
		m.menu, cmd = m.menu.Activate()
		return m, cmd
	}
	return m, nil
}

func (m *appModel) setSizes() {
	header := m.headerView()
	footer := m.footerView()
	h := max(0, m.ctx.height-lipgloss.Height(header)-lipgloss.Height(footer))
	w := max(0, m.ctx.width-2)
	m.menu.SetSize(w, h)
	m.projects.SetSize(w, h)
	m.experience.SetSize(w, h)
	m.stack.SetSize(w, h)
	m.blogs.SetSize(w, h)
	m.help.SetWidth(w)
}

func pageForDigit(s string) (Page, bool) {
	switch s {
	case "1":
		return PageHome, true
	case "2":
		return PageProjects, true
	case "3":
		return PageExperience, true
	case "4":
		return PageStack, true
	case "5":
		return PageBlogs, true
	}
	return "", false
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.width = msg.Width
		m.ctx.height = msg.Height
		m.setSizes()
		return m, nil

	case navigateMsg:
		next, cmd := m.navigateTo(msg.page)
		next.setSizes()
		return next, cmd

	case closeMenuMsg:
		m = m.closePageSelect()
		m.setSizes()
		return m.activateCurrentPage()

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.ctx.zone != nil {
				m.ctx.zone.Close()
			}
			return m, tea.Quit
		case "tab":
			m = m.openPageSelect()
			m.setSizes()
			return m.activateCurrentPage()
		default:
			if page, ok := pageForDigit(msg.String()); ok && m.page != PageSelect {
				next, cmd := m.navigateTo(page)
				next.setSizes()
				return next, cmd
			}
		}

	case tea.MouseClickMsg:
		if m.ctx.zone != nil {
			for _, page := range navPages() {
				id := "nav-" + string(page)
				if m.ctx.zone.Get(id).InBounds(msg) {
					next, cmd := m.navigateTo(page)
					next.setSizes()
					return next, cmd
				}
			}
		}
	}

	switch m.page {
	case PageHome:
		var cmd tea.Cmd
		m.home, cmd = m.home.Update(msg)
		return m, cmd
	case PageSelect:
		var cmd tea.Cmd
		m.menu, cmd = m.menu.Update(msg)
		return m, cmd
	case PageProjects:
		var cmd tea.Cmd
		m.projects, cmd = m.projects.Update(msg)
		return m, cmd
	case PageExperience:
		var cmd tea.Cmd
		m.experience, cmd = m.experience.Update(msg)
		return m, cmd
	case PageStack:
		var cmd tea.Cmd
		m.stack, cmd = m.stack.Update(msg)
		return m, cmd
	case PageBlogs:
		var cmd tea.Cmd
		m.blogs, cmd = m.blogs.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m appModel) headerView() string {
	p := profile()
	brand := styleBrand().Render(fmt.Sprintf(">_ %s", p.Name))

	current := m.effectivePage()
	tabs := make([]string, 0, len(navPages()))
	for i, page := range navPages() {
		label := fmt.Sprintf("%d %s", i+1, page.Title())
		var styled string
		if page == current {
			styled = styleNavActive().Render(label)
		} else {
			styled = styleNavIdle().Render(label)
		}
		if m.ctx.zone != nil {
			styled = m.ctx.zone.Mark("nav-"+string(page), styled)
		}
		tabs = append(tabs, styled)
	}

	nav := lipgloss.JoinHorizontal(lipgloss.Center, joinWithGap(tabs, "  ")...)
	row := lipgloss.JoinHorizontal(lipgloss.Center, brand, "   ", nav)
	return styleChromeBorder().Width(max(0, m.ctx.width)).Render(row)
}

func joinWithGap(parts []string, gap string) []string {
	if len(parts) == 0 {
		return parts
	}
	out := make([]string, 0, len(parts)*2-1)
	for i, part := range parts {
		if i > 0 {
			out = append(out, gap)
		}
		out = append(out, part)
	}
	return out
}

func (m appModel) footerView() string {
	return styleFooter().Width(max(0, m.ctx.width)).Render(m.help.View(m.keys))
}

func (m appModel) pageContent() string {
	switch m.page {
	case PageHome:
		return m.home.View()
	case PageSelect:
		return m.menu.View()
	case PageProjects:
		return m.projects.View()
	case PageExperience:
		return m.experience.View()
	case PageStack:
		return m.stack.View()
	case PageBlogs:
		return m.blogs.View()
	default:
		return "Unknown page"
	}
}

func (m appModel) View() tea.View {
	header := m.headerView()
	footer := m.footerView()
	bodyH := max(0, m.ctx.height-lipgloss.Height(header)-lipgloss.Height(footer))

	content := lipgloss.NewStyle().
		Width(max(0, m.ctx.width)).
		Height(bodyH).
		Padding(0, 1).
		Render(m.pageContent())

	output := lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
	if m.ctx.zone != nil {
		output = m.ctx.zone.Scan(output)
	}

	v := tea.NewView(output)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.WindowTitle = profile().Name
	return v
}
