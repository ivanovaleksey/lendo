package applicationsRepo

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/api/config"
	apiModels "github.com/ivanovaleksey/lendo/api/models"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/test"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestImpl_GetList(t *testing.T) {
	t.Run("should return list and total", func(t *testing.T) {
		const num = 5

		fx := newFixture(t)
		defer fx.Finish()

		var applications []models.Application
		for i := 0; i < num; i++ {
			applications = append(applications, fx.createApplication())
		}

		params := GetListParams{}
		list, total, err := fx.repo.GetList(fx.ctx, params)

		require.NoError(t, err)
		assert.Equal(t, applications, list)
		assert.Equal(t, num, total)
	})

	t.Run("should filter by status", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		applications := []models.Application{
			fx.createApplicationWithStatus(models.ApplicationStatusNew),
			fx.createApplicationWithStatus(models.ApplicationStatusNew),
			fx.createApplicationWithStatus(models.ApplicationStatusPending),
			fx.createApplicationWithStatus(models.ApplicationStatusCompleted),
		}

		params := GetListParams{
			Status: models.ApplicationStatusPending,
		}
		list, total, err := fx.repo.GetList(fx.ctx, params)

		require.NoError(t, err)
		assert.Equal(t, []models.Application{applications[2]}, list)
		assert.Equal(t, 1, total)
	})

	t.Run("should paginate", func(t *testing.T) {
		const num = 5

		fx := newFixture(t)
		defer fx.Finish()

		var applications []models.Application
		for i := 0; i < num; i++ {
			applications = append(applications, fx.createApplication())
			time.Sleep(100 * time.Millisecond)
		}

		params := GetListParams{
			PaginationParams: apiModels.PaginationParams{
				Offset: 0,
				Limit:  2,
			},
		}
		list, total, err := fx.repo.GetList(fx.ctx, params)

		require.NoError(t, err)
		assert.Equal(t, applications[:2], list)
		assert.Equal(t, num, total)

		params.Offset += 2
		list, total, err = fx.repo.GetList(fx.ctx, params)

		require.NoError(t, err)
		assert.Equal(t, applications[2:4], list)
		assert.Equal(t, num, total)

		params.Offset += 2
		list, total, err = fx.repo.GetList(fx.ctx, params)

		require.NoError(t, err)
		assert.Equal(t, applications[4:], list)
		assert.Equal(t, num, total)
	})
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
		item.ID = fx.insertApplication(item)

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

func (fx *fixture) createApplication() models.Application {
	return fx.createApplicationWithStatus(randomStatus())
}

func (fx *fixture) createApplicationWithStatus(status models.ApplicationStatus) models.Application {
	item := models.Application{
		NewApplication: models.NewApplication{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		},
		Status: status,
	}
	item.ID = fx.insertApplication(item)
	return item
}

func (fx *fixture) insertApplication(item models.Application) (id uuid.UUID) {
	const q = `
		INSERT INTO applications(first_name, last_name, status, created_at)
		VALUES ($1, $2, $3, clock_timestamp())
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
