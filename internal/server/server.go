package server

import (
	"goproject/internal/config"
	handler "goproject/internal/waterService/delivery/http" 
	"goproject/internal/waterService/repository"
	"goproject/internal/waterService/usecase"
	"goproject/internal/middleware"
	"goproject/internal/rabbitmq"

	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	cfg    *config.Config
	db     *sqlx.DB
	rabbit *rabbitmq.RabbitMQ
	log    *slog.Logger
	srv    *http.Server
}

func NewServer(cfg *config.Config, db *sqlx.DB, rabbit *rabbitmq.RabbitMQ, log *slog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		db:     db,
		rabbit: rabbit,
		log:    log,
	}
}

func (s *Server) Run(errCh chan error) error {
	userRepo := repository.NewUserRepo(s.db, s.log)
	sessionRepo := repository.NewSessionRepo(s.db, s.log)
	orderRepo := repository.NewOrderRepo(s.db, s.rabbit, s.log)
	tariffRepo := repository.NewTariffRepo(s.db, s.log)
	
	userUC := usecase.NewUserUsecase(userRepo, sessionRepo, s.log)
	
	orderUC := usecase.NewOrderUsecase(orderRepo, tariffRepo, s.log)

	userH := handler.NewUserHandler(userUC, s.log)
	orderH := handler.NewOrderHandler(orderUC, s.log)

	middlewareManager := middleware.NewMiddlewareManager(s.log, sessionRepo)

	r := mux.NewRouter()
	
	ordersRouter := r.PathPrefix("/orders").Subrouter()
	ordersRouter.Use(middlewareManager.JWTMiddleware)
	ordersRouter.Use(middlewareManager.SessionMiddleware)
	ordersRouter.HandleFunc("/add_order", orderH.AddNewOrder).Methods(http.MethodPost)
	ordersRouter.HandleFunc("/get_orders", orderH.GetOrdersList).Methods(http.MethodGet)

	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", userH.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", userH.LoginUser).Methods(http.MethodPost)
	

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler: r,
	}

	go func() {
		s.log.Info("Starting server", slog.String("port", s.cfg.Server.Port))
		err := s.srv.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("listen and server error: %w", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}