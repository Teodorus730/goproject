package usecase

import (
    "crypto/sha256"
    "fmt"
    "goproject/internal/domain"
)

// Юзкейсы для юзера
type userUsecase struct {
    userRepo    domain.UserRepository
    sessionRepo domain.SessionRepository
}

func NewUserUsecase(ur domain.UserRepository, sr domain.SessionRepository) domain.UserUsecase {
    return &userUsecase{userRepo: ur, sessionRepo: sr}
}

func (u *userUsecase) Register(name, email, phone, password string) (*domain.User, error) {
    hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
    user := &domain.User{Name: name, Email: email, PasswordHash: hash, Phone: phone}
    err := u.userRepo.CreateUser(user)
    return user, err
}

func (u *userUsecase) Login(email, password string) (*domain.Session, error) {
    hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
    _ = hash
    token := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))
    session := &domain.Session{Token: token}
    err := u.sessionRepo.CreateSession(session)
    return session, err
}

// Юзкейсы для заказов
type orderUsecase struct {
    orderRepo domain.OrderRepository
	tariffRepo domain.TariffRepository
}

func NewOrderUsecase(or domain.OrderRepository, tr domain.TariffRepository) domain.OrderUsecase {
    return &orderUsecase{orderRepo: or, tariffRepo: tr}
}

func (o *orderUsecase) PlaceOrder(userID string, tariffID int, amount int, address string) (*domain.Order, error) {
	tariff, err := o.tariffRepo.GetTariffByID(tariffID)
    if err != nil {
        return nil, err
    }

    total := amount * tariff.Price
    order := &domain.Order{UserID: userID, TariffID: tariffID, Amount: amount, Address: address, TotalPrice: total, Status: "NEW"}
    err = o.orderRepo.CreateOrder(order)
    return order, err
}

func (o *orderUsecase) GetOrders(userID string) ([]domain.Order, error) {
    return o.orderRepo.GetOrdersByUserID(userID)
}
