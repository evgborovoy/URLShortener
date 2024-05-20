package main

import (
	"URLShortener/internal/config"
	"URLShortener/internal/http-server/handlers/url/save"
	"URLShortener/internal/lib/logger/sl"
	"URLShortener/internal/storage/sqlite"
	"log/slog"
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
	router.Post("/url", save.New(logger, storage))

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
