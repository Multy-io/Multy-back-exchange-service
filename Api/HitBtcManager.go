package Api

import (

	"log"
	"fmt"
	"encoding/json"
	"time"
)



type HitBtcManager struct {

}

type HitBtcTicker struct {
	Params  struct {
		Rate        string    `json:"last"`
		Symbol      string    `json:"symbol"`
	} `json:"params"`
}

func (hitBtcTicker HitBtcTicker) IsFilled() bool {
	return (len(hitBtcTicker.Params.Symbol) > 0 && len(hitBtcTicker.Params.Rate) > 0)
}

var tickers = make(map[string]Ticker)

var hitBtcApi = HitBtcApi{}

func (b HitBtcManager) StartListen(callback func(tickerCollection TickerCollection, error error)) {

	go hitBtcApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("error:", error)
			//callback(nil, error)
		} else if message != nil {
			//fmt.Printf("%s \n", message)
			var hitBtcTicker HitBtcTicker
			error := json.Unmarshal(message, &hitBtcTicker)
			if error == nil && hitBtcTicker.IsFilled()  {
				add(hitBtcTicker)
			} else {
				fmt.Println( "error parsing hitBtc ticker:", error)
			}
		}
	})

	for range time.Tick(1 * time.Second) {
		//TODO: add check if data is old and don't sent it ti callback
		func() {
			values := []Ticker{}
			for _, value := range tickers {
				values = append(values, value)
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = values
			callback(tickerCollection, nil)
		}()
	}
}

func add(hitBtcTicker HitBtcTicker) {
	var ticker = Ticker{}
	ticker.Rate = hitBtcTicker.Params.Rate
	ticker.Symbol = hitBtcTicker.Params.Symbol
	tickers[ticker.Symbol] = ticker
}