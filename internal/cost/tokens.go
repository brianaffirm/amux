package cost

type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	Source       string // "jsonl-parsed", "estimated", "unavailable"
}
