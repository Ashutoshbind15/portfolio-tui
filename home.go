package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type homeItemKind int

const (
	homeItemProject homeItemKind = iota
	homeItemExperience
	homeItemStack
)

type homeItem struct {
	kind homeItemKind
	id   string
}

type homeTickMsg time.Time

func tickIST() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return homeTickMsg(t)
	})
}

func formatIST(t time.Time) string {
	ist := time.FixedZone("IST", 5*3600+30*60)
	return t.In(ist).Format("3:04 pm MST")
}

type homeModel struct {
	ctx       *Context
	viewport  viewport.Model
	cursor    int
	expanded  map[string]bool
	items     []homeItem
	itemLines [][2]int
	now       time.Time
}

func newHomeModel(ctx *Context) homeModel {
	items := make([]homeItem, 0, len(projects())+len(experiences())+len(stackCategories()))
	for _, p := range projects() {
		items = append(items, homeItem{kind: homeItemProject, id: p.ID})
	}
	for _, e := range experiences() {
		items = append(items, homeItem{kind: homeItemExperience, id: e.ID})
	}
	for _, cat := range stackCategories() {
		items = append(items, homeItem{kind: homeItemStack, id: cat.Name})
	}

	vp := viewport.New()
	vp.FillHeight = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.KeyMap.Up = key.NewBinding()
	vp.KeyMap.Down = key.NewBinding()
	vp.KeyMap.Left = key.NewBinding()
	vp.KeyMap.Right = key.NewBinding()
	vp.KeyMap.PageDown = key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	)
	vp.KeyMap.PageUp = key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	)
	vp.KeyMap.HalfPageDown = key.NewBinding()
	vp.KeyMap.HalfPageUp = key.NewBinding()

	return homeModel{
		ctx:      ctx,
		viewport: vp,
		expanded: map[string]bool{},
		items:    items,
		now:      time.Now(),
	}
}

func (m homeModel) Init() tea.Cmd { return tickIST() }

func (m homeModel) Activate() (homeModel, tea.Cmd) {
	m.now = time.Now()
	m.refresh()
	return m, tickIST()
}

func (m *homeModel) SetSize(width, height int) {
	sizeViewport(&m.viewport, width, height)
	m.refresh()
}

func (it homeItem) key() string {
	return fmt.Sprintf("%d:%s", it.kind, it.id)
}

func homeZoneID(i int) string {
	return "home-item-" + strconv.Itoa(i)
}

func (m homeModel) itemAt(msg tea.MouseMsg) int {
	if m.ctx.zone == nil {
		return -1
	}
	for i := range m.items {
		if m.ctx.zone.Get(homeZoneID(i)).InBounds(msg) {
			return i
		}
	}
	return -1
}

func (m *homeModel) toggle(i int) {
	if i < 0 || i >= len(m.items) {
		return
	}
	k := m.items[i].key()
	m.expanded[k] = !m.expanded[k]
}

func (m *homeModel) isOpen(i int) bool {
	if i < 0 || i >= len(m.items) {
		return false
	}
	return m.expanded[m.items[i].key()]
}

func (m *homeModel) moveCursor(delta int) {
	n := len(m.items)
	if n == 0 {
		return
	}
	next := m.cursor + delta
	if next < 0 {
		next = 0
	}
	if next >= n {
		next = n - 1
	}
	if next == m.cursor {
		return
	}
	m.cursor = next
	m.refresh()
	m.ensureCursorVisible()
}

func (m *homeModel) ensureCursorVisible() {
	if m.cursor < 0 || m.cursor >= len(m.itemLines) {
		return
	}
	start, end := m.itemLines[m.cursor][0], m.itemLines[m.cursor][1]
	h := m.viewport.Height()
	if h <= 0 {
		return
	}
	top := m.viewport.YOffset()
	if start < top {
		m.viewport.SetYOffset(start)
		return
	}
	if end >= top+h {
		m.viewport.SetYOffset(min(start, max(0, end-h+1)))
	}
}

func (m *homeModel) refresh() {
	if m.viewport.Width() <= 0 {
		return
	}
	offset := m.viewport.YOffset()
	content, lines := m.renderPage()
	m.itemLines = lines
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(offset)
}

func (m homeModel) Update(msg tea.Msg) (homeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case homeTickMsg:
		m.now = time.Time(msg)
		m.refresh()
		return m, tickIST()

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
			return m, nil
		case "down", "j":
			m.moveCursor(1)
			return m, nil
		case "enter", " ", "space":
			m.toggle(m.cursor)
			m.refresh()
			m.ensureCursorVisible()
			return m, nil
		}

	case tea.MouseClickMsg:
		if i := m.itemAt(msg); i >= 0 {
			m.cursor = i
			m.toggle(i)
			m.refresh()
			return m, nil
		}

	case tea.MouseMotionMsg:
		if i := m.itemAt(msg); i >= 0 && i != m.cursor {
			m.cursor = i
			m.refresh()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m homeModel) View() string {
	return viewWithScrollbar(m.viewport)
}

