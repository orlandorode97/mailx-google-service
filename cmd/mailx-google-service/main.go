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
	"github.com/orlandorode97/mailx-google-service/google"
	"github.com/orlandorode97/mailx-google-service/labels"
	repopg "github.com/orlandorode97/mailx-google-service/repos/postgres"
)

func main() {
	logger := kitlog.With(kitlog.NewLogfmtLogger(os.Stdout), "ts", kitlog.DefaultTimestampUTC)

	db, err := sql.Open("postgres", "")
	if err != nil {
		logger.Log(
			"message", "it was not possible to open a new connection with the database",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
	}

	repo := repopg.New(sqlx.NewDb(db, "postgres"))
	gmailClient, err := google.New()
	if err != nil {
		logger.Log(
			"message", "could not create a gmail client",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
	}
	var labelService labels.Service
	labelService = labels.NewService(logger, repo, gmailClient)

	mux := http.NewServeMux()
	mux.Handle("/", labels.MakeHandler(labelService, logger))
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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
			"message", "stopping mailx-google-service",
			"severity", "NOTICE",
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*45)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Log(
				"message", "mailx-google-service server has stopped",
				"err", err.Error(),
				"severity", "CRITICAL",
			)
		}
		close(connClosed)
	}()

	logger.Log(
		"message", fmt.Sprintf("listening for HTTP connections on %s", server.Addr),
		"severity", "NOTICE",
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log(
			"message", err.Error(),
			"severity", "CRITICAL",
		)
	} else {
		logger.Log(
			"message", "mailx-google-service stopped",
			"severity", "NOTICE",
		)
	}
	<-connClosed
}
