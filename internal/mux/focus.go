package mux

import "fmt"

// BuildFocusCommands returns tmux commands to focus a pane and resize the layout.
// paneIdx is the tmux pane index within the window.
// totalPanes is the total number of panes.
// termWidth is the terminal width in columns.
//
// Layout: pane 0 = control (20%), pane 1+ = agent/shell panes.
// The focused agent pane gets 60% width, others share the remaining ~20%.
func BuildFocusCommands(session, window string, paneIdx, totalPanes, termWidth int) []TmuxCmd {
	target := fmt.Sprintf("%s:%s.%d", session, window, paneIdx)

	var cmds []TmuxCmd

	// Select the pane.
	cmds = append(cmds, TmuxCmd{Args: []string{
		"select-pane", "-t", target,
	}})

	if totalPanes <= 2 {
		// Only control + 1 pane — no resizing needed.
		return cmds
	}

	// Calculate widths.
	controlW := termWidth * 20 / 100
	if controlW < 30 {
		controlW = 30
	}

	if paneIdx == 0 {
		// Focusing control pane — keep standard layout.
		// Control 20%, rest split evenly.
		cmds = append(cmds, TmuxCmd{Args: []string{
			"resize-pane", "-t", fmt.Sprintf("%s:%s.0", session, window), "-x", fmt.Sprintf("%d", controlW),
		}})
		return cmds
	}

	// Focusing an agent pane — give it 60%.
	focusW := termWidth * 60 / 100
	cmds = append(cmds, TmuxCmd{Args: []string{
		"resize-pane", "-t", target, "-x", fmt.Sprintf("%d", focusW),
	}})

	// Resize control pane back to 20%.
	cmds = append(cmds, TmuxCmd{Args: []string{
		"resize-pane", "-t", fmt.Sprintf("%s:%s.0", session, window), "-x", fmt.Sprintf("%d", controlW),
	}})

	return cmds
}
