package ledgerstore

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/formancehq/ledger/pkg/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Store struct {
	db       *bun.DB
	onDelete func(ctx context.Context) error

	once sync.Once

	isInitialized bool
	name          string
}

func (s *Store) Name() string {
	return s.name
}

func (d *Store) GetDatabase() *bun.DB {
	return d.db
}

func (s *Store) Delete(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "delete schema ? cascade", s.name)
	if err != nil {
		return err
	}
	return errors.Wrap(s.onDelete(ctx), "deleting ledger store")
}

func (s *Store) IsInitialized() bool {
	return s.isInitialized
}

func (s *Store) prepareTransaction(ctx context.Context) (bun.Tx, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return tx, err
	}
	if _, err := tx.Exec(fmt.Sprintf(`set search_path = "%s"`, s.Name())); err != nil {
		return tx, err
	}
	return tx, nil
}

func (s *Store) withTransaction(ctx context.Context, callback func(tx bun.Tx) error) error {
	tx, err := s.prepareTransaction(ctx)
	if err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		_ = tx.Rollback()
		return storage.PostgresError(err)
	}
	return tx.Commit()
}

func New(
	db *bun.DB,
	name string,
	onDelete func(ctx context.Context) error,
) (*Store, error) {
	return &Store{
		db:       db,
		name:     name,
		onDelete: onDelete,
	}, nil
}
