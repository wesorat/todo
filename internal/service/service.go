package service

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
)

type UserService interface {
	CreateUser(domain.CreateUser) (int, error)
	GetUser(_, _ string) (domain.User, error)
}

type Service struct {
	User UserService
}

func NewService(repo *repository.Repository, log *slog.Logger) *Service {
	return &Service{User: NewUserService(repo.User, log)}
}
