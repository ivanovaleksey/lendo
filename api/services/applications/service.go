package applicationsSrv

import (
	"context"
	applicationsRepo "github.com/ivanovaleksey/lendo/api/repos/applications"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type GetListParams = applicationsRepo.GetListParams

type Service struct {
	repo     Repo
	notifier Notifier
}

type Repo interface {
	GetList(ctx context.Context, params GetListParams) ([]models.Application, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Application, error)
	Create(ctx context.Context, item models.Application) (uuid.UUID, error)
}

type Notifier interface {
	NewApplication(ctx context.Context, application models.Application) error
}

func New(repo Repo, notifier Notifier) *Service {
	srv := &Service{
		repo:     repo,
		notifier: notifier,
	}
	return srv
}

func (srv *Service) GetList(ctx context.Context, params GetListParams) ([]models.Application, int, error) {
	return srv.repo.GetList(ctx, params)
}

func (srv *Service) GetByID(ctx context.Context, id uuid.UUID) (models.Application, error) {
	return srv.repo.GetByID(ctx, id)
}

func (srv *Service) Create(ctx context.Context, item models.NewApplication) (uuid.UUID, error) {
	application := models.Application{
		NewApplication: item,
		Status:         models.ApplicationStatusNew,
	}
	id, err := srv.repo.Create(ctx, application)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "can't create application")
	}

	application.ID = id
	log.Debugf("notify applications.new: %v", application)
	err = srv.notifier.NewApplication(ctx, application)
	if err != nil {
		log.Error(errors.Wrap(err, "can't notify"))
	}

	return id, nil
}
