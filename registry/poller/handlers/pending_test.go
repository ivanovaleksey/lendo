package handlers

import (
	"github.com/brianvoe/gofakeit"
	commonModels "github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/registry/bank"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPendingJobHandler_Handle(t *testing.T) {
	pendingJob := models.Job{
		ID: uuid.NewV4(),
		Application: commonModels.Application{
			NewApplication: commonModels.NewApplication{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
			},
			ID:     uuid.NewV4(),
			Status: commonModels.ApplicationStatusPending,
		},
		Status: models.JobStatusPending,
	}

	t.Run("when cannot get application status", func(t *testing.T) {
		fx := newPendingHandlerFixture(t)
		defer fx.Finish()

		bankErr := bank.Error{
			Code:    400,
			Message: gofakeit.Sentence(3),
		}
		fx.bank.On("GetApplicationStatus", fx.ctx, pendingJob.Application.ID).Return(commonModels.ApplicationStatus(""), bankErr)

		err := fx.handler.Handle(fx.ctx, dummyExecer{}, pendingJob)

		assert.Equal(t, bankErr, errors.Cause(err))
	})

	t.Run("when application status has not changed", func(t *testing.T) {
		fx := newPendingHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := pendingJob.Application.Status
		fx.bank.On("GetApplicationStatus", fx.ctx, pendingJob.Application.ID).Return(applicationStatus, nil)

		err := fx.handler.Handle(fx.ctx, dummyExecer{}, pendingJob)

		assert.NoError(t, err)
	})

	t.Run("when cannot update job", func(t *testing.T) {
		fx := newPendingHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("GetApplicationStatus", fx.ctx, pendingJob.Application.ID).Return(applicationStatus, nil)

		repoErr := errors.New(gofakeit.Sentence(3))
		job := pendingJob
		job.Status = models.JobStatusDone
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(repoErr)

		err := fx.handler.Handle(fx.ctx, fx.tx, pendingJob)

		assert.Equal(t, repoErr, errors.Cause(err))
	})

	t.Run("when cannot notify about status change", func(t *testing.T) {
		fx := newPendingHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("GetApplicationStatus", fx.ctx, pendingJob.Application.ID).Return(applicationStatus, nil)

		job := pendingJob
		job.Status = models.JobStatusDone
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(nil)

		notification := commonModels.StatusChange{
			ID:     pendingJob.Application.ID,
			Status: applicationStatus,
		}
		notifierErr := errors.New(gofakeit.Sentence(3))
		fx.notifier.On("ApplicationStatusChanged", fx.ctx, notification).Return(notifierErr)

		err := fx.handler.Handle(fx.ctx, fx.tx, pendingJob)

		assert.Equal(t, notifierErr, errors.Cause(err))
	})

	t.Run("when everything is fine should move to done", func(t *testing.T) {
		fx := newPendingHandlerFixture(t)
		defer fx.Finish()

		applicationStatus := commonModels.ApplicationStatus(gofakeit.Word())
		fx.bank.On("GetApplicationStatus", fx.ctx, pendingJob.Application.ID).Return(applicationStatus, nil)

		job := pendingJob
		job.Status = models.JobStatusDone
		job.Application.Status = applicationStatus
		fx.repo.On("UpdateJobTx", fx.ctx, fx.tx, job).Return(nil)

		notification := commonModels.StatusChange{
			ID:     pendingJob.Application.ID,
			Status: applicationStatus,
		}
		fx.notifier.On("ApplicationStatusChanged", fx.ctx, notification).Return(nil)

		err := fx.handler.Handle(fx.ctx, fx.tx, pendingJob)

		assert.NoError(t, err)
	})
}

func newPendingHandlerFixture(t *testing.T) *fixture{
	fx := newFixture(t)
	fx.handler = NewPendingJobHandler(fx.bank, fx.repo, fx.notifier)
	return fx
}
