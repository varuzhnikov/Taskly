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

type ProjectHandler struct {
	projects *service.ProjectService
}

func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

// List godoc
// @Summary      List all projects for the authenticated user
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array}  projectResponse
// @Failure      401 {object} errorResponse
// @Router       /api/v1/projects [get]
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	projects, err := h.projects.List(r.Context(), userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	resp := make([]projectResponse, len(projects))
	for i, p := range projects {
		resp[i] = toProjectResponse(p)
	}
	respond(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Create a new project
// @Tags         projects
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body createProjectRequest true "Project details"
// @Success      201 {object} projectResponse
// @Failure      400 {object} errorResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	p, err := h.projects.Create(r.Context(), service.CreateProjectInput{
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusCreated, toProjectResponse(p))
}

// GetByID godoc
// @Summary      Get a project by ID
// @Tags         projects
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Project UUID"
// @Success      200 {object} projectResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	p, err := h.projects.GetByID(r.Context(), id, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toProjectResponse(p))
}

// Update godoc
// @Summary      Update a project
// @Tags         projects
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string              true "Project UUID"
// @Param        body body updateProjectRequest true "Updated fields"
// @Success      200 {object} projectResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/projects/{id} [put]
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project id")
		return
	}
	var req updateProjectRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	p, err := h.projects.Update(r.Context(), service.UpdateProjectInput{
		ID:        id,
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toProjectResponse(p))
}

// Delete godoc
// @Summary      Delete a project
// @Tags         projects
// @Security     BearerAuth
// @Param        id path string true "Project UUID"
// @Success      204
// @Failure      404 {object} errorResponse
// @Router       /api/v1/projects/{id} [delete]
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid project id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.projects.Delete(r.Context(), id, userID); err != nil {
		mapDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DTOs -------------------------------------------------------------------

type projectResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type createProjectRequest struct {
	Name      string `json:"name"       validate:"required,max=255"`
	Color     string `json:"color"      validate:"omitempty,max=20"`
	SortOrder int    `json:"sort_order"`
}

type updateProjectRequest struct {
	Name      string `json:"name"       validate:"required,max=255"`
	Color     string `json:"color"      validate:"omitempty,max=20"`
	SortOrder int    `json:"sort_order"`
}

func toProjectResponse(p *domain.Project) projectResponse {
	return projectResponse{
		ID:        p.ID,
		Name:      p.Name,
		Color:     p.Color,
		SortOrder: p.SortOrder,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
