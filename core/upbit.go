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

type UpbitTicker struct {
	Code             string  `json:"code"`
	TradeDate        string  `json:"tradeDate"`
	TradeTime        string  `json:"tradeTime"`
	TradeDateKst     string  `json:"tradeDateKst"`
	TradeTimeKst     string  `json:"tradeTimeKst"`
	TradeTimestamp   int64   `json:"tradeTimestamp"`
	TradePrice       float64 `json:"tradePrice"`
	TradeVolume      float64 `json:"tradeVolume"`
	PrevClosingPrice float64 `json:"prevClosingPrice"`
	Change           string  `json:"change"`
	ChangePrice      float64 `json:"changePrice"`
	AskBid           string  `json:"askBid"`
	SequentialID     int64   `json:"sequentialId"`
	Timestamp        int64   `json:"timestamp"`
}


type UpbitManager struct {
	BasicManager
	upbitApi    *api.UpbitApi
}


func (b *UpbitManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.upbitApi = api.NewUpbitApi()
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

				var upbitTickers []UpbitTicker
				json.Unmarshal(response.Message, &upbitTickers)
				if len(upbitTickers) > 0 {
					upbitTicker := upbitTickers[0]
						var ticker Ticker
						ticker.Rate = upbitTicker.TradePrice
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

func (b *UpbitManager) listen(pairs []currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error) {
	go func() {
		for range time.Tick(5 * time.Second) {
			for _, paiar := range pairs {
				go b.upbitApi.GetTicker(paiar, responseCh, errorCh)
			}
		}
	}()
}

func (b *UpbitManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
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