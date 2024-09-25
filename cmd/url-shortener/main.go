package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"usekit-go/internal/config"
	"usekit-go/internal/http-server/handlers/url/save"
	"usekit-go/internal/lib/logger/sl"
	"usekit-go/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger: sl import "log/sl"
	logger := setLogger(cfg.Env)

	logger.Info("Starting URL shortener", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	// ---------------
	//getUrl, err := storage.GetUrl("youtube")
	//if err != nil {
	//	logger.Error("failed to get url", sl.Err(err))
	//	os.Exit(1)
	//}
	//logger.Info("getUrl", slog.String("getUrl", getUrl))
	// ------------------
	//id, err := storage.SaveURL("https://youtube.com", "youtube")
	//if err != nil {
	//	logger.Error("failed to save url", sl.Err(err))
	//	os.Exit(1)
	//}
	//
	//logger.Info("saved url", slog.Int64("id", id))
	// ------------------
	//deletedId, err := storage.DeleteUrl("5")
	//if err != nil {
	//	logger.Error("failed to delete url", sl.Err(err))
	//	os.Exit(1)
	//}
	//logger.Info("deleted element with id", slog.String("deletedId", deletedId))

	// TODO: init router: chi, "chi, render"
	router := chi.NewRouter()

	// mv for adding requestId to every request
	router.Use(middleware.RequestID)
	// mv for logging
	// TODO: write custom logger
	router.Use(middleware.Logger)
	// mw for recover app after panic in a specific function
	router.Use(middleware.Recoverer)
	// mv for comfortable using url addresses
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// TODO: run server
}

func setLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
