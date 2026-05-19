package repository

import (
    "github.com/jmoiron/sqlx"
    "goproject/internal/domain"
    "context"
    "log/slog"
	"fmt"
)

type orderRepo struct { 
    postgresql *sqlx.DB 
    log *slog.Logger
}

func NewOrderRepo(db *sqlx.DB, log *slog.Logger) domain.OrderRepository {
    return &orderRepo{postgresql: db, log: log}
}


func (r *orderRepo) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
    var o domain.Order
	err := r.postgresql.QueryRowContext(ctx, createOrderQuery, order.UserID, order.TariffID, order.Amount, order.Address, order.TotalPrice, order.Status).Scan(
		&o.ID, &o.UserID, &o.TariffID, &o.Amount, &o.Address, &o.TotalPrice, &o.Status, 
	)
	if err != nil {
		r.log.Error("failed to create order", slog.Any("error", err))
		return nil, err
	}

	return &o, nil
}

func (r *orderRepo) GetOrdersByUserID(ctx context.Context, userID string) ([]domain.Order, error) {
    var orders []domain.Order

    rows, err := r.postgresql.QueryxContext(
		ctx,
		getOrdersByUserIDQuery,
		userID,
	)
	if err != nil {
		r.log.Error("failed to get orders", slog.Any("error", err))
		return nil, err
	}

    for rows.Next() {
        var order domain.Order
        if err := rows.StructScan(&order); err != nil {
            return nil, fmt.Errorf("failed to scan orders for client rows: %w", err)
        }
        orders = append(orders, order)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate orders rows: %w", err)
    }

    return orders, nil
}
