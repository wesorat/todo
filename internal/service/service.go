package service

import (
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
)

// сервисы списков и элементов
// абстрактные интерфейсы

type Tokens struct {
	JWT          string
	RefreshToken string
}

type AuthService interface {
	CreateUser(domain.CreateUser) (int, error)
	GetUser(_, _ string) (domain.User, error)
	SignIn(_, _ string) (Tokens, error)
	ParseJWT(string) (int, error)
	Logout(string) error
	LogoutAll(string) error
	RenewalJWT(string) (string, error)
}

type ListService interface {
	Create(list domain.CreateList) (int, error)
	Get(user_id, list_id int) (domain.List, error)
	GetAll(user_id int) ([]domain.List, error)
	Update(user_id, list_id int, title, description *string) error
	Delete(user_id, list_id int) error
}

type Service struct {
	Auth  AuthService
	Lists ListService
}

func NewService(repo *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Auth:  NewAuthService(repo.Auth, log),
		Lists: NewListService(repo.Lists, log),
	}
}
