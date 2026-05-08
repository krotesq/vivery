package account

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krotesq/vivery/internal/auth"
)

func RoutesWithPool(pool *pgxpool.Pool, jwtSecret string, jwtIssuer string, jwtExpMin int, refreshExpDays int, bcryptCost int) chi.Router {

	router := chi.NewRouter()

	repository := newRepository(pool)
	service := newService(repository, jwtSecret, jwtIssuer, jwtExpMin, refreshExpDays, bcryptCost)
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
		r.Use(auth.Auth(jwtSecret))
		r.Post("/", handler.create)
		r.Get("/{id}", handler.findByID)
		r.Patch("/{id}/deactivate", handler.deactivateByID)
		r.Patch("/{id}/activate", handler.activateByID)
		r.Delete("/{id}", handler.deleteByID)
		r.Get("/me", handler.me)
	})

	return router
}
