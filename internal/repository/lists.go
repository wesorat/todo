package repository

import (
	"example/todo/internal/domain"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type listsRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewListsRepository(db *sqlx.DB, log *slog.Logger) *listsRepository {
	return &listsRepository{db: db, log: log}
}

func (r *listsRepository) Create(user_id int, list *domain.CreateList) (int, error)
func (r *listsRepository) Get(user_id, list_id int) (domain.List, error)
func (r *listsRepository) GetAll(user_id int) ([]domain.List, error)
func (r *listsRepository) Update(user_id, list_id int, title, description string) error
func (r *listsRepository) Delete(user_id, list_id int) error
