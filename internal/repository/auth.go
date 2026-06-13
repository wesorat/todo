package repository

import (
	"database/sql"
	"errors"
	"example/todo/internal/domain"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// сохранение токена в бд
// получение токена из бд
// отзыв токена из бд по айди и токену

var (
	ErrNameExists           = errors.New("username already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrRefreshNotAdded      = errors.New("refresh_hash not added")
	ErrRefreshHashExists    = errors.New("refresh_hash already exists")
	ErrRefreshTokenNotFound = errors.New("refresh_hash not found")
)

type authRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthRepository(db *sqlx.DB, log *slog.Logger) *authRepository {
	return &authRepository{db: db, log: log}
}

func (r *authRepository) CreateUser(user domain.CreateUser) (int, error) {
	var id int
	query := `INSERT INTO users (name, password_hash)  VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(query, user.Name, user.Password)
	r.log.Debug("", slog.Any("row", row))
	if err := row.Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, ErrNameExists
		}
		r.log.Error("Create user failed", slog.Any("err", err))
		return 0, err
	}
	return id, nil
}

func (r *authRepository) GetUser(name string, password string) (domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE name = $1`
	err := r.db.Get(&user, query, name)
	if err != nil {
		r.log.Error(err.Error())
		if err == sql.ErrNoRows {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *authRepository) SaveRefresh(user_id int, refresh_hash string, expired_at time.Time) error {
	query := `INSERT INTO refresh (refresh_hash, expired_at, user_id) VALUES ($1, $2, $3)`
	result, err := r.db.Exec(query, refresh_hash, expired_at, user_id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrRefreshHashExists
		}
		r.log.Error("Create refresh failed", slog.Any("err", err))
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRefreshNotAdded
	}
	return nil
}

func (r *authRepository) RevokeRefreshByHash(refresh_hash string) error {
	query := `DELETE FROM refresh WHERE refresh_hash = $1`
	result, err := r.db.Exec(query, refresh_hash)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (r *authRepository) RevokeRefreshByUserID(user_id int) error {
	query := `DELETE FROM refresh WHERE user_id = $1`
	result, err := r.db.Exec(query, user_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (r *authRepository) GetUserIDByRefresh(refresh_hash string) (int, error) {
	// проверка на годность
	var user_id int
	query := `SELECT user_id FROM refresh WHERE refresh_hash = $1`
	row := r.db.QueryRow(query, refresh_hash)
	r.log.Debug("", slog.Any("row", row))
	if err := row.Scan(&user_id); err != nil {
		return 0, ErrRefreshTokenNotFound
	}
	return user_id, nil

}
