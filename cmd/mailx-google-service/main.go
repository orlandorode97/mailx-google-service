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

	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/auth"
	"github.com/orlandorode97/mailx-google-service/labels"
	"github.com/orlandorode97/mailx-google-service/messages"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	repopg "github.com/orlandorode97/mailx-google-service/pkg/repos/postgres"
	"github.com/orlandorode97/mailx-google-service/users"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {
	logger := log.With(log.NewLogfmtLogger(os.Stdout), "ts", log.DefaultTimestampUTC)
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

	if repo == nil {
		logger.Log(
			"message", "it was not possible to open a new connection with the database.",
			"severity", "CRITICAL",
		)
		return
	}

	oauthConfig := google.NewConfig()

	mailxSvc := mailx.New(logger, repo, oauthConfig)

	authSvc := auth.New(logger, oauthConfig, repo, mailxSvc)
	labelsSvc := labels.New(logger, repo, mailxSvc)
	usersSvc := users.New(logger, repo, mailxSvc)
	messagesSvc := messages.New(logger, repo, mailxSvc)

	mux := http.NewServeMux()
	mux.Handle("/labels/", labels.MakeHandler(labelsSvc, logger))
	mux.Handle("/auth/", auth.MakeHandler(authSvc, logger))
	mux.Handle("/users/", users.MakeHandler(usersSvc, logger))
	mux.Handle("/messages/", messages.MakeHandler(messagesSvc, logger))

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://localhost:3000", "http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(mux),
	}

	listenAndServe(server, logger)
}

// listenAndServe gracefully shutdowns the mailx-google-service
func listenAndServe(server *http.Server, logger log.Logger) {
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
