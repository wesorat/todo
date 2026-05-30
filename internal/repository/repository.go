package repository

import (
	"example/todo/internal/domain"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user domain.CreateUser) (int, error)
}

type Repository struct {
	User UserRepository
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{User: NewUserRepository(db, log)}
}
