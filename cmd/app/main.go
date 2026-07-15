package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"shortLink/internal/config"
	"shortLink/internal/handler"
	"shortLink/internal/service"
	"shortLink/internal/storage"
	"shortLink/internal/storage/memory"
	"shortLink/internal/storage/postgres"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	ctx := context.Background()
	var store storage.Storage

	switch cfg.StorageType {
	case "memory":
		store = memory.New()
		logger.Info("using in-memory storage")
	case "postgres":
		store, err = postgres.New(ctx, cfg.DatabaseDSN)
		logger.Info("using postgres storage")
	}
	defer store.Close()

	svc := service.NewService(store)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	go func() {
		logger.Info("Server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	shutDownCtx, shutDownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutDownRelease()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	logger.Info("Server exited Properly")
}
