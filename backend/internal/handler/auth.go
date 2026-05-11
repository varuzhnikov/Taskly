package handler

import (
	"net/http"

	"github.com/todoist/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register godoc
// @Summary      Register a new account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body registerRequest true "Credentials"
// @Success      201 {object} tokenResponse
// @Failure      400 {object} errorResponse
// @Failure      409 {object} errorResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	pair, err := h.auth.Register(r.Context(), service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respond(w, http.StatusCreated, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		User:         tokenUser{ID: pair.UserID, Email: pair.UserEmail},
	})
}

// Login godoc
// @Summary      Authenticate and obtain tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body loginRequest true "Credentials"
// @Success      200 {object} tokenResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	pair, err := h.auth.Login(r.Context(), service.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
		IP:        r.RemoteAddr,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respond(w, http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		User:         tokenUser{ID: pair.UserID, Email: pair.UserEmail},
	})
}

// Refresh godoc
// @Summary      Rotate a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body refreshRequest true "Current refresh token"
// @Success      200 {object} tokenResponse
// @Failure      400 {object} errorResponse
// @Failure      401 {object} errorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	pair, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		mapDomainError(w, err)
		return
	}

	respond(w, http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		User:         tokenUser{ID: pair.UserID, Email: pair.UserEmail},
	})
}

// Logout godoc
// @Summary      Revoke a refresh token
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body refreshRequest true "Token to revoke"
// @Success      204
// @Failure      400 {object} errorResponse
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.auth.Logout(r.Context(), req.RefreshToken); err != nil {
		mapDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- request / response DTOs ------------------------------------------------

type registerRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type tokenUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type tokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         tokenUser `json:"user"`
}

type errorResponse struct {
	Error string `json:"error"`
}
