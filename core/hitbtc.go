package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

type HitBtcManager struct {
	BasicManager
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

func (b *HitBtcManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.hitBtcApi = &api.HitBtcApi{}

	var apiCurrenciesConfiguration = api.ApiCurrenciesConfiguration{}
	apiCurrenciesConfiguration.TargetCurrencies = exchangeConfiguration.TargetCurrencies
	apiCurrenciesConfiguration.ReferenceCurrencies = exchangeConfiguration.ReferenceCurrencies

	ch := make(chan api.Reposponse)

	go b.hitBtcApi.StartListen(apiCurrenciesConfiguration, ch)
	go b.startSendingDataBack(exchangeConfiguration, resultChan)

	for {
		select {
		case response := <-ch:

			if *response.Err != nil {
				log.Errorf("StartListen:HitBtcManager:error:", response.Err)
				//callback(nil, error)
			} else if response.Message != nil {
				//fmt.Printf("%s \n", message)
				var hitBtcTicker HitBtcTicker
				err := json.Unmarshal(*response.Message, &hitBtcTicker)
				if err == nil && hitBtcTicker.IsFilled() {
					b.add(hitBtcTicker)
				} else {
					log.Errorf("StartListen:HitBtcManager:error parsing hitBtc ticker")
				}
			}

		//default:
			//fmt.Println("no activity")
		}
	}

}

func (b *HitBtcManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	for range time.Tick(1 * time.Second) {
		//TODO: add check if data is old and don't sent it ti callback
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
				//fmt.Println(tickerCollection)
				resultChan <- Result{exchangeConfiguration.Exchange.String(), &tickerCollection, nil}
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
	b.Lock()
	b.tickers[ticker.Symbol] = ticker
	b.Unlock()
}
