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

//go:embed migrations/0-init-schema.sql
var initSchema string

func registerMigrations(migrator *migrations.Migrator, schema storage.Schema) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Up: func(tx bun.Tx) error {

				v1SchemaExists := false
				row := tx.QueryRow(`
					SELECT EXISTS (
						SELECT FROM
							pg_tables
						WHERE
							schemaname = ? AND
							tablename  = 'log'
						);
				`, schema.Name())
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
