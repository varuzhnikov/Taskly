package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
)

type ProjectService struct {
	projects *repository.ProjectRepo
	tasks    *repository.TaskRepo
}

func NewProjectService(projects *repository.ProjectRepo, tasks *repository.TaskRepo) *ProjectService {
	return &ProjectService{projects: projects, tasks: tasks}
}

type CreateProjectInput struct {
	UserID    uuid.UUID
	Name      string
	Color     string
	SortOrder int
}

type UpdateProjectInput struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Color     string
	SortOrder int
}

func (s *ProjectService) Create(ctx context.Context, in CreateProjectInput) (*domain.Project, error) {
	p := &domain.Project{
		ID:        uuid.New(),
		UserID:    in.UserID,
		Name:      in.Name,
		Color:     in.Color,
		SortOrder: in.SortOrder,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.projects.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Project, error) {
	return s.projects.GetByID(ctx, id, userID)
}

func (s *ProjectService) List(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	return s.projects.ListByUser(ctx, userID)
}

func (s *ProjectService) Update(ctx context.Context, in UpdateProjectInput) (*domain.Project, error) {
	p, err := s.projects.GetByID(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}
	p.Name = in.Name
	p.Color = in.Color
	p.SortOrder = in.SortOrder
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if err := s.projects.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
