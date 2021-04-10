package worker

import (
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/ticker"
	"github.com/ivanovaleksey/lendo/registry/models"
)

type Option func(*Worker)

func WithID(id int) Option {
	return func(w *Worker) {
		w.id = id
	}
}

func WithTxFactory(f db.TxFactory) Option {
	return func(w *Worker) {
		w.txFactory = f
	}
}

func WithTicker(t ticker.Ticker) Option {
	return func(w *Worker) {
		w.ticker = t
	}
}

func WithHandler(s models.JobStatus, h Handler) Option {
	return func(w *Worker) {
		w.handlers[s] = h
	}
}
