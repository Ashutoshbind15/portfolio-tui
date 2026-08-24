package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestNavToggleAlignsWithIcons(t *testing.T) {
	m := sizedApp(80, 24)
	lines := viewLines(m)
	iconCol := -1
	for _, line := range lines {
		i := strings.Index(line, "🏠")
		if i < 0 {
			continue
		}
		iconCol = ansi.StringWidth(line[:i])
		break
	}
	if iconCol < 0 {
		t.Fatal("missing home icon")
	}

	toggleCol := -1
	for _, line := range lines {
		i := strings.LastIndex(line, "<")
		if i < 0 {
			continue
		}
		toggleCol = ansi.StringWidth(line[:i])
	}
	if toggleCol < 0 {
		t.Fatal("missing expanded nav toggle")
	}
	if toggleCol != iconCol && toggleCol != iconCol+1 {
		t.Fatalf("toggle col %d, icon col %d", toggleCol, iconCol)
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
