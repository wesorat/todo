package main

import (
	"example/todo/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	log := setupLogger()

	log.Info("Starting todo service", slog.Any("cfg", cfg))

}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
