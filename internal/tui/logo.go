package tui

import (
	_ "embed"
	"strings"
)

//go:embed logo.txt
var logoRaw string

// scaleLogo returns a scaled-down version of the logo that fits within maxW
// columns. It strips surrounding whitespace and samples every other row to
// reduce height by half.
func scaleLogo(maxW int) string {
	lines := strings.Split(strings.TrimRight(logoRaw, "\n"), "\n")

	// Find the first and last non-empty lines so we skip blank margins.
	first, last := 0, len(lines)-1
	for first < len(lines) && strings.TrimSpace(lines[first]) == "" {
		first++
	}
	for last > first && strings.TrimSpace(lines[last]) == "" {
		last--
	}
	lines = lines[first : last+1]

	var out []string
	for i, line := range lines {
		// Sample every other line to halve the height.
		if i%2 != 0 {
			continue
		}
		// Strip leading spaces so the logo left-aligns in the pane.
		stripped := strings.TrimLeft(line, " ")
		runes := []rune(stripped)
		if len(runes) > maxW {
			stripped = string(runes[:maxW])
		}
		out = append(out, stripped)
	}
	return strings.Join(out, "\n")
}
