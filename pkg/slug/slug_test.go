package slug

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "lowercases basic text",
			in:   "Hello World",
			want: "hello-world",
		},
		{
			name: "removes punctuation",
			in:   "Hello, World!",
			want: "hello-world",
		},
		{
			name: "collapses spaces and hyphens",
			in:   "Already---Slug   Here",
			want: "already-slug-here",
		},
		{
			name: "keeps digits",
			in:   "Version 2 Release 10",
			want: "version-2-release-10",
		},
		{
			name: "trims leading and trailing separators",
			in:   "  -- Hello World --  ",
			want: "hello-world",
		},
		{
			name: "drops unsupported characters",
			in:   "Go@#$%^&*()+=Lang",
			want: "golang",
		},
		{
			name: "empty when nothing remains",
			in:   "!!!",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.in)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
