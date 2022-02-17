package repos

import (
	"context"
	"fmt"

	"github.com/orlandorode97/mailx-google-service/models"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

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
	CreateUser(context.Context, *models.User) error
	GetUserByID(context.Context, string) (*models.User, error)
	GetTokenByUserId(context.Context, string) (*models.Token, error)
	SaveAccessToken(context.Context, string, *oauth2.Token) error
	UpdateAccessToken(context.Context, string, *oauth2.Token) error
}
