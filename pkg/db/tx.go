package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type TxFactory interface {
	Begin(context.Context) (Tx, error)
}

type Tx interface {
	SQLTx
	Do(context.Context, TxFunc) error
}

type TxFunc func(context.Context, SQLTx) error

type SQLTx interface {
	sqlx.QueryerContext
	sqlx.ExecerContext

	Rollback() error
	Commit() error
}

func NewTxFactory(db *DB) TxFactory {
	return txFactory{db: db}
}

type txFactory struct {
	db *DB
}

func (f txFactory) Begin(ctx context.Context) (Tx, error) {
	tx, err := f.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return NewTx(tx), nil
}

func NewTx(sqlTx SQLTx) Tx {
	return tx{SQLTx: sqlTx}
}

type tx struct {
	SQLTx
}

func (tx tx) Do(ctx context.Context, fn TxFunc) (err error) {
	defer func() {
		rec := recover()
		if rec != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Errorf("can't rollback tx after panic: %v", rErr)
			}
			err = errors.Errorf("panic in tx: %v", rec)
		}
	}()

	err = fn(ctx, tx)
	if err == nil {
		err = tx.Commit()
	} else {
		if rErr := tx.Rollback(); rErr != nil {
			log.Errorf("can't rollback tx: %v", rErr)
		}
	}
	return err
}
