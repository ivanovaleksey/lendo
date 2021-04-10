package worker

import (
	"context"
	"database/sql"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/ticker"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	id        int
	txFactory db.TxFactory
	logger    log.FieldLogger
	ticker    ticker.Ticker
	handlers  map[models.JobStatus]Handler
}

type Handler interface {
	Handle(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error
}

func New(opts ...Option) *Worker {
	w := &Worker{
		handlers: make(map[models.JobStatus]Handler),
	}
	for _, opt := range opts {
		opt(w)
	}
	w.logger = log.WithFields(log.Fields{
		"component": "worker",
		"id":        w.id,
	})
	return w
}

func (w *Worker) Run(ctx context.Context) error {
	defer w.ticker.Stop()

	for {
		select {
		case <-w.ticker.Tick():
			if err := w.doWork(ctx); err != nil {
				w.logger.Error(err)
			}
		case <-ctx.Done():
			w.logger.Debug("context cancelled")
			return ctx.Err()
		}
	}
}

func (w *Worker) doWork(ctx context.Context) error {
	w.logger.Debug("do work")

	tx, err := w.txFactory.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "can't begin tx")
	}

	return tx.Do(ctx, w.doWorkTx)
}

func (w *Worker) doWorkTx(ctx context.Context, tx db.SQLTx) error {
	const query = `
		SELECT id, application, status
		FROM jobs
		WHERE status IN ('new', 'pending')
		ORDER BY created_at
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var job models.Job
	err := tx.QueryRowxContext(ctx, query).StructScan(&job)
	switch {
	case err == sql.ErrNoRows:
		w.logger.Debug("no work")
		return nil
	case err != nil:
		return errors.Wrap(err, "can't get job")
	}

	handler, ok := w.handlers[job.Status]
	if !ok {
		w.logger.Errorf("no handler for status %q", job.Status)
		return nil
	}

	return handler.Handle(ctx, tx, job)
}
