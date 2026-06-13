package repository

import (
	"example/todo/internal/domain"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

// абстрактные интерфейсы
// репо списка и элементов

type AuthRepository interface {
	CreateUser(user domain.CreateUser) (int, error)
	GetUser(string, string) (domain.User, error)
	RevokeRefreshByHash(string) error
	RevokeAllRefreshByUserID(int) error
	SaveRefresh(int, string, time.Time) error
	GetUserIDByRefresh(string) (int, error)
}

type Repository struct {
	User AuthRepository
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{User: NewAuthRepository(db, log)}
}
