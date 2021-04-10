package poller

import (
	"github.com/ivanovaleksey/lendo/pkg/ticker"
	"time"
)

type TickerFactory interface {
	NewTicker() ticker.Ticker
}

type stdTickerFactory struct {
	duration time.Duration
}

func (p stdTickerFactory) NewTicker() ticker.Ticker {
	return ticker.NewTicker(p.duration)
}
