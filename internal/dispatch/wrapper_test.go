package dispatch

import (
	"strings"
	"testing"
)

func TestBuildWrapper(t *testing.T) {
	script := BuildWrapper("ws-1", "dispatch-abc", "/tmp/comms/ws-1")

	checks := []struct {
		name    string
		contain string
	}{
		{"shebang", "#!/usr/bin/env bash"},
		{"workspace id", `WS_ID="ws-1"`},
		{"dispatch id", `DISPATCH_ID="dispatch-abc"`},
		{"comms dir", `COMMS_DIR="/tmp/comms/ws-1"`},
		{"prompt file", `PROMPT_FILE="$COMMS_DIR/prompt.md"`},
		{"result file", `RESULT_FILE="$COMMS_DIR/result.json"`},
		{"heartbeat touch", "touch \"$COMMS_DIR/heartbeat\""},
		{"heartbeat sleep", "sleep 30"},
		{"trap", "trap"},
		{"report started", "towr report --dispatch-id \"$DISPATCH_ID\" --status started"},
		{"claude call", "claude -p"},
		{"output format", "--output-format json"},
		{"report success", "--status success"},
		{"report failed", "--status failed"},
	}

	for _, c := range checks {
		if !strings.Contains(script, c.contain) {
			t.Errorf("%s: wrapper missing %q", c.name, c.contain)
		}
	}
}

func TestBuildRunCommand(t *testing.T) {
	cmd := BuildRunCommand("/tmp/comms/ws-1")
	want := "bash /tmp/comms/ws-1/run.sh"
	if cmd != want {
		t.Errorf("BuildRunCommand = %q, want %q", cmd, want)
	}
}
