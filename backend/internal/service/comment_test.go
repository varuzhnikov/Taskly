package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/internal/service"
	"github.com/todoist/backend/testutil"
)

func TestCommentService_TaskRouteMustMatchCommentTask(t *testing.T) {
	pool := testutil.NewPostgresPool(t)
	ctx := testutil.TxContext(t, pool)

	user := testutil.CreateUser(t, ctx, pool, "comment-owner@example.com", "password123")
	taskA := testutil.CreateTask(t, ctx, pool, user.ID, "Task A")
	taskB := testutil.CreateTask(t, ctx, pool, user.ID, "Task B")

	commentRepo := repository.NewCommentRepo(pool)
	taskRepo := repository.NewTaskRepo(pool)
	svc := service.NewCommentService(commentRepo, taskRepo)

	comment := &domain.Comment{
		ID:        uuid.New(),
		TaskID:    taskA.ID,
		UserID:    user.ID,
		Body:      "Route mismatch",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, commentRepo.Create(ctx, comment))

	_, err := svc.GetByID(ctx, taskB.ID, comment.ID, user.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = svc.Update(ctx, service.UpdateCommentInput{
		TaskID: taskB.ID,
		ID:     comment.ID,
		UserID: user.ID,
		Body:   "Still wrong route",
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = svc.Delete(ctx, taskB.ID, comment.ID, user.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
