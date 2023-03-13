package data

import (
	"time"

	"github.com/gofrs/uuid"
)

type OrderType string

const (
	SellOrder OrderType = "sell"
	BuyOrder  OrderType = "buy"
)

type Order struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	USDPrice  float64   `db:"usd_price"`
	Amount    float64   `db:"eth_price"`
	Type      OrderType `db:"type"`
	CreatedAt time.Time `db:"created_at"`
}
