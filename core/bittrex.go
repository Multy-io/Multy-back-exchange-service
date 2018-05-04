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


type BittrexTicker struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Bid  float64 `json:"Bid"`
		Ask  float64 `json:"Ask"`
		Last float64 `json:"Last"`
	} `json:"result"`
	Pair currencies.CurrencyPair
}


type BittrexManager struct {
	BasicManager
	bittrexApi    *api.BittrexApi
}


func (b *BittrexManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.bittrexApi = api.NewBittrexApi()
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

				var bittrexTicker BittrexTicker
				json.Unmarshal(response.Message, &bittrexTicker)
				if bittrexTicker.Success {
					bittrexTicker.Pair = response.Pair
					var ticker Ticker
					ticker.Rate = bittrexTicker.Result.Last
					ticker.TimpeStamp = time.Now()
					ticker.Pair = bittrexTicker.Pair
					b.Lock()
					b.tickers[ticker.Pair.Symbol()] = ticker
					b.Unlock()
				}
			}

		}


	}

}

func (b *BittrexManager) listen(pairs []currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error) {
	go func() {
		for range time.Tick(5 * time.Second) {
			for _, paiar := range pairs {
				go b.bittrexApi.GetTicker(paiar, responseCh, errorCh)
			}
		}
	}()
}

func (b *BittrexManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
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