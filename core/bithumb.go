package core

import (
	"time"
	"github.com/Appscrunch/Multy-back-exchange-service/api"
	"fmt"
	"encoding/json"
	"github.com/Appscrunch/Multy-back-exchange-service/currencies"
	//"strings"
	"strconv"
	)

type BithumbCoinResult struct {
	OpeningPrice string `json:"opening_price"`
	ClosingPrice string `json:"closing_price"`
	MinPrice     string `json:"min_price"`
	MaxPrice     string `json:"max_price"`
	AveragePrice string `json:"average_price"`
	UnitsTraded  string `json:"units_traded"`
	Volume1Day   string `json:"volume_1day"`
	Volume7Day   string `json:"volume_7day"`
	BuyPrice     string `json:"buy_price"`
	SellPrice    string `json:"sell_price"`
}

type BithumbTicker struct {
	Status  string `json:"status"`
	Data map[string]*BithumbCoinResult `json:"data"`
}


type BithumbManager struct {
	BasicManager
	bithumbApi    *api.BithumbApi
}


func (b *BithumbManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.bithumbApi = api.NewBithumbApi()
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

				var bithumbTicker BithumbTicker
				json.Unmarshal(response.Message, &bithumbTicker)



				for k, v :=range  bithumbTicker.Data {

					if  v.ClosingPrice != "" {
						//fmt.Println(k,v)
						var ticker Ticker
						ticker.Rate, _ = strconv.ParseFloat(v.ClosingPrice, 64)
						ticker.Rate = ticker.Rate * 0.0009125181247
						ticker.TimpeStamp = time.Now()
						targetCurrency := currencies.NewCurrencyWithCode(k)
						referenceCurrency := currencies.Tether
						ticker.Pair = currencies.CurrencyPair{TargetCurrency:targetCurrency, ReferenceCurrency:referenceCurrency}
						//fmt.Println(ticker.Pair.Symbol())
						b.Lock()
						b.tickers[ticker.Pair.Symbol()] = ticker
						b.Unlock()
					}

				}


			}

		}


	}

}

func (b *BithumbManager) listen(pairs []currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error) {
	go func() {
		for range time.Tick(5 * time.Second) {
			go b.bithumbApi.GetTicker(currencies.CurrencyPair{}, responseCh, errorCh)
		}
	}()
}

func (b *BithumbManager) startSendingDataBack(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {
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