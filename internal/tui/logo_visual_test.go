package tui

import (
	"fmt"
	"strings"
	"testing"
)

// TestLogoVisual prints the narrow dashboard to stdout so you can eyeball
// the logo. Run with: go test ./internal/tui/ -run TestLogoVisual -v
func TestLogoVisual(t *testing.T) {
	for _, width := range []int{30, 40, 50} {
		m := DashboardModel{
			width: width,
			workspaces: []WorkspaceRow{
				{ID: "refactor-auth", Status: "RUNNING"},
				{ID: "fix-login-bug", Status: "READY"},
			},
			planName: "inbox-060",
		}
		out := m.renderNarrowDashboard()
		fmt.Printf("\n=== width=%d ===\n%s\n", width, out)
		fmt.Println(strings.Repeat("─", width))
	}
}
