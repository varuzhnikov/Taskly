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

type CommentHandler struct {
	comments *service.CommentService
}

func NewCommentHandler(comments *service.CommentService) *CommentHandler {
	return &CommentHandler{comments: comments}
}

// ListByTask godoc
// @Summary      List comments for a task
// @Tags         comments
// @Security     BearerAuth
// @Produce      json
// @Param        taskId path string true "Task UUID"
// @Success      200 {array}  commentResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{taskId}/comments [get]
func (h *CommentHandler) ListByTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	comments, err := h.comments.ListByTask(r.Context(), taskID, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	resp := make([]commentResponse, len(comments))
	for i, c := range comments {
		resp[i] = toCommentResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Add a comment to a task
// @Tags         comments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        taskId path string             true "Task UUID"
// @Param        body   body createCommentRequest true "Comment body"
// @Success      201 {object} commentResponse
// @Failure      404 {object} errorResponse
// @Failure      422 {object} errorResponse
// @Router       /api/v1/tasks/{taskId}/comments [post]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	var req createCommentRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	c, err := h.comments.Create(r.Context(), service.CreateCommentInput{
		TaskID: taskID,
		UserID: userID,
		Body:   req.Body,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusCreated, toCommentResponse(c))
}

// GetByID godoc
// @Summary      Get a single comment
// @Tags         comments
// @Security     BearerAuth
// @Produce      json
// @Param        taskId    path string true "Task UUID"
// @Param        commentId path string true "Comment UUID"
// @Success      200 {object} commentResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{taskId}/comments/{commentId} [get]
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	commentID, err := uuid.Parse(chi.URLParam(r, "commentId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid comment id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	c, err := h.comments.GetByID(r.Context(), taskID, commentID, userID)
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toCommentResponse(c))
}

// Update godoc
// @Summary      Update a comment
// @Tags         comments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        taskId    path string             true "Task UUID"
// @Param        commentId path string             true "Comment UUID"
// @Param        body      body updateCommentRequest true "Updated body"
// @Success      200 {object} commentResponse
// @Failure      403 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{taskId}/comments/{commentId} [put]
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	commentID, err := uuid.Parse(chi.URLParam(r, "commentId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid comment id")
		return
	}
	var req updateCommentRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	c, err := h.comments.Update(r.Context(), service.UpdateCommentInput{
		TaskID: taskID,
		ID:     commentID,
		UserID: userID,
		Body:   req.Body,
	})
	if err != nil {
		mapDomainError(w, err)
		return
	}
	respond(w, http.StatusOK, toCommentResponse(c))
}

// Delete godoc
// @Summary      Delete a comment
// @Tags         comments
// @Security     BearerAuth
// @Param        taskId    path string true "Task UUID"
// @Param        commentId path string true "Comment UUID"
// @Success      204
// @Failure      403 {object} errorResponse
// @Failure      404 {object} errorResponse
// @Router       /api/v1/tasks/{taskId}/comments/{commentId} [delete]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	commentID, err := uuid.Parse(chi.URLParam(r, "commentId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid comment id")
		return
	}
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.comments.Delete(r.Context(), taskID, commentID, userID); err != nil {
		mapDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DTOs -------------------------------------------------------------------

type commentResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type createCommentRequest struct {
	Body string `json:"body" validate:"required,min=1,max=10000"`
}

type updateCommentRequest struct {
	Body string `json:"body" validate:"required,min=1,max=10000"`
}

func toCommentResponse(c *domain.Comment) commentResponse {
	return commentResponse{
		ID:        c.ID,
		TaskID:    c.TaskID,
		UserID:    c.UserID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
