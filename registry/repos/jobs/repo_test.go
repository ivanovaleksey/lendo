package jobsRepo

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/pkg/db"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/test"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepo_CreateJob(t *testing.T) {
	fx := newFixture(t)
	defer fx.Finish()

	application := commonModels.Application{
		NewApplication: commonModels.NewApplication{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		},
		ID:     uuid.NewV4(),
		Status: commonModels.ApplicationStatus(gofakeit.Word()),
	}
	item := models.Job{
		Application: application,
		Status:      models.JobStatus(gofakeit.Word()),
	}

	id, err := fx.repo.CreateJob(fx.ctx, item)

	require.NoError(t, err)
	item.ID = id
	job := fx.getJob(id)
	assert.Equal(t, item, job)
}

type fixture struct {
	t   *testing.T
	ctx context.Context
	db  *db.DB

	repo *Repo
}

func newFixture(t *testing.T) *fixture {
	test.LoadRegistryEnv(t)

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

func (fx *fixture) getJob(id uuid.UUID) (job models.Job) {
	const query = `SELECT id, application, status FROM jobs WHERE id = $1`

	err := fx.db.GetContext(fx.ctx, &job, query, id)
	require.NoError(fx.t, err)
	return
}
