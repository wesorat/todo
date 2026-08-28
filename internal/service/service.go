package service

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/wesorat/todo/internal/domain"
	"github.com/wesorat/todo/internal/repository"
)

type Tokens struct {
	JWT          string
	RefreshToken string
}

type AuthService interface {
	CreateUser(domain.CreateUser) (int, error)
	GetUser(_, _ string) (domain.User, error)
	SignIn(ctx context.Context, name, password string) (Tokens, error)
	ParseJWT(string) (int, error)
	Logout(ctx context.Context, refresh string) error
	LogoutAll(ctx context.Context, refresh string) error
	RenewalJWT(ctx context.Context, refresh string) (string, error)
}

type ListService interface {
	Create(list domain.CreateList) (int, error)
	Get(user_id, list_id int) (domain.List, error)
	GetAll(user_id int) ([]domain.List, error)
	Update(user_id, list_id int, title, description *string) error
	Delete(user_id, list_id int) error
}

type ItemService interface {
	Create(user_id int, item domain.CreateItem) (int, error)
	Get(user_id, item_id int) (domain.Item, error)
	GetAll(user_id, list_id int) ([]domain.Item, error)
	Update(user_id, item_id int, title, description *string, done *bool) error
	Delete(user_id, item_id int) error
}

type Service struct {
	Auth AuthService
	List ListService
	Item ItemService
}

func NewService(repo *repository.Repository, log *slog.Logger, redis *redis.Client, signingKey, refreshPapper string) *Service {
	return &Service{
		Auth: NewAuthService(repo.Auth, log, redis, signingKey, refreshPapper),
		List: NewListService(repo.List, log),
		Item: NewItemService(repo.Item, log),
	}
}
