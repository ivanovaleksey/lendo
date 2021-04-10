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

func WithRepo(r handlers.Repo) Option {
	return func(p *Poller) {
		p.repo = r
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

func WithTickerFactory(f TickerFactory) Option {
	return func(p *Poller) {
		p.tickerFactory = f
	}
}

func WithWorkerFactory(f WorkerFactory) Option {
	return func(p *Poller) {
		p.workerFactory = f
	}
}
