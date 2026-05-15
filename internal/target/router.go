package target

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krotesq/vivery/internal/config"
)

func RoutesWithPool(pool *pgxpool.Pool, c *config.Config) chi.Router {
	router := chi.NewRouter()

	repository := newRepository(pool)
	service := newService(repository, c)
	handler := newHandler(service)

	router.Get("/", handler.findAll)
	router.Get("/rtmp/{id}", handler.findWithRtmpByID)
	router.Post("/rtmp", handler.createWithRtmp)
	router.Delete("/{id}", handler.deleteByID)

	return router
}
