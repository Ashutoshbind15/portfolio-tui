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
	hover  string
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
	showSplash   bool
	splash       splashModel
	sel          textSel
	frame        *string

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
		showSplash: true,
		splash:     newSplashModel(),
		frame:      new(string),
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
	return m.splash.Init()
}

func (m appModel) enterPortfolio() (appModel, tea.Cmd) {
	m.showSplash = false
	next, cmd := m.activateCurrentPage()
	next.setSizes()
	return next, cmd
}

func (m appModel) openSplash() (appModel, tea.Cmd) {
	m.showSplash = true
	m.splash.Reset(m.ctx.width, m.ctx.height)
	return m, m.splash.Init()
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
	m.splash.SetSize(m.ctx.width, m.ctx.height)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.width = msg.Width
		m.ctx.height = msg.Height
		m.setSizes()
		return m, nil

	case splashFrameMsg:
		if !m.showSplash {
			return m, nil
		}
		var cmd tea.Cmd
		m.splash, cmd = m.splash.Update(msg)
		return m, cmd

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
		}
		if m.showSplash {
			switch msg.String() {
			case "enter", " ", "space":
				return m.enterPortfolio()
			case "]":
				next := m.cycleTheme(1)
				next.setSizes()
				return next, nil
			}
			return m, nil
		}
		switch msg.String() {
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
		if msg.Button != tea.MouseLeft {
			break
		}
		if m.showSplash {
			m.sel.press(msg.X, msg.Y)
			return m, nil
		}
		m.trackHover(msg)
		if m.hitNavZone(msg) {
			m.sel.clear()
			return m.handleNavClick(msg)
		}
		m.sel.press(msg.X, msg.Y)

	case tea.MouseMotionMsg:
		if m.sel.down {
			m.sel.drag(msg.X, msg.Y)
			if m.sel.dragging {
				m.ctx.hover = ""
				return m, nil
			}
		}
		if !m.showSplash {
			m.trackHover(msg)
		}

	case tea.MouseWheelMsg:
		if !m.showSplash && m.hitNavZone(msg) {
			return m, nil
		}

	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			break
		}
		wasDown := m.sel.down
		moved := m.sel.dragging
		m.sel.down = false
		m.sel.dragging = false
		if m.showSplash {
			if wasDown && !moved {
				m.sel.clear()
				return m.enterPortfolio()
			}
			if moved {
				return m.copySelection()
			}
			break
		}
		if moved {
			return m.copySelection()
		}
		m.sel.active = false
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
	brandSt := styleBrand()
	if m.hovering("brand") {
		brandSt = brandSt.Foreground(lipgloss.Color(colorAccent))
	}
	brand := brandSt.Render(fmt.Sprintf(">_ %s", p.Name))
	if m.ctx.zone != nil {
		brand = m.ctx.zone.Mark("brand", brand)
	}
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
	if m.showSplash {
		return m.splashView()
	}

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

	return m.finishView(output)
}

func (m appModel) splashView() tea.View {
	output := m.splash.View()
	if m.ctx.width > 0 && m.ctx.height > 0 {
		output = styleApp().
			Width(m.ctx.width).
			Height(m.ctx.height).
			MaxHeight(m.ctx.height).
			Render(output)
	}
	return m.finishView(output)
}

func (m *appModel) trackHover(msg tea.MouseMsg) {
	m.ctx.hover = m.hitZoneID(msg)
}

func (m appModel) navZoneIDs() []string {
	ids := []string{"brand", "theme-cycle", "nav-toggle"}
	for _, th := range allThemes() {
		ids = append(ids, "theme-"+th.ID)
	}
	for _, page := range navPages() {
		ids = append(ids, "nav-"+string(page))
	}
	return ids
}

func (m appModel) hitZoneID(msg tea.MouseMsg) string {
	if m.ctx.zone == nil {
		return ""
	}
	for _, id := range m.navZoneIDs() {
		if m.ctx.zone.Get(id).InBounds(msg) {
			return id
		}
	}
	return ""
}

func (m appModel) hitNavZone(msg tea.MouseMsg) bool {
	return m.hitZoneID(msg) != ""
}

func (m appModel) hovering(id string) bool {
	return m.ctx.hover == id
}

func (m appModel) handleNavClick(msg tea.MouseClickMsg) (appModel, tea.Cmd) {
	if m.ctx.zone.Get("brand").InBounds(msg) {
		return m.openSplash()
	}
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
		if m.ctx.zone.Get("nav-" + string(page)).InBounds(msg) {
			next, cmd := m.navigateTo(page)
			next.setSizes()
			return next, cmd
		}
	}
	return m, nil
}

func (m appModel) copySelection() (appModel, tea.Cmd) {
	m.sel.active = true
	if m.frame == nil {
		return m, nil
	}
	_, text := applySelection(*m.frame, m.sel.originX, m.sel.originY, m.sel.x, m.sel.y)
	if text != "" {
		return m, tea.SetClipboard(text)
	}
	return m, nil
}

func (m appModel) finishView(output string) tea.View {
	if m.frame != nil {
		*m.frame = output
	}
	if m.sel.live() {
		output, _ = applySelection(output, m.sel.originX, m.sel.originY, m.sel.x, m.sel.y)
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
