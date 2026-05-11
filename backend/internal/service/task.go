package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/render"
	"github.com/todoist/backend/internal/repository"
)

type TaskService struct {
	tasks     *repository.TaskRepo
	projects  *repository.ProjectRepo
	labels    *repository.LabelRepo
	txManager *db.TxManager
}

func NewTaskService(
	tasks *repository.TaskRepo,
	projects *repository.ProjectRepo,
	labels *repository.LabelRepo,
	txManager *db.TxManager,
) *TaskService {
	return &TaskService{
		tasks:     tasks,
		projects:  projects,
		labels:    labels,
		txManager: txManager,
	}
}

type CreateTaskInput struct {
	UserID       uuid.UUID
	ProjectID    *uuid.UUID
	ParentTaskID *uuid.UUID
	Title        string
	Description  string
	Links        []domain.TaskLink
	Priority     domain.Priority
	DueAt        *time.Time
	SortOrder    int
	LabelIDs     []uuid.UUID
}

type UpdateTaskInput struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	ProjectID    *uuid.UUID
	ParentTaskID *uuid.UUID
	Title        string
	Description  string
	Links        []domain.TaskLink
	Priority     domain.Priority
	DueAt        *time.Time
	SortOrder    int
	LabelIDs     []uuid.UUID
}

func (s *TaskService) Create(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	if err := s.validateOwnership(ctx, in.UserID, in.ProjectID, in.ParentTaskID, in.LabelIDs); err != nil {
		return nil, err
	}

	html, err := render.ToHTML(in.Description)
	if err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}

	t := &domain.Task{
		ID:              uuid.New(),
		UserID:          in.UserID,
		ProjectID:       in.ProjectID,
		ParentTaskID:    in.ParentTaskID,
		Title:           in.Title,
		DescriptionMD:   in.Description,
		DescriptionHTML: html,
		Links:           in.Links,
		Priority:        in.Priority,
		DueAt:           in.DueAt,
		SortOrder:       in.SortOrder,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.tasks.Create(ctx, t); err != nil {
			return err
		}
		if len(in.LabelIDs) > 0 {
			return s.tasks.SetLabels(ctx, t.ID, in.UserID, in.LabelIDs)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

func (s *TaskService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Task, error) {
	t, err := s.tasks.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	labels, err := s.tasks.GetLabels(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task labels: %w", err)
	}
	_ = labels // caller may embed via a richer response DTO; the repo method is tested separately
	return t, nil
}

func (s *TaskService) List(ctx context.Context, userID uuid.UUID, f domain.TaskFilter) ([]*domain.Task, error) {
	return s.tasks.List(ctx, userID, f)
}

func (s *TaskService) Update(ctx context.Context, in UpdateTaskInput) (*domain.Task, error) {
	t, err := s.tasks.GetByID(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.validateOwnership(ctx, in.UserID, in.ProjectID, in.ParentTaskID, in.LabelIDs); err != nil {
		return nil, err
	}

	html, err := render.ToHTML(in.Description)
	if err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}

	t.ProjectID = in.ProjectID
	t.ParentTaskID = in.ParentTaskID
	t.Title = in.Title
	t.DescriptionMD = in.Description
	t.DescriptionHTML = html
	t.Links = in.Links
	t.Priority = in.Priority
	t.DueAt = in.DueAt
	t.SortOrder = in.SortOrder

	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.tasks.Update(ctx, t); err != nil {
			return err
		}
		return s.tasks.SetLabels(ctx, t.ID, in.UserID, in.LabelIDs)
	})
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return t, nil
}

func (s *TaskService) Complete(ctx context.Context, id, userID uuid.UUID) (*domain.Task, error) {
	t, err := s.tasks.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	t.CompletedAt = &now
	if err := s.tasks.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("complete task: %w", err)
	}
	return t, nil
}

func (s *TaskService) Uncomplete(ctx context.Context, id, userID uuid.UUID) (*domain.Task, error) {
	t, err := s.tasks.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	t.CompletedAt = nil
	if err := s.tasks.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("uncomplete task: %w", err)
	}
	return t, nil
}

func (s *TaskService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.tasks.Delete(ctx, id, userID)
}

func (s *TaskService) GetLabels(ctx context.Context, taskID, userID uuid.UUID) ([]*domain.Label, error) {
	// Ownership check before exposing labels.
	if _, err := s.tasks.GetByID(ctx, taskID, userID); err != nil {
		return nil, err
	}
	return s.tasks.GetLabels(ctx, taskID)
}

func (s *TaskService) validateOwnership(
	ctx context.Context,
	userID uuid.UUID,
	projectID, parentTaskID *uuid.UUID,
	labelIDs []uuid.UUID,
) error {
	if projectID != nil {
		if _, err := s.projects.GetByID(ctx, *projectID, userID); err != nil {
			return err
		}
	}
	if parentTaskID != nil {
		if _, err := s.tasks.GetByID(ctx, *parentTaskID, userID); err != nil {
			return err
		}
	}
	for _, labelID := range labelIDs {
		if _, err := s.labels.GetByID(ctx, labelID, userID); err != nil {
			return err
		}
	}
	return nil
}
