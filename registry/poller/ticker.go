package poller

import (
	"github.com/ivanovaleksey/lendo/pkg/ticker"
	"time"
)

type TickerProvider interface {
	NewTicker() ticker.Ticker
}

type stdTickerProvider struct {
	duration time.Duration
}

func (p stdTickerProvider) NewTicker() ticker.Ticker {
	return ticker.NewTicker(p.duration)
}
