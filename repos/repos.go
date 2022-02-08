package repos

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

func BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetInt("POSTGRES_PORT"),
		viper.GetString("POSTGRES_DB_NAME"),
		viper.GetString("POSTGRES_SSL_MODE"),
		viper.GetString("POSTGRES_DB_SCHEMA"))
}

type Repository interface {
	CreateUser(context.Context, *User) error
	GetUserByID(context.Context, string) (*User, error)
}
