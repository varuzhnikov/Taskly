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

type ProjectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{pool: pool}
}

func (r *ProjectRepo) Create(ctx context.Context, p *domain.Project) error {
	const q = `
		INSERT INTO projects (id, user_id, name, color, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		p.ID, p.UserID, p.Name, p.Color, p.SortOrder, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *ProjectRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Project, error) {
	const q = `
		SELECT id, user_id, name, color, sort_order, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2`

	p, err := scanProject(db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, id, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get project: %w", err)
	}
	return p, nil
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	const q = `
		SELECT id, user_id, name, color, sort_order, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY sort_order ASC, created_at ASC`

	rows, err := db.ExtractDB(ctx, r.pool).Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepo) Update(ctx context.Context, p *domain.Project) error {
	const q = `
		UPDATE projects
		SET name = $1, color = $2, sort_order = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		p.Name, p.Color, p.SortOrder, time.Now(), p.ID, p.UserID,
	)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	const q = `DELETE FROM projects WHERE id = $1 AND user_id = $2`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// scanProject accepts anything with a Scan method — works for both pgx.Row and pgx.Rows.
func scanProject(row interface{ Scan(...any) error }) (*domain.Project, error) {
	var p domain.Project
	err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.Color, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
