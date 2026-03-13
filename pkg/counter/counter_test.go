package counter

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
}

func TestIncDec(t *testing.T) {
	c := New()
	c.Inc()
	c.Inc()
	c.Inc()
	if c.Value() != 3 {
		t.Fatalf("expected 3, got %d", c.Value())
	}
	c.Dec()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestConcurrentInc(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != n {
		t.Fatalf("expected %d, got %d", n, c.Value())
	}
}

func TestConcurrentIncDec(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	n := 500
	wg.Add(n * 2)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
		go func() {
			defer wg.Done()
			c.Dec()
		}()
	}
	wg.Wait()
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
}
