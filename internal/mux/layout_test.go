package mux

import (
	"testing"
)

func TestComputeLayout_SinglePane(t *testing.T) {
	rects := ComputeLayout(200, 50, 1, 0, false)

	if len(rects) != 1 {
		t.Fatalf("expected 1 rect, got %d", len(rects))
	}
	r := rects[0]
	if r.Role != RoleFocus {
		t.Errorf("expected RoleFocus, got %q", r.Role)
	}
	if r.PaneIndex != 0 {
		t.Errorf("expected PaneIndex 0, got %d", r.PaneIndex)
	}
	// Single pane should fill the whole terminal
	if r.X != 0 || r.Y != 0 || r.W != 200 || r.H != 50 {
		t.Errorf("expected full screen (0,0,200,50), got (%d,%d,%d,%d)", r.X, r.Y, r.W, r.H)
	}
}

func TestComputeLayout_SinglePaneWithControl(t *testing.T) {
	rects := ComputeLayout(200, 50, 1, 0, true)

	// Should have focus pane + control pane
	if len(rects) != 2 {
		t.Fatalf("expected 2 rects, got %d", len(rects))
	}

	focus := findRole(rects, RoleFocus)
	ctrl := findRole(rects, RoleControl)
	if focus == nil {
		t.Fatal("no focus rect found")
	}
	if ctrl == nil {
		t.Fatal("no control rect found")
	}
	if ctrl.PaneIndex != -1 {
		t.Errorf("control pane should have PaneIndex -1, got %d", ctrl.PaneIndex)
	}
}

func TestComputeLayout_TwoPanes_Wide(t *testing.T) {
	// Wide: >=160 cols → focus 60%, sidebar 40%
	rects := ComputeLayout(200, 50, 2, 0, false)

	focus := findRole(rects, RoleFocus)
	sidebar := findByRole(rects, RoleSidebar)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	if len(sidebar) != 1 {
		t.Fatalf("expected 1 sidebar pane, got %d", len(sidebar))
	}

	// Focus should be 60% of 200 = 120
	expectedFocusW := 120
	if focus.W != expectedFocusW {
		t.Errorf("wide focus width: expected %d, got %d", expectedFocusW, focus.W)
	}
	if focus.X != 0 {
		t.Errorf("focus should start at X=0, got %d", focus.X)
	}

	// Sidebar takes remaining width
	expectedSidebarW := 200 - expectedFocusW
	if sidebar[0].W != expectedSidebarW {
		t.Errorf("wide sidebar width: expected %d, got %d", expectedSidebarW, sidebar[0].W)
	}
	if sidebar[0].X != expectedFocusW {
		t.Errorf("sidebar should start at X=%d, got %d", expectedFocusW, sidebar[0].X)
	}
}

func TestComputeLayout_TwoPanes_Medium(t *testing.T) {
	// Medium: 100-159 cols → focus 65%, sidebar 35%
	rects := ComputeLayout(140, 50, 2, 0, false)

	focus := findRole(rects, RoleFocus)
	sidebar := findByRole(rects, RoleSidebar)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	if len(sidebar) != 1 {
		t.Fatalf("expected 1 sidebar, got %d", len(sidebar))
	}

	expectedFocusW := 91 // 65% of 140 = 91
	if focus.W != expectedFocusW {
		t.Errorf("medium focus width: expected %d, got %d", expectedFocusW, focus.W)
	}
	expectedSidebarW := 140 - expectedFocusW
	if sidebar[0].W != expectedSidebarW {
		t.Errorf("medium sidebar width: expected %d, got %d", expectedSidebarW, sidebar[0].W)
	}
}

func TestComputeLayout_TwoPanes_Narrow(t *testing.T) {
	// Narrow: <100 cols → no sidebar, full-screen focus
	rects := ComputeLayout(80, 40, 2, 0, false)

	focus := findRole(rects, RoleFocus)
	sidebar := findByRole(rects, RoleSidebar)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	// Narrow: no sidebar visible
	if len(sidebar) != 0 {
		t.Errorf("narrow mode should have no sidebar panes, got %d", len(sidebar))
	}
	// Focus fills full screen
	if focus.W != 80 || focus.H != 40 {
		t.Errorf("narrow focus should be full screen (80,40), got (%d,%d)", focus.W, focus.H)
	}
}

func TestComputeLayout_NarrowWithControl(t *testing.T) {
	// Narrow with control: control becomes 2-line bottom bar
	rects := ComputeLayout(80, 40, 2, 0, true)

	focus := findRole(rects, RoleFocus)
	ctrl := findRole(rects, RoleControl)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	if ctrl == nil {
		t.Fatal("no control rect in narrow mode")
	}

	// Control pane: 2-line bottom bar spanning full width
	if ctrl.H != 2 {
		t.Errorf("narrow control height should be 2, got %d", ctrl.H)
	}
	if ctrl.W != 80 {
		t.Errorf("narrow control width should be full width (80), got %d", ctrl.W)
	}
	if ctrl.Y != 38 {
		t.Errorf("narrow control Y should be 38 (40-2), got %d", ctrl.Y)
	}

	// Focus takes remaining height
	if focus.H != 38 {
		t.Errorf("narrow focus height should be 38 (40-2), got %d", focus.H)
	}
}

