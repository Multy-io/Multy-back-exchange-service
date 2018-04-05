package core

import (

	"log"
	"fmt"
	"encoding/json"
	"time"
	"strconv"
	"Multy-back-exchange-service/api"
)

type PoloniexTicker struct {
	CurrencyPair  string  `json:"currencyPair"`
	Last          string `json:"last"`
}


type PoloniexManager struct {
	tickers map[string]Ticker
	poloniexApi *api.PoloniexApi
	channelsByID map[string]string
}


func (poloniexTicker PoloniexTicker) IsFilled() bool {
	return (len(poloniexTicker.CurrencyPair) > 0 && len(poloniexTicker.Last) > 0)
}


func (b *PoloniexManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.poloniexApi = &api.PoloniexApi{}
	b.channelsByID = map[string]string{"121":"USDT_BTC", "149":"USDT_ETH", "168":"BTC_STEEM", "123":"USDT_LTC","191":"USDT_BCH","173":"USDT_ETC","122":"USDT_DASH"}

	go b.poloniexApi.StartListen( func(message []byte, error error) {
		if error != nil {
			log.Println("error:", error)
			callback(TickerCollection{}, error)
		} else if message != nil {
			var unmarshaledMessage []interface{}

			err := json.Unmarshal(message, &unmarshaledMessage)
			if err != nil {
				fmt.Println(err)
				callback(TickerCollection{}, err)
			} else if len(unmarshaledMessage) > 2 {
				var poloniexTicker PoloniexTicker
				args := unmarshaledMessage[2].([]interface{})
				poloniexTicker, err = b.convertArgsToTicker(args)
				//fmt.Println(poloniexTicker)

				if error == nil && poloniexTicker.IsFilled()  {

					var ticker Ticker
					ticker.Rate = poloniexTicker.Last
					ticker.Symbol = poloniexTicker.CurrencyPair
					b.tickers[ticker.Symbol] = ticker
				}
			} else {
				fmt.Println( "error parsing Poloniex ticker:", error)
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
			//fmt.Println(tickerCollection)
			callback(tickerCollection, nil)
		}()
	}
}


func (b *PoloniexManager) convertArgsToTicker(args []interface{}) (wsticker PoloniexTicker, err error) {
	wsticker.CurrencyPair = b.channelsByID[strconv.FormatFloat(args[0].(float64), 'f', 0, 64)]
	wsticker.Last = args[1].(string)
	return
}