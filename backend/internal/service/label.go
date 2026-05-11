package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
)

type LabelService struct {
	labels *repository.LabelRepo
}

func NewLabelService(labels *repository.LabelRepo) *LabelService {
	return &LabelService{labels: labels}
}

type CreateLabelInput struct {
	UserID uuid.UUID
	Name   string
	Color  string
}

type UpdateLabelInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
	Color  string
}

func (s *LabelService) Create(ctx context.Context, in CreateLabelInput) (*domain.Label, error) {
	l := &domain.Label{
		ID:        uuid.New(),
		UserID:    in.UserID,
		Name:      in.Name,
		Color:     in.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.labels.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}
	return l, nil
}

func (s *LabelService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Label, error) {
	return s.labels.GetByID(ctx, id, userID)
}

func (s *LabelService) List(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error) {
	return s.labels.ListByUser(ctx, userID)
}

func (s *LabelService) Update(ctx context.Context, in UpdateLabelInput) (*domain.Label, error) {
	l, err := s.labels.GetByID(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}
	l.Name = in.Name
	l.Color = in.Color
	if err := s.labels.Update(ctx, l); err != nil {
		return nil, fmt.Errorf("update label: %w", err)
	}
	return l, nil
}

func (s *LabelService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.labels.Delete(ctx, id, userID)
}
