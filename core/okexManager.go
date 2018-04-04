package core

import (

	"log"
	//"fmt"
	"encoding/json"
	//"github.com/Appscrunch/Multy-back/client"
	"time"
	//"strconv"
	"Multy-back-exchange-service/api"
	)



type OkexManager struct {
	tickers map[string]Ticker
	okexApi *api.OkexApi
}

type OkexTicker struct {
	Binary  int    `json:"binary"`
	Symbol string `json:"channel"`
	Data    struct {
		High       string `json:"high"`
		LimitLow   string `json:"limitLow"`
		Vol        string `json:"vol"`
		Last       string `json:"last"`
		Low        string `json:"low"`
		Buy        string `json:"buy"`
		HoldAmount string `json:"hold_amount"`
		Sell       string `json:"sell"`
		ContractID int64  `json:"contractId"`
		UnitAmount string `json:"unitAmount"`
		LimitHigh  string `json:"limitHigh"`
	} `json:"data"`
}



func (ticker OkexTicker) IsFilled() bool {
	return (len(ticker.Symbol) > 0 && len(ticker.Data.Last) > 0)
}



func (b *OkexManager) StartListen(callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.okexApi = &api.OkexApi{}

	go b.okexApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("error:", error)
			//callback(nil, error)
		} else if message != nil {
			//fmt.Printf("%s \n", message)
			b.addMessage(message)
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

func (b *OkexManager) add(okexTicker OkexTicker) {
	var ticker = Ticker{}
	ticker.Rate = okexTicker.Data.Last
	ticker.Symbol = okexTicker.Symbol
	b.tickers[ticker.Symbol] = ticker
}



func (b *OkexManager) addMessage (message []byte) {

	var okexTickers =[]OkexTicker{}
	json.Unmarshal(message, &okexTickers)

	for _, okexTicker := range okexTickers {
		if okexTicker.IsFilled() {
			var ticker = Ticker{}
			ticker.Symbol = okexTicker.Symbol
			ticker.Rate = okexTicker.Data.Last
			b.tickers[ticker.Symbol] = ticker
		}
	}


		//fmt.Println(b.tickers)

}