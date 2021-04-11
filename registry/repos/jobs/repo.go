package jobsRepo

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

const (
	tableName = "jobs"
)

type Repo struct {
	db      *db.DB
	builder squirrel.StatementBuilderType
}

func New(database *db.DB) *Repo {
	repo := &Repo{
		db:      database,
		builder: db.Builder,
	}
	return repo
}

func (repo *Repo) CreateJob(ctx context.Context, job models.Job) (uuid.UUID, error) {
	const query = `
		INSERT INTO ` + tableName + ` (application, status)
		VALUES ($1, $2)
		RETURNING id
	`

	var id uuid.UUID
	err := repo.db.GetContext(ctx, &id, query, job.Application, job.Status)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (repo *Repo) UpdateJob(ctx context.Context, job models.Job) error {
	return repo.UpdateJobTx(ctx, repo.db, job)
}

func (repo *Repo) UpdateJobTx(ctx context.Context, tx sqlx.ExecerContext, job models.Job) error {
	const query = `
		UPDATE ` + tableName + `
		SET status = $2, application = $3, updated_at = now()
		WHERE id = $1
	`
	_, err := tx.ExecContext(ctx, query, job.ID, job.Status, job.Application)
	return err
}
