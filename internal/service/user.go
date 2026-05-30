package service

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
)

type userService struct {
	repo repository.UserRepository
	log  *slog.Logger
}

func NewUserService(repo repository.UserRepository, log *slog.Logger) *userService {
	return &userService{repo: repo, log: log}
}

func (s *userService) CreateUser(user domain.CreateUser) (int, error) {
	password_hash, err := generatePasswordHash(user.Password)
	if err != nil {
		s.log.Error("The password cannot be empty", slog.Any("err", err))
		return 0, err
	}
	user.Password = password_hash
	id, err := s.repo.CreateUser(user)
	if err != nil {
		s.log.Error("Cannot create user", slog.Any("err", err))
		return 0, err
	}

	return id, nil
}

// TODO дописать хэширование
func generatePasswordHash(password string) (string, error) {
	return password + "_hash", nil
}
