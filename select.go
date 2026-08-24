package main

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type textSel struct {
	originX, originY int
	x, y             int
	down             bool
	dragging         bool
	active           bool
}

func (s textSel) live() bool {
	return s.dragging || s.active
}

func (s *textSel) press(x, y int) {
	s.originX, s.originY = x, y
	s.x, s.y = x, y
	s.down = true
	s.dragging = false
	s.active = false
}

func (s *textSel) drag(x, y int) {
	s.x, s.y = x, y
	if abs(x-s.originX)+abs(y-s.originY) >= 1 {
		s.dragging = true
		s.active = true
	}
}

func (s *textSel) clear() {
	s.down = false
	s.dragging = false
	s.active = false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func applySelection(content string, ax, ay, bx, by int) (string, string) {
	if ay > by || (ay == by && ax > bx) {
		ax, ay, bx, by = bx, by, ax, ay
	}
	lines := strings.Split(content, "\n")
	out := make([]string, len(lines))
	var selected []string
	hi := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorAccent))

	for i, line := range lines {
		w := ansi.StringWidth(line)
		if i < ay || i > by {
			out[i] = line
			continue
		}

		start, end := 0, w
		switch {
		case ay == by:
			start, end = ax, bx+1
		case i == ay:
			start, end = ax, w
		case i == by:
			start, end = 0, bx+1
		}
		if start < 0 {
			start = 0
		}
		if end > w {
			end = w
		}
		if start >= end {
			out[i] = line
			continue
		}

		left := ansi.Cut(line, 0, start)
		mid := ansi.Cut(line, start, end)
		right := ansi.Cut(line, end, w)
		plain := ansi.Strip(mid)
		selected = append(selected, strings.TrimRight(plain, " "))
		out[i] = left + hi.Render(plain) + right
	}

	text := strings.TrimSpace(strings.Join(selected, "\n"))
	return strings.Join(out, "\n"), text
}
