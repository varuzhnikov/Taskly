package handler

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/todoist/backend/internal/config"
	"github.com/todoist/backend/internal/middleware"
	"github.com/todoist/backend/internal/service"
)

// Deps holds all constructed handlers. It is populated once in main and passed here.
type Deps struct {
	Auth     *AuthHandler
	Users    *UserHandler
	Projects *ProjectHandler
	Tasks    *TaskHandler
	Labels   *LabelHandler
	Comments *CommentHandler

	AuthService *service.AuthService
}

// NewRouter builds and returns the fully configured chi router.
func NewRouter(deps Deps, logger *slog.Logger, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Global middleware — order matters.
	r.Use(cors.Handler(newCORSOptions(cfg)))
	r.Use(middleware.SetRequestID)
	r.Use(middleware.Logger(logger))
	r.Use(chiMiddleware.Recoverer)

	// Swagger UI — accessible without authentication.
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes.
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.Auth.Register)
			r.Post("/login", deps.Auth.Login)
			r.Post("/refresh", deps.Auth.Refresh)
			r.Post("/logout", deps.Auth.Logout)
		})

		// All routes below require a valid JWT.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(deps.AuthService))

			r.Get("/users/me", deps.Users.Me)

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", deps.Projects.List)
				r.Post("/", deps.Projects.Create)
				r.Get("/{id}", deps.Projects.GetByID)
				r.Put("/{id}", deps.Projects.Update)
				r.Delete("/{id}", deps.Projects.Delete)
			})

			r.Route("/labels", func(r chi.Router) {
				r.Get("/", deps.Labels.List)
				r.Post("/", deps.Labels.Create)
				r.Get("/{id}", deps.Labels.GetByID)
				r.Put("/{id}", deps.Labels.Update)
				r.Delete("/{id}", deps.Labels.Delete)
			})

			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", deps.Tasks.List)
				r.Post("/", deps.Tasks.Create)
				r.Get("/{id}", deps.Tasks.GetByID)
				r.Put("/{id}", deps.Tasks.Update)
				r.Delete("/{id}", deps.Tasks.Delete)
				r.Patch("/{id}/complete", deps.Tasks.Complete)
				r.Patch("/{id}/uncomplete", deps.Tasks.Uncomplete)

				r.Route("/{taskId}/comments", func(r chi.Router) {
					r.Get("/", deps.Comments.ListByTask)
					r.Post("/", deps.Comments.Create)
					r.Get("/{commentId}", deps.Comments.GetByID)
					r.Put("/{commentId}", deps.Comments.Update)
					r.Delete("/{commentId}", deps.Comments.Delete)
				})
			})
		})
	})

	return r
}

func newCORSOptions(cfg *config.Config) cors.Options {
	options := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}

	allowedOrigins := trimOrigins(cfg.CORSAllowedOrigins)
	if len(allowedOrigins) > 0 {
		options.AllowedOrigins = allowedOrigins
		return options
	}

	if strings.EqualFold(cfg.Env, "development") {
		options.AllowOriginFunc = func(r *http.Request, origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}

			if u.Scheme != "http" {
				return false
			}

			return u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1"
		}
		return options
	}

	options.AllowedOrigins = []string{"http://localhost:3000"}
	return options
}

func trimOrigins(origins []string) []string {
	trimmed := make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		trimmed = append(trimmed, origin)
	}
	return trimmed
}
