package account

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krotesq/vivery/internal/auth"
	"github.com/krotesq/vivery/internal/config"
)

func RoutesWithPool(p *pgxpool.Pool, c *config.Config) chi.Router {

	router := chi.NewRouter()

	r := newRepository(p)
	s := newService(r, c)
	h := newHandler(s)

	// public
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(50, time.Minute))
		r.Post("/logout", h.logout)
	})

	// public
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/login", h.login)
		r.Post("/refresh", h.refresh)
	})

	// protected
	router.Group(func(r chi.Router) {
		r.Use(auth.Auth(c.JSONWebTokenSecret))
		r.Get("/me", h.me)
		r.Post("/", h.create)
		r.Get("/{id}", h.findByID)
		r.Patch("/{id}/deactivate", h.deactivateByID)
		r.Patch("/{id}/activate", h.activateByID)
		r.Delete("/{id}", h.deleteByID)
	})

	return router
}
