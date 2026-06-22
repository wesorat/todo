package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wesorat/todo/internal/domain"

	"github.com/jmoiron/sqlx"
)

var (
	ErrListNotUpdatesField = errors.New("list has not fields for update")
	ErrListNotUpdate       = errors.New("list not update")
	ErrListNotDeleted      = errors.New("list not deleted")
)

type listsRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewListsRepository(db *sqlx.DB, log *slog.Logger) *listsRepository {
	return &listsRepository{db: db, log: log}
}

func (r *listsRepository) Create(list domain.CreateList) (int, error) {
	var list_id int
	query := `INSERT INTO lists (user_id, title, description) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, list.UserID, list.Title, list.Description)
	if err := row.Scan(&list_id); err != nil {
		return 0, err
	}
	return list_id, nil

}
func (r *listsRepository) Get(user_id, list_id int) (domain.List, error) {
	var list domain.List
	query := `SELECT * FROM lists WHERE id = $1 AND user_id = $2`
	err := r.db.Get(&list, query, list_id, user_id)
	if err != nil {
		return domain.List{}, err
	}
	return list, nil
}
func (r *listsRepository) GetAll(user_id int) ([]domain.List, error) {
	var lists []domain.List
	query := `SELECT * FROM lists WHERE user_id = $1`
	err := r.db.Select(&lists, query, user_id)
	if err != nil {
		return []domain.List{}, err
	}
	return lists, nil
}
func (r *listsRepository) Update(user_id, list_id int, title, description *string) error {
	query := "UPDATE lists SET"
	updates := make(map[string]any)
	args := []any{}

	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = *description
	}
	if len(updates) == 0 {
		return ErrListNotUpdatesField
	}
	i := 1
	for key, val := range updates {
		query += fmt.Sprintf(" %v = $%v,", key, i)
		args = append(args, val)
		i++
	}
	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE user_id = $%v AND id = $%v", i, i+1)
	args = append(args, user_id, list_id)
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrListNotUpdate
	}
	return nil
}
func (r *listsRepository) Delete(user_id, list_id int) error {
	query := `DELETE FROM lists WHERE user_id = $1 AND id = $2`
	result, err := r.db.Exec(query, user_id, list_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrListNotDeleted
	}
	return nil
}
