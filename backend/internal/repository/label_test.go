package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/testutil"
)

func TestLabelRepo_CRUD(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("creates and retrieves label", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lbl-user@example.com", "pass")
		l := testutil.CreateLabel(t, ctx, pool, user.ID, "urgent")

		repo := repository.NewLabelRepo(pool)
		got, err := repo.GetByID(ctx, l.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "urgent", got.Name)
	})

	t.Run("other user cannot access label", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "lbl-own@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "lbl-oth@example.com", "pass")
		l := testutil.CreateLabel(t, ctx, pool, owner.ID, "private-label")

		repo := repository.NewLabelRepo(pool)
		_, err := repo.GetByID(ctx, l.ID, other.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("list returns only user labels sorted by name", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lbl-list@example.com", "pass")
		testutil.CreateLabel(t, ctx, pool, user.ID, "zebra")
		testutil.CreateLabel(t, ctx, pool, user.ID, "apple")

		repo := repository.NewLabelRepo(pool)
		labels, err := repo.ListByUser(ctx, user.ID)
		require.NoError(t, err)
		require.Len(t, labels, 2)
		assert.Equal(t, "apple", labels[0].Name)
		assert.Equal(t, "zebra", labels[1].Name)
	})

	t.Run("updates label name and color", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lbl-upd@example.com", "pass")
		l := testutil.CreateLabel(t, ctx, pool, user.ID, "old-name")

		l.Name = "new-name"
		l.Color = "#ff0000"

		repo := repository.NewLabelRepo(pool)
		require.NoError(t, repo.Update(ctx, l))

		got, err := repo.GetByID(ctx, l.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "new-name", got.Name)
		assert.Equal(t, "#ff0000", got.Color)
	})

	t.Run("deletes label", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lbl-del@example.com", "pass")
		l := testutil.CreateLabel(t, ctx, pool, user.ID, "doomed")

		repo := repository.NewLabelRepo(pool)
		require.NoError(t, repo.Delete(ctx, l.ID, user.ID))

		_, err := repo.GetByID(ctx, l.ID, user.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
