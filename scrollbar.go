package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

const (
	scrollbarGutter = 1
	scrollbarThumb  = "▐"
	scrollbarTrack  = "▕"
)

func sizeViewport(vp *viewport.Model, width, height int) {
	vp.SetWidth(max(0, width-scrollbarGutter))
	vp.SetHeight(max(0, height))
	vp.MouseWheelDelta = 1
}

func viewWithScrollbar(vp viewport.Model) string {
	h := vp.Height()
	w := vp.Width()
	view := padBlock(vp.View(), w, h)
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
		thumb := styleScrollbarThumb().Render(scrollbarThumb)
		track := styleScrollbarTrack().Render(scrollbarTrack)
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

// padBlock makes a rectangular block so a right-edge gutter stays pinned
// instead of clinging to the last glyph of each line.
func padBlock(s string, width, height int) string {
	if height <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	if width <= 0 {
		return strings.Join(lines, "\n")
	}
	for i, line := range lines {
		n := lipgloss.Width(line)
		if n < width {
			lines[i] = line + strings.Repeat(" ", width-n)
		}
	}
	return strings.Join(lines, "\n")
}
