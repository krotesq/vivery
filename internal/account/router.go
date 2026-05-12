package account

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krotesq/vivery/internal/auth"
	"github.com/krotesq/vivery/internal/config"
)

func RoutesWithPool(pool *pgxpool.Pool, cfg *config.Config) chi.Router {

	router := chi.NewRouter()

	repository := newRepository(pool)
	service := newService(repository, cfg)
	handler := newHandler(service)

	// public — rate limited
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/login", handler.login)
		r.Post("/logout", handler.logout)
		r.Post("/refresh", handler.refresh)
	})

	// protected
	router.Group(func(r chi.Router) {
		r.Use(auth.Auth(cfg.JSONWebTokenSecret))
		r.Get("/me", handler.me)
		r.Post("/", handler.create)
		r.Get("/{id}", handler.findByID)
		r.Patch("/{id}/deactivate", handler.deactivateByID)
		r.Patch("/{id}/activate", handler.activateByID)
		r.Delete("/{id}", handler.deleteByID)
	})

	return router
}
