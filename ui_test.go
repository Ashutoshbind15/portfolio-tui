package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestNavToggleCenteredInSidebar(t *testing.T) {
	m := sizedApp(80, 24)
	lines := viewLines(m)
	mid := m.columnWidth() + sidebarBorderW + sidebarPadX + sidebarExpandedInner/2
	found := false
	for i := len(lines) - 1; i >= 0; i-- {
		idx := strings.LastIndex(lines[i], "<")
		if idx < 0 {
			continue
		}
		col := ansi.StringWidth(lines[i][:idx])
		if col < m.columnWidth() {
			continue
		}
		found = true
		if abs(col-mid) > 1 {
			t.Fatalf("toggle col %d, want ~%d\n%s", col, mid, lines[i])
		}
		break
	}
	if !found {
		t.Fatal("missing expanded nav toggle")
	}
}

func TestContactUsesBrandIcons(t *testing.T) {
	m := sizedApp(80, 24)
	plain := ansi.Strip(m.home.viewport.GetContent())
	for _, needle := range []string{
		"📧", "🐙", "𝕏", "🔗", "🌐", ">_",
		"dev@ashutoshbind.com",
		"github.com/Ashutoshbind15",
		"x.com/Ashutosh_Bind15",
		"linkedin.com/in/ashutosh-bind-56806b22b",
		"ashutoshbind.com",
		"go run github.com/Ashutoshbind15/portfolio-tui@main",
	} {
		if !strings.Contains(plain, needle) {
			t.Errorf("contact missing %q", needle)
		}
	}
	if strings.Contains(plain, "$  ssh") || strings.Contains(plain, "@  dev@") {
		t.Fatal("old contact glyphs should be gone")
	}
}

func TestProjectDetailSectionHeadings(t *testing.T) {
	m := sizedApp(80, 24)
	m, _ = m.navigateTo(PageProjects)
	m.setSizes()
	m.projects, _ = m.projects.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	plain := ansi.Strip(m.projects.viewport.GetContent())
	for _, heading := range []string{"WHAT", "WHY", "WHERE", "TECHNICALS", "UPCOMING"} {
		if !strings.Contains(plain, heading) {
			t.Errorf("project detail missing heading %q", heading)
		}
	}
	if !strings.Contains(plain, "varypane.com") {
		t.Error("VaryPane detail missing site link")
	}
	if !strings.Contains(strings.ToLower(plain), "site") {
		t.Error("VaryPane detail missing site label")
	}
}

func TestSelectionExtractsText(t *testing.T) {
	content := "hello world\nfoo bar baz"
	_, got := applySelection(content, 6, 0, 10, 0)
	if got != "world" {
		t.Fatalf("got %q want world", got)
	}

	highlighted, got := applySelection(content, 0, 0, 2, 1)
	if !strings.Contains(got, "hello") || !strings.Contains(got, "foo") {
		t.Fatalf("multi-line got %q", got)
	}
	lines := strings.Split(highlighted, "\n")
	if ansi.StringWidth(lines[0]) != ansi.StringWidth(strings.Split(content, "\n")[0]) {
		t.Fatalf("highlight changed line width")
	}
}
