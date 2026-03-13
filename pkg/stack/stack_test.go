package stack

import "testing"

func TestPushAndPop(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	s.Push(3)

	v, ok := s.Pop()
	if !ok || v != 3 {
		t.Fatalf("Pop() = (%d, %v), want (3, true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 2 {
		t.Fatalf("Pop() = (%d, %v), want (2, true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 1 {
		t.Fatalf("Pop() = (%d, %v), want (1, true)", v, ok)
	}
}

func TestPeek(t *testing.T) {
	var s Stack[string]
	s.Push("a")
	s.Push("b")

	v, ok := s.Peek()
	if !ok || v != "b" {
		t.Fatalf("Peek() = (%q, %v), want (\"b\", true)", v, ok)
	}
	if s.Len() != 2 {
		t.Fatalf("Peek should not remove elements; Len() = %d, want 2", s.Len())
	}
}

func TestLen(t *testing.T) {
	var s Stack[int]
	if s.Len() != 0 {
		t.Fatalf("new stack Len() = %d, want 0", s.Len())
	}
	s.Push(42)
	if s.Len() != 1 {
		t.Fatalf("after Push Len() = %d, want 1", s.Len())
	}
	s.Pop()
	if s.Len() != 0 {
		t.Fatalf("after Pop Len() = %d, want 0", s.Len())
	}
}

func TestIsEmpty(t *testing.T) {
	var s Stack[int]
	if !s.IsEmpty() {
		t.Fatal("new stack should be empty")
	}
	s.Push(1)
	if s.IsEmpty() {
		t.Fatal("stack with one element should not be empty")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Fatal("stack after popping last element should be empty")
	}
}

func TestPopEmpty(t *testing.T) {
	var s Stack[int]
	v, ok := s.Pop()
	if ok {
		t.Fatal("Pop on empty stack should return false")
	}
	if v != 0 {
		t.Fatalf("Pop on empty stack should return zero value, got %d", v)
	}
}

func TestPeekEmpty(t *testing.T) {
	var s Stack[string]
	v, ok := s.Peek()
	if ok {
		t.Fatal("Peek on empty stack should return false")
	}
	if v != "" {
		t.Fatalf("Peek on empty stack should return zero value, got %q", v)
	}
}

func TestStringType(t *testing.T) {
	var s Stack[string]
	s.Push("hello")
	s.Push("world")

	v, ok := s.Pop()
	if !ok || v != "world" {
		t.Fatalf("Pop() = (%q, %v), want (\"world\", true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != "hello" {
		t.Fatalf("Pop() = (%q, %v), want (\"hello\", true)", v, ok)
	}
}
