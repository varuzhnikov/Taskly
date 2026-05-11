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

type CommentRepo struct {
	pool *pgxpool.Pool
}

func NewCommentRepo(pool *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{pool: pool}
}

func (r *CommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	const q = `
		INSERT INTO comments (id, task_id, user_id, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		c.ID, c.TaskID, c.UserID, c.Body, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	const q = `
		SELECT id, task_id, user_id, body, created_at, updated_at
		FROM comments
		WHERE id = $1`

	var c domain.Comment
	err := db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, id).Scan(
		&c.ID, &c.TaskID, &c.UserID, &c.Body, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return &c, nil
}

func (r *CommentRepo) ListByTask(ctx context.Context, taskID, userID uuid.UUID) ([]*domain.Comment, error) {
	// userID is used to verify the task belongs to the caller (join check in service layer).
	const q = `
		SELECT id, task_id, user_id, body, created_at, updated_at
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC`

	rows, err := db.ExtractDB(ctx, r.pool).Query(ctx, q, taskID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Body, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, &c)
	}
	return comments, rows.Err()
}

func (r *CommentRepo) Update(ctx context.Context, c *domain.Comment) error {
	const q = `
		UPDATE comments
		SET body = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		c.Body, time.Now(), c.ID, c.UserID,
	)
	if err != nil {
		return fmt.Errorf("update comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	const q = `DELETE FROM comments WHERE id = $1 AND user_id = $2`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
