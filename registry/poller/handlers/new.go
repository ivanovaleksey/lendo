package handlers

import (
	"context"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// NewJobHandler registers a new application in a bank system
// and moves the job to a 'pending' status.
type NewJobHandler struct {
	Handler
}

func NewNewJobHandler(bank Bank, repo Repo, notifier Notifier) *NewJobHandler {
	return &NewJobHandler{
		Handler: Handler{
			bank:     bank,
			repo:     repo,
			notifier: notifier,
			logger:   log.WithField("handler", "new"),
		},
	}
}

func (h *NewJobHandler) Handle(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error {
	logger := h.logger.WithField("job_id", job.ID.String())

	status, err := h.bank.CreateApplication(ctx, job.Application)
	if err != nil {
		return errors.Wrap(err, "can't create application in bank")
	}

	job.Status = models.JobStatusPending
	job.Application.Status = status
	err = h.repo.UpdateJobTx(ctx, tx, job)
	if err != nil {
		return errors.Wrap(err, "can't update application status")
	}

	notification := commonModels.StatusChange{
		ID:     job.Application.ID,
		Status: job.Application.Status,
	}
	err = h.notifier.ApplicationStatusChanged(ctx, notification)
	if err != nil {
		logger.Errorf("can't send notification: %v", err)
	}

	return nil
}
