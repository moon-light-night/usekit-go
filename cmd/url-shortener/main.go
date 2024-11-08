package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	authgrpc "usekit-go/internal/clients/auth/grpc"
	"usekit-go/internal/config"
	deleteClient "usekit-go/internal/http-server/handlers/url/delete"
	"usekit-go/internal/http-server/handlers/url/redirect"
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

	authClient, err := authgrpc.New(
		context.Background(),
		logger,
		//cfg.Clients.AUTH.Address,
		"localhost:8082",
		//cfg.Clients.AUTH.Timeout,
		4,
		//cfg.Clients.AUTH.RetriesCount,
		5,
	)
	if err != nil {
		logger.Error("failed to init auth client", sl.Err(err))
		os.Exit(1)
	}

	resp, err := authClient.IsAdmin(context.Background(), 2)
	if err != nil {
		logger.Error("failed to check admin", sl.Err(err))
	}
	logger.Info("resp", resp)

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

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("/url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(logger, storage))
		r.Delete("/{alias}", deleteClient.New(logger, storage))
	})
	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	// TODO: run server
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start server", sl.Err(err))
	}

	logger.Info("shutting down server")
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
