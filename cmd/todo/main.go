package main

import (
	"example/todo/pkg/config"
	"example/todo/pkg/database"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	log := setupLogger()

	log.Info("Starting todo service", slog.Any("cfg", cfg))

	_, err := database.New(cfg.Database)
	if err != nil {
		log.Error("Failed connect to db", slog.Any("err", err))
		return
	}

	log.Info("Ending todo service")

}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
