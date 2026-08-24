package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

// Palette leans on Term Chess (title pill, selected rail) while staying
// as quiet as the web portfolio.
const (
	colorAccent     = "62"
	colorTitleFg    = "230"
	colorHighlight  = "213"
	colorSelected   = "212"
	colorText       = "252"
	colorMuted      = "245"
	colorFaint      = "241"
	colorBorder     = "238"
)

func styleHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorTitleFg)).
		Background(lipgloss.Color(colorAccent)).
		Padding(0, 1)
}

func styleBrand() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229"))
}

func styleNavIdle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
}

func styleNavActive() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorHighlight)).
		Bold(true)
}

func styleChromeBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color(colorAccent)).
		Padding(0, 1)
}

func styleFooter() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorFaint)).
		Padding(0, 1)
}

func styleHeading() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Bold(true)
}

func styleMuted() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
}

func styleFaint() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorFaint))
}

func styleBody() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
}

func styleAccent() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorHighlight))
}

func styleChip() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(0, 1).
		MarginRight(1)
}

func newItemList(title string, items []list.Item, width, height int) list.Model {
	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	l := list.New(items, d, width, height)
	l.Title = title
	l.Styles = list.DefaultStyles(true)
	l.Styles.Title = styleHeader()
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return l
}

func joinChips(labels []string) string {
	if len(labels) == 0 {
		return ""
	}
	parts := make([]string, 0, len(labels))
	for _, label := range labels {
		parts = append(parts, styleChip().Render(label))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
