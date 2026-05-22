package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
)

// TODO: rename to InsertUser or MustInsertUser to make the DB write explicit.
// CreateUser inserts a user with the given email and plaintext password,
// returning the created domain.User. The insert runs inside whatever
// transaction is embedded in ctx.
func CreateUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, email, password string) *domain.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4) // TODO: replace magic number with bcrypt.MinCost; low cost for tests
	// TODO: make helper flow more explicit by separating entity creation,
	// error capture, and require.NoError verification.
	require.NoError(t, err)

	u := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// TODO: make helper flow more explicit by storing the Create error first,
	// then verifying it with require.NoError.
	require.NoError(t, repository.NewUserRepo(pool).Create(ctx, u))
	return u
}

// TODO: rename to InsertProject or MustInsertProject to make the DB write explicit.
// CreateProject inserts a project owned by userID and returns it.
func CreateProject(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, name string) *domain.Project {
	t.Helper()
	p := &domain.Project{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Color:     "#cccccc",
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// TODO: make helper flow more explicit by storing the Create error first,
	// then verifying it with require.NoError.
	require.NoError(t, repository.NewProjectRepo(pool).Create(ctx, p))
	return p
}

// TODO: rename to InsertTask or MustInsertTask to make the DB write explicit.
// CreateTask inserts a minimal task and returns it.
func CreateTask(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, title string) *domain.Task {
	t.Helper()
	task := &domain.Task{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Links:     []domain.TaskLink{},
		Priority:  domain.PriorityNone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// TODO: make helper flow more explicit by storing the Create error first,
	// then verifying it with require.NoError.
	require.NoError(t, repository.NewTaskRepo(pool).Create(ctx, task))
	return task
}

// TODO: rename to InsertLabel or MustInsertLabel to make the DB write explicit.
// CreateLabel inserts a label owned by userID and returns it.
func CreateLabel(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, name string) *domain.Label {
	t.Helper()
	l := &domain.Label{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Color:     "#aabbcc",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// TODO: make helper flow more explicit by storing the Create error first,
	// then verifying it with require.NoError.
	require.NoError(t, repository.NewLabelRepo(pool).Create(ctx, l))
	return l
}
