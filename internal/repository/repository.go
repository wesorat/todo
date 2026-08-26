package repository

import (
	"log/slog"
	"time"

	"github.com/wesorat/todo/internal/domain"

	"github.com/jmoiron/sqlx"
)

// абстрактные интерфейсы
// репо списка и элементов

type AuthRepository interface {
	CreateUser(user domain.CreateUser) (int, error)
	GetUser(name string) (domain.User, error)
	RevokeRefreshByHash(refresh_hash string) error
	RevokeAllRefreshByUserID(user_id int) error
	SaveRefresh(user_id int, refresh_hash string, expired_at time.Time) error
	GetUserIDByRefresh(refresh_hash string) (int, error)
}

type ListsRepository interface {
	Create(list domain.CreateList) (int, error)
	Get(user_id, list_id int) (domain.List, error)
	GetAll(user_id int) ([]domain.List, error)
	Update(user_id, list_id int, title, description *string) error
	Delete(user_id, list_id int) error
}

type ItemsRepository interface {
	Create(user_id int, item domain.CreateItem) (int, error)
	Get(user_id, item_id int) (domain.Item, error)
	GetAll(user_id, list_id int) ([]domain.Item, error)
	Update(user_id, item_id int, title, description *string, done *bool) error
	Delete(user_id, item_id int) error
}

type Repository struct {
	Auth AuthRepository
	List ListsRepository
	Item ItemsRepository
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Auth: NewAuthRepository(db, log),
		List: NewListsRepository(db, log),
		Item: NewItemsRepository(db, log),
	}
}
