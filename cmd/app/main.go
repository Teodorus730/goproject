package main

import (
    "log"
    "net/http"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "goproject/internal/config"
    "goproject/internal/waterService/delivery/http"
    "goproject/internal/waterService/repository"
    "goproject/internal/waterService/usecase"
)

func main() {
    cfg, err := config.New()
    if err != nil {
        log.Fatal(err)
    }
    db, err := sqlx.Connect("postgres", cfg.DB.DSN())
    if err != nil {
        log.Fatal(err)
    }
    userRepo := repository.NewUserRepo(db)
    sessionRepo := repository.NewSessionRepo(db)
    orderRepo := repository.NewOrderRepo(db)
    tariffRepo := repository.NewTariffRepo(db)
    userUC := usecase.NewUserUsecase(userRepo, sessionRepo)
    orderUC := usecase.NewOrderUsecase(orderRepo, tariffRepo)
    userH := handler.NewUserHandler(userUC)
    _ = orderUC
    mux := http.NewServeMux()
    mux.HandleFunc("/register", userH.Register)
    log.Printf("server started on :%s", cfg.Server.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, mux))
}
