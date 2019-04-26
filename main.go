package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/orionprotocol/price-feed/exchanges/poloniex"

	"github.com/orionprotocol/price-feed/exchanges/bittrex"

	"github.com/orionprotocol/price-feed/api"
	"github.com/orionprotocol/price-feed/config"
	"github.com/orionprotocol/price-feed/exchanges/binance"
	"github.com/orionprotocol/price-feed/logger"
	"github.com/orionprotocol/price-feed/storage"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	cfg, err := config.FromFile()
	if err != nil {
		log.Fatalf("Could not read config: %v. Exiting", err)
	}

	l := logger.New(cfg.Logger)
	defer func() {
		if err = l.Close(); err != nil {
			log.Printf("Could not close logger: %v", err)
		}
	}()

	database := storage.New(cfg.Storage, l)
	pong, err := database.Check()
	if err != nil {
		l.Fatalf("Can't establish connection to database: %v", err)
	}
	l.Infof("Database check reply: %v", pong)

	if err := database.Flush(); err != nil {
		l.Fatalf("Could not flush database")
	}

	binanceWorker, err := binance.NewWorker(cfg.Binance, l, database, quit)
	if err != nil {
		l.Fatalf("Could not connect to Binance: %v", err)
	}

	binanceWorker.Start()

	bittrexWorker, err := bittrex.NewWorker(cfg.Bittrex, l, database, quit)
	if err != nil {
		l.Fatalf("Could not connect to Bittrex: %v", err)
	}

	bittrexWorker.Start()

	poloniexWorker, err := poloniex.NewWorker(cfg.Poloniex, l, database, quit)
	if err != nil {
		l.Fatalf("Could not connect to Bittrex: %v", err)
	}

	poloniexWorker.Start()

	apiServer := api.New(cfg.API, l, database, binanceWorker, bittrexWorker, poloniexWorker)

	go func() {
		if err = apiServer.Start(); err != nil {
			l.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
}
