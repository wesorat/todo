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
	ErrListNotExist        = errors.New("list not exist")
	ErrListNotOwn          = errors.New("not own list")
	ErrNotOwnItem          = errors.New("not own item")
	ErrItemNotDeleted      = errors.New("item not deleted")
	ErrItemNotUpdatesField = errors.New("item has not field for update")
	ErrItemNotUpdate       = errors.New("item not update")
)

type itemsRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewItemsRepository(db *sqlx.DB, log *slog.Logger) *itemsRepository {
	return &itemsRepository{db: db, log: log}
}

func (r *itemsRepository) Create(user_id int, item domain.CreateItem) (int, error) {
	if err := r.ownLists(user_id, item.ListID); err != nil {
		r.log.Error(err.Error())
		return 0, err
	}
	var item_id int
	query := `INSERT INTO items (list_id, title, description, done) VALUES ($1, $2, $3, $4) RETURNING id`
	row := r.db.QueryRow(query, item.ListID, item.Title, item.Description, item.Done)
	if err := row.Scan(&item_id); err != nil {
		r.log.Error(err.Error())
		return 0, err
	}
	return item_id, nil
}
func (r *itemsRepository) Get(user_id, item_id int) (domain.Item, error) {
	var item domain.Item
	query := `SELECT i.id, i.list_id, i.title, i.description, i.done, i.created_at
				FROM items i
				INNER JOIN lists l
					ON l.id = i.list_id
				WHERE l.user_id = $1 AND i.id = $2
			`
	if err := r.db.Get(&item, query, user_id, item_id); err != nil {
		r.log.Error(err.Error())
		return domain.Item{}, err
	}
	return item, nil

}
func (r *itemsRepository) GetAll(user_id, list_id int) ([]domain.Item, error) {
	var items []domain.Item
	query := `SELECT i.id, i.list_id, i.title, i.description, i.done, i.created_at
				FROM lists l
				INNER JOIN items i
					ON i.list_id = l.id
				WHERE l.user_id = $1 AND l.id = $2
			`
	if err := r.db.Select(&items, query, user_id, list_id); err != nil {
		r.log.Error(err.Error())
		return []domain.Item{}, err
	}
	return items, nil
}
func (r *itemsRepository) Update(user_id, item_id int, title, description *string, done *bool) error {
	if err := r.ownItem(user_id, item_id); err != nil {
		return err
	}
	updates := make(map[string]interface{})
	args := []interface{}{}

	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = *description
	}
	if done != nil {
		updates["done"] = *done
	}
	if len(updates) == 0 {
		return ErrItemNotUpdatesField
	}
	query := "UPDATE items SET"
	i := 1
	for key, val := range updates {
		query += fmt.Sprintf(" %v = $%v,", key, i)
		args = append(args, val)
		i++
	}
	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%v", i)
	args = append(args, item_id)
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrItemNotUpdate
	}
	return nil

}
func (r *itemsRepository) Delete(user_id, item_id int) error {
	if err := r.ownItem(user_id, item_id); err != nil {
		r.log.Error(err.Error())
		return err
	}
	query := `DELETE FROM items WHERE id = $1`
	result, err := r.db.Exec(query, item_id)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	if num == 0 {
		r.log.Error(ErrItemNotDeleted.Error())
		return ErrItemNotDeleted
	}
	return nil

}

func (r *itemsRepository) ownLists(user_id, list_id int) error {
	own := false
	checkQuery := `SELECT 1 FROM lists WHERE user_id = $1 AND id = $2`
	checkRow := r.db.QueryRow(checkQuery, user_id, list_id)
	if err := checkRow.Scan(&own); err != nil {
		return ErrListNotExist
	}
	if !own {
		return ErrListNotOwn
	}
	return nil
}

func (r *itemsRepository) ownItem(user_id, item_id int) error {
	own := false
	checkQuery := `SELECT 1
					FROM lists l
					INNER JOIN items i
						ON l.id = i.list_id
					WHERE user_id = $1 AND i.id = $2`
	checkRow := r.db.QueryRow(checkQuery, user_id, item_id)
	if err := checkRow.Scan(&own); err != nil {
		return err
	}
	if !own {
		return ErrNotOwnItem
	}
	return nil
}
