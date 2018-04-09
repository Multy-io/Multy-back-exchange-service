package core

import (

	"log"
	"fmt"
	"encoding/json"
	"time"
	"strconv"
	"Multy-back-exchange-service/api"
	//"sync"
	"strings"
	"Multy-back-exchange-service/currencies"
)

const (
	TICKER     = "1002" /* Ticker Channel Id */
	SUBSBUFFER = 24     /* Subscriptions Buffer */
)

type PoloniexTicker struct {
	Symbol  string  `json:"currencyPair"`
	Last          string `json:"last"`
}

func (b *PoloniexTicker) getCurriences() (currencies.Currency, currencies.Currency) {

	if len(b.Symbol) > 0 {
		var symbol = b.Symbol
		var currencyCodes = strings.Split(symbol, "_")
		if len(currencyCodes) > 1 {
			return currencies.NewCurrencyWithCode(currencyCodes[0]), currencies.NewCurrencyWithCode(currencyCodes[1])
		}
	}
	return currencies.NotAplicable, currencies.NotAplicable
}


type PoloniexManager struct {
	tickers map[string]Ticker
	poloniexApi *api.PoloniexApi
	channelsByID map[string]string
	channelsByName map[string]string
	marketChannels []string
	symbolsToParse map[string]bool
}


func (poloniexTicker PoloniexTicker) IsFilled() bool {
	return (len(poloniexTicker.Symbol) > 0 && len(poloniexTicker.Last) > 0)
}





func (b *PoloniexManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.poloniexApi = api.NewPoloniexApi()
	b.channelsByID = make(map[string]string)
	b.channelsByName = make(map[string]string)
	b.marketChannels = []string{}
	b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	b.setchannelids()

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
				//fmt.Println(poloniexTicker.CurrencyPair)

				if error == nil && poloniexTicker.IsFilled() && b.symbolsToParse[poloniexTicker.Symbol]  {

					var ticker Ticker
					ticker.Rate = poloniexTicker.Last
					ticker.Symbol = poloniexTicker.Symbol
					ticker.TimpeStamp = time.Now()
					targetCurrency, referenceCurrency  := poloniexTicker.getCurriences()
					ticker.TargetCurrency = targetCurrency
					ticker.ReferenceCurrency = referenceCurrency

					b.tickers[ticker.Symbol] = ticker
				}
			} else {
				fmt.Println( "error parsing Poloniex ticker:", error)
			}
		}
	})

	for range time.Tick(1 * time.Second) {
		func() {
			values := []Ticker{}
			for _, value := range b.tickers {
				if value.TimpeStamp.After(time.Now().Add(-maxTickerAge * time.Second)) {
					values = append(values, value)
				}
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = values
			//fmt.Println(tickerCollection)
			if len(tickerCollection.Tickers) > 0 {
				callback(tickerCollection, nil)
			}
		}()
	}
}


func (b *PoloniexManager) convertArgsToTicker(args []interface{}) (wsticker PoloniexTicker, err error) {

	if len(b.channelsByID) < 1 {
		b.setchannelids()
	}

	wsticker.Symbol = b.channelsByID[strconv.FormatFloat(args[0].(float64), 'f', 0, 64)]
	wsticker.Last = args[1].(string)
	return
}



func (b *PoloniexManager) setchannelids() (err error) {

	resp, err := b.poloniexApi.PubReturnTickers()
	if err != nil {
		return err
	}

	for k, v := range resp {
		chid := strconv.Itoa(v.ID)
		b.channelsByName[k] = chid
		b.channelsByID[chid] = k
		b.marketChannels = append(b.marketChannels, chid)
	}

	b.channelsByName["TICKER"] = TICKER
	b.channelsByID[TICKER] = "TICKER"
	return
}

func (b *PoloniexManager)  composeSybolsToParse(exchangeConfiguration ExchangeConfiguration) map[string]bool {
	var symbolsToParse = map[string]bool{}
	for _, targetCurrency := range exchangeConfiguration.TargetCurrencies {
		for _, referenceCurrency := range exchangeConfiguration.ReferenceCurrencies {

			if referenceCurrency == "USD" {
				referenceCurrency = "USDT"
			}

			symbol := referenceCurrency  + "_" + targetCurrency
			symbolsToParse[symbol] = true
		}
	}
	return symbolsToParse

}