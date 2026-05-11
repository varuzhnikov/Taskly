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

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: email already registered", domain.ErrAlreadyExists)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	user, err := scanUser(db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	user, err := scanUser(db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users
		SET email = $1, password_hash = $2, updated_at = $3
		WHERE id = $4`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		user.Email, user.PasswordHash, time.Now(), user.ID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: email already taken", domain.ErrAlreadyExists)
		}
		return fmt.Errorf("update user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
