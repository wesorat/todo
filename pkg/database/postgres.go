package database

import (
	"example/todo/pkg/config"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(s config.Database) (*sqlx.DB, error) {
	password := os.Getenv("DB_PASSWORD")

	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			s.Host, s.Port, s.Username, password, s.DBname, s.SSLMode),
	)

	if err != nil {
		return nil, err
	}
	return db, err
}
