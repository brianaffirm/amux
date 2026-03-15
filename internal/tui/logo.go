package tui

import (
	_ "embed"
	"strings"
)

//go:embed logo.txt
var logoRaw string

// scaleLogo centers the logo within maxW columns.
func scaleLogo(maxW int) string {
	lines := strings.Split(strings.TrimRight(logoRaw, "\n"), "\n")

	// Strip blank top/bottom margins.
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Find the minimum leading whitespace so we can strip it uniformly,
	// preserving the relative shape of the logo.
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}

	// Strip the common indent and measure the widest content line.
	stripped := make([]string, len(lines))
	maxContent := 0
	for i, line := range lines {
		if len(line) > minIndent {
			line = line[minIndent:]
		}
		stripped[i] = line
		if w := len([]rune(line)); w > maxContent {
			maxContent = w
		}
	}

	// Center the block within maxW.
	pad := (maxW - maxContent) / 2
	if pad < 0 {
		pad = 0
	}
	padding := strings.Repeat(" ", pad)
	var out []string
	for _, line := range stripped {
		out = append(out, padding+line)
	}
	return strings.Join(out, "\n")
}
