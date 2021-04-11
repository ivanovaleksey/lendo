package db

import (
	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

type DB struct {
	*sqlx.DB
}

func New(cfg Config) (*DB, error) {
	db, err := sqlx.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func NewTestDB(t *testing.T, cfg Config) *DB {
	driver := uuid.NewV4().String()
	txdb.Register(driver, "postgres", cfg.URL)

	db, err := sqlx.Open(driver, cfg.URL)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	return &DB{DB: db}
}

func (db *DB) ComponentName() string {
	return "db"
}
