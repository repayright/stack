package systemstore

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	storageerrors "github.com/formancehq/ledger/pkg/storage"
	"github.com/uptrace/bun"
)

type configuration struct {
	bun.BaseModel `bun:"_system.configuration,alias:configuration"`

	Key     string    `bun:"key,type:varchar(255),pk"` // Primary key
	Value   string    `bun:"value,type:text"`
	AddedAt core.Time `bun:"addedAt,type:timestamp"`
}

func (s *Store) CreateConfigurationTable(ctx context.Context) error {
	_, err := s.db.NewCreateTable().
		Model((*configuration)(nil)).
		IfNotExists().
		Exec(ctx)

	return storageerrors.PostgresError(err)
}

func (s *Store) GetConfiguration(ctx context.Context, key string) (string, error) {
	query := s.db.NewSelect().
		Model((*configuration)(nil)).
		Column("value").
		Where("key = ?", key).
		Limit(1).
		String()

	row := s.db.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return "", storageerrors.PostgresError(row.Err())
	}
	var value string
	if err := row.Scan(&value); err != nil {
		return "", storageerrors.PostgresError(err)
	}

	return value, nil
}

func (s *Store) InsertConfiguration(ctx context.Context, key, value string) error {
	config := &configuration{
		Key:     key,
		Value:   value,
		AddedAt: core.Now(),
	}

	_, err := s.db.NewInsert().
		Model(config).
		Exec(ctx)

	return storageerrors.PostgresError(err)
}
