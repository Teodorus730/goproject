package usecase

import (
    "goproject/internal/domain"
    "context"
    "log/slog"
)

type orderUsecase struct {
    orderRepo domain.OrderRepository
	tariffRepo domain.TariffRepository
    log *slog.Logger
}

func NewOrderUsecase(or domain.OrderRepository, tr domain.TariffRepository, log *slog.Logger) domain.OrderUsecase {
    return &orderUsecase{orderRepo: or, tariffRepo: tr, log: log}
}

func (u *orderUsecase) PlaceOrder(ctx context.Context, userID string, tariffID int, amount int, address string) (*domain.Order, error) {
	tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
    if err != nil {
        u.log.Error("Failed to place order", slog.Any("error", err))
        return nil, err
    }

    total := amount * tariff.Price
    order := &domain.Order{UserID: userID, TariffID: tariffID, Amount: amount, Address: address, TotalPrice: total, Status: "NEW"}
    created_order, err1 := u.orderRepo.CreateOrder(ctx, order)
    return created_order, err1
}

func (u *orderUsecase) GetOrders(ctx context.Context, userID string) ([]domain.Order, error) {
    return u.orderRepo.GetOrdersByUserID(ctx, userID)
}
