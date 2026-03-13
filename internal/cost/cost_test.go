package cost

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brianaffirm/towr/internal/dispatch"
)

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

func TestParseClaudeTokens(t *testing.T) {
	dir := t.TempDir()
	old := dispatch.GetProjectsDirOverride()
	dispatch.SetProjectsDirOverride(dir)
	t.Cleanup(func() { dispatch.SetProjectsDirOverride(old) })

	worktreePath := "/Users/test/.towr/worktrees/towr/costtest"
	encoded := dispatch.ClaudeProjectDir(worktreePath)
	projDir := filepath.Join(dir, encoded)
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}

	t.Run("result entry with usage", func(t *testing.T) {
		jsonlFile := filepath.Join(projDir, "session.jsonl")
		content := "{\"type\":\"user\",\"timestamp\":\"2026-03-13T00:00:00Z\"}\n{\"type\":\"result\",\"result\":{\"usage\":{\"input_tokens\":12450,\"output_tokens\":38200}}}\n"
		if err := os.WriteFile(jsonlFile, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		usage, err := ParseClaudeTokens(worktreePath)
		if err != nil {
			t.Fatalf("ParseClaudeTokens: %v", err)
		}
		if usage.InputTokens != 12450 {
			t.Errorf("input = %d, want 12450", usage.InputTokens)
		}
		if usage.OutputTokens != 38200 {
			t.Errorf("output = %d, want 38200", usage.OutputTokens)
		}
		if usage.Source != "jsonl-parsed" {
			t.Errorf("source = %q, want jsonl-parsed", usage.Source)
		}
	})

	t.Run("no result entry returns unavailable", func(t *testing.T) {
		// Separate worktree path to avoid mtime race
		wt2 := "/Users/test/.towr/worktrees/towr/costtest2"
		enc2 := dispatch.ClaudeProjectDir(wt2)
		pd2 := filepath.Join(dir, enc2)
		if err := os.MkdirAll(pd2, 0o755); err != nil {
			t.Fatal(err)
		}
		jf := filepath.Join(pd2, "session.jsonl")
		content := "{\"type\":\"user\"}\n{\"type\":\"last-prompt\",\"lastPrompt\":\"hello\"}\n"
		if err := os.WriteFile(jf, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		usage, err := ParseClaudeTokens(wt2)
		if err != nil {
			t.Fatalf("ParseClaudeTokens: %v", err)
		}
		if usage.Source != "unavailable" {
			t.Errorf("source = %q, want unavailable", usage.Source)
		}
	})
}

func TestEstimateTokens(t *testing.T) {
	usage := EstimateTokens("Write a simple function that adds two numbers")
	if usage.InputTokens == 0 {
		t.Error("estimated input should be > 0")
	}
	if usage.Source != "estimated" {
		t.Errorf("source = %q, want estimated", usage.Source)
	}
}
