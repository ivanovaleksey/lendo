package handlers

import (
	"context"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// PendingJobHandler polls a bank system for the application status,
// updates job and application status, notifies queue.
type PendingJobHandler struct {
	Handler
}

func NewPendingJobHandler(bank Bank, repo Repo, notifier Notifier) *PendingJobHandler {
	return &PendingJobHandler{
		Handler: Handler{
			bank:     bank,
			repo:     repo,
			notifier: notifier,
			logger:   log.WithField("handler", "pending"),
		},
	}
}

func (h *PendingJobHandler) Handle(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error {
	logger := h.logger.WithField("job_id", job.ID.String())

	status, err := h.bank.GetApplicationStatus(ctx, job.Application.ID)
	if err != nil {
		return errors.Wrap(err, "can't get application status")
	}

	if status == job.Application.Status {
		logger.Debug("not ready yet")
		return nil
	}

	job.Status = models.JobStatusDone
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
		return errors.Wrap(err, "can't send notification")
	}

	return nil
}
