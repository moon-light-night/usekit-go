package main

import (
	"fmt"
	"log/slog"
	"os"
	"usekit-go/internal/config"
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

	// TODO: init router: chi, "chi, render"

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
