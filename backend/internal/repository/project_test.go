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

func TestProjectRepo_CRUD(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("create and retrieve project", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "proj-user@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, user.ID, "Work")

		repo := repository.NewProjectRepo(pool)
		got, err := repo.GetByID(ctx, p.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Work", got.Name)
	})

	t.Run("other user cannot see the project", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "proj-owner@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "proj-other@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, owner.ID, "Secret Project")

		repo := repository.NewProjectRepo(pool)
		_, err := repo.GetByID(ctx, p.ID, other.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("list returns only user projects in sort order", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "proj-list@example.com", "pass")
		testutil.CreateProject(t, ctx, pool, user.ID, "B Project")
		testutil.CreateProject(t, ctx, pool, user.ID, "A Project")

		repo := repository.NewProjectRepo(pool)
		projects, err := repo.ListByUser(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, projects, 2)
	})

	t.Run("update project name", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "proj-upd@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, user.ID, "Old Name")

		p.Name = "New Name"
		repo := repository.NewProjectRepo(pool)
		require.NoError(t, repo.Update(ctx, p))

		got, err := repo.GetByID(ctx, p.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "New Name", got.Name)
	})

	t.Run("delete project", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "proj-del@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, user.ID, "To Delete")

		repo := repository.NewProjectRepo(pool)
		require.NoError(t, repo.Delete(ctx, p.ID, user.ID))

		_, err := repo.GetByID(ctx, p.ID, user.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("delete by wrong user returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "proj-del-own@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "proj-del-oth@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, owner.ID, "Protected")

		repo := repository.NewProjectRepo(pool)
		err := repo.Delete(ctx, p.ID, other.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestProjectRepo_Update_NonExistent(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("updating a non-existent project returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		repo := repository.NewProjectRepo(pool)

		ghost := &domain.Project{ID: uuid.New(), UserID: uuid.New(), Name: "Ghost"}
		err := repo.Update(ctx, ghost)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
