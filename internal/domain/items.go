package domain

import "time"

type Item struct {
	ID int `json:"id" db:"id"`
	CreateItem
}

type CreateItem struct {
	ListID      int       `json:"list_id" db:"list_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Done        bool      `json:"done" db:"done"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UpdateItem struct {
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Done        *bool   `json:"done" db:"done"`
}
