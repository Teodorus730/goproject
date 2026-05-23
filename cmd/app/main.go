package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"goproject/internal/config"
	"goproject/internal/server"
	"goproject/internal/rabbitmq"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.New()
	if err != nil {
		log.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	
	termCtx, termCancel := context.WithCancel(context.Background())
	go waitSigterm(termCancel, log)

	rabbit, err := rabbitmq.NewRabbitMQ(&cfg.RabbitMQ, termCtx, log)
	if err != nil {
		log.Error("failed to create rabbit connection", slog.Any("error", err))
		return
	}

	errCh := make(chan error, 1)

	srv := server.NewServer(cfg, db, rabbit, log)
	err = srv.Run(errCh)
	if err != nil {
		log.Error("failed to run server", slog.Any("error", err))
	}

	select {
	case err := <-errCh:
		log.Error("got error from server", slog.Any("error", err))
	case <-termCtx.Done():
		err := srv.Stop()
		if err != nil {
			log.Error("shutdown failed", slog.Any("error", err))
		}
	}

	log.Info("service terminated")
}

func waitSigterm(terminate context.CancelFunc, log *slog.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	caughtSignal := <-sigCh

	log.Warn("service starts termination", slog.String("signal", caughtSignal.String()))

	signal.Stop(sigCh)
	close(sigCh)
	terminate()
}