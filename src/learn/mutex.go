package learn

import (
	"sync"
)

type SafeCounter struct {
	V   map[string]int
	mux sync.Mutex
}

func (c *SafeCounter) Inc(key string)  {
	c.mux.Lock()
	c.V[key]++
	c.mux.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.V[key]
}

