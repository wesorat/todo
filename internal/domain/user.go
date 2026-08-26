package domain

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"-" db:"password_hash" binding:"required"`
}

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
