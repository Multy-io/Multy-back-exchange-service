package core

import (
	"log"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"Multy-back-exchange-service/api"
)

type BitfinexManager struct {
	bitfinexTickers map[int]BitfinexTicker
	api *api.BitfinexApi
}

type BitfinexTicker struct {
	ChanID  int    `json:"chanId"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Pair    string `json:"pair"`
	Symbol  string `json:"symbol"`
	Rate	string
}

	func (bitfinexTicker BitfinexTicker) IsFilled() bool {
	return (len(bitfinexTicker.Symbol) > 0 && len(bitfinexTicker.Rate) > 0)
}


func (b *BitfinexManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {
	b.bitfinexTickers = make(map[int]BitfinexTicker)
	b.api = api.NewBitfinexApi()

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	go b.api.StartListen(apiCurrenciesConfiguration, func(message []byte, error error) {
		//fmt.Println(0)
		if error != nil {
			log.Println("error:", error)
			callback(TickerCollection{}, error)
		} else if message != nil {
			//fmt.Printf("%s \n", message)
			//fmt.Println(1)
			b.addMessage(message)
			//fmt.
			} else {
				fmt.Println( "error parsing Bitfinex ticker:", error)
			}
		})

	for range time.Tick(1 * time.Second) {
		//TODO: add check if data is old and don't sent it to callback
		func() {
			tickers := []Ticker{}
			for _, value := range b.bitfinexTickers {
				var ticker = Ticker{}
				ticker.Rate = value.Rate
				ticker.Symbol = value.Symbol
				tickers = append(tickers, ticker)
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = tickers
			//fmt.Println(tickerCollection)
			callback(tickerCollection, nil)
		}()
	}
}


func (b *BitfinexManager) addMessage (message []byte) {

	var bitfinexTicker BitfinexTicker
	json.Unmarshal(message, &bitfinexTicker)

	if bitfinexTicker.ChanID > 0 {
		//fmt.Println(bitfinexTicker)
		b.bitfinexTickers[bitfinexTicker.ChanID] = bitfinexTicker
	} else {

		var unmarshaledTickerMessage []interface{}
		json.Unmarshal(message, &unmarshaledTickerMessage)


		if len(unmarshaledTickerMessage) > 1 {
			var chanId = int(unmarshaledTickerMessage[0].(float64))
			//var unmarshaledTicker []interface{}
			if v, ok := unmarshaledTickerMessage[1].([]interface{}); ok {
				var sub = b.bitfinexTickers[chanId]
				sub.Rate = strconv.FormatFloat(v[0].(float64), 'f', 8, 64)
				b.bitfinexTickers[chanId] = sub
				//fmt.Println(b.bitfinexTickers)
			}
		}
	}
}

//func (b PoloniexManager) convertArgsToTicker(args []interface{}) (wsticker PoloniexTicker, err error) {
//	wsticker.CurrencyPair = b.channelsByID[strconv.FormatFloat(args[0].(float64), 'f', 0, 64)]
//	wsticker.Last = args[1].(string)
//	return
//}