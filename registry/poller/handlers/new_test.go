package handlers

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/bank"
	"github.com/ivanovaleksey/lendo/registry/models"
	mockHandlers "github.com/ivanovaleksey/lendo/registry/poller/handlers/mocks"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJobHandler_Handle(t *testing.T) {
	newJob := models.Job{
		ID: uuid.NewV4(),
		Application: commonModels.Application{
			NewApplication: commonModels.NewApplication{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
			},
			ID:     uuid.NewV4(),
			Status: commonModels.ApplicationStatusNew,
		},
		Status: models.JobStatusNew,
	}

	t.Run("when cannot register application in bank", func(t *testing.T) {
		fx := newNewHandlerFixture(t)
		defer fx.Finish()

		bankErr := bank.Error{
			Code:    400,
			Message: gofakeit.Sentence(3),
		}
		fx.bank.On("CreateApplication", fx.ctx, newJob.Application).Return(commonModels.ApplicationStatus(""), bankErr)

		err := fx.handler.Handle(fx.ctx, dummyExecer{}, newJob)

		assert.Equal(t, bankErr, errors.Cause(err))
	})

	t.Run("when cannot update job", func(t *testing.T) {
		fx := newNewHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("CreateApplication", fx.ctx, newJob.Application).Return(applicationStatus, nil)

		repoErr := errors.New(gofakeit.Sentence(3))
		job := newJob
		job.Status = models.JobStatusPending
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(repoErr)

		err := fx.handler.Handle(fx.ctx, fx.tx, newJob)

		assert.Equal(t, repoErr, errors.Cause(err))
	})

	t.Run("when cannot notify about status change should not fail", func(t *testing.T) {
		fx := newNewHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("CreateApplication", fx.ctx, newJob.Application).Return(applicationStatus, nil)

		job := newJob
		job.Status = models.JobStatusPending
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(nil)

		notification := commonModels.StatusChange{
			ID:     newJob.Application.ID,
			Status: applicationStatus,
		}
		notifierErr := errors.New(gofakeit.Sentence(3))
		fx.notifier.On("ApplicationStatusChanged", fx.ctx, notification).Return(notifierErr)

		err := fx.handler.Handle(fx.ctx, fx.tx, newJob)

		assert.NoError(t, err)
	})

	t.Run("when everything is fine should move to pending", func(t *testing.T) {
		fx := newNewHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("CreateApplication", fx.ctx, newJob.Application).Return(applicationStatus, nil)

		job := newJob
		job.Status = models.JobStatusPending
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(nil)

		notification := commonModels.StatusChange{
			ID:     newJob.Application.ID,
			Status: applicationStatus,
		}
		fx.notifier.On("ApplicationStatusChanged", fx.ctx, notification).Return(nil)

		err := fx.handler.Handle(fx.ctx, fx.tx, newJob)

		assert.NoError(t, err)
	})
}

type fixture struct {
	t   *testing.T
	ctx context.Context
	tx  sqlx.ExecerContext

	bank     *mockHandlers.Bank
	repo     *mockHandlers.Repo
	notifier *mockHandlers.Notifier

	handler handler
}

type handler interface {
	Handle(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error
}

func newFixture(t *testing.T) *fixture {
	fx := &fixture{
		t:   t,
		ctx: context.Background(),
		tx:  dummyExecer{},

		bank:     &mockHandlers.Bank{},
		repo:     &mockHandlers.Repo{},
		notifier: &mockHandlers.Notifier{},
	}
	return fx
}

func newNewHandlerFixture(t *testing.T) *fixture {
	fx := newFixture(t)
	fx.handler = NewNewJobHandler(fx.bank, fx.repo, fx.notifier)
	return fx
}

func (fx *fixture) Finish() {
	fx.bank.AssertExpectations(fx.t)
	fx.repo.AssertExpectations(fx.t)
	fx.notifier.AssertExpectations(fx.t)
}

type dummyExecer struct{}

func (dummyExecer) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, nil
}
