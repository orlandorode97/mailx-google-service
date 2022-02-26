package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/log"
	_ "github.com/lib/pq"
	_ "github.com/orlandorode97/mailx-google-service/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

func main() {
	flag.Parse()
	command := flag.Args()[0]
	var db *sql.DB
	logger := log.With(log.NewLogfmtLogger(os.Stdout), "ts", log.DefaultTimestampUTC)

	flag.Parse()
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Log(
			"message", "it was not possible to read the .env file",
			"error", err.Error(),
			"severity", "CRITITAL",
		)
		return
	}
	psInfo := buildDSN()

	db, err := goose.OpenDBWithDriver("postgres", psInfo)
	if err != nil {
		logger.Log(
			"message", "it was not possible to open the database",
			"error", err.Error(),
			"severity", "CRITITAL",
		)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Log(
				"message", "it was not possible to close the database",
				"error", err.Error(),
				"severity", "CRITITAL",
			)
			return
		}
	}()

	if err := goose.Run(command, db, "./migrations/"); err != nil {
		logger.Log(
			"message", "it was not possible to run the migrations",
			"error", err.Error(),
			"severity", "CRITITAL",
		)
		return
	}

}

func buildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetInt("POSTGRES_PORT"),
		viper.GetString("POSTGRES_DB_NAME"),
		viper.GetString("POSTGRES_SSL_MODE"),
	)
}
