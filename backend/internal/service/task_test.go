package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/internal/service"
	"github.com/todoist/backend/testutil"
)

func newTestTaskService(t *testing.T) *service.TaskService {
	t.Helper()
	pool := testutil.NewPostgresPool(t)
	return service.NewTaskService(
		repository.NewTaskRepo(pool),
		repository.NewProjectRepo(pool),
		repository.NewLabelRepo(pool),
		db.NewTxManager(pool),
	)
}

func TestTaskService_Create_RejectsForeignReferences(t *testing.T) {
	pool := testutil.NewPostgresPool(t)
	ctx := testutil.TxContext(t, pool)

	owner := testutil.CreateUser(t, ctx, pool, "task-owner@example.com", "password123")
	other := testutil.CreateUser(t, ctx, pool, "task-other@example.com", "password123")
	foreignProject := testutil.CreateProject(t, ctx, pool, other.ID, "Foreign Project")
	foreignParent := testutil.CreateTask(t, ctx, pool, other.ID, "Foreign Parent")
	foreignLabel := testutil.CreateLabel(t, ctx, pool, other.ID, "foreign-label")

	svc := service.NewTaskService(
		repository.NewTaskRepo(pool),
		repository.NewProjectRepo(pool),
		repository.NewLabelRepo(pool),
		db.NewTxManager(pool),
	)

	_, err := svc.Create(ctx, service.CreateTaskInput{
		UserID:       owner.ID,
		ProjectID:    &foreignProject.ID,
		ParentTaskID: &foreignParent.ID,
		Title:        "Should fail",
		LabelIDs:     []uuid.UUID{foreignLabel.ID},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
