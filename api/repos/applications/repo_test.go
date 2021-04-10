package applicationsRepo

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/api/config"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/test"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestImpl_GetList(t *testing.T) {
	t.Skip("not implemented")
}

func TestImpl_GetByID(t *testing.T) {
	t.Run("when record does not exist", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		id := uuid.NewV4()

		application, err := fx.repo.GetByID(fx.ctx, id)

		require.Equal(t, ErrNotFound, err)
		assert.Empty(t, application)
	})

	t.Run("when record exists", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		item := models.Application{
			NewApplication: models.NewApplication{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
			},
			Status: models.ApplicationStatusPending,
		}
		item.ID = fx.createApplication(item)

		application, err := fx.repo.GetByID(fx.ctx, item.ID)

		require.NoError(t, err)
		assert.Equal(t, item, application)
	})
}

func TestImpl_Create(t *testing.T) {
	fx := newFixture(t)
	defer fx.Finish()

	item := models.Application{
		NewApplication: models.NewApplication{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		},
		Status: randomStatus(),
	}

	id, err := fx.repo.Create(fx.ctx, item)

	require.NoError(t, err)
	item.ID = id
	application := fx.getApplication(id)
	assert.Equal(t, item, application)
}

type fixture struct {
	t   *testing.T
	ctx context.Context
	db  *db.DB

	repo *Repo
}

func newFixture(t *testing.T) *fixture {
	test.LoadAPIEnv(t)

	cfg, err := config.New()
	require.NoError(t, err)

	fx := &fixture{
		t:   t,
		ctx: context.Background(),
		db:  db.NewTestDB(t, cfg.DB),
	}
	fx.repo = New(fx.db)
	return fx
}

func (fx *fixture) Finish() {
	require.NoError(fx.t, fx.db.Close())
}

func (fx *fixture) createApplication(item models.Application) (id uuid.UUID) {
	const q = `
		INSERT INTO applications(first_name, last_name, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := fx.db.GetContext(fx.ctx, &id, q, item.FirstName, item.LastName, item.Status)
	require.NoError(fx.t, err)
	return
}

func (fx *fixture) getApplication(id uuid.UUID) (item models.Application) {
	const q = `SELECT row_to_json(applications) FROM applications WHERE id = $1`
	err := fx.db.GetContext(fx.ctx, &item, q, id)
	require.NoError(fx.t, err)
	return
}

func randomStatus() models.ApplicationStatus {
	all := []models.ApplicationStatus{
		models.ApplicationStatusNew,
		models.ApplicationStatusPending,
		models.ApplicationStatusCompleted,
		models.ApplicationStatusRejected,
	}
	return all[rand.Intn(len(all))]
}
