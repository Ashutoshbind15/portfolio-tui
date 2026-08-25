package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// Active palette; applyTheme swaps these. Seeded with Pumice so styles
// work before the first model is constructed.
var (
	colorAccent  = "#6A9E96"
	colorBright  = "#E8E4D9"
	colorText    = "#C9C3B6"
	colorMuted   = "#8F897C"
	colorFaint   = "#5E5A52"
	colorBorder  = "#3F3C37"
	colorCode    = "#C9A36A"
	colorCodeBg  = "#2A2824"
	colorBg      = "#1A1917"
	colorSurface = "#2A2824"
	colorWarn    = "#C17B6A"
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
		Background(lipgloss.Color(colorSurface))
}

func styleNavIcon(active, hover bool) lipgloss.Style {
	s := lipgloss.NewStyle().
		Width(navIconWidth).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center)
	if active || hover {
		s = s.Background(lipgloss.Color(colorSurface))
	}
	return s
}

func styleNavGap(active, hover bool) lipgloss.Style {
	s := lipgloss.NewStyle().Width(1).Height(1)
	if active || hover {
		s = s.Background(lipgloss.Color(colorSurface))
	}
	return s
}

func styleNavLabel(active, hover bool) lipgloss.Style {
	if active {
		return styleNavActive()
	}
	if hover {
		return styleNavIdle().
			Foreground(lipgloss.Color(colorBright)).
			Background(lipgloss.Color(colorSurface))
	}
	return styleNavIdle()
}

func styleHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorBg)).
		Padding(1, 2, 1, 2)
}

func styleFooter() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorFaint)).
		Background(lipgloss.Color(colorBg)).
		Padding(0, 2)
}

func styleScrollbarThumb() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
}

func styleScrollbarTrack() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorder))
}

func styleApp() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorBg))
}

func styleHeading() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBright)).
		Bold(true)
}

func styleSection() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
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
	var s help.Styles
	if currentTheme().Light {
		s = help.DefaultLightStyles()
	} else {
		s = help.DefaultDarkStyles()
	}
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
	s := list.NewDefaultItemStyles(!currentTheme().Light)
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

type zonedDelegate struct {
	list.DefaultDelegate
	zone   *zone.Manager
	prefix string
}

func (d zonedDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var b strings.Builder
	d.DefaultDelegate.Render(&b, m, index, item)
	out := b.String()
	if d.zone != nil {
		out = d.zone.Mark(listZoneID(d.prefix, index), out)
	}
	fmt.Fprint(w, out)
}

func newItemList(title string, items []list.Item, width, height int, z *zone.Manager, prefix string) list.Model {
	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	d.Styles = itemStyles()
	l := list.New(items, zonedDelegate{DefaultDelegate: d, zone: z, prefix: prefix}, width, height)
	l.Title = title
	l.Styles = list.DefaultStyles(!currentTheme().Light)
	l.Styles.Title = styleHeading()
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return l
}

func restyleList(l *list.Model, z *zone.Manager, prefix string) {
	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	d.Styles = itemStyles()
	l.SetDelegate(zonedDelegate{DefaultDelegate: d, zone: z, prefix: prefix})
	l.Styles = list.DefaultStyles(!currentTheme().Light)
	l.Styles.Title = styleHeading()
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
