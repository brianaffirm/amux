// Package mux implements a terminal multiplexer with dynamic focus tiling.
package mux

// Role identifies the purpose of a pane rectangle in the layout.
type Role string

const (
	RoleFocus   Role = "focus"
	RoleSidebar Role = "sidebar"
	RoleControl Role = "control"
)

// Rect describes the position, size, and role of a pane in the terminal grid.
type Rect struct {
	X, Y, W, H int
	Role        Role
	PaneIndex   int // -1 for control pane
}

// Breakpoints for responsive layout.
const (
	wideMinCols   = 160
	mediumMinCols = 100
)

// Control pane dimensions.
const (
	controlRows       = 10 // default height within 8-14 range
	narrowControlRows = 2  // bottom bar in narrow mode
)

// ComputeLayout returns rectangles for each pane given terminal dimensions,
// pane count, focused pane index, and whether the control pane is visible.
//
// Layout rules:
//   - Wide (>=160 cols): focus 60% left, sidebar 40% right
//   - Medium (100-159 cols): focus 65% left, sidebar 35% right
//   - Narrow (<100 cols): full-screen focus, no sidebar, control becomes 2-line bottom bar
//   - Sidebar panes divide height equally (minus control pane if shown)
//   - Control pane: pinned bottom-right, 8-14 rows, collapsible
func ComputeLayout(termW, termH, paneCount, focusIdx int, showControl bool) []Rect {
	if paneCount <= 0 {
		return nil
	}

	// Clamp focusIdx.
	if focusIdx < 0 || focusIdx >= paneCount {
		focusIdx = 0
	}

	if termW < mediumMinCols {
		return layoutNarrow(termW, termH, paneCount, focusIdx, showControl)
	}
	return layoutWideOrMedium(termW, termH, paneCount, focusIdx, showControl)
}

// layoutNarrow handles <100 cols: full-screen focus, optional 2-line control bar.
func layoutNarrow(termW, termH, paneCount, focusIdx int, showControl bool) []Rect {
	var rects []Rect

	focusH := termH
	if showControl {
		focusH = termH - narrowControlRows
	}

	rects = append(rects, Rect{
		X: 0, Y: 0, W: termW, H: focusH,
		Role:      RoleFocus,
		PaneIndex: focusIdx,
	})

	if showControl {
		rects = append(rects, Rect{
			X: 0, Y: focusH, W: termW, H: narrowControlRows,
			Role:      RoleControl,
			PaneIndex: -1,
		})
	}

	return rects
}

// layoutWideOrMedium handles >=100 cols with focus left and sidebar right.
func layoutWideOrMedium(termW, termH, paneCount, focusIdx int, showControl bool) []Rect {
	var rects []Rect

	// Determine focus/sidebar split ratio.
	var focusPct int
	if termW >= wideMinCols {
		focusPct = 60
	} else {
		focusPct = 65
	}

	sidebarCount := paneCount - 1

	// Single pane: no sidebar needed.
	if sidebarCount == 0 {
		focusW := termW
		focusH := termH

		if showControl {
			// Control pane bottom-right; focus takes the rest.
			ctrlH := controlRows
			rects = append(rects, Rect{
				X: 0, Y: 0, W: focusW, H: termH - ctrlH,
				Role:      RoleFocus,
				PaneIndex: focusIdx,
			})
			rects = append(rects, Rect{
				X: 0, Y: termH - ctrlH, W: focusW, H: ctrlH,
				Role:      RoleControl,
				PaneIndex: -1,
			})
		} else {
			rects = append(rects, Rect{
				X: 0, Y: 0, W: focusW, H: focusH,
				Role:      RoleFocus,
				PaneIndex: focusIdx,
			})
		}
		return rects
	}

	focusW := termW * focusPct / 100
	sidebarW := termW - focusW

	// Focus pane: full height on the left.
	rects = append(rects, Rect{
		X: 0, Y: 0, W: focusW, H: termH,
		Role:      RoleFocus,
		PaneIndex: focusIdx,
	})

	// Right column: sidebar panes + optional control pane.
	sidebarAvailH := termH
	if showControl {
		sidebarAvailH = termH - controlRows
	}

	// Distribute sidebar panes equally in available height.
	sidebarPaneH := sidebarAvailH / sidebarCount
	curY := 0

	sidebarIdx := 0
	for i := 0; i < paneCount; i++ {
		if i == focusIdx {
			continue
		}
		h := sidebarPaneH
		// Last sidebar pane absorbs rounding remainder.
		if sidebarIdx == sidebarCount-1 {
			h = sidebarAvailH - curY
		}
		rects = append(rects, Rect{
			X: focusW, Y: curY, W: sidebarW, H: h,
			Role:      RoleSidebar,
			PaneIndex: i,
		})
		curY += h
		sidebarIdx++
	}

	if showControl {
		rects = append(rects, Rect{
			X: focusW, Y: sidebarAvailH, W: sidebarW, H: controlRows,
			Role:      RoleControl,
			PaneIndex: -1,
		})
	}

	return rects
}
