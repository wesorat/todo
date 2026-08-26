package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/wesorat/todo/internal/handler"
	"github.com/wesorat/todo/internal/repository"
	"github.com/wesorat/todo/internal/server"
	"github.com/wesorat/todo/internal/service"
	"github.com/wesorat/todo/pkg/config"
	"github.com/wesorat/todo/pkg/database"
	"github.com/wesorat/todo/pkg/redis"
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
	} else {
		log.Info("Connected to db")
	}
	defer db.Close()

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Error("Failed connect to redis", slog.Any("err", err))
		return
	} else {
		log.Info("Connected to redis")
	}
	defer rdb.Close()

	repo := repository.NewRepository(db, log)
	service := service.NewService(repo, log, signingKey, refreshPapper)
	handlers := handler.NewHandler(service, log)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPServer.Port, handlers.InitRoutes()); err != nil {
			log.Error("Error when runnning the http server", slog.Any("err", err))
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	log.Info("Shutting down service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown server error", slog.Any("err", err))
	} else {
		log.Info("db closed")
	}
	if err := db.Close(); err != nil {
		log.Error("shutdown db error", slog.Any("err", err))
	}
	if err := rdb.Close(); err != nil {
		log.Error("shutdown redis error", slog.Any("err", err))
	} else {
		log.Info("redis closed")
	}
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
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if signingKey == "" {
		missing = append(missing, signingKey)
	}
	if refreshPepper == "" {
		missing = append(missing, refreshPepper)
	}
	if dbPassword == "" {
		missing = append(missing, dbPassword)
	}
	if redisPassword == "" {
		missing = append(missing, redisPassword)
	}
	if len(missing) != 0 {
		return "", "", fmt.Errorf("not set required field, %v", strings.Join(missing, ", "))
	}
	return signingKey, refreshPepper, nil
}
