package domain

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Color     string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
}
