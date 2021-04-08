package handlers

import (
	"context"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Bank interface {
	CreateApplication(ctx context.Context, application commonModels.Application) (commonModels.ApplicationStatus, error)
	GetApplicationStatus(ctx context.Context, id uuid.UUID) (commonModels.ApplicationStatus, error)
}

type Notifier interface {
	ApplicationStatusChanged(ctx context.Context, change commonModels.StatusChange) error
}

type Handler struct {
	bank     Bank
	notifier Notifier
	logger   log.FieldLogger
}

func (h *Handler) UpdateJob(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error {
	const query = `
		UPDATE jobs
		SET status = $2, application = $3, updated_at = now()
		WHERE id = $1
	`
	_, err := tx.ExecContext(ctx, query, job.ID, job.Status, job.Application)
	return err
}
