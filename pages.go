package main

import tea "charm.land/bubbletea/v2"

// Page is a top-level route in the portfolio TUI.
type Page string

const (
	PageHome       Page = "home"
	PageSelect     Page = "select"
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
	case PageSelect:
		return "Pages"
	default:
		return string(p)
	}
}

func (p Page) Description() string {
	switch p {
	case PageHome:
		return "Profile and contact"
	case PageProjects:
		return "Open-source and shipped work"
	case PageExperience:
		return "Roles and internships"
	case PageStack:
		return "Tools grouped by area"
	case PageBlogs:
		return "Short notes"
	default:
		return ""
	}
}

func navPages() []Page {
	return []Page{PageHome, PageProjects, PageExperience, PageStack, PageBlogs}
}

// navigateMsg asks the root model to switch pages.
type navigateMsg struct {
	page Page
}

func navigateTo(page Page) tea.Cmd {
	return func() tea.Msg { return navigateMsg{page: page} }
}

// closeMenuMsg pops the page-select overlay.
type closeMenuMsg struct{}
