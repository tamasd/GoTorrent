package util

import (
	"sync"
	"testing"
)

func TestCounterConcurrency(t *testing.T) {
	c := NewCounter()

	var wg sync.WaitGroup

	num := 128

	wg.Add(num * 2)

	for i := 0; i < num; i++ {
		go func() {
			c.Add(1)
			wg.Done()
		}()
		go func() {
			c.Value()
			wg.Done()
		}()
	}

	wg.Wait()
	val := c.Value()
	if val != uint64(num) {
		t.Errorf("got %d expected %d", val, num)
	}
}
