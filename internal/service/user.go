package service

import (
	"errors"
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
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

func (s *userService) GetUser(name, password string) (domain.User, error) {
	user, err := s.repo.GetUser(name, password)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, repository.ErrUserNotFound
		}
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidPassword
	}
	return user, nil

}

func generatePasswordHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
