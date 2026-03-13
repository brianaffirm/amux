package cost

const Version = "2026-03"

type ModelPricing struct {
	InputPerMillion  float64
	OutputPerMillion float64
}

// Pricing reflects Anthropic's current model pricing as of March 2026.
// Source: https://docs.anthropic.com/en/docs/about-claude/models
//   Opus 4.6:  $5/M input, $25/M output
//   Sonnet 4.6: $3/M input, $15/M output
//   Haiku 4.5:  $1/M input, $5/M output
var Pricing = map[string]ModelPricing{
	"opus":   {5.00, 25.00},
	"sonnet": {3.00, 15.00},
	"haiku":  {1.00, 5.00},
}

func Calculate(model string, usage TokenUsage) float64 {
	p, ok := Pricing[model]
	if !ok {
		return 0
	}
	inputCost := float64(usage.InputTokens) / 1_000_000 * p.InputPerMillion
	outputCost := float64(usage.OutputTokens) / 1_000_000 * p.OutputPerMillion
	return inputCost + outputCost
}
