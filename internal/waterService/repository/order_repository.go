package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/google/uuid"
	"github.com/streadway/amqp"

    "goproject/internal/domain"
	"goproject/internal/rabbitmq"
    "context"
    "log/slog"
	"encoding/json"
	"fmt"
	"time"
)

type orderRepo struct { 
    postgresql *sqlx.DB 
    rabbit     *rabbitmq.RabbitMQ
    log *slog.Logger
}

func NewOrderRepo(db *sqlx.DB, rabbit *rabbitmq.RabbitMQ, log *slog.Logger) domain.OrderRepository {
    return &orderRepo{postgresql: db, rabbit: rabbit, log: log}
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

func (s *orderRepo) ChangeOrderStatus(ctx context.Context, newStatusData *domain.ChangeOrderStatusData) error {
	var err error
	newStatusData.UUIDOrderID, err = uuid.Parse(newStatusData.OrderID)
	if err != nil{
		s.log.Error("failed to parse order uuid", slog.Any("error", err))
		return err
	}
	result, err := s.postgresql.NamedExecContext(
		ctx,
		changeOrderStatusQuery,
		newStatusData,
	)
	if err != nil {
		s.log.Error("failed to change order status", slog.Any("error", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Error("failed to get delete rows affected", slog.Any("error", err))
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no orders status changed")
	}

	return nil
}

func (r *orderRepo) PublishNewOrder(orderId string) error {
	err := r.rabbit.Channel.ExchangeDeclare(r.rabbit.PubConf.Name, r.rabbit.PubConf.Kind, r.rabbit.PubConf.Durable, r.rabbit.PubConf.AutoDelete, r.rabbit.PubConf.Internal, r.rabbit.PubConf.NoWait, r.rabbit.PubConf.Args)
	if err != nil {
		r.log.Error("failed to declare an exchange", slog.Any("errors", err))
		return err
	}

	body, err := json.Marshal(orderId)
	if err != nil {
		r.log.Error("failed to marshall error message", slog.Any("error", err))
		return err
	}

	err = r.rabbit.Channel.Publish(
		r.rabbit.PubConf.Name, r.rabbit.PubConf.QueueName, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		r.log.Error("failed to publish", slog.Any("error", err))
		return err
	}

	r.log.Info("successfully send new order to processing service")
	return nil
}

func (r *orderRepo) StartStatusChangeConsumer(queueName string) {
	q, err := r.rabbit.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		r.log.Error("failed to declare a queue", slog.Any("errors", err))
		return
	}
	msgs, err := r.rabbit.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false, false, false, nil)
	if err != nil {
		r.log.Error("failed to start consumer", slog.Any("errors", err))
		return
	}

	go func() {
		r.log.Info("Consumer started")
		for {
			select {
			case <-r.rabbit.Ctx.Done():
				r.log.Info("Consumer stopped")
				return
			case d, ok := <-msgs:
				if !ok {
					r.log.Error("rabbit channel closed")
					return
				}

				var newOrderStatus domain.ChangeOrderStatusData
				if err := json.Unmarshal(d.Body, &newOrderStatus); err != nil {
					r.log.Error("failed to parse rabbit message", slog.Any("errors", err))
					d.Nack(false, true)
					continue
				}
				newOrderStatus.UpdatedAt = time.Now()
				if err := r.changeStatus(&newOrderStatus); err != nil {
					r.log.Error("failed to change order status", slog.Any("errors", err))
					d.Nack(false, true)
					continue
				}
				d.Ack(false)
			}
		}
	}()
}

func (r *orderRepo) changeStatus(newOrderStatus *domain.ChangeOrderStatusData) error {
	changeOrderStatusCtx, changeOrderStatusCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer changeOrderStatusCancel()
	return r.ChangeOrderStatus(changeOrderStatusCtx, newOrderStatus)
}
