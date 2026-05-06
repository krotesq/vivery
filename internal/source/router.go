package source

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RoutesWithPool(pool *pgxpool.Pool) chi.Router {
	router := chi.NewRouter()

	repository := newRepository(pool)
	service := newService(repository)
	handler := newHandler(service)

	router.Get("/", handler.findAll)
	router.Get("/rtmp/{id}", handler.findWithRtmpByID)
	router.Post("/rtmp", handler.createWithRtmp)
	router.Delete("/{id}", handler.deleteByID)

	return router
}
