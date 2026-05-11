package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/testutil"
)

func TestUserRepo_Create(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("creates user successfully", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		u := &domain.User{
			ID:           uuid.New(),
			Email:        "alice@example.com",
			PasswordHash: "$2a$12$hashed",
		}
		require.NoError(t, repo.Create(ctx, u))
	})

	t.Run("duplicate email returns ErrAlreadyExists", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		u := &domain.User{ID: uuid.New(), Email: "bob@example.com", PasswordHash: "x"}
		require.NoError(t, repo.Create(ctx, u))

		dup := &domain.User{ID: uuid.New(), Email: "bob@example.com", PasswordHash: "y"}
		err := repo.Create(ctx, dup)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestUserRepo_GetByEmail(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("finds existing user", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		want := testutil.CreateUser(t, ctx, pool, "charlie@example.com", "password123")
		got, err := repo.GetByEmail(ctx, want.Email)
		require.NoError(t, err)
		assert.Equal(t, want.ID, got.ID)
		assert.Equal(t, want.Email, got.Email)
	})

	t.Run("unknown email returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		_, err := repo.GetByEmail(ctx, "nobody@nowhere.com")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_GetByID(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("finds existing user by id", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		want := testutil.CreateUser(t, ctx, pool, "dave@example.com", "password123")
		got, err := repo.GetByID(ctx, want.ID)
		require.NoError(t, err)
		assert.Equal(t, want.ID, got.ID)
	})

	t.Run("random uuid returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		_, err := repo.GetByID(ctx, uuid.New())
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUserRepo_Update(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("updates email", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		u := testutil.CreateUser(t, ctx, pool, "eve@example.com", "password123")
		u.Email = "eve-new@example.com"
		require.NoError(t, repo.Update(ctx, u))

		got, err := repo.GetByID(ctx, u.ID)
		require.NoError(t, err)
		assert.Equal(t, "eve-new@example.com", got.Email)
	})

	t.Run("updating non-existent user returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewUserRepo(pool)

		ghost := &domain.User{ID: uuid.New(), Email: "ghost@example.com"}
		assert.ErrorIs(t, repo.Update(ctx, ghost), domain.ErrNotFound)
	})
}
