package main

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
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
	service.User.CreateUser(domain.CreateUser{
		Name:     "John123",
		Password: "pass",
	})
	// repo.User.CreateUser(domain.CreateUser{
	// 	Name:     "John",
	// 	Password: "pass",
	// })
	_ = repo

	log.Info("Ending todo service")

}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
