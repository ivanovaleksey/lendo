package applicationsPubSub

import (
	"context"
	"encoding/json"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type NewApplicationHandler struct {
	repo   Repo
	logger log.FieldLogger
}

type Repo interface {
	CreateJob(ctx context.Context, job models.Job) (uuid.UUID, error)
}

func NewNewApplicationHandler(repo Repo) *NewApplicationHandler {
	h := &NewApplicationHandler{
		repo:   repo,
		logger: log.WithField("handler", "applications-new"),
	}
	return h
}

func (h *NewApplicationHandler) Handle(msg *nats.Msg) {
	h.logger.Debugf("new application %s", string(msg.Data))

	var application commonModels.Application
	err := json.Unmarshal(msg.Data, &application)
	if err != nil {
		err = errors.Wrap(err, "can't parse application")
		h.logger.Error(err)
		return
	}

	job := models.Job{
		Status:      models.JobStatusNew,
		Application: application,
	}
	// todo: think about ctx
	id, err := h.repo.CreateJob(context.TODO(), job)
	if err != nil {
		err = errors.Wrap(err, "can't create job")
		h.logger.Error(err)
		return
	}

	h.logger.Debugf("job created %s", id.String())
}
