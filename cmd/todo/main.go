package main

import (
	"example/todo/internal/handler"
	"example/todo/internal/repository"
	"example/todo/internal/server"
	"example/todo/internal/service"
	"example/todo/pkg/config"
	"example/todo/pkg/database"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	log := setupLogger()

	log.Info("Starting todo service", slog.Any("cfg", cfg))

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Error("Failed connect to db", slog.Any("err", err))
		return
	}

	repo := repository.NewRepository(db, log)
	service := service.NewService(repo, log)
	handlers := handler.NewHandler(service, log)
	srv := new(server.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Error("Error when runnning the http server", slog.Any("err", err))
		return
	}

	_ = service
	_ = repo

	log.Info("Ending todo service")

}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
