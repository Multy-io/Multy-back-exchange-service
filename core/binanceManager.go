package core

import (
	"log"
	"encoding/json"
	"time"
	"Multy-back-exchange-service/api"
)


type BinanceTicker struct {
	Symbol 	string `json:"s"`
	Rate	string `json:"c"`
	EventTime 	float64 `json:"E"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
	StatisticCloseTime float64 `json:"C"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
}

type BinanceManager struct {
	binanceApi *api.BinanceApi
	symbolsToParse map[string]bool
}

func NewBinanceManager() *BinanceManager {
	var manger = BinanceManager{}
	manger.symbolsToParse = map[string]bool{}
	manger.binanceApi = &api.BinanceApi{}
	return &manger
}

func (b *BinanceManager)  StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {
	b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	b.binanceApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("binance error:", error)
			callback(TickerCollection{}, error)
		} else if message != nil {
			var binanceTickers []BinanceTicker
			json.Unmarshal(message, &binanceTickers)

			var tickers = []Ticker{}

			for _, binanceTicker := range binanceTickers {
				if b.symbolsToParse[binanceTicker.Symbol] {
					var ticker= Ticker{binanceTicker.Symbol, binanceTicker.Rate}
					tickers = append(tickers, ticker)
				}
			}

			var tickerCollection TickerCollection
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = tickers
			callback(tickerCollection, nil)
		}
	} )

}

func (b *BinanceManager)  composeSybolsToParse(exchangeConfiguration ExchangeConfiguration) map[string]bool {
	var symbolsToParse = map[string]bool{}
	for _, targetCurrency := range exchangeConfiguration.TargetCurrencies {
		for _, referenceCurrency := range exchangeConfiguration.ReferenceCurrencies {

			if referenceCurrency == "USD" {
				referenceCurrency = "USDT"
			}

			symbol := targetCurrency + referenceCurrency
			symbolsToParse[symbol] = true
		}
	}
	return symbolsToParse

}
