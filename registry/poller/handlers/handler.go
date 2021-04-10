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

type Repo interface {
	UpdateJobTx(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error
}

type Handler struct {
	repo     Repo
	bank     Bank
	notifier Notifier
	logger   log.FieldLogger
}
