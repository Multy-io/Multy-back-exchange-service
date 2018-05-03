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

type KRWTicker struct {
	KRWUSD struct {
		Val float64 `json:"val"`
	} `json:"KRW_USD"`
}



type BithumbManager struct {
	BasicManager
	bithumbApi    *api.BithumbApi
	fiatApi    *api.FiatApi
	currentKrwRate float64
}


func (b *BithumbManager) StartListen(exchangeConfiguration ExchangeConfiguration, resultChan chan Result) {

	b.tickers = make(map[string]Ticker)
	b.bithumbApi = api.NewBithumbApi()
	b.fiatApi = api.NewFiatApi()
	//b.symbolsToParse = b.composeSybolsToParse(exchangeConfiguration)
	//b.setchannelids()


	krwPair := currencies.CurrencyPair{currencies.SouthKoreanWon, currencies.Tether}

	//ch := make(chan api.Reposponse)

	responseCh := make(chan api.RestApiReposponse)
	errorCh := make(chan error)

	b.listenFiat(krwPair, responseCh, errorCh, 60)
	b.listen(currencies.CurrencyPair{}, responseCh, errorCh, 5)
	b.startSendingDataBack(exchangeConfiguration, resultChan)

	for {
		select {
		case err := <-errorCh:
			fmt.Println(err)
		case response := <-responseCh:
			//fmt.Printf("%s %@ %@ \n", response.Message, response.Pair.TargetCurrency.CurrencyCode(), response.Pair.ReferenceCurrency.CurrencyCode())
			if response.Message != nil {

				if response.Pair.TargetCurrency == currencies.SouthKoreanWon {
					var krwTicker KRWTicker
					json.Unmarshal(response.Message, &krwTicker)
					b.currentKrwRate = krwTicker.KRWUSD.Val
				} else if b.currentKrwRate > 0 {
					var bithumbTicker BithumbTicker
					json.Unmarshal(response.Message, &bithumbTicker)

					for k, v :=range  bithumbTicker.Data {

						if  v.ClosingPrice != "" {
							//fmt.Println(k,v)
							var ticker Ticker
							ticker.Rate, _ = strconv.ParseFloat(v.ClosingPrice, 64)
							ticker.Rate = ticker.Rate * b.currentKrwRate
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

}

func (b *BithumbManager) listen(pair currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error, refreshInterval time.Duration) {
	go func() {
		for range time.Tick(refreshInterval * time.Second) {
			go b.bithumbApi.GetTicker(pair, responseCh, errorCh)
		}
	}()
}

func (b *BithumbManager) listenFiat(pair currencies.CurrencyPair, responseCh chan api.RestApiReposponse, errorCh chan error, refreshInterval time.Duration) {
	go func() {
		for range time.Tick(refreshInterval * time.Second) {
			go b.fiatApi.GetTicker(pair, responseCh, errorCh)
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