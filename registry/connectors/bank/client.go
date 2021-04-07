package bank

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
)

type Client interface {
	CreateApplication(ctx context.Context, application models.Application) (models.ApplicationStatus, error)
	GetApplicationStatus(ctx context.Context, id uuid.UUID) (models.ApplicationStatus, error)
}
