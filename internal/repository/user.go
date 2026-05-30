package repository

import (
	"example/todo/internal/domain"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewUserRepository(db *sqlx.DB, log *slog.Logger) *userRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) CreateUser(user domain.CreateUser) (int, error) {
	var id int
	query := `INSERT INTO users (name, password_hash)  VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(query, user.Name, user.Password)
	r.log.Debug("", slog.Any("row", row))
	if err := row.Scan(&id); err != nil {
		r.log.Error("Create user failed", slog.Any("err", err))
		return 0, err
	}
	return id, nil
}
