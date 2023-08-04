package ledgerstore

import (
	_ "embed"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

//go:embed migrations/0-init-schema.sql
var initSchema string

func registerMigrations(migrator *migrations.Migrator) {
	migrator.RegisterMigrations(
		migrations.Migration{
			Up: func(tx bun.Tx) error {

				_, err := tx.Exec(initSchema)
				if err != nil {
					return err
				}
				// TODO: Call v1 upgrade
				//v1migrations.UpgradeLogs(context.Background(), )

				return nil
			},
		},
	)
}
