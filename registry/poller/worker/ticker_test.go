package worker

import (
	"testing"
	"time"
)

// fixedTicker ticks cap(ch) times
type fixedTicker struct {
	ch chan time.Time
}

func newFixedTicker(ticks int) *fixedTicker {
	ch := make(chan time.Time, ticks)
	for i := 0; i < cap(ch); i++ {
		ch <- time.Now()
	}
	return &fixedTicker{ch: ch}
}

func (h *fixedTicker) Tick() <-chan time.Time {
	return h.ch
}

func (h *fixedTicker) Stop() {
	close(h.ch)
}

func TestFixedTicker(t *testing.T) {
	const ticks = 3

	tick := newFixedTicker(ticks)
	defer tick.Stop()

	sink := make(chan time.Time, 3)
	defer close(sink)

	for i := 0; i < ticks; i++ {
		v := <-tick.Tick()
		sink <- v
	}
}

