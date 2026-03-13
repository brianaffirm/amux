package wc

import (
	"reflect"
	"testing"
)

func TestCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "counts case insensitive words",
			text: "Go go GO stop",
			want: map[string]int{
				"go":   3,
				"stop": 1,
			},
		},
		{
			name: "splits on punctuation and whitespace",
			text: "Hello, world!\nHello\tworld... hello?",
			want: map[string]int{
				"hello": 3,
				"world": 2,
			},
		},
		{
			name: "supports unicode letters and numbers",
			text: "Cafe cafe CAFE 2026 café",
			want: map[string]int{
				"2026": 1,
				"cafe": 3,
				"café": 1,
			},
		},
		{
			name: "returns empty map for no words",
			text: " \n\t!!!",
			want: map[string]int{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Count(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Count(%q) = %#v, want %#v", tt.text, got, tt.want)
			}
		})
	}
}
