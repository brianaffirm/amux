package mux

import (
	"fmt"
	"strings"
)

// PaneView abstracts over Pane for testability.
type PaneView interface {
	Render() string
	Write([]byte) (int, error)
	Resize(cols, rows int) error
	Exited() bool
	ExitCode() int
	ID() string
	Close() error
	CursorPosition() (row, col int)
	Notify() <-chan struct{}
}

// RenderControlContent renders the control pane body (no header).
// The caller is responsible for adding the header above this content.
func RenderControlContent(panes []PaneView, w, h int) string {
	var lines []string

	// Count states.
	var running, completed, errored int
	for _, p := range panes {
		if !p.Exited() {
			running++
		} else if p.ExitCode() == 0 {
			completed++
		} else {
			errored++
		}
	}

	// SESSIONS section
	lines = append(lines, dimStyle.Render("SESSIONS"))
	lines = append(lines, "  "+controlLabelStyle.Render("running")+controlValueStyle.Render(fmt.Sprintf("%d", running)))
	lines = append(lines, "  "+controlLabelStyle.Render("completed")+controlValueStyle.Render(fmt.Sprintf("%d", completed)))
	if errored > 0 {
		lines = append(lines, "  "+controlLabelStyle.Render("errored")+controlValueStyle.Render(fmt.Sprintf("%d", errored)))
	}

	// PANES section (if height allows)
	// Need at least: current lines + 1 blank + 1 header + 1 pane entry
	if len(lines)+3 <= h && len(panes) > 0 {
		lines = append(lines, "") // blank separator
		lines = append(lines, dimStyle.Render("PANES"))
		for _, p := range panes {
			if len(lines) >= h {
				break
			}
			id := p.ID()
			if len(id) > 12 {
				id = id[:12]
			}

			var dot, status string
			if !p.Exited() {
				dot = dotRunning
				status = "running"
			} else if p.ExitCode() == 0 {
				dot = dotReady
				status = "done"
			} else {
				dot = dotError
				status = fmt.Sprintf("exit %d", p.ExitCode())
			}
			lines = append(lines, fmt.Sprintf("  %s %s %s", dot, id, dimStyle.Render(status)))
		}
	}

	// Truncate to h lines.
	if len(lines) > h {
		lines = lines[:h]
	}

	// Pad to exactly h lines.
	for len(lines) < h {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
