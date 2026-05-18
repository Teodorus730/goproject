package repository

import (
    "github.com/jmoiron/sqlx"
    "goproject/internal/domain"
)


// Юзеры
type userRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) domain.UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) CreateUser(u *domain.User) error {
    _, err := r.db.Exec(
        "INSERT INTO users(name,email,phone,password_hash) VALUES($1,$2,$3,$4)",
        u.Name, u.Email, u.Phone, u.PasswordHash,
    )
    return err
}

func (r *userRepo) GetUserByID(id string) (*domain.User, error) {
    var u domain.User
    err := r.db.Get(&u, "SELECT * FROM users WHERE id=$1", id)
    return &u, err
}


// Сессии
type sessionRepo struct{ db *sqlx.DB }

func NewSessionRepo(db *sqlx.DB) domain.SessionRepository {
    return &sessionRepo{db: db}
}

func (r *sessionRepo) CreateSession(s *domain.Session) error {
    _, err := r.db.Exec(
        "INSERT INTO sessions(user_id,token,expires_at) VALUES($1,$2,$3)",
        s.UserID, s.Token, s.ExpiresAt,
    )
    return err
}

func (r *sessionRepo) GetSessionByToken(token string) (*domain.Session, error) {
    var s domain.Session
    err := r.db.Get(&s, "SELECT * FROM sessions WHERE token=$1", token)
    return &s, err
}


// Заказы
type orderRepo struct{ db *sqlx.DB }

func NewOrderRepo(db *sqlx.DB) domain.OrderRepository {
    return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(o *domain.Order) error {
    _, err := r.db.Exec(
        "INSERT INTO orders(user_id,tarif_id,amount,address,total_price,status) VALUES($1,$2,$3,$4,$5,$6)", 
        o.UserID, o.TariffID, o.Amount, o.Address, o.TotalPrice, o.Status,
    )
    return err
}

func (r *orderRepo) GetOrdersByUserID(userID string) ([]domain.Order, error) {
    var orders []domain.Order
    err := r.db.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", userID)
    return orders, err
}

// Тарифы
type tariffRepo struct{ db *sqlx.DB }

func NewTariffRepo(db *sqlx.DB) domain.TariffRepository {
    return &tariffRepo{db: db}
}

func (r *tariffRepo) GetTariffByID(id int) (*domain.Tariff, error) {
    var t domain.Tariff
    err := r.db.Get(&t, "SELECT * FROM tariffes WHERE id=$1", id)
    return &t, err
}