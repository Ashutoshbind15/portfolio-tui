package main

import (
	"fmt"
	"strings"

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
	innerW int
}

type keyMap struct {
	Next      key.Binding
	Prev      key.Binding
	Open      key.Binding
	Back      key.Binding
	Move      key.Binding
	ToggleNav key.Binding
	Theme     key.Binding
	Quit      key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Next:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
		Prev:      key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
		Open:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Move:      key.NewBinding(key.WithKeys("up", "down", "j", "k"), key.WithHelp("j/k", "move")),
		ToggleNav: key.NewBinding(key.WithKeys("["), key.WithHelp("[", "nav")),
		Theme:     key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "theme")),
		Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Move, k.Open, k.Back, k.ToggleNav, k.Theme, k.Quit}
}

func (k keyMap) forPage(page Page) keyMap {
	switch page {
	case PageHome:
		k.Open = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "expand"))
		k.Move = key.NewBinding(key.WithKeys("up", "down", "j", "k"), key.WithHelp("j/k", "move"))
	case PageStack:
		k.Open = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "toggle"))
	}
	return k
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

type appModel struct {
	ctx          *Context
	page         Page
	navCollapsed bool
	themeID      string
	keys         keyMap
	help         help.Model

	home       homeModel
	projects   projectsModel
	experience experienceModel
	stack      stackModel
	blogs      blogsModel
}

func newAppModel() appModel {
	ctx := &Context{zone: zone.New()}
	m := appModel{
		ctx:        ctx,
		page:       PageHome,
		themeID:    defaultThemeID,
		keys:       newKeyMap(),
		help:       newHelp(),
		home:       newHomeModel(ctx),
		projects:   newProjectsModel(ctx),
		experience: newExperienceModel(ctx),
		stack:      newStackModel(ctx),
		blogs:      newBlogsModel(ctx),
	}
	m.applyThemeID(defaultThemeID)
	return m
}

func (m appModel) Init() tea.Cmd {
	return m.home.Init()
}

func (m appModel) navigateTo(page Page) (appModel, tea.Cmd) {
	m.page = page
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
	}
	return m, nil
}

const contentPadX = 2

func (m appModel) columnWidth() int {
	return max(0, m.ctx.width-sidebarFrameWidth(m.navCollapsed))
}

func (m *appModel) contentSize() (width, height int) {
	width = max(0, m.columnWidth()-contentPadX)
	height = max(0, m.ctx.height-lipgloss.Height(m.headerView())-lipgloss.Height(m.footerView()))
	return width, height
}

func (m *appModel) setSizes() {
	w, h := m.contentSize()
	m.ctx.innerW = w
	m.home.SetSize(w, h)
	m.projects.SetSize(w, h)
	m.experience.SetSize(w, h)
	m.stack.SetSize(w, h)
	m.blogs.SetSize(w, h)
	m.help.SetWidth(max(0, m.ctx.width-4))
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

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.ctx.zone != nil {
				m.ctx.zone.Close()
			}
			return m, tea.Quit
		case "tab":
			next, cmd := m.navigateTo(cyclePage(m.page, 1))
			next.setSizes()
			return next, cmd
		case "shift+tab":
			next, cmd := m.navigateTo(cyclePage(m.page, -1))
			next.setSizes()
			return next, cmd
		case "[":
			next := m.toggleNav()
			next.setSizes()
			return next, nil
		case "]":
			next := m.cycleTheme(1)
			next.setSizes()
			return next, nil
		}

	case tea.MouseClickMsg:
		if m.ctx.zone != nil {
			if m.ctx.zone.Get("theme-cycle").InBounds(msg) {
				next := m.cycleTheme(1)
				next.setSizes()
				return next, nil
			}
			for _, th := range allThemes() {
				if m.ctx.zone.Get("theme-" + th.ID).InBounds(msg) {
					next := m.setTheme(th.ID)
					next.setSizes()
					return next, nil
				}
			}
			if m.ctx.zone.Get("nav-toggle").InBounds(msg) {
				next := m.toggleNav()
				next.setSizes()
				return next, nil
			}
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

func (m appModel) splitRule(junction string) string {
	col := m.columnWidth()
	side := sidebarFrameWidth(m.navCollapsed)
	line := strings.Repeat("─", max(0, col)) + junction + strings.Repeat("─", max(0, side-1))
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBorder)).
		Background(lipgloss.Color(colorBg)).
		Width(max(0, m.ctx.width)).
		Render(line)
}

func (m appModel) headerView() string {
	p := profile()
	brand := styleBrand().Render(fmt.Sprintf(">_ %s", p.Name))
	inner := styleHeader().Width(m.ctx.width).Render(brand)
	return lipgloss.JoinVertical(lipgloss.Left, inner, m.splitRule("┬"))
}

func (m appModel) footerView() string {
	help := styleFooter().Width(m.ctx.width).Render(m.help.View(m.keys.forPage(m.page)))
	return lipgloss.JoinVertical(lipgloss.Left, m.splitRule("┴"), help)
}

func (m appModel) pageContent() string {
	switch m.page {
	case PageHome:
		return m.home.View()
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
	_, bodyH := m.contentSize()

	content := lipgloss.NewStyle().
		Width(m.columnWidth()).
		Height(bodyH).
		Background(lipgloss.Color(colorBg)).
		Padding(0, 0, 0, contentPadX).
		Render(m.pageContent())

	body := lipgloss.JoinHorizontal(lipgloss.Top, content, m.sidebarView(bodyH))
	output := lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
	if m.ctx.zone != nil {
		output = m.ctx.zone.Scan(output)
	}
	if m.ctx.width > 0 && m.ctx.height > 0 {
		output = styleApp().
			Width(m.ctx.width).
			Height(m.ctx.height).
			MaxHeight(m.ctx.height).
			Render(output)
	}

	v := tea.NewView(output)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.WindowTitle = profile().Name
	if c := parseHexColor(colorBg); c != nil {
		v.BackgroundColor = c
	}
	if c := parseHexColor(colorText); c != nil {
		v.ForegroundColor = c
	}
	return v
}

func (m appModel) setTheme(id string) appModel {
	m.applyThemeID(id)
	return m
}

func (m appModel) cycleTheme(delta int) appModel {
	return m.setTheme(cycleThemeID(m.themeID, delta))
}

func (m *appModel) applyThemeID(id string) {
	t, ok := themeByID(id)
	if !ok {
		t, _ = themeByID(defaultThemeID)
	}
	m.themeID = t.ID
	applyTheme(t)
	m.help = newHelp()
	m.help.SetWidth(max(0, m.ctx.width-4))
	m.home.refresh()
	m.projects.applyTheme()
	m.experience.applyTheme()
	m.stack.applyTheme()
	m.blogs.applyTheme()
}
