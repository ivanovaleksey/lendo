package applicationsPubSub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/pubsub/applications/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPub_ApplicationStatusChanged(t *testing.T) {
	change := models.StatusChange{
		ID:     uuid.NewV4(),
		Status: models.ApplicationStatus(gofakeit.Word()),
	}
	data, _ := json.Marshal(change)

	t.Run("without error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		fx.nats.On("Publish", "applications.changed", data).Return(nil)

		err := fx.pub.ApplicationStatusChanged(fx.ctx, change)

		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		clientErr := errors.New(gofakeit.Sentence(3))
		fx.nats.On("Publish", "applications.changed", data).Return(clientErr)

		err := fx.pub.ApplicationStatusChanged(fx.ctx, change)

		assert.Equal(t, clientErr, err)
	})
}

type fixture struct {
	t    *testing.T
	ctx  context.Context
	nats *mocks.PubClient

	pub *Pub
}

func newFixture(t *testing.T) *fixture {
	fx := &fixture{
		t:    t,
		ctx:  context.Background(),
		nats: &mocks.PubClient{},
	}
	fx.pub = NewPub(fx.nats)
	return fx
}

func (fx *fixture) Finish() {
	fx.nats.AssertExpectations(fx.t)
}
