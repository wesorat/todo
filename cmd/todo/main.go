package main

import (
	"log/slog"
	"os"

	"github.com/wesorat/todo/internal/handler"
	"github.com/wesorat/todo/internal/repository"
	"github.com/wesorat/todo/internal/server"
	"github.com/wesorat/todo/internal/service"
	"github.com/wesorat/todo/pkg/config"
	"github.com/wesorat/todo/pkg/database"
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
	if err := srv.Run(cfg.HTTPServer.Port, handlers.InitRoutes()); err != nil {
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
