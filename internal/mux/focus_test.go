package mux

import "testing"

func TestFocusResizeCommands(t *testing.T) {
	// 3 panes: control (0), master (1), agent (2)
	// Focus on master (1) — should get 60% width
	cmds := BuildFocusCommands("towr-mux", "mux", 1, 3, 200)

	// Should select pane 1.
	hasSelect := false
	for _, c := range cmds {
		if c.Args[0] == "select-pane" {
			hasSelect = true
		}
	}
	if !hasSelect {
		t.Error("should have select-pane command")
	}

	// Should resize focused pane.
	hasResize := false
	for _, c := range cmds {
		if c.Args[0] == "resize-pane" {
			hasResize = true
		}
	}
	if !hasResize {
		t.Error("should have resize-pane command")
	}
}

func TestFocusOnControlPaneKeepsLayout(t *testing.T) {
	// Focus on control pane (0) — control stays 20%, others share 80%
	cmds := BuildFocusCommands("towr-mux", "mux", 0, 3, 200)

	// Should still select pane 0.
	hasSelect := false
	for _, c := range cmds {
		if c.Args[0] == "select-pane" {
			for _, arg := range c.Args {
				if arg == "towr-mux:mux.0" {
					hasSelect = true
				}
			}
		}
	}
	if !hasSelect {
		t.Error("should select control pane")
	}
}

func TestFocusTwoPanesNoResize(t *testing.T) {
	// Only 2 panes (control + master) — no resizing needed.
	cmds := BuildFocusCommands("towr-mux", "mux", 1, 2, 200)

	hasResize := false
	for _, c := range cmds {
		if c.Args[0] == "resize-pane" {
			hasResize = true
		}
	}
	if hasResize {
		t.Error("should not resize with only 2 panes")
	}

	// Should still have select-pane.
	if len(cmds) != 1 {
		t.Errorf("expected 1 command (select-pane), got %d", len(cmds))
	}
}

func TestFocusAgentPaneWidths(t *testing.T) {
	// 3 panes, focus on pane 1, 200 cols.
	cmds := BuildFocusCommands("towr-mux", "mux", 1, 3, 200)

	// Should have: select-pane, resize focused to 60%, resize control to 20%.
	if len(cmds) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(cmds))
	}

	// Resize focused pane to 120 (60% of 200).
	resizeFocus := cmds[1]
	if resizeFocus.Args[0] != "resize-pane" {
		t.Errorf("expected resize-pane, got %s", resizeFocus.Args[0])
	}
	foundWidth := false
	for i, arg := range resizeFocus.Args {
		if arg == "-x" && i+1 < len(resizeFocus.Args) && resizeFocus.Args[i+1] == "120" {
			foundWidth = true
		}
	}
	if !foundWidth {
		t.Errorf("focused pane should be resized to 120 (60%% of 200), got args: %v", resizeFocus.Args)
	}

	// Resize control to 40 (20% of 200).
	resizeControl := cmds[2]
	foundControlWidth := false
	for i, arg := range resizeControl.Args {
		if arg == "-x" && i+1 < len(resizeControl.Args) && resizeControl.Args[i+1] == "40" {
			foundControlWidth = true
		}
	}
	if !foundControlWidth {
		t.Errorf("control pane should be resized to 40 (20%% of 200), got args: %v", resizeControl.Args)
	}
}

func TestFocusMinControlWidth(t *testing.T) {
	// Small terminal: 100 cols. 20% = 20, but min is 30.
	cmds := BuildFocusCommands("towr-mux", "mux", 0, 3, 100)

	// Control pane resize should use minimum width of 30.
	foundMinWidth := false
	for _, c := range cmds {
		if c.Args[0] == "resize-pane" {
			for i, arg := range c.Args {
				if arg == "-x" && i+1 < len(c.Args) && c.Args[i+1] == "30" {
					foundMinWidth = true
				}
			}
		}
	}
	if !foundMinWidth {
		t.Error("control pane should use minimum width of 30 when 20% < 30")
	}
}
