package poller

import (
	"context"
	"database/sql"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type worker struct {
	db       *db.DB
	logger   log.FieldLogger
	handlers map[models.JobStatus]JobHandler
}

type JobHandler interface {
	Handle(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error
}

func (w *worker) Run(ctx context.Context) error {
	const tickTime = 10 * time.Second

	ticker := time.NewTicker(tickTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.doWork(ctx); err != nil {
				w.logger.Error(err)
			}
		case <-ctx.Done():
			w.logger.Debug("context cancelled")
			return ctx.Err()
		}
	}
}

func (w *worker) doWork(ctx context.Context) error {
	w.logger.Debug("do work")

	tx, err := w.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "can't begin tx")
	}

	if err := w.doWorkTx(ctx, tx); err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			w.logger.Error(errors.Wrap(rErr, "can't rollback tx"))
		}
		return err
	}

	return tx.Commit()
}

func (w *worker) doWorkTx(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
		SELECT id, application, status
		FROM jobs
		WHERE status IN ('new', 'pending')
		ORDER BY created_at
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var job models.Job
	err := tx.GetContext(ctx, &job, query)
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
