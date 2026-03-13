package agent

import "sort"

// registry maps agent names to their implementations.
var registry = map[string]Agent{
	"claude-code": &ClaudeCode{},
}

// Get returns the agent by name. Returns Default() if name is empty or not found.
func Get(name string) Agent {
	if name == "" {
		return Default()
	}
	if a, ok := registry[name]; ok {
		return a
	}
	return Default()
}

// GetWithModel returns an agent for the given model shorthand.
// Maps: "opus" → ClaudeCode{opus}, "sonnet" → ClaudeCode{sonnet},
// "cursor" → CursorAgent, "codex" → CodexAgent.
// Falls back to Get(name) for unknown models.
func GetWithModel(model, agentName string) Agent {
	// Model takes precedence if specified.
	if model != "" {
		switch model {
		case "opus", "sonnet", "haiku":
			return &ClaudeCode{ModelFlag: model}
		case "cursor":
			return Get("cursor")
		case "codex":
			return Get("codex")
		default:
			// Assume it's a full model ID for claude-code.
			return &ClaudeCode{ModelFlag: model}
		}
	}
	return Get(agentName)
}

// Register adds an agent to the registry.
func Register(name string, a Agent) {
	registry[name] = a
}

// List returns all registered agent names in sorted order.
func List() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
