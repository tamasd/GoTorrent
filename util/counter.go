package util

import "sync"

type Counter struct {
	value uint64
	mtx   sync.RWMutex
}

func NewCounter() *Counter {
	c := new(Counter)
	return c
}

func (c *Counter) Reset() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.value = 0
}

func (c *Counter) Add(i uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.value += i
}

func (c *Counter) Value() uint64 {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.value
}
