package ledgerstore

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/formancehq/ledger/pkg/storage"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)


func (s *Store) getMigrator() *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(s.Name(), true))
	registerMigrations(migrator, s.schema)
	return migrator
}

func (s *Store) Migrate(ctx context.Context) (bool, error) {
	migrator := s.getMigrator()

	if err := migrator.Up(ctx, s.schema.IDB); err != nil {
		return false, err
	}

	// TODO: Update migrations package to return modifications
	return false, nil
}

func (s *Store) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return s.getMigrator().GetMigrations(ctx, s.schema.IDB)
}

//go:embed migrations/0-init-schema.sql
var initSchema string

func registerMigrations(migrator *migrations.Migrator, schema storage.Schema) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			Up: func(tx bun.Tx) error {

				v1SchemaExists := false
				row := tx.QueryRow(`select exists (
					select from pg_tables
					where schemaname = ? and tablename  = 'log'
				)`, schema.Name())
				if row.Err() != nil {
					return row.Err()
				}
				var ret string
				if err := row.Scan(&ret); err != nil {
					panic(err)
				}
				v1SchemaExists = ret != "false"

				if v1SchemaExists {
					_, err := tx.Exec(`alter schema rename ? to ?`, schema.Name(), fmt.Sprintf(schema.Name()+oldSchemaRenameSuffix))
					if err != nil {
						return errors.Wrap(err, "renaming old schema")
					}
					if err := schema.Create(context.Background()); err != nil {
						return errors.Wrap(err, "creating new schema")
					}
				}

				_, err := tx.Exec(initSchema)
				if err != nil {
					return errors.Wrap(err, "initializing new schema")
				}

				if v1SchemaExists {
					if err := migrateLogs(context.Background(), fmt.Sprintf(schema.Name()+oldSchemaRenameSuffix), schema.Name(), tx); err != nil {
						return errors.Wrap(err, "migrating logs")
					}
				}

				return nil
			},
		},
	)
}
