package binance

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adshao/go-binance"
)

type OrderBook struct {
	AggTradesC            chan *binance.WsAggTradeEvent
	TradesC               chan *binance.WsTradeEvent
	KlinesC               chan *binance.WsKlineEvent
	AllMarketMiniTickersC chan binance.WsAllMiniMarketsStatEvent
	AllMarketTickersC     chan binance.WsAllMarketsStatEvent
	PartialBookDepthsC    chan *binance.WsPartialDepthEvent
	DiffDepthsC           chan *binance.WsDepthEvent
	StopC                 chan struct{}
	stops                 []chan struct{}
	dones                 []chan struct{}
}

func New() *OrderBook {
	return &OrderBook{
		AggTradesC:            make(chan *binance.WsAggTradeEvent),
		TradesC:               make(chan *binance.WsTradeEvent),
		KlinesC:               make(chan *binance.WsKlineEvent),
		AllMarketMiniTickersC: make(chan binance.WsAllMiniMarketsStatEvent),
		AllMarketTickersC:     make(chan binance.WsAllMarketsStatEvent),
		PartialBookDepthsC:    make(chan *binance.WsPartialDepthEvent),
		DiffDepthsC:           make(chan *binance.WsDepthEvent, 10000),
		StopC:                 make(chan struct{}),
	}
}

func (b *OrderBook) AggTrades(symbol string) error {
	wsAggTradesHandler := func(event *binance.WsAggTradeEvent) {
		b.AggTradesC <- event
	}

	doneC, stopC, err := binance.WsAggTradeServe(symbol, wsAggTradesHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) Klines(symbol, interval string) error {
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		b.KlinesC <- event
	}
	doneC, stopC, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) Trades(symbol string) error {
	wsTradesHandler := func(event *binance.WsTradeEvent) {
		b.TradesC <- event
	}
	doneC, stopC, err := binance.WsTradeServe(symbol, wsTradesHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) AllMarketMiniTickers() error {
	wsAllMarketMiniTickersHandler := func(event binance.WsAllMiniMarketsStatEvent) {
		b.AllMarketMiniTickersC <- event
	}
	doneC, stopC, err := binance.WsAllMiniMarketsStatServe(wsAllMarketMiniTickersHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) AllMarketTickers() error {
	wsAllMarketTickersHandler := func(event binance.WsAllMarketsStatEvent) {
		b.AllMarketTickersC <- event
	}
	doneC, stopC, err := binance.WsAllMarketsStatServe(wsAllMarketTickersHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) PartialBookDepths(symbol, levels string) error {
	wsPartialBookDepthsHandler := func(event *binance.WsPartialDepthEvent) {
		b.PartialBookDepthsC <- event
	}
	doneC, stopC, err := binance.WsPartialDepthServe(symbol, levels, wsPartialBookDepthsHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) DiffDepths(symbol string) error {
	wsDiffDepthsHandler := func(event *binance.WsDepthEvent) {
		b.DiffDepthsC <- event
	}
	doneC, stopC, err := binance.WsDepthServe(symbol, wsDiffDepthsHandler, errHandler)
	if err != nil {
		return err
	}

	b.dones = append(b.dones, doneC)
	b.stops = append(b.stops, stopC)

	return nil
}

func (b *OrderBook) StopAll() {
	for _, c := range b.stops {
		c <- struct{}{}
	}

	for _, c := range b.dones {
		<-c
	}

	b.StopC <- struct{}{}
}

func errHandler(err error) {
	log.Println(err)
}

func GetSymbols() ([]string, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		return nil, err
	}

	var data []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	symbols := make([]string, len(data))

	for i, item := range data {
		symbols[i] = item.Symbol
	}

	return symbols, nil
}
