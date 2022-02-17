package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/auth"
	"github.com/orlandorode97/mailx-google-service/google"
	"github.com/orlandorode97/mailx-google-service/labels"
	"github.com/orlandorode97/mailx-google-service/repos"
	repopg "github.com/orlandorode97/mailx-google-service/repos/postgres"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {
	logger := kitlog.With(kitlog.NewLogfmtLogger(os.Stdout), "ts", kitlog.DefaultTimestampUTC)
	err := setViperConfig()
	if err != nil {
		logger.Log(
			"message", "it was not possible to read the .env file.",
			"error", err.Error(),
			"severity", "CRITITAL",
		)
		return
	}
	db, err := sql.Open("postgres", repos.BuildDSN())
	if err != nil {
		logger.Log(
			"message", "it was not possible to open a new connection with the database.",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
		return
	}

	repo := repopg.New(sqlx.NewDb(db, "postgres"))

	oauthConfig := google.NewConfig()

	mailxSvc := mailx.NewService(logger, oauthConfig)

	var labelsSvc labels.Service
	labelsSvc = labels.NewService(logger, repo, mailxSvc)

	var authSvc auth.Service
	authSvc = auth.NewService(logger, repo, mailxSvc)

	mux := http.NewServeMux()
	mux.Handle("/labels/", labels.MakeHandler(labelsSvc, logger))
	mux.Handle("/login/", auth.MakeHandler(authSvc, logger))

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}))

	c := cors.AllowAll()

	server := &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(mux),
	}

	listenAndServe(server, logger)
}

// listenAndServe gracefully shutdowns the mailx-google-service
func listenAndServe(server *http.Server, logger kitlog.Logger) {
	connClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		logger.Log(
			"message", "stopping mailx-google-service.",
			"severity", "NOTICE",
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*45)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Log(
				"message", "mailx-google-service server has stopped.",
				"err", err.Error(),
				"severity", "CRITICAL",
			)
		}
		close(connClosed)
	}()

	logger.Log(
		"message", fmt.Sprintf("listening for HTTP connections on %s.", server.Addr),
		"severity", "NOTICE",
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log(
			"message", err.Error(),
			"severity", "CRITICAL",
		)
	} else {
		logger.Log(
			"message", "mailx-google-service stopped.",
			"severity", "NOTICE",
		)
	}
	<-connClosed
}

func setViperConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}
