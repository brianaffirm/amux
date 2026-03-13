package math2

// Clamp returns val clamped to [min, max].
func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Lerp linearly interpolates between a and b by t.
// t=0 returns a, t=1 returns b.
func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
