package main

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

const scrollbarGutter = 1

func sizeViewport(vp *viewport.Model, width, height int) {
	vp.SetWidth(max(0, width-scrollbarGutter))
	vp.SetHeight(max(0, height))
}

func viewWithScrollbar(vp viewport.Model) string {
	view := vp.View()
	h := vp.Height()
	if h <= 0 {
		return view
	}

	bars := make([]string, h)
	total := vp.TotalLineCount()
	if total > h {
		thumbSize := max(1, h*h/total)
		if thumbSize > h {
			thumbSize = h
		}
		maxStart := h - thumbSize
		start := int(vp.ScrollPercent() * float64(maxStart))
		if vp.AtBottom() {
			start = maxStart
		}
		if start < 0 {
			start = 0
		}
		thumb := styleScrollbarThumb().Render("█")
		track := styleScrollbarTrack().Render("░")
		for i := 0; i < h; i++ {
			if i >= start && i < start+thumbSize {
				bars[i] = thumb
			} else {
				bars[i] = track
			}
		}
	} else {
		for i := 0; i < h; i++ {
			bars[i] = " "
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, view, lipgloss.JoinVertical(lipgloss.Left, bars...))
}
