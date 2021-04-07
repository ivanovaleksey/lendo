package applicationsRepo

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/models"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Repo interface {
	GetList(ctx context.Context, params GetListParams) ([]models.Application, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Application, error)
	Create(ctx context.Context, item models.NewApplication) (uuid.UUID, error)
}
