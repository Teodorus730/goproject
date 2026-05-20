package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"goproject/internal/config"
	
	handler "goproject/internal/waterService/delivery/http" 
	"goproject/internal/waterService/repository"
	"goproject/internal/waterService/usecase"
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

	userRepo := repository.NewUserRepo(db, log)
	sessionRepo := repository.NewSessionRepo(db, log)
	orderRepo := repository.NewOrderRepo(db, log)
	tariffRepo := repository.NewTariffRepo(db, log)
	
	userUC := usecase.NewUserUsecase(userRepo, sessionRepo, log)
	
	orderUC := usecase.NewOrderUsecase(orderRepo, tariffRepo, log)

	userH := handler.NewUserHandler(userUC, log)
	orderH := handler.NewOrderHandler(orderUC, log)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", userH.Register)
	mux.HandleFunc("/login", userH.LoginUser)
	mux.HandleFunc("/add_order", orderH.AddNewOrder)
	mux.HandleFunc("/get_orders", orderH.GetOrdersList)

	log.Info("server started", slog.String("port", cfg.Server.Port))
	
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Error("server stopped unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}
}