package main

import (
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	zoneProjects = "projects-item-"
	zoneBlogs    = "blogs-item-"
	zoneWork     = "experience-item-"
	zoneStack    = "stack-line-"
)

func listZoneID(prefix string, i int) string {
	return prefix + strconv.Itoa(i)
}

func stackZoneID(i int) string {
	return zoneStack + strconv.Itoa(i)
}

func hitListIndex(z *zone.Manager, prefix string, n int, msg tea.MouseMsg) int {
	if z == nil {
		return -1
	}
	for i := 0; i < n; i++ {
		if z.Get(listZoneID(prefix, i)).InBounds(msg) {
			return i
		}
	}
	return -1
}

func handleListPointer(z *zone.Manager, prefix string, l *list.Model, msg tea.Msg) (consumed, open bool) {
	switch msg := msg.(type) {
	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelDown:
			l.CursorDown()
			return true, false
		case tea.MouseWheelUp:
			l.CursorUp()
			return true, false
		}
	case tea.MouseMotionMsg:
		idx := hitListIndex(z, prefix, len(l.Items()), msg)
		if idx < 0 {
			return false, false
		}
		if l.Index() != idx {
			l.Select(idx)
		}
		return true, false
	case tea.MouseClickMsg:
		idx := hitListIndex(z, prefix, len(l.Items()), msg)
		if idx < 0 {
			return false, false
		}
		l.Select(idx)
		return true, true
	}
	return false, false
}
