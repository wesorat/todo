package domain

type List struct {
	ID int `json:"id" db:"id"`
	CreateList
}

type CreateList struct {
	UserID      int    `json:"user_id" db:"user_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}
