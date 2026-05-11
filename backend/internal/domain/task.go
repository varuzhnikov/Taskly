package domain

import (
	"time"

	"github.com/google/uuid"
)

type Priority int8

const (
	PriorityNone   Priority = 0
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
	PriorityUrgent Priority = 4
)

func (p Priority) Valid() bool {
	return p >= PriorityNone && p <= PriorityUrgent
}

type TaskLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Task struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	ProjectID       *uuid.UUID
	ParentTaskID    *uuid.UUID
	Title           string
	DescriptionMD   string
	DescriptionHTML string
	Links           []TaskLink
	Priority        Priority
	DueAt           *time.Time
	CompletedAt     *time.Time
	SortOrder       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TaskFilter struct {
	ProjectID *uuid.UUID
	LabelID   *uuid.UUID
	Priority  *Priority
	Completed *bool
	DueBefore *time.Time
	DueAfter  *time.Time
	ParentID  *uuid.UUID
	Search    string
	Limit     int
	Offset    int
}
