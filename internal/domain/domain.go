package domain

import (
    "context"
)

type UserRepository interface {
    CreateUser(ctx context.Context, user *User) (*User, error)
    GetUserByID(ctx context.Context, id string) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type SessionRepository interface {
    CreateSession(ctx context.Context, sessionID, userID string) error
    GetSessionByUserID(ctx context.Context, userID string) (*Session, error)
    UpdateSessionExpiry(ctx context.Context, sessionID string) error
}

type OrderRepository interface {
    CreateOrder(ctx context.Context, order *Order) (*Order, error)
    GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error)
}

type TariffRepository interface {
    GetTariffByID(ctx context.Context, id int) (*Tariff, error)
}

type UserUsecase interface {
    CreateUser(ctx context.Context, name, email, phone, password string) (*User, error)
    CreateSession(ctx context.Context, sessionID, userID string) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type OrderUsecase interface {
    PlaceOrder(ctx context.Context, userID string, tariffID int, amount int, address string) (*Order, error)
    GetOrders(ctx context.Context, userID string) ([]Order, error)
}

