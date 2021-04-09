package ticker

import "time"

type Ticker interface {
	Tick() <-chan time.Time
	Stop()
}

type ticker struct {
	*time.Ticker
}

func NewTicker(dur time.Duration) Ticker {
	t := &ticker{
		Ticker: time.NewTicker(dur),
	}
	return t
}

func (t *ticker) Tick() <-chan time.Time {
	return t.C
}
