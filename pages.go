package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Page is a top-level route in the portfolio TUI.
type Page string

const (
	PageHome       Page = "home"
	PageProjects   Page = "projects"
	PageExperience Page = "experience"
	PageStack      Page = "stack"
	PageBlogs      Page = "blogs"
)

func (p Page) Title() string {
	switch p {
	case PageHome:
		return "Home"
	case PageProjects:
		return "Projects"
	case PageExperience:
		return "Experience"
	case PageStack:
		return "Stack"
	case PageBlogs:
		return "Blogs"
	default:
		return string(p)
	}
}

// NavLabel is the short sidebar text. "experience" is too wide for the rail.
func (p Page) NavLabel() string {
	switch p {
	case PageExperience:
		return "work"
	default:
		return strings.ToLower(p.Title())
	}
}

// Icon is a Unicode emoji (no extra font license). All are 2 cells wide.
func (p Page) Icon() string {
	switch p {
	case PageHome:
		return "🏠"
	case PageProjects:
		return "📦"
	case PageExperience:
		return "💼"
	case PageStack:
		return "💻"
	case PageBlogs:
		return "📝"
	default:
		return "•"
	}
}

func navPages() []Page {
	return []Page{PageHome, PageProjects, PageExperience, PageStack, PageBlogs}
}

func cyclePage(current Page, delta int) Page {
	pages := navPages()
	idx := 0
	for i, page := range pages {
		if page == current {
			idx = i
			break
		}
	}
	n := len(pages)
	return pages[(idx+delta%n+n)%n]
}

// navigateMsg asks the root model to switch pages.
type navigateMsg struct {
	page Page
}

func navigateTo(page Page) tea.Cmd {
	return func() tea.Msg { return navigateMsg{page: page} }
}
