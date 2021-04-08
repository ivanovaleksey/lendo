package applicationsPubSub

import (
	"context"
	"encoding/json"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ApplicationStatusChangedHandler struct {
	repo   Repo
	logger log.FieldLogger
}

type Repo interface {
	UpdateStatus(ctx context.Context, change commonModels.StatusChange) error
}

func NewApplicationStatusChangedHandler(repo Repo) *ApplicationStatusChangedHandler {
	h := &ApplicationStatusChangedHandler{
		repo:   repo,
		logger: log.WithField("handler", "applications-changed"),
	}
	return h
}

func (h *ApplicationStatusChangedHandler) Handle(msg *nats.Msg) {
	h.logger.Debugf("status changed %s", string(msg.Data))

	var change commonModels.StatusChange
	err := json.Unmarshal(msg.Data, &change)
	if err != nil {
		err = errors.Wrap(err, "can't parse message")
		h.logger.Error(err)
		return
	}

	// todo: think about ctx
	err = h.repo.UpdateStatus(context.TODO(), change)
	if err != nil {
		err = errors.Wrap(err, "can't update status")
		h.logger.Error(err)
		return
	}

	h.logger.Debugf("status changed %s", change.ID.String())
}
