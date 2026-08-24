package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func sizedApp(w, h int) appModel {
	m := newAppModel()
	m.showSplash = false
	mod, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return mod.(appModel)
}

func sizedSplash(w, h int) appModel {
	m := newAppModel()
	mod, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return mod.(appModel)
}

func viewLines(m appModel) []string {
	plain := ansi.Strip(m.View().Content)
	return strings.Split(strings.TrimRight(plain, "\n"), "\n")
}

func TestLayoutFitsTerminal(t *testing.T) {
	const w, h = 80, 24
	m := sizedApp(w, h)
	lines := viewLines(m)
	if len(lines) != h {
		t.Fatalf("height: got %d want %d\n%s", len(lines), h, strings.Join(lines, "\n"))
	}
	for i, line := range lines {
		if got := ansi.StringWidth(line); got != w {
			t.Errorf("line %d width=%d want=%d: %q", i, got, w, line)
		}
	}

	m = m.toggleNav()
	m.setSizes()
	collapsed := viewLines(m)
	if len(collapsed) != h {
		t.Fatalf("collapsed height: got %d want %d", len(collapsed), h)
	}
	for i, line := range collapsed {
		if got := ansi.StringWidth(line); got != w {
			t.Errorf("collapsed line %d width=%d want=%d: %q", i, got, w, line)
		}
	}
}

func TestNavIconsSameWidth(t *testing.T) {
	want := -1
	for _, page := range navPages() {
		got := lipgloss.Width(page.Icon())
		if want < 0 {
			want = got
		} else if got != want {
			t.Errorf("%s icon width=%d want=%d (%q)", page, got, want, page.Icon())
		}
		if got != navIconWidth {
			t.Errorf("%s icon width=%d, navIconWidth=%d (%q)", page, got, navIconWidth, page.Icon())
		}
	}
}

func TestNavLabelsShareColumn(t *testing.T) {
	m := sizedApp(80, 24)
	lines := viewLines(m)
	col := -1
	for _, label := range []string{"home", "projects", "work", "stack", "blogs"} {
		found := false
		for _, line := range lines {
			i := strings.Index(line, label)
			if i < 0 {
				continue
			}
			found = true
			start := ansi.StringWidth(line[:i])
			if col < 0 {
				col = start
			} else if start != col {
				t.Errorf("%q starts at col %d, want %d\n%s", label, start, col, line)
			}
			break
		}
		if !found {
			t.Errorf("missing nav label %q", label)
		}
	}
}

func TestCollapsedNavNarrower(t *testing.T) {
	exp := sidebarFrameWidth(false)
	col := sidebarFrameWidth(true)
	if col >= exp {
		t.Fatalf("collapsed %d should be < expanded %d", col, exp)
	}
}

func TestHomeHasWebsiteSections(t *testing.T) {
	m := sizedApp(80, 24)
	plain := strings.ToLower(ansi.Strip(m.home.viewport.GetContent()))
	for _, section := range []string{"projects", "experience", "education", "stack", "contact"} {
		if !strings.Contains(plain, section) {
			t.Errorf("home missing section %q", section)
		}
	}
	for _, needle := range []string{"scribblesvg", "varipane", "indie hacking", "national institute of technology patna", "navsari"} {
		if !strings.Contains(plain, needle) {
			t.Errorf("home missing %q", needle)
		}
	}
}

func TestHomeScrollbarWhenOverflow(t *testing.T) {
	m := sizedApp(80, 24)
	if m.home.viewport.TotalLineCount() <= m.home.viewport.Height() {
		t.Fatalf("home should overflow at 80x24, got %d lines in height %d",
			m.home.viewport.TotalLineCount(), m.home.viewport.Height())
	}
	plain := ansi.Strip(m.View().Content)
	if !strings.Contains(plain, scrollbarThumb) {
		t.Fatalf("expected scrollbar thumb on overflowing home")
	}
}

