package set

// Set is a generic unordered collection of unique elements.
type Set[T comparable] map[T]struct{}

// New returns an empty Set.
func New[T comparable](elems ...T) Set[T] {
	s := make(Set[T], len(elems))
	for _, e := range elems {
		s[e] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(v T)          { s[v] = struct{}{} }
func (s Set[T]) Remove(v T)       { delete(s, v) }
func (s Set[T]) Contains(v T) bool { _, ok := s[v]; return ok }
func (s Set[T]) Len() int         { return len(s) }

// Union returns a new set containing all elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := New[T]()
	for v := range s {
		out.Add(v)
	}
	for v := range other {
		out.Add(v)
	}
	return out
}

// Intersection returns a new set containing only elements present in both sets.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	out := New[T]()
	// iterate the smaller set for efficiency
	small, big := s, other
	if len(small) > len(big) {
		small, big = big, small
	}
	for v := range small {
		if big.Contains(v) {
			out.Add(v)
		}
	}
	return out
}
