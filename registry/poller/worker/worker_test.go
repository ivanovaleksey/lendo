package worker

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/pkg/db"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/test"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/ivanovaleksey/lendo/registry/poller/worker/mocks"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWorker_Run(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("should handle new job", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		jobs := fx.buildJobs()
		newJob := jobs[0]
		fx.insertJobs(jobs)

		fx.newJobHandler.On("Handle", fx.ctx, mock.AnythingOfType("db.tx"), newJob).Return(nil)

		time.AfterFunc(100*time.Millisecond, fx.cancel)
		err := fx.worker.Run(fx.ctx)

		require.Equal(t, context.Canceled, err)
	})

	t.Run("should handle pending job", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		jobs := fx.buildJobs()
		pendingJob := jobs[1]
		fx.insertJobs(jobs[1:])

		fx.pendingJobHandler.On("Handle", fx.ctx, mock.AnythingOfType("db.tx"), pendingJob).Return(nil)

		time.AfterFunc(100*time.Millisecond, fx.cancel)
		err := fx.worker.Run(fx.ctx)

		require.Equal(t, context.Canceled, err)
	})
}

type fixture struct {
	t      *testing.T
	ctx    context.Context
	cancel context.CancelFunc
	db     *db.DB

	newJobHandler     *mocks.Handler
	pendingJobHandler *mocks.Handler

	worker *Worker
}

func newFixture(t *testing.T, opts ...Option) *fixture {
	test.LoadRegistryEnv(t)

	cfg, err := config.New()
	require.NoError(t, err)

	fx := &fixture{
		t:  t,
		db: db.NewTestDB(t, cfg.DB),

		newJobHandler:     &mocks.Handler{},
		pendingJobHandler: &mocks.Handler{},
	}
	fx.ctx, fx.cancel = context.WithCancel(context.Background())

	baseOpts := []Option{
		WithTicker(newFixedTicker(1)),
		WithTxFactory(db.NewTxFactory(fx.db)),
		WithHandler(models.JobStatusNew, fx.newJobHandler),
		WithHandler(models.JobStatusPending, fx.pendingJobHandler),
	}
	opts = append(baseOpts, opts...)
	fx.worker = New(opts...)
	return fx
}

func (fx *fixture) Finish() {
	fx.cancel()
	fx.newJobHandler.AssertExpectations(fx.t)
	fx.pendingJobHandler.AssertExpectations(fx.t)
}

func (fx *fixture) buildJobs() []models.Job {
	jobs := []models.Job{
		{
			Application: commonModels.Application{
				NewApplication: commonModels.NewApplication{
					FirstName: gofakeit.FirstName(),
					LastName:  gofakeit.LastName(),
				},
				ID:     uuid.NewV4(),
				Status: commonModels.ApplicationStatus(gofakeit.Word()),
			},
			Status: models.JobStatusNew,
			ID:     uuid.NewV4(),
		},
		{
			Application: commonModels.Application{
				NewApplication: commonModels.NewApplication{
					FirstName: gofakeit.FirstName(),
					LastName:  gofakeit.LastName(),
				},
				ID:     uuid.NewV4(),
				Status: commonModels.ApplicationStatus(gofakeit.Word()),
			},
			Status: models.JobStatusPending,
			ID:     uuid.NewV4(),
		},
		{
			Application: commonModels.Application{
				NewApplication: commonModels.NewApplication{
					FirstName: gofakeit.FirstName(),
					LastName:  gofakeit.LastName(),
				},
				ID:     uuid.NewV4(),
				Status: commonModels.ApplicationStatus(gofakeit.Word()),
			},
			Status: models.JobStatusDone,
			ID:     uuid.NewV4(),
		},
	}
	return jobs
}

func (fx *fixture) insertJobs(jobs []models.Job) {
	for _, job := range jobs {
		fx.insertJob(job)
	}
}

func (fx *fixture) insertJob(job models.Job) {
	const q = `
		INSERT INTO jobs (id, application, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	_, err := fx.db.ExecContext(fx.ctx, q, job.ID, job.Application, job.Status)
	require.NoError(fx.t, err)
}

func (fx *fixture) getJob(id uuid.UUID) (job models.Job) {
	const q = `SELECT id, application, status FROM jobs WHERE id = $1`
	err := fx.db.GetContext(fx.ctx, &job, q, id)
	require.NoError(fx.t, err)
	return
}
