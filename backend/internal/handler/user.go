package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/todoist/backend/internal/middleware"
	"github.com/todoist/backend/internal/repository"
)

type UserHandler struct {
	users *repository.UserRepo
}

func NewUserHandler(users *repository.UserRepo) *UserHandler {
	return &UserHandler{users: users}
}

// Me godoc
// @Summary      Get current user profile
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} userResponse
// @Failure      401 {object} errorResponse
// @Router       /api/v1/users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toUserResponse(user.ID, user.Email, user.CreatedAt))
}

// --- DTOs -------------------------------------------------------------------

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(id uuid.UUID, email string, createdAt time.Time) userResponse {
	return userResponse{ID: id, Email: email, CreatedAt: createdAt}
}
