package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
)

type CommentService struct {
	comments *repository.CommentRepo
	tasks    *repository.TaskRepo
}

func NewCommentService(comments *repository.CommentRepo, tasks *repository.TaskRepo) *CommentService {
	return &CommentService{comments: comments, tasks: tasks}
}

type CreateCommentInput struct {
	TaskID uuid.UUID
	UserID uuid.UUID
	Body   string
}

type UpdateCommentInput struct {
	TaskID  uuid.UUID
	ID     uuid.UUID
	UserID uuid.UUID
	Body   string
}

func (s *CommentService) Create(ctx context.Context, in CreateCommentInput) (*domain.Comment, error) {
	// Verify the task belongs to the caller before allowing a comment.
	if _, err := s.tasks.GetByID(ctx, in.TaskID, in.UserID); err != nil {
		return nil, err
	}

	c := &domain.Comment{
		ID:        uuid.New(),
		TaskID:    in.TaskID,
		UserID:    in.UserID,
		Body:      in.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.comments.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return c, nil
}

func (s *CommentService) ListByTask(ctx context.Context, taskID, userID uuid.UUID) ([]*domain.Comment, error) {
	// Verify task ownership before exposing comments.
	if _, err := s.tasks.GetByID(ctx, taskID, userID); err != nil {
		return nil, err
	}
	return s.comments.ListByTask(ctx, taskID, userID)
}

func (s *CommentService) Update(ctx context.Context, in UpdateCommentInput) (*domain.Comment, error) {
	c, err := s.comments.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if c.TaskID != in.TaskID {
		return nil, domain.ErrNotFound
	}
	if c.UserID != in.UserID {
		return nil, domain.ErrForbidden
	}
	c.Body = in.Body
	if err := s.comments.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}
	return c, nil
}

func (s *CommentService) Delete(ctx context.Context, taskID, id, userID uuid.UUID) error {
	c, err := s.comments.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c.TaskID != taskID {
		return domain.ErrNotFound
	}
	if c.UserID != userID {
		return domain.ErrForbidden
	}
	return s.comments.Delete(ctx, id, userID)
}

// GetByID returns a single comment, verifying ownership.
func (s *CommentService) GetByID(ctx context.Context, taskID, id, userID uuid.UUID) (*domain.Comment, error) {
	c, err := s.comments.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.TaskID != taskID {
		return nil, domain.ErrNotFound
	}
	if c.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return c, nil
}
