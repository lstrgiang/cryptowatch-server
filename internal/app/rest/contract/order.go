package contract

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lstrgiang/cryptowatch-server/internal/data"
)

//easyjson:json
type Order struct {
	ID        uuid.UUID      `json:"id"`
	Type      data.OrderType `json:"type"`
	USDPrice  float64        `json:"usd_price"`
	Amount    float64        `json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
}

//easyjson:json
type OrderSlice []Order

func NewOrderContract(data data.Order) Order {
	return Order{
		ID:        data.ID,
		Type:      data.Type,
		USDPrice:  data.USDPrice,
		Amount:    data.Amount,
		CreatedAt: data.CreatedAt,
	}
}

func NewOrdersContract(orders []data.Order) OrderSlice {
	result := make([]Order, 0)
	for _, order := range orders {
		result = append(result, NewOrderContract(order))
	}
	return result
}