func TestComputeLayout_ThreePanes_Wide(t *testing.T) {
	// 3 panes in wide: focus left, 2 sidebar panes stacked on right
	rects := ComputeLayout(200, 50, 3, 0, false)

	focus := findRole(rects, RoleFocus)
	sidebar := findByRole(rects, RoleSidebar)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	if len(sidebar) != 2 {
		t.Fatalf("expected 2 sidebar panes, got %d", len(sidebar))
	}

	// Sidebar panes divide height equally
	expectedH := 50 / 2
	for i, s := range sidebar {
		if s.H != expectedH {
			t.Errorf("sidebar[%d] height: expected %d, got %d", i, expectedH, s.H)
		}
	}
	// Second sidebar pane starts below first
	if sidebar[1].Y != sidebar[0].Y+sidebar[0].H {
		t.Errorf("sidebar panes should be vertically stacked")
	}
}

func TestComputeLayout_ThreePanes_WideWithControl(t *testing.T) {
	// 3 panes + control: sidebar panes share right column with control at bottom
	rects := ComputeLayout(200, 50, 3, 0, true)

	focus := findRole(rects, RoleFocus)
	sidebar := findByRole(rects, RoleSidebar)
	ctrl := findRole(rects, RoleControl)

	if focus == nil {
		t.Fatal("no focus rect")
	}
	if ctrl == nil {
		t.Fatal("no control rect")
	}

	// Control pane should be bottom-right, 8-14 rows
	if ctrl.H < 8 || ctrl.H > 14 {
		t.Errorf("control height should be 8-14, got %d", ctrl.H)
	}

	// Sidebar panes share remaining height above control
	sidebarTotalH := 0
	for _, s := range sidebar {
		sidebarTotalH += s.H
	}
	expectedSidebarTotal := 50 - ctrl.H
	if sidebarTotalH != expectedSidebarTotal {
		t.Errorf("sidebar total height: expected %d, got %d", expectedSidebarTotal, sidebarTotalH)
	}
}

func TestComputeLayout_FocusIndex(t *testing.T) {
	// When focus is on pane 1 (second pane), it should be the focus pane
	rects := ComputeLayout(200, 50, 3, 1, false)

	focus := findRole(rects, RoleFocus)
	if focus == nil {
		t.Fatal("no focus rect")
	}
	if focus.PaneIndex != 1 {
		t.Errorf("focus pane should be index 1, got %d", focus.PaneIndex)
	}

	// Other panes should be sidebar
	sidebar := findByRole(rects, RoleSidebar)
	if len(sidebar) != 2 {
		t.Fatalf("expected 2 sidebar panes, got %d", len(sidebar))
	}
	indices := map[int]bool{}
	for _, s := range sidebar {
		indices[s.PaneIndex] = true
	}
	if !indices[0] || !indices[2] {
		t.Errorf("sidebar should contain panes 0 and 2, got %v", indices)
	}
}

func TestComputeLayout_BreakpointBoundaries(t *testing.T) {
	// Exactly 160 cols → wide
	rects := ComputeLayout(160, 50, 2, 0, false)
	focus := findRole(rects, RoleFocus)
	if focus.W != 96 { // 60% of 160
		t.Errorf("at 160 cols (wide boundary), focus width should be 96, got %d", focus.W)
	}

	// 159 cols → medium
	rects = ComputeLayout(159, 50, 2, 0, false)
	focus = findRole(rects, RoleFocus)
	if focus.W != 103 { // 65% of 159
		t.Errorf("at 159 cols (medium), focus width should be 103, got %d", focus.W)
	}

	// Exactly 100 cols → medium
	rects = ComputeLayout(100, 50, 2, 0, false)
	focus = findRole(rects, RoleFocus)
	if focus.W != 65 { // 65% of 100
		t.Errorf("at 100 cols (medium boundary), focus width should be 65, got %d", focus.W)
	}

	// 99 cols → narrow
	rects = ComputeLayout(99, 40, 2, 0, false)
	focus = findRole(rects, RoleFocus)
	if focus.W != 99 { // full width in narrow
		t.Errorf("at 99 cols (narrow), focus width should be 99, got %d", focus.W)
	}
}

// helpers

func findRole(rects []Rect, role Role) *Rect {
	for i := range rects {
		if rects[i].Role == role {
			return &rects[i]
		}
	}
	return nil
}

func findByRole(rects []Rect, role Role) []Rect {
	var out []Rect
	for _, r := range rects {
		if r.Role == role {
			out = append(out, r)
		}
	}
	return out
}
