package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/todoist/backend/internal/domain"
)

var validate = validator.New()

// respond serialises v as JSON with the given status code.
func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", slog.String("err", err.Error()))
	}
}

// respondError writes a machine-readable {"error": "..."} JSON body.
func respondError(w http.ResponseWriter, status int, msg string) {
	respond(w, status, map[string]string{"error": msg})
}

// decodeJSON reads and validates the request body into dst.
// Returns false and writes an appropriate error response if decoding or
// validation fails — the caller should return immediately.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return false
	}
	if err := validate.Struct(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			respondError(w, http.StatusUnprocessableEntity, formatValidationErrors(ve))
			return false
		}
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return false
	}
	return true
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "validation failed"
	}
	// Return the first error in a human-readable format.
	e := errs[0]
	return e.Field() + ": failed " + e.Tag() + " validation"
}

// mapDomainError translates domain sentinel errors into HTTP status codes.
func mapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		respondError(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		respondError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrForbidden):
		respondError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrTokenExpired):
		respondError(w, http.StatusUnauthorized, "token expired")
	case errors.Is(err, domain.ErrTokenInvalid):
		respondError(w, http.StatusUnauthorized, "token invalid")
	default:
		slog.Error("unhandled error", slog.String("err", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
