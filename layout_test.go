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
	for _, section := range []string{"software", "experience", "education", "stack", "contact"} {
		if !strings.Contains(plain, section) {
			t.Errorf("home missing section %q", section)
		}
	}
	for _, needle := range []string{"varypane", "indie hacking", "national institute of technology patna", "navsari"} {
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
	if !strings.Contains(plain, "█") {
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
	if strings.Contains(before, "starter templates") {
		t.Fatalf("collapsed home should not show VaryPane detail")
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	after := ansi.Strip(m.home.viewport.GetContent())
	if !strings.Contains(after, "starter templates") {
		t.Fatalf("enter should expand the selected project")
	}

	m.home, _ = m.home.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	closed := ansi.Strip(m.home.viewport.GetContent())
	if strings.Contains(closed, "starter templates") {
		t.Fatalf("second enter should collapse the project")
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
	if !strings.Contains(plain, "█") {
		t.Fatalf("expected scrollbar thumb in overflowing blog view")
	}
}
