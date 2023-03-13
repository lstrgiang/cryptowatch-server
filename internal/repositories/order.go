package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lstrgiang/cryptowatch-server/internal/data"
)

type (
	NewOrderRepositoryFn func() OrderRepository
	OrderRepository      interface {
		ListOrderByUserID(ctx context.Context, db *sqlx.DB, userId int64) ([]data.Order, error)
		Create(ctx context.Context, db *sqlx.DB, newData *data.Order) (*data.Order, error)
	}
	orderRepository struct {
	}
)

func NewOrderRepository() OrderRepository {
	return orderRepository{}
}

func (r orderRepository) ListOrderByUserID(ctx context.Context, db *sqlx.DB, userId int64) ([]data.Order, error) {
	var data []data.Order
	err := db.Select(&data, `SELECT * FROM orders WHERE user_id = ?`, userId)
	return data, err
}

func (r orderRepository) Create(ctx context.Context, db *sqlx.DB, newData *data.Order) (*data.Order, error) {
	nstmt, err := db.PrepareNamed(`
		INSERT INTO orders (user_id, usd_price, amount, type) 
		VALUES (:user_id, :usd_price, :amount, :type)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	var createdData *data.Order
	if err := nstmt.Select(createdData, *newData); err != nil {
		return nil, err
	}
	return createdData, nil
}
