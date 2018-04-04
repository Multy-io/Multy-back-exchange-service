package core

import (
	"sync"
	"time"
	"fmt"
)

type TickerCollection struct {
	TimpeStamp time.Time
	Tickers []Ticker
}

type Ticker struct {
	Symbol 	string
	Rate	string
}



type Agregator struct {
	sync.Mutex
	allTickers map[string]TickerCollection
}


func NewAgregator() *Agregator {
	var agregator = Agregator{}
	agregator.allTickers = make(map[string]TickerCollection)
	return &agregator
}


func (b *Agregator) add(tickerCollection TickerCollection, forExchange string) {
	b.Lock()
	b.allTickers[forExchange] = tickerCollection
	fmt.Println(b.allTickers)
	b.Unlock()
}

func (b *Agregator) getTickers(startDate time.Time) map[string]TickerCollection {
	var filteredTickers = make(map[string]TickerCollection)
	b.Lock()
	for exhange, tickerCollection := range b.allTickers {
		if tickerCollection.TimpeStamp.After(startDate) {
			filteredTickers[exhange] = tickerCollection
		}
	}
	b.Unlock()
	return filteredTickers
}
