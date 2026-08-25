package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	signingKey, refreshPapper, err := loadSecrets()
	if err != nil {
		log.Error("Failed to load secrets", slog.Any("err", err))
		return
	}

	log.Info("Starting todo service", slog.Any("cfg", cfg))



	db, err := database.New(cfg.Database)
	if err != nil {
		log.Error("Failed connect to db", slog.Any("err", err))
		return
	}

	repo := repository.NewRepository(db, log)
	service := service.NewService(repo, log, signingKey, refreshPapper)
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


func loadSecrets() (string, string, error) {
	missing := []string{}
	signingKey := os.Getenv("signingKey")
	refreshPepper := os.Getenv("refreshPepper")
	dbPassword := os.Getenv("DB_PASSWORD")

	if signingKey == "" {
		missing = append(missing, signingKey)
	}
	if refreshPepper == "" {
		missing = append(missing, refreshPepper)
	}
	if dbPassword== "" {
		missing = append(missing, dbPassword)
	}
	if len(missing) != 0 {
		return "", "", fmt.Errorf("not set required field, %v", strings.Join(missing, ", "))
	}
	return signingKey, refreshPepper, nil
}