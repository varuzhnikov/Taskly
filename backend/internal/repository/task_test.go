package repository_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/testutil"
)

func TestTaskRepo_Create(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("creates task with all fields", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "task-user@example.com", "pass")
		repo := repository.NewTaskRepo(pool)

		due := time.Now().Add(24 * time.Hour)
		task := &domain.Task{
			ID:       uuid.New(),
			UserID:   user.ID,
			Title:    "Buy groceries",
			Links:    []domain.TaskLink{{URL: "https://example.com", Title: "Store"}},
			Priority: domain.PriorityHigh,
			DueAt:    &due,
		}
		require.NoError(t, repo.Create(ctx, task))
	})

	t.Run("creates task with nil optional fields", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "task-nil@example.com", "pass")
		repo := repository.NewTaskRepo(pool)

		task := &domain.Task{
			ID:     uuid.New(),
			UserID: user.ID,
			Title:  "Minimal task",
			Links:  []domain.TaskLink{},
		}
		require.NoError(t, repo.Create(ctx, task))
	})
}

func TestTaskRepo_GetByID(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("returns task for correct owner", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "owner@example.com", "pass")
		created := testutil.CreateTask(t, ctx, pool, user.ID, "My Task")

		repo := repository.NewTaskRepo(pool)
		got, err := repo.GetByID(ctx, created.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, "My Task", got.Title)
	})

	t.Run("different user cannot see the task", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "owner2@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "other2@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, owner.ID, "Private task")

		repo := repository.NewTaskRepo(pool)
		_, err := repo.GetByID(ctx, task.ID, other.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestTaskRepo_List(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("lists only tasks belonging to user", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "listuser@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "listother@example.com", "pass")

		testutil.CreateTask(t, ctx, pool, user.ID, "Task A")
		testutil.CreateTask(t, ctx, pool, user.ID, "Task B")
		testutil.CreateTask(t, ctx, pool, other.ID, "Other's task")

		repo := repository.NewTaskRepo(pool)
		tasks, err := repo.List(ctx, user.ID, domain.TaskFilter{Limit: 50})
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("filter by completion status", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "compuser@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, user.ID, "Done task")

		// Mark as complete.
		now := time.Now()
		task.CompletedAt = &now
		repo := repository.NewTaskRepo(pool)
		require.NoError(t, repo.Update(ctx, task))

		completed := true
		tasks, err := repo.List(ctx, user.ID, domain.TaskFilter{Completed: &completed, Limit: 50})
		require.NoError(t, err)
		assert.Len(t, tasks, 1)

		incomplete := false
		tasks, err = repo.List(ctx, user.ID, domain.TaskFilter{Completed: &incomplete, Limit: 50})
		require.NoError(t, err)
		assert.Len(t, tasks, 0)
	})
}

func TestTaskRepo_SetGetLabels(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("attaches and reads labels", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lbluser@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, user.ID, "Labelled task")
		l1 := testutil.CreateLabel(t, ctx, pool, user.ID, "work")
		l2 := testutil.CreateLabel(t, ctx, pool, user.ID, "urgent")

		repo := repository.NewTaskRepo(pool)
		require.NoError(t, repo.SetLabels(ctx, task.ID, user.ID, []uuid.UUID{l1.ID, l2.ID}))

		labels, err := repo.GetLabels(ctx, task.ID)
		require.NoError(t, err)
		assert.Len(t, labels, 2)
	})

	t.Run("clearing labels leaves zero labels", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "lblclear@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, user.ID, "Task with labels")
		l := testutil.CreateLabel(t, ctx, pool, user.ID, "old")

		repo := repository.NewTaskRepo(pool)
		require.NoError(t, repo.SetLabels(ctx, task.ID, user.ID, []uuid.UUID{l.ID}))
		require.NoError(t, repo.SetLabels(ctx, task.ID, user.ID, nil)) // clear

		labels, err := repo.GetLabels(ctx, task.ID)
		require.NoError(t, err)
		assert.Empty(t, labels)
	})
}

func TestTaskRepo_Delete(t *testing.T) {
	pool := testutil.NewPostgresPool(t)

	t.Run("owner can delete their task", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		user := testutil.CreateUser(t, ctx, pool, "delowner@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, user.ID, "To be deleted")

		repo := repository.NewTaskRepo(pool)
		require.NoError(t, repo.Delete(ctx, task.ID, user.ID))

		_, err := repo.GetByID(ctx, task.ID, user.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("other user cannot delete the task", func(t *testing.T) {
		ctx := testutil.TxContext(t, pool)
		owner := testutil.CreateUser(t, ctx, pool, "delowner2@example.com", "pass")
		other := testutil.CreateUser(t, ctx, pool, "delother2@example.com", "pass")
		task := testutil.CreateTask(t, ctx, pool, owner.ID, "Not yours")

		repo := repository.NewTaskRepo(pool)
		err := repo.Delete(ctx, task.ID, other.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
