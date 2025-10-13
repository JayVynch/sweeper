package test_helpers

import (
	"context"
	"testing"

	"github.com/JayVynch/sweeper/database"
	"github.com/stretchr/testify/require"
)

func TearDownDB(ctx context.Context, t *testing.T, db *database.DB) {
	t.Helper()

	err := db.Truncate(ctx)

	require.NoError(t, err)
}
