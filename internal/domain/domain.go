package domain

type UserRepository interface {
    CreateUser(user *User) error
    GetUserByID(id string) (*User, error)
}

type SessionRepository interface {
    CreateSession(session *Session) error
    GetSessionByToken(token string) (*Session, error)
}

type OrderRepository interface {
    CreateOrder(order *Order) error
    GetOrdersByUserID(userID string) ([]Order, error)
}

type TariffRepository interface {
    GetTariffByID(id int) (*Tariff, error)
}

type UserUsecase interface {
    Register(name, email, phone, password string) (*User, error)
    Login(email, password string) (*Session, error)
}

type OrderUsecase interface {
    PlaceOrder(userID string, tariffID int, amount int, address string) (*Order, error)
    GetOrders(userID string) ([]Order, error)
}