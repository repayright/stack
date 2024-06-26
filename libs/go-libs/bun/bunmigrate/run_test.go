package bunmigrate

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestRunMigrate(t *testing.T) {
	require.NoError(t, pgtesting.CreatePostgresServer())
	t.Cleanup(func() {
		require.NoError(t, pgtesting.DestroyPostgresServer())
	})

	connectionOptions := &bunconnect.ConnectionOptions{
		DatabaseSourceName: pgtesting.Server().GetDatabaseDSN("testing"),
	}
	executor := func(args []string, db *bun.DB) error {
		return nil
	}

	err := run(logging.TestingContext(), os.Stdout, []string{}, connectionOptions, executor)
	require.NoError(t, err)

	// Must be idempotent
	err = run(logging.TestingContext(), os.Stdout, []string{}, connectionOptions, executor)
	require.NoError(t, err)
}
