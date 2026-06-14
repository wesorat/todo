package repository

import (
	"example/todo/internal/domain"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type itemsRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewItemsRepository(db *sqlx.DB, log *slog.Logger) *itemsRepository {
	return &itemsRepository{db: db, log: log}
}

func (r *itemsRepository) Create(user_id, list_id int, item *domain.CreateItem) (int, error)
func (r *itemsRepository) Get(user_id, item int) (domain.Item, error)
func (r *itemsRepository) GetAll(user_id, list_id int) ([]domain.Item, error)
func (r *itemsRepository) Update(user_id, item_id int, title, description string, done bool) error
func (r *itemsRepository) Delete(user_id, item_id int) error
