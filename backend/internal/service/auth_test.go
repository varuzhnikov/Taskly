package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
	"github.com/todoist/backend/internal/service"
	"github.com/todoist/backend/testutil"
)

// newTestAuthService wires a real AuthService backed by a fresh Postgres container.
func newTestAuthService(t *testing.T) *service.AuthService {
	t.Helper()
	pool := testutil.NewPostgresPool(t)
	return service.NewAuthService(
		repository.NewUserRepo(pool),
		repository.NewTokenRepo(pool),
		db.NewTxManager(pool),
		"super-secret-jwt-key-for-test-use-only",
		15*time.Minute,
		7*24*time.Hour,
		4,
	)
}

func TestAuthService_Register(t *testing.T) {
	svc := newTestAuthService(t)

	t.Run("returns token pair on success", func(t *testing.T) {
		pair, err := svc.Register(context.Background(), service.RegisterInput{
			Email:    "reg@example.com",
			Password: "strongpassword",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, pair.AccessToken)
		assert.NotEmpty(t, pair.RefreshToken)
		assert.Greater(t, pair.ExpiresIn, 0)
	})

	t.Run("duplicate email returns ErrAlreadyExists", func(t *testing.T) {
		_, err := svc.Register(context.Background(), service.RegisterInput{
			Email:    "dup@example.com",
			Password: "password1",
		})
		require.NoError(t, err)

		_, err = svc.Register(context.Background(), service.RegisterInput{
			Email:    "dup@example.com",
			Password: "password2",
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestAuthService_Login(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, service.RegisterInput{
		Email:    "login@example.com",
		Password: "correcthorse",
	})
	require.NoError(t, err)

	t.Run("correct credentials return token pair", func(t *testing.T) {
		pair, err := svc.Login(ctx, service.LoginInput{
			Email:    "login@example.com",
			Password: "correcthorse",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, pair.AccessToken)
	})

	t.Run("wrong password returns ErrUnauthorized", func(t *testing.T) {
		_, err := svc.Login(ctx, service.LoginInput{
			Email:    "login@example.com",
			Password: "wrongpassword",
		})
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("unknown email returns ErrUnauthorized to prevent user enumeration", func(t *testing.T) {
		_, err := svc.Login(ctx, service.LoginInput{
			Email:    "nobody@example.com",
			Password: "anything",
		})
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

func TestAuthService_Refresh(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	pair, err := svc.Register(ctx, service.RegisterInput{
		Email:    "refresh@example.com",
		Password: "password",
	})
	require.NoError(t, err)

	t.Run("valid refresh token issues new pair", func(t *testing.T) {
		newPair, err := svc.Refresh(ctx, pair.RefreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newPair.AccessToken)
	})

	t.Run("used refresh token is revoked (rotation)", func(t *testing.T) {
		// pair.RefreshToken was consumed by the previous subtest's Refresh call.
		_, err := svc.Refresh(ctx, pair.RefreshToken)
		assert.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("garbage token returns ErrTokenInvalid", func(t *testing.T) {
		_, err := svc.Refresh(ctx, "not-a-real-token")
		assert.ErrorIs(t, err, domain.ErrTokenInvalid)
	})
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	svc := newTestAuthService(t)

	pair, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "validate@example.com",
		Password: "password",
	})
	require.NoError(t, err)

	t.Run("valid token parses correctly", func(t *testing.T) {
		id, err := svc.ValidateAccessToken(pair.AccessToken)
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("tampered signature returns ErrTokenInvalid", func(t *testing.T) {
		_, err := svc.ValidateAccessToken(pair.AccessToken + "x")
		assert.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("garbage string returns ErrTokenInvalid", func(t *testing.T) {
		_, err := svc.ValidateAccessToken("not.a.jwt")
		assert.ErrorIs(t, err, domain.ErrTokenInvalid)
	})
}
