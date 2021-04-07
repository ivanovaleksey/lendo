package applicationsRepo

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
)

const (
	tableName = "applications"
)

type impl struct {
	db      *db.DB
	builder squirrel.StatementBuilderType
}

func New(db *db.DB) Repo {
	repo := &impl{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
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

func (impl *impl) GetList(ctx context.Context, params GetListParams) ([]models.Application, int, error) {
	qb := impl.builder.
		Select("*").
		From(tableName).
		OrderBy("created_at").
		Offset(uint64(params.Offset)).
		Limit(uint64(params.GetLimit()))

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

func (impl *impl) GetByID(ctx context.Context, id uuid.UUID) (models.Application, error) {
	const query = `
		SELECT id, first_name, last_name, status
		FROM ` + tableName + `
		WHERE id = $1
	`
	var item models.Application
	err := impl.db.GetContext(ctx, &item, query, id)
	if err != nil {
		// todo: handle not found
		return models.Application{}, err
	}

	return item, nil
}

func (impl *impl) Create(ctx context.Context, item models.NewApplication) (uuid.UUID, error) {
	const query = `
		INSERT INTO ` + tableName + ` (first_name, last_name, status)
		VALUES ($1, $2, $3)
		RETURNING ID
	`

	var id uuid.UUID
	err := impl.db.GetContext(ctx, &id, query, item.FirstName, item.LastName, "")
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
