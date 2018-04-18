package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type HitBtcManager struct {
	tickers   map[string]Ticker
	hitBtcApi *api.HitBtcApi
}

type HitBtcTicker struct {
	Params struct {
		Rate   string `json:"last"`
		Symbol string `json:"symbol"`
	} `json:"params"`
}

func (b *HitBtcTicker) getCurriences() (currencies.Currency, currencies.Currency) {

	if len(b.Params.Symbol) > 0 {
		var symbol = b.Params.Symbol
		var damagedSymbol = TrimLeftChars(symbol, 1)
		for _, referenceCurrency := range currencies.DefaultReferenceCurrencies {
			//fmt.Println(damagedSymbol, referenceCurrency.CurrencyCode())

			referenceCurrencyCode := referenceCurrency.CurrencyCode()

			if referenceCurrencyCode == "USDT" {
				referenceCurrencyCode = "USD"
			}

			if strings.Contains(damagedSymbol, referenceCurrencyCode) {

				//fmt.Println("2",symbol, referenceCurrency.CurrencyCode())
				targetCurrencyString := strings.TrimSuffix(symbol, referenceCurrencyCode)
				//fmt.Println(targetCurrencyString)
				var targetCurrency = currencies.NewCurrencyWithCode(targetCurrencyString)
				return targetCurrency, referenceCurrency
			}
		}

	}
	return currencies.NotAplicable, currencies.NotAplicable
}

func (hitBtcTicker HitBtcTicker) IsFilled() bool {
	return (len(hitBtcTicker.Params.Symbol) > 0 && len(hitBtcTicker.Params.Rate) > 0)
}

func (b *HitBtcManager) StartListen(exchangeConfiguration ExchangeConfiguration, callback func(tickerCollection TickerCollection, error error)) {

	b.tickers = make(map[string]Ticker)
	b.hitBtcApi = &api.HitBtcApi{}

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	go b.hitBtcApi.StartListen(apiCurrenciesConfiguration, func(message []byte, error error) {
		if error != nil {
			log.Println("error:", error)
			//callback(nil, error)
		} else if message != nil {
			//fmt.Printf("%s \n", message)
			var hitBtcTicker HitBtcTicker
			error := json.Unmarshal(message, &hitBtcTicker)
			if error == nil && hitBtcTicker.IsFilled() {
				b.add(hitBtcTicker)
			} else {
				fmt.Println("error parsing hitBtc ticker:", error)
			}
		}
	})

	for range time.Tick(1 * time.Second) {
		//TODO: add check if data is old and don't sent it ti callback
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
			if len(tickerCollection.Tickers) > 0 {
				//fmt.Println(tickerCollection)
				callback(tickerCollection, nil)
			}
		}()
	}
}

func (b *HitBtcManager) add(hitBtcTicker HitBtcTicker) {
	var ticker = Ticker{}
	ticker.Rate = hitBtcTicker.Params.Rate
	ticker.Symbol = hitBtcTicker.Params.Symbol
	//fmt.Println(hitBtcTicker)
	targetCurrency, referenceCurrency := hitBtcTicker.getCurriences()
	ticker.TargetCurrency = targetCurrency
	ticker.ReferenceCurrency = referenceCurrency
	ticker.TimpeStamp = time.Now()

	b.tickers[ticker.Symbol] = ticker
}
