package repository

import (
	"database/sql"
	"errors"
	"example/todo/internal/domain"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNameExists      = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type userRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewUserRepository(db *sqlx.DB, log *slog.Logger) *userRepository {
	return &userRepository{db: db, log: log}
}

// TODO проверка на уникальность имени
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

func (r *userRepository) GetUser(name string, password string) (domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, name)
	if err != nil {
		r.log.Error(err.Error())
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, err
	}
	//TODO check password
	return user, nil
}
