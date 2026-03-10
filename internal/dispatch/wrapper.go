package dispatch

import "fmt"

// BuildWrapper generates a shell wrapper script that orchestrates a dispatch task.
// It sets up environment variables, starts a heartbeat, runs claude, and reports results.
func BuildWrapper(workspaceID, dispatchID, commsDir string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail

# Dispatch metadata
WS_ID=%q
DISPATCH_ID=%q
COMMS_DIR=%q
PROMPT_FILE="$COMMS_DIR/prompt.md"
RESULT_FILE="$COMMS_DIR/result.json"

# Heartbeat: touch file every 30s so the orchestrator knows we're alive
(
  while true; do
    touch "$COMMS_DIR/heartbeat"
    sleep 30
  done
) &
HEARTBEAT_PID=$!
trap 'kill $HEARTBEAT_PID 2>/dev/null || true' EXIT

# Report started
towr report "$WS_ID" --dispatch-id "$DISPATCH_ID" --status started

# Run claude
if claude -p "$(cat "$PROMPT_FILE")" --output-format json > "$RESULT_FILE" 2>"$RESULT_FILE.err"; then
  towr report "$WS_ID" --dispatch-id "$DISPATCH_ID" --status success --file "$RESULT_FILE"
else
  towr report "$WS_ID" --dispatch-id "$DISPATCH_ID" --status failed --file "$RESULT_FILE.err"
fi
`, workspaceID, dispatchID, commsDir)
}

// BuildRunCommand returns the shell command to execute the wrapper script.
func BuildRunCommand(commsDir string) string {
	return fmt.Sprintf("bash %s/run.sh", commsDir)
}
