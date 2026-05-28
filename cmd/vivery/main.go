package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/krotesq/vivery/internal/account"
	"github.com/krotesq/vivery/internal/auth"
	"github.com/krotesq/vivery/internal/config"
	"github.com/krotesq/vivery/internal/db"
	"github.com/krotesq/vivery/internal/mediamtx"
	"github.com/krotesq/vivery/internal/source"
	"github.com/krotesq/vivery/internal/target"
	"github.com/krotesq/vivery/ui"
)

func main() {

	// load config
	// this config struct should be passed down to all functions that need it
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// background context for database connections
	ctx := context.Background()

	// connect to postgres and redis database
	pool, err := db.ConnectPostgres(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %s", err.Error())
	}
	defer pool.Close()
	log.Println("Connected to postgres")

	rdb, err := db.ConnectRedis(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	log.Println("Connected to redis")
	defer rdb.Close()

	// create main router
	router := chi.NewRouter()

	// enable middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// create api router
	routerApi := chi.NewRouter()

	// the account router has public and protected routes so we
	// mount the account router seperate and configure auth inside
	// we could move the /login, /register, /logout & /reset to the auth package sometime
	routerApi.Mount("/account", account.RoutesWithPool(pool, cfg))

	// protected routes
	routerApi.Group(func(r chi.Router) {
		r.Use(auth.Auth(cfg.JSONWebTokenSecret))
		r.Mount("/source", source.RoutesWithPool(pool))
		r.Mount("/target", target.RoutesWithPool(pool, cfg))
		r.Mount("/mediamtx", mediamtx.RoutesWithPool(pool))
	})

	// mount all routers to main router
	router.Mount("/api", routerApi)

	// embed ui
	sub, err := ui.FS()
	if err != nil {
		log.Fatal(err)
	}

	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := sub.Open(r.URL.Path[1:]) // try to open requested file
			if err != nil {
					// serve index.html if path is not a real file
					http.ServeFileFS(w, r, sub, "index.html")
					return
			}
			http.FileServer(http.FS(sub)).ServeHTTP(w, r)
	}))

	// run server with graceful shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// we run the server in a seperate go routine, so we can handle signals from the os
	go func() {
		log.Printf("Server running at %s:%s", cfg.Host, cfg.Port)
		log.Printf("Allowed origins: %v", cfg.AllowedOrigins)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err.Error())
	}

	log.Println("Server stopped")
}