type pageBuilder struct {
	lines []string
}

func (p *pageBuilder) add(s string) {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return
	}
	p.lines = append(p.lines, strings.Split(s, "\n")...)
}

func (p *pageBuilder) blank() {
	p.lines = append(p.lines, "")
}

func (p *pageBuilder) height() int { return len(p.lines) }

func (p *pageBuilder) string() string {
	return strings.Join(p.lines, "\n")
}

func wrapText(style lipgloss.Style, text string, width int) string {
	rendered := style.Width(max(1, width)).Render(text)
	lines := strings.Split(strings.TrimRight(rendered, "\n"), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

func wrapIndent(style lipgloss.Style, indent int, text string, width int) string {
	inner := max(1, width-indent)
	wrapped := wrapText(style, text, inner)
	if indent <= 0 {
		return wrapped
	}
	pad := strings.Repeat(" ", indent)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = pad + strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

func wrapBullet(style lipgloss.Style, text string, width int) string {
	const prefix = "•  "
	pw := lipgloss.Width(prefix)
	wrapped := wrapText(style, text, max(1, width-pw))
	pad := strings.Repeat(" ", pw)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		line = strings.TrimRight(line, " ")
		if i == 0 {
			lines[i] = prefix + line
		} else {
			lines[i] = pad + line
		}
	}
	return strings.Join(lines, "\n")
}

func (m *homeModel) renderPage() (string, [][2]int) {
	w := max(20, m.viewport.Width())
	p := profile()
	var b pageBuilder
	itemLines := make([][2]int, len(m.items))
	idx := 0

	b.add(styleHeading().Render(p.Name))
	if uw := min(w, lipgloss.Width(p.Name)); uw > 0 {
		b.add(styleAccent().Render(strings.Repeat("─", uw)))
	}
	b.add(styleMuted().Render(p.Role))
	b.blank()
	for _, line := range p.Bullets {
		b.add(wrapBullet(styleBody(), line, w))
	}
	b.blank()
	b.add(styleMuted().Render(p.Location + "  ·  " + formatIST(m.now)))
	b.blank()
	b.add(sectionRule("software", w))
	b.blank()

	for _, project := range projects() {
		start := b.height()
		b.add(m.renderProjectRow(idx, project, w))
		itemLines[idx] = [2]int{start, max(start, b.height()-1)}
		idx++
		b.blank()
	}

	b.add(sectionRule("experience", w))
	b.blank()
	for _, exp := range experiences() {
		start := b.height()
		b.add(m.renderExperienceRow(idx, exp, w))
		itemLines[idx] = [2]int{start, max(start, b.height()-1)}
		idx++
		b.blank()
	}

	b.add(sectionRule("education", w))
	b.blank()
	b.add(m.renderEducation(w))
	b.blank()

	b.add(sectionRule("stack", w))
	b.blank()
	b.add(styleFaint().Render(fmt.Sprintf("%d tools  ·  enter a row", stackToolCount())))
	b.blank()
	for i, cat := range stackCategories() {
		start := b.height()
		b.add(m.renderStackRow(idx, i+1, cat, w))
		itemLines[idx] = [2]int{start, max(start, b.height()-1)}
		idx++
		b.blank()
	}

	b.add(sectionRule("contact", w))
	b.blank()
	if len(p.Contact) > 0 {
		for _, line := range p.Contact {
			b.add(wrapBullet(styleBody(), line, w))
		}
		b.blank()
	}
	for _, row := range contactRows(p) {
		b.add(renderContactRow(row, w))
	}

	return b.string(), itemLines
}

func sectionRule(title string, width int) string {
	label := " " + strings.ToUpper(title) + " "
	fill := max(0, width-2-lipgloss.Width(label))
	var b strings.Builder
	b.WriteString(styleFaint().Render("──"))
	b.WriteString(styleSection().Render(label))
	if fill > 0 {
		b.WriteString(styleFaint().Render(strings.Repeat("─", fill)))
	}
	return b.String()
}

func (m homeModel) caret(i int) string {
	open := m.isOpen(i)
	glyph := "▸"
	if open {
		glyph = "▾"
	}
	if i == m.cursor {
		return styleAccent().Render(glyph)
	}
	return styleFaint().Render(glyph)
}

func (m homeModel) markBlock(i int, block string, width int) string {
	lines := strings.Split(strings.TrimRight(block, "\n"), "\n")
	for j, line := range lines {
		n := lipgloss.Width(line)
		if width > 0 && n < width {
			lines[j] = line + strings.Repeat(" ", width-n)
		}
	}
	block = strings.Join(lines, "\n")
	if m.ctx.zone != nil {
		return m.ctx.zone.Mark(homeZoneID(i), block)
	}
	return block
}

func padRow(left, right string, width int) string {
	lw := lipgloss.Width(left)
	rw := lipgloss.Width(right)
	if rw == 0 {
		return left
	}
	gap := width - lw - rw
	if gap < 1 {
		return left
	}
	return left + strings.Repeat(" ", gap) + right
}

func (m homeModel) itemTitle(i int, name string) string {
	st := styleHeading()
	if i == m.cursor {
		st = styleAccent().Bold(true)
	}
	return st.Render(name)
}

func (m homeModel) renderProjectRow(i int, p Project, width int) string {
	indent := 3
	title := m.itemTitle(i, p.Name)
	kind := styleMuted().Render(p.Kind)
	head := padRow(m.caret(i)+"  "+title, kind, width)

	var b strings.Builder
	b.WriteString(head)
	b.WriteString("\n")
	if p.Tagline != "" {
		b.WriteString(wrapIndent(styleFaint(), indent, p.Tagline, width))
		b.WriteString("\n")
	}
	b.WriteString(wrapIndent(styleMuted(), indent, p.Summary, width))
	if !m.isOpen(i) {
		return m.markBlock(i, b.String(), width)
	}

	if projectHasLinks(p) {
		b.WriteString("\n")
	}
	for _, link := range projectLinks(p) {
		b.WriteString("\n")
		b.WriteString(wrapIndent(styleMuted(), indent, padLinkLabel(link.Label)+displayLinkValue(link.Value), width))
	}
	return m.markBlock(i, b.String(), width)
}

func projectHasLinks(p Project) bool {
	return len(p.Links) > 0
}

func projectLinks(p Project) []ProjectLink {
	return p.Links
}

func padLinkLabel(label string) string {
	const w = 16
	if n := lipgloss.Width(label); n >= w {
		return label + "  "
	}
	return fmt.Sprintf("%-*s", w, label)
}

func displayLinkValue(v string) string {
	v = strings.TrimPrefix(v, "https://")
	v = strings.TrimPrefix(v, "http://")
	return v
}

func fitIcon(s string) string {
	w := lipgloss.Width(s)
	if w >= 2 {
		return s
	}
	return s + strings.Repeat(" ", 2-w)
}

func renderContactRow(row contactRow, width int) string {
	icon := styleAccent().Render(fitIcon(row.Icon))
	prefix := icon + "  "
	pw := lipgloss.Width(prefix)
	wrapped := wrapText(styleBody(), row.Value, max(1, width-pw))
	lines := strings.Split(wrapped, "\n")
	pad := strings.Repeat(" ", pw)
	for i, line := range lines {
		line = strings.TrimRight(line, " ")
		if i == 0 {
			lines[i] = prefix + line
		} else {
			lines[i] = pad + line
		}
	}
	return strings.Join(lines, "\n")
}

func (m homeModel) renderExperienceRow(i int, e Experience, width int) string {
	indent := 3
	title := e.Role
	if e.Org != "" {
		title += " at " + e.Org
	}
	badge := ""
	if e.Current {
		badge = styleAccent().Render("current")
	}
	head := padRow(m.caret(i)+"  "+m.itemTitle(i, title), badge, width)

	var b strings.Builder
	b.WriteString(head)
	b.WriteString("\n")
	b.WriteString(wrapIndent(styleMuted(), indent, e.Period, width))
	if !m.isOpen(i) {
		return m.markBlock(i, b.String(), width)
	}

	b.WriteString("\n\n")
	b.WriteString(wrapIndent(styleBody(), indent, e.Summary, width))
	if len(e.Highlights) > 0 {
		b.WriteString("\n")
		for _, h := range e.Highlights {
			b.WriteString("\n")
			b.WriteString(wrapIndent(styleMuted(), indent, "•  "+h, width))
		}
	}
	if len(e.Tech) > 0 {
		b.WriteString("\n\n")
		b.WriteString(wrapIndent(styleMuted(), indent, strings.Join(e.Tech, "  ·  "), width))
	}
	return m.markBlock(i, b.String(), width)
}

func (m homeModel) renderEducation(width int) string {
	e := education()
	indent := 3
	var b strings.Builder
	b.WriteString(wrapIndent(styleHeading(), indent, e.Degree, width))
	b.WriteString("\n")
	b.WriteString(wrapIndent(styleBody(), indent, e.School, width))
	b.WriteString("\n")
	meta := e.Period
	if e.Note != "" {
		meta += "  ·  " + e.Note
	}
	b.WriteString(wrapIndent(styleMuted(), indent, meta, width))
	return b.String()
}

func (m homeModel) renderStackRow(i, n int, cat StackCategory, width int) string {
	indent := 3
	label := fmt.Sprintf("%02d  %s", n, cat.Name)
	count := styleMuted().Render(fmt.Sprintf("×%d", len(cat.Items)))
	head := padRow(m.caret(i)+"  "+m.itemTitle(i, label), count, width)
	if !m.isOpen(i) {
		return m.markBlock(i, head, width)
	}
	var b strings.Builder
	b.WriteString(head)
	b.WriteString("\n")
	b.WriteString(wrapIndent(styleMuted(), indent, strings.Join(cat.Items, "  ·  "), width))
	return m.markBlock(i, b.String(), width)
}
