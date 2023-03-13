package main

import (
	"context"
	"log"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest"
	"github.com/lstrgiang/cryptowatch-server/internal/clients/cryptowatch"
	"github.com/lstrgiang/cryptowatch-server/internal/infra/cache"
)

func main() {
	cfg := rest.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	// create new context with sigterm/sigint receiver
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// get initiaal ETH value
	initialPrice, err := cryptowatch.GetPrice()
	if err != nil {
		initialPrice = 0
	}

	// create and init websocket
	log.Printf("Initial ETH price %f", initialPrice)
	cache := cache.New(initialPrice)
	cw := cryptowatch.New(func(updates cryptowatch.IntervalUpdate) error {
		if len(updates.Intervals) <= 0 {
			return nil
		}
		// ensured there is at least one interval
		// get smallest closeTime
		closeTime := updates.Intervals[0].CloseTime
		latestPrice := updates.Intervals[0].OHLC.CloseStr
		for _, update := range updates.Intervals {
			if update.CloseTime < closeTime {
				closeTime = update.CloseTime
				latestPrice = update.OHLC.CloseStr
			}
		}
		floatVal, err := strconv.ParseFloat(latestPrice, 64)
		if err != nil {
			return err
		}
		log.Printf("Updated ETH price %f", floatVal)
		// update latest price
		cache.Update(floatVal)
		return nil
	})
	if err := cw.Init(cryptowatch.Config{
		APIKey: "NL2IXFR1GWJDR93RGBOM",
	}, []string{cryptowatch.ResourceETHUSD}); err != nil {
		panic(err)
	}
	//start consuming
	go cw.Consume(ctx, stop)

	server := rest.NewServer(cfg)
	server.Register(cache)
	server.Start(ctx, stop)
	//on sigterm
	<-ctx.Done()
	stop()
	cw.Close()
	server.Shutdown()
}
