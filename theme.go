package main

import (
	"fmt"
	"image/color"
	"strings"
)

// Original palettes for this app. Names and hex values are not taken from
// trademarked terminal themes.
type Theme struct {
	ID      string
	Name    string
	Light   bool
	Accent  string
	Bright  string
	Text    string
	Muted   string
	Faint   string
	Border  string
	Code    string
	CodeBg  string
	Bg      string
	Surface string
	Warn    string
}

const defaultThemeID = "pumice"

var activeThemeID = defaultThemeID

func allThemes() []Theme {
	return []Theme{
		{
			ID:      "pumice",
			Name:    "Pumice",
			Accent:  "#6A9E96",
			Bright:  "#E8E4D9",
			Text:    "#C9C3B6",
			Muted:   "#8F897C",
			Faint:   "#5E5A52",
			Border:  "#3F3C37",
			Code:    "#C9A36A",
			CodeBg:  "#2A2824",
			Bg:      "#1A1917",
			Surface: "#2A2824",
			Warn:    "#C17B6A",
		},
		{
			ID:      "phosphor",
			Name:    "Phosphor",
			Accent:  "#4FD67A",
			Bright:  "#D2F5DC",
			Text:    "#9BC9A6",
			Muted:   "#5E8A68",
			Faint:   "#3A5740",
			Border:  "#1E3324",
			Code:    "#B4E06A",
			CodeBg:  "#122016",
			Bg:      "#0A110C",
			Surface: "#122016",
			Warn:    "#D4A056",
		},
		{
			ID:      "ember",
			Name:    "Ember",
			Accent:  "#E39456",
			Bright:  "#F4E6D6",
			Text:    "#D2C0AC",
			Muted:   "#9A8470",
			Faint:   "#685848",
			Border:  "#3A2E26",
			Code:    "#E8C07A",
			CodeBg:  "#241C18",
			Bg:      "#171210",
			Surface: "#241C18",
			Warn:    "#E07050",
		},
		{
			ID:      "harbor",
			Name:    "Harbor",
			Accent:  "#5AAFBE",
			Bright:  "#E4F0F4",
			Text:    "#B4C8D2",
			Muted:   "#748A96",
			Faint:   "#4A5E68",
			Border:  "#283840",
			Code:    "#7EC8A0",
			CodeBg:  "#162028",
			Bg:      "#0D141A",
			Surface: "#162028",
			Warn:    "#D08060",
		},
		{
			ID:      "orchid",
			Name:    "Orchid",
			Accent:  "#C486B4",
			Bright:  "#F0E6F2",
			Text:    "#CEC2D4",
			Muted:   "#948A9C",
			Faint:   "#605868",
			Border:  "#383040",
			Code:    "#E0A878",
			CodeBg:  "#1E1824",
			Bg:      "#141018",
			Surface: "#1E1824",
			Warn:    "#D07070",
		},
		{
			ID:      "parchment",
			Name:    "Parchment",
			Light:   true,
			Accent:  "#2F6F66",
			Bright:  "#1F1C18",
			Text:    "#3A352E",
			Muted:   "#6E675C",
			Faint:   "#9A9286",
			Border:  "#D0C8B8",
			Code:    "#8B5A2B",
			CodeBg:  "#E7E0D2",
			Bg:      "#F3EEE3",
			Surface: "#E7E0D2",
			Warn:    "#A04030",
		},
	}
}

func themeByID(id string) (Theme, bool) {
	for _, t := range allThemes() {
		if t.ID == id {
			return t, true
		}
	}
	return Theme{}, false
}

func currentTheme() Theme {
	if t, ok := themeByID(activeThemeID); ok {
		return t
	}
	t, _ := themeByID(defaultThemeID)
	return t
}

func cycleThemeID(id string, delta int) string {
	themes := allThemes()
	idx := 0
	for i, t := range themes {
		if t.ID == id {
			idx = i
			break
		}
	}
	n := len(themes)
	return themes[(idx+delta%n+n)%n].ID
}

func applyTheme(t Theme) {
	activeThemeID = t.ID
	colorAccent = t.Accent
	colorBright = t.Bright
	colorText = t.Text
	colorMuted = t.Muted
	colorFaint = t.Faint
	colorBorder = t.Border
	colorCode = t.Code
	colorCodeBg = t.CodeBg
	colorBg = t.Bg
	colorSurface = t.Surface
	colorWarn = t.Warn
}

func parseHexColor(s string) color.Color {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 6 {
		return nil
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b); err != nil {
		return nil
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
