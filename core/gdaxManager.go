package core

import (

	"log"
	"fmt"
	"encoding/json"
	//"github.com/Appscrunch/Multy-back/client"
	"time"
	"Multy-back-exchange-service/api"
)



type GdaxManager struct {
	tickers map[string]Ticker
	gdaxApi *api.GdaxApi
}

type GdaxTicker struct {
	BestAsk   string `json:"best_ask"`
	BestBid   string `json:"best_bid"`
	High24h   string `json:"high_24h"`
	Low24h    string `json:"low_24h"`
	Open24h   string `json:"open_24h"`
	Rate     string `json:"price"`
	Symbol string `json:"product_id"`
	Sequence  int    `json:"sequence"`
	Type      string `json:"type"`
	Volume24h string `json:"volume_24h"`
	Volume30d string `json:"volume_30d"`
}

func (ticker GdaxTicker) IsFilled() bool {
	return (len(ticker.Symbol) > 0 && len(ticker.Rate) > 0)
}



func (b *GdaxManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.gdaxApi = &api.GdaxApi{}

	go b.gdaxApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("error:", error)
			//callback(nil, error)
		} else if message != nil {
			//fmt.Printf("%s \n", message)
			var gdaxTicker GdaxTicker
			error := json.Unmarshal(message, &gdaxTicker)
			if error == nil && gdaxTicker.IsFilled()  {
				b.add(gdaxTicker)
				//fmt.Println(gdaxTicker)
			} else {
				fmt.Println( "error parsing hitBtc ticker:", error)
			}
		}
	})

	for range time.Tick(1 * time.Second) {
		//TODO: add check if data is old and don't sent it ti callback
		func() {
			values := []Ticker{}
			for _, value := range b.tickers {
				values = append(values, value)
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = values
			callback(tickerCollection, nil)
		}()
	}
}

func (b *GdaxManager) add(gdaxTicker GdaxTicker) {
	var ticker = Ticker{}
	ticker.Rate = gdaxTicker.Rate
	ticker.Symbol = gdaxTicker.Symbol
	b.tickers[ticker.Symbol] = ticker
}