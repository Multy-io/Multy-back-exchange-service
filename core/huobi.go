package core

import (
"time"
"github.com/Appscrunch/Multy-back-exchange-service/api"
"fmt"
"encoding/json"
"github.com/Appscrunch/Multy-back-exchange-service/currencies"
)

//type BittrexTicker struct {
//	Symbol string `json:"currencyPair"`
//	Last   string `json:"last"`
//}



type HuobiTicker struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		Amount  float64   `json:"amount"`
		Open    float64   `json:"open"`
		Close   float64   `json:"close"`
		High    float64   `json:"high"`
		ID      int64     `json:"id"`
		Count   int       `json:"count"`
		Low     float64   `json:"low"`
		Version int64     `json:"version"`
		Ask     []float64 `json:"ask"`
		Vol     float64   `json:"vol"`
		Bid     []float64 `json:"bid"`
	} `json:"tick"`
}


type HuobiManager struct {
	BasicManager
	huobyApi    *api.HuobiApi
}


func (b *HuobiManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.huobyApi = api.NewHuobiApi()
	//b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	//b.setchannelids()

	pairs := exchangeConfiguration.Pairs

	//ch := make(chan api.Reposponse)

	responseCh := make(chan api.RestApiReposponse)
	errorCh := make(chan error)

	b.listen(pairs, responseCh, errorCh)
	b.startSendingDataBack(exchangeConfiguration, resultChan)

	for {
		select {
		case err := <-errorCh:
			fmt.Println(err)
		case response := <-responseCh:
			//fmt.Printf("%s %@ %@ \n", response.Message, response.Pair.TargetCurrency.CurrencyCode(), response.Pair.ReferenceCurrency.CurrencyCode())
			if response.Message != nil {

				var huobiTicker HuobiTicker
				json.Unmarshal(response.Message, &huobiTicker)
				if huobiTicker.Status == "ok" {
					var ticker Ticker
					ticker.Rate = huobiTicker.Tick.Bid[0]
					ticker.TimpeStamp = time.Now()
					ticker.Pair = response.Pair
					b.Lock()
					b.tickers[ticker.Pair.Symbol()] = ticker
					b.Unlock()
				}
			}

		}


	}

}

func (b *HuobiManager) listen(pairs []currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error) {
	go func() {
		for range time.Tick(5 * time.Second) {
			for _, paiar := range pairs {
				go b.huobyApi.GetTicker(paiar, responseCh, errorCh)
			}
		}
	}()
}

func (b *HuobiManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
	go func() {
		for range time.Tick(1 * time.Second) {
			func() {
				tickers := []Ticker{}

				b.Lock()
				for _, value := range b.tickers {
					if value.TimpeStamp.After(time.Now().Add(-maxTickerAge * time.Second)) {
						tickers = append(tickers, value)
					}
				}
				b.Unlock()

				var tickerCollection = TickerCollection{}
				tickerCollection.TimpeStamp = time.Now()
				tickerCollection.Tickers = tickers
				//fmt.Println(tickerCollection)
				if len(tickerCollection.Tickers) > 0 {
					resultChan <- Result{exchangeConfiguration.Exchange.String(), &tickerCollection, nil}
				}
			}()
		}
	}()
}