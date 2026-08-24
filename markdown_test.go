package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarkdownCodeFollowsTheme(t *testing.T) {
	t.Cleanup(func() {
		applyTheme(mustTheme(t, defaultThemeID))
	})

	src := "```go\nfunc main() {}\n```"
	applyTheme(mustTheme(t, "pumice"))
	pumice := renderMarkdown(src, 80)
	applyTheme(mustTheme(t, "phosphor"))
	phosphor := renderMarkdown(src, 80)

	if pumice == phosphor {
		t.Fatal("fenced code highlighting did not change with theme")
	}
	if !strings.Contains(pumice, ansiRGB(mustTheme(t, "pumice").Accent)) {
		t.Fatalf("pumice highlight missing accent %s\n%s", mustTheme(t, "pumice").Accent, pumice)
	}
	if !strings.Contains(phosphor, ansiRGB(mustTheme(t, "phosphor").Accent)) {
		t.Fatalf("phosphor highlight missing accent %s\n%s", mustTheme(t, "phosphor").Accent, phosphor)
	}
}

func mustTheme(t *testing.T, id string) Theme {
	t.Helper()
	th, ok := themeByID(id)
	if !ok {
		t.Fatalf("unknown theme %q", id)
	}
	return th
}

func ansiRGB(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return ""
	}
	return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
}
