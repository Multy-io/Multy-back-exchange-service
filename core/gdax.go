package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type GdaxManager struct {
	BasicManager
	gdaxApi *api.GdaxApi
}

type GdaxTicker struct {
	BestAsk   string `json:"best_ask"`
	BestBid   string `json:"best_bid"`
	High24h   string `json:"high_24h"`
	Low24h    string `json:"low_24h"`
	Open24h   string `json:"open_24h"`
	Rate      string `json:"price"`
	Symbol    string `json:"product_id"`
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

func (b *GdaxManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.gdaxApi = &api.GdaxApi{}

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	ch := make(chan api.Reposponse)

	go b.gdaxApi.StartListen(apiCurrenciesConfiguration, ch)
	go b.startSendingDataBack(exchangeConfiguration, resultChan)

	for {
		select {
		case response := <-ch:

			if *response.Err != nil {
				log.Errorf("StartListen:GdaxManager:error:", response.Err)
				//callback(nil, error)
			} else if response.Message != nil {
				//fmt.Printf("%s \n", message)
				var gdaxTicker GdaxTicker
				err := json.Unmarshal(*response.Message, &gdaxTicker)
				if err == nil && gdaxTicker.IsFilled() {
					b.add(gdaxTicker)
					//fmt.Println(gdaxTicker)
				} else {
					log.Errorf("StartListen:error parsing hitBtc ticker:")
				}
			}

		//default:
			//fmt.Println("no activity")
		}
	}

}

func (b *GdaxManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	for range time.Tick(1 * time.Second) {
		func() {
			values := []Ticker{}
			b.Lock()
			tickers := b.tickers
			b.Unlock()

			for _, value := range tickers {
				if value.TimpeStamp.After(time.Now().Add(-maxTickerAge * time.Second)) {
					values = append(values, value)
				}
			}

			var tickerCollection = TickerCollection{}
			tickerCollection.TimpeStamp = time.Now()
			tickerCollection.Tickers = values
			if len(tickerCollection.Tickers) > 0 {
				resultChan <- Result{exchangeConfiguration.Exchange.String(), &tickerCollection, nil}
			}
		}()
	}
}

func (b *GdaxManager) add(gdaxTicker GdaxTicker) {
	var ticker = Ticker{}
	ticker.Rate = gdaxTicker.Rate
	ticker.Symbol = gdaxTicker.Symbol

	targetCurrency, referenceCurrency := gdaxTicker.getCurriences()
	ticker.TargetCurrency = targetCurrency
	ticker.ReferenceCurrency = referenceCurrency
	ticker.TimpeStamp = time.Now()
	b.Lock()
	b.tickers[ticker.Symbol] = ticker
	b.Unlock()
}
