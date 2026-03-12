package agent

import "fmt"

// CursorAgent implements the Agent interface for the Cursor CLI (cursor-agent).
type CursorAgent struct{}

// Name returns "cursor".
func (c *CursorAgent) Name() string { return "cursor" }

// LaunchCommand returns the shell command to start the Cursor agent.
func (c *CursorAgent) LaunchCommand() string {
	return "cursor-agent"
}

// LaunchEnv returns environment variables for launching Cursor.
func (c *CursorAgent) LaunchEnv() map[string]string {
	return nil
}

// IdlePattern returns the pattern used to detect Cursor's idle prompt.
// Cursor shows "→ Add a follow-up" in a bordered input box when idle.
func (c *CursorAgent) IdlePattern() string {
	return "Add a follow-up"
}

// DialogIndicators returns strings that indicate a permission/confirmation dialog.
// Cursor auto-approves file edits but may prompt for shell commands.
func (c *CursorAgent) DialogIndicators() []string {
	return []string{
		"Trust this workspace",
		"Use arrow keys to navigate",
		"Do you want to",
	}
}

// StartupDialogs returns patterns to auto-dismiss during launch.
// Cursor shows a trust dialog on first run — press 'a' to accept.
func (c *CursorAgent) StartupDialogs() []string {
	return []string{"Trust this workspace"}
}

// StartupKey returns the key to press for the startup trust dialog.
// Cursor uses 'a' for "Trust this workspace", not Enter.
func (c *CursorAgent) StartupKey() string {
	return "a"
}

// CompletionMode returns "idle_pattern" since Cursor doesn't have JSONL events.
func (c *CursorAgent) CompletionMode() string {
	return "idle_pattern"
}

// DetectActivity returns an error since Cursor doesn't support JSONL-based detection.
func (c *CursorAgent) DetectActivity(worktreePath string) (string, string, error) {
	return "", "", fmt.Errorf("cursor agent does not support JSONL detection")
}

func init() {
	Register("cursor", &CursorAgent{})
}
