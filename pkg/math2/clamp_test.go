package math2

import "testing"

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		val, min, max float64
		want     float64
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -3, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 10},
		{"min equals max", 5, 7, 7, 7},
		{"negative range within", -5, -10, -1, -5},
		{"negative range below", -15, -10, -1, -10},
		{"negative range above", 0, -10, -1, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Clamp(tc.val, tc.min, tc.max)
			if got != tc.want {
				t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tc.val, tc.min, tc.max, got, tc.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		name    string
		a, b, t float64
		want    float64
	}{
		{"t=0 returns a", 0, 10, 0, 0},
		{"t=1 returns b", 0, 10, 1, 10},
		{"t=0.5 midpoint", 0, 10, 0.5, 5},
		{"t=0.25", 0, 100, 0.25, 25},
		{"negative values", -10, 10, 0.5, 0},
		{"a equals b", 5, 5, 0.7, 5},
		{"extrapolate t>1", 0, 10, 2, 20},
		{"extrapolate t<0", 0, 10, -1, -10},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Lerp(tc.a, tc.b, tc.t)
			if got != tc.want {
				t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tc.a, tc.b, tc.t, got, tc.want)
			}
		})
	}
}
