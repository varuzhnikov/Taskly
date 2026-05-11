package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/middleware"
	"github.com/todoist/backend/internal/service"
)

type LabelHandler struct {
	labels *service.LabelService
}

func NewLabelHandler(labels *service.LabelService) *LabelHandler {
	return &LabelHandler{labels: labels}
}

// List godoc
// @Summary      List all labels for the authenticated user
// @Tags         labels
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array}  labelResponse
// @Failure      401 {object} errorResponse
// @Router       /api/v1/labels [get]
func (h *LabelHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	labels, err := h.labels.List(r.Context(), userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	resp := make([]labelResponse, len(labels))
	for i, l := range labels {
		resp[i] = toLabelResponse(l)
	}
	respond(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Create a label
// @Tags         labels
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body createLabelRequest true "Label details"
// @Success      201 {object} labelResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/labels [post]
func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createLabelRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	l, err := h.labels.Create(r.Context(), service.CreateLabelInput{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusCreated, toLabelResponse(l))
}

// GetByID godoc
// @Summary      Get a label by ID
// @Tags         labels
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Label UUID"
// @Success      200 {object} labelResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/labels/{id} [get]
func (h *LabelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid label id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	l, err := h.labels.GetByID(r.Context(), id, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toLabelResponse(l))
}

// Update godoc
// @Summary      Update a label
// @Tags         labels
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string            true "Label UUID"
// @Param        body body updateLabelRequest true "Updated fields"
// @Success      200 {object} labelResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/labels/{id} [put]
func (h *LabelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid label id")
		return
	}
	var req updateLabelRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	l, err := h.labels.Update(r.Context(), service.UpdateLabelInput{
		ID:     id,
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toLabelResponse(l))
}

// Delete godoc
// @Summary      Delete a label
// @Tags         labels
// @Security     BearerAuth
// @Param        id path string true "Label UUID"
// @Success      204
// @Failure      404 {object} errorResponse
// @Router       /api/v1/labels/{id} [delete]
func (h *LabelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid label id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.labels.Delete(r.Context(), id, userID); err != nil {
		mapDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DTOs -------------------------------------------------------------------

type labelResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type createLabelRequest struct {
	Name  string `json:"name"  validate:"required,max=100"`
	Color string `json:"color" validate:"omitempty,max=20"`
}

type updateLabelRequest struct {
	Name  string `json:"name"  validate:"required,max=100"`
	Color string `json:"color" validate:"omitempty,max=20"`
}

func toLabelResponse(l *domain.Label) labelResponse {
	return labelResponse{
		ID:        l.ID,
		Name:      l.Name,
		Color:     l.Color,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}