func TestHomeMovesSelectionWithoutWrap(t *testing.T) {
	m := sizedApp(80, 24)
	if m.home.cursor != 0 {
		t.Fatalf("start cursor=%d want 0", m.home.cursor)
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.home.cursor != 1 {
		t.Fatalf("j should move to the next row, cursor=%d", m.home.cursor)
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	if m.home.cursor != 0 {
		t.Fatalf("k at first row should stay put, cursor=%d", m.home.cursor)
	}

	last := len(m.home.items) - 1
	m.home.cursor = last
	m.home.refresh()
	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.home.cursor != last {
		t.Fatalf("j at last row should stay put, cursor=%d want %d", m.home.cursor, last)
	}
}

func TestHomeExpandTogglesDetail(t *testing.T) {
	m := sizedApp(80, 24)
	before := ansi.Strip(m.home.viewport.GetContent())
	if strings.Contains(before, "@scribblesvg/core") {
		t.Fatalf("collapsed home should not show ScribbleSVG npm links")
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	after := ansi.Strip(m.home.viewport.GetContent())
	if !strings.Contains(after, "@scribblesvg/core") {
		t.Fatalf("enter should expand the selected project")
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	closed := ansi.Strip(m.home.viewport.GetContent())
	if strings.Contains(closed, "@scribblesvg/core") {
		t.Fatalf("second enter should collapse the project")
	}
}

func TestProjectDetailMatchesSiteCopy(t *testing.T) {
	m := sizedApp(80, 24)
	m, _ = m.navigateTo(PageProjects)
	m.setSizes()
	m.projects, _ = m.projects.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	plain := strings.ToLower(ansi.Strip(m.projects.viewport.GetContent()))
	needles := []string{
		"scribblesvg",
		"hand-drawn diagramming toolkit",
		"what",
		"a featherlight diagramming toolkit",
		"why",
		"i plan and reason systems by drawing",
		"where",
		"scribblesvg.ashutoshbind.com",
		"technicals",
		"layered svg on a simple state-driven canvas",
		"upcoming",
		"bring-your-own fonts",
	}
	for _, needle := range needles {
		if !strings.Contains(plain, needle) {
			t.Errorf("project detail missing %q", needle)
		}
	}
	if strings.Contains(plain, "gallery") {
		t.Fatal("project detail should not include site gallery copy")
	}
}

func TestProjectDetailFitsTerminal(t *testing.T) {
	const w, h = 80, 24
	m := sizedApp(w, h)
	m, _ = m.navigateTo(PageProjects)
	m.setSizes()
	m.projects, _ = m.projects.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.setSizes()
	lines := viewLines(m)
	if len(lines) != h {
		t.Fatalf("project detail height: got %d want %d\n%s", len(lines), h, strings.Join(lines, "\n"))
	}
	for i, line := range lines {
		if got := ansi.StringWidth(line); got != w {
			t.Errorf("line %d width=%d want=%d: %q", i, got, w, line)
		}
	}
}

func TestBlogScrollbarWhenOverflow(t *testing.T) {
	m := sizedApp(80, 24)
	m, _ = m.navigateTo(PageBlogs)
	m.setSizes()
	if it, ok := m.blogs.list.SelectedItem().(blogItem); ok {
		m.blogs.openID = it.blog.ID
		m.blogs.refreshDetail()
	}
	if m.blogs.viewport.TotalLineCount() <= m.blogs.viewport.Height() {
		t.Skip("blog fits without a scrollbar in this size")
	}
	plain := ansi.Strip(m.View().Content)
	if !strings.Contains(plain, scrollbarThumb) {
		t.Fatalf("expected scrollbar thumb in overflowing blog view")
	}
}

func TestScrollbarHugsContentEdge(t *testing.T) {
	const w, h = 80, 24
	m := sizedApp(w, h)
	wantCol := m.columnWidth() - 1
	found := false
	for _, line := range viewLines(m) {
		idx := strings.Index(line, scrollbarThumb)
		if idx < 0 {
			continue
		}
		found = true
		col := ansi.StringWidth(line[:idx])
		if col != wantCol {
			t.Errorf("thumb at col %d, want %d (right edge of content)\n%s", col, wantCol, line)
		}
	}
	if !found {
		t.Fatal("expected overlay scrollbar thumb on overflowing home")
	}
}

func TestThemeSwitcherInSidebar(t *testing.T) {
	m := sizedApp(80, 24)
	plain := strings.ToLower(ansi.Strip(m.View().Content))
	if !strings.Contains(plain, "theme") {
		t.Fatal("expected theme label in sidebar")
	}
	if !strings.Contains(plain, "pumice") {
		t.Fatal("expected current theme name in sidebar")
	}
	if !strings.Contains(ansi.Strip(m.View().Content), "●") {
		t.Fatal("expected theme swatches in sidebar")
	}
	if !strings.Contains(ansi.Strip(m.View().Content), "◆") {
		t.Fatal("expected active theme marker in sidebar")
	}
}

func TestSplashFitsTerminal(t *testing.T) {
	sizes := [][2]int{{80, 24}, {100, 30}, {60, 20}, {40, 16}}
	for _, size := range sizes {
		w, h := size[0], size[1]
		m := sizedSplash(w, h)
		if !m.showSplash {
			t.Fatal("app should open on the intro")
		}
		lines := viewLines(m)
		if len(lines) != h {
			t.Fatalf("%dx%d splash height: got %d want %d\n%s", w, h, len(lines), h, strings.Join(lines, "\n"))
		}
		for i, line := range lines {
			if got := ansi.StringWidth(line); got != w {
				t.Errorf("%dx%d splash line %d width=%d want=%d: %q", w, h, i, got, w, line)
			}
		}
	}
}

func TestSplashShowsNameAndEnterHint(t *testing.T) {
	m := sizedSplash(80, 24)
	plain := ansi.Strip(m.View().Content)
	if !strings.Contains(plain, "/___/") && !strings.Contains(plain, "_____") {
		t.Fatalf("expected big-name banner on splash\n%s", plain)
	}
	if !strings.Contains(strings.ToLower(plain), "enter") {
		t.Fatal("expected enter hint on splash")
	}
	if !strings.Contains(plain, "\u2800") && !hasBraille(plain) {
		t.Fatal("expected braille-dot waves on splash")
	}
}

func hasBraille(s string) bool {
	for _, r := range s {
		if r >= 0x2800 && r <= 0x28FF {
			return true
		}
	}
	return false
}

func TestEnterLeavesSplash(t *testing.T) {
	m := sizedSplash(80, 24)
	mod, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = mod.(appModel)
	if m.showSplash {
		t.Fatal("enter should leave the intro")
	}
	plain := strings.ToLower(ansi.Strip(m.View().Content))
	if !strings.Contains(plain, "home") {
		t.Fatal("expected portfolio nav after enter")
	}
}

func TestBrandReopensSplash(t *testing.T) {
	m := sizedApp(80, 24)
	if m.showSplash {
		t.Fatal("sizedApp should start in the portfolio")
	}
	next, _ := m.openSplash()
	if !next.showSplash {
		t.Fatal("brand click should return to the intro")
	}
}

func TestThemeCycleKeepsLayout(t *testing.T) {
	const w, h = 80, 24
	m := sizedApp(w, h)
	mod, _ := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
	m = mod.(appModel)
	if m.themeID == defaultThemeID {
		t.Fatal("] should cycle off the default theme")
	}
	plain := strings.ToLower(ansi.Strip(m.View().Content))
	if !strings.Contains(plain, strings.ToLower(currentTheme().Name)) {
		t.Fatalf("expected cycled theme name %q in sidebar", currentTheme().Name)
	}
	lines := viewLines(m)
	if len(lines) != h {
		t.Fatalf("height after theme cycle: got %d want %d", len(lines), h)
	}
	for i, line := range lines {
		if got := ansi.StringWidth(line); got != w {
			t.Errorf("line %d width=%d want=%d: %q", i, got, w, line)
		}
	}

	m = m.toggleNav()
	m.setSizes()
	collapsed := ansi.Strip(m.View().Content)
	if !strings.Contains(collapsed, "◆") {
		t.Fatal("expected collapsed theme swatch")
	}

	for range allThemes() {
		mod, _ := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
		m = mod.(appModel)
		lines := viewLines(m)
		if len(lines) != h {
			t.Fatalf("theme %s height=%d want %d", m.themeID, len(lines), h)
		}
		for i, line := range lines {
			if got := ansi.StringWidth(line); got != w {
				t.Fatalf("theme %s line %d width=%d want=%d", m.themeID, i, got, w)
			}
		}
	}
}
