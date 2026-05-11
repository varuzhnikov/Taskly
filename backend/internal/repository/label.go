package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
)

type LabelRepo struct {
	pool *pgxpool.Pool
}

func NewLabelRepo(pool *pgxpool.Pool) *LabelRepo {
	return &LabelRepo{pool: pool}
}

func (r *LabelRepo) Create(ctx context.Context, l *domain.Label) error {
	const q = `
		INSERT INTO labels (id, user_id, name, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		l.ID, l.UserID, l.Name, l.Color, l.CreatedAt, l.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create label: %w", err)
	}
	return nil
}

func (r *LabelRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Label, error) {
	const q = `
		SELECT id, user_id, name, color, created_at, updated_at
		FROM labels
		WHERE id = $1 AND user_id = $2`

	var l domain.Label
	err := db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, id, userID).Scan(
		&l.ID, &l.UserID, &l.Name, &l.Color, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get label: %w", err)
	}
	return &l, nil
}

func (r *LabelRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error) {
	const q = `
		SELECT id, user_id, name, color, created_at, updated_at
		FROM labels
		WHERE user_id = $1
		ORDER BY name ASC`

	rows, err := db.ExtractDB(ctx, r.pool).Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	defer rows.Close()

	var labels []*domain.Label
	for rows.Next() {
		var l domain.Label
		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.Color, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan label: %w", err)
		}
		labels = append(labels, &l)
	}
	return labels, rows.Err()
}

func (r *LabelRepo) Update(ctx context.Context, l *domain.Label) error {
	const q = `
		UPDATE labels
		SET name = $1, color = $2, updated_at = $3
		WHERE id = $4 AND user_id = $5`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		l.Name, l.Color, time.Now(), l.ID, l.UserID,
	)
	if err != nil {
		return fmt.Errorf("update label: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *LabelRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	const q = `DELETE FROM labels WHERE id = $1 AND user_id = $2`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete label: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
