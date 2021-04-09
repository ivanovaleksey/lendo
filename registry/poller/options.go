package poller

import (
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/registry/poller/handlers"
)

type Option func(*Poller)

func WithBank(b handlers.Bank) Option {
	return func(p *Poller) {
		p.bank = b
	}
}

func WithNotifier(n handlers.Notifier) Option {
	return func(p *Poller) {
		p.notifier = n
	}
}

func WithDB(db *db.DB) Option {
	return func(p *Poller) {
		p.db = db
	}
}

func WithNumWorkers(num int) Option {
	return func(p *Poller) {
		p.numWorkers = num
	}
}

func WithTickerProvider(provider TickerProvider) Option {
	return func(p *Poller) {
		p.tickerProvider = provider
	}
}
