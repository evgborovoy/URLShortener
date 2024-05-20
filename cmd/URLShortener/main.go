package main

import (
	"URLShortener/internal/config"
	"URLShortener/internal/http-server/handlers/redirect"
	"URLShortener/internal/http-server/handlers/url/save"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	config := config.MustLoad()
	logger := setupLogger(config.Env)
	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		slog.Error("failed to init storage", sl.Err(err))
		return
	}
	logger.Info("storage init")
	_ = storage
	err = storage.DeleteURL("go")
	if err != nil {
		slog.Error("failed to init storage", sl.Err(err))
		return
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			config.HttpServer.User: config.HttpServer.Password,
		}))
		r.Post("/", save.New(logger, storage))
	})

	router.Get("/{alias}", redirect.New(logger, storage))

	slog.Info("starting server", slog.String("address", config.Adress))

	svr := &http.Server{
		Addr:         config.Adress,
		Handler:      router,
		ReadTimeout:  config.HttpServer.Timeout,
		WriteTimeout: config.HttpServer.Timeout,
		IdleTimeout:  config.HttpServer.IdleTimeout,
	}
	if err := svr.ListenAndServe(); err != nil {
		slog.Error("failed to start server")
	}
	slog.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
