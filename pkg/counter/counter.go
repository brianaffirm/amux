package counter

import "sync"

// Counter is a thread-safe integer counter.
type Counter struct {
	mu  sync.Mutex
	val int
}

// New returns a new Counter initialized to zero.
func New() *Counter {
	return &Counter{}
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

// Dec decrements the counter by 1.
func (c *Counter) Dec() {
	c.mu.Lock()
	c.val--
	c.mu.Unlock()
}

// Value returns the current counter value.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}
