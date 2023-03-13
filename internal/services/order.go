package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lstrgiang/cryptowatch-server/internal/data"
	"github.com/lstrgiang/cryptowatch-server/internal/repositories"
)

type (
	NewOrderServiceFn func(db *sqlx.DB) OrderService
	OrderService      interface {
		NewOrder(ctx context.Context, order *data.Order) (*data.Order, error)
		ListOrder(ctx context.Context, userId int64) ([]data.Order, error)
	}
	orderService struct {
		db        *sqlx.DB
		orderRepo repositories.OrderRepository
	}
)

func NewOrderService(db *sqlx.DB) OrderService {
	return orderService{
		db:        db,
		orderRepo: repositories.NewOrderRepository(),
	}
}

func (s orderService) NewOrder(ctx context.Context, order *data.Order) (*data.Order, error) {
	return s.orderRepo.Create(ctx, order)
}

func (s orderService) ListOrder(ctx context.Context, userId int64) ([]data.Order, error) {
	return s.orderRepo.ListOrderByUserID(ctx, userId)
}
