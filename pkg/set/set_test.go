package set

import "testing"

func TestAddContains(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	if !s.Contains(1) {
		t.Fatal("expected set to contain 1")
	}
	if s.Contains(3) {
		t.Fatal("expected set not to contain 3")
	}
}

func TestRemove(t *testing.T) {
	s := New[string]("a", "b", "c")
	s.Remove("b")
	if s.Contains("b") {
		t.Fatal("expected b to be removed")
	}
	if s.Len() != 2 {
		t.Fatalf("expected len 2, got %d", s.Len())
	}
}

func TestLen(t *testing.T) {
	s := New[int]()
	if s.Len() != 0 {
		t.Fatalf("expected len 0, got %d", s.Len())
	}
	s.Add(10)
	s.Add(20)
	s.Add(10) // duplicate
	if s.Len() != 2 {
		t.Fatalf("expected len 2, got %d", s.Len())
	}
}

func TestNewWithElements(t *testing.T) {
	s := New(1, 2, 3, 2)
	if s.Len() != 3 {
		t.Fatalf("expected len 3, got %d", s.Len())
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Contains(v) {
			t.Fatalf("expected set to contain %d", v)
		}
	}
}

func TestUnion(t *testing.T) {
	a := New(1, 2, 3)
	b := New(3, 4, 5)
	u := a.Union(b)

	if u.Len() != 5 {
		t.Fatalf("expected union len 5, got %d", u.Len())
	}
	for _, v := range []int{1, 2, 3, 4, 5} {
		if !u.Contains(v) {
			t.Fatalf("expected union to contain %d", v)
		}
	}
}

func TestUnionEmpty(t *testing.T) {
	a := New(1, 2)
	b := New[int]()
	u := a.Union(b)
	if u.Len() != 2 {
		t.Fatalf("expected union len 2, got %d", u.Len())
	}
}

func TestIntersection(t *testing.T) {
	a := New(1, 2, 3, 4)
	b := New(3, 4, 5, 6)
	i := a.Intersection(b)

	if i.Len() != 2 {
		t.Fatalf("expected intersection len 2, got %d", i.Len())
	}
	for _, v := range []int{3, 4} {
		if !i.Contains(v) {
			t.Fatalf("expected intersection to contain %d", v)
		}
	}
}

func TestIntersectionDisjoint(t *testing.T) {
	a := New(1, 2)
	b := New(3, 4)
	i := a.Intersection(b)
	if i.Len() != 0 {
		t.Fatalf("expected intersection len 0, got %d", i.Len())
	}
}

func TestRemoveNonExistent(t *testing.T) {
	s := New(1, 2)
	s.Remove(99)
	if s.Len() != 2 {
		t.Fatalf("expected len 2 after removing non-existent, got %d", s.Len())
	}
}

func TestStringSet(t *testing.T) {
	s := New("hello", "world")
	if !s.Contains("hello") {
		t.Fatal("expected set to contain 'hello'")
	}
	s.Add("world") // duplicate
	if s.Len() != 2 {
		t.Fatalf("expected len 2, got %d", s.Len())
	}
}
