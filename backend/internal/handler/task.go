package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/todoist/backend/internal/domain"
	"github.com/todoist/backend/internal/middleware"
	"github.com/todoist/backend/internal/service"
)

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

// List godoc
// @Summary      List tasks with optional filters
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        project_id  query  string false "Filter by project UUID"
// @Param        label_id    query  string false "Filter by label UUID"
// @Param        priority    query  int    false "Filter by priority (0-4)"
// @Param        completed   query  bool   false "Filter by completion status"
// @Param        due_before  query  string false "ISO-8601 upper bound for due_at"
// @Param        due_after   query  string false "ISO-8601 lower bound for due_at"
// @Param        parent_id   query  string false "Filter subtasks by parent UUID"
// @Param        search      query  string false "Full-text search on title"
// @Param        limit       query  int    false "Page size (default 50, max 200)"
// @Param        offset      query  int    false "Page offset"
// @Success      200 {array}  taskResponse
// @Failure      401 {object} errorResponse
// @Router       /api/v1/tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	f := parseTaskFilter(r)
	tasks, err := h.tasks.List(r.Context(), userID, f)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	resp := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}
	respond(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Create a task
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body createTaskRequest true "Task details"
// @Success      201 {object} taskResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	t, err := h.tasks.Create(r.Context(), service.CreateTaskInput{
		UserID:       userID,
		ProjectID:    req.ProjectID,
		ParentTaskID: req.ParentTaskID,
		Title:        req.Title,
		Description:  req.Description,
		Links:        req.Links,
		Priority:     domain.Priority(req.Priority),
		DueAt:        req.DueAt,
		SortOrder:    req.SortOrder,
		LabelIDs:     req.LabelIDs,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusCreated, toTaskResponse(t))
}

// GetByID godoc
// @Summary      Get a task by ID
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Task UUID"
// @Success      200 {object} taskResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	t, err := h.tasks.GetByID(r.Context(), id, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toTaskResponse(t))
}

// Update godoc
// @Summary      Replace a task
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string            true "Task UUID"
// @Param        body body updateTaskRequest true "Task fields"
// @Success      200 {object} taskResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{id} [put]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	var req updateTaskRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	t, err := h.tasks.Update(r.Context(), service.UpdateTaskInput{
		ID:           id,
		UserID:       userID,
		ProjectID:    req.ProjectID,
		ParentTaskID: req.ParentTaskID,
		Title:        req.Title,
		Description:  req.Description,
		Links:        req.Links,
		Priority:     domain.Priority(req.Priority),
		DueAt:        req.DueAt,
		SortOrder:    req.SortOrder,
		LabelIDs:     req.LabelIDs,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toTaskResponse(t))
}

// Complete godoc
// @Summary      Mark a task as complete
// @Tags         tasks
// @Security     BearerAuth
// @Param        id path string true "Task UUID"
// @Success      200 {object} taskResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{id}/complete [patch]
func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	t, err := h.tasks.Complete(r.Context(), id, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toTaskResponse(t))
}

// Uncomplete godoc
// @Summary      Mark a task as incomplete
// @Tags         tasks
// @Security     BearerAuth
// @Param        id path string true "Task UUID"
// @Success      200 {object} taskResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{id}/uncomplete [patch]
func (h *TaskHandler) Uncomplete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	t, err := h.tasks.Uncomplete(r.Context(), id, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toTaskResponse(t))
}

// Delete godoc
// @Summary      Delete a task
// @Tags         tasks
// @Security     BearerAuth
// @Param        id path string true "Task UUID"
// @Success      204
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{id} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.tasks.Delete(r.Context(), id, userID); err != nil {
		mapDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DTOs -------------------------------------------------------------------

type createTaskRequest struct {
	ProjectID    *uuid.UUID        `json:"project_id"`
	ParentTaskID *uuid.UUID        `json:"parent_task_id"`
	Title        string            `json:"title"       validate:"required,max=500"`
	Description  string            `json:"description" validate:"max=50000"`
	Links        []domain.TaskLink `json:"links"       validate:"dive"`
	Priority     int8              `json:"priority"    validate:"min=0,max=4"`
	DueAt        *time.Time        `json:"due_at"`
	SortOrder    int               `json:"sort_order"`
	LabelIDs     []uuid.UUID       `json:"label_ids"`
}

type updateTaskRequest struct {
	ProjectID    *uuid.UUID        `json:"project_id"`
	ParentTaskID *uuid.UUID        `json:"parent_task_id"`
	Title        string            `json:"title"       validate:"required,max=500"`
	Description  string            `json:"description" validate:"max=50000"`
	Links        []domain.TaskLink `json:"links"       validate:"dive"`
	Priority     int8              `json:"priority"    validate:"min=0,max=4"`
	DueAt        *time.Time        `json:"due_at"`
	SortOrder    int               `json:"sort_order"`
	LabelIDs     []uuid.UUID       `json:"label_ids"`
}

type taskResponse struct {
	ID              uuid.UUID         `json:"id"`
	ProjectID       *uuid.UUID        `json:"project_id"`
	ParentTaskID    *uuid.UUID        `json:"parent_task_id"`
	Title           string            `json:"title"`
	DescriptionMD   string            `json:"description_md"`
	DescriptionHTML string            `json:"description_html"`
	Links           []domain.TaskLink `json:"links"`
	Priority        domain.Priority   `json:"priority"`
	DueAt           *time.Time        `json:"due_at"`
	CompletedAt     *time.Time        `json:"completed_at"`
	SortOrder       int               `json:"sort_order"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func toTaskResponse(t *domain.Task) taskResponse {
	links := t.Links
	if links == nil {
		links = []domain.TaskLink{}
	}
	return taskResponse{
		ID:              t.ID,
		ProjectID:       t.ProjectID,
		ParentTaskID:    t.ParentTaskID,
		Title:           t.Title,
		DescriptionMD:   t.DescriptionMD,
		DescriptionHTML: t.DescriptionHTML,
		Links:           links,
		Priority:        t.Priority,
		DueAt:           t.DueAt,
		CompletedAt:     t.CompletedAt,
		SortOrder:       t.SortOrder,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

func parseTaskFilter(r *http.Request) domain.TaskFilter {
	q := r.URL.Query()
	f := domain.TaskFilter{
		Search: q.Get("search"),
	}

	if v := q.Get("project_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.ProjectID = &id
		}
	}
	if v := q.Get("label_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.LabelID = &id
		}
	}
	if v := q.Get("parent_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.ParentID = &id
		}
	}
	if v := q.Get("priority"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p := domain.Priority(n)
			if p.Valid() {
				f.Priority = &p
			}
		}
	}
	if v := q.Get("completed"); v != "" {
		b := v == "true"
		f.Completed = &b
	}
	if v := q.Get("due_before"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DueBefore = &t
		}
	}
	if v := q.Get("due_after"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DueAfter = &t
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			f.Limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			f.Offset = n
		}
	}
	return f
}
