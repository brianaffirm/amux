package mux

import (
	"strings"
	"testing"
)

// testPane is a mock PaneView for testing without real PTYs.
type testPane struct {
	id       string
	input    []byte
	exited   bool
	exitCode int
}

func (p *testPane) Render() string                { return "[" + p.id + " output]" }
func (p *testPane) Write(b []byte) (int, error)   { p.input = append(p.input, b...); return len(b), nil }
func (p *testPane) Resize(c, r int) error         { return nil }
func (p *testPane) Exited() bool                  { return p.exited }
func (p *testPane) ExitCode() int                 { return p.exitCode }
func (p *testPane) ID() string                    { return p.id }
func (p *testPane) Close() error                  { return nil }
func (p *testPane) CursorPosition() (int, int)    { return 0, 0 }
func (p *testPane) Notify() <-chan struct{}        { return make(chan struct{}) }

var _ PaneView = (*testPane)(nil)

func TestRenderControlContent_Basic(t *testing.T) {
	panes := []PaneView{
		&testPane{id: "agent-1", exited: false, exitCode: 0},
		&testPane{id: "agent-2", exited: false, exitCode: 0},
		&testPane{id: "agent-3", exited: true, exitCode: 0},
	}

	result := RenderControlContent(panes, 40, 20)
	lines := strings.Split(result, "\n")

	if len(lines) != 20 {
		t.Errorf("expected 20 lines, got %d", len(lines))
	}

	// Should contain running count of 2
	found := false
	for _, line := range lines {
		if strings.Contains(line, "running") && strings.Contains(line, "2") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find 'running' with count 2 in output:\n%s", result)
	}

	// Should contain completed count of 1
	found = false
	for _, line := range lines {
		if strings.Contains(line, "completed") && strings.Contains(line, "1") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find 'completed' with count 1 in output:\n%s", result)
	}
}

func TestRenderControlContent_MinHeight(t *testing.T) {
	panes := []PaneView{
		&testPane{id: "agent-1", exited: false, exitCode: 0},
	}

	result := RenderControlContent(panes, 40, 3)
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// Should still render something meaningful
	nonEmpty := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmpty++
		}
	}
	if nonEmpty == 0 {
		t.Errorf("expected at least one non-empty line with h=3")
	}
}

func TestRenderControlContent_WithErrors(t *testing.T) {
	panes := []PaneView{
		&testPane{id: "agent-1", exited: true, exitCode: 1},
		&testPane{id: "agent-2", exited: true, exitCode: 0},
	}

	result := RenderControlContent(panes, 40, 20)

	if !strings.Contains(result, "errored") {
		t.Errorf("expected 'errored' in output when there are error exits:\n%s", result)
	}
}

func TestRenderControlContent_NoErrors(t *testing.T) {
	panes := []PaneView{
		&testPane{id: "agent-1", exited: true, exitCode: 0},
	}

	result := RenderControlContent(panes, 40, 20)

	if strings.Contains(result, "errored") {
		t.Errorf("should not show 'errored' when count is 0:\n%s", result)
	}
}

func TestRenderControlContent_LongID(t *testing.T) {
	panes := []PaneView{
		&testPane{id: "very-long-agent-identifier", exited: false, exitCode: 0},
	}

	result := RenderControlContent(panes, 40, 20)

	// ID should be truncated to 12 chars
	if strings.Contains(result, "very-long-agent-identifier") {
		t.Errorf("long ID should be truncated:\n%s", result)
	}
	if !strings.Contains(result, "very-long-ag") {
		t.Errorf("expected truncated ID 'very-long-ag' in output:\n%s", result)
	}
}
