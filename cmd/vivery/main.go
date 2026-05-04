package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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
)

func main() {

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// connect to database
	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err.Error())
	}
	defer pool.Close()
	log.Println("connected to database")

	// create main router
	router := chi.NewRouter()

	// enable middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
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
	routerApi.Mount("/account", account.RoutesWithPool(pool, cfg.JSONWebTokenSecret, cfg.JSONWebTokenIssuer, cfg.JSONWebTokenExpireMinutes, cfg.RefreshTokenExpireDays, cfg.BcryptCost))

	// protected routes
	routerApi.Group(func(r chi.Router) {
		r.Use(auth.Auth(cfg.JSONWebTokenSecret))
		r.Mount("/source", source.RoutesWithPool(pool))
		r.Mount("/target", target.RoutesWithPool(pool))
		r.Mount("/mediamtx", mediamtx.RoutesWithPool(pool))
	})

	// create web router
	routerWeb := chi.NewRouter()

	// add fs to web router
	webDir := filepath.Join(".", "web")
	fileServer := http.StripPrefix("/", http.FileServer(http.Dir(webDir)))
	routerWeb.Handle("/*", fileServer)

	// mount all routers to main router
	router.Mount("/api", routerApi)
	router.Mount("/", routerWeb)

	// run server with graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Server running at %s:%s", cfg.Host, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error: %s", err.Error())
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err.Error())
	}

	log.Println("server stopped")
}
