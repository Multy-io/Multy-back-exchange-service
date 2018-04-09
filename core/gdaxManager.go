package core

import (

	"log"
	"fmt"
	"encoding/json"
	//"github.com/Appscrunch/Multy-back/client"
	"time"
	"Multy-back-exchange-service/api"
	"strings"
	"Multy-back-exchange-service/currencies"
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

func (b *GdaxTicker) getCurriences() (currencies.Currency, currencies.Currency) {

	if len(b.Symbol) > 0 {
		var symbol = b.Symbol
		var currencyCodes = strings.Split(symbol, "-")
		if len(currencyCodes) > 1 {
			return currencies.NewCurrencyWithCode(currencyCodes[0]), currencies.NewCurrencyWithCode(currencyCodes[1])
		}
	}
	return currencies.NotAplicable, currencies.NotAplicable
}



func (b *GdaxManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.gdaxApi = &api.GdaxApi{}

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	go b.gdaxApi.StartListen(apiCurrenciesConfiguration, func(message []byte, error error) {
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
		func() {
			values := []Ticker{}
			for _, value := range b.tickers {
				if value.TimpeStamp.After(time.Now().Add(-3 * time.Second)) {
					values = append(values, value)
				}
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = values
			if len(tickerCollection.Tickers) > 0 {
				callback(tickerCollection, nil)
			}
		}()
	}
}

func (b *GdaxManager) add(gdaxTicker GdaxTicker) {
	var ticker = Ticker{}
	ticker.Rate = gdaxTicker.Rate
	ticker.Symbol = gdaxTicker.Symbol

	targetCurrency, referenceCurrency  := gdaxTicker.getCurriences()
	ticker.TargetCurrency = targetCurrency
	ticker.ReferenceCurrency = referenceCurrency
	ticker.TimpeStamp = time.Now()

	b.tickers[ticker.Symbol] = ticker
}