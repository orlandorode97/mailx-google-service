package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Repository struct {
	db *sqlx.DB
}

func BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=%s",
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetInt("POSTGRES_PORT"),
		viper.GetString("POSTGRES_SSL_MODE"))
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
