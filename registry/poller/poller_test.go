package poller

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/pkg/db"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/ticker"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/ivanovaleksey/lendo/registry/poller/handlers/mocks"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPoller_Run(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("should move new to pending", func(t *testing.T) {
		const numTicks = 1

		fx := newFixture(t, numTicks)
		defer fx.Finish()

		jobs := fx.buildJobs()
		newJob := jobs[0]
		fx.insertJob(newJob)

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())

		{
			fx.bank.On("CreateApplication", mock.Anything, newJob.Application).Return(applicationStatus, nil)

			notification := commonModels.StatusChange{
				ID:     newJob.Application.ID,
				Status: applicationStatus,
			}
			fx.notifier.On("ApplicationStatusChanged", mock.Anything, notification).Return(nil)
		}

		err := fx.poller.Run(fx.ctx)
		require.NoError(t, err)

		fx.Wait()

		job1 := fx.getJob(newJob.ID)
		assert.Equal(t, models.JobStatusPending, job1.Status)
		newJob.Application.Status = applicationStatus
		assert.Equal(t, newJob.Application, job1.Application)
	})

	t.Run("should move pending to done", func(t *testing.T) {
		const numTicks = 1

		fx := newFixture(t, numTicks)
		defer fx.Finish()

		jobs := fx.buildJobs()
		pendingJob := jobs[1]
		fx.insertJob(pendingJob)

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())

		{
			fx.bank.On("GetApplicationStatus", mock.Anything, pendingJob.Application.ID).Return(applicationStatus, nil)

			notification := commonModels.StatusChange{
				ID:     pendingJob.Application.ID,
				Status: applicationStatus,
			}
			fx.notifier.On("ApplicationStatusChanged", mock.Anything, notification).Return(nil)
		}

		err := fx.poller.Run(fx.ctx)
		require.NoError(t, err)

		fx.Wait()

		job := fx.getJob(pendingJob.ID)
		assert.Equal(t, models.JobStatusDone, job.Status)
		pendingJob.Application.Status = applicationStatus
		assert.Equal(t, pendingJob.Application, job.Application)
	})

	t.Run("should move new to pending to done", func(t *testing.T) {
		const numTicks = 2

		fx := newFixture(t, numTicks)
		defer fx.Finish()

		jobs := fx.buildJobs()
		newJob := jobs[0]
		fx.insertJob(newJob)

		status1 := commonModels.ApplicationStatus(gofakeit.Word())
		status2 := commonModels.ApplicationStatus(gofakeit.Word())

		{
			fx.bank.On("CreateApplication", mock.Anything, newJob.Application).Return(status1, nil)

			notification := commonModels.StatusChange{
				ID:     newJob.Application.ID,
				Status: status1,
			}
			fx.notifier.On("ApplicationStatusChanged", mock.Anything, notification).Return(nil)
		}
		{
			fx.bank.On("GetApplicationStatus", mock.Anything, newJob.Application.ID).Return(status2, nil)

			notification := commonModels.StatusChange{
				ID:     newJob.Application.ID,
				Status: status2,
			}
			fx.notifier.On("ApplicationStatusChanged", mock.Anything, notification).Return(nil)
		}

		err := fx.poller.Run(fx.ctx)
		require.NoError(t, err)

		fx.Wait()

		job := fx.getJob(newJob.ID)
		assert.Equal(t, models.JobStatusDone, job.Status)
		newJob.Application.Status = status2
		assert.Equal(t, newJob.Application, job.Application)
	})

	t.Run("when application status has not changed", func(t *testing.T) {
		const numTicks = 1

		fx := newFixture(t, numTicks)
		defer fx.Finish()

		jobs := fx.buildJobs()
		pendingJob := jobs[1]
		fx.insertJob(pendingJob)

		applicationStatus := pendingJob.Application.Status

		{
			fx.bank.On("GetApplicationStatus", mock.Anything, pendingJob.Application.ID).Return(applicationStatus, nil)
		}

		err := fx.poller.Run(fx.ctx)
		require.NoError(t, err)

		fx.Wait()

		job := fx.getJob(pendingJob.ID)
		assert.Equal(t, models.JobStatusPending, job.Status)
		pendingJob.Application.Status = applicationStatus
		assert.Equal(t, pendingJob.Application, job.Application)
	})
}

type fixture struct {
	t   *testing.T
	ctx context.Context
	db  *db.DB

	bank     *mocks.Bank
	notifier *mocks.Notifier

	poller *Poller
}

func newFixture(t *testing.T, numTicks int) *fixture {
	cfg, err := config.New()
	require.NoError(t, err)

	fx := &fixture{
		t:   t,
		ctx: context.Background(),
		db:  db.NewTestDB(t, cfg.DB),

		bank:     &mocks.Bank{},
		notifier: &mocks.Notifier{},
	}

	opts := []Option{
		WithDB(fx.db),
		WithBank(fx.bank),
		WithNotifier(fx.notifier),
		WithNumWorkers(1),
		WithTickerProvider(FixedTickerProvider{count: numTicks}),
	}
	fx.poller = New(opts...)

	return fx
}

func (fx *fixture) Finish() {
	require.NoError(fx.t, fx.poller.Close())
	fx.bank.AssertExpectations(fx.t)
	fx.notifier.AssertExpectations(fx.t)
}

func (fx *fixture) Wait() {
	time.AfterFunc(1*time.Second, func() {
		fx.poller.Close()
	})
	fx.poller.workersWg.Wait()
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
		log.Debugf("job %v", job)
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

// FixedTicker ticks cap(ch) times
type FixedTicker struct {
	ch chan time.Time
}

func (h *FixedTicker) Tick() <-chan time.Time {
	return h.ch
}

func (h *FixedTicker) Stop() {
	close(h.ch)
}

type FixedTickerProvider struct {
	count int
}

func (p FixedTickerProvider) NewTicker() ticker.Ticker {
	ch := make(chan time.Time, p.count)
	for i := 0; i < p.count; i++ {
		ch <- time.Now()
	}
	return &FixedTicker{ch: ch}
}

func TestFixedTickerProvider(t *testing.T) {
	const size = 3

	p := FixedTickerProvider{count: size}
	tick := p.NewTicker()
	defer tick.Stop()

	sink := make(chan time.Time, 3)
	defer close(sink)

	for i := 0; i < size; i++ {
		v := <-tick.Tick()
		sink <- v
	}
}
