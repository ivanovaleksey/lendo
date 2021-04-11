package applicationsRepo

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	tableName = "applications"
)

var (
	ErrNotFound = errors.New("application not found")
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

type GetListParams struct {
	Pagination
	Status string
}

type Pagination struct {
	Offset int
	Limit  int
}

func (params Pagination) GetLimit() int {
	const defaultLimit = 10

	if params.Limit > 0 {
		return params.Limit
	}
	return defaultLimit
}

func (impl *Repo) GetList(ctx context.Context, params GetListParams) ([]models.Application, int, error) {
	qb := impl.builder.
		Select("id, first_name", "last_name", "status", "count(*) over () AS total").
		From(tableName).
		OrderBy("created_at").
		Offset(uint64(params.Offset)).
		Limit(uint64(params.GetLimit()))

	if params.Status != "" {
		qb = qb.Where(squirrel.Eq{"status": params.Status})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var (
		items = make([]models.Application, 0)
		total int
	)
	rows, err := impl.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var row struct {
			models.Application
			Total int `db:"total"`
		}
		err := rows.StructScan(&row)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, row.Application)
		total = row.Total
	}

	return items, total, nil
}

func (impl *Repo) GetByID(ctx context.Context, id uuid.UUID) (models.Application, error) {
	const query = `
		SELECT row_to_json(t)
		FROM ` + tableName + ` AS t
		WHERE t.id = $1
	`
	var item models.Application
	err := impl.db.QueryRowxContext(ctx, query, id).Scan(&item)
	switch {
	case err == sql.ErrNoRows:
		return models.Application{}, ErrNotFound
	case err != nil:
		return models.Application{}, err
	}

	return item, nil
}

func (impl *Repo) Create(ctx context.Context, item models.Application) (uuid.UUID, error) {
	const query = `
		INSERT INTO ` + tableName + ` (first_name, last_name, status)
		VALUES ($1, $2, $3)
		RETURNING ID
	`

	var id uuid.UUID
	err := impl.db.GetContext(ctx, &id, query, item.FirstName, item.LastName, item.Status)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (impl *Repo) UpdateStatus(ctx context.Context, change models.StatusChange) error {
	const query = `
		UPDATE applications
		SET status = $2, updated_at = now()
		WHERE id = $1
	`

	_, err := impl.db.ExecContext(ctx, query, change.ID, change.Status)
	return err
}
