package cryptowatch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

type (
	Subscription struct {
		StreamSubscription `json:"streamSubscription"`
	}

	StreamSubscription struct {
		Resource string `json:"resource"`
	}

	SubscribeRequest struct {
		Subscriptions []Subscription `json:"subscriptions"`
	}
	PriceResult struct {
		Result struct {
			Price float64 `json:"price"`
		} `json:"result"`
	}
	Update struct {
		MarketUpdate struct {
			Market struct {
				MarketId int `json:"marketId,string"`
			} `json:"market"`
			IntervalUpdate IntervalUpdate `json:"intervalsUpdate"`
		} `json:"marketUpdate"`
	}
	IntervalUpdate struct {
		Intervals []struct {
			CloseTime string `json:"closetime"`
			Period    int    `json:"period"`
			OHLC      struct {
				OpenStr  string `json:"openStr"`
				HighStr  string `json:"highStr"`
				LowStr   string `json:"lowStr"`
				CloseStr string `json:"closeStr"`
			} `json:"ohlc"`
			VolumeBaseStr  string `json:"volumeBaseStr"`
			VolumeQuoteStr string `json:"volumeQuoteStr"`
		} `json:"intervals"`
	}
	Config struct {
		APIKey string
	}
	CryptoWatch interface {
		Init(cfg Config, resources []string) error
		Consume(ctx context.Context, stop context.CancelFunc)
		Close()
	}
	cryptoWatch struct {
		conn             *websocket.Conn
		processMessageFn func(IntervalUpdate) error
		c                chan os.Signal
	}
)

const (
	webSocketURL = "wss://stream.cryptowat.ch/connect?apikey=%s"

	ResourceETHUSD = "markets:68:ohlc"
)

func GetPrice() (float64, error) {
	response, err := http.Get(
		"https://api.cryptowat.ch/markets/kraken/ethusd/price",
	)
	if err != nil {
		return 0, err
	}
	defer func() {
		// close request
		io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	var result PriceResult
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	return result.Result.Price, nil
}

func New(processMessageFn func(IntervalUpdate) error) CryptoWatch {
	return &cryptoWatch{
		c:                make(chan os.Signal, 1),
		processMessageFn: processMessageFn,
	}
}

func (cw cryptoWatch) Close() {
	// stop consuming if any
	signal.Notify(cw.c, os.Interrupt, syscall.SIGTERM)
	// close the connection
	cw.conn.Close()
}

// initialize cryptowatch connection
func (cw *cryptoWatch) Init(cfg Config, resources []string) error {
	dialUrl := fmt.Sprintf(webSocketURL, cfg.APIKey)
	c, _, err := websocket.DefaultDialer.Dial(dialUrl, nil)
	if err != nil {
		c.Close()
		return err
	}

	// Read first message, which should be an authentication response
	_, message, err := c.ReadMessage()
	var authResult struct {
		AuthenticationResult struct {
			Status string `json:"status"`
		} `json:"authenticationResult"`
	}
	err = json.Unmarshal(message, &authResult)
	if err != nil {
		c.Close()
		return err
	}

	// Send a JSON payload to subscribe to a list of resources
	// Read more about resources here: https://docs.cryptowat.ch/websocket-api/data-subscriptions#resources

	subMessage := struct {
		Subscribe SubscribeRequest `json:"subscribe"`
	}{}
	// No map function in golang :-(
	for _, resource := range resources {
		subMessage.Subscribe.Subscriptions = append(subMessage.Subscribe.Subscriptions, Subscription{StreamSubscription: StreamSubscription{Resource: resource}})
	}
	msg, err := json.Marshal(subMessage)
	err = c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		c.Close()
		return err
	}
	cw.conn = c
	return nil
}

// process update message, store to any database or caching
func (cw cryptoWatch) ProcessUpdate(update Update) {
	if err := cw.processMessageFn(update.MarketUpdate.IntervalUpdate); err != nil {
		fmt.Println(fmt.Sprintf("cannot process update with error %s", err.Error()))
	}
	for _, interval := range update.MarketUpdate.IntervalUpdate.Intervals {
		log.Printf(
			"ETH/USD trade on market %d: %s",
			update.MarketUpdate.Market.MarketId,
			interval.OHLC.CloseStr,
		)
	}
}

// start consuming
func (cw cryptoWatch) Consume(ctx context.Context, stop context.CancelFunc) {
	for {
		// Process incoming ETH/USD trades
		_, message, err := cw.conn.ReadMessage()
		log.Println("processing update message")
		if err != nil {
			log.Fatal("Error reading from connection", err)
			return
		}
		// unmarshal message
		var update Update
		err = json.Unmarshal(message, &update)
		if err != nil {
			panic(err)
		}
		// process message if no err occur
		cw.ProcessUpdate(update)
	}
}
