package rev

import "testing"

func TestReverseString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "single char", in: "a", want: "a"},
		{name: "ascii", in: "hello", want: "olleh"},
		{name: "palindrome", in: "racecar", want: "racecar"},
		{name: "spaces", in: "a b c", want: "c b a"},
		{name: "unicode CJK", in: "日本語", want: "語本日"},
		{name: "accented", in: "café", want: "éfac"},
		{name: "emoji", in: "Go🚀Fast", want: "tsaF🚀oG"},
		{name: "multi emoji", in: "😀🎉🔥", want: "🔥🎉😀"},
		{name: "mixed unicode", in: "Hello, 世界! 🌍", want: "🌍 !界世 ,olleH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseString(tt.in)
			if got != tt.want {
				t.Errorf("ReverseString(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
