package domain

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Phone string `json:"phone" db:"phone"`
	PasswordHash string `json:"-" db:"password_hash"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Session struct {
	SessionID string `db:"session_id" json:"sessionID"`
	UserID string `db:"user_id" json:"userID"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type Order struct {
	ID string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	TariffID int `json:"tariff_id" db:"tariff_id"`
	Amount int `json:"amount" db:"amount"`
	Address string `json:"address" db:"address"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
	Status string `json:"status" db:"status"`
}

type Tariff struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}

type ChangeOrderStatusData struct {
	OrderID     string    `json:"order_id" validate:"required"`
	UUIDOrderID uuid.UUID `db:"order_id" json:"-"`
	NewStatus   string    `db:"status" json:"status" validate:"required,oneof=UNDEFINED NEW ARRIVING COMPLETED CANCELED"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}