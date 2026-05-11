package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/repository"
)

type AuthService struct {
	users       *repository.UserRepo
	tokens      *repository.TokenRepo
	txManager   *db.TxManager
	jwtSecret   []byte
	accessExp   time.Duration
	refreshExp  time.Duration
	bcryptCost  int
}

func NewAuthService(
	users *repository.UserRepo,
	tokens *repository.TokenRepo,
	txManager *db.TxManager,
	jwtSecret string,
	accessExp, refreshExp time.Duration,
	bcryptCost int,
) *AuthService {
	return &AuthService{
		users:      users,
		tokens:     tokens,
		txManager:  txManager,
		jwtSecret:  []byte(jwtSecret),
		accessExp:  accessExp,
		refreshExp: refreshExp,
		bcryptCost: bcryptCost,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // access token TTL in seconds
	UserID       string
	UserEmail    string
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
	// Optional client hints stored with the refresh token for audit.
	UserAgent string
	IP        string
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        in.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	var pair *TokenPair
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.users.Create(ctx, user); err != nil {
			return err
		}
		p, err := s.issueTokenPair(ctx, user.ID, "", "")
		if err != nil {
			return err
		}
		pair = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	pair.UserID = user.ID.String()
	pair.UserEmail = user.Email
	return pair, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*TokenPair, error) {
	user, err := s.users.GetByEmail(ctx, in.Email)
	if err != nil {
		// Return the same error shape whether email or password is wrong
		// to prevent user-enumeration attacks.
		return nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	pair, err := s.issueTokenPair(ctx, user.ID, in.UserAgent, in.IP)
	if err != nil {
		return nil, err
	}
	pair.UserID = user.ID.String()
	pair.UserEmail = user.Email
	return pair, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*TokenPair, error) {
	hash := hashToken(rawToken)

	stored, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}
	if time.Now().After(stored.ExpiresAt) {
		_ = s.tokens.Delete(ctx, stored.ID)
		return nil, domain.ErrTokenExpired
	}

	// Rotate: delete the used token before issuing a new pair.
	if err := s.tokens.Delete(ctx, stored.ID); err != nil {
		return nil, fmt.Errorf("rotate refresh token: %w", err)
	}

	pair, err := s.issueTokenPair(ctx, stored.UserID, stored.UserAgent, stored.IP)
	if err != nil {
		return nil, err
	}
	user, err := s.users.GetByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}
	pair.UserID = user.ID.String()
	pair.UserEmail = user.Email
	return pair, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	hash := hashToken(rawToken)
	stored, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		// Token not found — treat as already logged out.
		return nil
	}
	return s.tokens.Delete(ctx, stored.ID)
}

// issueTokenPair mints a new JWT access token and a refresh token, persisting
// the refresh token hash to the database.
func (s *AuthService) issueTokenPair(ctx context.Context, userID uuid.UUID, userAgent, ip string) (*TokenPair, error) {
	accessToken, err := s.mintAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("mint access token: %w", err)
	}

	rawRefresh, err := generateOpaqueToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.refreshExp),
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IP:        ip,
	}
	if err := s.tokens.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("persist refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int(s.accessExp.Seconds()),
	}, nil
}

func (s *AuthService) mintAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(s.accessExp).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateAccessToken parses and validates a JWT, returning the subject (user ID).
func (s *AuthService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return uuid.Nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, domain.ErrTokenInvalid
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, domain.ErrTokenInvalid
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, domain.ErrTokenInvalid
	}
	return id, nil
}

// generateOpaqueToken returns a cryptographically random 32-byte hex string.
func generateOpaqueToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken returns the SHA-256 hex digest of a raw token value.
// Only the hash is stored in the database; the raw value goes to the client.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
