package cost

import "testing"

func TestCalculate(t *testing.T) {
	tests := []struct {
		model string
		usage TokenUsage
		want  float64
	}{
		{"opus", TokenUsage{InputTokens: 10000, OutputTokens: 30000}, 2.40},
		{"sonnet", TokenUsage{InputTokens: 10000, OutputTokens: 30000}, 0.48},
		{"haiku", TokenUsage{InputTokens: 10000, OutputTokens: 30000}, 0.04},
		{"opus", TokenUsage{InputTokens: 0, OutputTokens: 0}, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			got := Calculate(tt.model, tt.usage)
			if diff := got - tt.want; diff > 0.01 || diff < -0.01 {
				t.Errorf("Calculate(%q, %+v) = %.4f, want %.4f", tt.model, tt.usage, got, tt.want)
			}
		})
	}
}

func TestCalculate_UnknownModel(t *testing.T) {
	got := Calculate("gpt-4", TokenUsage{InputTokens: 1000, OutputTokens: 1000})
	if got != 0 {
		t.Errorf("unknown model should return 0, got %f", got)
	}
}
