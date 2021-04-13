package applicationsPubSub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/api/pubsub/applications/mocks"
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPub_NewApplication(t *testing.T) {
	application := models.Application{
		ID:     uuid.NewV4(),
		Status: models.ApplicationStatus(gofakeit.Word()),
		NewApplication: models.NewApplication{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		},
	}
	data, _ := json.Marshal(application)

	t.Run("without error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		fx.nats.On("Publish", "applications.new", data).Return(nil)

		err := fx.pub.NewApplication(fx.ctx, application)

		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		clientErr := errors.New(gofakeit.Sentence(3))
		fx.nats.On("Publish", "applications.new", data).Return(clientErr)

		err := fx.pub.NewApplication(fx.ctx, application)

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
