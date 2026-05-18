package domain

import "time"

type User struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Phone string `json:"phone" db:"phone"`
	PasswordHash string `json:"-" db:"password_hash"`
}

type Session struct {
	ID int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
	Token string `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type Order struct {
	ID int `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	TariffID int `json:"tarif_id" db:"tarif_id"`
	Amount int `json:"amount" db:"amount"`
	Address string `json:"address" db:"address"`
	TotalPrice int `json:"total_price" db:"total_price"`
	Status string `json:"status" db:"status"`
}

type Tariff struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Price int `json:"price" db:"price"`
}
