package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

// Hex palette: one accent, warm paper text, stone secondaries.
// Mid-saturation so it downsamples cleanly on 256-color and 16-color
// terminals instead of blowing out to neon.
const (
	colorAccent = "#6A9E96"
	colorBright = "#E8E4D9"
	colorText   = "#C9C3B6"
	colorMuted  = "#8F897C"
	colorFaint  = "#5E5A52"
	colorBorder = "#3F3C37"
	colorCode   = "#C9A36A"
	colorCodeBg = "#2A2824"
)

func styleBrand() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorBright))
}

func styleNavIdle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
}

func styleNavActive() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		Background(lipgloss.Color(colorCodeBg))
}

func styleNavIcon(active bool) lipgloss.Style {
	s := lipgloss.NewStyle().
		Width(navIconWidth).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center)
	if active {
		s = s.Background(lipgloss.Color(colorCodeBg))
	}
	return s
}

func styleNavGap(active bool) lipgloss.Style {
	s := lipgloss.NewStyle().Width(1).Height(1)
	if active {
		s = s.Background(lipgloss.Color(colorCodeBg))
	}
	return s
}

func styleNavLabel(active bool) lipgloss.Style {
	if active {
		return styleNavActive()
	}
	return styleNavIdle()
}

func styleHeader() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 1, 2)
}

func styleFooter() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorFaint)).
		Padding(0, 2)
}

func styleScrollbarThumb() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
}

func styleScrollbarTrack() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorder))
}

func styleHeading() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBright)).
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
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
}

func newHelp() help.Model {
	h := help.New()
	s := help.DefaultDarkStyles()
	s.ShortKey = s.ShortKey.Foreground(lipgloss.Color(colorMuted))
	s.ShortDesc = s.ShortDesc.Foreground(lipgloss.Color(colorFaint))
	s.ShortSeparator = s.ShortSeparator.Foreground(lipgloss.Color(colorBorder))
	s.FullKey = s.ShortKey
	s.FullDesc = s.ShortDesc
	s.FullSeparator = s.ShortSeparator
	s.Ellipsis = s.ShortSeparator
	h.Styles = s
	return h
}

func styleChip() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorBorder)).
		Padding(0, 1).
		MarginRight(1)
}

func itemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles(true)
	s.NormalTitle = s.NormalTitle.Foreground(lipgloss.Color(colorText))
	s.NormalDesc = s.NormalDesc.Foreground(lipgloss.Color(colorMuted))
	s.SelectedTitle = s.SelectedTitle.
		BorderForeground(lipgloss.Color(colorAccent)).
		Foreground(lipgloss.Color(colorBright))
	s.SelectedDesc = s.SelectedDesc.
		BorderForeground(lipgloss.Color(colorAccent)).
		Foreground(lipgloss.Color(colorMuted))
	s.DimmedTitle = s.DimmedTitle.Foreground(lipgloss.Color(colorFaint))
	s.DimmedDesc = s.DimmedDesc.Foreground(lipgloss.Color(colorFaint))
	return s
}

func newItemList(title string, items []list.Item, width, height int) list.Model {
	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	d.Styles = itemStyles()
	l := list.New(items, d, width, height)
	l.Title = title
	l.Styles = list.DefaultStyles(true)
	l.Styles.Title = styleHeading()
	l.SetShowTitle(false)
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
