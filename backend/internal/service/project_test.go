package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/internal/service"
	"github.com/todoist/backend/testutil"
)

func TestProjectService_Delete(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("delete removes the project", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-del@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, user.ID, "To Delete")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		err := svc.Delete(ctx, p.ID, user.ID)
		require.NoError(t, err)

		_, err = repository.NewProjectRepo(pool).GetByID(ctx, p.ID, user.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("delete by wrong user returns error", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "svc-del-own@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "svc-del-oth@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, owner.ID, "Protected")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		err := svc.Delete(ctx, p.ID, other.ID)
		assert.Error(t, err)

		// project must still exist for the owner
		_, err = repository.NewProjectRepo(pool).GetByID(ctx, p.ID, owner.ID)
		assert.NoError(t, err)
	})

	t.Run("delete non-existent project returns error", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-del-ghost@example.com", "pass")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		err := svc.Delete(ctx, uuid.New(), user.ID)
		assert.Error(t, err)
	})
}

func TestProjectService_Create(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("create returns project with all input fields set", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-create@example.com", "pass")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		p, err := svc.Create(ctx, service.CreateProjectInput{
			UserID:    user.ID,
			Name:      "Work",
			Color:     "#db4035",
			SortOrder: 5,
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, p.ID)
		assert.Equal(t, user.ID, p.UserID)
		assert.Equal(t, "Work", p.Name)
		assert.Equal(t, "#db4035", p.Color)
		assert.Equal(t, 5, p.SortOrder)
	})

	t.Run("created project is retrievable from the database", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-create2@example.com", "pass")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		p, err := svc.Create(ctx, service.CreateProjectInput{
			UserID: user.ID,
			Name:   "Persisted",
			Color:  "#cccccc",
		})
		require.NoError(t, err)

		got, err := repository.NewProjectRepo(pool).GetByID(ctx, p.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Persisted", got.Name)
	})
}

func TestProjectService_Update(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("update changes name and color", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-upd@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, user.ID, "Old Name")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		updated, err := svc.Update(ctx, service.UpdateProjectInput{
			ID:     p.ID,
			UserID: user.ID,
			Name:   "New Name",
			Color:  "#ff0000",
		})
		require.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, "#ff0000", updated.Color)
	})

	t.Run("update non-existent project returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "svc-upd-ghost@example.com", "pass")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		_, err := svc.Update(ctx, service.UpdateProjectInput{
			ID:     uuid.New(),
			UserID: user.ID,
			Name:   "Ghost",
		})
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("update by wrong user returns ErrNotFound", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "svc-upd-own@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "svc-upd-oth@example.com", "pass")
		p := testutil.CreateProject(t, ctx, pool, owner.ID, "Owner's Project")

		svc := service.NewProjectService(
			repository.NewProjectRepo(pool),
			repository.NewTaskRepo(pool),
		)

		_, err := svc.Update(ctx, service.UpdateProjectInput{
			ID:     p.ID,
			UserID: other.ID,
			Name:   "Hijacked",
		})
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
