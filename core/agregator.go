package core

import (
	"sync"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
	"github.com/KristinaEtc/slf"
)

var log = slf.WithContext("core")

type TickerCollection struct {
	TimpeStamp time.Time
	Tickers    []Ticker
}

type Ticker struct {
	TargetCurrency    currencies.Currency
	ReferenceCurrency currencies.Currency
	Symbol            string
	Rate              string
	TimpeStamp        time.Time
}

type Agregator struct {
	sync.Mutex
	allTickers map[string]*TickerCollection
}

func NewAgregator() *Agregator {
	var agregator = Agregator{}
	agregator.allTickers = make(map[string]*TickerCollection)
	return &agregator
}

func (b *Agregator) add(tickerCollection TickerCollection, forExchange string) {
	b.Lock()
	b.allTickers[forExchange] = &tickerCollection
	//fmt.Println("added:", tickerCollection)
	b.Unlock()
}

func (b *Agregator) getTickers(startDate time.Time) map[string]TickerCollection {
	var filteredTickers = make(map[string]TickerCollection)
	b.Lock()
	allTickers := b.allTickers
	b.Unlock()
	for exhange, tickerCollection := range allTickers {
		if tickerCollection.TimpeStamp.After(startDate) {
			filteredTickers[exhange] = *tickerCollection
		}
	}

	return filteredTickers
}
