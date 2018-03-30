package Api

import (
	"log"
	"encoding/json"
	"time"
)

//main level, shoud go to main
type TickerCollection struct {
	TimpeStamp time.Time
	Tickers []Ticker `json:"array"`
}

//main level, shoud go to main
type Ticker struct {
	Symbol 	string `json:"s"`
	Rate	string `json:"c"`
}


type BinanceTicker struct {
	Symbol 	string `json:"s"`
	Rate	string `json:"c"`
	EventTime 	float64 `json:"E"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
	StatisticCloseTime float64 `json:"C"` // field is not needed but it's a workaround because unmarshal is case insensitive and without this filed json can't be parsed
}

type BinanceManager struct {

}


var binanceApi = BinanceApi{}

func (b BinanceManager)  StartListen(callback func(tickerCollection TickerCollection, error error)) {
	binanceApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("binance error:", error)
			callback(TickerCollection{}, error)
		} else if message != nil {
			var binanceTickers []BinanceTicker
			json.Unmarshal(message, &binanceTickers)

			var tickers = []Ticker{}

			for _, binanceTicker := range binanceTickers {
				var ticker = Ticker{binanceTicker.Symbol, binanceTicker.Rate}
				tickers = append(tickers, ticker)
			}

			var tickerCollection TickerCollection
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = tickers
			callback(tickerCollection, nil)
		}
	} )

}

