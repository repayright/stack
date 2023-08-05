package ledgerstore

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/formancehq/ledger/pkg/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type Store struct {
	schema   storage.Schema
	onDelete func(ctx context.Context) error

	once sync.Once

	isInitialized bool
}

func (s *Store) Schema() storage.Schema {
	return s.schema
}

func (s *Store) Name() string {
	return s.schema.Name()
}

func (s *Store) Delete(ctx context.Context) error {
	if err := s.schema.Delete(ctx); err != nil {
		return err
	}
	return errors.Wrap(s.onDelete(ctx), "deleting ledger store")
}

func (s *Store) IsInitialized() bool {
	return s.isInitialized
}

func (s *Store) prepareTransaction(ctx context.Context) (*storage.Tx, error) {
	tx, err := s.schema.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return tx, err
	}
	if _, err := tx.Exec(fmt.Sprintf(`set search_path = "%s"`, s.Name())); err != nil {
		return tx, err
	}
	return tx, nil
}

func (s *Store) withTransaction(ctx context.Context, callback func(tx *storage.Tx) error) error {
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
	schema storage.Schema,
	onDelete func(ctx context.Context) error,
) (*Store, error) {
	return &Store{
		schema:   schema,
		onDelete: onDelete,
	}, nil
}
