package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
)

type TokenRepo struct {
	pool *pgxpool.Pool
}

func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{pool: pool}
}

func (r *TokenRepo) Create(ctx context.Context, t *domain.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		t.ID, t.UserID, t.TokenHash, t.ExpiresAt, t.CreatedAt, t.UserAgent, t.IP,
	)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepo) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, created_at, user_agent, ip
		FROM refresh_tokens
		WHERE token_hash = $1`

	var t domain.RefreshToken
	err := db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.UserAgent, &t.IP,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return &t, nil
}

func (r *TokenRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM refresh_tokens WHERE id = $1`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepo) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	const q = `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("delete all refresh tokens for user: %w", err)
	}
	return nil
}

func (r *TokenRepo) DeleteExpired(ctx context.Context) error {
	const q = `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q)
	if err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}
